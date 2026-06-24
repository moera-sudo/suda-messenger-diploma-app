package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
)

func (s *ChatService) SetContactName(ctx context.Context, ownerID, contactID uuid.UUID, name string) error {
	if ownerID == contactID {
		return fmt.Errorf("cannot set custom name for yourself")
	}

	blocked, err := s.mods.IsBlocked(ctx, contactID, ownerID)
	if err != nil {
		return err
	}
	if blocked {
		return ErrBlocked
	}

	if err := s.mods.SetContactName(ctx, ownerID, contactID, name); err != nil {
		log.Error().Err(err).Str("owner_id", ownerID.String()).Str("contact_id", contactID.String()).Msg("Failed to set contact name")
		return err
	}

	log.Debug().Str("owner_id", ownerID.String()).Str("contact_id", contactID.String()).Msg("Contact name set")
	return nil
}

func (s *ChatService) RemoveContactName(ctx context.Context, ownerID, contactID uuid.UUID) error {
	if err := s.mods.RemoveContactName(ctx, ownerID, contactID); err != nil {
		log.Error().Err(err).Str("owner_id", ownerID.String()).Str("contact_id", contactID.String()).Msg("Failed to remove contact name")
		return fmt.Errorf("remove contact name: %w", err)
	}

	log.Debug().Str("owner_id", ownerID.String()).Str("contact_id", contactID.String()).Msg("Contact name removed")
	return nil
}

func (s *ChatService) GetContacts(ctx context.Context, ownerID uuid.UUID) ([]chat.Contact, error) {
	return s.mods.GetContacts(ctx, ownerID)
}