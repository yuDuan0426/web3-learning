package database

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config 数据库配置
type Config struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	Params   string
}

// NewMySQLConnection 创建MySQL连接
func NewMySQLConnection(config *Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		config.Params,
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// InitFromViper 从Viper配置初始化数据库连接
func InitFromViper() (*sqlx.DB, error) {
	config := &Config{
		Driver:   viper.GetString("database.driver"),
		Host:     viper.GetString("database.host"),
		Port:     viper.GetInt("database.port"),
		Username: viper.GetString("database.username"),
		Password: viper.GetString("database.password"),
		DBName:   viper.GetString("database.dbname"),
		Params:   viper.GetString("database.params"),
	}

	logrus.Infof("Connecting to %s database at %s:%d/%s", config.Driver, config.Host, config.Port, config.DBName)
	return NewMySQLConnection(config)
}