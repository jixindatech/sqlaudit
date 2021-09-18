package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jixindatech/sqlaudit/pkg/storage"
	"time"
)

type Model struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updateAt"`
	// DeletedAt *time.Time `json:"deletedAt"`
}

type Database struct {
	DB string
}

var db *gorm.DB
var Storage storage.Storage

func OpenDatabase(database string, _storage storage.Storage) error {
	var err error
	db, err = gorm.Open("sqlite3", database)
	if err != nil {
		return err
	}

	db.SingularTable(true)
	db.LogMode(false)
	/*
		db.DB().SetMaxIdleConns(10)
		db.DB().SetMaxOpenConns(100)
	*/

	db.AutoMigrate(Rule{})

	Storage = _storage
	return nil
}
