package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/auth"
)

type Repository interface {
	CreateUser(ctx context.Context, user *auth.User) (uuid.UUID, error)
	CreateVerification(ctx context.Context, email, codeHash, purpose string) error
	GetVerification(ctx context.Context, email, purpose string) (string, uuid.UUID, error)
	MarkUserVerified(ctx context.Context, userID uuid.UUID, email string) error
	UseVerificationCode(ctx context.Context, email string, code_hash string) error
	GetUserByEmail(ctx context.Context, email string) (*auth.User, error)
	SaveSession(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, refreshToken string, expiresAt time.Time, user_agent string, client_ip string) error
	RefreshSession(ctx context.Context, sessionID uuid.UUID, oldRefreshToken string, newRefreshToken string, expiresAt time.Time, userAgent string, clientIP string) error
	UpdatePassword(ctx context.Context, email string, passwordHash string) error
	UpdatePasswordByUserID(ctx context.Context, userID uuid.UUID, passwordHash string) error
	GetUserByID(ctx context.Context, user_id uuid.UUID) (*auth.User, error)
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error

	// Session management (Active Sessions feature)
	ListActiveSessions(ctx context.Context, userID uuid.UUID) ([]SessionRow, error)
	RevokeSessionByID(ctx context.Context, userID, sessionID uuid.UUID) error
	RevokeSessionsExceptToken(ctx context.Context, userID uuid.UUID, keepToken string) (int64, error)
}

// SessionRow — internal session row used by ListActiveSessions / Active Sessions API.
// Maps to messenger_sessions; device_name is best-effort join with messenger_user_devices.
type SessionRow struct {
	ID         uuid.UUID
	UserAgent  string
	ClientIP   string
	CreatedAt  time.Time
	LastUsedAt *time.Time
	DeviceName string
}



type authRepo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository{
	return &authRepo{db: db}
}

func (r *authRepo) CreateUser(ctx context.Context, user *auth.User) (uuid.UUID,  error) {
	query := `
		INSERT INTO messenger_users (username, email, password_hash, display_name, role, is_verified)
		VALUES ($1, $2, $3, $4, 'USER', FALSE)
		RETURNING id
	`

	var id uuid.UUID
	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.PasswordHash, user.DisplayName).Scan(&id)
	if err == nil{
		return id, nil
	}

	// если юзер уже существует
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		if pgErr.ConstraintName == "messenger_users_username_key" {
			return r.handleConflict(ctx, user)
		}

		if pgErr.ConstraintName == "messenger_users_email_key" {
			return r.handleConflict(ctx, user)
		}
		

	}
	return uuid.Nil, err

}

func (r *authRepo) handleConflict(ctx context.Context, user *auth.User) (uuid.UUID, error) {
	var existingID uuid.UUID
	var isVerified bool
	var username string
	var email string

	err := r.db.QueryRow(ctx, `SELECT id, is_verified, username, email FROM messenger_users WHERE email = $1 or username = $2`, user.Email, user.Username).Scan(&existingID, &isVerified, &username, &email)

	if err != nil {
		return uuid.Nil, err
	}

	if isVerified && email == user.Email && username == user.Username {
		return uuid.Nil, errors.New("Email and username are already taken")
	}
	if isVerified && email == user.Email {
		return uuid.Nil, errors.New("Email is already taken")
	}
	if isVerified && username == user.Username {
		return uuid.Nil, errors.New("Username is already taken")
	}


	updateQuery := `
		UPDATE messenger_users
		SET username = $2, password_hash = $3, display_name = $4, created_at = NOW()
		WHERE id = $1
		RETURNING id
	`

	var updatedID uuid.UUID
	err = r.db.QueryRow(ctx, updateQuery, existingID, user.Username, user.PasswordHash, user.DisplayName).Scan(&updatedID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "messenger_users_username_key" {
			return uuid.Nil, errors.New("Username is already taken")
		}
		return uuid.Nil, err
	}

	return updatedID, nil
}

