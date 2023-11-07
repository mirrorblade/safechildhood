package service

import (
	"context"
	"safechildhood/internal/app/domain"
	"safechildhood/pkg/storage"

	geojson "github.com/paulmach/go.geojson"
	"google.golang.org/api/drive/v3"
)

type Complaints interface {
	GetEarly(ctx context.Context) ([]*domain.Complaint, error)
	Get(ctx context.Context, complaintId any) (domain.Complaint, error)
	Create(ctx context.Context, complaint domain.Complaint) error
	Delete(ctx context.Context, complaintId any) error
}

type Playgrounds interface {
	GetPlaygrounds() *geojson.FeatureCollection
	Refresh() chan any
	SetPlaygroundsMap(playgroundsMap map[string]*MapProperties)
	UpdatePlaygroundsMap(playgroundsMap map[string]*MapProperties)
}

type Storage interface {
	SetFoldersIdsMap(foldersIdsMap map[string]string)
	GetSavedFolderId(key string) string
	GetByParams(ctx context.Context, params storage.GoogleDriveParameters) ([]*drive.File, error)
	Create(ctx context.Context, params storage.GoogleDriveParameters) (*drive.File, error)
}

type Service struct {
	Playgrounds
	Complaints
	Storage
}

func New(c Complaints, p Playgrounds, s Storage) *Service {
	return &Service{
		Complaints:  c,
		Playgrounds: p,
		Storage:     s,
	}
}
