// Package donation — донаты SUDA: P2P (юзеру) или каналу (на treasury канала).
//
// Технически донат — это обычный SudaToken.Transfer, но:
//   - tx_pending помечается expected_type="DONATION" + donation_message;
//   - observer при индексации пишет дополнительную строку в tx_donations
//     (плоская история для статистики treasury) и создаёт system message
//     типа DONATION в чате/канале.
package donation

// EventType для NotifyUserEvent (согласовано с messenger-service).
const (
	EventDonationSent     = "DONATION_SENT"
	EventDonationReceived = "DONATION_RECEIVED"
)

// SystemMessageDonation — тип system message, который messenger создаёт в чате.
const SystemMessageDonation = "DONATION"

// TxTypeDonation — значение tx_transactions.type.
const TxTypeDonation = "DONATION"
