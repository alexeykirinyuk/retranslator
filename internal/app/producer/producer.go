package producer

import (
	"context"
	"sync"
	"time"

	"github.com/alexeykirinyuk/retranslator/internal/model"
)

type eventSender interface {
	Send(subdomain *model.ProductEvent) error
}

type Producer struct {
	routines uint64
	inCh     <-chan *model.ProductEvent

	okCh  chan *model.ProductEvent
	errCh chan *model.ProductEvent

	batchSize uint64
	timeout   time.Duration

	sender eventSender
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewProducer(
	routines uint64,
	inCh <-chan *model.ProductEvent,
	sender eventSender,
	batchSize uint64,
	timeout time.Duration,
) *Producer {
	ctx, cancel := context.WithCancel(context.Background())

	return &Producer{
		routines:  routines,
		inCh:      inCh,
		okCh:      make(chan *model.ProductEvent),
		errCh:     make(chan *model.ProductEvent),
		batchSize: batchSize,
		timeout:   timeout,
		sender:    sender,
		wg:        &sync.WaitGroup{},
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (p *Producer) Produce() {
	for i := uint64(0); i < p.routines; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()

			for event := range p.inCh {
				sendEvent(p, event)
			}
		}()
	}
}

func sendEvent(p *Producer, event *model.ProductEvent) {
	tmpCh := make(chan bool, 1)

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		defer close(tmpCh)

		err := p.sender.Send(event)
		tmpCh <- err != nil
	}()

	select {
	case <-p.ctx.Done():
		return
	case val, ok := <-tmpCh:
		if ok && val {
			p.okCh <- event
		} else {
			p.errCh <- event
		}
	}
}

func (p *Producer) GetOkChannel() <-chan *model.ProductEvent {
	return p.okCh
}

func (p *Producer) GetErrChannel() <-chan *model.ProductEvent {
	return p.errCh
}

func (p *Producer) Stop() {
	p.cancel()
	p.wg.Wait()

	close(p.okCh)
	close(p.errCh)
}
