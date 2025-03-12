package migrations

import (
	"backend/internal/infrastructures/entities"
	"backend/pkg/database"
)

func GetModels() []any {
	return []any{
		&entities.User{},
		&entities.Role{},
		&entities.UserRole{},
	}
}
func Migrate(dbEngine database.DBEngine) error {
	return dbEngine.Migrate(GetModels()...)
}
