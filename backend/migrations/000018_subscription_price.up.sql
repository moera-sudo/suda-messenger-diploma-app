-- 18_subscription_price.sql
-- Token gating becomes a PAID subscription for PUBLIC channels: subscribing
-- charges this price in SUDA (wei) to the channel treasury. A PUBLIC channel
-- with a rule and subscription_price_wei > 0 is a paid channel.

ALTER TABLE tx_gating_rules
    ADD COLUMN IF NOT EXISTS subscription_price_wei NUMERIC(78, 0) NOT NULL DEFAULT 0;
