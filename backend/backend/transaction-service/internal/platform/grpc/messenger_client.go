package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/moera-sudo/backend/backend/proto/transaction_access"
)

// MessengerClient — клиент к gRPC-сервису TransactionAccessService,
// реализованному в messenger-service. Используется:
//
//   - Wallet feature: ResolveUsername при переводе по @username
//   - Observer: NotifyUserEvent при on-chain событиях
//   - Channel ops: CheckChannelPermission перед withdraw из treasury канала
//
// Соединение однократное, держится открытым в течение жизни процесса.
// Reconnect происходит автоматически на уровне gRPC.
type MessengerClient struct {
	conn   *grpc.ClientConn
	client transaction_access.TransactionAccessServiceClient
}

// NewMessengerClient устанавливает соединение и сохраняет client.
// Если messenger недоступен в момент старта — НЕ блокирует init
// (gRPC устанавливает соединение лениво на первом запросе).
func NewMessengerClient(addr string) (*MessengerClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc: dial messenger %s: %w", addr, err)
	}

	log.Info().Str("addr", addr).Msg("messenger gRPC client created")
	return &MessengerClient{
		conn:   conn,
		client: transaction_access.NewTransactionAccessServiceClient(conn),
	}, nil
}

// Close закрывает соединение. Вызывать при graceful shutdown.
func (mc *MessengerClient) Close() error {
	if mc.conn != nil {
		return mc.conn.Close()
	}
	return nil
}

// ────────────────────────────────────────────────────────────
//  ResolveUsername — @username → user_id + wallet_address
// ────────────────────────────────────────────────────────────

type ResolvedUser struct {
	Found         bool
	UserID        string
	DisplayName   string
	WalletAddress string
}

// ResolveUsername спрашивает messenger «есть ли такой юзер по username
// и если да — какой у него EVM-адрес». Используется в Wallet WebApp
// при отправке SUDA по @username.
//
// username — без префикса @.
func (mc *MessengerClient) ResolveUsername(ctx context.Context, username string) (*ResolvedUser, error) {
	callCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := mc.client.ResolveUsername(callCtx, &transaction_access.ResolveUsernameRequest{
		Username: username,
	})
	if err != nil {
		return nil, fmt.Errorf("grpc: ResolveUsername %s: %w", username, err)
	}

	return &ResolvedUser{
		Found:         resp.GetFound(),
		UserID:        resp.GetUserId(),
		DisplayName:   resp.GetDisplayName(),
		WalletAddress: resp.GetWalletAddress(),
	}, nil
}

// ────────────────────────────────────────────────────────────
//  NotifyUserEvent — WS + опциональный system message в чате
// ────────────────────────────────────────────────────────────

// NotifyEventRequest — параметры уведомления.
// Если TargetChatID и SystemMessageType заданы — Messenger также создаёт
// system message в чате (DONATION / NFT_GIFT / QUEST / FUNDRAISE).
type NotifyEventRequest struct {
	UserID            string
	EventType         string // SUDA_RECEIVED, NFT_RECEIVED, DONATION_RECEIVED и т.д.
	PayloadJSON       string
	TargetChatID      string // опционально
	SystemMessageType string // опционально (DONATION / NFT_GIFT / QUEST / FUNDRAISE)
}

type NotifyEventResult struct {
	WSDelivered     bool
	MessageCreated  bool
	SystemMessageID string
}

// NotifyUserEvent шлёт уведомление через messenger. Используется observer'ом
// после того как он зафиксировал on-chain событие.
//
// Метод НЕ блокирующий для пользовательских операций: если messenger лежит,
// мы залогируем ошибку и продолжим. on-chain истина уже зафиксирована,
// уведомление — best-effort UX.
func (mc *MessengerClient) NotifyUserEvent(ctx context.Context, req NotifyEventRequest) (*NotifyEventResult, error) {
	callCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := mc.client.NotifyUserEvent(callCtx, &transaction_access.NotifyUserEventRequest{
		UserId:            req.UserID,
		EventType:         req.EventType,
		PayloadJson:       req.PayloadJSON,
		TargetChatId:      req.TargetChatID,
		SystemMessageType: req.SystemMessageType,
	})
	if err != nil {
		return nil, fmt.Errorf("grpc: NotifyUserEvent (user=%s event=%s): %w", req.UserID, req.EventType, err)
	}

	return &NotifyEventResult{
		WSDelivered:     resp.GetWsDelivered(),
		MessageCreated:  resp.GetMessageCreated(),
		SystemMessageID: resp.GetSystemMessageId(),
	}, nil
}

