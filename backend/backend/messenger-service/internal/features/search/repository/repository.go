package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/search"
	appcrypto "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/crypto"
)

type Repository interface {
	SearchUsers(ctx context.Context, query string, limit int) ([]search.SearchResult, error)
	SearchChats(ctx context.Context, query string, limit int) ([]search.SearchResult, error)
	SearchMessagesGlobal(ctx context.Context, userID uuid.UUID, query string, limit int) ([]search.SearchResult, error)
	SearchMessagesInChat(ctx context.Context, chatID uuid.UUID, query string, limit int) ([]search.SearchResult, error)
}

type searchRepo struct {
	db     *pgxpool.Pool
	cipher appcrypto.ContentCipher
}

func NewRepository(db *pgxpool.Pool, cipher appcrypto.ContentCipher) Repository {
	return &searchRepo{db: db, cipher: cipher}
}

// toTsQuery — преобразует пользовательский ввод в tsquery
// "привет мир" → "привет:* & мир:*"
func toTsQuery(query string) string {
	// Простой подход: разбиваем по пробелам, добавляем prefix matching
	words := strings.Fields(query)
	if len(words) == 0 {
		return ""
	}

	parts := make([]string, len(words))
	for i, w := range words {
		// Экранируем спецсимволы tsquery
		w = strings.ReplaceAll(w, "'", "")
		w = strings.ReplaceAll(w, "\\", "")
		parts[i] = w + ":*"
	}
	return strings.Join(parts, " & ")
}

func (r *searchRepo) SearchUsers(ctx context.Context, query string, limit int) ([]search.SearchResult, error) {
	tsq := toTsQuery(query)
	if tsq == "" {
		return nil, nil
	}

	rows, err := r.db.Query(ctx,
		`SELECT
			id::text, username, display_name,
			COALESCE(avatar_media_id::text, ''),
			ts_rank(search_vector, to_tsquery('simple', $1)) as rank
		 FROM messenger_users
		 WHERE search_vector @@ to_tsquery('simple', $1) AND is_verified = TRUE
		 ORDER BY rank DESC
		 LIMIT $2`,
		tsq, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []search.SearchResult
	for rows.Next() {
		var r search.SearchResult
		var username, displayName, avatar string
		if err := rows.Scan(&r.ID, &username, &displayName, &avatar, &r.Rank); err != nil {
			return nil, err
		}
		r.Type = search.ResultTypeUser
		r.Title = displayName
		r.Subtitle = username
		r.ImageURL = avatar
		results = append(results, r)
	}
	return results, rows.Err()
}

func (r *searchRepo) SearchChats(ctx context.Context, query string, limit int) ([]search.SearchResult, error) {
	tsq := toTsQuery(query)
	// Exact @username lets PRIVATE channels be discovered (so a user can request to
	// join) without exposing them through the full-text index, which only contains
	// public channels and groups.
	exact := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(query), "@")))
	if tsq == "" && exact == "" {
		return nil, nil
	}
	ftsParam := tsq
	if ftsParam == "" {
		ftsParam = "____no_match____" // valid single token that matches nothing
	}

	rows, err := r.db.Query(ctx,
		`SELECT
			id::text, type, COALESCE(name, ''), COALESCE(description, ''),
			COALESCE(avatar_media_id::text, ''),
			COALESCE(ts_rank(search_vector, to_tsquery('simple', $1)), 0)
				+ CASE WHEN username IS NOT NULL AND lower(username) = $3 THEN 1.0 ELSE 0 END as rank
		 FROM messenger_chats
		 WHERE type IN ('GROUP', 'CHANNEL')
		   AND (
		        search_vector @@ to_tsquery('simple', $1)
		        OR (type = 'CHANNEL' AND username IS NOT NULL AND lower(username) = $3)
		   )
		 ORDER BY rank DESC
		 LIMIT $2`,
		ftsParam, limit, exact,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []search.SearchResult
	for rows.Next() {
		var r search.SearchResult
		var chatType, name, desc, avatar string
		if err := rows.Scan(&r.ID, &chatType, &name, &desc, &avatar, &r.Rank); err != nil {
			return nil, err
		}
		r.Type = search.ResultTypeGroup
		if chatType == "CHANNEL" {
			r.Type = search.ResultTypeChannel
		}
		r.Title = name
		r.Subtitle = desc
		r.ImageURL = avatar
		results = append(results, r)
	}
	return results, rows.Err()
}

