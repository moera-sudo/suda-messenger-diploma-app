package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type lastSeenHiddenRepo struct {
	db *pgxpool.Pool
}

func (r *lastSeenHiddenRepo) Add(ctx context.Context, ownerID, hiddenID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO messenger_last_seen_hidden (owner_id, hidden_id)
		 VALUES ($1, $2) ON CONFLICT (owner_id, hidden_id) DO NOTHING`,
		ownerID, hiddenID,
	)
	return err
}

func (r *lastSeenHiddenRepo) Remove(ctx context.Context, ownerID, hiddenID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM messenger_last_seen_hidden WHERE owner_id = $1 AND hidden_id = $2`,
		ownerID, hiddenID,
	)
	return err
}

func (r *lastSeenHiddenRepo) IsHidden(ctx context.Context, ownerID, hiddenID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT 1 FROM messenger_last_seen_hidden
			WHERE owner_id = $1 AND hidden_id = $2
		)`, ownerID, hiddenID,
	).Scan(&exists)
	return exists, err
}

func (r *lastSeenHiddenRepo) GetAll(ctx context.Context, ownerID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx,
		`SELECT hidden_id FROM messenger_last_seen_hidden WHERE owner_id = $1`,
		ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}