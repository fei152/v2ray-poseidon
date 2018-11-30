package v2ray_ssrpanel_plugin

import (
	"time"

	"github.com/jinzhu/gorm"
)

type UserModel struct {
	ID      uint
	VmessID string
	Email   string `gorm:"column:username"`
}

func (UserModel) TableName() string {
	return "user"
}

type UserTrafficLog struct {
	ID        uint `gorm:"primary_key"`
	UserID    uint
	Uplink    uint64 `gorm:"column:u"`
	Downlink  uint64 `gorm:"column:d"`
	NodeID    uint
	Rate      float64
	Traffic   string
	CreatedAt time.Time `gorm:"column:log_time"`
}

type DB struct {
	DB *gorm.DB
}

func (db *DB) GetAllUsers() ([]UserModel, error) {
	users := make([]UserModel, 0)
	db.DB.Select("id, vmess_id, username").Where("enable = 1 AND u + d < transfer_enable").Find(&users)
	return users, nil
}

func (db *DB) CreateUserTrafficLog(log *UserTrafficLog) bool {
	db.DB.Create(log)
	return db.DB.NewRecord(log)
}
