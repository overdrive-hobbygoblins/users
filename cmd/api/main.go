package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("api: couldn't load env config: %v", err)
	}

	api := http.Server{
		Addr:         cfg.Web.DebugHost,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}
	apiErrors := make(chan error, 1)

	go func() {
		log.Printf("api: started on %s", api.Addr)
		apiErrors <- api.ListenAndServe()
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-apiErrors:
		log.Printf("api: error occurred on starting http-listener: %v", err)
	case <-osSignals:
		log.Print("api: starting shutdown")

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			log.Panicf("api: graceful shutdown failed: %v", err)
			_ = api.Close()
		}
	}

	log.Print("api: shutdown")
}
