package service

import (
	"context"
	"safechildhood/internal/app/domain"
	"safechildhood/internal/app/repository"
)

type ComplaintsService struct {
	repo repository.Complaints
}

func NewComplaintsService(repo repository.Complaints) *ComplaintsService {
	return &ComplaintsService{
		repo: repo,
	}
}

func (c *ComplaintsService) GetEarly(ctx context.Context) ([]*domain.Complaint, error) {
	return c.repo.GetEarly(ctx)
}

func (c *ComplaintsService) Get(ctx context.Context, complaintId any) (domain.Complaint, error) {
	return c.repo.Get(ctx, complaintId)
}

func (c *ComplaintsService) Create(ctx context.Context, complaint domain.Complaint) error {
	return c.repo.Create(ctx, complaint)
}

func (c *ComplaintsService) Delete(ctx context.Context, complaintId any) error {
	return c.repo.Delete(ctx, complaintId)
}
