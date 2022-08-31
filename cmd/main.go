package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/alexeykirinyuk/retranslator/internal/app/retranslator"
	"github.com/alexeykirinyuk/retranslator/internal/storage/kafka"
	"github.com/alexeykirinyuk/retranslator/internal/storage/psql"
)

func main() {
	cfg := retranslator.Cfg{
		ConsumerCount: 100,
		ProducerCount: 100,
		BatchSize:     100,
		Timeout:       time.Second,
		Sender:        &kafka.KafkaEventSender{},
		Repo:          &psql.PsqlEventRepo{},
	}

	r := retranslator.NewRetranslator(cfg)

	fmt.Println("r.Run()")
	r.Run()

	fmt.Println("time.Sleep(time.Second * 3)")
	time.Sleep(time.Second * 3)

	fmt.Println("r.Stop()")
	r.Stop()

	fmt.Printf("runtime.NumGoroutine()=%d\n", runtime.NumGoroutine())
}
