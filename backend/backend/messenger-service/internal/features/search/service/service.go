package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat/repository"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/search"
	searchRepo "github.com/moera-sudo/backend/backend/messenger-service/internal/features/search/repository"
)

type Service struct {
	repo    searchRepo.Repository
	members repository.MemberRepo
}

func NewService(repo searchRepo.Repository, members repository.MemberRepo) *Service {
	return &Service{repo: repo, members: members}
}

func (s *Service) GlobalSearch(ctx context.Context, userID uuid.UUID, query string) (*search.GlobalSearchResponse, error) {
	if len(query) < 2 {
		return &search.GlobalSearchResponse{}, nil
	}

	users, err := s.repo.SearchUsers(ctx, query, 10)
	if err != nil {
		log.Warn().Err(err).Str("query", query).Msg("Global search: failed to search users")
		users = nil
	}

	chats, err := s.repo.SearchChats(ctx, query, 10)
	if err != nil {
		log.Warn().Err(err).Str("query", query).Msg("Global search: failed to search chats")
		chats = nil
	}

	messages, err := s.repo.SearchMessagesGlobal(ctx, userID, query, 20)
	if err != nil {
		log.Warn().Err(err).Str("query", query).Msg("Global search: failed to search messages")
		messages = nil
	}

	log.Debug().Str("user_id", userID.String()).Str("query", query).Int("users", len(users)).Int("chats", len(chats)).Int("messages", len(messages)).Msg("Global search completed")

	// Ensure non-nil slices for JSON
	if users == nil {
		users = []search.SearchResult{}
	}
	if chats == nil {
		chats = []search.SearchResult{}
	}
	if messages == nil {
		messages = []search.SearchResult{}
	}

	return &search.GlobalSearchResponse{
		Users:    users,
		Chats:    chats,
		Messages: messages,
	}, nil
}

func (s *Service) SearchInChat(ctx context.Context, userID uuid.UUID, chatID uuid.UUID, query string) (*search.ChatSearchResponse, error) {
	if len(query) < 2 {
		return &search.ChatSearchResponse{Messages: []search.SearchResult{}}, nil
	}

	// Проверяем membership
	isMember, err := s.members.IsMember(ctx, chatID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, fmt.Errorf("forbidden: not a member of this chat")
	}

	messages, err := s.repo.SearchMessagesInChat(ctx, chatID, query, 50)
	if err != nil {
		log.Error().Err(err).Str("chat_id", chatID.String()).Str("query", query).Msg("Failed to search messages in chat")
		return nil, err
	}

	log.Debug().Str("chat_id", chatID.String()).Str("query", query).Int("results", len(messages)).Msg("Chat search completed")

	if messages == nil {
		messages = []search.SearchResult{}
	}

	return &search.ChatSearchResponse{
		Messages: messages,
		Total:    len(messages),
	}, nil
}
