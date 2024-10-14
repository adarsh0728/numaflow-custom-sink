package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
)

// Slow sink implementation to test metric data in debugging
type SlowSink struct {
	minSleepDuration int
	maxSleepDuration int
	startTime        time.Time
}

func newSlowSink() *SlowSink {

	minSleepDurationStr, ok := os.LookupEnv("MIN_SLEEP_TIME")
	if !ok {
		minSleepDurationStr = "10"
	}
	maxSleepDurationStr, ok := os.LookupEnv("MAX_SLEEP_TIME")
	if !ok {
		maxSleepDurationStr = "20"
	}
	minSleepDuration, _ := strconv.Atoi(minSleepDurationStr)
	maxSleepDuration, _ := strconv.Atoi(maxSleepDurationStr)
	return &SlowSink{
		minSleepDuration: minSleepDuration,
		maxSleepDuration: maxSleepDuration,
		startTime:        time.Now(),
	}
}

func (l *SlowSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	result := sinksdk.ResponsesBuilder()
	min := l.minSleepDuration
	max := l.maxSleepDuration

	fmt.Println("min sleep time: ", min)
	fmt.Println("max sleep time: ", max)

	for d := range datumStreamCh {
		id := d.ID()
		// first 5 minutes(configurable), sink should work normally. B/w 5 and 8 minutes based on event times of mssgs,
		// if event time is multiple of 3 (configurable) we are introducing sleep to mimic slow sink
		// sleep duration is in a range based on set env variables (min and max time)
		if d.EventTime().Nanosecond()%3 == 0 && time.Since(l.startTime) >= 5*time.Minute && time.Since(l.startTime) <= 8*time.Minute {
			randomNumber := rand.Intn(max-min+1) + min
			sleepDuration := time.Duration(randomNumber) * time.Second
			fmt.Println("id: ", id, "sleep time: ", sleepDuration, "event time: ", d.EventTime().Nanosecond())
			time.Sleep(sleepDuration)
		} else {
			fmt.Println("id: ", id, "event time: ", d.EventTime().Nanosecond())
		}
		result = result.Append(sinksdk.ResponseOK(id))
	}
	return result
}

func main() {
	slow_sink := newSlowSink()
	err := sinksdk.NewServer(slow_sink).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start sink function server: ", err)
	}
}
