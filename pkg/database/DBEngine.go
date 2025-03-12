package database

import (
	"time"

	"gorm.io/gorm"
)

type Option func(*Database)

func ConnAttempts(attempts int) Option {
	return func(p *Database) {
		p.connAttempts = attempts
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(p *Database) {
		p.connTimeout = timeout
	}
}

type DBEngine interface {
	GetDatabase() *gorm.DB
	Configure(...Option) DBEngine
	Close()
	Migrate(types ...any) error
}
