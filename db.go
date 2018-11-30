package v2ray_ssrpanel_plugin

import (
	"github.com/jinzhu/gorm"
)

type UserModel struct {
	ID      int
	VmessID string
}

func (UserModel) TableName() string {
	return "user"
}

type DB struct {
	*gorm.DB
}

func (db *DB) GetAllUsers() ([]UserModel, error) {
	users := make([]UserModel, 0)
	db.Select("id, vmess_id").Where("enable = 1 AND u + d < transfer_enable").Find(&users)
	return users, nil
}
