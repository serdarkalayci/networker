package data

import "time"

type ConsolidatedLog struct {
	TotalCalls      int
	AverageDuration float64
	ErrorCount      int
	T1Count         int
	T2Count         int
	T3Count         int
	OutlierRecords  []LogRecord
	ErrorRecords    []ErrorRecord
}

type LogRecord struct {
	Timestamp time.Time
	Duration  int64
}

type ErrorRecord struct {
	Timestamp time.Time
	Duration  int64
	Error     string
}
