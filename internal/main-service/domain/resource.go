package domain

import "time"

type Resource struct {
    ID         int64     `json:"id" db:"id"`
    Name       string    `json:"name" db:"name"`
    Type       string    `json:"type" db:"type"`
    Region     string    `json:"region" db:"region"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
    UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}