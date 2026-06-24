DROP INDEX IF EXISTS idx_message_hidden_user;
DROP INDEX IF EXISTS idx_blocked_users_blocker;
DROP INDEX IF EXISTS idx_blocked_users_blocked;
DROP INDEX IF EXISTS idx_pinned_messages_chat;
DROP INDEX IF EXISTS idx_reactions_message;
DROP INDEX IF EXISTS idx_messages_deleted;
DROP INDEX IF EXISTS idx_messages_edited;
DROP INDEX IF EXISTS idx_chats_creator;
DROP INDEX IF EXISTS idx_chats_type;
DROP INDEX IF EXISTS idx_contacts_owner;

DROP TABLE IF EXISTS messenger_stickers;
DROP TABLE IF EXISTS messenger_sticker_packs;
DROP TABLE IF EXISTS messenger_reactions;
DROP TABLE IF EXISTS messenger_pinned_messages;
DROP TABLE IF EXISTS messenger_blocked_users;
DROP TABLE IF EXISTS messenger_contacts;
DROP TABLE IF EXISTS messenger_message_hidden;

ALTER TABLE messenger_chat_members
    DROP COLUMN IF EXISTS is_muted,
    DROP COLUMN IF EXISTS muted_until,
    DROP COLUMN IF EXISTS notifications_enabled;

ALTER TABLE messenger_messages
    DROP COLUMN IF EXISTS edited_at,
    DROP COLUMN IF EXISTS deleted_at;

ALTER TABLE messenger_chats
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS avatar_media_id,
    DROP COLUMN IF EXISTS creator_id,
    DROP COLUMN IF EXISTS invite_link;
