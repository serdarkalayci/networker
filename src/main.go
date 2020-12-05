package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/serdarkalayci/networker/handlers"

	"github.com/gorilla/mux"
	"github.com/serdarkalayci/networker/env"
)

var bindAddress = env.String("BASE_URL", ":5600")

func main() {
	l := log.New(os.Stdout, "networker ", log.LstdFlags)

	// create the handlers
	apiContext := handlers.NewAPIContext()

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	// handlers for API
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/", apiContext.Index)
	getR.HandleFunc("/health/live", apiContext.Live)
	getR.HandleFunc("/health/ready", apiContext.Ready)
	getR.HandleFunc("/log", apiContext.Logs)
	getR.HandleFunc("/reset", apiContext.Reset)

	// create a new server
	s := http.Server{
		Addr:         bindAddress,       // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Println(fmt.Sprintf("Starting server on %s", bindAddress))

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	go checkAddress(apiContext)

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
