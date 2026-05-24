package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
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
	cfg     ConsumerConfig
	handler Handler
}

func NewConsumer(cfg ConsumerConfig, handler Handler) *Consumer {
	return &Consumer{cfg: cfg, handler: handler}
}

func (c *Consumer) newReader() *kafka.Reader {
	minBytes := c.cfg.MinBytes
	if minBytes == 0 {
		minBytes = 1e3
	}
	maxBytes := c.cfg.MaxBytes
	if maxBytes == 0 {
		maxBytes = 10e6
	}

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:           c.cfg.Brokers,
		GroupID:           c.cfg.GroupID,
		Topic:             c.cfg.Topic,
		MinBytes:          minBytes,
		MaxBytes:          maxBytes,
		CommitInterval:    0, // manual commit
		SessionTimeout:    30 * time.Second,
		HeartbeatInterval: 3 * time.Second,
		MaxWait:           10 * time.Second,
		StartOffset:       kafka.LastOffset,
		ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			log.Printf("[kafka-go:error] "+msg, args...)
		}),
	})
}

// Run loops forever, reconnecting on any non-context error.
func (c *Consumer) Run(ctx context.Context) error {
	for {
		if err := c.runOnce(ctx); err != nil {
			return err // only context cancellation propagates up
		}
	}
}

func (c *Consumer) runOnce(ctx context.Context) error {
	reader := c.newReader()
	defer reader.Close()

	log.Printf("kafka consumer: connecting to topic=%s group=%s", c.cfg.Topic, c.cfg.GroupID)

	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			log.Printf("kafka consumer: FetchMessage error on topic=%s, reconnecting in 2s: %v", c.cfg.Topic, err)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(2 * time.Second):
			}
			return nil // return nil → outer loop recreates the reader
		}

		if err := c.handler(ctx, msg); err != nil {
			log.Printf("kafka consumer: handler error on topic=%s (will retry msg): %v", c.cfg.Topic, err)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(250 * time.Millisecond):
			}
			continue
		}

		if err := reader.CommitMessages(ctx, msg); err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			log.Printf("kafka consumer: CommitMessages error on topic=%s, reconnecting: %v", c.cfg.Topic, err)
			return nil // recreate reader
		}
	}
}
