ALTER TABLE messenger_user_preferences
    ADD COLUMN who_can_add_to_groups VARCHAR(20) DEFAULT 'EVERYONE',
    ADD COLUMN preview_in_notifications BOOLEAN DEFAULT TRUE,
    ADD COLUMN nft_showcase_enabled BOOLEAN DEFAULT FALSE,
    ADD COLUMN accept_gifts_from VARCHAR(20) DEFAULT 'EVERYONE',
    ADD COLUMN wallpaper_id UUID;
