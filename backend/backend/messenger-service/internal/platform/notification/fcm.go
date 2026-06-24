package notification

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)


type Service interface {
	SendPush(ctx context.Context, userIDs []uuid.UUID, title, body string, data map[string]string) error
	SaveToken(ctx context.Context, userID uuid.UUID, token, deviceName string) error
}

type fcmService struct {
	client *messaging.Client
	db *pgxpool.Pool
}

func NewFCMService(credentialsFile string, db *pgxpool.Pool) (Service, error) {
	if credentialsFile == "" {
		return &fcmService{client: nil, db: db}, nil
	}
	
	opt := option.WithCredentialsFile(credentialsFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("Error initializing firebase app: %v", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Error getting messaging client: %v", err)
	}

	return &fcmService{
		client: client,
		db: db,
	}, nil
}


func (s *fcmService) SendPush(ctx context.Context, userIDs []uuid.UUID, title, body string, data map[string]string) error {
	if s.client == nil{
		return nil
	}
	
	tokens, err := s.getTokensForUsers(ctx, userIDs)
	if err != nil {
		log.Error().Err(err).Int("user_count", len(userIDs)).Msg("Failed to get FCM tokens")
		return err
	}

	if len(tokens) == 0 {
		log.Debug().Int("user_count", len(userIDs)).Msg("No FCM tokens found, skipping push")
		return nil
	}

	// 2. Отправляем батчами по 500 (ограничение Firebase), но для диплома шлем всё сразу
	// Firebase Multicast поддерживает до 500 токенов за раз.
	// Если токенов больше 500, тут нужен цикл.

	message := &messaging.MulticastMessage{
		Tokens: tokens,
		Data:   data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{Title: title, Body: body}, // Если хотим системный пуш
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{ContentAvailable: true},
			},
		},
	}

	br, err := s.client.SendEachForMulticast(ctx, message)

	if err != nil {
		log.Error().Err(err).Int("token_count", len(tokens)).Msg("Failed to send FCM multicast")
		return err
	}

	log.Debug().Int("success", br.SuccessCount).Int("failure", br.FailureCount).Str("title", title).Msg("FCM push sent")

	if br.FailureCount > 0 {
		go s.cleanInvalidTokens(context.Background(), tokens, br.Responses)
	}

	return nil
}

func (s *fcmService) SaveToken(ctx context.Context, userID uuid.UUID, token, deviceName string) error {

	query := `
		INSERT INTO messenger_user_devices (user_id, fcm_token, device_name, last_used_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, fcm_token) DO UPDATE
		SET last_used_at = NOW(), device_name = EXCLUDED.device_name
	`

	_, err := s.db.Exec(ctx, query, userID, token, deviceName)
	return err
}


func (s *fcmService) getTokensForUsers(ctx context.Context, userIDs []uuid.UUID) ([]string, error) {
	query := `
		SELECT fcm_token
		FROM messenger_user_devices
		WHERE user_id = ANY($1)
			AND last_used_at > NOW() - INTERVAL '6 months'
	`

	rows, err := s.db.Query(ctx, query, userIDs)
	if err != nil{
		return nil, err
	}

	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err == nil{
			tokens = append(tokens, t)
		}
	}
	return tokens, nil
}

func (s *fcmService) cleanInvalidTokens(ctx context.Context, tokens []string, responses []*messaging.SendResponse) {
	var tokensToDelete []string

	for i, resp := range responses{
		if !resp.Success {
			if messaging.IsUnregistered(resp.Error) || 
				messaging.IsSenderIDMismatch(resp.Error) {
					tokensToDelete = append(tokensToDelete, tokens[i])
				}
		}
	}

	if len(tokensToDelete) > 0{
		log.Info().Int("count", len(tokensToDelete)).Msg("Cleaning invalid FCM tokens")

		query := `DELETE FROM messenger_user_devices WHERE fcm_token = ANY($1)`
		_, err := s.db.Exec(ctx, query, tokensToDelete)
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete invalid tokens")
		}
	}
}