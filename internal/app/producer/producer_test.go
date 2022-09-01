package producer_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/alexeykirinyuk/retranslator/internal/app/producer"
	"github.com/alexeykirinyuk/retranslator/internal/model"
	"github.com/stretchr/testify/assert"
)

type mockEventSender struct {
	send bool
	mtx  *sync.Mutex
}

func newMockEventRepo() *mockEventSender {
	return &mockEventSender{
		send: false,
		mtx:  &sync.Mutex{},
	}
}

func (m *mockEventSender) Send(event model.ProductEvent) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if m.send {
		return errors.New("test err")
	}

	m.send = true
	return nil
}

func TestProducer(t *testing.T) {

	in := make(chan model.ProductEvent)
	defer close(in)

	p := producer.NewProducer(
		uint64(2),
		in,
		newMockEventRepo(),
		10,
		time.Millisecond*100,
	)

	productOk := model.ProductEvent{
		ID: 100,
	}
	productErr := model.ProductEvent{
		ID: 101,
	}

	t.Run("producer->produce", func(t *testing.T) {
		p.Produce()
	})

	t.Run("producer->ok-channel", func(t *testing.T) {
		in <- productOk
		res := <-p.GetOkChannel()
		assert.Equal(t, productOk, res)
	})

	t.Run("producer->err-channel", func(t *testing.T) {
		in <- productErr
		res := <-p.GetErrChannel()
		assert.Equal(t, productErr, res)
	})

	t.Run("producer->stop", func(t *testing.T) {
		in <- productOk
		in <- productErr
		p.Stop()
	})

	t.Run("producer->produce", func(t *testing.T) {
		p.Produce()
	})

	t.Run("producer->close", func(t *testing.T) {
		defer p.Close()

		p.Stop()
	})
}
