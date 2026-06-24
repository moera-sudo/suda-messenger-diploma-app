
DROP TRIGGER IF EXISTS trg_messages_search ON messenger_messages;
DROP TRIGGER IF EXISTS trg_users_search ON messenger_users;
DROP TRIGGER IF EXISTS trg_chats_search ON messenger_chats;



DROP FUNCTION IF EXISTS messages_search_trigger();
DROP FUNCTION IF EXISTS users_search_trigger();
DROP FUNCTION IF EXISTS chats_search_trigger();


DROP INDEX IF EXISTS idx_messages_search_vector;
DROP INDEX IF EXISTS idx_users_search_vector;
DROP INDEX IF EXISTS idx_chats_search_vector;

DROP INDEX IF EXISTS idx_last_seen_hidden_owner;



ALTER TABLE messenger_messages DROP COLUMN IF EXISTS search_vector;
ALTER TABLE messenger_users DROP COLUMN IF EXISTS search_vector;
ALTER TABLE messenger_chats DROP COLUMN IF EXISTS search_vector;



DROP TABLE IF EXISTS messenger_last_seen_hidden;
DROP TABLE IF EXISTS messenger_user_preferences;
