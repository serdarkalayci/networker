package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/serdarkalayci/networker/data"
	"github.com/serdarkalayci/networker/env"
	"github.com/serdarkalayci/networker/handlers"
)

func checkAddress(apiContext *handlers.APIContext) {
	l := log.New(os.Stdout, "networker ", log.LstdFlags)
	targetAddress := env.String("TARGET_ADDR", "http://google.com")
	sleepDuration := env.Int("SLEEP_DURATION", 2)

	t1 := env.Int("T1", 2)
	t2 := env.Int("T2", 5)
	t3 := env.Int("T3", 10)
	ot := env.Int("OT", 10)
	for {
		timeStart := time.Now()
		resp, err := http.Get(targetAddress)
		timeEnd := time.Now()
		duration := timeEnd.Sub(timeStart).Milliseconds()
		totalTime := apiContext.LogRecord.AverageDuration * float64(apiContext.LogRecord.TotalCalls)
		apiContext.LogRecord.AverageDuration = (totalTime + float64(duration)) / float64(apiContext.LogRecord.TotalCalls+1)
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
		time.Sleep(time.Duration(sleepDuration) * time.Second)
	}

}
