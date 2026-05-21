package main

import (
	"context"
	"log"
	"os"

	rkboot "github.com/rookie-ninja/rk-boot/v2"
	rkpostgres "github.com/rookie-ninja/rk-db/postgres"
	rkgin "github.com/rookie-ninja/rk-gin/v2/boot"

	app_client "developerApplication/internal/adapters/grpc/application_service"
	av_client "developerApplication/internal/adapters/grpc/auto_verification"
	"developerApplication/internal/adapters/repository/postgres/developer_application"
	resthandlers "developerApplication/internal/adapters/rest"
	"developerApplication/internal/adapters/s3"
	"developerApplication/internal/core/service"
)

func main() {
	raw, err := os.ReadFile("boot.yaml")
	if err != nil {
		log.Fatalf("failed to read boot.yaml: %v", err)
	}

	boot := rkboot.NewBoot(rkboot.WithBootConfigRaw([]byte(os.ExpandEnv(string(raw)))))
	boot.Bootstrap(context.TODO())

	gin := rkgin.GetGinEntry("developer_application")

	pgEntry := rkpostgres.GetPostgresEntry("developer_application_postgres")
	if pgEntry == nil {
		log.Fatal("postgres entry not found")
	}

	dbEntry := pgEntry.GetDB("developer_application")
	db, err := dbEntry.DB()
	if err != nil {
		panic(err)
	}

	minioEndpoint := getEnv("MINIO_ENDPOINT", "localhost:9000")
	minioAccessKey := getEnv("MINIO_ACCESS_KEY", "minio")
	minioSecretKey := getEnv("MINIO_SECRET_KEY", "minio123")
	minioBucket := getEnv("MINIO_BUCKET", "app-builds")
	minioPublicEndpoint := getEnv("MINIO_PUBLIC_ENDPOINT", minioEndpoint)
	minioRegion := getEnv("MINIO_REGION", "us-east-1")

	autoVerificationAddr := getEnv("AUTO_VERIFICATION_ADDR", "auto:8080")
	applicationServiceAddr := getEnv("APPLICATION_SERVICE_ADDR", "application:9090")

	storage := s3.NewMinio(minioEndpoint, minioAccessKey, minioSecretKey, minioBucket, minioPublicEndpoint, minioRegion)

	verifier, err := av_client.NewClient(autoVerificationAddr)
	if err != nil {
		log.Fatalf("failed to connect to auto verification: %v", err)
	}
	defer verifier.Close()

	appServiceClient, err := app_client.NewClient(applicationServiceAddr)
	if err != nil {
		log.Fatalf("failed to connect to application service: %v", err)
	}
	defer appServiceClient.Close()

	repository := developer_application.New(db)
	svc := service.NewService(repository, storage, verifier, appServiceClient)
	handler := resthandlers.NewHandler(svc)

	gin.Router.POST("/applications", handler.UploadBuild)
	gin.Router.GET("/applications", handler.ListApplications)
	gin.Router.GET("/applications/:id", handler.GetApplication)
	gin.Router.PATCH("/applications/:id", handler.UpdateApplication)
	gin.Router.POST("/applications/:id/verify", handler.StartVerification)
	gin.Router.GET("/applications/:id/verify", handler.GetVerificationStatus)
	gin.Router.POST("/applications/:id/publish", handler.PublishApplication)

	boot.WaitForShutdownSig(context.TODO())
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
