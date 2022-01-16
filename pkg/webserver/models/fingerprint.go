package models

import "github.com/jinzhu/gorm"

type FingerPrint struct {
	Model

	FingerPrint string `json:"fingerprint" gorm:"column:finger_print;default:''"`
	Remark      string `json:"remark" gorm:"column:remark;"`
}

func AddFingerPrint(data map[string]interface{}) error {
	fingerprint := FingerPrint{
		FingerPrint: data["fingerprint"].(string),
		Remark:      data["remark"].(string),
	}

	err := db.Create(&fingerprint).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteFingerPrint(id uint) error {
	err := db.Where("id = ?", id).Delete(FingerPrint{}).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateFingerPrint(id uint, data map[string]interface{}) error {
	if err := db.Model(&FingerPrint{}).Where("id = ?", id).Update(data).Error; err != nil {
		return err
	}
	return nil
}

func GetFingerPrint(id uint) (*FingerPrint, error) {
	var fingerprint FingerPrint
	if err := db.Where("id = ? ", id).Find(&fingerprint).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &fingerprint, nil
}

func GetFingerPrints(query map[string]interface{}, page int, pageSize int) ([]*FingerPrint, error) {
	var fingerprints []*FingerPrint
	var err error

	if page == 0 || pageSize == 0 {
		err = db.Find(&fingerprints).Error
	} else {
		pageNum := (page - 1) * pageSize
		err = db.Offset(pageNum).Limit(pageSize).Order("created_at DESC").Find(&fingerprints).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return fingerprints, nil
}

func GetFingerPrintCount(query map[string]interface{}) (int, error) {
	var err error
	count := 0
	err = db.Model(&FingerPrint{}).Count(&count).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return count, err
}
