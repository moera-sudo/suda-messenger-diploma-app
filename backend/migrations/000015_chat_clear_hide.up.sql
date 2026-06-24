-- 15_chat_clear_hide.sql
-- Per-user "delete chat for me": clear history + hide from chat list.
-- cleared_at — messages created at/before this are hidden from this member.
-- hidden_at  — chat is hidden from this member's list until a message arrives
--              newer than hidden_at (Telegram-style reappear-on-new-message).

ALTER TABLE messenger_chat_members
    ADD COLUMN IF NOT EXISTS cleared_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS hidden_at  TIMESTAMPTZ;
