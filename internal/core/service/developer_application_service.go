package service

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	app_client "developerApplication/internal/adapters/grpc/application_service"
	app_gen "developerApplication/internal/adapters/grpc/application_service/gen"
	av_client "developerApplication/internal/adapters/grpc/auto_verification"
	"developerApplication/internal/adapters/repository/postgres"
	"developerApplication/internal/adapters/repository/postgres/developer_application"
	"developerApplication/internal/adapters/s3"

	"github.com/google/uuid"
)

const androidDeviceTypeID = 1

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

func (s *Service) UploadBuild(ctx context.Context, p UploadBuildParams) (developer_application.DeveloperApplication, error) {
	if err := s.storage.UploadFile(ctx, p.FileName, p.Reader, p.Size, "application/vnd.android.package-archive"); err != nil {
		return developer_application.DeveloperApplication{}, fmt.Errorf("upload to storage: %w", err)
	}

	return s.repo.CreateDeveloperApplication(ctx, developer_application.CreateDeveloperApplicationParams{
		DeveloperID:          uuid.MustParse(p.DeveloperID),
		CodeName:             p.CodeName,
		CategoryID:           p.CategoryID,
		AndroidPackageName:   p.AndroidPackageName,
		DefaultLocale:        p.DefaultLocale,
		Name:                 p.Name,
		ShortTitle:           p.ShortTitle,
		Description:          nullString(p.Description),
		Goals:                p.Goals,
		Tasks:                p.Tasks,
		Results:              nullString(p.Results),
		Challenges:           nullString(p.Challenges),
		Location:             nullString(p.Location),
		VideoCover:           nullString(p.VideoCover),
		Safety:               nullString(p.Safety),
		WebVideo:             p.WebVideo,
		InappVideo:           p.InappVideo,
		WebBackgroundImage:   p.WebBackgroundImage,
		InappBackgroundImage: nullString(p.InappBackgroundImage),
		ApkFilename:          p.FileName,
		Version:              p.Version,
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

func (s *Service) UpdateApplication(ctx context.Context, id string, p UpdateApplicationParams) (developer_application.DeveloperApplication, error) {
	app, err := s.repo.GetDeveloperApplication(ctx, uuid.MustParse(id))
	if err != nil {
		return developer_application.DeveloperApplication{}, err
	}

	applyString := func(dst *string, src *string) {
		if src != nil {
			*dst = *src
		}
	}
	applyInt64 := func(dst *int64, src *int64) {
		if src != nil {
			*dst = *src
		}
	}

	applyString(&app.CodeName, p.CodeName)
	applyInt64(&app.CategoryID, p.CategoryID)
	applyString(&app.AndroidPackageName, p.AndroidPackageName)
	applyString(&app.DefaultLocale, p.DefaultLocale)
	applyString(&app.WebVideo, p.WebVideo)
	applyString(&app.InappVideo, p.InappVideo)
	applyString(&app.WebBackgroundImage, p.WebBackgroundImage)
	applyString(&app.Name, p.Name)
	applyString(&app.ShortTitle, p.ShortTitle)
	applyString(&app.Goals, p.Goals)
	applyString(&app.Tasks, p.Tasks)
	applyString(&app.Version, p.Version)

	if p.InappBackgroundImage != nil {
		app.InappBackgroundImage = nullString(*p.InappBackgroundImage)
	}
	if p.Description != nil {
		app.Description = nullString(*p.Description)
	}
	if p.Results != nil {
		app.Results = nullString(*p.Results)
	}
	if p.Challenges != nil {
		app.Challenges = nullString(*p.Challenges)
	}
	if p.Location != nil {
		app.Location = nullString(*p.Location)
	}
	if p.VideoCover != nil {
		app.VideoCover = nullString(*p.VideoCover)
	}
	if p.Safety != nil {
		app.Safety = nullString(*p.Safety)
	}

	return s.repo.UpdateDeveloperApplication(ctx, developer_application.UpdateDeveloperApplicationParams{
		ID:                   app.ID,
		CodeName:             app.CodeName,
		CategoryID:           app.CategoryID,
		AndroidPackageName:   app.AndroidPackageName,
		DefaultLocale:        app.DefaultLocale,
		WebVideo:             app.WebVideo,
		InappVideo:           app.InappVideo,
		WebBackgroundImage:   app.WebBackgroundImage,
		InappBackgroundImage: app.InappBackgroundImage,
		Name:                 app.Name,
		ShortTitle:           app.ShortTitle,
		Description:          app.Description,
		Goals:                app.Goals,
		Tasks:                app.Tasks,
		Results:              app.Results,
		Challenges:           app.Challenges,
		Location:             app.Location,
		VideoCover:           app.VideoCover,
		Safety:               app.Safety,
		Version:              app.Version,
	})
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

	return s.repo.SetVerificationProcess(ctx, developer_application.SetVerificationProcessParams{
		ID:                    app.ID,
		VerificationProcessID: uuid.NullUUID{UUID: uuid.MustParse(processID), Valid: true},
	})
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

func (s *Service) PublishApplication(ctx context.Context, id string) error {
	app, err := s.repo.GetDeveloperApplication(ctx, uuid.MustParse(id))
	if err != nil {
		return err
	}

	if !app.VerificationProcessID.Valid || app.VerificationStatus.String != "completed" {
		return ErrNotVerified
	}

	major, minor, micro, err := parseVersion(app.Version)
	if err != nil {
		return fmt.Errorf("parse version: %w", err)
	}

	appID, err := s.appService.CreateApplication(ctx, &app_gen.CreateApplicationRequest{
		DeveloperID:          app.DeveloperID.String(),
		CodeName:             app.CodeName,
		CategoryID:           app.CategoryID,
		DefaultLocale:        app.DefaultLocale,
		WebVideo:             app.WebVideo,
		InappVideo:           app.InappVideo,
		WebBackgroundImage:   app.WebBackgroundImage,
		InappBackgroundImage: app.InappBackgroundImage.String,
		VideoCover:           ptrString(app.VideoCover),
		Name:                 app.Name,
		ShortTitle:           app.ShortTitle,
		Goals:                app.Goals,
		Tasks:                app.Tasks,
		Description:          ptrString(app.Description),
		Results:              ptrString(app.Results),
		Challenges:           ptrString(app.Challenges),
		Location:             ptrString(app.Location),
		Safety:               ptrString(app.Safety),
	})
	if err != nil {
		return fmt.Errorf("create application: %w", err)
	}

	_, err = s.appService.CreateRepository(ctx, &app_gen.CreateRepositoryRequest{
		ApplicationID: appID,
		DeviceTypeID:  androidDeviceTypeID,
		MajorVersion:  int32(major),
		MinorVersion:  int32(minor),
		MicroVersion:  int32(micro),
		BuildTypeID:   "release",
		LaunchUrl:     app.AndroidPackageName,
	})
	if err != nil {
		return fmt.Errorf("create repository: %w", err)
	}

	return s.repo.MarkAsPublished(ctx, app.ID)
}

func parseVersion(version string) (int, int, int, error) {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("expected major.minor.micro format, got %q", version)
	}
	nums := [3]int{}
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("invalid version component %q: %w", p, err)
		}
		nums[i] = n
	}
	return nums[0], nums[1], nums[2], nil
}

func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func ptrString(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}
