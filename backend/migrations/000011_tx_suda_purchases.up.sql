-- tx_suda_purchases — записи об "имитированных" покупках SUDA за фиатные деньги.
-- Реальной финансовой транзакции нет: status проходит PENDING -> PROCESSING -> COMPLETED
-- через имитацию (sleep 2-3 сек), потом treasury переводит SUDA пользователю.
CREATE TABLE tx_suda_purchases (
    id              UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID         NOT NULL,
    package_code    VARCHAR(50)  NOT NULL,
    suda_amount_wei NUMERIC(78, 0) NOT NULL,
    fiat_amount     NUMERIC(10, 2) NOT NULL,
    fiat_currency   VARCHAR(3)   NOT NULL,
    status          VARCHAR(20)  NOT NULL,    -- PENDING | PROCESSING | COMPLETED | FAILED
    payment_method  VARCHAR(20),              -- CARD | APPLE_PAY | GOOGLE_PAY
    tx_hash         VARCHAR(66),              -- on-chain transfer treasury -> user
    failure_reason  TEXT,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    completed_at    TIMESTAMPTZ
);

CREATE INDEX idx_tx_purchases_user
    ON tx_suda_purchases(user_id, created_at DESC);

-- Observer ищет purchase по tx_hash, чтобы помечать tx_transactions.type=PURCHASE
-- и слать PURCHASE_COMPLETED event.
CREATE INDEX idx_tx_purchases_tx_hash
    ON tx_suda_purchases(tx_hash)
    WHERE tx_hash IS NOT NULL;

-- Помогает находить "застрявшие" PENDING/PROCESSING для cleanup-задач в будущем.
CREATE INDEX idx_tx_purchases_active
    ON tx_suda_purchases(status, created_at)
    WHERE status IN ('PENDING', 'PROCESSING');
