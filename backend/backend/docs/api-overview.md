# Suda Messenger — API Overview

Карта всех HTTP-эндпоинтов и WebSocket-событий. Источник правды — Swagger UI каждого сервиса; этот документ — для быстрого поиска.

| Сервис | Swagger UI |
|--------|------------|
| Messenger | http://localhost:8081/swagger/index.html |
| Transaction | http://localhost:8082/swagger/index.html |
| Media | http://localhost:8084/swagger/index.html (legacy spec) |

> Пути в этом документе указаны **без gateway-префикса**. Снаружи (через api-gateway на :8080) ко всем messenger-путям добавляется `/api/v1/messenger`, к transaction — `/api/v1/tx`, к media — `/api/v1/media`. Пример: `POST /chats` ⇒ `POST /api/v1/messenger/chats` через gateway.
>
> Все эндпоинты (кроме `/auth/*`) требуют `Authorization: Bearer <access_token>`.

---

## Messenger Service

### Auth

| Method | Path | Description |
|--------|------|-------------|
| POST   | `/auth/register`             | Register a new user |
| POST   | `/auth/verify`               | Verify email with code |
| POST   | `/auth/login`                | Log in with email and password |
| POST   | `/auth/refresh`              | Refresh access and refresh tokens |
| POST   | `/auth/forgot-password`      | Request password reset code |
| POST   | `/auth/reset-password`       | Reset password with verification code |

### User

| Method | Path | Description |
|--------|------|-------------|
| GET    | `/user/me`                   | Get current user profile |
| PUT    | `/user/profile`              | Update user profile |
| PUT    | `/user/avatar`               | Update user avatar by media ID |
| GET    | `/user/id/{id}`              | Get user profile by ID |
| GET    | `/user/id/{id}/status`       | Get user online status |
| POST   | `/user/device`               | Register a device for push notifications |
| POST   | `/user/init-data`            | Generate init data for mini-app integration |
| POST   | `/user/logout`               | Log out and revoke session |

### Users (contacts, block-list)

| Method | Path | Description |
|--------|------|-------------|
| POST   | `/users/block`               | Block a user |
| POST   | `/users/unblock`             | Unblock a user |
| GET    | `/users/blocked`             | Get list of blocked user IDs |
| GET    | `/users/contacts`            | Get all contacts with custom names |
| POST   | `/users/contacts`            | Set a custom name for a contact |
| DELETE | `/users/contacts/{id}`       | Remove a custom contact name |

### Preferences

| Method | Path | Description |
|--------|------|-------------|
| GET    | `/user/preferences`              | Get user preferences |
| PUT    | `/user/preferences`              | Update user preferences |
| POST   | `/user/last-seen-hidden`         | Hide my last-seen from a specific user |
| DELETE | `/user/last-seen-hidden/{id}`    | Unhide my last-seen for a specific user |
| GET    | `/user/last-seen-hidden`         | Get list of users from whom my last-seen is hidden |

### UserPins

| Method | Path | Description |
|--------|------|-------------|
| GET    | `/user/pins`            | Get user pins |
| POST   | `/user/pins`            | Create a user pin |
| PUT    | `/user/pins/reorder`    | Reorder user pins |
| DELETE | `/user/pins/{id}`       | Delete a user pin |

### Chats

| Method | Path | Description |
|--------|------|-------------|
| GET    | `/chats`                | Get all chats for current user |
| POST   | `/chats`                | Create a new chat |
| GET    | `/chats/{id}/info`      | Get detailed chat info with members |
| PUT    | `/chats/{id}`           | Update chat name, description or avatar |
| DELETE | `/chats/{id}`           | Delete a chat |
| GET    | `/chats/{id}/media`     | Get media files shared in chat |

### Chat Members

| Method | Path | Description |
|--------|------|-------------|
| POST   | `/chats/{id}/members`               | Add a member to a group chat |
| DELETE | `/chats/{id}/members/{uid}`         | Remove a member from a group chat |
| PUT    | `/chats/{id}/members/{uid}/role`    | Update a member's role in chat |
| POST   | `/chats/{id}/leave`                 | Leave a group chat |
| POST   | `/chats/{id}/mute`                  | Mute chat notifications |
| POST   | `/chats/{id}/unmute`                | Unmute chat notifications |
| POST   | `/chats/{id}/pin`                   | Pin a message in chat |
| DELETE | `/chats/{id}/pin/{mid}`             | Unpin a message in chat |

### Messages

| Method | Path | Description |
|--------|------|-------------|
| GET    | `/chats/{id}/messages`      | Get messages from a chat |
| POST   | `/chats/{id}/messages`      | Send a message to a chat |
| PUT    | `/messages/{id}`            | Edit a message |
| DELETE | `/messages/{id}`            | Delete a message |
| GET    | `/messages/{id}/readers`    | Get users who read a message |
| POST   | `/chats/forward`            | Forward a message to another chat |

### Channels

