package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
)

// Slow sink implementation to test metric data in debugging
type SlowSink struct{}

func newSlowSink() *SlowSink {
	return &SlowSink{}
}

func (l *SlowSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	result := sinksdk.ResponsesBuilder()
	for d := range datumStreamCh {
		if d.EventTime().Nanosecond()%10 == 0 {
			sleepDuration := time.Duration(rand.Intn(5)) * time.Second
			time.Sleep(sleepDuration)
		}
		id := d.ID()
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
