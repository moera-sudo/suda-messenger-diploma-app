package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/channel"
)

type appsRepo struct {
	db *pgxpool.Pool
}

func (r *appsRepo) Link(ctx context.Context, link *channel.ChannelAppLink) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO messenger_channel_apps (channel_id, app_id, added_by)
		 VALUES ($1, $2, $3) ON CONFLICT (channel_id, app_id) DO NOTHING`,
		link.ChannelID, link.AppID, link.AddedBy,
	)
	return err
}

func (r *appsRepo) Unlink(ctx context.Context, channelID, appID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM messenger_channel_apps WHERE channel_id = $1 AND app_id = $2`,
		channelID, appID,
	)
	return err
}

func (r *appsRepo) GetByChannel(ctx context.Context, channelID uuid.UUID) ([]channel.ChannelAppLink, error) {
	rows, err := r.db.Query(ctx,
		`SELECT channel_id, app_id, added_by, added_at
		 FROM messenger_channel_apps WHERE channel_id = $1
		 ORDER BY added_at DESC`,
		channelID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []channel.ChannelAppLink
	for rows.Next() {
		var l channel.ChannelAppLink
		if err := rows.Scan(&l.ChannelID, &l.AppID, &l.AddedBy, &l.AddedAt); err != nil {
			return nil, err
		}
		links = append(links, l)
	}
	return links, rows.Err()
}