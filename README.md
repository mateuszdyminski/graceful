#### Graceful shutdown

Graceful shutdown allows you to shutdown the http.Server without interrupting any active connections. Shutdown works by first closing all open listeners, then closing all idle connections, and then waiting indefinitely for connections to return to idle and then shut down.

#### Graceful shutdown in Golang

In following repository you can see the Golang build-in Graceful Shutdown:
 
```go
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
	log.Println("shutting down server. waiting to drain the ongoing requests...")

	// shut down gracefully, but wait no longer than the Timeout value.
	ctx, _ := context.WithTimeout(context.Background(), Timeout*time.Second)

	// shutdown the http server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("error while shutdown http server: %v\n", err)
	}

	log.Println("server gracefully stopped")
}
```

##### Requirements

Golang in version 1.8 or higher

##### Run 

```
go run graceful.go
```

##### Verification

Open browser and go to: [http://localhost:8080/heavy](https://localhost:8080/heavy)

In the meantime you could type Ctrl+c in the terminal and see that the server returns the response and shutdown the server gracefully without closing current active connection.
