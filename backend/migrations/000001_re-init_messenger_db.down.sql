DROP INDEX IF EXISTS idx_messages_reply;
DROP INDEX IF EXISTS idx_messages_attachment;
DROP INDEX IF EXISTS idx_sessions_refresh_token;
DROP INDEX IF EXISTS idx_chat_members_user_id;
DROP INDEX IF EXISTS idx_messages_chat_id;
DROP INDEX IF EXISTS idx_verifications_code;
DROP INDEX IF EXISTS idx_verifications_user_purpose;

DROP TABLE IF EXISTS messenger_messages;
DROP TABLE IF EXISTS messenger_chat_members;
DROP TABLE IF EXISTS messenger_chats;
DROP TABLE IF EXISTS messenger_user_devices;
DROP TABLE IF EXISTS messenger_sessions;
DROP TABLE IF EXISTS messenger_verifications;
DROP TABLE IF EXISTS messenger_users;

DROP EXTENSION IF EXISTS "uuid-ossp";