// ────────────────────────────────────────────────────────────
//  CheckChannelPermission — проверка прав на действия от лица канала
// ────────────────────────────────────────────────────────────

// Уровни прав, известные обоим сервисам. Совпадают с константами
// в messenger.channel.PermissionLevel.
const (
	PermissionOwner        = "OWNER"
	PermissionOwnerOrAdmin = "OWNER_OR_ADMIN"
	PermissionMember       = "MEMBER"
)

type ChannelPermissionResult struct {
	Granted bool
	Reason  string // "ok" | "not_a_member" | "not_admin" | "not_owner" | "channel_not_found"
}

// CheckChannelPermission проверяет, что юзер имеет указанный уровень прав
// в канале. Используется например при withdraw'е из treasury канала —
// нужно подтвердить что вызывающий = OWNER канала.
func (mc *MessengerClient) CheckChannelPermission(
	ctx context.Context, userID, channelID, permission string,
) (*ChannelPermissionResult, error) {
	callCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := mc.client.CheckChannelPermission(callCtx, &transaction_access.ChannelPermissionRequest{
		UserId:     userID,
		ChannelId:  channelID,
		Permission: permission,
	})
	if err != nil {
		return nil, fmt.Errorf("grpc: CheckChannelPermission (user=%s channel=%s perm=%s): %w",
			userID, channelID, permission, err)
	}

	return &ChannelPermissionResult{
		Granted: resp.GetGranted(),
		Reason:  resp.GetReason(),
	}, nil
}

// ────────────────────────────────────────────────────────────
//  CheckChatMembership — batch-проверка членства в чате
// ────────────────────────────────────────────────────────────

// UserMembership — результат проверки членства одного юзера.
type UserMembership struct {
	UserID   string
	IsMember bool
	Role     string // OWNER | ADMIN | MEMBER | SUBSCRIBER (если IsMember=true)
}

// ChatMembershipResult — результат batch-проверки.
type ChatMembershipResult struct {
	ChatExists bool
	Members    []UserMembership // в том же порядке, что userIDs в запросе
}

// CheckChatMembership спрашивает messenger «эти юзеры — члены этого чата?»
//
// Используется wallet-сервисом перед broadcast'ом transfer-in-chat,
// чтобы убедиться что и отправитель, и получатель в чате.
//
// Возвращает результат для каждого user_id в том же порядке.
// Если чат не существует — ChatExists=false и все IsMember=false.
func (mc *MessengerClient) CheckChatMembership(
	ctx context.Context, chatID string, userIDs []string,
) (*ChatMembershipResult, error) {
	callCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := mc.client.CheckChatMembership(callCtx, &transaction_access.ChatMembershipRequest{
		ChatId:  chatID,
		UserIds: userIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("grpc: CheckChatMembership (chat=%s users=%d): %w", chatID, len(userIDs), err)
	}

	members := make([]UserMembership, 0, len(resp.GetMembers()))
	for _, m := range resp.GetMembers() {
		members = append(members, UserMembership{
			UserID:   m.GetUserId(),
			IsMember: m.GetIsMember(),
			Role:     m.GetRole(),
		})
	}

	return &ChatMembershipResult{
		ChatExists: resp.GetChatExists(),
		Members:    members,
	}, nil
}

// ────────────────────────────────────────────────────────────
//  GetUsersByIDs — batch-резолв user_id → username
// ────────────────────────────────────────────────────────────

// UserBrief — краткая инфа о юзере (username для UI).
type UserBrief struct {
	UserID      string
	Username    string
	DisplayName string
}

// GetUsersByIDs резолвит список user_id в username + display_name.
//
// Используется treasury-статистикой канала (top-донатеры, список донатов).
// Возвращает map user_id → UserBrief; отсутствующие в ответе id означают
// что юзер не найден. Пустой список ids → пустая map без gRPC-вызова.
func (mc *MessengerClient) GetUsersByIDs(ctx context.Context, userIDs []string) (map[string]UserBrief, error) {
	if len(userIDs) == 0 {
		return map[string]UserBrief{}, nil
	}

	callCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := mc.client.GetUsersByIDs(callCtx, &transaction_access.GetUsersByIDsRequest{
		UserIds: userIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("grpc: GetUsersByIDs (count=%d): %w", len(userIDs), err)
	}

	out := make(map[string]UserBrief, len(resp.GetUsers()))
	for _, u := range resp.GetUsers() {
		out[u.GetUserId()] = UserBrief{
			UserID:      u.GetUserId(),
			Username:    u.GetUsername(),
			DisplayName: u.GetDisplayName(),
		}
	}
	return out, nil
}