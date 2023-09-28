package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type ObjectMode int

const (
	FILE ObjectMode = iota + 1
	FOLDER
)

const (
	GOOGLE_DRIVE_MAIN_FOLDER = "root"

	GOOGLE_DRIVE_FOLDER_MIMETYPE = "application/vnd.google-apps.folder"
	GOOGLE_DRIVE_FILE_MIMETYPE   = "application/octet-stream"
)

var (
	ErrInvalidParams       = errors.New("enterned function parameters are invalid")
	ErrInvalidObjectMode   = errors.New("object mode is invalid")
	ErrObjectAlreadyExists = errors.New("object with same name already exists")
)

type GoogleDriveParameters struct {
	Name                   string
	ObjectMode             ObjectMode
	Content                io.Reader
	ParentId               string
	SkipAlreadyExistsCheck bool
}

type GoogleDrive struct {
	service *drive.Service
}

func NewGoogleDrive(ctx context.Context, pathToCredentials string) (*GoogleDrive, error) {
	service, err := drive.NewService(ctx, option.WithCredentialsFile(pathToCredentials))
	if err != nil {
		return &GoogleDrive{}, err
	}

	return &GoogleDrive{
		service: service,
	}, nil
}

func (g *GoogleDrive) Get(ctx context.Context, objectId string) (*drive.File, error) {
	return g.service.Files.Get(objectId).Context(ctx).Do()
}

func (g *GoogleDrive) GetByParams(ctx context.Context, params GoogleDriveParameters) ([]*drive.File, error) {
	if params.Name == "" && params.ObjectMode == 0 && params.ParentId == "" {
		return []*drive.File{}, ErrInvalidParams
	}

	querySlice := make([]string, 0, 3)

	if params.Name != "" {
		querySlice = append(querySlice, fmt.Sprintf(`name="%s"`, params.Name))
	}

	if params.ObjectMode == FILE {
		querySlice = append(querySlice, fmt.Sprintf(`mimeType="%s"`, GOOGLE_DRIVE_FILE_MIMETYPE))
	} else if params.ObjectMode == FOLDER {
		querySlice = append(querySlice, fmt.Sprintf(`mimeType="%s"`, GOOGLE_DRIVE_FOLDER_MIMETYPE))
	}

	if params.ParentId != "" {
		querySlice = append(querySlice, fmt.Sprintf(`"%s" in parents`, params.ParentId))
	}

	list, err := g.service.Files.
		List().
		Q(
			strings.Join(querySlice, " and "),
		).
		Context(ctx).
		Do()
	if err != nil {
		return []*drive.File{}, err
	}

	return list.Files, nil

}

func (g *GoogleDrive) Create(ctx context.Context, params GoogleDriveParameters) (*drive.File, error) {
	var object *drive.File

	if params.Name == "" || params.ParentId == "" || params.ObjectMode == 0 {
		return &drive.File{}, ErrInvalidParams
	}

	switch params.ObjectMode {
	case FILE:
		fileShell := &drive.File{
			MimeType: "application/octet-stream",
			Name:     params.Name,
			Parents:  []string{params.ParentId},
		}

		if !params.SkipAlreadyExistsCheck {
			list, _ := g.service.Files.
				List().
				Q(
					fmt.Sprintf(`"%s" in parents and name="%s" and mimeType="%s"`, params.ParentId, params.Name, GOOGLE_DRIVE_FILE_MIMETYPE),
				).
				Context(ctx).
				Do()

			if len(list.Files) != 0 {
				return &drive.File{}, ErrObjectAlreadyExists
			}
		}

		file, err := g.service.Files.Create(fileShell).Media(params.Content).Context(ctx).Do()
		if err != nil {
			return &drive.File{}, err
		}

		object = file

	case FOLDER:
		folderShell := &drive.File{
			Name:     params.Name,
			MimeType: GOOGLE_DRIVE_FOLDER_MIMETYPE,
			Parents:  []string{params.ParentId},
		}

		if !params.SkipAlreadyExistsCheck {
			list, _ := g.service.Files.
				List().
				Q(
					fmt.Sprintf(`"%s" in parents and name="%s" and mimeType="%s"`, params.ParentId, params.Name, GOOGLE_DRIVE_FOLDER_MIMETYPE),
				).
				Context(ctx).
				Do()

			if len(list.Files) != 0 {
				return &drive.File{}, ErrObjectAlreadyExists
			}
		}

		folder, err := g.service.Files.Create(folderShell).Context(ctx).Do()
		if err != nil {
			return &drive.File{}, err
		}

		object = folder

	default:
		return &drive.File{}, ErrInvalidObjectMode
	}

	return object, nil
}

func (g *GoogleDrive) Rename(ctx context.Context, changedName, objectId string) (*drive.File, error) {
	return g.service.Files.
		Update(objectId, &drive.File{Name: changedName}).
		Context(ctx).
		Do()
}

func (g *GoogleDrive) Copy(ctx context.Context, copiedFromId string, copiedTo *drive.File) (*drive.File, error) {
	return g.service.Files.Copy(copiedFromId, copiedTo).Context(ctx).Do()
}

func (g *GoogleDrive) Delete(ctx context.Context, objectId string) error {
	return g.service.Files.Delete(objectId).Context(ctx).Do()
}
