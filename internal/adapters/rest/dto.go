package rest

import "developerApplication/internal/adapters/repository/postgres/developer_application"

type UploadBuildInDTO struct {
	DeveloperID string `form:"developer_id" binding:"required"`
	Version     string `form:"version"      binding:"required"`
}

type PublishApplicationInDTO struct {
	DeveloperID string `json:"developer_id" binding:"required"`
	CodeName    string `json:"code_name"    binding:"required"`
	CategoryID  int64  `json:"category_id"  binding:"required"`
}

type ListApplicationsInDTO struct {
	DeveloperID string `form:"developer_id" binding:"required"`
}

type ApplicationOutDTO struct {
	developer_application.DeveloperApplication
}

type ListApplicationsOutDTO struct {
	Items []developer_application.DeveloperApplication `json:"items"`
	Page  int                                          `json:"page"`
	Size  int                                          `json:"size"`
}

type VerificationStatusOutDTO struct {
	Status string `json:"status"`
}

type ErrorResponseDTO struct {
	Error string `json:"error"`
}
