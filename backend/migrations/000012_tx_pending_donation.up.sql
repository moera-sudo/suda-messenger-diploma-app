-- donation_message — текст доната, который пишет sender при вызове /donate.
-- Хранится в tx_pending до индексации Transfer event'а observer'ом; затем
-- переезжает в tx_donations.message и tx_pending row удаляется.
ALTER TABLE tx_pending
    ADD COLUMN donation_message TEXT;
