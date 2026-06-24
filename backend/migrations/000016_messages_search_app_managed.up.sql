-- 16_messages_search_app_managed.sql
-- Message content is now encrypted at rest, so the DB trigger that builds
-- search_vector FROM content would index ciphertext. The application now builds
-- search_vector from the plaintext at write time (SaveMessage/EditMessage), so
-- drop the trigger. The function is kept so the down migration can restore it.

DROP TRIGGER IF EXISTS trg_messages_search ON messenger_messages;
