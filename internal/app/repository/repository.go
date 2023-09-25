package repository

import (
	"context"
	"safechildhood/internal/app/domain"
	"safechildhood/internal/app/repository/postgre"

	"github.com/jackc/pgx/v5"
)

type Complaints interface {
	Get(ctx context.Context, complaintId any) (domain.Complaint, error)
	GetEarly(ctx context.Context) ([]domain.Complaint, error)
	Create(ctx context.Context, complaint *domain.Complaint) error
	Delete(ctx context.Context, complaintId any) error
}

type Repository struct {
	Complaints
}

func New(conn *pgx.Conn) *Repository {
	return &Repository{
		Complaints: postgre.NewComplaints(conn),
	}
}
