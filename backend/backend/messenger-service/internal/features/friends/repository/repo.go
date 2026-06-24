// Package repository — raw-SQL access to messenger_friend_requests.
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/friends"
)

type Repository interface {
	// Lookup
	FindBetween(ctx context.Context, a, b uuid.UUID) (*friends.FriendRequest, error)
	GetByID(ctx context.Context, id uuid.UUID) (*friends.FriendRequest, error)

	// Mutations
	Insert(ctx context.Context, requesterID, targetID uuid.UUID) (*friends.FriendRequest, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	DeleteBetween(ctx context.Context, a, b uuid.UUID) error

	// Listing
	ListRequests(ctx context.Context, userID uuid.UUID, direction string, limit, offset int) ([]FriendRequestRow, error)
	ListFriends(ctx context.Context, userID uuid.UUID, limit, offset int) ([]FriendRow, error)
}

// FriendRow — row returned by ListFriends (join with messenger_users).
type FriendRow struct {
	UserID        uuid.UUID
	Username      string
	DisplayName   string
	AvatarMediaID *uuid.UUID
	BecameAt      string // updated_at when ACCEPTED
}

// FriendRequestRow — row returned by ListRequests.
type FriendRequestRow struct {
	RequestID     uuid.UUID
	Direction     string // "incoming" (current user is target) | "outgoing" (requester)
	Status        string
	OtherUserID   uuid.UUID
	OtherUsername string
	OtherDisplay  string
	OtherAvatarID *uuid.UUID
	CreatedAt     string
}

type repo struct{ db *pgxpool.Pool }

func New(db *pgxpool.Pool) Repository { return &repo{db: db} }

func (r *repo) FindBetween(ctx context.Context, a, b uuid.UUID) (*friends.FriendRequest, error) {
	const q = `
		SELECT id, requester_id, target_id, status, created_at, updated_at
		FROM   messenger_friend_requests
		WHERE  (requester_id = $1 AND target_id = $2) OR (requester_id = $2 AND target_id = $1)
		LIMIT  1
	`
	var fr friends.FriendRequest
	err := r.db.QueryRow(ctx, q, a, b).Scan(&fr.ID, &fr.RequesterID, &fr.TargetID, &fr.Status, &fr.CreatedAt, &fr.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("friend req find: %w", err)
	}
	return &fr, nil
}

func (r *repo) GetByID(ctx context.Context, id uuid.UUID) (*friends.FriendRequest, error) {
	const q = `
		SELECT id, requester_id, target_id, status, created_at, updated_at
		FROM   messenger_friend_requests
		WHERE  id = $1
	`
	var fr friends.FriendRequest
	err := r.db.QueryRow(ctx, q, id).Scan(&fr.ID, &fr.RequesterID, &fr.TargetID, &fr.Status, &fr.CreatedAt, &fr.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, friends.ErrRequestNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("friend req get: %w", err)
	}
	return &fr, nil
}

func (r *repo) Insert(ctx context.Context, requesterID, targetID uuid.UUID) (*friends.FriendRequest, error) {
	const q = `
		INSERT INTO messenger_friend_requests (requester_id, target_id, status)
		VALUES ($1, $2, 'PENDING')
		ON CONFLICT (requester_id, target_id) DO UPDATE
		   SET status = 'PENDING', updated_at = NOW()
		RETURNING id, requester_id, target_id, status, created_at, updated_at
	`
	var fr friends.FriendRequest
	if err := r.db.QueryRow(ctx, q, requesterID, targetID).Scan(&fr.ID, &fr.RequesterID, &fr.TargetID, &fr.Status, &fr.CreatedAt, &fr.UpdatedAt); err != nil {
		return nil, fmt.Errorf("friend req insert: %w", err)
	}
	return &fr, nil
}

func (r *repo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	const q = `UPDATE messenger_friend_requests SET status = $1, updated_at = NOW() WHERE id = $2`
	tag, err := r.db.Exec(ctx, q, status, id)
	if err != nil {
		return fmt.Errorf("friend req update: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return friends.ErrRequestNotFound
	}
	return nil
}

func (r *repo) DeleteBetween(ctx context.Context, a, b uuid.UUID) error {
	const q = `
		DELETE FROM messenger_friend_requests
		WHERE (requester_id = $1 AND target_id = $2) OR (requester_id = $2 AND target_id = $1)
	`
	_, err := r.db.Exec(ctx, q, a, b)
	return err
}

func (r *repo) ListFriends(ctx context.Context, userID uuid.UUID, limit, offset int) ([]FriendRow, error) {
	const q = `
		SELECT u.id, u.username, u.display_name, u.avatar_media_id, fr.updated_at::text
		FROM   messenger_friend_requests fr
		JOIN   messenger_users u
		       ON u.id = CASE WHEN fr.requester_id = $1 THEN fr.target_id ELSE fr.requester_id END
		WHERE  (fr.requester_id = $1 OR fr.target_id = $1)
		  AND  fr.status = 'ACCEPTED'
		ORDER BY fr.updated_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list friends: %w", err)
	}
	defer rows.Close()
	out := make([]FriendRow, 0)
	for rows.Next() {
		var f FriendRow
		if err := rows.Scan(&f.UserID, &f.Username, &f.DisplayName, &f.AvatarMediaID, &f.BecameAt); err != nil {
			return nil, fmt.Errorf("scan friend: %w", err)
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func (r *repo) ListRequests(ctx context.Context, userID uuid.UUID, direction string, limit, offset int) ([]FriendRequestRow, error) {
	var where string
	switch direction {
	case "incoming":
		where = "fr.target_id = $1 AND fr.status = 'PENDING'"
	case "outgoing":
		where = "fr.requester_id = $1 AND fr.status = 'PENDING'"
	default:
		where = "(fr.target_id = $1 OR fr.requester_id = $1) AND fr.status = 'PENDING'"
	}
	q := fmt.Sprintf(`
		SELECT fr.id,
		       CASE WHEN fr.target_id = $1 THEN 'incoming' ELSE 'outgoing' END AS direction,
		       fr.status,
		       u.id, u.username, u.display_name, u.avatar_media_id,
		       fr.created_at::text
		FROM   messenger_friend_requests fr
		JOIN   messenger_users u
		       ON u.id = CASE WHEN fr.requester_id = $1 THEN fr.target_id ELSE fr.requester_id END
		WHERE  %s
		ORDER BY fr.created_at DESC
		LIMIT $2 OFFSET $3
	`, where)
	rows, err := r.db.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list requests: %w", err)
	}
	defer rows.Close()
	out := make([]FriendRequestRow, 0)
	for rows.Next() {
		var row FriendRequestRow
		if err := rows.Scan(&row.RequestID, &row.Direction, &row.Status, &row.OtherUserID, &row.OtherUsername, &row.OtherDisplay, &row.OtherAvatarID, &row.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan request: %w", err)
		}
		out = append(out, row)
	}
	return out, rows.Err()
}
