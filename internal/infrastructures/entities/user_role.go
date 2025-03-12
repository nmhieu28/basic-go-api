package entities

import (
	"backend/pkg/entity"

	"github.com/google/uuid"
)

type UserRole struct {
	entity.BaseAuditTrackingEntity
	UserId uuid.NullUUID `json:"userId,omitempty" gorm:"type:uuid;"`
	RoleId uuid.NullUUID `json:"roleId,omitempty" gorm:"type:uuid;"`
}

func (UserRole) TableName() string {
	return "authentication.user_roles"
}
