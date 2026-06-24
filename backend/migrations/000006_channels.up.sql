-- 6_channels.sql

-- ══════════════════════════════════════════════════════════
--  Расширения messenger_chats для каналов
-- ══════════════════════════════════════════════════════════

ALTER TABLE messenger_chats
    ADD COLUMN visibility VARCHAR(10) DEFAULT 'PRIVATE',
    ADD COLUMN subscriber_count INT DEFAULT 0,
    ADD COLUMN comments_enabled BOOLEAN DEFAULT TRUE,
    ADD COLUMN username VARCHAR(32) UNIQUE;

CREATE INDEX idx_chats_username ON messenger_chats(username) WHERE username IS NOT NULL;
CREATE INDEX idx_chats_visibility ON messenger_chats(visibility) WHERE type = 'CHANNEL';

-- ══════════════════════════════════════════════════════════
--  Комментарии под постами канала
-- ══════════════════════════════════════════════════════════

CREATE TABLE messenger_channel_post_comments (
    id BIGSERIAL PRIMARY KEY,
    channel_id           UUID NOT NULL REFERENCES messenger_chats(id) ON DELETE CASCADE,
    post_id              BIGINT NOT NULL REFERENCES messenger_messages(id) ON DELETE CASCADE,
    sender_id            UUID NOT NULL REFERENCES messenger_users(id) ON DELETE CASCADE,
    content              TEXT NOT NULL,
    reply_to_comment_id  BIGINT REFERENCES messenger_channel_post_comments(id) ON DELETE SET NULL,
    edited_at            TIMESTAMPTZ,
    deleted_at           TIMESTAMPTZ,
    created_at           TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_channel_comments_post ON messenger_channel_post_comments(post_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_channel_comments_channel ON messenger_channel_post_comments(channel_id);
CREATE INDEX idx_channel_comments_sender ON messenger_channel_post_comments(sender_id);

-- ══════════════════════════════════════════════════════════
--  Привязки мини-приложений к каналам (заготовка)
-- ══════════════════════════════════════════════════════════

CREATE TABLE messenger_channel_apps (
    channel_id UUID NOT NULL REFERENCES messenger_chats(id) ON DELETE CASCADE,
    app_id     UUID NOT NULL,
    added_by   UUID NOT NULL REFERENCES messenger_users(id),
    added_at   TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (channel_id, app_id)
);

CREATE INDEX idx_channel_apps_channel ON messenger_channel_apps(channel_id);

-- ══════════════════════════════════════════════════════════
--  Триггер: автообновление subscriber_count
-- ══════════════════════════════════════════════════════════

CREATE OR REPLACE FUNCTION update_subscriber_count() RETURNS trigger AS $$
DECLARE
    chat_type VARCHAR(20);
BEGIN
    IF TG_OP = 'INSERT' THEN
        SELECT type INTO chat_type FROM messenger_chats WHERE id = NEW.chat_id;
        IF chat_type = 'CHANNEL' THEN
            UPDATE messenger_chats
            SET subscriber_count = subscriber_count + 1
            WHERE id = NEW.chat_id;
        END IF;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        SELECT type INTO chat_type FROM messenger_chats WHERE id = OLD.chat_id;
        IF chat_type = 'CHANNEL' THEN
            UPDATE messenger_chats
            SET subscriber_count = GREATEST(subscriber_count - 1, 0)
            WHERE id = OLD.chat_id;
        END IF;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_subscriber_count
    AFTER INSERT OR DELETE ON messenger_chat_members
    FOR EACH ROW
    EXECUTE FUNCTION update_subscriber_count();

-- ══════════════════════════════════════════════════════════
--  Поиск: только публичные каналы попадают в выдачу
-- ══════════════════════════════════════════════════════════

DROP TRIGGER IF EXISTS trg_chats_search ON messenger_chats;

CREATE OR REPLACE FUNCTION chats_search_trigger() RETURNS trigger AS $$
BEGIN
    IF NEW.type = 'GROUP' OR (NEW.type = 'CHANNEL' AND NEW.visibility = 'PUBLIC') THEN
        NEW.search_vector :=
            setweight(to_tsvector('simple', COALESCE(NEW.name, '')), 'A') ||
            setweight(to_tsvector('simple', COALESCE(NEW.username, '')), 'A') ||
            setweight(to_tsvector('simple', COALESCE(NEW.description, '')), 'B');
    ELSE
        NEW.search_vector := NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_chats_search
    BEFORE INSERT OR UPDATE OF name, description, username, visibility
    ON messenger_chats
    FOR EACH ROW
    EXECUTE FUNCTION chats_search_trigger();
