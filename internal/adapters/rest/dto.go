package rest

import (
	"developerApplication/internal/adapters/repository/postgres/developer_application"
	"time"
)

type UploadBuildInDTO struct {
	DeveloperID          string `form:"developer_id"           binding:"required"`
	CodeName             string `form:"code_name"              binding:"required"`
	CategoryID           int64  `form:"category_id"`
	AndroidPackageName   string `form:"android_package_name"   binding:"required"`
	DefaultLocale        string `form:"default_locale"         binding:"required"`
	Name                 string `form:"name"                   binding:"required"`
	ShortTitle           string `form:"short_title"            binding:"required"`
	Version              string `form:"version"                binding:"required"`
	Description          string `form:"description"`
	Goals                string `form:"goals"`
	Tasks                string `form:"tasks"`
	Results              string `form:"results"`
	Challenges           string `form:"challenges"`
	Location             string `form:"location"`
	VideoCover           string `form:"video_cover"`
	Safety               string `form:"safety"`
	WebVideo             string `form:"web_video"`
	InappVideo           string `form:"inapp_video"`
	WebBackgroundImage   string `form:"web_background_image"`
	InappBackgroundImage string `form:"inapp_background_image"`
}

type UpdateApplicationInDTO struct {
	CodeName             *string `json:"code_name"`
	CategoryID           *int64  `json:"category_id"`
	AndroidPackageName   *string `json:"android_package_name"`
	DefaultLocale        *string `json:"default_locale"`
	WebVideo             *string `json:"web_video"`
	InappVideo           *string `json:"inapp_video"`
	WebBackgroundImage   *string `json:"web_background_image"`
	InappBackgroundImage *string `json:"inapp_background_image"`
	Name                 *string `json:"name"`
	ShortTitle           *string `json:"short_title"`
	Description          *string `json:"description"`
	Goals                *string `json:"goals"`
	Tasks                *string `json:"tasks"`
	Results              *string `json:"results"`
	Challenges           *string `json:"challenges"`
	Location             *string `json:"location"`
	VideoCover           *string `json:"video_cover"`
	Safety               *string `json:"safety"`
	Version              *string `json:"version"`
}

type ListApplicationsInDTO struct {
	DeveloperID string `form:"developer_id" binding:"required"`
}

type ApplicationOutDTO struct {
	ID                    string    `json:"id"`
	DeveloperID           string    `json:"developer_id"`
	CodeName              string    `json:"code_name"`
	CategoryID            int64     `json:"category_id"`
	AndroidPackageName    string    `json:"android_package_name"`
	DefaultLocale         string    `json:"default_locale"`
	WebVideo              string    `json:"web_video"`
	InappVideo            string    `json:"inapp_video"`
	WebBackgroundImage    string    `json:"web_background_image"`
	InappBackgroundImage  *string   `json:"inapp_background_image,omitempty"`
	Name                  string    `json:"name"`
	ShortTitle            string    `json:"short_title"`
	Description           *string   `json:"description,omitempty"`
	Goals                 string    `json:"goals"`
	Tasks                 string    `json:"tasks"`
	Results               *string   `json:"results,omitempty"`
	Challenges            *string   `json:"challenges,omitempty"`
	Location              *string   `json:"location,omitempty"`
	VideoCover            *string   `json:"video_cover,omitempty"`
	Safety                *string   `json:"safety,omitempty"`
	ApkFilename           string    `json:"apk_filename"`
	Version               string    `json:"version"`
	VerificationProcessID *string   `json:"verification_process_id,omitempty"`
	VerificationStatus    *string   `json:"verification_status,omitempty"`
	Published             bool      `json:"published"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type ListApplicationsOutDTO struct {
	Items []ApplicationOutDTO `json:"items"`
	Page  int                 `json:"page"`
	Size  int                 `json:"size"`
}

type VerificationStatusOutDTO struct {
	Status     string  `json:"status"`
	FailedStep *string `json:"failed_step,omitempty"`
}

type ErrorResponseDTO struct {
	Error string `json:"error"`
}

func toApplicationOutDTO(a developer_application.DeveloperApplication) ApplicationOutDTO {
	dto := ApplicationOutDTO{
		ID:                 a.ID.String(),
		DeveloperID:        a.DeveloperID.String(),
		CodeName:           a.CodeName,
		CategoryID:         a.CategoryID,
		AndroidPackageName: a.AndroidPackageName,
		DefaultLocale:      a.DefaultLocale,
		WebVideo:           a.WebVideo,
		InappVideo:         a.InappVideo,
		WebBackgroundImage: a.WebBackgroundImage,
		Name:               a.Name,
		ShortTitle:         a.ShortTitle,
		Goals:              a.Goals,
		Tasks:              a.Tasks,
		ApkFilename:        a.ApkFilename,
		Version:            a.Version,
		Published:          a.Published,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}
	if a.InappBackgroundImage.Valid {
		dto.InappBackgroundImage = &a.InappBackgroundImage.String
	}
	if a.Description.Valid {
		dto.Description = &a.Description.String
	}
	if a.Results.Valid {
		dto.Results = &a.Results.String
	}
	if a.Challenges.Valid {
		dto.Challenges = &a.Challenges.String
	}
	if a.Location.Valid {
		dto.Location = &a.Location.String
	}
	if a.VideoCover.Valid {
		dto.VideoCover = &a.VideoCover.String
	}
	if a.Safety.Valid {
		dto.Safety = &a.Safety.String
	}
	if a.VerificationProcessID.Valid {
		s := a.VerificationProcessID.UUID.String()
		dto.VerificationProcessID = &s
	}
	if a.VerificationStatus.Valid {
		dto.VerificationStatus = &a.VerificationStatus.String
	}
	return dto
}
