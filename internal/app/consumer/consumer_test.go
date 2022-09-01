package consumer_test

import (
	"testing"
	"time"

	"github.com/alexeykirinyuk/retranslator/internal/app/consumer"
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

func TestConsume(t *testing.T) {
	t.Parallel()

	repo := newMockEventRepo()

	c := consumer.NewConsumer(
		uint64(2),
		5, // buffer
		&repo,
		uint64(100),
		time.Microsecond,
	)

	t.Run("consumer->consume", func(t *testing.T) {
		c.Consume()
	})

	t.Run("consumer->read", func(t *testing.T) {
		<-c.GetOutChannel()
	})

	t.Run("consumer->stop", func(t *testing.T) {
		c.Stop()
	})

	t.Run("consumer->consume 2", func(t *testing.T) {
		c.Consume()
	})

	t.Run("consumer->close", func(t *testing.T) {
		defer c.Close()
		c.Stop()
	})
}
