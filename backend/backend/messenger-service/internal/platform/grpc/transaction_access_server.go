package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
	chatRepo "github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat/repository"
	userRepo "github.com/moera-sudo/backend/backend/messenger-service/internal/features/user/repository"
	socketmodel "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket/hub"

	pb "github.com/moera-sudo/backend/backend/proto/transaction_access"
)

// TransactionAccessServer — gRPC-сервер, к которому обращается transaction-service.
//
// Три метода:
//  1. ResolveUsername — поиск юзера по @username для P2P-перевода SUDA.
//  2. NotifyUserEvent — слать WS-event юзеру (+ опционально создать system message
//     в чате; будет реализовано в Этапе 3).
//  3. CheckChannelPermission — проверка прав юзера в канале (OWNER/ADMIN/MEMBER).
type TransactionAccessServer struct {
	pb.UnimplementedTransactionAccessServiceServer

	userRepo    userRepo.Repository
	memberRepo  chatRepo.MemberRepo
	messageRepo chatRepo.MessageRepo
	hub         *hub.Hub
}

// NewTransactionAccessServer создаёт сервер с зависимостями.
// userRepo  — для GetUserByExactUsername.
// memberRepo — для GetMemberRole в CheckChannelPermission.
// hub — для WS-доставки в NotifyUserEvent.
func NewTransactionAccessServer(
	userRepo userRepo.Repository,
	memberRepo chatRepo.MemberRepo,
	messageRepo chatRepo.MessageRepo,
	hub *hub.Hub,
) *TransactionAccessServer {
	return &TransactionAccessServer{
		userRepo:    userRepo,
		memberRepo:  memberRepo,
		messageRepo: messageRepo,
		hub:         hub,
	}
}

// ────────────────────────────────────────────────────────────
//  1. ResolveUsername
// ────────────────────────────────────────────────────────────

// ResolveUsername — точный lookup по @username.
//
// Вызывается из transaction-service когда юзер делает P2P-перевод по username.
// Возвращает found=false без ошибки, если юзера нет — это нормальный кейс.
// Возвращает found=true с пустым wallet_address, если юзер существует но
// кошелька ещё нет (не верифицирован или CreateWallet упал в Verify).
func (s *TransactionAccessServer) ResolveUsername(
	ctx context.Context, req *pb.ResolveUsernameRequest,
) (*pb.ResolveUsernameResponse, error) {
	username := req.GetUsername()
	if username == "" {
		return &pb.ResolveUsernameResponse{Found: false}, nil
	}

	u, err := s.userRepo.GetUserByExactUsername(ctx, username)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("ResolveUsername: db error")
		return nil, err
	}
	if u == nil {
		// Юзер не найден — это не ошибка, просто кейс.
		return &pb.ResolveUsernameResponse{Found: false}, nil
	}

	return &pb.ResolveUsernameResponse{
		Found:         true,
		UserId:        u.ID.String(),
		DisplayName:   u.DisplayName,
		WalletAddress: u.WalletAddress, // пусто, если кошелька ещё нет
	}, nil
}

// ────────────────────────────────────────────────────────────
//  2. NotifyUserEvent
// ────────────────────────────────────────────────────────────

