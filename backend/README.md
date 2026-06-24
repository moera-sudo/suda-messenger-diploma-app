# Suda Messenger — технический отчет по проекту

Suda Messenger — микросервисный backend мессенджера с realtime-чатами, медиа-хранилищем, WebSocket-уведомлениями и встроенной Web3-частью на приватной EVM-сети Hyperledger Besu. Проект объединяет классические функции мессенджера: авторизацию, пользователей, личные и групповые чаты, каналы, сообщения, медиа, поиск, настройки приватности, pinned-сущности, friends/request flow, а также custodial Web3-кошелек, SUDA token transfers, донаты, покупки SUDA-пакетов и token-gating каналов.

Репозиторий устроен как monorepo:

- `backend/api-gateway` — единая HTTP-точка входа для REST API.
- `backend/messenger-service` — основной домен мессенджера: пользователи, auth, чаты, каналы, WebSocket, social-функции.
- `backend/media-service` — upload/download медиа через presigned URL и MinIO/S3.
- `backend/transaction-service` — Web3-кошельки, SUDA transfers, donations, purchases, gating, Besu observer.
- `backend/proto` — gRPC/protobuf контракты между сервисами.
- `contracts` — Solidity-контракты и Hardhat deployment.
- `migrations` — SQL-миграции PostgreSQL для messenger/media/transaction данных.
- `web/wallet` — Wallet WebView клиент на React/Vite.
- `bruno/Suda Messenger` — Bruno collection для ручного и сценарного тестирования API/WS.
- `besu/network` — genesis/QBFT конфигурация приватной Besu-сети.
- `nginx` — production reverse proxy.

## 1. Архитектура

### 1.1 Компоненты

| Компонент | Технологии | Назначение |
|---|---|---|
| Nginx | nginx:alpine | Внешняя HTTP-точка в production, прокси для API, WebSocket, Swagger и Wallet WebView. |
| API Gateway | Go, Echo | JWT-auth, HMAC-подпись downstream-запросов, reverse proxy на сервисы. |
| Messenger Service | Go, Echo, gRPC, Melody WS | Основная бизнес-логика мессенджера и realtime-слой. |
| Media Service | Go, Echo, gRPC, MinIO/S3 | Метаданные медиа, presigned upload/view/download URLs, ACL для приватных файлов. |
| Transaction Service | Go, Echo, gRPC, go-ethereum | Custodial wallets, SUDA token operations, Besu observer, purchase/donation/gating. |
| PostgreSQL | postgres:16-alpine | Основная реляционная БД всех backend-доменов. |
| Redis | redis:7-alpine | Online-status, кеши/локи, realtime-adjacent состояние. |
| MinIO | S3-compatible storage | Объектное хранилище медиа. |
| Hyperledger Besu | private EVM chain | Приватная EVM-сеть для SUDA token и контрактов проекта. |
| Contracts Deploy | Hardhat | One-shot deployment Solidity-контрактов в Besu. |
| Wallet Frontend | React, Vite, Zustand | WebView wallet client на `/wallet/`. |
| Portainer | portainer-ce | Web UI для управления Docker stack в deployment. |

### 1.2 Высокоуровневый поток запросов

Production-вход:

```text
Client / WebView / Mobile App
        |
        v
      Nginx :80
        |
        +--> /api/...              -> api-gateway:8080
        +--> /ws                   -> messenger-service:8081/ws
        +--> /wallet/              -> wallet-frontend:80
        +--> /swagger/messenger/   -> messenger-service:8081/swagger/
        +--> /swagger/transaction/ -> transaction-service:8082/swagger/
```

Gateway-внутри:

```text
/api/v1/messenger/* -> messenger-service:8081 /*
/api/v1/tx/*        -> transaction-service:8082 /*
/api/v1/media/*     -> media-service:8084 /api/media/*
```

Сервисные gRPC-связи:

```text
messenger-service -> media-service        MediaService
messenger-service -> transaction-service  TransactionService
media-service     -> messenger-service    MessengerAccessService
transaction-service -> messenger-service  TransactionAccessService
```

## 2. Deployment

### 2.1 Production stack

Основной production stack описан в `docker-compose.yml`. Он поднимает:

- `postgres`, `redis`, `minio`;
- `besu-node`;
- one-shot `migrate`;
- one-shot `contracts-deploy`;
- `api-gateway`, `messenger-service`, `transaction-service`, `media-service`;
- `wallet-frontend`;
- `nginx`;
- `portainer`.

Запуск:

```bash
docker compose -f docker-compose.yml up -d --build
```

Через Taskfile:

```bash
task docker-up-prod
```

Production compose не использует `docker-compose.override.yml`, поэтому наружу открыты только:

- `80` — Nginx;
- `9443` — Portainer HTTP UI, который предполагается держать за ZeroTier/VPN.

Внутренние сервисные порты доступны только внутри Docker network `app-network`.

### 2.2 Development stack

Обычный `docker compose up -d --build` автоматически применяет `docker-compose.override.yml`, который открывает внутренние порты на host:

| Сервис | Host port |
|---|---:|
| PostgreSQL | `55432` |
| Redis | `6379` |
| MinIO S3 API | `9001` |
| MinIO Console | `9091` |
| Besu RPC | `8545` |
| API Gateway | `8080` |
| Messenger HTTP/gRPC | `8081` / `9081` |
| Transaction HTTP/gRPC | `8082` / `9082` |
| Media HTTP/gRPC | `8084` / `9094` |

Запуск dev stack:

```bash
task docker-up
```

### 2.3 Secrets

Docker secrets лежат в `secrets/`:

- `wallet_encryption_key.txt` — AES-256-GCM key для приватных ключей кошельков;
- `message_encryption_key.txt` — AES-256-GCM key для at-rest шифрования содержимого сообщений (hex, 32 байта);
- `jwt_secret.txt` — HMAC secret для access/refresh JWT;
- `gateway_signature_secret.txt` — HMAC secret для подписи gateway -> services;
- `init_data_sign_token.txt` — token для init-data flow;
- `eth_private_key.txt` — treasury EVM private key;
- `firebase-adminsdk.json` — Firebase service account для push-notifications.

Генерация базовых secrets:

```bash
task secrets-generate
```

После генерации вручную добавляются treasury private key и Firebase JSON. Сервисы читают secrets через `*_FILE` переменные; при Docker Compose запуске secrets не берутся напрямую из `.env`, если есть file-based вариант.

### 2.4 Migrations

Миграции выполняет one-shot контейнер `migrate`:

```text
migrate/migrate:v4.17.0
  -path=/migrations
  -database=postgres://...
  up
```

Все Go-сервисы зависят от успешного завершения миграций. Это защищает сервисы от старта на неподготовленной схеме.

Полезные команды:

```bash
task migrate-new name=create_indexes
task migrate-up
task migrate-down
task migrate-version
task docker-migrate-up
```

### 2.5 Contracts deploy

`contracts-deploy` — one-shot контейнер на базе Hardhat. Он подключается к `besu-node:8545`, деплоит контракты в порядке:

1. `SudaToken`
2. `SudaNFT`
3. `SudaMarketplace`
4. `SudaEscrow`
5. `SudaFundraising`

Адреса сохраняются в deployment artifacts и могут попадать в `secrets/contracts.env`. Backend на текущем реализованном уровне использует бизнес-логику вокруг `SudaToken`: balance checks, transfer, purchases, donations и observer.

