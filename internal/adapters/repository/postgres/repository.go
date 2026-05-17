package postgres

import (
	"context"
	"developerApplication/internal/adapters/repository/postgres/developer_application"

	"github.com/google/uuid"
)

type DeveloperApplicationRepository interface {
	CreateDeveloperApplication(ctx context.Context, arg developer_application.CreateDeveloperApplicationParams) (developer_application.DeveloperApplication, error)
	GetDeveloperApplication(ctx context.Context, id uuid.UUID) (developer_application.DeveloperApplication, error)
	ListDeveloperApplications(ctx context.Context, arg developer_application.ListDeveloperApplicationsParams) ([]developer_application.DeveloperApplication, error)
	SetVerificationProcess(ctx context.Context, arg developer_application.SetVerificationProcessParams) (developer_application.DeveloperApplication, error)
	UpdateVerificationStatus(ctx context.Context, arg developer_application.UpdateVerificationStatusParams) error
	MarkAsPublished(ctx context.Context, id uuid.UUID) error
}
