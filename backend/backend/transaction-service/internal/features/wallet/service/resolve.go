package service

import (
	"context"
	"fmt"
	"strings"
)

// ResolvedUser — результат разрешения @username.
// Совпадает по форме с MessengerClient.ResolveUsername, но живёт
// в нашем service-слое чтобы handler не знал про platgrpc типы.
type ResolvedUser struct {
	Found         bool
	UserID        string
	DisplayName   string
	WalletAddress string
}

// ResolveUsername разрешает @username в (user_id, address, display_name).
// Префикс "@" удаляется автоматически, юзер может прислать что в "@john"
// что в "john".
func (s *Service) ResolveUsername(ctx context.Context, username string) (*ResolvedUser, error) {
	clean := strings.TrimPrefix(strings.TrimSpace(username), "@")
	if clean == "" {
		return nil, fmt.Errorf("%w: username is empty", ErrInvalidInput)
	}

	resolved, err := s.messenger.ResolveUsername(ctx, clean)
	if err != nil {
		return nil, fmt.Errorf("resolve username: %w", err)
	}

	return &ResolvedUser{
		Found:         resolved.Found,
		UserID:        resolved.UserID,
		DisplayName:   resolved.DisplayName,
		WalletAddress: resolved.WalletAddress,
	}, nil
}