### 2.6 Besu

Besu запускается как приватная EVM-сеть:

- image: `hyperledger/besu:25.12.0`;
- `chainId=1337`;
- QBFT consensus;
- `blockPeriodSeconds=2`;
- `zeroBaseFee=true`;
- `--min-gas-price=0`;
- HTTP RPC включен внутри Docker network;
- RPC API: `ETH,NET,QBFT,WEB3`.

Genesis и validator key лежат в `besu/network`. Данные chain state хранятся в Docker volume `besu_data`.

### 2.7 Nginx

Nginx слушает `:80` и работает как внешний reverse proxy:

- `/api/` -> `api-gateway:8080`;
- `/ws` и `/ws/` -> `messenger-service:8081/ws`;
- `/swagger/messenger/` -> `messenger-service:8081/swagger/`;
- `/swagger/transaction/` -> `transaction-service:8082/swagger/`;
- `/wallet/` -> `wallet-frontend:80/`;
- `/` возвращает health-like текст `Suda backend is up`.

WebSocket проксируется напрямую в messenger-service, потому что browser/WebView WebSocket API не позволяет надежно поставить custom headers в handshake. JWT передается через `?token=<jwt>` и валидируется самим WS Hub.

### 2.8 ZeroTier и Portainer

Deployment guide рассчитан на сервер без публичного IP. Доступ к `http://<server-zt-ip>/...` и `http://<server-zt-ip>:9443/` идет через ZeroTier VPN.

Portainer монтирует `/var/run/docker.sock`, поэтому фактически получает root-доступ к Docker engine. Его порт должен быть доступен только через доверенную сеть/VPN.

## 3. API Gateway

API Gateway — первый backend-слой для REST API. Он делает:

- CORS;
- request logging;
- JWT validation;
- пропуск public routes;
- добавление user context headers;
- HMAC-подпись запроса к downstream сервисам;
- reverse proxy.

### 3.1 Auth middleware

Gateway ищет JWT:

- в `Authorization: Bearer <token>`;
- в query parameter `?token=<jwt>` для WS-compatible сценариев.

Public routes:

- `/api/v1/messenger/auth/*`;
- `/ws` и `/ws/*`;
- `/health`;
- `/swagger/*`.

После успешной проверки access token gateway добавляет:

- `X-User-ID`;
- `X-User-Role`;
- `X-Session-ID`, если в claims есть `sid`;
- `X-Gateway-Timestamp`;
- `X-Gateway-Path`;
- `X-Gateway-Signature`.

`X-Gateway-Signature` считается как HMAC-SHA256 от строки:

```text
userID:role:timestamp:path
```

Downstream сервисы используют gateway middleware, чтобы принимать защищенные REST-запросы только от gateway.

### 3.2 Swagger aggregation

Gateway содержит Swagger aggregator на `/swagger/doc.json` и `/swagger/*`. Сейчас агрегируются messenger и media specs. Transaction Swagger доступен напрямую через Nginx `/swagger/transaction/` и напрямую на service port в dev.

## 4. Messenger Service

Messenger Service — центральный домен мессенджера. Он предоставляет HTTP API, WebSocket endpoint и два gRPC server контракта:

- `MessengerAccessService` для media ACL;
- `TransactionAccessService` для transaction-service callbacks, username resolving и permission checks.

### 4.1 Основные подсистемы

| Подсистема | Что делает |
|---|---|
| Auth | Register, verify email, login, refresh, forgot/reset password, webapp session. |
| User | `me`, profile, avatar, logout, devices, status, sessions, init-data. |
| Chat | DIRECT/GROUP/SAVED/CHANNEL creation, list, info, update, delete (scope=me/everyone), media overview. |
| Messages | Send, get messages, edit, delete, readers, forward. |
| Members | Add/remove members, roles, leave, mute/unmute, pin/unpin. |
| Moderation | Block/unblock users, blocked list, contact names. |
| Channels | Subscribe/unsubscribe (платная и бесплатная), private join requests, invites, comments, linked apps, settings. |
| Preferences | User privacy settings, hidden last-seen list. |
| Search | Global search and chat-local search. |
| User pins | User-level pinned chats с сортировкой по дате закрепления. |
| Friends | Friend requests, accept/reject/cancel, unfriend, relationship status. |
| WebSocket | Realtime messages, typing, reads, online/offline and domain events. |
| At-rest encryption | AES-256-GCM шифрование содержимого сообщений DIRECT и GROUP чатов. |

### 4.2 At-rest encryption

Пакет `backend/messenger-service/internal/platform/crypto` реализует шифрование содержимого сообщений на уровне репозитория.

Интерфейс `ContentCipher`:

- `EncryptContent(plaintext) string` — возвращает ciphertext с sentinel-prefix `enc:v1:`.
- `DecryptContent(stored) string` — если значение начинается с `enc:v1:`, расшифровывает; иначе возвращает как есть.

Реализации:

- `aesContentCipher` — AES-256-GCM; случайный nonce на каждое сообщение; key version в sentinel позволяет ротировать ключ.
- `noopContentCipher` — identity-преобразование; используется если `MESSAGE_ENCRYPTION_KEY` не задан (dev-режим).

Поведение:

- Шифруются тела сообщений только в чатах типа `DIRECT` и `GROUP`. Системные сообщения и чаты типа `CHANNEL`/`SAVED` не шифруются.
- `search_vector` строится из plaintext в application-коде до шифрования, поэтому полнотекстовый поиск работает независимо от включения шифрования.
- Сообщения без sentinel-prefix (`enc:v1:`) читаются как plaintext — обратная совместимость без backfill.
- При ошибке шифрования сообщение сохраняется в открытом виде; при ошибке расшифровки — возвращается сохранённое значение.

### 4.3 HTTP routes

Через gateway все routes доступны с prefix `/api/v1/messenger`.

Auth:

```text
POST /auth/register
POST /auth/verify
POST /auth/login
POST /auth/refresh
POST /auth/forgot-password
POST /auth/reset-password
POST /auth/webapp-session
```

User:

```text
GET    /user/me
GET    /user/id/:id
GET    /user/id/:id/status
POST   /user/logout
PUT    /user/profile
PUT    /user/avatar
POST   /user/device
POST   /user/init-data
PUT    /user/password
GET    /user/sessions
DELETE /user/sessions/:id
POST   /user/sessions/terminate-others
```

Chats/messages:

```text
POST   /chats
GET    /chats
GET    /chats/:id/info
PUT    /chats/:id
DELETE /chats/:id
GET    /chats/:id/media
POST   /chats/:id/messages
GET    /chats/:id/messages
POST   /chats/forward
PUT    /messages/:id
DELETE /messages/:id
GET    /messages/:id/readers
```

Members/moderation:

```text
POST   /chats/:id/members
DELETE /chats/:id/members/:uid
PUT    /chats/:id/members/:uid/role
POST   /chats/:id/leave
POST   /chats/:id/mute
POST   /chats/:id/unmute
POST   /chats/:id/pin
DELETE /chats/:id/pin/:mid
POST   /users/block
POST   /users/unblock
GET    /users/blocked
POST   /users/contacts
DELETE /users/contacts/:id
GET    /users/contacts
```

Channels:

