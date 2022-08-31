package psql

import (
	"time"

	"github.com/alexeykirinyuk/retranslator/internal/model"
)

type PsqlEventRepo struct {
}

var events = []*model.ProductEvent{
	{
		ID:          1,
		EventType:   model.EventType_ProductCreated,
		EventStatus: model.EventStatus_Created,
		OccuredAt:   time.Now().UTC(),
		Data:        &model.Product{},
	},
	{
		ID:          2,
		EventType:   model.EventType_ProductUpdated,
		EventStatus: model.EventStatus_Created,
		OccuredAt:   time.Now().UTC(),
		Data:        &model.Product{},
	},
}

func (p *PsqlEventRepo) Lock(batchSize uint64) ([]*model.ProductEvent, error) {
	return events, nil
}
