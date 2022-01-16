package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jixindatech/sqlaudit/pkg/config"
	"github.com/jixindatech/sqlaudit/pkg/storage"
	"time"
)

type Model struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updateAt"`
	// DeletedAt *time.Time `json:"deletedAt"`
}

var db *gorm.DB
var Storage storage.Storage

func OpenDatabase(cfg *config.DataBase, _storage storage.Storage) error {
	var err error
	if cfg.Type == "mysql" {
		db, err = gorm.Open(cfg.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Name))
		if err != nil {
			return err
		}
	}

	if len(cfg.TablePrefix) > 0 {
		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			return cfg.TablePrefix + defaultTableName
		}
	}

	db.SingularTable(true)
	db.LogMode(false)

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	db.AutoMigrate(Rule{})
	db.AutoMigrate(FingerPrint{})

	Storage = _storage
	return nil
}
