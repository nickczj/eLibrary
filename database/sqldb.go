package database

import (
	"eLibrary/global"
	"eLibrary/model"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func NewPostgres() {
	handler := dbHandler(start)
	dsn := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=verify-full",
		"nickczj",
		"SjjaKpb06K0W8hF66iTwKw",
		"free-tier6.gcp-asia-southeast1.cockroachlabs.cloud",
		"26257",
		"last-koala-2636.defaultdb",
	)
	handler.handleDb(postgres.Open(dsn))
}

func start(dialect gorm.Dialector) error {
	// initialize connection
	conf := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	database, err := gorm.Open(dialect, conf)
	if err != nil {
		return err
	}

	// settings
	sqlDB, err := database.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(3)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(2 * time.Hour)

	// test connection
	if err = sqlDB.Ping(); err != nil {
		return err
	}

	// synchronize DB schemas
	err = database.AutoMigrate(&model.BookDetail{})
	err = database.AutoMigrate(&model.LoanDetail{})
	err = database.AutoMigrate(&model.User{})
	if err != nil {
		return err
	}

	global.Database = database
	return nil
}

type dbHandler func(dialect gorm.Dialector) error

func (fn dbHandler) handleDb(dialect gorm.Dialector) {
	if err := fn(dialect); err != nil {
		log.Error(err)
	}
}
