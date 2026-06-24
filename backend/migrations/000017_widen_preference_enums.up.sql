-- 17_widen_preference_enums.sql
-- theme VARCHAR(10) overflowed on client values like "teaChatsDark" (SQLSTATE 22001).
-- Widen theme + a couple of other short enum-ish columns to match the API validation (max 64).

ALTER TABLE messenger_user_preferences
    ALTER COLUMN theme          TYPE VARCHAR(32),
    ALTER COLUMN chat_font_size TYPE VARCHAR(20),
    ALTER COLUMN lang_code      TYPE VARCHAR(10);
