package service

import (
	"context"
	"safechildhood/internal/app/domain"
	"safechildhood/internal/app/repository"
)

type ComplaintsService struct {
	repo repository.Complaints

	playgrounds Playgrounds
}

func NewComplaintsService(repo repository.Complaints, playgrounds Playgrounds) *ComplaintsService {
	return &ComplaintsService{
		repo:        repo,
		playgrounds: playgrounds,
	}
}

func (c *ComplaintsService) Get(ctx context.Context, complaintId any) (domain.Complaint, error) {
	return c.repo.Get(ctx, complaintId)
}

func (c *ComplaintsService) Create(ctx context.Context, complaint *domain.Complaint) error {
	if err := c.repo.Create(ctx, complaint); err != nil {
		return err
	}

	complaints, err := c.repo.GetEarly(ctx)
	if err != nil {
		return err
	}

	c.playgrounds.UpdatePlaygroundsMap(complaints)

	return nil
}

func (c *ComplaintsService) Delete(ctx context.Context, complaintId any) error {
	if err := c.repo.Delete(ctx, complaintId); err != nil {
		return err
	}

	complaints, err := c.repo.GetEarly(ctx)
	if err != nil {
		return err
	}

	c.playgrounds.UpdatePlaygroundsMap(complaints)

	return nil
}
