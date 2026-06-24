-- ════════════════════════════════════════════════════════════
--  Transaction Service — initial schema (single migration)
--
--  АРХИТЕКТУРНЫЕ ПРАВИЛА (Custodial @wallet-style, web3 source-of-truth):
--
--  1. БЛОКЧЕЙН — ИСТОЧНИК ПРАВДЫ для всех финансовых данных.
--     PostgreSQL — read model, обновляется ТОЛЬКО observer'ом из event logs.
--
--  2. tx_transactions — append-only. Заполняется ТОЛЬКО observer'ом после
--     подтверждения блока. HTTP-handler'ы НЕ ПИШУТ сюда напрямую.
--     Pending-индикатор для UI живёт в tx_pending и удаляется, когда
--     observer создаёт соответствующую запись в tx_transactions.
--
--  3. tx_wallets.suda_balance_cache — кеш с коротким TTL (~30 сек).
--     Реальный баланс — SudaToken.balanceOf(address). Кеш используется
--     только для списков (топ-донатеры и т.п.), не для отображения юзеру.
--
--  4. Любая операция дешифровки приватного ключа пишет строку
--     в tx_signing_audit. subject_type определяет ROOT владельца ключа:
--     USER (юзер) или CHANNEL (treasury канала).
--
--  5. Observer хранит cursor в tx_observer_state — позицию (last_processed_block)
--     per contract, чтобы пережить рестарт сервиса без потери/дубля событий.
-- ════════════════════════════════════════════════════════════

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


-- ──────────────────────────────────────────────────────────
--  USER WALLETS
-- ──────────────────────────────────────────────────────────

-- user_id            : ссылка на messenger_users.id (FK не ставим — другой сервис)
-- encrypted_private_key : AES-256-GCM(nonce || ciphertext || tag), base64
-- key_version        : версия мастер-ключа (под будущую ротацию)
-- suda_balance_cache : ТОЛЬКО кеш, не источник правды
CREATE TABLE tx_wallets (
    user_id               UUID PRIMARY KEY,
    address               VARCHAR(42) NOT NULL UNIQUE,
    encrypted_private_key TEXT NOT NULL,
    key_version           SMALLINT NOT NULL DEFAULT 1,
    suda_balance_cache    NUMERIC(78, 0) DEFAULT 0,
    balance_updated_at    TIMESTAMPTZ,
    created_at            TIMESTAMPTZ DEFAULT NOW()
);


-- ──────────────────────────────────────────────────────────
--  CHANNEL WALLETS (treasury каждого канала)
--  Отдельный кошелёк per channel — для приёма донатов, paid subscription'ов,
--  выплат админам, и т.п. Owner канала (через CheckChannelPermission)
--  может инициировать withdraw.
-- ──────────────────────────────────────────────────────────
CREATE TABLE tx_channel_wallets (
    channel_id            UUID PRIMARY KEY,
    address               VARCHAR(42) NOT NULL UNIQUE,
    encrypted_private_key TEXT NOT NULL,
    key_version           SMALLINT NOT NULL DEFAULT 1,
    suda_balance_cache    NUMERIC(78, 0) DEFAULT 0,
    balance_updated_at    TIMESTAMPTZ,
    created_at            TIMESTAMPTZ DEFAULT NOW()
);


