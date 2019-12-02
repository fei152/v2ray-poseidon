package ssrpanel

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type MySQLConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

func (c *MySQLConfig) FormatDSN() (string, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return "", err
	}

	cc := &mysql.Config{
		Collation:            "utf8mb4_unicode_ci",
		User:                 c.User,
		Passwd:               c.Password,
		Loc:                  loc,
		DBName:               c.DBName,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%d", c.Host, c.Port),
		AllowNativePasswords: true,
	}

	return cc.FormatDSN(), nil
}

func NewMySQLConn(config *MySQLConfig) (*DB, error) {
	newError("Connecting database...").AtInfo().WriteToLog()

	dsn, err := config.FormatDSN()
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SingularTable(true)
	defer newError("Connected").AtInfo().WriteToLog()

	return &DB{DB: db}, nil
}
