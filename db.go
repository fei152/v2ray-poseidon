package v2ray_ssrpanel_plugin

import (
	"github.com/jinzhu/gorm"
	"time"
)

type UserModel struct {
	ID      uint
	VmessID string
	Email   string `gorm:"column:username"`
}

func (*UserModel) TableName() string {
	return "user"
}

type UserTrafficLog struct {
	ID       uint `gorm:"primary_key"`
	UserID   uint
	Uplink   uint64 `gorm:"column:u"`
	Downlink uint64 `gorm:"column:d"`
	NodeID   uint
	Rate     float64
	Traffic  string
	LogTime  int64
}

func (l *UserTrafficLog) BeforeCreate(scope *gorm.Scope) error {
	l.LogTime = time.Now().Unix()
	return nil
}

type DB struct {
	DB *gorm.DB
}

func (db *DB) GetAllUsers() ([]UserModel, error) {
	users := make([]UserModel, 0)
	db.DB.Select("id, vmess_id, username").Where("enable = 1 AND u + d < transfer_enable").Find(&users)
	return users, nil
}
