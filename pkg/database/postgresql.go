package database

import (
	"fmt"
	"log"

	configs "backend/pkg/config"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Database     *gorm.DB
	connAttempts int
	connTimeout  time.Duration
}

const (
	_defaultConnAttempts = 3
	_defaultConnTimeout  = time.Second
)

type DatabaseConnectionString string

func NewDatabase(appConfig *configs.AppConfig) (DBEngine, func(), error) {

	config := appConfig.Postgresql

	connectionString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		config.Host,
		config.Port,
		config.UserName,
		config.DBName,
		config.Password,
	)

	instance := &Database{
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	for instance.connAttempts > 0 {
		database, err := gorm.Open(postgres.Open(connectionString))

		if err == nil {
			instance.Database = database
			break
		}
		log.Printf("Postgres is trying to connect, attempts left: %d", instance.connAttempts)
		time.Sleep(instance.connTimeout)
		instance.connAttempts--

	}
	return instance, instance.Close, nil
}
func (instance *Database) Close() {
	if instance.Database != nil {
		sqlDB, _ := instance.Database.DB()
		sqlDB.Close()
		fmt.Println("Authen DB connection closed.")
	}
}
func (p *Database) Configure(opts ...Option) DBEngine {
	for _, opt := range opts {
		opt(p)
	}

	return p
}
func (p *Database) GetDatabase() *gorm.DB {
	return p.Database
}

func GetConnectionString(databaseSetting configs.PostgresConfig) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		databaseSetting.Host, databaseSetting.UserName, databaseSetting.Password, databaseSetting.DBName, strconv.Itoa(databaseSetting.Port), "UTC")
}
func (p *Database) Migrate(types ...interface{}) error {

	if p.Database == nil {
		return fmt.Errorf("Database connection is nil")
	}

	for _, t := range types {
		err := p.Database.AutoMigrate(t)
		if err != nil {
			return err
		}
	}
	return nil
}
