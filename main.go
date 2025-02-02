package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/inaryzen/prio_cards/common"
	"github.com/inaryzen/prio_cards/consts"
	"github.com/inaryzen/prio_cards/csv"
	"github.com/inaryzen/prio_cards/db"
	"github.com/inaryzen/prio_cards/handlers"
)

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

	fmt.Println("shutting down the server...")

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}

func configureServerMux() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/cards", http.StatusFound) // 302
	})

	http.HandleFunc("GET /cards", handlers.GetCards)
	http.HandleFunc("POST /cards", handlers.PostCardHandler)
	http.HandleFunc("PUT /cards", handlers.PutCardHandler)
	http.HandleFunc("POST /cards/{id}/toggle-completed", handlers.PostCardToggleCompleted)

	http.HandleFunc("POST /toggle-completed-filter", handlers.PostToggleCompletedFilter)
	http.HandleFunc("POST "+consts.URL_TOGGLE_SORT_TABLE, handlers.PostToggleSortTable)

	http.HandleFunc("GET /view/card/{id}", handlers.GetViewCardByIdHandler)
	http.HandleFunc("GET /view/new-card", handlers.GetViewEmptyCard)

	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
}

func startServer(s *http.Server) {
	fmt.Println("starting the server...")
	fmt.Printf("http://localhost%s \n", s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
