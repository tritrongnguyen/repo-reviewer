package domain

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password,omitempty"`
	Role         int       `json:"role:"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserProfile struct {
	FullName  string `json:"full_name"`
	AvatarUrl string `json:"avatar_url"`
	Phone     string `json:"phone"`
}

type UserConfig struct {
	APIKey string `json:"string,omitempty"`
}

type UserDetails struct {
	User
	Profile UserProfile `json:"profile"`
	Config  UserConfig  `json:"config"`
}