func (r *authRepo) CreateVerification(ctx context.Context, email, codeHash, purpose string) error {
	query := `
		INSERT INTO messenger_verifications (user_id, email, code_hash, purpose, expires_at)
		VALUES ((SELECT id FROM messenger_users WHERE email = $1), $1, $2, $3, NOW() + INTERVAL '15 minutes')
	`

	_, err := r.db.Exec(ctx, query, email, codeHash, purpose)
	return err
}

func (r *authRepo) GetVerification(ctx context.Context, email, purpose string) (string, uuid.UUID, error) {
	query := `
		SELECT code_hash, user_id
		FROM messenger_verifications
		WHERE email = $1 AND purpose = $2 AND used_at IS NULL AND expires_at > NOW()
		ORDER BY created_at DESC LIMIT 1
	`
	var hash string
	var userID uuid.UUID
	err := r.db.QueryRow(ctx, query, email, purpose).Scan(&hash, &userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", uuid.Nil, errors.New("Code not found or expired")
	}
	return hash, userID, err
}

func (r *authRepo) MarkUserVerified(ctx context.Context, userID uuid.UUID, email string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `UPDATE messenger_users SET is_verified = TRUE WHERE id = $1`, userID); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `UPDATE messenger_verifications SET used_at = NOW() WHERE email = $1 AND used_at IS NULL`, email); err != nil {
		return err
	}

	return tx.Commit(ctx)

}

func (r *authRepo) UseVerificationCode(ctx context.Context, email string, code_hash string) error {
	query := `UPDATE messenger_verifications SET used_at = NOW() WHERE email = $1 AND used_at IS NULL AND code_hash = $2`
	_, err := r.db.Exec(ctx, query, email, code_hash)
	return err
}
 
func (r *authRepo) GetUserByEmail(ctx context.Context, email string) (*auth.User, error) {
	query := `SELECT id, username, email, password_hash, display_name, is_verified, role FROM messenger_users WHERE email = $1`
	var u auth.User
	err := r.db.QueryRow(ctx, query, email).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.DisplayName, &u.IsVerified, &u.Role)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return &u, err
}

func (r *authRepo) GetUserByID(ctx context.Context, user_id uuid.UUID) (*auth.User, error) {
	query := `SELECT id, username, email, password_hash, display_name, is_verified, role, avatar_media_id FROM messenger_users WHERE id = $1`
	var u auth.User
	err := r.db.QueryRow(ctx, query, user_id).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.DisplayName, &u.IsVerified, &u.Role, &u.AvatarMediaID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return &u, err
}

func (r *authRepo) SaveSession(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, refreshToken string, expiresAt time.Time, userAgent string, clientIP string) error {
	_, err := r.db.Exec(ctx, `INSERT INTO messenger_sessions (id, user_id, refresh_token, expires_at, user_agent, client_ip) VALUES ($1, $2, $3, $4, $5, $6)`, sessionID, userID, refreshToken, expiresAt, userAgent, clientIP)
	return err
}

func (r *authRepo) RefreshSession(ctx context.Context, sessionID uuid.UUID, oldRefreshToken, newRefreshToken string, expiresAt time.Time, userAgent, clientIP string) error {
    query := `
        WITH revoked_session AS (
            UPDATE messenger_sessions
            SET
                is_revoked = TRUE,
                revoked_at = NOW()
            WHERE refresh_token = $1
              AND is_revoked = FALSE
            RETURNING user_id -- Возвращаем ID юзера, чей токен мы отозвали
        )
        INSERT INTO messenger_sessions (
            id,
            user_id,
            refresh_token,
            expires_at,
            user_agent,
            client_ip
        )
        SELECT
            $6, -- sessionID (matches the "sid" claim in the new tokens)
            user_id,
            $2, -- newRefreshToken
            $3, -- expiresAt
            $4, -- userAgent
            $5  -- clientIP
        FROM revoked_session -- Если UPDATE ничего не нашел, SELECT будет пустым, и INSERT не сработает
        RETURNING id;
    `

    var newID uuid.UUID
    // Используем QueryRow. Если старый токен не найден, INSERT не выполнится,
    // и Scan вернет ошибку pgx.ErrNoRows.
    err := r.db.QueryRow(ctx, query, oldRefreshToken, newRefreshToken, expiresAt, userAgent, clientIP, sessionID).Scan(&newID)

    if errors.Is(err, pgx.ErrNoRows) {
        // Это значит: либо токен не найден, либо уже отозван
        return errors.New("Invalid or revoked session")
    }
    
    return err
}

