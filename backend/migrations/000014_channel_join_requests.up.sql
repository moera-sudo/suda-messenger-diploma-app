-- 14_channel_join_requests.sql
--
-- Join requests (user -> private channel, admin approves) and invites
-- (admin -> user, user accepts). One row per (channel, user); `kind` tells the
-- direction, `status` tracks the lifecycle. On approve/accept the user becomes a
-- SUBSCRIBER in messenger_chat_members (handled in the service).

CREATE TABLE messenger_channel_join_requests (
    channel_id UUID NOT NULL REFERENCES messenger_chats(id)  ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES messenger_users(id)  ON DELETE CASCADE,
    kind       VARCHAR(10) NOT NULL,                       -- 'REQUEST' | 'INVITE'
    status     VARCHAR(10) NOT NULL DEFAULT 'PENDING',     -- 'PENDING' | 'APPROVED' | 'REJECTED'
    created_by UUID NOT NULL REFERENCES messenger_users(id) ON DELETE CASCADE,
    decided_by UUID REFERENCES messenger_users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    decided_at TIMESTAMPTZ,
    PRIMARY KEY (channel_id, user_id)
);

CREATE INDEX idx_cjr_channel_pending ON messenger_channel_join_requests(channel_id) WHERE status = 'PENDING';
CREATE INDEX idx_cjr_user_pending    ON messenger_channel_join_requests(user_id)    WHERE status = 'PENDING';
