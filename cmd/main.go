package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Maxim-Ba/cv-backend/config"
	"github.com/Maxim-Ba/cv-backend/internal/dbconn"
)

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	cfg := config.GetConfig()
	db, err := dbconn.New(*cfg)

	if err != nil {
		panic(err.Error())
	}
	router, err := initApplication(ctx,db, cfg)
	if err != nil {
		panic(err.Error())
	}
	var wg sync.WaitGroup
	server := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: router.R,
		
	}
	go func() {
		wg.Add(1)
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {

			os.Exit(1)
		}
	}()

	select {
		case <-exit:
		case <-ctx.Done():
	}
//TODO shutdown actions
	wg.Wait()
}
