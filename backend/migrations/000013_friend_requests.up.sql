-- Friend Requests / Friendships (Facebook-like model).
-- Lifecycle: PENDING -> ACCEPTED  |  PENDING -> REJECTED  |  PENDING -> CANCELLED
-- Friendship == single row with status = 'ACCEPTED' (direction kept for audit).
CREATE TABLE messenger_friend_requests (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    requester_id UUID NOT NULL REFERENCES messenger_users(id) ON DELETE CASCADE,
    target_id    UUID NOT NULL REFERENCES messenger_users(id) ON DELETE CASCADE,
    status       TEXT NOT NULL DEFAULT 'PENDING'
                 CHECK (status IN ('PENDING','ACCEPTED','REJECTED','CANCELLED')),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (requester_id <> target_id),
    UNIQUE (requester_id, target_id)
);

CREATE INDEX idx_friend_requests_target_status    ON messenger_friend_requests(target_id, status);
CREATE INDEX idx_friend_requests_requester_status ON messenger_friend_requests(requester_id, status);
