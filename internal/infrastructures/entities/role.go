package entities

import "backend/pkg/entity"

type Role struct {
	entity.BaseAuditTrackingEntity
	Name string `json:"name" gorm:"type:varchar(100);not null;"`
	Code string `json:"code" gorm:"type:varchar(100);not null;"`
}

func (Role) TableName() string {
	return "authentication.roles"
}
