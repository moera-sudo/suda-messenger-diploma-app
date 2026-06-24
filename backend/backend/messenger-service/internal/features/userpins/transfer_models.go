package userpins

import "github.com/google/uuid"

type CreatePinReq struct {
	PinType    string    `json:"pin_type"    validate:"required,oneof=SIDEBAR CHATLIST APP_HUB"`
	TargetType string    `json:"target_type" validate:"required,oneof=CHAT APP"`
	TargetID   uuid.UUID `json:"target_id"   validate:"required"`
}

type ReorderItem struct {
	PinID     int64 `json:"pin_id"     validate:"required"`
	SortOrder int   `json:"sort_order"`
}

type ReorderPinsReq struct {
	Items []ReorderItem `json:"items" validate:"required,min=1,dive"`
}