package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/inaryzen/priotasks/common"
	"github.com/inaryzen/priotasks/consts"
	"github.com/inaryzen/priotasks/csv"
	"github.com/inaryzen/priotasks/db"
	"github.com/inaryzen/priotasks/handlers"
	"github.com/inaryzen/priotasks/services"
)

//go:embed assets/*
var assets embed.FS

func main() {
	printVersion()

	common.InitConfig()

	var port string = fmt.Sprintf(":%v", common.Conf.ServerPort)
	server := &http.Server{Addr: port}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	newDb := db.NewDbSQLite()
	db.SetDB(newDb)
	db.DB().Init("")

	defer db.DB().Close()

	if common.Conf.DumpOnStartup {
		if common.IsDebug() {
			log.Println("dump db on startup...")
		}
		err := csv.Dump()
		if err != nil {
			log.Printf("failed to dump db: %v", err)
			return
		}
	}

	if common.Conf.LoadDumpOnStartup != "" {
		csv.Load(common.Conf.LoadDumpOnStartup)
	}

	if common.Conf.AutomaticDump {
		log.Println("auto-dump enabled...")
		s := csv.NewDumpScheduler(2400) // 40min
		defer s.Release()
	}

	services.Init()
	configureServerMux()
	go startServer(server)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println()
	log.Println("shutting down the server...")

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}

func printVersion() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("Unable to determine version information.")
		return
	}
	fmt.Printf("version:%s\n", buildInfo.Main.Version)
	fmt.Printf("sum:%s\n", buildInfo.Main.Sum)
}

func configureServerMux() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, consts.URL_TASKS, http.StatusFound) // 302
	})

	http.HandleFunc("GET "+consts.URL_TASKS, handlers.GetTasks)
	http.HandleFunc("POST "+consts.URL_TASKS, handlers.PostTaskHandler)
	http.HandleFunc("PUT "+consts.URL_TASKS, handlers.PutTaskHandler)
	http.HandleFunc("POST /tasks/{id}/toggle-completed", handlers.PostTaskToggleCompleted)
	http.HandleFunc("DELETE "+consts.URL_TASKS_ID, handlers.DeleteTasksId)
	http.HandleFunc("POST /filter/{name}", handlers.PostFilterName)
	http.HandleFunc("DELETE /filter/tag/{name}", handlers.DeleteTagName)
	http.HandleFunc("POST /prepared-query/{name}", handlers.PostPreparedQuery)
	http.HandleFunc("POST "+consts.URL_TOGGLE_SORT_TABLE, handlers.PostToggleSortTable)
	http.HandleFunc("GET /view/task/{id}", handlers.GetViewTaskByIdHandler)
	http.HandleFunc("GET /view/new-task", handlers.GetViewEmptyTask)
	http.HandleFunc("POST /tags", handlers.PostTagsHandler)
	http.HandleFunc("DELETE /tags/{name}", handlers.DeleteTagHandler)
	http.HandleFunc("POST /tasks/reduce-priority", handlers.PostReducePriorityHandler)
	http.Handle("/assets/", http.FileServer(http.FS(assets)))
}

func startServer(s *http.Server) {
	log.Println("starting the server...")
	log.Printf("http://localhost%s \n", s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
