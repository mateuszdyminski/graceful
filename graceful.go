package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Timeout is the duration to allow outstanding requests to survive before forcefully terminating them.
const Timeout = 10

// START OMIT
func main() {
	// subscribe to SIGINT signals
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)

	// configure http handlers
	mux := http.NewServeMux()
	mux.Handle("/heavy", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("doing some heavy business logic!")
		time.Sleep((Timeout - 1) * time.Second)
		log.Printf("done!")
		io.WriteString(w, "done!")
	}))

	// create and start http server in new goroutine
	srv := &http.Server{Addr: ":8080", Handler: mux}
	go func() {
		// we can't use log.Fatal here!
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("http server stoped: %s\n", err)
		}
	}()
	// blocks the execution until os.Interrupt or os.Kill signal appears
	<-quit
	// END OMIT
	// START2 OMIT
	log.Println("shutting down server. waiting to drain the ongoing requests...")

	// shut down gracefully, but wait no longer than the Timeout value.
	ctx, _ := context.WithTimeout(context.Background(), Timeout*time.Second)

	// shutdown the http server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("error while shutdown http server: %v\n", err)
	}

	log.Println("server gracefully stopped")
}
// END2 OMIT
