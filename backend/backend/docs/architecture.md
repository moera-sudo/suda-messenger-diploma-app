# Suda Messenger — Architecture

> Бэкенд мессенджера с интегрированным blockchain-кошельком (custodial). Документ — обзорный, описывает сервисы, протоколы взаимодействия, инфраструктуру, blockchain-flow и БД. Конкретный API см. в Swagger UI каждого сервиса (`/swagger/`) и в [api-overview.md](api-overview.md).

---

## 1. Общая картина

Suda — Go-микросервисный backend для мобильного мессенджера (Flutter) с встроенным non-custodial-ощущением, но фактически **custodial**-кошельком на приватной EVM-сети (Hyperledger Besu, QBFT). Каждый юзер получает Ethereum-адрес и зашифрованный приватный ключ в БД; от его имени бэк подписывает и broadcast'ит транзакции (`SudaToken.transfer`, NFT-минт, donations). Подтверждение приходит асинхронно через **observer** — фоновый процесс, который опрашивает блокчейн и пишет confirmed-транзакции в append-only таблицу.

Поверх блокчейна построены фичи: P2P-переводы SUDA, донаты юзерам и каналам, имитация покупки SUDA, token-gating подписки на канал (минимальный баланс / NFT-владение), treasury-статистика каналов.

---

## 2. Сервисы и порты

| Сервис | HTTP | gRPC | Зачем |
|--------|-----:|-----:|-------|
| **api-gateway** | 8080 | — | Единственная внешняя точка входа. Проверяет JWT, подписывает запросы `X-Gateway-Signature` и проксирует их к downstream-сервисам. |
| **messenger-service** | 8081 | 9081 | Пользователи, аутентификация (JWT + email-verify), чаты, сообщения, каналы (комментарии, апы), поиск, preferences, user-pins, WebSocket-hub для real-time. |
| **transaction-service** | 8082 | 9082 | Кошельки (custodial), переводы SUDA, donations, token-gating, purchases (имитация), treasury-stats. Содержит **observer** — фоновый poll-loop блокчейна. |
| **media-service** | 8084 | 9094 | Presigned URL'ы к MinIO (S3) для upload/download. Проверяет права доступа через gRPC к messenger. |

Все сервисы билдятся multi-stage Dockerfile: `golang:alpine` → `alpine:3.20`. Конфиг загружается через [godotenv]: пробует `.env` в cwd, затем `../../.env`, затем системные env'ы.

### Лейаут модуля

Каждый сервис — отдельный Go-модуль, объединяется в `go.work` на корне репозитория. Внутри клин-архитектура:

```
internal/
  app/            # DI-композиция, lifecycle (Init / Run / Shutdown)
  features/<f>/
    delivery/http # echo handlers + DTO + swagger annotations
    service/      # бизнес-логика, доменные ошибки
    repository/   # raw SQL (pgx)
    models.go     # доменные структуры
  platform/       # инфра-адаптеры: postgres, redis, grpc, blockchain, crypto
  pkg/            # утилиты: http envelope, middlewares, logger, utils
  observer/       # (только tx) фоновый watcher событий блокчейна
```

DI — конструкторный, никаких контейнеров. Все зависимости передаются явно в `New(deps)`. Это даёт тестируемость через small per-feature interfaces (см. `internal/mocks/` в transaction-service).

---

## 3. Инфраструктура

| Компонент | Версия / Образ | Назначение |
|-----------|----------------|------------|
| **PostgreSQL** | 16 | Единая БД для всех сервисов (схемно разделено префиксами `messenger_`, `media`, `tx_`). pgx-драйвер, raw SQL, без ORM. |
| **Redis** | alpine | Кеш + Pub/Sub. Используется messenger'ом для WS-fanout и transaction-service для rate-limit (purchase confirm). |
| **MinIO** | latest | S3-совместимое объектное хранилище для медиа (avatars, attachments, stickers). |
| **Besu (Hyperledger)** | 25.12.0 | Приватная EVM-сеть, консенсус **QBFT**, blockPeriod=2s. Один узел в dev, multi-node QBFT планируется на prod. |
| **Firebase Cloud Messaging** | — | Push-нотификации на мобильные клиенты. |
| **Gmail SMTP** | — | Email-верификация при регистрации + password-reset. |