// UpdatePasswordByUserID — updates password by user id (used by /user/password change flow).
func (r *authRepo) UpdatePasswordByUserID(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	query := `UPDATE messenger_users SET password_hash = $1, updated_at = now() WHERE id = $2`
	if _, err := r.db.Exec(ctx, query, passwordHash, userID); err != nil {
		return fmt.Errorf("update password by user_id: %w", err)
	}
	return nil
}

func (r *authRepo) UpdatePassword(ctx context.Context, email string, passwordHash string) error {
	query := `UPDATE messenger_users SET password_hash = $1, updated_at = now() WHERE email = $2`
	tag, err := r.db.Exec(ctx, query, passwordHash, email)

	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("User not found")
	}
	return nil
}

func (r *authRepo) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE messenger_sessions
		SET is_revoked = TRUE, revoked_at = NOW()
		WHERE user_id = $1 AND is_revoked = FALSE
	`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

// ListActiveSessions — все не-revoked, не-просроченные сессии. device_name
// подтягивается best-effort из messenger_user_devices по user_id (последняя по
// last_used_at) — это не строгая 1:1 связь сессии и устройства, но достаточно
// для UI «активные сессии».
func (r *authRepo) ListActiveSessions(ctx context.Context, userID uuid.UUID) ([]SessionRow, error) {
	const q = `
		SELECT s.id,
		       COALESCE(s.user_agent, '') AS user_agent,
		       COALESCE(s.client_ip,  '') AS client_ip,
		       s.created_at,
		       -- На refresh строка пересоздаётся, поэтому created_at ≈ время
		       -- последней активности сессии. Используем его как last_used_at.
		       s.created_at              AS last_used_at,
		       COALESCE(d.device_name, '') AS device_name
		FROM   messenger_sessions s
		LEFT JOIN LATERAL (
		           SELECT device_name FROM messenger_user_devices
		           WHERE user_id = s.user_id
		           ORDER BY last_used_at DESC
		           LIMIT 1
		       ) d ON TRUE
		WHERE  s.user_id = $1
		  AND  s.is_revoked = FALSE
		  AND  s.expires_at > NOW()
		ORDER BY s.created_at DESC
	`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}
	defer rows.Close()
	out := make([]SessionRow, 0)
	for rows.Next() {
		var sr SessionRow
		if err := rows.Scan(&sr.ID, &sr.UserAgent, &sr.ClientIP, &sr.CreatedAt, &sr.LastUsedAt, &sr.DeviceName); err != nil {
			return nil, fmt.Errorf("scan session row: %w", err)
		}
		out = append(out, sr)
	}
	return out, rows.Err()
}

// RevokeSessionByID — soft-revoke single session, only if it belongs to userID.
func (r *authRepo) RevokeSessionByID(ctx context.Context, userID, sessionID uuid.UUID) error {
	const q = `
		UPDATE messenger_sessions
		SET    is_revoked = TRUE, revoked_at = NOW()
		WHERE  id = $1 AND user_id = $2 AND is_revoked = FALSE
	`
	tag, err := r.db.Exec(ctx, q, sessionID, userID)
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errors.New("session not found")
	}
	return nil
}

// RevokeSessionsExceptToken — revoke every non-revoked session for userID except
// the one that holds `keepToken` as refresh_token. Returns number revoked.
func (r *authRepo) RevokeSessionsExceptToken(ctx context.Context, userID uuid.UUID, keepToken string) (int64, error) {
	const q = `
		UPDATE messenger_sessions
		SET    is_revoked = TRUE, revoked_at = NOW()
		WHERE  user_id = $1 AND is_revoked = FALSE AND refresh_token <> $2
	`
	tag, err := r.db.Exec(ctx, q, userID, keepToken)
	if err != nil {
		return 0, fmt.Errorf("revoke others: %w", err)
	}
	return tag.RowsAffected(), nil
}