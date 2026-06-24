package grpc

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	pb "github.com/moera-sudo/backend/backend/proto/messenger_access"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat/repository"
)

type AccessServer struct {
	pb.UnimplementedMessengerAccessServiceServer
	memberRepo repository.MemberRepo
}

func NewAccessServer(memberRepo repository.MemberRepo) *AccessServer {
	return &AccessServer{memberRepo: memberRepo}
}

func (s *AccessServer) CheckEntityAccess(ctx context.Context, req *pb.CheckEntityAccessRequest) (*pb.CheckEntityAccessResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &pb.CheckEntityAccessResponse{
			HasAccess: false,
			Reason:    "invalid_user_id",
		}, nil
	}

	switch req.EntityType {
	case "MESSAGE":
		return s.checkMessageAccess(ctx, userID, req.EntityId)

	case "CHAT_AVATAR":
		return s.checkChatAccess(ctx, userID, req.EntityId)

	case "USER_AVATAR":
		// The profile layer (GetUserWithBlockCheck: avatar visibility + block checks)
		// already gates who receives a user's avatar_media_id. Possession of the
		// unguessable media_id therefore implies authorization, so don't 403 here —
		// otherwise a private user avatar would be unviewable by everyone else.
		return &pb.CheckEntityAccessResponse{HasAccess: true}, nil

	default:
		log.Warn().
			Str("entity_type", req.EntityType).
			Msg("Unknown entity type in access check")

		return &pb.CheckEntityAccessResponse{
			HasAccess: false,
			Reason:    "unknown_entity_type",
		}, nil
	}
}

func (s *AccessServer) checkMessageAccess(ctx context.Context, userID uuid.UUID, entityID string) (*pb.CheckEntityAccessResponse, error) {
	messageID, err := strconv.ParseInt(entityID, 10, 64)
	if err != nil {
		return &pb.CheckEntityAccessResponse{
			HasAccess: false,
			Reason:    "invalid_message_id",
		}, nil
	}

	isMember, err := s.memberRepo.IsMemberByMessageID(ctx, userID, messageID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check message access")
		return &pb.CheckEntityAccessResponse{
			HasAccess: false,
			Reason:    "internal_error",
		}, nil
	}

	if !isMember {
		return &pb.CheckEntityAccessResponse{
			HasAccess: false,
			Reason:    "not_a_member",
		}, nil
	}

	return &pb.CheckEntityAccessResponse{HasAccess: true}, nil
}

func (s *AccessServer) checkChatAccess(ctx context.Context, userID uuid.UUID, entityID string) (*pb.CheckEntityAccessResponse, error) {
	chatID, err := uuid.Parse(entityID)
	if err != nil {
		return &pb.CheckEntityAccessResponse{
			HasAccess: false,
			Reason:    "invalid_chat_id",
		}, nil
	}

	// IsMember signature is (chatID, userID) — pass in that order, not (userID, chatID),
	// otherwise non-uploader members get 403 when resolving a CHAT_AVATAR media_id.
	isMember, err := s.memberRepo.IsMember(ctx, chatID, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check chat access")
		return &pb.CheckEntityAccessResponse{
			HasAccess: false,
			Reason:    "internal_error",
		}, nil
	}

	if !isMember {
		return &pb.CheckEntityAccessResponse{
			HasAccess: false,
			Reason:    "not_a_member",
		}, nil
	}

	return &pb.CheckEntityAccessResponse{HasAccess: true}, nil
}