---

## 4. Композиция: запрос «снаружи»

```
Mobile/Web client
       │  HTTPS + Bearer JWT
       ▼
┌──────────────────┐
│  api-gateway     │  валидирует JWT, добавляет X-Gateway-Signature
│  (port 8080)     │  + проксирует HTTP
└────────┬─────────┘
         │
   ┌─────┴───────┬──────────────┐
   ▼             ▼              ▼
┌────────┐  ┌────────────┐  ┌────────┐
│messenger│  │transaction │  │ media  │
│  :8081 │  │   :8082    │  │ :8084  │
└────────┘  └────────────┘  └────────┘
   │  ▲         │   ▲          │
   │  │ gRPC    │   │ gRPC     │ gRPC
   │  └─────────┴───┴──────────┘
   │     │
   │     ▼
   │  PostgreSQL / Redis / MinIO
   │
   ▼
WebSocket (melody)
```

**HTTP middleware stack** (одинаковый во всех сервисах):
1. RequestID
2. Recover (panic → 500)
3. CORS
4. **GatewaySignature** — проверяет `X-Gateway-Signature` (HMAC от gateway). Без него — 401. (В dev можно обойти через `DEV_SKIP_GATEWAY=true`.)
5. Validator (echo + go-playground/validator)

---

## 5. gRPC-связи между сервисами

```
                ResolveUsername(username) → {user_id, wallet_address}
                CheckChannelPermission(user, channel, perm) → granted?
                CheckChatMembership(chat, [users]) → membership[]
                NotifyUserEvent(user, event_type, payload, [target_chat_id]) → ok
                GetUsersByIDs([uuid]) → [{user_id, username, display_name}]
                ─────────────────────────────────────────────────────────────▶
   ┌──────────┐                                                    ┌───────────┐
   │messenger │                                                    │transaction│
   │  :9081   │ ◀────────────────────────────────────────────────  │   :9082   │
   └──────────┘                                                    └───────────┘
                CheckTokenGating(chat, user) → {passed, required,
                                                min_balance, user_balance, reason}
                CreateWalletForChannel(channel_id) → {address}
```

Связь **двусторонняя**, оба сервиса — клиенты и серверы друг друга:

- **messenger → transaction** (`TransactionAccessServer` в tx):
  - `CheckTokenGating` — при `Channel.Subscribe`. Если правило не выполнено → клиент получает 402 `INSUFFICIENT_BALANCE` / 403 `NO_WALLET` / `NFT_REQUIRED`.
  - `CreateWalletForChannel` — best-effort goroutine при создании канала (chat type=Channel).

- **transaction → messenger** (`TransactionAccessServer` в messenger):
  - `ResolveUsername` — при `Wallet.Transfer` / `Donate` для P2P (узнать `user_id` и `wallet_address`).
  - `CheckChannelPermission(OWNER | OWNER_OR_ADMIN)` — для `Gating.CreateRule/DeleteRule`, `Treasury.Get*`.
  - `CheckChatMembership` — при transfer-in-chat / P2P-donation в чате (проверить что оба участника — члены чата).
  - `NotifyUserEvent` — observer шлёт WS-события (`SUDA_RECEIVED`, `DONATION_SENT`, `PURCHASE_COMPLETED`, ...) и опционально создаёт system message в чате (`TargetChatID + SystemMessageType`).
  - `GetUsersByIDs` — batch-резолв username'ов для treasury top-donors / донат-листа.

- **media → messenger** — проверка доступа к чатам перед выдачей presigned URL.

