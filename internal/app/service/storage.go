package service

import (
	"context"
	"safechildhood/pkg/storage"

	"google.golang.org/api/drive/v3"
)

type StorageService struct {
	service *storage.GoogleDrive

	foldersIdsMap map[string]string
}

func NewStorageService(service *storage.GoogleDrive) (*StorageService, error) {
	return &StorageService{
		service: service,
	}, nil
}

func (s *StorageService) SetFoldersIdsMap(foldersIdsMap map[string]string) {
	s.foldersIdsMap = make(map[string]string)

	for k, v := range foldersIdsMap {
		s.foldersIdsMap[k] = v
	}
}

func (s *StorageService) GetSavedFolderId(key string) string {
	return s.foldersIdsMap[key]
}

func (s *StorageService) GetByParams(ctx context.Context, params storage.GoogleDriveParameters) ([]*drive.File, error) {
	return s.service.GetByParams(ctx, params)
}

func (s *StorageService) Create(ctx context.Context, params storage.GoogleDriveParameters) (*drive.File, error) {
	return s.service.Create(ctx, params)
}
