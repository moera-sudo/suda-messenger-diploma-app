package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/user"
)

type Repository interface {
	UpdateProfile(ctx context.Context, userID uuid.UUID, displayName, bio, firstName, lastName string) error
	UpdateAvatar(ctx context.Context, userID uuid.UUID, mediaID uuid.UUID) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*user.UserResponse, error)
	UserLogout(ctx context.Context, userID uuid.UUID, refreshToken string) error
	RemoveDeviceToken(ctx context.Context, userID uuid.UUID, fcmToken string) error

	UpdateLastSeen(ctx context.Context, userID uuid.UUID) error
	GetAllInterlocutorIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	GetLastSeenAt(ctx context.Context, userID uuid.UUID) (*time.Time, error)

	GetUserByExactUsername(ctx context.Context, username string) (*user.UserResponse, error)
	GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]user.UserResponse, error)
	GetWalletAddress(ctx context.Context, userID uuid.UUID) (string, error)
	SetWalletAddress(ctx context.Context, userID uuid.UUID, address string) error
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &userRepo{db: db}
}

func (r *userRepo) UpdateProfile(ctx context.Context, userID uuid.UUID, displayName, bio, firstName, lastName string) error {
	query := `
		UPDATE messenger_users
		SET display_name = $1, bio = $2, first_name = $3, last_name = $4
		WHERE id = $5
	`
	_, err := r.db.Exec(ctx, query, displayName, bio, firstName, lastName, userID)
	return err
}

// ФИКС: было &1/&2, стало $1/$2. Теперь сохраняет media_id вместо URL
func (r *userRepo) UpdateAvatar(ctx context.Context, userID uuid.UUID, mediaID uuid.UUID) error {
	query := `
		UPDATE messenger_users
		SET avatar_media_id = $1
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, mediaID, userID)
	return err
}

// ФИКС: добавлены first_name, last_name, avatar_media_id в SELECT
func (r *userRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (*user.UserResponse, error) {
	query := `
		SELECT id, username, display_name,
		       COALESCE(first_name, ''), COALESCE(last_name, ''),
		       COALESCE(bio, ''),
		       avatar_media_id
		FROM messenger_users
		WHERE id = $1
	`
	var u user.UserResponse
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&u.ID, &u.Username, &u.DisplayName,
		&u.FirstName, &u.LastName,
		&u.Bio, &u.AvatarMediaID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &u, err
}

// GetUsersByIDs — batch-резолв юзеров по списку id. Возвращает только
// найденных; отсутствующие id просто не попадают в результат.
func (r *userRepo) GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]user.UserResponse, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	const query = `
		SELECT id, username, display_name
		FROM messenger_users
		WHERE id = ANY($1)
	`
	rows, err := r.db.Query(ctx, query, userIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []user.UserResponse
	for rows.Next() {
		var u user.UserResponse
		if err := rows.Scan(&u.ID, &u.Username, &u.DisplayName); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func (r *userRepo) UserLogout(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	query := `
		UPDATE messenger_sessions
		SET is_revoked = TRUE, revoked_at = NOW()
		WHERE user_id = $1 AND refresh_token = $2
	`
	_, err := r.db.Exec(ctx, query, userID, refreshToken)
	return err
}

func (r *userRepo) RemoveDeviceToken(ctx context.Context, userID uuid.UUID, fcmToken string) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM messenger_user_devices WHERE user_id = $1 AND fcm_token = $2`,
		userID, fcmToken)
	return err
}


func (r *userRepo) UpdateLastSeen(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE messenger_users SET last_seen_at = NOW() WHERE id = $1`, userID)
	return err
}

func (r *userRepo) GetAllInterlocutorIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx,
		`SELECT DISTINCT cm2.user_id
		 FROM messenger_chat_members cm1
		 JOIN messenger_chat_members cm2 ON cm1.chat_id = cm2.chat_id
		 WHERE cm1.user_id = $1 AND cm2.user_id != $1`,
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

func (r *userRepo) GetLastSeenAt(ctx context.Context, userID uuid.UUID) (*time.Time, error) {
	var lastSeen *time.Time
	err := r.db.QueryRow(ctx,
		`SELECT last_seen_at FROM messenger_users WHERE id = $1`, userID,
	).Scan(&lastSeen)
	if err != nil {
		return nil, err
	}
	return lastSeen, nil
}

func (r *userRepo) GetUserByExactUsername(ctx context.Context, username string) (*user.UserResponse, error) {
	// Usernames are persisted with a leading '@'. Callers pass the handle with
	// or without it, so normalize to the stored canonical form before matching.
	handle := "@" + strings.TrimPrefix(strings.TrimSpace(username), "@")

	const query = `
		SELECT
			id, username, display_name,
			COALESCE(first_name, ''), COALESCE(last_name, ''),
			COALESCE(bio, ''), COALESCE(avatar_media_id::text, ''),
			COALESCE(wallet_address, ''),
			is_verified
		FROM messenger_users
		WHERE username = $1
	`

	var u user.UserResponse
	var avatarID string
	err := r.db.QueryRow(ctx, query, handle).Scan(
		&u.ID, &u.Username, &u.DisplayName,
		&u.FirstName, &u.LastName,
		&u.Bio, &avatarID,
		&u.WalletAddress,
		&u.IsVerified,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if avatarID != "" {
		if id, parseErr := uuid.Parse(avatarID); parseErr == nil {
			u.AvatarMediaID = &id
		}
	}

	return &u, nil
}

func (r *userRepo) GetWalletAddress(ctx context.Context, userID uuid.UUID) (string, error) {
	const query = `SELECT COALESCE(wallet_address, '') FROM messenger_users WHERE id = $1`

	var addr string
	err := r.db.QueryRow(ctx, query, userID).Scan(&addr)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", errors.New("user not found")
	}
	if err != nil {
		return "", err
	}
	return addr, nil
}

func (r *userRepo) SetWalletAddress(ctx context.Context, userID uuid.UUID, address string) error {
	const query = `UPDATE messenger_users SET wallet_address = $1 WHERE id = $2`

	tag, err := r.db.Exec(ctx, query, address, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}