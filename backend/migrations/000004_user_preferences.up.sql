ALTER TABLE messenger_messages
    ADD COLUMN search_vector tsvector;

UPDATE messenger_messages
SET search_vector =
    setweight(to_tsvector('russian', COALESCE(content, '')), 'A') ||
    setweight(to_tsvector('simple',  COALESCE(content, '')), 'B')
WHERE type != 'SYSTEM' AND deleted_at IS NULL;

CREATE INDEX idx_messages_search_vector ON messenger_messages USING GIN(search_vector);

CREATE OR REPLACE FUNCTION messages_search_trigger() RETURNS trigger AS $$
BEGIN
    IF NEW.type != 'SYSTEM' THEN
        NEW.search_vector :=
            setweight(to_tsvector('russian', COALESCE(NEW.content, '')), 'A') ||
            setweight(to_tsvector('simple',  COALESCE(NEW.content, '')), 'B');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;



CREATE TRIGGER trg_messages_search
    BEFORE INSERT OR UPDATE OF content
    ON messenger_messages
    FOR EACH ROW
    EXECUTE FUNCTION messages_search_trigger();

ALTER TABLE messenger_users
    ADD COLUMN search_vector tsvector;

UPDATE messenger_users
SET search_vector =
    setweight(to_tsvector('simple', COALESCE(username, '')), 'A') ||
    setweight(to_tsvector('simple', COALESCE(display_name, '')), 'B') ||
    setweight(to_tsvector('simple', COALESCE(first_name, '')), 'C') ||
    setweight(to_tsvector('simple', COALESCE(last_name, '')), 'C');

CREATE INDEX idx_users_search_vector ON messenger_users USING GIN(search_vector);

CREATE OR REPLACE FUNCTION users_search_trigger() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('simple', COALESCE(NEW.username, '')), 'A') ||
        setweight(to_tsvector('simple', COALESCE(NEW.display_name, '')), 'B') ||
        setweight(to_tsvector('simple', COALESCE(NEW.first_name, '')), 'C') ||
        setweight(to_tsvector('simple', COALESCE(NEW.last_name, '')), 'C');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_search
    BEFORE INSERT OR UPDATE OF username, display_name, first_name, last_name
    ON messenger_users
    FOR EACH ROW
    EXECUTE FUNCTION users_search_trigger();

ALTER TABLE messenger_chats
    ADD COLUMN search_vector tsvector;

UPDATE messenger_chats
SET search_vector =
    setweight(to_tsvector('simple', COALESCE(name, '')), 'A') ||
    setweight(to_tsvector('simple', COALESCE(description, '')), 'B')
WHERE type IN ('GROUP', 'CHANNEL');

CREATE INDEX idx_chats_search_vector ON messenger_chats USING GIN(search_vector);

CREATE OR REPLACE FUNCTION chats_search_trigger() RETURNS trigger AS $$
BEGIN
    IF NEW.type IN ('GROUP', 'CHANNEL') THEN
        NEW.search_vector :=
            setweight(to_tsvector('simple', COALESCE(NEW.name, '')), 'A') ||
            setweight(to_tsvector('simple', COALESCE(NEW.description, '')), 'B');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_chats_search
    BEFORE INSERT OR UPDATE OF name, description
    ON messenger_chats
    FOR EACH ROW
    EXECUTE FUNCTION chats_search_trigger();

CREATE TABLE messenger_user_preferences (
    user_id UUID PRIMARY KEY REFERENCES messenger_users(id) ON DELETE CASCADE,

    lang_code VARCHAR(5) DEFAULT 'en',
    theme VARCHAR(10) DEFAULT 'system',

    last_seen_visibility VARCHAR(20) DEFAULT 'EVERYONE',
    avatar_visibility VARCHAR(20) DEFAULT 'EVERYONE',
    who_can_message VARCHAR(20) DEFAULT 'EVERYONE',

    notification_messages BOOLEAN DEFAULT TRUE,
    notification_groups BOOLEAN DEFAULT TRUE,
    notification_channels BOOLEAN DEFAULT TRUE,
    notification_sound BOOLEAN DEFAULT TRUE,
    notification_vibration BOOLEAN DEFAULT TRUE,

    chat_font_size VARCHAR(10) DEFAULT 'medium',
    send_by_enter BOOLEAN DEFAULT TRUE, -- useless in current mobile version, but why not
    auto_download_media BOOLEAN DEFAULT TRUE,

    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE messenger_last_seen_hidden (
    owner_id  UUID NOT NULL REFERENCES messenger_users(id) ON DELETE CASCADE,
    hidden_id UUID NOT NULL REFERENCES messenger_users(id) ON DELETE CASCADE,
    PRIMARY KEY (owner_id, hidden_id)
);

CREATE INDEX idx_last_seen_hidden_owner ON messenger_last_seen_hidden(owner_id);
