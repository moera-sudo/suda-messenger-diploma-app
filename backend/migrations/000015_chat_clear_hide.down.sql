ALTER TABLE messenger_chat_members
    DROP COLUMN IF EXISTS cleared_at,
    DROP COLUMN IF EXISTS hidden_at;