```text
POST   /channels/:id/subscribe
POST   /channels/:id/unsubscribe
GET    /channels/:id/subscribers
GET    /channels/by-username/:username
GET    /channels/:id/view
GET    /channels/:id/settings
PUT    /channels/:id/settings
POST   /channels/:id/join-request
DELETE /channels/:id/join-request
GET    /channels/:id/join-requests
POST   /channels/:id/join-requests/:uid/approve
POST   /channels/:id/join-requests/:uid/reject
GET    /channels/invites
POST   /channels/:id/invites
DELETE /channels/:id/invites/:uid
POST   /channels/:id/invites/accept
POST   /channels/:id/invites/decline
GET    /channels/:id/posts/:msg_id/comments
POST   /channels/:id/posts/:msg_id/comments
PUT    /channels/comments/:comment_id
DELETE /channels/comments/:comment_id
GET    /channels/:id/apps
POST   /channels/:id/apps
DELETE /channels/:id/apps/:app_id
```

Preferences, search, pins, friends:

```text
GET    /preferences
PUT    /preferences
POST   /last-seen-hidden
DELETE /last-seen-hidden/:id
GET    /last-seen-hidden

GET /search
GET /chats/:id/search

POST   /user-pins
DELETE /user-pins/:id
PUT    /user-pins/reorder
GET    /user-pins

GET    /friends
DELETE /friends/:user_id
GET    /friends/:user_id/status
POST   /friends/requests
GET    /friends/requests
POST   /friends/requests/:id/accept
POST   /friends/requests/:id/reject
DELETE /friends/requests/:id
```

## 5. Чаты

### 5.1 Типы чатов

`messenger_chats.type` поддерживает:

| Type | Назначение |
|---|---|
| `DIRECT` | Личный чат между двумя пользователями. Создание должно быть идемпотентным: повторный запрос на тот же pair возвращает существующий чат. Учитываются блокировки пользователей. |
| `GROUP` | Групповой чат с owner/admin/member ролями. Поддерживает управление участниками, роли, mute, leave, pin/unpin, системные сообщения. |
| `CHANNEL` | Канал на базе общей chat-модели. Имеет visibility `PUBLIC`/`PRIVATE`, username, subscriber count, comments, join requests, invites, linked apps и token-gating. |
| `SAVED` | Личное хранилище сообщений пользователя. Используется как Saved Messages. |

### 5.2 Роли

Базовые роли участников:

- `OWNER` — владелец чата/канала;
- `ADMIN` — администратор;
- `MEMBER` — обычный участник;
- `SUBSCRIBER` — роль подписчика канала.

### 5.3 Сообщения

Типы сообщений:

- `TEXT` — обычный текст;
- `IMAGE`, `VIDEO`, `AUDIO`, `FILE`, `VOICE` — сообщения с media attachment;
- `SYSTEM` — системное сообщение о действии;
- `SUDA_TRANSFER` — системный/специализированный тип для transfer-карточек;
- `DONATION` — системный/специализированный тип для donation-карточек.

Статусы:

- `SENT`;
- `READ`.

Системные actions:

- `USER_JOINED`;
- `USER_LEFT`;
- `USER_REMOVED`;
- `GROUP_CREATED`;
- `GROUP_RENAMED`;
- `AVATAR_CHANGED`;
- `MESSAGE_PINNED`;
- `MESSAGE_UNPINNED`.

### 5.4 Алгоритм отправки сообщения

Сообщение можно отправить двумя путями:

- HTTP: `POST /chats/:id/messages`;
- WebSocket: event `SEND_MESSAGE`.

Общий алгоритм:

1. Handler/router получает `SendMessageReq`: `chat_id`, `content`, `type`, optional `attachment_media_id`, optional `reply_to_message_id`, `client_side_id`.
2. `ChatService.SendMessage` вызывает `canSendMessage`.
3. Проверяется, что пользователь может писать в чат. Для обычных чатов это membership. Для direct учитываются блокировки. Для каналов права зависят от channel/member rules.
4. Сообщение сохраняется в `messenger_messages`.
5. Если есть `attachment_media_id`, messenger-service асинхронно вызывает gRPC `MediaService.LinkMediaToEntity(media_id, "MESSAGE", message_id)`.
6. Через `broadcastMessageWithPush` сообщение доставляется онлайн-участникам как WS event `NEW_MESSAGE`.
7. Для offline-получателей, кроме отправителя, используется push notifier. Если preview отключен настройками приватности, push отправляется без текста сообщения.
8. История читается через `GET /chats/:id/messages`; для public channels чтение разрешено без подписки, для остальных типов нужен доступ.

### 5.5 Reply, edit, delete, forward

- Reply хранится через `reply_to_message_id`.
- Edit разрешен после проверки membership и прав автора/правил чата, затем рассылается `MESSAGE_EDITED`.
- Delete сообщения поддерживает два режима: hide for me через `messenger_message_hidden` и delete for everyone через `deleted_at`, после чего рассылается `MESSAGE_DELETED`.
- Delete чата (`DELETE /chats/:id`) принимает query param `?scope=me|everyone`:
  - `scope=me`: записывает `cleared_at` и `hidden_at` в `messenger_chat_members` для вызывающего — история очищается и чат скрывается из списка; рассылает `CHAT_DELETED {scope:"me"}` только устройствам вызывающего.
  - `scope=everyone` (default): полное каскадное удаление чата и всех сообщений для всех участников; рассылает `CHAT_DELETED {scope:"everyone"}` перед удалением.
- Forward проверяет доступ к source chat и право писать в target chat. Пересланное сообщение сохраняет `forwarded_from_chat` и `forwarded_from_msg`. Если у original было media attachment, media повторно линкуется к новому message id.

### 5.6 Read receipts

`READ_MESSAGE` по WS или service-вызов обновляет read cursor участника (`last_read_message_id`) и рассылает `MESSAGES_READ` остальным участникам чата.

### 5.7 Pinned messages

Pin/unpin требует соответствующих прав. Для группы это owner/admin-like действие. После pin/unpin:

- пишется/удаляется запись в `messenger_pinned_messages`;
- создается system message;
- рассылается `MESSAGE_PINNED` или `MESSAGE_UNPINNED`.

### 5.8 Каналы

Канал реализован как `messenger_chats` с `type='CHANNEL'`.

Особенности:

- `visibility=PUBLIC` — канал можно читать публично, но write/admin actions требуют роли.
- `visibility=PRIVATE` — доступ к чтению и комментариям требует membership/subscription.
- `username` позволяет найти канал через `/channels/by-username/:username`.
- `comments_enabled` включает комментарии к channel posts.
- `subscriber_count` хранит счетчик подписчиков.

Private channel flow:

- пользователь отправляет `POST /channels/:id/join-request`;
- администраторы видят pending requests;
- admin/owner вызывает approve/reject;
- также поддерживается invite flow: admin создает invite, пользователь принимает или отклоняет.

Token gating / платная подписка:

- Владелец задаёт правило через `POST /gating/rule` с полем `subscription_price_wei`.
- Если `subscription_price_wei > 0`: подписка платная — при вызове `subscribe` сервис через gRPC `ChargeChannelSubscription` списывает указанную сумму SUDA с кошелька пользователя в treasury канала; при успехе пользователь немедленно добавляется как подписчик.
- Повторный subscribe идемпотентен: уже подписанный пользователь не проходит через списание повторно.
- Не-подписчик не может читать платный PUBLIC-канал.
- Если правила нет или `subscription_price_wei = 0`: PUBLIC-канал открыт для чтения и бесплатной подписки без транзакций.
- NFT-requirement предусмотрен в модели данных (`tx_gating_rules`), но полноценный NFT business-flow не реализован.

