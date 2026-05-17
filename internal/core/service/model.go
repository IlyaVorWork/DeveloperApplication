package service

import (
	"errors"
	"io"
)

var ErrNotVerified = errors.New("application must be successfully verified before publishing")

type UploadBuildParams struct {
	DeveloperID          string
	CodeName             string
	CategoryID           int64
	AndroidPackageName   string
	DefaultLocale        string
	Name                 string
	ShortTitle           string
	Description          string
	Goals                string
	Tasks                string
	Results              string
	Challenges           string
	Location             string
	VideoCover           string
	Safety               string
	WebVideo             string
	InappVideo           string
	WebBackgroundImage   string
	InappBackgroundImage string
	Version              string
	FileName             string
	Reader               io.Reader
	Size                 int64
}
