package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/inaryzen/priotasks/common"
	"github.com/inaryzen/priotasks/consts"
	"github.com/inaryzen/priotasks/csv"
	"github.com/inaryzen/priotasks/db"
	"github.com/inaryzen/priotasks/handlers"
)

//go:embed assets/*
var assets embed.FS

func main() {
	common.InitConfig()

	var port string = fmt.Sprintf(":%v", common.Conf.ServerPort)
	server := &http.Server{Addr: port}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	db.DbInit()
	defer db.DbClose()

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
		fmt.Println("auto-dump enabled...")
		s := csv.NewDumpScheduler(2400) // 40min
		defer s.Release()
	}

	configureServerMux()
	go startServer(server)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("\nshutting down the server...")

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}

func configureServerMux() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, consts.URL_TASKS, http.StatusFound) // 302
	})

	http.HandleFunc("GET "+consts.URL_TASKS, handlers.GetTasks)
	http.HandleFunc("POST "+consts.URL_TASKS, handlers.PostTaskHandler)
	http.HandleFunc("PUT "+consts.URL_TASKS, handlers.PutTaskHandler)
	http.HandleFunc("POST /tasks/{id}/toggle-completed", handlers.PostTaskToggleCompleted)

	http.HandleFunc("POST /toggle-completed-filter", handlers.PostToggleCompletedFilter)
	http.HandleFunc("POST "+consts.URL_TOGGLE_SORT_TABLE, handlers.PostToggleSortTable)

	http.HandleFunc("GET /view/task/{id}", handlers.GetViewTaskByIdHandler)
	http.HandleFunc("GET /view/new-task", handlers.GetViewEmptyTask)

	http.Handle("/assets/", http.FileServer(http.FS(assets)))
}

func startServer(s *http.Server) {
	fmt.Println("starting the server...")
	fmt.Printf("http://localhost%s \n", s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