Comments:

- комментарии хранятся в `messenger_channel_post_comments`;
- комментировать может подписчик, если comments enabled;
- author может редактировать свой комментарий;
- удалять может author или owner/admin канала;
- новый комментарий рассылается как `CHANNEL_NEW_COMMENT`.

## 6. WebSocket

### 6.1 Endpoint и auth

Endpoint:

```text
GET /ws?token=<access_jwt>
```

В production Nginx проксирует `/ws` напрямую в `messenger-service:8081/ws`. Gateway тоже имеет `/ws` proxy group, но основной production path идет через Nginx напрямую.

WS Hub построен на `github.com/olahol/melody`. Максимальный размер сообщения — 1 MB.

Auth priority в Hub:

1. `X-User-ID` header, если соединение пришло из server-to-server/gateway сценария.
2. `?token=<jwt>`, если header отсутствует.

JWT валидируется локально тем же `jwt_secret`. Принимаются access tokens с `type=access` или без явного type. User id берется из `user_id` или `sub`.

### 6.2 Online/offline

При connect:

- session сохраняется в `userSessions[userID]`;
- Redis key `user:<uuid>:online` ставится в `"1"` на 10 минут;
- вызывается `onConnect`, через который user feature может разослать `USER_ONLINE`.

При disconnect:

- session удаляется;
- если у пользователя не осталось активных WS sessions, Redis online key удаляется;
- вызывается `onDisconnect`, через который обновляется last seen и рассылается `USER_OFFLINE`.

Один пользователь может иметь несколько одновременных WS sessions.

### 6.3 Формат события

Client -> server:

```json
{
  "type": "SEND_MESSAGE",
  "payload": {
    "chat_id": "uuid",
    "content": "hello",
    "type": "TEXT",
    "client_side_id": "local-id"
  }
}
```

Server -> client:

```json
{
  "type": "NEW_MESSAGE",
  "payload": {
    "id": 123,
    "chat_id": "uuid"
  }
}
```

### 6.4 Client -> server events

| Event | Payload | Что делает |
|---|---|---|
| `SEND_MESSAGE` | `SendMessageReq` | Отправляет сообщение через `ChatService.SendMessage`. |
| `EDIT_MESSAGE` | `message_id`, `content` | Редактирует сообщение. |
| `DELETE_MESSAGE` | `message_id`, `for_everyone` | Удаляет для всех или скрывает для текущего пользователя. |
| `READ_MESSAGE` | `chat_id`, `message_id` | Обновляет read cursor. |
| `TYPING` | `chat_id` | Рассылает typing notification. |
| `FORWARD_MESSAGE` | `ForwardMessageReq` | Пересылает сообщение. |
| `REACTION` | `message_id`, `emoji` | В router принято, но текущая реализация только логирует событие как not implemented. |

### 6.5 Server -> client events

Основные события:

- `NEW_MESSAGE`;
- `MESSAGE_EDITED`;
- `MESSAGE_DELETED`;
- `MESSAGES_READ`;
- `TYPING`;
- `USER_ONLINE`;
- `USER_OFFLINE`;
- `MEMBER_ADDED`;
- `MEMBER_REMOVED`;
- `MEMBER_LEFT`;
- `CHAT_UPDATED`;
- `CHAT_DELETED`;
- `MESSAGE_PINNED`;
- `MESSAGE_UNPINNED`;
- `ERROR`;
- `PREFERENCES_UPDATED`;
- `CHANNEL_NEW_POST`;
- `CHANNEL_NEW_COMMENT`;
- `CHANNEL_SUBSCRIBED`;
- `CHANNEL_UNSUBSCRIBED`;
- `CHANNEL_JOIN_REQUEST`;
- `CHANNEL_JOIN_APPROVED`;
- `CHANNEL_JOIN_REJECTED`;
- `CHANNEL_INVITE`;
- `CHANNEL_INVITE_ACCEPTED`;
- `FRIEND_REQUEST_RECEIVED`;
- `FRIEND_REQUEST_ACCEPTED`;
- `FRIEND_REQUEST_REJECTED`;
- `FRIEND_REQUEST_CANCELLED`;
- `FRIEND_REMOVED`.

Transaction Service через `TransactionAccessService.NotifyUserEvent` также инициирует user events для wallet/transfer/donation/purchase сценариев. Они доставляются через messenger WS Hub и могут дополнительно создавать system message в чате.

## 7. Media Service

Media Service отделяет binary upload от сообщений. Backend не принимает файл целиком через messenger-service: клиент получает presigned URL и загружает объект напрямую в MinIO/S3.

### 7.1 HTTP routes

Через gateway routes доступны с prefix `/api/v1/media`, который переписывается в `/api/media` на media-service.

```text
POST   /api/v1/media/upload/init
POST   /api/v1/media/:id/confirm
GET    /api/v1/media/:id/url
GET    /api/v1/media/:id/download-url
GET    /api/v1/media/:id
POST   /api/v1/media/:id/link
DELETE /api/v1/media/:id
```

### 7.2 Media lifecycle

```text
PENDING -> READY -> DELETED
```

- `PENDING` — metadata создана, upload URL выдан, объект еще не подтвержден.
- `READY` — объект найден в S3/MinIO, media можно использовать.
- `DELETED` — soft delete в БД; объект также удаляется из storage best-effort.

### 7.3 Upload algorithm

1. Клиент вызывает `POST /api/v1/media/upload/init` с `kind`, `content_type`, `size_bytes`, `original_name`, privacy flag.
2. Media Service валидирует kind/content type.
3. Создается запись `media` со статусом `PENDING`.
4. Генерируется object key и presigned PUT URL.
5. Клиент загружает файл напрямую в MinIO/S3.
6. Клиент вызывает `POST /api/v1/media/:id/confirm`.
7. Service проверяет owner и наличие object в storage.
8. Статус переводится в `READY`.
9. При отправке сообщения messenger сохраняет `attachment_media_id`.
10. Messenger async вызывает gRPC `LinkMediaToEntity(media_id, "MESSAGE", message_id)`, чтобы media-service мог проверять ACL.

### 7.4 Media kinds

Код поддерживает доменные kinds вроде:

- `AVATAR`;
- `ATTACHMENT`;
- `VOICE_MESSAGE`.

Для avatar принимаются image content types. Для voice message ожидается audio. Attachment допускает широкий набор content types.

### 7.5 Private media ACL

Если media public, presigned view/download URL выдается после базовых проверок статуса.

Если `is_private=true`:

1. Media Service достает links из `media_entity_links`.
2. Для link `MESSAGE` вызывает messenger gRPC:

```text
MessengerAccessService.CheckEntityAccess(user_id, "MESSAGE", message_id)
```

3. Messenger находит message -> chat -> membership/read access.
4. Если доступ подтвержден, media-service выдает presigned URL.
5. Если links отсутствуют или доступ не подтвержден, возвращается forbidden.

## 8. Transaction Service и Web3