// NotifyUserEvent — доставка WS-события юзеру + (опционально) создание
// system message в чате.
//
// Логика:
//  1. Шлём WS-event юзеру (через hub.SendToUser).
//  2. Если задан target_chat_id и system_message_type — дополнительно:
//     a. INSERT system message в messenger_messages с type=<system_message_type>
//     и content=payload_json. sender_id = NULL (системное сообщение).
//     b. Broadcast этого сообщения всем участникам чата через WS
//     (event NEW_MESSAGE), чтобы UI у всех показал карточку.
//
// Возвращает:
//
//	ws_delivered      — true если WS-event адресату прошёл (юзер был онлайн)
//	message_created   — true если system message создан в чате
//	system_message_id — id созданного сообщения (если был создан)
func (s *TransactionAccessServer) NotifyUserEvent(
	ctx context.Context, req *pb.NotifyUserEventRequest,
) (*pb.NotifyUserEventResponse, error) {
	userIDStr := req.GetUserId()
	eventType := req.GetEventType()
	payloadJSON := req.GetPayloadJson()

	if userIDStr == "" || eventType == "" {
		return nil, errors.New("user_id and event_type are required")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("user_id must be uuid")
	}

	// Парсим payload JSON в interface{}, чтобы WS получил структурированный объект.
	var payload any
	if payloadJSON != "" {
		if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
			log.Warn().Err(err).
				Str("event_type", eventType).
				Str("payload", payloadJSON).
				Msg("NotifyUserEvent: payload is not valid JSON, sending as string")
			payload = payloadJSON
		}
	}

	// 1. Слать WS-event адресату.
	wsMsg := socketmodel.WSMessage{
		Type:    eventType,
		Payload: payload,
	}
	err = s.hub.SendToUser(userID, wsMsg)
	wsDelivered := err == nil
	if !wsDelivered {
		log.Debug().
			Str("user_id", userIDStr).
			Str("event_type", eventType).
			Err(err).
			Msg("NotifyUserEvent: WS delivery to recipient failed (user offline)")
	}

	// 2. Если задан target_chat_id + system_message_type — создаём system message.
	var messageCreated bool
	var systemMessageID string

	targetChatIDStr := req.GetTargetChatId()
	systemMsgType := req.GetSystemMessageType()

	if targetChatIDStr != "" && systemMsgType != "" {
		msgID, err := s.createAndBroadcastSystemMessage(
			ctx, userID, targetChatIDStr, systemMsgType, payloadJSON,
		)
		if err != nil {
			log.Error().Err(err).
				Str("chat_id", targetChatIDStr).
				Str("msg_type", systemMsgType).
				Msg("NotifyUserEvent: failed to create system message")
			// Не возвращаем ошибку — WS-event уже доставлен, это best-effort.
		} else {
			messageCreated = true
			systemMessageID = msgID
		}
	}

	log.Info().
		Str("user_id", userIDStr).
		Str("event_type", eventType).
		Bool("ws_delivered", wsDelivered).
		Bool("message_created", messageCreated).
		Str("system_message_id", systemMessageID).
		Msg("NotifyUserEvent done")

	return &pb.NotifyUserEventResponse{
		WsDelivered:     wsDelivered,
		MessageCreated:  messageCreated,
		SystemMessageId: systemMessageID,
	}, nil
}

// ────────────────────────────────────────────────────────────
//  3. CheckChannelPermission
// ────────────────────────────────────────────────────────────

// CheckChannelPermission — проверка прав юзера в чате/канале.
//
// permission values:
//
//	"OWNER"          — только владелец
//	"OWNER_OR_ADMIN" — владелец или администратор
//	"MEMBER"         — любой участник
//
// Возвращает granted=false с reason если права не подходят:
//
//	"not_a_member"      — юзер не в чате
//	"not_admin"         — юзер в чате, но не админ
//	"not_owner"         — юзер админ/мембер, но не владелец
//	"channel_not_found" — чат не существует
//	"ok"                — права подходят
func (s *TransactionAccessServer) CheckChannelPermission(
	ctx context.Context, req *pb.ChannelPermissionRequest,
) (*pb.ChannelPermissionResponse, error) {
	userIDStr := req.GetUserId()
	channelIDStr := req.GetChannelId()
	permission := req.GetPermission()

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return &pb.ChannelPermissionResponse{Granted: false, Reason: "invalid_user_id"}, nil
	}
	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		return &pb.ChannelPermissionResponse{Granted: false, Reason: "invalid_channel_id"}, nil
	}

	// Получаем роль юзера в чате/канале.
	role, err := s.memberRepo.GetMemberRole(ctx, channelID, userID)
	if err != nil {
		// GetMemberRole возвращает error если юзер не member или чат не существует.
		// Мы не можем точно отличить причины, поэтому возвращаем not_a_member
		// (это самый частый кейс).
		log.Debug().Err(err).
			Str("user_id", userIDStr).
			Str("channel_id", channelIDStr).
			Msg("CheckChannelPermission: not a member or channel not found")
		return &pb.ChannelPermissionResponse{Granted: false, Reason: "not_a_member"}, nil
	}

	// Маппим permission на роли.
	granted, reason := evaluatePermission(role, permission)

	log.Debug().
		Str("user_id", userIDStr).
		Str("channel_id", channelIDStr).
		Str("permission", permission).
		Str("role", role).
		Bool("granted", granted).
		Msg("CheckChannelPermission")

	return &pb.ChannelPermissionResponse{
		Granted: granted,
		Reason:  reason,
	}, nil
}

