-- 7_user_pins.sql

CREATE TABLE messenger_user_pins (
    id           BIGSERIAL PRIMARY KEY,
    user_id      UUID NOT NULL REFERENCES messenger_users(id) ON DELETE CASCADE,
    pin_type     VARCHAR(20) NOT NULL,    -- SIDEBAR / CHATLIST / APP_HUB
    target_type  VARCHAR(20) NOT NULL,    -- CHAT / APP
    target_id    UUID NOT NULL,           -- chat_id или app_id
    sort_order   INT NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE (user_id, pin_type, target_id)
);

CREATE INDEX idx_user_pins_user_type ON messenger_user_pins(user_id, pin_type, sort_order);
CREATE INDEX idx_user_pins_target ON messenger_user_pins(target_type, target_id);

-- Валидация комбинаций (чтобы не закрепить APP в CHATLIST например)
ALTER TABLE messenger_user_pins
    ADD CONSTRAINT valid_pin_combinations CHECK (
        (pin_type = 'CHATLIST' AND target_type = 'CHAT') OR
        (pin_type = 'SIDEBAR'  AND target_type IN ('CHAT', 'APP')) OR
        (pin_type = 'APP_HUB'  AND target_type = 'APP')
    );
