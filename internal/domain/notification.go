package domain

import "time"

type Notification struct {
    Event string `json:"event" db:"event"`
    ID        int64     `json:"id" db:"id"`
    UserID    int64     `json:"user_id" db:"user_id"`
    Message   string    `json:"message" db:"message"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}