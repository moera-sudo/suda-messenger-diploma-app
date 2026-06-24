CREATE TABLE media (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_user_id   UUID NOT NULL,
    bucket          TEXT NOT NULL,
    object_key      TEXT NOT NULL UNIQUE,
    original_name   TEXT,
    kind            TEXT NOT NULL,
    content_type    TEXT NOT NULL,
    size_bytes      BIGINT,
    checksum_sha256 TEXT,
    status          TEXT NOT NULL DEFAULT 'PENDING',
    is_private      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ready_at        TIMESTAMPTZ,
    deleted_at      TIMESTAMPTZ
);

-- Привязка медиа к сущностям, управляющим доступом
-- entity_type = 'MESSAGE'      entity_id = message_id  -> доступ через chat membership
-- entity_type = 'CHAT_AVATAR'  entity_id = chat_id     -> доступ через chat membership
-- entity_type = 'USER_AVATAR'  entity_id = user_id     -> публичный (is_private=false)
CREATE TABLE media_entity_links (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    media_id    UUID NOT NULL REFERENCES media(id) ON DELETE CASCADE,
    entity_type TEXT NOT NULL,
    entity_id   TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (media_id, entity_type, entity_id)
);

CREATE INDEX idx_media_owner         ON media(owner_user_id);
CREATE INDEX idx_media_status        ON media(status) WHERE status != 'DELETED';
CREATE INDEX idx_media_bucket_key    ON media(bucket, object_key);
CREATE INDEX idx_media_entity_media  ON media_entity_links(media_id);
CREATE INDEX idx_media_entity_lookup ON media_entity_links(entity_type, entity_id);
