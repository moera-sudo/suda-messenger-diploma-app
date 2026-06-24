package auth

import "github.com/google/uuid"

type User struct {
	ID            uuid.UUID
	Username      string
	Email         string
	PasswordHash  string
	DisplayName   string
	IsVerified    bool
	Role          string
	AvatarMediaID *uuid.UUID
}
