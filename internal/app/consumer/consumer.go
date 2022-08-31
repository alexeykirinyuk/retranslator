package consumer

import (
	"context"
	"sync"
	"time"

	"github.com/alexeykirinyuk/retranslator/internal/model"
)

type eventRepo interface {
	Lock(batchSize uint64) ([]*model.ProductEvent, error)
}

type Consumer struct {
	routines uint64

	batchSize uint64
	timeout   time.Duration

	repo eventRepo
	wg   *sync.WaitGroup

	ch     chan *model.ProductEvent
	ctx    context.Context
	cancel context.CancelFunc
}

func NewConsumer(
	routines uint64,
	repo eventRepo,
	batchSize uint64,
	timeout time.Duration,
) *Consumer {
	ctx, cancel := context.WithCancel(context.Background())

	return &Consumer{
		routines:  routines,
		ch:        make(chan *model.ProductEvent),
		repo:      repo,
		batchSize: batchSize,
		wg:        &sync.WaitGroup{},
		ctx:       ctx,
		cancel:    cancel,
		timeout:   timeout,
	}
}

func (c *Consumer) GetOutChannel() <-chan *model.ProductEvent {
	return c.ch
}

func (c *Consumer) Consume() {
	for i := uint64(0); i < c.routines; i++ {
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()

			for {
				timer := time.NewTicker(c.timeout)

				select {
				case <-c.ctx.Done():
					return
				case <-timer.C:
					consumeEvents(c)
				}
			}
		}()
	}
}

func consumeEvents(c *Consumer) {
	syncChannel := make(chan []*model.ProductEvent, 1)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		defer close(syncChannel)

		events, err := c.repo.Lock(c.batchSize)
		if err != nil {
			syncChannel <- events
		}
	}()

	select {
	case <-c.ctx.Done():
		return
	case events, ok := <-syncChannel:
		if ok {
			for _, event := range events {
				c.ch <- event
			}
		}
	}
}

func (c *Consumer) Stop() {
	c.cancel()
	c.wg.Wait()
	close(c.ch)
}
