package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/userpins"
)

type repo struct {
	db *pgxpool.Pool
}

func (r *repo) Create(ctx context.Context, pin *userpins.UserPin) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO messenger_user_pins
			(user_id, pin_type, target_type, target_id, sort_order)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (user_id, pin_type, target_id) DO UPDATE
		     SET sort_order = EXCLUDED.sort_order
		 RETURNING id, created_at`,
		pin.UserID, pin.PinType, pin.TargetType, pin.TargetID, pin.SortOrder,
	).Scan(&pin.ID, &pin.CreatedAt)
}

func (r *repo) Delete(ctx context.Context, userID uuid.UUID, pinID int64) error {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM messenger_user_pins WHERE id = $1 AND user_id = $2`,
		pinID, userID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("pin not found")
	}
	return nil
}

func (r *repo) DeleteByTarget(ctx context.Context, userID, targetID uuid.UUID, pinType string) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM messenger_user_pins
		 WHERE user_id = $1 AND target_id = $2 AND pin_type = $3`,
		userID, targetID, pinType,
	)
	return err
}

func (r *repo) GetByUser(ctx context.Context, userID uuid.UUID, pinType string) ([]userpins.UserPin, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, pin_type, target_type, target_id, sort_order, created_at
		 FROM messenger_user_pins
		 WHERE user_id = $1 AND pin_type = $2
		 ORDER BY sort_order ASC, created_at ASC`,
		userID, pinType,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pins []userpins.UserPin
	for rows.Next() {
		var p userpins.UserPin
		if err := rows.Scan(&p.ID, &p.UserID, &p.PinType, &p.TargetType,
			&p.TargetID, &p.SortOrder, &p.CreatedAt); err != nil {
			return nil, err
		}
		pins = append(pins, p)
	}
	return pins, rows.Err()
}

// GetPinnedChatIDs — для GetUserChats, чтобы знать какие чаты закреплены в CHATLIST
func (r *repo) GetPinnedChatIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx,
		`SELECT target_id FROM messenger_user_pins
		 WHERE user_id = $1 AND pin_type = 'CHATLIST' AND target_type = 'CHAT'
		 ORDER BY sort_order ASC`,
		userID,
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

// GetPinnedChats — закреплённые чаты (CHATLIST) с временем закрепления,
// отсортированные по дате закрепления (новые сверху) — для GetUserChats.
func (r *repo) GetPinnedChats(ctx context.Context, userID uuid.UUID) ([]userpins.PinnedChat, error) {
	rows, err := r.db.Query(ctx,
		`SELECT target_id, created_at FROM messenger_user_pins
		 WHERE user_id = $1 AND pin_type = 'CHATLIST' AND target_type = 'CHAT'
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []userpins.PinnedChat
	for rows.Next() {
		var p userpins.PinnedChat
		if err := rows.Scan(&p.ChatID, &p.PinnedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *repo) Reorder(ctx context.Context, userID uuid.UUID, items []userpins.ReorderItem) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, item := range items {
		_, err := tx.Exec(ctx,
			`UPDATE messenger_user_pins
			 SET sort_order = $1
			 WHERE id = $2 AND user_id = $3`,
			item.SortOrder, item.PinID, userID,
		)
		if err != nil {
			return fmt.Errorf("update pin %d: %w", item.PinID, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *repo) GetMaxSortOrder(ctx context.Context, userID uuid.UUID, pinType string) (int, error) {
	var max *int
	err := r.db.QueryRow(ctx,
		`SELECT MAX(sort_order) FROM messenger_user_pins
		 WHERE user_id = $1 AND pin_type = $2`,
		userID, pinType,
	).Scan(&max)
	if err != nil {
		return 0, err
	}
	if max == nil {
		return 0, nil
	}
	return *max, nil
}