-- Restore the DB-side search_vector trigger (messages_search_trigger() still exists).
CREATE TRIGGER trg_messages_search
    BEFORE INSERT OR UPDATE OF content
    ON messenger_messages
    FOR EACH ROW
    EXECUTE FUNCTION messages_search_trigger();
