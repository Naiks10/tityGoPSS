package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type service struct {
	channels     []*channel
	processed    atomic.Value
	que          *queue
	finishSignal chan *channel

	rangeProcessing int
	timeLogger      time.Time
}

func NewService(processing, rangeProc, channels int, t time.Time) *service {
	que := NewQueue()

	serv := &service{
		finishSignal:    make(chan *channel),
		rangeProcessing: rangeProc,
		que:             que,
		channels:        MakeChannels(channels, processing, que),
		timeLogger:      t,
	}

	serv.processed.Store(0)

	<-time.After(1 * time.Second)
	go serv.QueListen()

	return serv
}

func (s *service) FilledChannels() int {
	counter := 0

	for _, channel := range s.channels {
		if channel.IsActive() {
			counter++
		}
	}

	return counter
}

func (s *service) Process() {
	stored := false

	for _, channel := range s.channels {
		if !channel.IsActive() && s.que.Length() == 0 {
			channel.active.Store(true)

			go func() {
				s.finishSignal <- channel.SetActive()
			}()

			stored = true
			break
		}
	}

	if !stored {
		s.que.AppendRecord()
	}
}

func (s *service) QueListen() {
	for channel := range s.finishSignal {
		s.processed.Store(s.processed.Load().(int) + 1)

		if s.que.Length() > 0 {
			if !channel.IsActive() {
				channel.active.Store(true)

				s.que.DepartRecord()

				go func() {
					s.finishSignal <- channel.SetActive()
				}()
			}
		}
	}
}

func (s *service) Processed(records int, ttt float64) {
	fmt.Printf("processed records: %d; queued records: %d; all records: %d; \ntime: %d; \nutil: %.2f, expectedUtil: %.2f; \nfilledChannels: %d \n",
		s.processed.Load().(int),
		s.que.Length(),
		records,
		time.Since(s.timeLogger).Milliseconds(),
		float64(s.processed.Load().(int))/float64(records),
		ttt,
		s.FilledChannels())
}
