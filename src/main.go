package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/serdarkalayci/networker/data"
	"github.com/serdarkalayci/networker/handlers"

	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
)

var bindAddress = *env.String("BASE_URL", true, ":5500", "Bind address for the server")

var t1 = *env.Int("T1", false, 2, "The first treshold to be counted")
var t2 = *env.Int("T2", false, 5, "The second treshold to be counted")
var t3 = *env.Int("T3", false, 10, "The third treshold to be counted")
var ot = *env.Int("OT", false, 10, "The treshold to be logged specifically")

func main() {
	l := log.New(os.Stdout, "networker ", log.LstdFlags)

	if bindAddress == "" {
		bindAddress = ":5500"
	}
	if t1 == 0 {
		t1 = 2
		t2 = 5
		t3 = 10
		ot = 10
	}

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

	go func() {
		for {
			timeStart := time.Now()
			resp, err := http.Get("http://google.com")
			timeEnd := time.Now()
			duration := timeEnd.Sub(timeStart).Milliseconds()
			totalTime := apiContext.LogRecord.AverageDuration * float64(apiContext.LogRecord.TotalCalls)
			apiContext.LogRecord.AverageDuration = totalTime / float64(apiContext.LogRecord.TotalCalls+1)
			apiContext.LogRecord.TotalCalls++
			if err != nil {
				apiContext.LogRecord.ErrorCount++
				errorRecord := data.ErrorRecord{
					Timestamp: timeStart,
					Duration:  duration,
					Error:     err.Error(),
				}
				apiContext.LogRecord.ErrorRecords = append(apiContext.LogRecord.ErrorRecords, errorRecord)
			} else {
				resp.Body.Close()
				if duration > int64(t1*1000) {
					apiContext.LogRecord.T1Count++
				}
				if duration > int64(t2*1000) {
					apiContext.LogRecord.T2Count++
				}
				if duration > int64(t3*1000) {
					apiContext.LogRecord.T3Count++
				}
				if duration > int64(ot*1000) {
					logRecord := data.LogRecord{
						Timestamp: timeStart,
						Duration:  duration,
					}
					apiContext.LogRecord.OutlierRecords = append(apiContext.LogRecord.OutlierRecords, logRecord)
				}
			}
			l.Println(fmt.Sprintf("Duration %d milliseconds", duration))
			time.Sleep(2 * time.Second)
		}
	}()

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
