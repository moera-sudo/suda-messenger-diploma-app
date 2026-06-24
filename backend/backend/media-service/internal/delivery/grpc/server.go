// Package grpc exposes media-service over gRPC for other services.
//
// messenger-service calls LinkMediaToEntity right after a message with an
// attachment is saved. Without this server running, that call fails silently
// and the media never gets an entity link — which makes the ACL deny every
// non-owner (checkAccess sees no links → 403). So this server is what lets
// chat members (not just the uploader) open/play attachments.
package grpc

import (
	"context"

	pb "github.com/moera-sudo/backend/backend/proto/media"

	"github.com/moera-sudo/backend/backend/media-service/internal/service"
)

// Server implements pb.MediaServiceServer on top of the media service.
// Only LinkMediaToEntity is used today (by messenger-service); the rest fall
// back to UnimplementedMediaServiceServer.
type Server struct {
	pb.UnimplementedMediaServiceServer
	svc *service.MediaService
}

func NewServer(svc *service.MediaService) *Server {
	return &Server{svc: svc}
}

// LinkMediaToEntity links an uploaded media to a domain entity (e.g. a chat
// MESSAGE). Subsequent ACL checks rely on this link to grant access to chat
// members.
func (s *Server) LinkMediaToEntity(ctx context.Context, req *pb.LinkMediaToEntityRequest) (*pb.LinkMediaToEntityResponse, error) {
	if err := s.svc.LinkToEntity(ctx, req.GetMediaId(), req.GetEntityType(), req.GetEntityId()); err != nil {
		return nil, err
	}
	return &pb.LinkMediaToEntityResponse{Success: true}, nil
}
