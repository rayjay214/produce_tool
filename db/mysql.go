package db

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var MysqlConn *gorm.DB

type TestRecord struct {
	RecordID    uint      `gorm:"column:record_id;primaryKey;autoIncrement"`
	Pass        string    `gorm:"column:pass"`
	Version     string    `gorm:"column:version"`
	Sim         string    `gorm:"column:sim"`
	Imei        string    `gorm:"column:imei"`
	Sn          string    `gorm:"column:sn"`
	Signal      string    `gorm:"column:signal"`
	Gps         string    `gorm:"column:gps"`
	Gsensor     string    `gorm:"column:gsensor"`
	Wifi        string    `gorm:"column:wifi"`
	Light       string    `gorm:"column:light"`
	MainIp      string    `gorm:"column:main_ip"`
	ViceIp      string    `gorm:"column:vice_ip"`
	SetType     string    `gorm:"column:set_type"`
	Power       string    `gorm:"column:power"`
	Protocol    string    `gorm:"column:protocol"`
	SetMainIp   string    `gorm:"column:set_main_ip"`
	SetViceIp   string    `gorm:"column:set_vice_ip"`
	SetProtocol string    `gorm:"column:set_protocol"`
	Operator    string    `gorm:"column:operator"`
	UploadWay   string    `gorm:"column:upload_way"`
	PlanId      uint      `gorm:"column:plan_id"`
	CreateTime  time.Time `gorm:"column:create_time;primaryKey"`
}

func (TestRecord) TableName() string {
	return "test_record"
}

func InitMysql() {
	dsn := "admin:qxyc@tcp(8.130.23.234:8000)/factory?charset=utf8mb4&parseTime=True&loc=Local&timeout=3s"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Errorf("create mysql conn failed %v\n", err)
		return
	}
	MysqlConn = db
	log.Info("create mysql conn success")
}

func InsertRecordMysql(record TestRecord) {
	if MysqlConn == nil {
		log.Error("mysql conn invalid")
		return
	}
	result := MysqlConn.Create(&record)
	if result.Error != nil {
		log.Errorf("insert failed:%v", result.Error)
	} else {
		log.Info("insert success")
	}
}
