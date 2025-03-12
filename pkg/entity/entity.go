package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DbContext interface {
	Connect() error
	Close()
	Database() *gorm.DB
	MigrateTable(types ...interface{})
}

const (
	KEY_AUTHENTICATION = "auth"
	KEY_CUSTOMER       = "customer"
)
const (
	DRIVER_POSTGRESQL = "postgresql"
)

type DateTimeTracking struct {
	CreatedDateTimeUtc *time.Time `json:"createdDateTimeUtc,omitempty" gorm:"not null;"`
	UpdatedDateTimeUtc *time.Time `json:"updatedDateTimeUtc,omitempty" gorm:"not null;"`
}
type Entity struct {
	Id uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;"`
}
type AuditTracking struct {
	CreatedBy uuid.NullUUID `json:"createdBy,omitempty" gorm:"type:uuid;"`
	UpdatedBy uuid.NullUUID `json:"updatedBy,omitempty" gorm:"type:uuid;"`
}
type SoftDelete struct {
	IsDeleted          bool       `json:"isDeleted" gorm:"default:false;not null;"`
	DeletedDateTimeUtc *time.Time `json:"deletedDateTimeUtc,omitempty" gorm:"null;"`
}
type Multitenant struct {
	Id uuid.UUID `json:"tenantId" gorm:"type:uuid;null;"`
}
type BaseAuditTrackingEntity struct {
	Entity
	DateTimeTracking `gorm:"embedded"`
	AuditTracking    `gorm:"embedded"`
}
type BaseEntitySoftDelete struct {
	BaseAuditTrackingEntity `gorm:"embedded"`
	SoftDelete              `gorm:"embedded"`
}

func NewSQLModel() BaseAuditTrackingEntity {
	now := time.Now().UTC()
	id := uuid.New()
	return BaseAuditTrackingEntity{
		Entity: Entity{
			Id: id,
		},
		DateTimeTracking: DateTimeTracking{
			CreatedDateTimeUtc: &now,
			UpdatedDateTimeUtc: &now,
		},
		AuditTracking: AuditTracking{
			CreatedBy: uuid.NullUUID{Valid: false},
			UpdatedBy: uuid.NullUUID{Valid: false},
		},
	}
}
