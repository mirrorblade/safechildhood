package repository

import (
	"context"
	"safechildhood/internal/app/domain"
	"safechildhood/internal/app/repository/postgre"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Complaints interface {
	GetEarly(ctx context.Context) ([]*domain.Complaint, error)
	Get(ctx context.Context, complaintId any) (domain.Complaint, error)
	Create(ctx context.Context, complaint domain.Complaint) error
	Delete(ctx context.Context, complaintId any) error
}

type Repository struct {
	Complaints
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		Complaints: postgre.NewComplaints(pool),
	}
}
