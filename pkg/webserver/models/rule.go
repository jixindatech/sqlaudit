package models

import (
	"github.com/jinzhu/gorm"
)

/*
export const sqlOptions = [
  { key: 1, value: 'SELECT' },
  { key: 2, value: 'UNION' },
  { key: 3, value: 'INSERT' },
  { key: 4, value: 'UPDATE' },
  { key: 5, value: 'DELETE' },
  { key: 6, value: 'DDL' },
  { key: 7, value: 'SHOW' },
  { key: 8, value: 'USE' },
  { key: 9, value: 'SET' },
  { key: 10, value: 'BEGIN' },
  { key: 11, value: 'COMMIT' },
  { key: 12, value: 'ROLLBACK' },
  { key: 13, value: 'OTHERREAD' },
  { key: 14, value: 'OTHERADMIN' }
]
*/
type Rule struct {
	Model

	Name  string `json:"name" gorm:"column:name;not null"`
	Time  int    `json:"time" gorm:"column:time"`
	Type  int    `json:"type" gorm:"column:type;default 0"`
	User  string `json:"user" gorm:"column:user;default:''"`
	IP    string `json:"ip"   gorm:"column:ip;default:''"`
	Db    string `json:"db"   gorm:"column:db;default:''"`
	Op    int    `json:"op" gorm:"column:op;default:0"`
	Alert int    `json:"alert" gorm:"column:alert;default:0"`

	Sql      string `json:"sql" gorm:"column:sql;default:''"`
	Match    int    `json:"match" gorm:"column:match;default:0"`
	Priority int    `json:"priority" gorm:"column:priority;default:0"`

	Remark string `json:"remark" gorm:"column:remark;"`
}

func AddRule(data map[string]interface{}) error {
	rule := Rule{
		Name:     data["name"].(string),
		Time:     data["time"].(int),
		Type:     data["type"].(int),
		User:     data["user"].(string),
		IP:       data["ip"].(string),
		Db:       data["db"].(string),
		Op:       data["op"].(int),
		Sql:      data["sql"].(string),
		Match:    data["match"].(int),
		Priority: data["priority"].(int),
		Alert:    data["alert"].(int),
		Remark:   data["remark"].(string),
	}

	err := db.Create(&rule).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteRule(id uint) error {
	err := db.Where("id = ?", id).Delete(Rule{}).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateRule(id uint, data map[string]interface{}) error {
	if err := db.Model(&Rule{}).Where("id = ?", id).Update(data).Error; err != nil {
		return err
	}
	return nil
}

func GetRule(id uint) (*Rule, error) {
	var rule Rule
	if err := db.Where("id = ? ", id).Find(&rule).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &rule, nil
}

const (
	SORT_UNKNOWN = iota
	TIME_DESC
	TIME_ASC
	PRIORITY_DESC
	PRIORITY_ASC
)

func GetRules(query map[string]interface{}, page int, pageSize int) ([]*Rule, error) {
	var rules []*Rule
	var err error

	if page == 0 || pageSize == 0 {
		err = db.Order("priority DESC").Find(&rules).Error
	} else {
		pageNum := (page - 1) * pageSize
		name := query["name"].(string)
		switch query["sort"].(int) {
		case TIME_DESC:
			if len(name) > 0 {
				name = "%" + name + "%"
				err = db.Where("name like ?", name).Order("created_at DESC").Offset(pageNum).Limit(pageSize).Find(&rules).Error
			} else {
				err = db.Offset(pageNum).Limit(pageSize).Order("created_at DESC").Find(&rules).Error
			}
		case TIME_ASC:
			if len(name) > 0 {
				name = "%" + name + "%"
				err = db.Where("name like ?", name).Order("created_at ASC").Offset(pageNum).Limit(pageSize).Find(&rules).Error
			} else {
				err = db.Offset(pageNum).Limit(pageSize).Order("created_at ASC").Find(&rules).Error
			}
		case PRIORITY_DESC:
			if len(name) > 0 {
				name = "%" + name + "%"
				err = db.Where("name like ?", name).Order("priority DESC").Offset(pageNum).Limit(pageSize).Find(&rules).Error
			} else {
				err = db.Offset(pageNum).Limit(pageSize).Order("priority DESC").Find(&rules).Error
			}
		case PRIORITY_ASC:
			if len(name) > 0 {
				name = "%" + name + "%"
				err = db.Where("name like ?", name).Order("priority ASC").Offset(pageNum).Limit(pageSize).Find(&rules).Error
			} else {
				err = db.Offset(pageNum).Limit(pageSize).Order("priority ASC").Find(&rules).Error
			}
		default:
			if len(name) > 0 {
				name = "%" + name + "%"
				err = db.Where("name like ?", name).Offset(pageNum).Limit(pageSize).Find(&rules).Error
			} else {
				err = db.Offset(pageNum).Limit(pageSize).Find(&rules).Error
			}
		}
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return rules, nil
}

func GetRuleCount(query map[string]interface{}) (int, error) {
	var err error
	count := 0
	name := query["name"].(string)
	if len(name) > 0 {
		name = "%" + name + "%"
	}
	if len(name) > 0 {
		err = db.Model(&Rule{}).Where("name like ?", name).Count(&count).Error
	} else {
		err = db.Model(&Rule{}).Count(&count).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return count, err
}
