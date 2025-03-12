package repositories

import (
	"backend/internal/infrastructures/entities"
	"backend/pkg/database"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	database.RepositoryBase[entities.User, uuid.UUID]
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
}
type userRepository struct {
	database.Repository[entities.User, uuid.UUID]
}

func NewUserRepository(dbEngine database.DBEngine) UserRepository {
	DbContext := dbEngine.GetDatabase()
	return &userRepository{
		Repository: *database.NewRepository[entities.User, uuid.UUID](DbContext),
	}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	user, err := r.First(&entities.User{Email: email}, ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}
