package retranslator_test

import (
	"testing"
	"time"

	"github.com/alexeykirinyuk/retranslator/internal/app/retranslator"
	"github.com/alexeykirinyuk/retranslator/internal/model"
)

type mockEventRepo struct {
}

func newMockEventRepo() mockEventRepo {
	return mockEventRepo{}
}

func (m *mockEventRepo) Lock(batchSize uint64) ([]model.ProductEvent, error) {
	events := make([]model.ProductEvent, 0, 20) // buffer - 5
	for i := uint64(0); i < 20; i++ {
		events = append(events, model.ProductEvent{ID: i})
	}

	return events, nil
}

type mockEventSender struct {
}

func (m *mockEventSender) Send(event model.ProductEvent) error {
	return nil
}

func TestRetranslator(t *testing.T) {
	t.Parallel()

	cfg := retranslator.Cfg{
		ConsumerCount: 100,
		ProducerCount: 100,
		BatchSize:     100,
		Timeout:       time.Second,
		Sender:        &mockEventSender{},
		Repo:          &mockEventRepo{},
	}

	r := retranslator.NewRetranslator(cfg)

	t.Run("run retranslator", func(t *testing.T) {
		r.Run()
	})

	t.Run("stop retranslator", func(t *testing.T) {
		time.Sleep(time.Second)
		r.Stop()
	})
}
