package database

import (
	"context"

	"gorm.io/gorm"
)

type RepositoryBase[T any, ID comparable] interface {
	// GetByID retrieves an entity by its ID
	GetByID(id ID, ctx context.Context) (*T, error)

	// Create inserts a new entity into the database
	Create(entity *T, ctx context.Context) (*T, error)

	// Update modifies an existing entity in the database
	Update(entity *T, ctx context.Context) error

	// Delete removes an entity from the database by its ID
	Delete(id ID, ctx context.Context) error

	// List retrieves all entities from the database
	List(ctx context.Context) (*[]T, error)

	// SkipTake implements pagination by skipping a number of records and taking a specified amount
	SkipTake(skip int, take int, ctx context.Context) (*[]T, error)

	// CountWhere counts the number of records matching the given parametersß
	CountWhere(params *T, ctx context.Context) int64

	// Where retrieves all entities matching the given parameters
	Where(params *T, ctx context.Context) (*[]T, error)

	// WhereNotDeleted retrieves all entities matching the given parameters and not deleted
	WhereNotDeleted(params *T, ctx context.Context) (*[]T, error)

	// First retrieves the first entity matching the given parameters
	First(params *T, ctx context.Context) (*T, error)
}

type Repository[T any, Id comparable] struct {
	DbContext *gorm.DB
}

func NewRepository[T any, Id comparable](db *gorm.DB) *Repository[T, Id] {
	return &Repository[T, Id]{DbContext: db}
}

func (r *Repository[T, Id]) GetByID(id Id, ctx context.Context) (*T, error) {
	var entity T
	result := r.DbContext.WithContext(ctx).First(&entity, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &entity, nil
}

func (r *Repository[T, Id]) Create(entity *T, ctx context.Context) (*T, error) {
	result := r.DbContext.WithContext(ctx).Create(entity)
	return entity, result.Error
}

func (r *Repository[T, Id]) Delete(id Id, ctx context.Context) error {
	var entity T
	result := r.DbContext.WithContext(ctx).Delete(&entity, id)
	return result.Error
}
func (r *Repository[T, Id]) List(ctx context.Context) (*[]T, error) {
	var entities *[]T
	result := r.DbContext.WithContext(ctx).Find(&entities)
	return entities, result.Error
}

func (r *Repository[T, ID]) Update(entity *T, ctx context.Context) error {
	result := r.DbContext.WithContext(ctx).Save(entity)
	return result.Error
}

func (r *Repository[T, ID]) SkipTake(skip int, take int, ctx context.Context) (*[]T, error) {
	var entities *[]T
	result := r.DbContext.WithContext(ctx).Offset(skip).Limit(take).Find(&entities)
	return entities, result.Error
}

func (r *Repository[T, ID]) Count(ctx context.Context) int64 {
	var entity T
	var count int64
	r.DbContext.WithContext(ctx).Model(&entity).Count(&count)
	return count
}

func (r *Repository[T, ID]) CountWhere(params *T, ctx context.Context) int64 {
	var entity T
	var count int64
	r.DbContext.WithContext(ctx).Model(&entity).Where(&params).Count(&count)
	return count
}

func (r *Repository[T, Id]) Where(params *T, ctx context.Context) (*[]T, error) {
	var entities []T
	err := r.DbContext.WithContext(ctx).Where(&params).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return &entities, nil
}

func (r *Repository[T, Id]) WhereNotDeleted(params *T, ctx context.Context) (*[]T, error) {
	var entities []T
	err := r.DbContext.WithContext(ctx).Where(params).Where("is_deleted = ?", false).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return &entities, nil
}

func (r *Repository[T, Id]) First(params *T, ctx context.Context) (*T, error) {
	var entity T
	err := r.DbContext.WithContext(ctx).Where(&params).FirstOrInit(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