func (r *searchRepo) SearchMessagesGlobal(ctx context.Context, userID uuid.UUID, query string, limit int) ([]search.SearchResult, error) {
	// Capture the cipher: inside the row loop `r` is shadowed by the result var.
	cipher := r.cipher
	tsq := toTsQuery(query)
	if tsq == "" {
		return nil, nil
	}

	// ts_headline is gone: content is encrypted at rest, so the snippet/highlight
	// is built in Go from the decrypted text. Matching/ranking still use
	// search_vector (built from plaintext at write time), so they're unchanged.
	rows, err := r.db.Query(ctx,
		`SELECT
			m.id, m.chat_id, m.content,
			COALESCE(u.display_name, 'System'),
			m.created_at,
			ts_rank(m.search_vector, to_tsquery('russian', $1)) as rank
		 FROM messenger_messages m
		 JOIN messenger_chat_members cm ON cm.chat_id = m.chat_id AND cm.user_id = $2
		 LEFT JOIN messenger_users u ON m.sender_id = u.id
		 WHERE m.search_vector @@ (to_tsquery('russian', $1) || to_tsquery('simple', $1))
		   AND m.deleted_at IS NULL
		   AND m.type != 'SYSTEM'
		   AND NOT EXISTS (
		       SELECT 1 FROM messenger_message_hidden h
		       WHERE h.message_id = m.id AND h.user_id = $2
		   )
		 ORDER BY rank DESC, m.created_at DESC
		 LIMIT $3`,
		tsq, userID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []search.SearchResult
	for rows.Next() {
		var r search.SearchResult
		var msgID int64
		var chatID uuid.UUID
		var content, senderName string
		var createdAt time.Time

		if err := rows.Scan(&msgID, &chatID, &content, &senderName, &createdAt, &r.Rank); err != nil {
			return nil, err
		}

		content = cipher.DecryptContent(content)
		r.Type = search.ResultTypeMessage
		r.ID = fmt.Sprintf("%d", msgID)
		r.MessageID = &msgID
		r.ChatID = &chatID
		r.Title = senderName
		r.Subtitle = content
		r.Highlight = buildHighlight(content, query)
		r.CreatedAt = &createdAt

		results = append(results, r)
	}
	return results, rows.Err()
}

func (r *searchRepo) SearchMessagesInChat(ctx context.Context, chatID uuid.UUID, query string, limit int) ([]search.SearchResult, error) {
	// Capture the cipher: inside the row loop `r` is shadowed by the result var.
	cipher := r.cipher
	tsq := toTsQuery(query)
	if tsq == "" {
		return nil, nil
	}

	// ts_headline dropped (content encrypted at rest) — snippet built in Go.
	// search_vector matching/ranking unchanged.
	rows, err := r.db.Query(ctx,
		`SELECT
			m.id, m.content,
			COALESCE(u.display_name, 'System'),
			m.created_at,
			ts_rank(m.search_vector, to_tsquery('russian', $1)) as rank
		 FROM messenger_messages m
		 LEFT JOIN messenger_users u ON m.sender_id = u.id
		 WHERE m.chat_id = $2
		   AND m.search_vector @@ to_tsquery('russian', $1)
		   AND m.deleted_at IS NULL
		   AND m.type != 'SYSTEM'
		 ORDER BY rank DESC, m.created_at DESC
		 LIMIT $3`,
		tsq, chatID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []search.SearchResult
	for rows.Next() {
		var r search.SearchResult
		var msgID int64
		var content, senderName string
		var createdAt time.Time

		if err := rows.Scan(&msgID, &content, &senderName, &createdAt, &r.Rank); err != nil {
			return nil, err
		}

		content = cipher.DecryptContent(content)
		r.Type = search.ResultTypeMessage
		r.ID = fmt.Sprintf("%d", msgID)
		r.MessageID = &msgID
		r.ChatID = &chatID
		r.Title = senderName
		r.Subtitle = content
		r.Highlight = buildHighlight(content, query)
		r.CreatedAt = &createdAt

		results = append(results, r)
	}
	return results, rows.Err()
}

// buildHighlight builds an HTML snippet (with <b>…</b> around query matches) from
// already-decrypted content. It replaces SQL ts_headline, which can't run on the
// now-encrypted content column. Best-effort: case-insensitive substring matching.
func buildHighlight(content, query string) string {
	content = strings.TrimSpace(content)
	if content == "" {
		return ""
	}

	const maxRunes = 200
	runes := []rune(content)
	truncated := false
	if len(runes) > maxRunes {
		runes = runes[:maxRunes]
		truncated = true
	}
	snippet := string(runes)

	for _, w := range strings.Fields(query) {
		w = strings.Trim(w, "*:'\"\\")
		if len([]rune(w)) >= 2 { // skip 1-char tokens to avoid mangling <b> tags
			snippet = highlightToken(snippet, w)
		}
	}
	if truncated {
		snippet += "…"
	}
	return snippet
}

// highlightToken wraps every case-insensitive occurrence of token in <b>…</b>,
// keeping the original casing of the matched text.
func highlightToken(s, token string) string {
	lowerS := strings.ToLower(s)
	lowerT := strings.ToLower(token)
	var b strings.Builder
	for {
		idx := strings.Index(lowerS, lowerT)
		if idx < 0 {
			b.WriteString(s)
			break
		}
		b.WriteString(s[:idx])
		b.WriteString("<b>")
		b.WriteString(s[idx : idx+len(token)])
		b.WriteString("</b>")
		s = s[idx+len(token):]
		lowerS = lowerS[idx+len(token):]
	}
	return b.String()
}