Transaction Service отвечает за реализованную Web3-часть проекта: custodial wallets, SUDA balances, transfers, donations, simulated purchases, channel treasury read model и token-gating.

### 8.1 HTTP routes

Через gateway доступны с prefix `/api/v1/tx`.

Wallet:

```text
GET  /wallet/me
POST /wallet/resolve
POST /wallet/transfer
GET  /wallet/history
GET  /wallet/channel/:channel_id
GET  /wallet/channel/:channel_id/treasury
GET  /wallet/channel/:channel_id/donations
POST /wallet/channel/:channel_id/withdraw
```

Purchase:

```text
GET  /purchase/packages
POST /purchase/initiate
POST /purchase/:id/confirm
GET  /purchase/history
```

Donation:

```text
POST /donate
```

Gating:

```text
POST   /gating/rule
GET    /gating/rule/:chat_id
DELETE /gating/rule/:chat_id
```

Зарегистрированные, но пустые на текущей итерации groups:

```text
/nft
/marketplace
/fundraise
/quests
```

### 8.2 Custodial wallets

Для пользователя:

- таблица `tx_wallets`;
- ключевая привязка: `user_id -> EVM address`;
- приватный ключ хранится как encrypted payload;
- encryption: AES-256-GCM;
- `key_version` позволяет версионировать ключ шифрования;
- `suda_balance_cache` — кеш, не источник истины.

Для канала:

- таблица `tx_channel_wallets`;
- ключевая привязка: `channel_id -> EVM address`;
- используется как treasury address канала;
- private key также зашифрован.

Создание кошелька идемпотентно. Если wallet уже существует, service возвращает существующий address и `existed=true`.

### 8.3 gRPC TransactionService

Реализуется transaction-service, вызывается messenger-service:

| RPC | Назначение |
|---|---|
| `CreateWalletForUser` | Создает custodial wallet при регистрации пользователя. |
| `CreateWalletForChannel` | Создает treasury wallet при создании канала. |
| `GetWallet` | Возвращает user wallet address из БД без Besu RPC. |
| `GetChannelWallet` | Возвращает channel treasury address. |
| `GetBalance` | Делает `SudaToken.balanceOf` на Besu и обновляет кеш. |
| `CheckTokenGating` | Проверяет gating rule; возвращает `required`, `passed`, `price_wei`. |
| `ChargeChannelSubscription` | Списывает `subscription_price_wei` с кошелька пользователя в treasury канала при оформлении платной подписки. |

### 8.4 Username resolving и обратные callbacks

Transaction Service не владеет username/display name. Для операций по username он вызывает messenger-service через `TransactionAccessService`.

| RPC | Кто вызывает | Для чего |
|---|---|---|
| `ResolveUsername` | transaction -> messenger | `@username` -> `user_id`, `display_name`, `wallet_address`. |
| `NotifyUserEvent` | transaction -> messenger | WS notification и optional system message в чате. |
| `CheckChannelPermission` | transaction -> messenger | Проверка owner/admin/member прав для channel treasury/gating. |
| `CheckChatMembership` | transaction -> messenger | Проверка, что отправитель и получатель состоят в чате для transfer-in-chat/donation-in-chat. |
| `GetUsersByIDs` | transaction -> messenger | Обогащение treasury/donation статистики username/display name. |

### 8.5 SUDA transfer

`POST /api/v1/tx/wallet/transfer` выполняет P2P перевод SUDA.

Алгоритм:

1. Парсится `amount_wei`, сумма должна быть положительной integer string.
2. `to_username` резолвится через messenger gRPC `ResolveUsername`.
3. Проверяется, что получатель найден и имеет wallet address.
4. Загружается wallet отправителя из `tx_wallets`.
5. Self-transfer запрещен.
6. Если передан `chat_id`, transaction-service вызывает `CheckChatMembership`, чтобы и отправитель, и получатель были участниками чата.
7. Приватный ключ отправителя расшифровывается.
8. Через Besu читается on-chain `SudaToken.balanceOf`.
9. При недостаточном балансе операция отклоняется до broadcast.
10. В `tx_signing_audit` пишется audit entry до отправки транзакции.
11. Broadcaster готовит nonce, gas settings и `TransactOpts`.
12. Подписывается и отправляется `SudaToken.Transfer`.
13. В `tx_pending` пишется pending row с tx hash и optional chat id.
14. API сразу возвращает `tx_hash`; подтверждение приходит асинхронно через observer.

### 8.6 Donation

`POST /api/v1/tx/donate` выполняет донат SUDA.

Поддерживаются два получателя:

- `to_username` — P2P донат пользователю;
- `to_channel_id` — донат на treasury кошелек канала.

Ровно один получатель должен быть задан.

Алгоритм похож на transfer:

1. Парсинг amount.
2. Выбор recipient mode: user или channel.
3. Для user recipient вызывается `ResolveUsername`; для channel recipient берется `tx_channel_wallets`.
4. Загружается wallet отправителя.
5. Self-donation запрещен.
6. Для P2P donation в чате проверяется membership обоих пользователей.
7. Делается on-chain balance check.
8. Пишется audit.
9. Отправляется `SudaToken.Transfer`.
10. В `tx_pending` пишется DONATION metadata: tx hash, sender, system message chat id, donation message.
11. Observer позже индексирует transfer, записывает donation read model и инициирует WS/system message через messenger.

### 8.7 Purchase

Purchase — имитация покупки SUDA за fiat. Реального платежного провайдера нет: card fields принимаются для UI, но не валидируются и не сохраняются.

Доступные packages зашиты в коде:

| Code | SUDA | Fiat |
|---|---:|---|
| `SMALL` | 100 | `1.99 USD` |
| `MEDIUM` | 500 | `8.99 USD` |
| `LARGE` | 1000 | `14.99 USD` |
| `MEGA` | 5000 | `49.99 USD` |

Flow:

1. `GET /purchase/packages` возвращает catalog.
2. `POST /purchase/initiate` создает `tx_suda_purchases` со статусом `PENDING`.
3. `POST /purchase/:id/confirm` ставит Redis lock на пользователя на 5 секунд.
4. Проверяется ownership и status.
5. Status CAS-переход: `PENDING -> PROCESSING`.
6. Имитируется processing delay 2-3 секунды.
7. Загружается wallet пользователя.
8. Treasury signer отправляет `SudaToken.Transfer` на wallet пользователя.
9. Purchase помечается `COMPLETED` с `tx_hash`.
10. В `tx_pending` пишется pending row.
11. Observer после индексации отправляет user event о завершении.

### 8.8 Token gating и платная подписка

CRUD правил:

- `POST /gating/rule`;
- `GET /gating/rule/:chat_id`;
- `DELETE /gating/rule/:chat_id`.

Создавать и удалять rule может только `OWNER` канала. Проверка прав делается через messenger gRPC `CheckChannelPermission`.

Rule хранится в `tx_gating_rules`:

- `chat_id`;
- `subscription_price_wei NUMERIC(78,0)` — цена подписки в wei; `0` = бесплатный;
- `min_suda_balance` — минимальный баланс для классического balance-check;
- optional `required_nft_collection_id`;
- `created_by`.

Логика при `subscribe`:

1. Messenger вызывает `CheckTokenGating(chat_id, user_id)` через gRPC.
2. Transaction читает rule; если rule отсутствует — `required=false`, вход свободный.
3. Если rule есть и `subscription_price_wei > 0`: возвращает `required=true`, `price_wei`.
4. Messenger проверяет, является ли пользователь уже участником (идемпотентность).
5. Если нет — вызывает `ChargeChannelSubscription(user_id, channel_id)`.
6. Transaction читает `subscription_price_wei`, кошелёк пользователя и treasury канала.
7. Выполняет `SudaToken.balanceOf` пользователя; при недостаточном балансе — отказ.
8. Подписывает и broadcast-ит `SudaToken.Transfer(treasury, price)` ключом пользователя.
9. Записывает `InsertPending("CHANNEL_SUBSCRIBE")`.
10. При успехе Messenger добавляет пользователя как подписчика (`SUBSCRIBER`).

`CheckTokenGating` также используется при чтении платного PUBLIC-канала не-участником — возвращает ошибку доступа до оформления подписки.

NFT-requirement предусмотрен в модели данных, полноценный NFT business-flow не реализован.

### 8.9 Channel treasury

Channel treasury endpoints:

- `GET /wallet/channel/:channel_id` — адрес treasury wallet канала;
- `GET /wallet/channel/:channel_id/treasury` — on-chain баланс казны, топ-доноры;
- `GET /wallet/channel/:channel_id/donations` — история донатов и подписок в казну;
- `POST /wallet/channel/:channel_id/withdraw` — вывод SUDA из казны на кошелёк владельца.

Перед выдачей treasury stats проверяются права канала через messenger-service. Статистика строится по channel wallet, donations read model и user enrichment через `GetUsersByIDs`.

**Вывод средств из казны (`WithdrawTreasury`):**

Доступен только `OWNER` канала. Алгоритм:

1. Парсинг `amount_wei` — должна быть положительная целая строка.
2. `CheckChannelPermission(callerID, channelID, OWNER)` через messenger gRPC — отказ если не owner.
3. Чтение `tx_channel_wallets` (treasury) и `tx_wallets` (owner как получатель).
4. Расшифровка приватного ключа казны (`encryptor.Decrypt`).
5. On-chain `SudaToken.balanceOf(treasuryAddr)` ≥ amount — иначе отказ.
6. `WriteAudit(SubjectChannel, OpChannelWithdraw)` до broadcast.
7. `SudaToken.Transfer(ownerAddr, amount)` подписывается ключом казны и broadcast-ится.
8. `InsertPending("CHANNEL_WITHDRAW")`.
9. Возвращает `{tx_hash, from_address, to_address, amount_wei}` немедленно; on-chain подтверждение — через observer.

Поддерживается частичный и полный вывод суммы.

### 8.10 Blockchain layer

Blockchain layer включает:

- `Client` — wrapper над `ethclient.Client`, проверяющий chain id;
- `Contracts` — Go bindings к Solidity contracts;
- `Reader` — чтение on-chain state вроде `SudaToken.balanceOf`;
- `Broadcaster` — подготовка и отправка signed transactions;
- `NonceManager` — per-address lock и nonce management;
- `TokenSender` — компактная обертка для отправки `SudaToken.Transfer`.

`Broadcaster.PrepareOpts`:

1. Лочит адрес signer-а в `NonceManager`.
2. Берет fresh nonce.
3. Создает `bind.TransactOpts`.
4. Ставит `GasPrice=0`, потому что Besu chain использует zero base fee.
5. Возвращает `release`, который обязательно вызывается после broadcast.

### 8.11 Observer

Observer — фоновый процесс внутри transaction-service. Он poll-ит Besu каждые 2 секунды.

State хранится в `tx_observer_state` отдельно по contract name:

- `SUDA_TOKEN`;
- `SUDA_NFT`;
- `SUDA_MARKETPLACE`;
- `SUDA_ESCROW`;
- `SUDA_FUNDRAISING`.

На текущей итерации полноценная обработка реализована для `SudaToken`. Остальные cursors двигаются вперед как placeholders без бизнес-обработки logs.

Алгоритм tick:

1. Получить current block.
2. Прочитать last processed block для `SUDA_TOKEN`.
3. Если есть новые блоки, обработать диапазон максимум 1000 блоков за раз.
4. Handler читает Transfer logs.
5. Обновляет read model `tx_transactions`.
6. Сопоставляет tx hash с `tx_pending`.
7. Для purchase/donation обновляет соответствующие tables.
8. Удаляет/закрывает pending.
9. Через messenger gRPC `NotifyUserEvent` отправляет WS event и optional system message.
10. Только после успешной обработки двигает cursor.

## 9. Wallet WebView client

Wallet client находится в `web/wallet`.

Технологии:

- React;
- Vite;
- TypeScript;
- Zustand;
- Axios;
- hash router.

Production path:

```text
/wallet/
```

Routes:

| Route | Назначение |
|---|---|
| `/` | Главный экран кошелька: address, balance, recent history, welcome bonus state. |
| `/send` | Отправка SUDA по username. |
| `/buy` | Покупка SUDA package через simulated checkout. |
| `/receive` | Получение адреса кошелька. |
| `/history` | История операций с infinite scroll/filtering. |
| `*` | Error page. |

API client использует base URL:

```text
/api/v1
```

Основные calls:

- `POST /api/v1/messenger/auth/webapp-session` — получить session/access token для WebView;
- `GET /api/v1/tx/wallet/me` — address + balance;
- `GET /api/v1/tx/wallet/history` — transaction history;
- `POST /api/v1/tx/wallet/resolve` — username resolving;
- `POST /api/v1/tx/wallet/transfer` — transfer;
- `GET /api/v1/tx/purchase/packages` — packages;
- `POST /api/v1/tx/purchase/initiate` — create pending purchase;
- `POST /api/v1/tx/purchase/:id/confirm` — simulated payment confirmation.

WS client подключается к:

```text
ws(s)://<current-host>/ws?token=<access_token>
```

Incoming WS events обновляют wallet store и history. На главном экране есть polling/soft provisioning behavior: если `/tx/wallet/me` временно возвращает 404, UI считает wallet still provisioning и повторяет запрос несколько раз.

## 10. Proto contracts

Proto-контракты лежат в `backend/proto`. Генерация:

```bash
task go-proto
```

### 10.1 `media/media.proto`

Service: `MediaService`.

Используется другими backend-сервисами для работы с media metadata и links.

RPC:

- `GetMedia(GetMediaRequest)`;
- `LinkMediaToEntity(LinkMediaToEntityRequest)`;
- `BatchGetMediaURLs(BatchGetMediaURLsRequest)`.

Ключевые модели:

- `MediaInfo`;
- `LinkMediaToEntityRequest`;
- `BatchGetMediaURLsResponse`, где `urls` — map `media_id -> presigned URL`.

### 10.2 `messenger_access/messenger_access.proto`

Service: `MessengerAccessService`.

Реализуется messenger-service, вызывается media-service.

RPC:

- `CheckEntityAccess(user_id, entity_type, entity_id)`.

Типичный flow: media-service получает запрос на private media URL, находит link `MESSAGE -> message_id`, спрашивает messenger-service, имеет ли user доступ к этому message/chat.

### 10.3 `transaction/transaction.proto`

Service: `TransactionService`.

Реализуется transaction-service, вызывается messenger-service.

RPC:

- `CreateWalletForUser`;
- `CreateWalletForChannel`;
- `GetWallet`;
- `GetChannelWallet`;
- `GetBalance`;
- `CheckTokenGating`.

Этот контракт нужен messenger-service при регистрации, создании каналов, отображении wallet data и проверке gated каналов.

### 10.4 `transaction_access/service.proto`

Service: `TransactionAccessService`.

Реализуется messenger-service, вызывается transaction-service.

RPC:

- `ResolveUsername`;
- `NotifyUserEvent`;
- `CheckChannelPermission`;
- `CheckChatMembership`;
- `GetUsersByIDs`.

Этот контракт закрывает обратное направление: transaction-service не хранит messenger identity/social graph, поэтому спрашивает messenger-service о username, membership и правах.

## 11. Database tables

### 11.1 Messenger и Media

| Таблица | Краткое назначение |
|---|---|
| `messenger_users` | Пользователи мессенджера: email/username/display profile, avatar, wallet address mirror, search fields. |
| `messenger_verifications` | Email/password verification codes и purpose-specific verification state. |
| `messenger_sessions` | Refresh sessions, session ids, expiration/revocation state. |
| `messenger_user_devices` | Device tokens/metadata для push notifications. |
| `messenger_chats` | Общая таблица DIRECT/GROUP/CHANNEL/SAVED чатов и channel metadata. |
| `messenger_chat_members` | Membership, roles, read cursor, mute/notification settings, `cleared_at` (время очистки истории у участника), `hidden_at` (время скрытия чата из списка). |
| `messenger_messages` | Сообщения, attachments, replies, forwards, edit/delete timestamps, search vector. |
| `messenger_message_hidden` | Per-user скрытие сообщений при delete-for-me. |
| `messenger_blocked_users` | Отношения block между пользователями. |
| `messenger_contacts` | Пользовательские contact names для других пользователей. |
| `messenger_pinned_messages` | Закрепленные сообщения в чатах. |
| `messenger_reactions` | Реакции на сообщения; модель есть, WS reaction handler пока не реализует бизнес-логику. |
| `messenger_sticker_packs` | Метаданные sticker packs. |
| `messenger_stickers` | Стикеры внутри sticker packs. |
| `messenger_user_preferences` | Настройки приватности и поведения пользователя. |
| `messenger_last_seen_hidden` | Список пользователей, от которых скрыт last seen. |
| `messenger_channel_post_comments` | Комментарии к постам каналов. |
| `messenger_channel_apps` | Связи каналов с external/app сущностями. |
| `messenger_user_pins` | Пользовательские pins с типом цели и сортировкой. |
| `messenger_friend_requests` | Friend request lifecycle между пользователями. |
| `messenger_channel_join_requests` | Private channel join requests и invites. |
| `media` | Метаданные файлов, S3 bucket/object key, owner, status, privacy. |
| `media_entity_links` | Связи media с MESSAGE/CHAT_AVATAR/USER_AVATAR и другими сущностями для ACL. |

### 11.2 Transaction

| Таблица | Краткое назначение |
|---|---|
| `tx_wallets` | Custodial wallets пользователей с encrypted private key и balance cache. |
| `tx_channel_wallets` | Custodial treasury wallets каналов. |
| `tx_transactions` | Индексированная read model on-chain операций. |
| `tx_pending` | Pending tx rows для операций, которые уже broadcast, но еще ждут observer confirmation. |
| `tx_signing_audit` | Audit journal до подписи/broadcast sensitive операций. |
| `tx_observer_state` | Last processed block cursor по контрактам observer-а. |
| `tx_nft_collections` | Модель NFT collections для будущих/частичных NFT сценариев. |
| `tx_nft_items` | Модель NFT items для будущих/частичных NFT сценариев. |
| `tx_marketplace_listings` | Модель marketplace listings; готовый HTTP marketplace flow сейчас не реализован. |
| `tx_donations` | Read model донатов после observer processing. |
| `tx_fundraisers` | Модель fundraising entities; готовый HTTP fundraising flow сейчас не реализован. |
| `tx_fundraiser_contributions` | Модель contributions к fundraising; готовый flow сейчас не реализован. |
| `tx_quests` | Модель quests; готовый HTTP quests flow сейчас не реализован. |
| `tx_gating_rules` | Token-gating правила каналов: `subscription_price_wei` (цена платной подписки в wei), `min_suda_balance`, optional NFT collection. |
| `tx_suda_purchases` | Simulated SUDA purchase records и их lifecycle. |

## 12. Реализованные алгоритмы и состояния для диаграмм

### 12.1 Auth/register -> wallet creation

```text
Client -> Gateway -> Messenger Auth
Messenger creates user
Messenger -> TransactionService.CreateWalletForUser
Transaction generates EVM keypair
Transaction encrypts private key
Transaction writes tx_wallets
Transaction returns address
Messenger stores/mirrors wallet address where needed
Client receives auth/session response
```

### 12.2 Message with media

```text
Client -> Media upload/init
Media -> media(PENDING) + presigned PUT
Client -> MinIO PUT object
Client -> Media confirm
Media -> media(READY)
Client -> Messenger SendMessage(attachment_media_id)
Messenger -> messenger_messages
Messenger -> MediaService.LinkMediaToEntity(MESSAGE, message_id)
Messenger -> WS NEW_MESSAGE
Recipient -> Media get URL
Media -> MessengerAccess.CheckEntityAccess
Media -> presigned GET URL
```

### 12.3 WS send/read/typing

```text
Client opens /ws?token=...
Hub validates JWT
Hub stores session and online key in Redis
Client SEND_MESSAGE / TYPING / READ_MESSAGE
WS Router dispatches to ChatService
ChatService updates DB/cache
Hub sends event to relevant user sessions
```

### 12.4 SUDA transfer

```text
Wallet WebView -> /api/v1/tx/wallet/transfer
Gateway validates JWT
Transaction resolves recipient username via Messenger gRPC
Transaction checks sender wallet and optional chat membership
Transaction decrypts private key
Transaction reads SudaToken.balanceOf
Transaction writes tx_signing_audit
Transaction signs and broadcasts SudaToken.Transfer
Transaction writes tx_pending
Observer reads SudaToken Transfer log
Observer writes tx_transactions
Observer clears pending
Observer calls Messenger NotifyUserEvent
Messenger sends WS event and optional system message
```

### 12.5 Purchase lifecycle

```text
PENDING -> PROCESSING -> COMPLETED
                  |
                  v
                FAILED
```

`PENDING` создается на initiate, `PROCESSING` ставится на confirm, `COMPLETED` ставится после treasury transfer broadcast, `FAILED` ставится при no wallet/broadcast/processing failure.

### 12.6 Media lifecycle

```text
PENDING -> READY -> DELETED
```

### 12.7 Pending transaction lifecycle

```text
API broadcasts tx
API writes tx_pending
Observer finds matching on-chain log
Observer creates tx_transactions/read model
Observer applies operation-specific side effects
Observer removes or resolves pending state
Messenger receives NotifyUserEvent
Client updates over WS
```

### 12.8 Paid channel subscription