-- ──────────────────────────────────────────────────────────
--  TRANSACTIONS (append-only read model — пишется ТОЛЬКО observer'ом)
--
--  UNIQUE(tx_hash, log_index) — один tx может породить несколько events:
--    например, SudaMarketplace.Sold внутри себя дёргает SudaToken.Transfer
--    и SudaNFT.Transfer. Каждое — отдельная запись с разным log_index.
--    Для P2P-перевода будет одна запись (log_index=0).
-- ──────────────────────────────────────────────────────────

-- type:
--   P2P_TRANSFER | DONATION | NFT_GIFT | NFT_MINT | NFT_BUY | NFT_LIST | NFT_CANCEL_LIST
--   QUEST_LOCK | QUEST_REWARD | QUEST_REFUND
--   FUNDRAISE_DONATE | FUNDRAISE_WITHDRAW | FUNDRAISE_REFUND
--   MINT | BURN
-- status:
--   CONFIRMED | FAILED   (PENDING сюда не попадает — он живёт в tx_pending)
-- related_entity_type:
--   NFT_ITEM | QUEST | FUNDRAISE | CHANNEL | CHAT_MESSAGE
CREATE TABLE tx_transactions (
    id                   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tx_hash              VARCHAR(66) NOT NULL,
    log_index            INTEGER NOT NULL,
    from_user_id         UUID,
    to_user_id           UUID,
    from_address         VARCHAR(42) NOT NULL,
    to_address           VARCHAR(42) NOT NULL,
    amount               NUMERIC(78, 0) NOT NULL,
    type                 VARCHAR(30) NOT NULL,
    status               VARCHAR(20) NOT NULL,
    related_entity_type  VARCHAR(30),
    related_entity_id    UUID,
    note                 TEXT,
    block_number         BIGINT NOT NULL,
    confirmed_at         TIMESTAMPTZ NOT NULL,
    created_at           TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE (tx_hash, log_index)
);

-- Append-only гарантия: блокируем UPDATE иммутабельных полей и любые DELETE.
-- Observer пишет одну запись и больше её не трогает.
CREATE OR REPLACE FUNCTION tx_transactions_immutable() RETURNS trigger AS $$
BEGIN
    IF TG_OP = 'UPDATE' THEN
        IF OLD.tx_hash      IS DISTINCT FROM NEW.tx_hash      OR
           OLD.log_index    IS DISTINCT FROM NEW.log_index    OR
           OLD.from_address IS DISTINCT FROM NEW.from_address OR
           OLD.to_address   IS DISTINCT FROM NEW.to_address   OR
           OLD.amount       IS DISTINCT FROM NEW.amount       OR
           OLD.type         IS DISTINCT FROM NEW.type         OR
           OLD.block_number IS DISTINCT FROM NEW.block_number OR
           OLD.confirmed_at IS DISTINCT FROM NEW.confirmed_at OR
           OLD.created_at   IS DISTINCT FROM NEW.created_at THEN
            RAISE EXCEPTION 'tx_transactions: immutable field modified';
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_tx_transactions_immutable
    BEFORE UPDATE ON tx_transactions
    FOR EACH ROW
    EXECUTE FUNCTION tx_transactions_immutable();

-- Запретить DELETE полностью.
CREATE RULE tx_transactions_no_delete AS
    ON DELETE TO tx_transactions DO INSTEAD NOTHING;


-- ──────────────────────────────────────────────────────────
--  PENDING TRANSACTIONS (UI индикатор «отправляется...»)
--
--  Создаётся в момент broadcast'а транзакции.
--  Удаляется observer'ом после того, как он записал tx_transactions.
--  Висящие записи > 5 минут считаем зависшими — фоновая чистка (TODO).
-- ──────────────────────────────────────────────────────────
CREATE TABLE tx_pending (
    tx_hash       VARCHAR(66) PRIMARY KEY,
    from_user_id  UUID NOT NULL,
    expected_type VARCHAR(30) NOT NULL,    -- что мы ожидаем (для UI hint)
    submitted_at  TIMESTAMPTZ DEFAULT NOW()
);


-- ──────────────────────────────────────────────────────────
--  SIGNING AUDIT (лог использования приватных ключей)
--
--  subject_type: USER | CHANNEL  — кому принадлежит ключ
--  subject_id  : user_id или channel_id соответственно
--  operation   : что делали (TRANSFER, DONATION, NFT_GIFT, и т.п.)
--  tx_hash     : NULL если запись сделана ДО broadcast'а (полный journal)
--
--  Каждая дешифровка приватного ключа => одна строка.
-- ──────────────────────────────────────────────────────────
CREATE TABLE tx_signing_audit (
    id           BIGSERIAL PRIMARY KEY,
    subject_type VARCHAR(20) NOT NULL,
    subject_id   UUID NOT NULL,
    operation    VARCHAR(50) NOT NULL,
    tx_hash      VARCHAR(66),
    request_ip   TEXT,
    user_agent   TEXT,
    occurred_at  TIMESTAMPTZ DEFAULT NOW()
);


-- ──────────────────────────────────────────────────────────
--  OBSERVER STATE (cursor для polling)
--
--  Per-contract last_processed_block. Observer на старте читает отсюда
--  и продолжает с last+1. После обработки окна — UPSERT.
--  Если записи нет — observer стартует с current block чейна (не индексируем
--  историю до запуска).
-- ──────────────────────────────────────────────────────────
CREATE TABLE tx_observer_state (
    contract_name        VARCHAR(50) PRIMARY KEY,
    last_processed_block BIGINT NOT NULL,
    updated_at           TIMESTAMPTZ DEFAULT NOW()
);


-- ──────────────────────────────────────────────────────────
--  NFT (collections + items)
--
--  On-chain source of truth — SudaNFT (ERC-721).
--  owner_user_id обновляет observer на каждый Transfer event.
-- ──────────────────────────────────────────────────────────
CREATE TABLE tx_nft_collections (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name           VARCHAR(100) NOT NULL,
    description    TEXT,
    creator_id     UUID,
    cover_media_id UUID,
    is_official    BOOLEAN DEFAULT FALSE,
    base_price     NUMERIC(78, 0) DEFAULT 0,
    created_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE tx_nft_items (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    token_id        NUMERIC(78, 0) NOT NULL UNIQUE,
    collection_id   UUID REFERENCES tx_nft_collections(id) ON DELETE SET NULL,
    owner_user_id   UUID NOT NULL,
    owner_address   VARCHAR(42) NOT NULL,
    name            VARCHAR(100) NOT NULL,
    description     TEXT,
    media_id        UUID NOT NULL,
    category        VARCHAR(30) DEFAULT 'STICKER',  -- STICKER | GIFT | AVATAR_FRAME
    metadata_uri    TEXT,
    minted_at       TIMESTAMPTZ DEFAULT NOW(),
    transferred_at  TIMESTAMPTZ
);


-- ──────────────────────────────────────────────────────────
--  MARKETPLACE LISTINGS
--  On-chain source of truth — SudaMarketplace.
-- ──────────────────────────────────────────────────────────
CREATE TABLE tx_marketplace_listings (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nft_item_id      UUID NOT NULL REFERENCES tx_nft_items(id) ON DELETE CASCADE,
    seller_user_id   UUID NOT NULL,
    seller_address   VARCHAR(42) NOT NULL,
    price            NUMERIC(78, 0) NOT NULL,
    status           VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',  -- ACTIVE | SOLD | CANCELLED
    listed_at        TIMESTAMPTZ DEFAULT NOW(),
    sold_at          TIMESTAMPTZ,
    buyer_user_id    UUID,
    buyer_address    VARCHAR(42),
    tx_hash          VARCHAR(66)
);


-- ──────────────────────────────────────────────────────────
--  DONATIONS (плоская история донатов)
--  to_address: либо channel treasury, либо адрес юзера (P2P донат)
-- ──────────────────────────────────────────────────────────
CREATE TABLE tx_donations (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    channel_id    UUID,                                       -- NULL для P2P донатов
    to_user_id    UUID,                                       -- NULL для donation в канал
    from_user_id  UUID NOT NULL,
    from_address  VARCHAR(42) NOT NULL,
    to_address    VARCHAR(42) NOT NULL,
    amount        NUMERIC(78, 0) NOT NULL,
    message       TEXT,
    tx_hash       VARCHAR(66) NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT NOW()
);


-- ──────────────────────────────────────────────────────────
--  FUNDRAISERS (сборы средств)
--  On-chain source of truth — SudaFundraising контракт.
-- ──────────────────────────────────────────────────────────
CREATE TABLE tx_fundraisers (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    channel_id        UUID NOT NULL,
    creator_id        UUID NOT NULL,
    creator_address   VARCHAR(42) NOT NULL,
    contract_id       NUMERIC(78, 0) NOT NULL UNIQUE,
    title             VARCHAR(200) NOT NULL,
    description       TEXT,
    goal              NUMERIC(78, 0) NOT NULL,
    raised            NUMERIC(78, 0) NOT NULL DEFAULT 0,
    deadline          TIMESTAMPTZ NOT NULL,
    status            VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    -- ACTIVE | GOAL_REACHED | WITHDRAWN | EXPIRED
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    completed_at      TIMESTAMPTZ
);

CREATE TABLE tx_fundraiser_contributions (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    fundraiser_id   UUID NOT NULL REFERENCES tx_fundraisers(id) ON DELETE CASCADE,
    from_user_id    UUID NOT NULL,
    from_address    VARCHAR(42) NOT NULL,
    amount          NUMERIC(78, 0) NOT NULL,
    tx_hash         VARCHAR(66) NOT NULL,
    refunded        BOOLEAN DEFAULT FALSE,
    refund_tx_hash  VARCHAR(66),
    created_at      TIMESTAMPTZ DEFAULT NOW()
);


-- ──────────────────────────────────────────────────────────
--  QUESTS / BOUNTIES (эскроу)
--  On-chain source of truth — SudaEscrow.
-- ──────────────────────────────────────────────────────────
CREATE TABLE tx_quests (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    channel_id       UUID NOT NULL,
    creator_id       UUID NOT NULL,
    creator_address  VARCHAR(42) NOT NULL,
    contract_id      NUMERIC(78, 0) NOT NULL UNIQUE,

    title            VARCHAR(200) NOT NULL,
    description      TEXT,
    reward           NUMERIC(78, 0) NOT NULL,
    deadline         TIMESTAMPTZ NOT NULL,

    status           VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    -- OPEN | CLAIMED | SUBMITTED | APPROVED | CANCELLED | EXPIRED

    assignee_id      UUID,
    assignee_address VARCHAR(42),
    claimed_at       TIMESTAMPTZ,
    submitted_at     TIMESTAMPTZ,
    submission_note  TEXT,
    completed_at     TIMESTAMPTZ,

    created_at       TIMESTAMPTZ DEFAULT NOW()
);


-- ──────────────────────────────────────────────────────────
--  GATING (правила доступа к чату/каналу)
--  Одно правило на чат. Условия по AND: min_suda И required_nft.
-- ──────────────────────────────────────────────────────────
CREATE TABLE tx_gating_rules (
    chat_id                     UUID PRIMARY KEY,
    min_suda_balance            NUMERIC(78, 0) NOT NULL DEFAULT 0,
    required_nft_collection_id  UUID REFERENCES tx_nft_collections(id) ON DELETE SET NULL,
    created_by                  UUID NOT NULL,
    created_at                  TIMESTAMPTZ DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ DEFAULT NOW()
);


-- ════════════════════════════════════════════════════════════
--  INDEXES
-- ════════════════════════════════════════════════════════════

-- Wallets
CREATE INDEX idx_tx_wallets_address          ON tx_wallets(address);
CREATE INDEX idx_tx_channel_wallets_address  ON tx_channel_wallets(address);

-- Transactions
CREATE INDEX idx_tx_txns_tx_hash    ON tx_transactions(tx_hash);
CREATE INDEX idx_tx_txns_from_user  ON tx_transactions(from_user_id) WHERE from_user_id IS NOT NULL;
CREATE INDEX idx_tx_txns_to_user    ON tx_transactions(to_user_id)   WHERE to_user_id   IS NOT NULL;
CREATE INDEX idx_tx_txns_type       ON tx_transactions(type);
CREATE INDEX idx_tx_txns_confirmed  ON tx_transactions(confirmed_at DESC);
CREATE INDEX idx_tx_txns_block      ON tx_transactions(block_number DESC);
CREATE INDEX idx_tx_txns_related    ON tx_transactions(related_entity_type, related_entity_id)
    WHERE related_entity_id IS NOT NULL;

-- Pending
CREATE INDEX idx_tx_pending_user      ON tx_pending(from_user_id);
CREATE INDEX idx_tx_pending_submitted ON tx_pending(submitted_at);  -- для cleanup >5 мин

-- Signing audit
CREATE INDEX idx_tx_audit_subject   ON tx_signing_audit(subject_type, subject_id, occurred_at DESC);
CREATE INDEX idx_tx_audit_operation ON tx_signing_audit(operation);
CREATE INDEX idx_tx_audit_tx_hash   ON tx_signing_audit(tx_hash) WHERE tx_hash IS NOT NULL;

-- NFT
CREATE INDEX idx_tx_collections_creator  ON tx_nft_collections(creator_id) WHERE creator_id IS NOT NULL;
CREATE INDEX idx_tx_collections_official ON tx_nft_collections(is_official) WHERE is_official = TRUE;
CREATE INDEX idx_tx_nft_owner            ON tx_nft_items(owner_user_id);
CREATE INDEX idx_tx_nft_collection       ON tx_nft_items(collection_id) WHERE collection_id IS NOT NULL;
CREATE INDEX idx_tx_nft_category         ON tx_nft_items(category);

-- Marketplace
CREATE INDEX idx_tx_listings_seller       ON tx_marketplace_listings(seller_user_id);
CREATE INDEX idx_tx_listings_nft          ON tx_marketplace_listings(nft_item_id);
CREATE INDEX idx_tx_listings_active       ON tx_marketplace_listings(status)         WHERE status = 'ACTIVE';
CREATE INDEX idx_tx_listings_active_price ON tx_marketplace_listings(price)          WHERE status = 'ACTIVE';
CREATE INDEX idx_tx_listings_listed_at    ON tx_marketplace_listings(listed_at DESC) WHERE status = 'ACTIVE';

-- Donations
CREATE INDEX idx_tx_donations_channel ON tx_donations(channel_id)  WHERE channel_id IS NOT NULL;
CREATE INDEX idx_tx_donations_to_user ON tx_donations(to_user_id)  WHERE to_user_id IS NOT NULL;
CREATE INDEX idx_tx_donations_from    ON tx_donations(from_user_id);
CREATE INDEX idx_tx_donations_created ON tx_donations(created_at DESC);

-- Fundraisers
CREATE INDEX idx_tx_fundraisers_channel ON tx_fundraisers(channel_id);
CREATE INDEX idx_tx_fundraisers_creator ON tx_fundraisers(creator_id);
CREATE INDEX idx_tx_fundraisers_status  ON tx_fundraisers(status);
CREATE INDEX idx_tx_fundraisers_active  ON tx_fundraisers(deadline) WHERE status = 'ACTIVE';

CREATE INDEX idx_tx_contribs_fundraiser     ON tx_fundraiser_contributions(fundraiser_id);
CREATE INDEX idx_tx_contribs_from           ON tx_fundraiser_contributions(from_user_id);
CREATE INDEX idx_tx_contribs_pending_refund ON tx_fundraiser_contributions(fundraiser_id)
    WHERE refunded = FALSE;

-- Quests
CREATE INDEX idx_tx_quests_channel  ON tx_quests(channel_id);
CREATE INDEX idx_tx_quests_creator  ON tx_quests(creator_id);
CREATE INDEX idx_tx_quests_assignee ON tx_quests(assignee_id) WHERE assignee_id IS NOT NULL;
CREATE INDEX idx_tx_quests_status   ON tx_quests(status);
CREATE INDEX idx_tx_quests_open     ON tx_quests(deadline)
    WHERE status IN ('OPEN', 'CLAIMED', 'SUBMITTED');

-- Gating
CREATE INDEX idx_tx_gating_created_by ON tx_gating_rules(created_by);
CREATE INDEX idx_tx_gating_nft        ON tx_gating_rules(required_nft_collection_id)
    WHERE required_nft_collection_id IS NOT NULL;