Сгенерированный proto-код закоммичен в [backend/proto/](../proto/). Регенерация — `task go-proto`.

---

## 6. WebSocket

Real-time события доставляются через WebSocket-hub на messenger-service:

- Базируется на [melody](https://github.com/olahol/melody) — простой in-memory hub с broadcast.
- Соединение апгрейдится с авторизованного HTTP (тот же JWT).
- Сообщения формата `{event_type, payload, target_chat_id?, system_message_type?}`.
- `NotifyUserEvent` (gRPC) → hub шлёт юзеру; если задан `target_chat_id` + `system_message_type` — создаётся system message в чате и broadcastится всем участникам.

**Типы событий** (отправляются observer'ом transaction-service):

| event_type | Кому | Когда | Создаёт system message? |
|------------|------|-------|------|
| `SUDA_RECEIVED` | получатель | подтверждение P2P-перевода | да (если transfer-in-chat) |
| `SUDA_SENT` | отправитель | то же | нет (sysmsg уже создан событием SUDA_RECEIVED) |
| `DONATION_SENT` | донатер | подтверждение доната (юзеру/каналу) | да (DONATION) |
| `DONATION_RECEIVED` | получатель | подтверждение P2P-доната (не канал) | нет |
| `PURCHASE_COMPLETED` | покупатель | подтверждение покупки SUDA (treasury → user) | нет |

---

## 7. Blockchain integration

### 7.1 Custodial-модель

Каждый юзер при регистрации (после email-verify) получает кошелёк:
1. Генерируется новая EOA пара (secp256k1).
2. Приватный ключ шифруется **AES-256-GCM** с key из `WALLET_ENCRYPTION_KEY` (32-байт).
3. В БД (`tx_wallets`) пишутся: `user_id, address, encrypted_private_key, key_version, created_at`.
4. Treasury переводит **welcome bonus 100 SUDA** на новый адрес (best-effort goroutine).

Юзер **не знает** свой PK. Все операции инициирует бэк по JWT-токену. Эта схема compromise между UX (никаких seed-фраз) и security (на проде ключи должны быть в KMS/HSM, не в БД).

### 7.2 Flow перевода (transfer / donate / purchase)

Эталон — `walletsvc.Transfer`. Универсальные шаги:

```
1. Parse amount (decimal wei string → *big.Int)
2. Validate input (exactly one recipient, positive amount, ...)
3. Resolve recipient (username → user_id+address  ИЛИ  channel_id → channel wallet)
4. Lookup sender's wallet, проверить self-transfer
5. (Optional) Membership check (transfer-in-chat / P2P-donate в чате)
6. Decrypt sender's PK → blockchain.Signer
7. On-chain balance check (reader.SudaBalanceOf)
8. WriteAudit (tx_signing_audit, tx_hash="") — ДО broadcast'а (полный journal)
9. broadcaster.SendTokenTransfer(signer, to, amount) → tx_hash
10. InsertPending(tx_hash, expected_type, chat_id, [donation_message])
11. Return tx_hash (HTTP 202 / 201)

→ Async (observer):
12. observer.HandleLogs() ловит Transfer event
13. InsertTransaction(...) (ON CONFLICT DO NOTHING — идемпотентно)
14. DeletePending(tx_hash)
15. (DONATION) InsertDonation в tx_donations
16. NotifyUserEvent → клиент видит подтверждение, создаётся system message в чате
```

Универсальный механизм — флаг `expected_type` в `tx_pending`:
- `P2P_TRANSFER` → обычный transfer
- `DONATION` → дополнительно пишется `tx_donations` row + DONATION_SENT/RECEIVED events
- `PURCHASE` → детектится observer'ом по `tx_suda_purchases.tx_hash`, шлётся `PURCHASE_COMPLETED` без `SUDA_RECEIVED`

### 7.3 Observer pattern

[backend/transaction-service/internal/observer/](../transaction-service/internal/observer/) — отдельная goroutine, стартует в `app.Run()`.

```
loop {
   from_block, to_block := lookup_window(last_processed_block)
   events := contracts.Token.FilterTransfer(from..to)
   for ev in events {
       handler.handleTransfer(ev) {
           inserted, err := repo.InsertTransaction(...)
           if !inserted { DeletePending; continue }  // идемпотентность
           // ... insert donation / send WS events ...
       }
   }
   UpsertObserverState(last_processed_block = to_block)
   sleep(poll_interval)
}
```

**Идемпотентность** — критична. Если observer крашится между `InsertTransaction` и `UpsertObserverState`, тот же event придёт снова. `InsertTransaction` использует `ON CONFLICT (tx_hash, log_index) DO NOTHING`; возвращает `(inserted bool, err error)`. Если `inserted == false` — пропускаем дальнейшую обработку (донат не пишется второй раз, WS-event не дублируется).

### 7.4 Контракты

- **SudaToken** (ERC20) — основной токен, mintable owner'ом (treasury).
- **SudaNFT** — для NFT-фичи (пока не реализовано).
- **SudaMarketplace / SudaEscrow / SudaFundraising** — stub'ы (планируется в будущих итерациях).

Genesis pre-allocates `GLOBAL_TREASURY_ADDR` максимальный баланс (для welcome bonus и purchases). Адреса всех контрактов — в `.env` (`SUDA_TOKEN_ADDRESS`, и т.д.).

---

## 8. Обзор БД (PostgreSQL)

Одна БД `suda`, схемно разделённая префиксами:

### 8.1 messenger_*

| Таблица | Зачем |
|---------|-------|
| `messenger_users` | юзеры (username, display_name, email, password_hash, avatar_media_id) |
| `messenger_verifications` | email-verify токены |
| `messenger_sessions` | JWT refresh-сессии |
| `messenger_user_devices` | FCM-токены устройств |
| `messenger_chats` | чаты (DIRECT/GROUP/CHANNEL/SAVED) |
| `messenger_chat_members` | участники + роли (OWNER/ADMIN/MEMBER) |
| `messenger_messages` | сообщения, sequential id внутри чата |
| `messenger_user_preferences` | per-user настройки (язык, нотификации, приватность) |
| `messenger_last_seen_hidden` | юзеры, от которых скрыт last_seen |
| `messenger_channel_post_comments` | комментарии к постам каналов |
| `messenger_channel_apps` | прилинкованные апы к каналам |
| `messenger_user_pins` | закреплённые чаты в списке |
| `messenger_blocked_users` | блок-листы |
| `messenger_contacts` | контакты |
| `messenger_pinned_messages` | закреплённые сообщения в чате |
| `messenger_reactions` | реакции на сообщения |
| `messenger_sticker_packs` / `messenger_stickers` | стикеры |

### 8.2 tx_*

| Таблица | Зачем |
|---------|-------|
| `tx_wallets` | user_id → encrypted_private_key, address, key_version |
| `tx_channel_wallets` | channel_id → address (treasury канала) |
| `tx_transactions` | **append-only** confirmed транзакции (tx_hash + log_index PK), type, related_entity |
| `tx_pending` | UI-индикатор «idem в обработке», удаляется observer'ом после confirm. Поля: tx_hash, expected_type, related_chat_id, donation_message |
| `tx_signing_audit` | ПОЛНЫЙ журнал каждой подписи (даже если broadcast упадёт) |
| `tx_observer_state` | last_processed_block (для resume после рестарта) |
| `tx_donations` | плоские строки донатов (для treasury-stats и истории) |
| `tx_suda_purchases` | имитация покупки SUDA (PENDING → PROCESSING → COMPLETED/FAILED) |
| `tx_gating_rules` | per-channel правила (min_suda_balance, required_nft_collection_id) |
| `tx_nft_collections` / `tx_nft_items` | NFT (фича planned) |
| `tx_marketplace_listings` / `tx_fundraisers` / `tx_quests` | stub-таблицы для будущих фич |

### 8.3 media

| Таблица | Зачем |
|---------|-------|
| `media` | uploaded-файлы (s3_key, mime, size, owner_id) |
| `media_entity_links` | связь media ↔ chat/message/user/channel |

---

## 9. Аутентификация и безопасность

- **JWT** (`golang-jwt/v5`): access (30 min) + refresh (720h). Подписан `JWT_SECRET`.
- **Email-verify**: при register'е высылается код через Gmail SMTP. Без verify юзер не получает кошелёк.
- **Gateway signature** (`GATEWAY_SIGNATURE_SECRET`) — HMAC, который api-gateway добавляет в каждый запрос к downstream. downstream проверяет; без подписи — 401. Защита от прямого доступа к 8081/8082/8084.
- **Init-data sign token** (`INIT_DATA_SIGN_TOKEN`) — для подписи данных мобильного клиента при первой загрузке (anti-tampering).
- **Wallet encryption key** (`WALLET_ENCRYPTION_KEY`, 32 байта) — AES-256-GCM для приватных ключей. Версионируется через `key_version` колонку (rotation плановая на prod).
- **ETH private key** (`ETH_PRIVATE_KEY`) — приватник treasury-кошелька, нужен для welcome bonus + purchase confirm.

На production все 5 чувствительных значений (`JWT_SECRET`, `GATEWAY_SIGNATURE_SECRET`, `WALLET_ENCRYPTION_KEY`, `ETH_PRIVATE_KEY`, `INIT_DATA_SIGN_TOKEN`) должны жить в KMS/Vault/Docker-secrets, не в `.env`. Сейчас (Stage 4) — все в `.env`, что приемлемо только для dev.

---

## 10. Тестирование

| Сервис | Покрытие | Подход |
|--------|----------|--------|
| messenger | chat-service + auth + некоторые user-features (`internal/features/chat/service/chat_service_test.go` и др.) | testify + hand-written mocks в `internal/mocks/` + miniredis |
| transaction | 42 теста: purchase / donation / gating / treasury / observer (idempotency) | то же — hand-written mocks, fakeReader/fakeEncryptor/fakeMessenger inline в test-файлах |
| media | минимум | — |

Запуск:
- `task go-test-tx` — все unit-тесты транзакции с `-race`
- `cd backend/messenger-service && go test ./...` — мессенджер

API-тестирование — через Bruno-коллекции в `bruno/Suda Messenger/`. Покрывают auth, chats, channels, preferences, userpins, transaction (wallet/purchase/donation/gating/treasury).

---

## 11. История стейджей

- **Stage 1** — базовая messenger-service: users, auth, chats, direct messaging.
- **Stage 2** — channels, comments, apps; preferences; user-pins; media-service.
- **Stage 3** — transaction-service v1: custodial wallets, P2P transfer, observer, transfer-in-chat (system messages).
- **Stage 4** — purchase (имитация), donations (P2P + channel), token-gating, treasury statistics. Idempotency bugfix observer'а.
- **Stage 5** (текущий, не финализирован) — **документация** (Swagger всех эндпоинтов + этот файл + api-overview.md + examples). **Deployment** (docker-compose с nginx + portainer + secrets) — отложен на отдельную итерацию.

---

## 12. Ссылки

- [api-overview.md](api-overview.md) — табличная карта всех эндпоинтов
- [examples/.env.example](examples/.env.example) — шаблон переменных окружения
- [examples/qbft-config-example.json](examples/qbft-config-example.json) — пример QBFT-конфига Besu
- [swagger/messenger/](swagger/messenger/) / [swagger/transaction/](swagger/transaction/) — сгенерированные OpenAPI-спецификации
- [stage-plan.md](stage-plan.md) — план разработки по стейджам
