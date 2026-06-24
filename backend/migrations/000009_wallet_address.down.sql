DROP INDEX IF EXISTS idx_messenger_users_wallet_address;
ALTER TABLE messenger_users DROP COLUMN IF EXISTS wallet_address;
