package models

import "time"

type UserRole string

const (
	RoleUser  UserRole = "USER"
	RoleAdmin UserRole = "ADMIN"
)

type User struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	FirebaseUID string    `bson:"firebase_uid" json:"firebase_uid"`
	Email       string    `bson:"email" json:"email"`
	Name        string    `bson:"name" json:"name"`
	Role        UserRole  `bson:"role" json:"role"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
