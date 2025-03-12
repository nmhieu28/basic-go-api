package responses

import (
	"time"

	"github.com/google/uuid"
)

type AuthenResponse struct {
	AccessToken        string `json:"accessToken"`
	RefreshToken       string `json:"-"`
	TokenExpire        int64  `json:"-"`
	RefreshTokenExpire int64  `json:"-"`
}

type UserResponse struct {
	Id          uuid.UUID  `json:"id"`
	Email       string     `json:"email"`
	FullName    string     `json:"fullName"`
	Avatar      string     `json:"avatar,omitempty"`
	DateOfBirth *time.Time `json:"dateOfBirth,omitempty"`
}
