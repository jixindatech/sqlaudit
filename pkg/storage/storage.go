package storage

import (
	"github.com/jixindatech/sqlaudit/pkg/config"
)

type Storage interface {
	InitStorage(cfg *config.EsConfig) error
	Query(query map[string]interface{}, page, size int) (map[string]interface{}, error)
	QueryInfo(map[string]interface{}) ([]byte, error)
	QueryFingerPrintInfo(map[string]interface{}, int, int) (map[string]interface{}, error)
	Save(body interface{}) error
}

func GetStorage(cfg *config.EsConfig) (Storage, error) {
	var storage Storage
	storage = &EsStorage{}
	err := storage.InitStorage(cfg)
	if err != nil {
		return nil, err
	}
	return storage, nil
}
