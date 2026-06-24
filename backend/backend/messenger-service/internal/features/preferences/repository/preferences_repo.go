package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/preferences"
)

type preferencesRepo struct {
	db *pgxpool.Pool
}

func (r *preferencesRepo) Get(ctx context.Context, userID uuid.UUID) (*preferences.UserPreferences, error) {
	var p preferences.UserPreferences
	err := r.db.QueryRow(ctx,
		`SELECT user_id, lang_code, theme,
		        last_seen_visibility, avatar_visibility, who_can_message,
		        who_can_add_to_groups, accept_gifts_from,
		        notification_messages, notification_groups, notification_channels,
		        notification_sound, notification_vibration, preview_in_notifications,
		        chat_font_size, send_by_enter, auto_download_media, wallpaper_id,
		        nft_showcase_enabled, updated_at
		 FROM messenger_user_preferences WHERE user_id = $1`,
		userID,
	).Scan(
		&p.UserID, &p.LangCode, &p.Theme,
		&p.LastSeenVisibility, &p.AvatarVisibility, &p.WhoCanMessage,
		&p.WhoCanAddToGroups, &p.AcceptGiftsFrom,
		&p.NotificationMessages, &p.NotificationGroups, &p.NotificationChannels,
		&p.NotificationSound, &p.NotificationVibration, &p.PreviewInNotifications,
		&p.ChatFontSize, &p.SendByEnter, &p.AutoDownloadMedia, &p.WallpaperID,
		&p.NFTShowcaseEnabled, &p.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &p, err
}

func (r *preferencesRepo) Create(ctx context.Context, p *preferences.UserPreferences) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO messenger_user_preferences
			(user_id, lang_code, theme,
			 last_seen_visibility, avatar_visibility, who_can_message,
			 who_can_add_to_groups, accept_gifts_from,
			 notification_messages, notification_groups, notification_channels,
			 notification_sound, notification_vibration, preview_in_notifications,
			 chat_font_size, send_by_enter, auto_download_media, wallpaper_id,
			 nft_showcase_enabled)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		 ON CONFLICT (user_id) DO NOTHING`,
		p.UserID, p.LangCode, p.Theme,
		p.LastSeenVisibility, p.AvatarVisibility, p.WhoCanMessage,
		p.WhoCanAddToGroups, p.AcceptGiftsFrom,
		p.NotificationMessages, p.NotificationGroups, p.NotificationChannels,
		p.NotificationSound, p.NotificationVibration, p.PreviewInNotifications,
		p.ChatFontSize, p.SendByEnter, p.AutoDownloadMedia, p.WallpaperID,
		p.NFTShowcaseEnabled,
	)
	return err
}

func (r *preferencesRepo) Update(ctx context.Context, userID uuid.UUID, req *preferences.UpdatePreferencesReq) error {
	var sets []string
	var args []interface{}
	argNum := 1

	if req.LangCode != nil {
		sets = append(sets, fmt.Sprintf("lang_code = $%d", argNum))
		args = append(args, *req.LangCode)
		argNum++
	}
	if req.Theme != nil {
		sets = append(sets, fmt.Sprintf("theme = $%d", argNum))
		args = append(args, *req.Theme)
		argNum++
	}
	if req.LastSeenVisibility != nil {
		sets = append(sets, fmt.Sprintf("last_seen_visibility = $%d", argNum))
		args = append(args, *req.LastSeenVisibility)
		argNum++
	}
	if req.AvatarVisibility != nil {
		sets = append(sets, fmt.Sprintf("avatar_visibility = $%d", argNum))
		args = append(args, *req.AvatarVisibility)
		argNum++
	}
	if req.WhoCanMessage != nil {
		sets = append(sets, fmt.Sprintf("who_can_message = $%d", argNum))
		args = append(args, *req.WhoCanMessage)
		argNum++
	}
	if req.NotificationMessages != nil {
		sets = append(sets, fmt.Sprintf("notification_messages = $%d", argNum))
		args = append(args, *req.NotificationMessages)
		argNum++
	}
	if req.NotificationGroups != nil {
		sets = append(sets, fmt.Sprintf("notification_groups = $%d", argNum))
		args = append(args, *req.NotificationGroups)
		argNum++
	}
	if req.NotificationChannels != nil {
		sets = append(sets, fmt.Sprintf("notification_channels = $%d", argNum))
		args = append(args, *req.NotificationChannels)
		argNum++
	}
	if req.NotificationSound != nil {
		sets = append(sets, fmt.Sprintf("notification_sound = $%d", argNum))
		args = append(args, *req.NotificationSound)
		argNum++
	}
	if req.NotificationVibration != nil {
		sets = append(sets, fmt.Sprintf("notification_vibration = $%d", argNum))
		args = append(args, *req.NotificationVibration)
		argNum++
	}
	if req.ChatFontSize != nil {
		sets = append(sets, fmt.Sprintf("chat_font_size = $%d", argNum))
		args = append(args, *req.ChatFontSize)
		argNum++
	}
	if req.SendByEnter != nil {
		sets = append(sets, fmt.Sprintf("send_by_enter = $%d", argNum))
		args = append(args, *req.SendByEnter)
		argNum++
	}
	if req.AutoDownloadMedia != nil {
		sets = append(sets, fmt.Sprintf("auto_download_media = $%d", argNum))
		args = append(args, *req.AutoDownloadMedia)
		argNum++
	}
	if req.WhoCanAddToGroups != nil {
		sets = append(sets, fmt.Sprintf("who_can_add_to_groups = $%d", argNum))
		args = append(args, *req.WhoCanAddToGroups)
		argNum++
	}
	if req.AcceptGiftsFrom != nil {
		sets = append(sets, fmt.Sprintf("accept_gifts_from = $%d", argNum))
		args = append(args, *req.AcceptGiftsFrom)
		argNum++
	}
	if req.PreviewInNotifications != nil {
		sets = append(sets, fmt.Sprintf("preview_in_notifications = $%d", argNum))
		args = append(args, *req.PreviewInNotifications)
		argNum++
	}
	if req.WallpaperID != nil {
		sets = append(sets, fmt.Sprintf("wallpaper_id = $%d", argNum))
		args = append(args, *req.WallpaperID)
		argNum++
	}
	if req.NFTShowcaseEnabled != nil {
		sets = append(sets, fmt.Sprintf("nft_showcase_enabled = $%d", argNum))
		args = append(args, *req.NFTShowcaseEnabled)
		argNum++
	}

	if len(sets) == 0 {
		return nil // ничего обновлять
	}

	sets = append(sets, "updated_at = NOW()")
	args = append(args, userID)

	query := fmt.Sprintf(
		`UPDATE messenger_user_preferences SET %s WHERE user_id = $%d`,
		strings.Join(sets, ", "), argNum,
	)

	_, err := r.db.Exec(ctx, query, args...)
	return err
}