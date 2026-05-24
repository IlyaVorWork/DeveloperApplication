package kafka

import (
	"context"
	"database/sql"
	"developerApplication/internal/adapters/repository/postgres/developer_application"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	kafkago "github.com/segmentio/kafka-go"
	"log"
	"strings"
	"time"
)

type VerificationEvent struct {
	EventID       string    `json:"EventID"`
	CorrelationID string    `json:"CorrelationID"`
	Type          string    `json:"Type"`
	FileName      string    `json:"FileName"`
	Timestamp     time.Time `json:"Timestamp"`
}

type VerificationStatusUpdater interface {
	UpdateVerificationStatus(ctx context.Context, arg developer_application.UpdateVerificationStatusParams) error
}

func NewVerificationCompletedHandler(repo VerificationStatusUpdater) Handler {
	return func(ctx context.Context, msg kafkago.Message) error {
		var event VerificationEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("verification.completed: failed to unmarshal event: %v", err)
			return err
		}

		log.Printf("verification.completed: received event_id=%s correlation_id=%s type=%s",
			event.EventID, event.CorrelationID, event.Type)

		processUUID, err := uuid.Parse(event.CorrelationID)
		if err != nil {
			return fmt.Errorf("verification.completed: invalid correlation_id %q: %w", event.CorrelationID, err)
		}

		status := "completed"
		if strings.Contains(event.Type, "failed") {
			status = "failed"
		}

		if err := repo.UpdateVerificationStatus(ctx, developer_application.UpdateVerificationStatusParams{
			VerificationProcessID: uuid.NullUUID{UUID: processUUID, Valid: true},
			VerificationStatus:    sql.NullString{String: status, Valid: true},
		}); err != nil {
			log.Printf("verification.completed: failed to update status for process %s: %v", event.CorrelationID, err)
			return err
		}

		log.Printf("verification.completed: updated status=%s for process %s", status, event.CorrelationID)
		return nil
	}
}
