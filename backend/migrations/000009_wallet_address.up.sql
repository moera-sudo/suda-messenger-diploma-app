-- Добавляем колонку wallet_address для денормализации.
-- Заполняется при verify пользователя (после CreateWalletForUser в transaction-service).
-- NULL означает: кошелёк ещё не создан (юзер либо не верифицирован,
-- либо CreateWalletForUser упал и его нужно повторить через "Активировать кошелёк" в UI).
ALTER TABLE messenger_users
    ADD COLUMN wallet_address VARCHAR(42);

-- Индекс для обратного поиска: по адресу → юзер.
-- Используется например когда фронту нужно узнать «чей это адрес».
-- WHERE-clause держит индекс компактным (не индексируем NULL'ы).
CREATE UNIQUE INDEX idx_messenger_users_wallet_address
    ON messenger_users(wallet_address)
    WHERE wallet_address IS NOT NULL;