| Method | Path | Description |
|--------|------|-------------|
| POST   | `/channels/{id}/subscribe`                            | Subscribe to a channel (token-gating enforced) |
| POST   | `/channels/{id}/unsubscribe`                          | Unsubscribe from a channel |
| GET    | `/channels/{id}/subscribers`                          | Get channel subscribers |
| GET    | `/channels/by-username/{username}`                    | Resolve channel by username |
| PUT    | `/channels/{id}/settings`                             | Update channel settings (OWNER/ADMIN) |
| GET    | `/channels/{id}/posts/{msg_id}/comments`              | Get comments for a channel post |
| POST   | `/channels/{id}/posts/{msg_id}/comments`              | Add a comment to a channel post |
| PUT    | `/channels/comments/{comment_id}`                     | Edit a comment |
| DELETE | `/channels/comments/{comment_id}`                     | Delete a comment |
| GET    | `/channels/{id}/apps`                                 | Get apps linked to a channel |
| POST   | `/channels/{id}/apps`                                 | Link an app to a channel |
| DELETE | `/channels/{id}/apps/{app_id}`                        | Unlink an app from a channel |

### Search

| Method | Path | Description |
|--------|------|-------------|
| GET    | `/search`                       | Search users, chats and messages globally |
| GET    | `/chats/{id}/search`            | Search messages within a specific chat |

---

## Transaction Service

### Wallet

| Method | Path | Description |
|--------|------|-------------|
| GET    | `/wallet/me`                                  | Get my wallet (address + on-chain SUDA balance) |
| POST   | `/wallet/resolve`                             | Resolve username → user_id + wallet_address |
| POST   | `/wallet/transfer`                            | Transfer SUDA to a user (optionally in chat) |
| GET    | `/wallet/history`                             | Get my transaction history (paginated) |
| GET    | `/wallet/channel/{channel_id}`                | Get channel treasury wallet (OWNER/ADMIN) |

### Treasury

| Method | Path | Description |
|--------|------|-------------|
| GET    | `/wallet/channel/{channel_id}/treasury`       | Channel treasury statistics (top donors + totals) |
| GET    | `/wallet/channel/{channel_id}/donations`      | Paginated list of channel donations |

### Purchase

| Method | Path | Description |
|--------|------|-------------|
| GET    | `/purchase/packages`              | List SUDA purchase packages |
| POST   | `/purchase/initiate`              | Initiate a SUDA purchase (PENDING) |
| POST   | `/purchase/{id}/confirm`          | Confirm a SUDA purchase (treasury → user) |
| GET    | `/purchase/history`               | Get my purchase history |

### Donation

| Method | Path | Description |
|--------|------|-------------|
| POST   | `/donate`                         | Donate SUDA to a user or channel |

### Gating

| Method | Path | Description |
|--------|------|-------------|
| POST   | `/gating/rule`                    | Create or update a gating rule (OWNER) |
| GET    | `/gating/rule/{chat_id}`          | Get gating rule for a channel |
| DELETE | `/gating/rule/{chat_id}`          | Delete gating rule (OWNER) |

---

## Media Service

| Method | Path | Description |
|--------|------|-------------|
| POST   | `/media/upload/url`               | Get a presigned URL for upload |
| POST   | `/media/{id}/confirm`             | Confirm completion of an upload |
| GET    | `/media/{id}`                     | Get media metadata + presigned download URL |
| DELETE | `/media/{id}`                     | Delete media (owner only) |

> Конкретные пути — в Swagger UI media-service. Список выше — типичный для S3-presigned flow.

---

## WebSocket Events (observer → client)

Все события приходят на WebSocket-соединение клиента (`/ws/`, JWT-авторизация). Полезная нагрузка — `payload_json` со стандартным набором полей: `tx_hash`, `from_address`, `to_address`, `amount_wei` + специфичные для типа.

| event_type | Кто получает | Когда отправляется | System message в чате? |
|------------|--------------|---------------------|------------------------|
| `SUDA_RECEIVED`       | получатель | Подтверждение P2P-перевода SUDA. Если transfer-in-chat — `target_chat_id` указан. | да (`SUDA_TRANSFER`), если transfer-in-chat |
| `SUDA_SENT`           | отправитель | То же событие, но для sender'а. | нет |
| `DONATION_SENT`       | донатер | Подтверждение доната (юзеру или каналу). | да (`DONATION`), если задан `chat_id` (P2P) или это донат каналу |
| `DONATION_RECEIVED`   | получатель | Подтверждение P2P-доната (не канал). | нет |
| `PURCHASE_COMPLETED`  | покупатель | Подтверждение покупки SUDA (treasury → user). `SUDA_RECEIVED` НЕ шлётся отдельно. | нет |

---

## Конвенция HTTP-ответов

**Success envelope** (для action endpoints — register, login, delete и т.п.):
```json
{ "action": "USER_CREATED", "message": "User registered, verification email sent" }
```

**Error envelope**:
```json
{ "error": "INVALID_REQUEST", "message": "amount_wei must be positive integer string" }
```

Resource endpoints (GET-ы возвращающие данные) отдают объект напрямую без envelope:
```json
{ "id": "...", "username": "alice", "display_name": "Alice", ... }
```

---

## Конвенция кодов ошибок

| HTTP | Когда |
|------|-------|
| 400 | invalid_request / validation_error |
| 401 | unauthorized (JWT, gateway signature) |
| 402 | INSUFFICIENT_BALANCE (token-gating) |
| 403 | forbidden / not_owner / sender_not_in_chat / NFT_REQUIRED |
| 404 | resource not found |
| 409 | already_handled (purchase double-confirm), already_taken (username) |
| 424 | wallet_not_found (failed dependency — нужен кошелёк) |
| 429 | rate_limited |
| 500 | internal_error |
