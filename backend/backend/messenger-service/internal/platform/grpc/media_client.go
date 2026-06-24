package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/moera-sudo/backend/backend/proto/media"
)

type MediaClient struct {
	client pb.MediaServiceClient
	conn   *grpc.ClientConn
}

func NewMediaClient(addr string) (*MediaClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewMediaServiceClient(conn)
	return &MediaClient{client: client, conn: conn}, nil
}

func (c *MediaClient) Close() {
	c.conn.Close()
}

func (c *MediaClient) LinkMediaToEntity(ctx context.Context, mediaID, entityType, entityID string) error {
	_, err := c.client.LinkMediaToEntity(ctx, &pb.LinkMediaToEntityRequest{
		MediaId:    mediaID,
		EntityType: entityType,
		EntityId:   entityID,
	})
	return err
}