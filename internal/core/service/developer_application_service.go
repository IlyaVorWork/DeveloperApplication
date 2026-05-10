package service

import (
	"context"
	"fmt"
	"io"

	app_client "developerApplication/internal/adapters/grpc/application_service"
	app_gen "developerApplication/internal/adapters/grpc/application_service/gen"
	av_client "developerApplication/internal/adapters/grpc/auto_verification"
	"developerApplication/internal/adapters/repository/postgres"
	"developerApplication/internal/adapters/repository/postgres/developer_application"
	"developerApplication/internal/adapters/s3"

	"github.com/google/uuid"
)

type Service struct {
	repo       postgres.DeveloperApplicationRepository
	storage    s3.FileStorage
	verifier   *av_client.Client
	appService *app_client.Client
}

func NewService(
	repo postgres.DeveloperApplicationRepository,
	storage s3.FileStorage,
	verifier *av_client.Client,
	appService *app_client.Client,
) *Service {
	return &Service{
		repo:       repo,
		storage:    storage,
		verifier:   verifier,
		appService: appService,
	}
}

func (s *Service) UploadBuild(ctx context.Context, developerID, version string, fileName string, reader io.Reader, size int64) (developer_application.DeveloperApplication, error) {
	if err := s.storage.UploadFile(ctx, fileName, reader, size, "application/vnd.android.package-archive"); err != nil {
		return developer_application.DeveloperApplication{}, fmt.Errorf("upload to storage: %w", err)
	}

	return s.repo.CreateDeveloperApplication(ctx, developer_application.CreateDeveloperApplicationParams{
		DeveloperID: uuid.MustParse(developerID),
		ApkFilename: fileName,
		Version:     version,
	})
}

func (s *Service) GetApplication(ctx context.Context, id string) (developer_application.DeveloperApplication, error) {
	return s.repo.GetDeveloperApplication(ctx, uuid.MustParse(id))
}

func (s *Service) ListApplications(ctx context.Context, developerID string, page, size int) ([]developer_application.DeveloperApplication, error) {
	list, err := s.repo.ListDeveloperApplications(ctx, developer_application.ListDeveloperApplicationsParams{
		DeveloperID: uuid.MustParse(developerID),
		Limit:       int32(size),
		Offset:      int32(page * size),
	})
	if err != nil {
		return nil, err
	}
	if list == nil {
		return []developer_application.DeveloperApplication{}, nil
	}
	return list, nil
}

func (s *Service) StartVerification(ctx context.Context, id string) (developer_application.DeveloperApplication, error) {
	app, err := s.repo.GetDeveloperApplication(ctx, uuid.MustParse(id))
	if err != nil {
		return developer_application.DeveloperApplication{}, err
	}

	processID, err := s.verifier.StartVerification(ctx, app.ApkFilename)
	if err != nil {
		return developer_application.DeveloperApplication{}, fmt.Errorf("start verification: %w", err)
	}

	return s.repo.SetVerificationProcess(ctx, app.ID, uuid.MustParse(processID))
}

func (s *Service) GetVerificationStatus(ctx context.Context, id string) (string, error) {
	app, err := s.repo.GetDeveloperApplication(ctx, uuid.MustParse(id))
	if err != nil {
		return "", err
	}

	if !app.VerificationProcessID.Valid {
		return "not_started", nil
	}

	processes, err := s.verifier.GetVerification(ctx, app.VerificationProcessID.UUID.String())
	if err != nil {
		return "", fmt.Errorf("get verification: %w", err)
	}

	if len(processes) == 0 {
		return app.VerificationStatus.String, nil
	}

	return processes[0].Status, nil
}

func (s *Service) PublishApplication(ctx context.Context, id, developerID, codeName string, categoryID int64) error {
	app, err := s.repo.GetDeveloperApplication(ctx, uuid.MustParse(id))
	if err != nil {
		return err
	}

	if !app.VerificationProcessID.Valid || app.VerificationStatus.String != "human.verify.succeeded" {
		return ErrNotVerified
	}

	appID, err := s.appService.CreateApplication(ctx, &app_gen.CreateApplicationRequest{
		DeveloperID: developerID,
		CodeName:    codeName,
		CategoryID:  categoryID,
	})
	if err != nil {
		return fmt.Errorf("create application: %w", err)
	}

	_, err = s.appService.CreateRepository(ctx, &app_gen.CreateRepositoryRequest{
		ApplicationID: appID,
		LaunchUrl:     app.ApkFilename,
	})
	if err != nil {
		return fmt.Errorf("create repository: %w", err)
	}

	return s.repo.MarkAsPublished(ctx, app.ID)
}
