package retranslator

import (
	"time"

	"github.com/alexeykirinyuk/retranslator/internal/app/consumer"
	"github.com/alexeykirinyuk/retranslator/internal/app/producer"
	"github.com/alexeykirinyuk/retranslator/internal/model"
)

type eventRepo interface {
	Lock(batchSize uint64) ([]model.ProductEvent, error)
}

type eventSender interface {
	Send(subdomain model.ProductEvent) error
}

type Retranslator struct {
	consumer *consumer.Consumer
	producer *producer.Producer
}

type Cfg struct {
	ConsumerCount uint64
	ConsumerSize  uint64
	ProducerCount uint64

	BatchSize uint64
	Timeout   time.Duration

	Sender eventSender
	Repo   eventRepo
}

func NewRetranslator(cfg Cfg) *Retranslator {
	consumer := consumer.NewConsumer(
		cfg.ConsumerCount,
		cfg.ConsumerSize,
		cfg.Repo,
		cfg.BatchSize,
		cfg.Timeout,
	)

	producer := producer.NewProducer(
		cfg.ProducerCount,
		consumer.GetOutChannel(),
		cfg.Sender,
		cfg.BatchSize,
		cfg.Timeout,
	)

	return &Retranslator{
		consumer: consumer,
		producer: producer,
	}
}

func (r *Retranslator) Run() {
	r.consumer.Consume()
	r.producer.Produce()
}

func (r *Retranslator) Stop() {
	r.producer.Stop()
	r.consumer.Stop()
}
