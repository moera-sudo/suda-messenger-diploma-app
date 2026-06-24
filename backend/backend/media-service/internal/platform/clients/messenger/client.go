package messenger

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/moera-sudo/backend/backend/proto/messenger_access"
)

type AccessChecker interface {
	CheckAccess(ctx context.Context, userID, entityType, entityID string) (bool, error)
	Close() error
}


type messengerClient struct {
	conn *grpc.ClientConn
	client pb.MessengerAccessServiceClient
}

func NewAccessChecker(addr string) (AccessChecker, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect to messenger-service at %s: %w", addr, err)
	}

	client := pb.NewMessengerAccessServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.CheckEntityAccess(ctx, &pb.CheckEntityAccessRequest{
		UserId:     "health-check",
		EntityType: "HEALTH",
		EntityId:   "0",
	})
	// Ошибка от сервера — это ок, значит он ответил
	// Ошибка соединения — значит сервис недоступен
	if err != nil && !isServerError(err) {
		conn.Close()
		return nil, fmt.Errorf("messenger-service not reachable: %w", err)
	}

	return &messengerClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *messengerClient) CheckAccess(ctx context.Context, userID, entityType, entityID string) (bool, error) {
	resp, err := c.client.CheckEntityAccess(ctx, &pb.CheckEntityAccessRequest{
		UserId:     userID,
		EntityType: entityType,
		EntityId:   entityID,
	})
	if err != nil {
		return false, fmt.Errorf("check entity access: %w", err)
	}

	return resp.HasAccess, nil
}

func (c *messengerClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// isServerError проверяет что ошибка пришла от сервера (а не от сети)
func isServerError(err error) bool {
	// gRPC server вернул ответ с ошибкой — значит он живой
	return err != nil && (
		// Unimplemented = сервер запущен, но метод не найден
		// InvalidArgument, Internal и т.д. = сервер обработал запрос
		true)
}


type NoopAccessChecker struct{}

func NewNoopAccessChecker() AccessChecker {
	return &NoopAccessChecker{}
}

func (n *NoopAccessChecker) CheckAccess(_ context.Context, _, _, _ string) (bool, error) {
	return true, nil
}

func (n *NoopAccessChecker) Close() error {
	return nil
}