```text
Client -> POST /channels/:id/subscribe
Messenger -> CheckTokenGating(channelID, userID)   [gRPC -> Transaction]
Transaction reads tx_gating_rules -> subscription_price_wei > 0
Transaction reads SudaToken.balanceOf(userWallet)
Transaction returns required=true, price_wei=X
[if already member -> return OK, skip charge]
Messenger -> ChargeChannelSubscription(userID, channelID)  [gRPC -> Transaction]
Transaction loads userWallet, channelTreasuryWallet
Transaction decrypts user private key
Transaction checks balance >= price_wei
Transaction broadcasts SudaToken.Transfer(treasury, price)
Transaction writes tx_pending("CHANNEL_SUBSCRIBE")
Messenger -> AddMember(channelID, userID, SUBSCRIBER)
Client receives 200 + CHANNEL_SUBSCRIBED WS event
```

### 12.9 At-rest message encryption

```text
Client -> POST /chats/:id/messages  {content: "hello"}
Messenger service sets ChatType on message struct
ContentCipher.EncryptContent("hello")
  -> AES-256-GCM(key, nonce) -> base64 -> "enc:v1:<ciphertext>"
DB: messenger_messages.content = "enc:v1:<ciphertext>"
DB: messenger_messages.search_vector = to_tsvector("hello")  [from plaintext, before encrypt]

Client -> GET /chats/:id/messages
Repo reads content = "enc:v1:<ciphertext>"
ContentCipher.DecryptContent("enc:v1:<ciphertext>")
  -> strip sentinel -> AES-256-GCM decrypt -> "hello"
Response: content = "hello"
```

## 13. Testing

### 13.1 Go tests

В репозитории есть unit tests по ключевым сервисным слоям:

| Путь | Что покрывает |
|---|---|
| `backend/messenger-service/internal/features/chat/service/chat_service_test.go` | Создание чатов, group behavior, permissions, pin/member/message-related logic. |
| `backend/messenger-service/internal/features/search/service/search_service_test.go` | Search service behavior. |
| `backend/messenger-service/internal/features/user/service/user_service_test.go` | User service behavior. |
| `backend/messenger-service/internal/platform/crypto/content_cipher_test.go` | AES-256-GCM roundtrip, legacy plaintext passthrough, bad sentinel, noop cipher, random nonce uniqueness. |
| `backend/media-service/internal/service/media_service_test.go` | Upload init, validation, confirm, private ACL, URLs, delete. |
| `backend/transaction-service/internal/features/donation/service/donation_service_test.go` | Donation validation, membership, recipient/channel behavior. |
| `backend/transaction-service/internal/features/gating/service/gating_service_test.go` | Gating rule CRUD и permission checks. |
| `backend/transaction-service/internal/features/purchase/service/purchase_service_test.go` | Purchase initiate/confirm/history behavior. |
| `backend/transaction-service/internal/features/wallet/service/treasury_test.go` | Channel treasury stats/donation listing behavior. |
| `backend/transaction-service/internal/observer/handlers/suda_token_test.go` | SudaToken observer handling. |

Команды:

```bash
cd backend/messenger-service && go test ./...
cd backend/media-service && go test ./...
cd backend/transaction-service && go test ./...
```

Transaction Service race-enabled shortcut:

```bash
task go-test-tx
```

### 13.2 Bruno collection

Bruno collection находится в:

```text
bruno/Suda Messenger
```

Она используется для ручного и сценарного тестирования REST/WS API. Основные папки:

- `Auth` — register/login/verify/refresh/password flows;
- `Gateway` — gateway health и auth negative cases;
- `User` — profile, avatar, sessions, status, device, init-data;
- `Chats` — direct/group/saved chats, messages, replies, forward, media, pins, members, block/contact flows;
- `Channels` — subscribe, comments, settings, apps;
- `Media` — upload init, upload to S3, confirm, metadata, URLs, delete, invalid cases;
- `Preferences` — privacy settings and hidden last seen;
- `Search` — global/chat search and negative cases;
- `UserPins` — create/delete/list/reorder pins;
- `Ws` — WS connect, send, typing, read, forward, invalid JSON/no token/unknown event;
- `Full Conversation Test` — end-to-end сценарии полного общения;
- `Transaction` — wallet, purchase, donations, gating, treasury и transaction healthcheck.

### 13.3 Swagger

Swagger docs лежат в `backend/docs/swagger`.

Generation:

```bash
task go-swagger-all
task go-swagger-tx
task go-swagger-messenger
task go-swagger-media
```

Production URLs через Nginx:

```text
/swagger/messenger/index.html
/swagger/transaction/index.html
```

Gateway aggregator:

```text
/swagger/doc.json
/swagger/
```

## 14. Suggested diagrams and UML

Для технического отчета/презентации хорошо подходят следующие диаграммы.

### 14.1 C4 Container diagram

Показать:

- Mobile/WebView client;
- Nginx;
- API Gateway;
- Messenger Service;
- Media Service;
- Transaction Service;
- PostgreSQL;
- Redis;
- MinIO;
- Besu node;
- Wallet frontend;
- Portainer as ops-only component.

### 14.2 Component diagram for backend services

Показать HTTP и gRPC зависимости:

- gateway -> services via REST proxy;
- messenger -> media via `MediaService`;
- messenger -> transaction via `TransactionService`;
- media -> messenger via `MessengerAccessService`;
- transaction -> messenger via `TransactionAccessService`;
- transaction -> Besu via JSON-RPC.

### 14.3 Sequence: register/login -> wallet creation

Показать создание пользователя, wallet provisioning через gRPC и сохранение EVM address.

### 14.4 Sequence: send message with attachment

Показать media upload init, direct PUT to MinIO, confirm, send message, async link, WS delivery и private media URL ACL.

### 14.5 Sequence: WebSocket realtime

Показать connect, JWT validation, Redis online state, `SEND_MESSAGE`, `TYPING`, `READ_MESSAGE`, broadcast to sessions.

### 14.6 Sequence: SUDA transfer

Показать WebView -> gateway -> transaction-service -> messenger resolving -> Besu transfer -> observer -> messenger notification -> WS.

### 14.7 Sequence: channel donation

Показать donation to channel treasury, pending donation metadata, observer, `tx_donations`, system message в channel chat.

### 14.8 Sequence: token-gated subscribe

Показать user subscribe request, messenger -> transaction `CheckTokenGating`, balanceOf на Besu, allow/deny result.

### 14.9 ERD

Разделить на три bounded contexts:

- Messenger: users/chats/members/messages/preferences/friends/channels;
- Media: media/media_entity_links;
- Transaction: wallets/transactions/pending/audit/observer/purchases/donations/gating.

### 14.10 State diagrams

Подходящие state machines:

- `media.status`: `PENDING -> READY -> DELETED`;
- `tx_suda_purchases.status`: `PENDING -> PROCESSING -> COMPLETED/FAILED`;
- pending transaction: `broadcasted -> pending -> observed -> notified`;
- private channel request: `PENDING -> APPROVED/REJECTED`.

## 15. Operational commands

Частые команды:

```bash
# Production stack
task docker-up-prod

# Dev stack
task docker-up

# Status
task docker-ps

# Logs
task docker-logs
docker compose logs -f messenger-service
docker compose logs -f transaction-service

# Rebuild one service
task docker-rebuild SVC=transaction-service

# Stop stack, keep volumes
task docker-down

# Destructive stop with volumes
task docker-down-volumes

# Postgres shell
docker compose exec postgres psql -U postgres suda

# Transaction tests
task go-test-tx
```
