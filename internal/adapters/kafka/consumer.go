package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

type Handler func(ctx context.Context, msg kafka.Message) error

type ConsumerConfig struct {
	Brokers  []string
	GroupID  string
	Topic    string
	MinBytes int
	MaxBytes int
}

type Consumer struct {
	reader  *kafka.Reader
	handler Handler
}

func NewConsumer(cfg ConsumerConfig, handler Handler) *Consumer {
	minBytes := cfg.MinBytes
	if minBytes == 0 {
		minBytes = 1e3 // 1KB
	}
	maxBytes := cfg.MaxBytes
	if maxBytes == 0 {
		maxBytes = 10e6 // 10MB
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		GroupID:        cfg.GroupID,
		Topic:          cfg.Topic,
		MinBytes:       minBytes,
		MaxBytes:       maxBytes,
		CommitInterval: 0, // manual commit
	})

	return &Consumer{
		reader:  reader,
		handler: handler,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return err
		}

		if err := c.handler(ctx, msg); err != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(250 * time.Millisecond):
			}
			continue
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return err
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
