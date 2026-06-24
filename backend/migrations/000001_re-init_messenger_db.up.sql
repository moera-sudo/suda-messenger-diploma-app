CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE messenger_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) NOT NULL UNIQUE,
    display_name VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    first_name VARCHAR(60),
    last_name VARCHAR(60),
    bio TEXT,
    role VARCHAR(20) NOT NULL DEFAULT 'USER', -- USER/ADMIN/SUPERADMIN
    avatar_media_id UUID,
    last_seen_at TIMESTAMPTZ DEFAULT NOW(),
    is_verified BOOLEAN DEFAULT FALSE, -- Подтвердил ли почту
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE messenger_verifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES messenger_users(id) ON DELETE CASCADE, 
    email TEXT NOT NULL,
    code_hash TEXT NOT NULL,          
    purpose VARCHAR(50) NOT NULL,    -- 'REGISTRATION', 'PASSWORD_RESET', 'CHANGE_EMAIL'
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,             -- Если NULL - еще не использован
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE messenger_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES messenger_users(id) ON DELETE CASCADE, 
    refresh_token TEXT NOT NULL,
    user_agent TEXT,
    client_ip TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    is_revoked BOOLEAN DEFAULT FALSE,
    revoked_at TIMESTAMPTZ
);

CREATE TABLE messenger_user_devices (
    user_id UUID NOT NULL REFERENCES messenger_users(id) ON DELETE CASCADE,
    fcm_token TEXT NOT NULL, -- Токен, который выдает Firebase на клиенте
    device_name TEXT,        -- "iPhone 13" (для красоты в настройках)
    last_used_at TIMESTAMPTZ DEFAULT NOW(),
    
    PRIMARY KEY (user_id, fcm_token) -- Один токен не может дублироваться у юзера
);

CREATE TABLE messenger_chats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), 
    type VARCHAR(20) NOT NULL, -- 'DIRECT', 'GROUP'
    name TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE messenger_chat_members (
    chat_id UUID NOT NULL REFERENCES messenger_chats(id) ON DELETE CASCADE, 
    user_id UUID NOT NULL REFERENCES messenger_users(id) ON DELETE CASCADE,
    role VARCHAR(20) DEFAULT 'MEMBER',
    joined_at TIMESTAMPTZ DEFAULT NOW(),
    last_read_message_id BIGINT DEFAULT 0,
    PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE messenger_messages (
    id BIGSERIAL PRIMARY KEY, -- возможно надо поставить UUID
    chat_id UUID NOT NULL REFERENCES messenger_chats(id) ON DELETE CASCADE, 
    sender_id UUID REFERENCES messenger_users(id) ON DELETE SET NULL,
    attachment_media_id UUID,
    reply_to_message_id BIGINT REFERENCES messenger_messages(id) ON DELETE SET NULL,
    forwarded_from_chat UUID,
    forwarded_from_msg BIGINT,       
    content TEXT,
    type VARCHAR(20) DEFAULT 'TEXT',
    status VARCHAR(20) DEFAULT 'SENT',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_verifications_user_purpose ON messenger_verifications(user_id, purpose);
CREATE INDEX idx_verifications_code ON messenger_verifications(code_hash); -- Если будем искать по хэшу
CREATE INDEX idx_messages_chat_id ON messenger_messages(chat_id);
CREATE INDEX idx_chat_members_user_id ON messenger_chat_members(user_id);
CREATE INDEX idx_sessions_refresh_token ON messenger_sessions(refresh_token);
CREATE INDEX idx_messages_attachment ON messenger_messages(attachment_media_id)
    WHERE attachment_media_id IS NOT NULL;

CREATE INDEX idx_messages_reply ON messenger_messages(reply_to_message_id)
    WHERE reply_to_message_id IS NOT NULL;
    