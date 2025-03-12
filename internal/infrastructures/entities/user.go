package entities

import (
	"backend/pkg/entity"
	"fmt"
	"time"
)

var (
	Account_Internal = 0
	Account_Client   = 100
)

type User struct {
	entity.BaseAuditTrackingEntity
	FirstName         string     `json:"firstName" gorm:"type:varchar(100);not null;"`
	LastName          string     `json:"lastName" gorm:"type:varchar(100);not null;"`
	UserName          string     `json:"userName" gorm:"type:varchar(100);not null;"`
	DateOfBirth       *time.Time `json:"dateOfBirth,omitempty" gorm:"null;"`
	Email             string     `json:"email" gorm:"type:varchar(256);not null;"`
	UserTypeID        int16      `json:"userTypeId" gorm:"type:smallint;default:1;not null;"`
	Avatar            string     `json:"avatar,omitempty" gorm:"type:varchar(1024);"`
	TwoFactorEnabled  bool       `json:"twoFactorEnabled" gorm:"default:false;not null;"`
	LockoutEnd        *time.Time `json:"lockoutEnd,omitempty" gorm:"null;"`
	LockoutEnabled    bool       `json:"lockoutEnabled" gorm:"default:false;not null;"`
	AccessFailedCount int16      `json:"accessFailedCount" gorm:"type:smallint;default:0;not null;"`
	EmailConfirm      bool       `json:"emailConfirm" gorm:"default:false;not null;"`
	PasswordHash      string     `json:"passwordHash" gorm:"type:varchar(100);not null;"`
	TimeZoneID        int16      `json:"timeZoneId,omitempty" gorm:"type:smallint;null;"`
}

func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

func (User) TableName() string {
	return "authentication.users"
}
