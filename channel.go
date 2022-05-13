package main

import (
	"sync/atomic"
	"time"
)

type channel struct {
	active     atomic.Value
	processing time.Duration
	que        *queue
}

func MakeChannels(i, t int, que *queue) []*channel {
	chans := make([]*channel, 0)

	for range make([]int, i) {
		chann := &channel{
			processing: time.Millisecond * time.Duration(t),
			que:        que,
		}
		chann.active.Store(false)
		chans = append(chans, chann)
	}

	return chans
}

func (c *channel) SetActive() *channel {
	c.proc()
	return c
}

func (c *channel) IsActive() bool {
	return c.active.Load().(bool)
}

func (c *channel) proc() {
	<-time.After(c.processing)
	c.active.Store(false)
}
