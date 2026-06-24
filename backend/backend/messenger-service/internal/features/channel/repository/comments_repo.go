package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/channel"
)

type commentsRepo struct {
	db *pgxpool.Pool
}

func (r *commentsRepo) Create(ctx context.Context, c *channel.ChannelComment) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO messenger_channel_post_comments
			(channel_id, post_id, sender_id, content, reply_to_comment_id)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at`,
		c.ChannelID, c.PostID, c.SenderID, c.Content, c.ReplyToCommentID,
	).Scan(&c.ID, &c.CreatedAt)
}

func (r *commentsRepo) GetByID(ctx context.Context, commentID int64) (*channel.ChannelComment, error) {
	var c channel.ChannelComment
	err := r.db.QueryRow(ctx,
		`SELECT id, channel_id, post_id, sender_id, content,
		        reply_to_comment_id, edited_at, deleted_at, created_at
		 FROM messenger_channel_post_comments
		 WHERE id = $1 AND deleted_at IS NULL`,
		commentID,
	).Scan(&c.ID, &c.ChannelID, &c.PostID, &c.SenderID, &c.Content,
		&c.ReplyToCommentID, &c.EditedAt, &c.DeletedAt, &c.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &c, err
}

func (r *commentsRepo) GetByPost(ctx context.Context, postID int64, limit, offset int) ([]channel.ChannelComment, error) {
	// LEFT JOIN the parent comment (pc) + its author (pu) so each reply carries a
	// hydrated preview and the client never depends on the parent being in the
	// loaded window.
	rows, err := r.db.Query(ctx,
		`SELECT c.id, c.channel_id, c.post_id, c.sender_id, c.content,
		        c.reply_to_comment_id, c.edited_at, c.created_at,
		        u.username, u.display_name,
		        pc.content, COALESCE(pu.display_name, pu.username, ''),
		        (pc.deleted_at IS NOT NULL)
		 FROM messenger_channel_post_comments c
		 JOIN messenger_users u ON c.sender_id = u.id
		 LEFT JOIN messenger_channel_post_comments pc ON pc.id = c.reply_to_comment_id
		 LEFT JOIN messenger_users pu ON pu.id = pc.sender_id
		 WHERE c.post_id = $1 AND c.deleted_at IS NULL
		 ORDER BY c.created_at ASC
		 LIMIT $2 OFFSET $3`,
		postID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []channel.ChannelComment
	for rows.Next() {
		var (
			c             channel.ChannelComment
			parentContent *string
			parentName    string
			parentDeleted *bool
		)
		if err := rows.Scan(&c.ID, &c.ChannelID, &c.PostID, &c.SenderID, &c.Content,
			&c.ReplyToCommentID, &c.EditedAt, &c.CreatedAt,
			&c.SenderUsername, &c.SenderDisplayName,
			&parentContent, &parentName, &parentDeleted); err != nil {
			return nil, err
		}
		// Populate the reply preview only when the parent exists and isn't deleted.
		if c.ReplyToCommentID != nil && parentContent != nil && (parentDeleted == nil || !*parentDeleted) {
			c.ReplyPreview = &channel.CommentReplyPreview{
				SenderDisplayName: strings.TrimPrefix(parentName, "@"),
				Content:           *parentContent,
			}
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (r *commentsRepo) Update(ctx context.Context, commentID int64, content string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE messenger_channel_post_comments
		 SET content = $1, edited_at = NOW()
		 WHERE id = $2 AND deleted_at IS NULL`,
		content, commentID,
	)
	return err
}

func (r *commentsRepo) SoftDelete(ctx context.Context, commentID int64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE messenger_channel_post_comments
		 SET deleted_at = NOW()
		 WHERE id = $1`,
		commentID,
	)
	return err
}

func (r *commentsRepo) CountByPost(ctx context.Context, postID int64) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM messenger_channel_post_comments
		 WHERE post_id = $1 AND deleted_at IS NULL`,
		postID,
	).Scan(&count)
	return count, err
}