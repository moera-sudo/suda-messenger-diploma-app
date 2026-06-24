// Package friends — Facebook-like friendship model.
//
// Lifecycle of messenger_friend_requests row:
//   PENDING -> ACCEPTED   (target user accepts)
//   PENDING -> REJECTED   (target user rejects)
//   PENDING -> CANCELLED  (requester revokes their own sent request)
//
// Friendship "exists" between A and B iff there is a row with status='ACCEPTED'
// for the unordered pair {A,B}. Direction (requester/target) is kept for audit
// but is irrelevant once accepted.
package friends

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Status of a friend-request row.
const (
	StatusPending   = "PENDING"
	StatusAccepted  = "ACCEPTED"
	StatusRejected  = "REJECTED"
	StatusCancelled = "CANCELLED"
)

// RelationStatus — derived 1-on-1 relationship between two users (returned by
// GET /users/friends/{user_id}/status).
const (
	RelationNone            = "NONE"
	RelationPendingSent     = "PENDING_SENT"     // current user has sent a request waiting on target
	RelationPendingReceived = "PENDING_RECEIVED" // target has sent a request waiting on current user
	RelationFriends         = "FRIENDS"
	RelationBlocked         = "BLOCKED"
)

// FriendRequest is the persisted record.
type FriendRequest struct {
	ID          uuid.UUID `json:"id"`
	RequesterID uuid.UUID `json:"requester_id"`
	TargetID    uuid.UUID `json:"target_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Domain errors — handler maps these to HTTP codes.
var (
	ErrSelfRequest         = errors.New("cannot send friend request to self")
	ErrBlocked             = errors.New("blocked")
	ErrAlreadyFriends      = errors.New("already friends")
	ErrAlreadyPending      = errors.New("request already pending")
	ErrRequestNotFound     = errors.New("friend request not found")
	ErrForbidden           = errors.New("forbidden")
	ErrInvalidTransition   = errors.New("invalid status transition")
)
