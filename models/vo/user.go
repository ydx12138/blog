package vo

import "time"

type UserProfile struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Phone     *string   `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}