// evaluatePermission применяет логику маппинга permission → роли.
// Использует константы chat.RoleOwner / chat.RoleAdmin / chat.RoleSubscriber.
func evaluatePermission(role, permission string) (granted bool, reason string) {
	switch permission {
	case "OWNER":
		if role == chat.RoleOwner {
			return true, "ok"
		}
		return false, "not_owner"

	case "OWNER_OR_ADMIN":
		if role == chat.RoleOwner || role == chat.RoleAdmin {
			return true, "ok"
		}
		return false, "not_admin"

	case "MEMBER":
		// Если роль есть (любая) — юзер уже member.
		return true, "ok"

	default:
		return false, "invalid_permission"
	}
}

// ────────────────────────────────────────────────────────────
//  4. CheckChatMembership (batch для transfer-in-chat)
// ────────────────────────────────────────────────────────────

// CheckChatMembership — batch-проверка членства группы юзеров в чате.
//
// Используется в transaction-service перед broadcast'ом transfer-in-chat,
// чтобы убедиться что и отправитель, и получатель — участники этого чата.
//
// Возвращает результат для каждого user_id в том же порядке. Если чат
// не существует — chat_exists=false и all members.is_member=false.
func (s *TransactionAccessServer) CheckChatMembership(
	ctx context.Context, req *pb.ChatMembershipRequest,
) (*pb.ChatMembershipResponse, error) {
	chatIDStr := req.GetChatId()
	userIDStrs := req.GetUserIds()

	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		log.Debug().Str("chat_id", chatIDStr).Msg("CheckChatMembership: invalid chat_id")
		return &pb.ChatMembershipResponse{
			ChatExists: false,
			Members:    buildEmptyMembers(userIDStrs),
		}, nil
	}

	if len(userIDStrs) == 0 {
		return &pb.ChatMembershipResponse{ChatExists: true, Members: nil}, nil
	}

	// Проверяем каждого юзера через memberRepo.
	// memberRepo.GetMemberRole возвращает (role, nil) если member, (_, err) иначе.
	// Мы не отличаем «чата нет» от «юзер не member» — оба случая = is_member=false.
	// Это нормально для нашего use case.
	members := make([]*pb.UserMembership, 0, len(userIDStrs))

	for _, userIDStr := range userIDStrs {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			members = append(members, &pb.UserMembership{
				UserId:   userIDStr,
				IsMember: false,
			})
			continue
		}

		role, err := s.memberRepo.GetMemberRole(ctx, chatID, userID)
		if err != nil {
			// Юзер не в чате (или чата нет). Не различаем — это OK для use case.
			members = append(members, &pb.UserMembership{
				UserId:   userIDStr,
				IsMember: false,
			})
			continue
		}

		members = append(members, &pb.UserMembership{
			UserId:   userIDStr,
			IsMember: true,
			Role:     role,
		})
	}

	log.Debug().
		Str("chat_id", chatIDStr).
		Int("users_checked", len(userIDStrs)).
		Msg("CheckChatMembership")

	return &pb.ChatMembershipResponse{
		ChatExists: true,
		Members:    members,
	}, nil
}

// ────────────────────────────────────────────────────────────
//  5. GetUsersByIDs
// ────────────────────────────────────────────────────────────

