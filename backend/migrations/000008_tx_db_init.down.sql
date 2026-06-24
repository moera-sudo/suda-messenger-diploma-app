-- ════════════════════════════════════════════════════════════
--  Transaction Service — REVERT initial schema
-- ════════════════════════════════════════════════════════════

-- 1. Сначала удаляем таблицы, которые ссылаются на другие (зависимые таблицы)
DROP TABLE IF EXISTS tx_gating_rules;
DROP TABLE IF EXISTS tx_marketplace_listings;
DROP TABLE IF EXISTS tx_nft_items;
DROP TABLE IF EXISTS tx_nft_collections;

DROP TABLE IF EXISTS tx_fundraiser_contributions;
DROP TABLE IF EXISTS tx_fundraisers;

DROP TABLE IF EXISTS tx_quests;
DROP TABLE IF EXISTS tx_donations;

-- 2. Удаляем независимые таблицы
DROP TABLE IF EXISTS tx_observer_state;
DROP TABLE IF EXISTS tx_signing_audit;
DROP TABLE IF EXISTS tx_pending;

-- 3. Удаляем таблицу транзакций (правило tx_transactions_no_delete 
--    и триггер trg_tx_transactions_immutable удалятся автоматически вместе с ней)
DROP TABLE IF EXISTS tx_transactions;

-- 4. Удаляем функцию для триггера
DROP FUNCTION IF EXISTS tx_transactions_immutable();

-- 5. Удаляем кошельки
DROP TABLE IF EXISTS tx_channel_wallets;
DROP TABLE IF EXISTS tx_wallets;

-- ════════════════════════════════════════════════════════════
-- ВНИМАНИЕ: Удаление расширения uuid-ossp закомментировано,
-- так как оно может использоваться другими сервисами или таблицами
-- в этой же базе данных. Раскомментируйте, если уверены, что 
-- эта БД используется исключительно данным сервисом.
-- ════════════════════════════════════════════════════════════
-- DROP EXTENSION IF EXISTS "uuid-ossp";
