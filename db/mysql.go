package db

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var MysqlConn *gorm.DB

type TestRecord struct {
	RecordID     uint      `gorm:"column:record_id;primaryKey;autoIncrement"`
	Pass         string    `gorm:"column:pass"`
	Version      string    `gorm:"column:version"`
	Uuid         string    `gorm:"column:uuid"`
	CaliBand     string    `gorm:"column:caliband"`
	Sim          string    `gorm:"column:sim"`
	Imei         string    `gorm:"column:imei"`
	Sn           string    `gorm:"column:sn"`
	Signal       string    `gorm:"column:signal"`
	Gps          string    `gorm:"column:gps"`
	Gsensor      string    `gorm:"column:gsensor"`
	Wifi         string    `gorm:"column:wifi"`
	Light        string    `gorm:"column:light"`
	MainIp       string    `gorm:"column:main_ip"`
	ViceIp       string    `gorm:"column:vice_ip"`
	Power        string    `gorm:"column:power"`
	Protocol     string    `gorm:"column:protocol"`
	SetType      string    `gorm:"column:set_type"`
	SetMainIp    string    `gorm:"column:set_main_ip"`
	SetViceIp    string    `gorm:"column:set_vice_ip"`
	SetProtocol  string    `gorm:"column:set_protocol"`
	CreateTime   time.Time `gorm:"column:create_time;primaryKey"`
	Operator     string    `gorm:"column:operator"`
	UploadWay    string    `gorm:"column:upload_way"`
	PlanId       int64     `gorm:"column:plan_id"`
	SetApn       string    `gorm:"column:set_apn"`
	SetWifiPrior string    `gorm:"column:set_wifi_prior"`
}

func (TestRecord) TableName() string {
	return "smt_test_record"
}

func InitMysql() {
	dsn := "admin:qxyc@tcp(8.130.23.234:8000)/factory?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s"
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
