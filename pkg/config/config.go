package configs

import (
	"backend/pkg/environment"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

func LoadConfig(fileName string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(fileName)
	v.SetConfigType("yml")
	fmt.Println(fileName)
	v.AddConfigPath("../config/")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

func InitAppConfig() (*AppConfig, error) {
	configPath := environment.GetEnvironment().String()
	v := viper.New()
	v.SetConfigName(configPath)
	v.SetConfigType("yml")
	fmt.Println(configPath)
	v.AddConfigPath("../config/")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	var c AppConfig

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}

type AppConfig struct {
	Server     ServerConfig     `mapstructure:"server"`
	Postgresql PostgresConfig   `mapstructure:"postgresql"`
	Logger     LoggerConfig     `mapstructure:"logger"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Jwt        JWTConfig        `mapstructure:"jwt"`
	Smtp       SMTPConfig       `mapstructure:"smtp"`
	Cors       CORSConfig       `mapstructure:"cors"`
	ServiceUrl ServiceUrlConfig `mapstructure:"serviceUrl"`
}
type PostgresConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	UserName        string `mapstructure:"userName"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbName"`
	MaxIdleConns    int    `mapstructure:"maxIdleConns"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns"`
	ConnMaxLifeTime int    `mapstructure:"connMaxLifeTime"`
	Driver          string `mapstructure:"driver"`
	Schema          string `mapstructure:"schema"`
	SSLMode         bool   `mapstructure:"sslMode"`
}

type LoggerConfig struct {
	Development       bool   `mapstructure:"development"`
	DisableCaller     bool   `mapstructure:"disableCaller"`
	DisableStacktrace bool   `mapstructure:"disableStacktrace"`
	Encoding          string `mapstructure:"encoding"`
	Level             string `mapstructure:"level"`
	MaxSize           int    `mapstructure:"maxSize"`
	MaxAge            int    `mapstructure:"maxAge"`
	MaxBackups        int    `mapstructure:"maxBackups"`
	FileName          string `mapstructure:"fileName"`
	Compress          bool   `mapstructure:"compress"`
}
type ServerConfig struct {
	AppVersion        string        `mapstructure:"appVersion"`
	ServiceName       string        `mapstructure:"serviceName"`
	Port              string        `mapstructure:"port"`
	Mode              string        `mapstructure:"mode"`
	ReadTimeout       time.Duration `mapstructure:"readTimeout"`
	WriteTimeout      time.Duration `mapstructure:"writeTimeout"`
	SSL               bool          `mapstructure:"ssl"`
	CtxDefaultTimeout time.Duration `mapstructure:"ctxDefaultTimeout"`
	CSRF              bool          `mapstructure:"csrf"`
	Debug             bool          `mapstructure:"debug"`
}

type JWTConfig struct {
	SecretKey              string `mapstructure:"secretKey"`
	TokenExpire            int    `mapstructure:"tokenExpire"`
	RefreshTokenExpire     int    `mapstructure:"refreshTokenExpire"`
	Audience               string `mapstructure:"audience"`
	Issuer                 string `mapstructure:"issuer"`
	RefreshSecretKey       string `mapstructure:"refreshSecretKey"`
	VerifyEmailSecretKey   string `mapstructure:"verifyEmailSecretKey"`
	VerifyEmailTokenExpire int    `mapstructure:"verifyEmailTokenExpire"`
}

type SMTPConfig struct {
	UserName string `mapstructure:"userName"`
	Password string `mapstructure:"password"`
	Service  string `mapstructure:"service"`
	Port     string `mapstructure:"port"`
	Host     string `mapstructure:"host"`
	From     string `mapstructure:"from"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	UserName string `mapstructure:"userName"`
	Password string `mapstructure:"password"`
}

type CORSConfig struct {
	Enable bool     `mapstructure:"enable"`
	Allows []string `mapstructure:"allows"`
}
type ServiceUrlConfig struct {
	Frontend string `mapstructure:"frontend"`
}
