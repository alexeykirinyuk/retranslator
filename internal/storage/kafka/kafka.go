package kafka

import "github.com/alexeykirinyuk/retranslator/internal/model"

type KafkaEventSender struct {
}

func (k *KafkaEventSender) Send(subdomain model.ProductEvent) error {
	return nil
}
