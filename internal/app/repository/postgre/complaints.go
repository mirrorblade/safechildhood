package postgre

import (
	"context"
	"errors"
	"fmt"
	"safechildhood/internal/app/domain"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Complaints struct {
	pool *pgxpool.Pool

	table string
}

func NewComplaints(pool *pgxpool.Pool) *Complaints {
	return &Complaints{
		pool:  pool,
		table: "complaints",
	}
}

func (c *Complaints) Get(ctx context.Context, complaintId any) (domain.Complaint, error) {
	startQuery := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", c.table)

	row, err := c.pool.Query(ctx, startQuery, complaintId)
	if err != nil {
		return domain.Complaint{}, err
	}

	complaint, err := pgx.CollectOneRow(row, pgx.RowToStructByName[domain.Complaint])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Complaint{}, domain.ErrUserNotFound
		}

		return domain.Complaint{}, err
	}

	formatedId, ok := complaint.ID.([16]byte)
	if !ok {
		return domain.Complaint{}, err
	}

	complaint.ID, err = uuid.FromBytes(formatedId[:])
	if err != nil {
		return domain.Complaint{}, err
	}

	return complaint, nil
}

func (c *Complaints) GetEarly(ctx context.Context) ([]domain.Complaint, error) {
	startQuery := fmt.Sprintf(`SELECT DISTINCT ON (coordinates) *
	 	FROM %s ORDER BY coordinates, created_at ASC`, c.table)

	rows, err := c.pool.Query(ctx, startQuery)
	if err != nil {
		return []domain.Complaint{}, err
	}

	complaints, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Complaint])
	if err != nil {
		return []domain.Complaint{}, err
	}

	return complaints, nil
}

func (c *Complaints) Create(ctx context.Context, complaint *domain.Complaint) error {
	startQuery := fmt.Sprintf(`INSERT INTO %s (
		id, 
		coordinates, 
		short_description,
		description,
		photos_path,
		created_at
	) VALUES ($1, $2, $3, $4, $5, $6)`, c.table)
	_, err := c.pool.Exec(ctx, startQuery,
		complaint.ID,
		complaint.Coordinates,
		complaint.ShortDescription,
		complaint.Description,
		complaint.PhotosPath,
		complaint.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (c *Complaints) Delete(ctx context.Context, complaintId any) error {
	startQuery := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, c.table)

	commandTag, err := c.pool.Exec(ctx, startQuery, complaintId)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != -1 {
		return domain.ErrUserNotFound
	}

	return nil
}