// GetUsersByIDs — batch-резолв user_id → username + display_name.
// Используется transaction-service для treasury-статистики (top-донатеры).
// Невалидные/несуществующие id просто отсутствуют в ответе.
func (s *TransactionAccessServer) GetUsersByIDs(
	ctx context.Context, req *pb.GetUsersByIDsRequest,
) (*pb.GetUsersByIDsResponse, error) {
	idStrs := req.GetUserIds()
	if len(idStrs) == 0 {
		return &pb.GetUsersByIDsResponse{}, nil
	}

	ids := make([]uuid.UUID, 0, len(idStrs))
	for _, s := range idStrs {
		if id, err := uuid.Parse(s); err == nil {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return &pb.GetUsersByIDsResponse{}, nil
	}

	users, err := s.userRepo.GetUsersByIDs(ctx, ids)
	if err != nil {
		log.Error().Err(err).Int("ids", len(ids)).Msg("GetUsersByIDs: db error")
		return nil, err
	}

	out := make([]*pb.UserBrief, 0, len(users))
	for _, u := range users {
		out = append(out, &pb.UserBrief{
			UserId:      u.ID.String(),
			Username:    u.Username,
			DisplayName: u.DisplayName,
		})
	}
	return &pb.GetUsersByIDsResponse{Users: out}, nil
}

// buildEmptyMembers возвращает массив с is_member=false для всех user_ids.
// Используется когда невалидный chat_id — все членства = false.
func buildEmptyMembers(userIDStrs []string) []*pb.UserMembership {
	out := make([]*pb.UserMembership, 0, len(userIDStrs))
	for _, uid := range userIDStrs {
		out = append(out, &pb.UserMembership{
			UserId:   uid,
			IsMember: false,
		})
	}
	return out
}

// createAndBroadcastSystemMessage:
//  1. Парсит chatID
//  2. Делает SaveMessage с type=msgType, content=payloadJSON, sender_id=NULL
//  3. Broadcast'ит NEW_MESSAGE всем участникам чата через WS
//
// Возвращает ID созданного сообщения как строку.
func (s *TransactionAccessServer) createAndBroadcastSystemMessage(
	ctx context.Context,
	triggerUserID uuid.UUID,
	chatIDStr, msgType, payloadJSON string,
) (string, error) {
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		return "", fmt.Errorf("invalid chat_id: %w", err)
	}

	// content для system message = payload JSON. UI парсит его в зависимости от type.
	msg := &chat.Message{
		ChatID:   chatID,
		SenderID: nil, // system message — без отправителя
		Content:  payloadJSON,
		Type:     msgType,
		// Status, ID, CreatedAt — заполнит SaveMessage из RETURNING.
	}

	if err := s.messageRepo.SaveMessage(ctx, msg); err != nil {
		return "", fmt.Errorf("save message: %w", err)
	}

	// Broadcast NEW_MESSAGE всем участникам чата.
	members, err := s.memberRepo.GetChatMembers(ctx, chatID)
	if err != nil {
		// Сообщение уже создано, но broadcast не получился — не critical,
		// юзеры увидят его при следующем GetMessages.
		log.Warn().Err(err).
			Str("chat_id", chatIDStr).
			Int64("message_id", msg.ID).
			Msg("createAndBroadcastSystemMessage: failed to get chat members for broadcast")
		return strconv.FormatInt(msg.ID, 10), nil
	}

	wsMsg := socketmodel.WSMessage{
		Type:    socketmodel.EventNewMessage,
		Payload: msg,
	}
	for _, memberID := range members {
		// SendToUser отдаёт error если юзер оффлайн — это нормально, не ошибка.
		_ = s.hub.SendToUser(memberID, wsMsg)
	}

	log.Info().
		Str("chat_id", chatIDStr).
		Int64("message_id", msg.ID).
		Str("type", msgType).
		Int("broadcast_to", len(members)).
		Msg("system message created and broadcasted")

	return strconv.FormatInt(msg.ID, 10), nil
}
