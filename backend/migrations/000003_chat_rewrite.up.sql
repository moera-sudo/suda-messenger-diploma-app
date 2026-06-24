ALTER TABLE messenger_chats
    ADD COLUMN description TEXT,
    ADD COLUMN avatar_media_id UUID,
    ADD COLUMN creator_id UUID REFERENCES messenger_users(id) ON DELETE SET NULL,
    ADD COLUMN invite_link TEXT UNIQUE;

ALTER TABLE messenger_messages
    ADD COLUMN edited_at TIMESTAMPTZ,
    ADD COLUMN deleted_at TIMESTAMPTZ;

ALTER TABLE messenger_chat_members
    ADD COLUMN is_muted BOOLEAN DEFAULT FALSE,
    ADD COLUMN muted_until TIMESTAMPTZ,
    ADD COLUMN notifications_enabled BOOLEAN DEFAULT TRUE;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'messenger_users' AND column_name = 'bio') THEN
        ALTER TABLE messenger_users ADD COLUMN bio TEXT;
    END IF;
END $$;


CREATE TABLE messenger_message_hidden (
    user_id    UUID NOT NULL REFERENCES messenger_users(id) ON DELETE CASCADE,
    message_id BIGINT NOT NULL REFERENCES messenger_messages(id) ON DELETE CASCADE,
    hidden_at  TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, message_id)
);

CREATE TABLE messenger_blocked_users (
    blocker_id UUID NOT NULL REFERENCES messenger_users(id) ON DELETE CASCADE,
    blocked_id UUID NOT NULL REFERENCES messenger_users(id) ON DELETE CASCADE,
    blocked_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (blocker_id, blocked_id)
);

CREATE TABLE messenger_contacts (
    owner_id     UUID NOT NULL REFERENCES messenger_users(id) ON DELETE CASCADE,
    contact_id   UUID NOT NULL REFERENCES messenger_users(id) ON DELETE CASCADE,
    custom_name  TEXT NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (owner_id, contact_id)
);

CREATE TABLE messenger_pinned_messages (
    chat_id    UUID NOT NULL REFERENCES messenger_chats(id) ON DELETE CASCADE,
    message_id BIGINT NOT NULL REFERENCES messenger_messages(id) ON DELETE CASCADE,
    pinned_by  UUID REFERENCES messenger_users(id) ON DELETE SET NULL,
    pinned_at  TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (chat_id, message_id)
);

CREATE TABLE messenger_reactions (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    message_id BIGINT NOT NULL REFERENCES messenger_messages(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES messenger_users(id) ON DELETE CASCADE,
    emoji      TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (message_id, user_id, emoji)
);

CREATE TABLE messenger_sticker_packs (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        TEXT NOT NULL,
    creator_id  UUID REFERENCES messenger_users(id) ON DELETE SET NULL,
    is_official BOOLEAN DEFAULT FALSE,
    is_nft      BOOLEAN DEFAULT FALSE,
    cover_media_id UUID,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE messenger_stickers (
    id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pack_id  UUID NOT NULL REFERENCES messenger_sticker_packs(id) ON DELETE CASCADE,
    media_id UUID NOT NULL,
    emoji    TEXT,
    nft_token_id TEXT,
    position INT DEFAULT 0
);

CREATE INDEX idx_message_hidden_user ON messenger_message_hidden(user_id);
CREATE INDEX idx_blocked_users_blocker ON messenger_blocked_users(blocker_id);
CREATE INDEX idx_blocked_users_blocked ON messenger_blocked_users(blocked_id);
CREATE INDEX idx_contacts_owner on messenger_contacts(owner_id);
CREATE INDEX idx_pinned_messages_chat ON messenger_pinned_messages(chat_id);
CREATE INDEX idx_reactions_message ON messenger_reactions(message_id);
CREATE INDEX idx_messages_deleted ON messenger_messages(deleted_at) WHERE deleted_at IS NOT NULL;
CREATE INDEX idx_messages_edited ON messenger_messages(edited_at) WHERE edited_at IS NOT NULL;
CREATE INDEX idx_chats_creator ON messenger_chats(creator_id);
CREATE INDEX idx_chats_type ON messenger_chats(type);
