package db

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
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

	//修改设备测试状态
	nSn, _ := strconv.ParseInt(record.Sn, 10, 64)
	result = MysqlConn.Table("device").Where("devno = ?", nSn).
		Update("test_status", "1")
	if result.Error != nil {
		log.Errorf("%v change status failed:%v", nSn, result.Error)
	}
}

func CheckSn(sn, planId int64) error {
	if MysqlConn == nil {
		return fmt.Errorf("数据库连接失败")
	}

	// 直接查询SN的详细信息
	var device struct {
		Devno       int64  `gorm:"column:devno"`
		PlanId      int64  `gorm:"column:plan_id"`
		WriteStatus string `gorm:"column:write_status"`
	}

	result := MysqlConn.Table("device").Select("devno, plan_id, write_status").Where("devno = ?", sn).First(&device)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("SN不存在")
		}
		return fmt.Errorf("查询SN失败: %v", result.Error)
	}

	if device.PlanId != planId {
		return fmt.Errorf("SN不属于当前计划")
	}

	if device.WriteStatus != "0" {
		return fmt.Errorf("SN已被使用")
	}

	// SN存在且状态正常，可以使用
	return nil
}

func WriteStatus(sn string) error {
	if MysqlConn == nil {
		return fmt.Errorf("数据库连接失败")
	}

	nSn, _ := strconv.ParseInt(sn, 10, 64)

	result := MysqlConn.Table("device").Where("devno = ?", nSn).
		Update("write_status", "1")

	if result.Error != nil {
		return fmt.Errorf("更新SN状态失败: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到SN为 %s 的记录", sn)
	}

	return nil
}

func CompareStatus(sn string) error {
	if MysqlConn == nil {
		return fmt.Errorf("数据库连接失败")
	}

	nSn, _ := strconv.ParseInt(sn, 10, 64)

	result := MysqlConn.Table("device").Where("devno = ?", nSn).
		Update("compare_status", "1")

	if result.Error != nil {
		return fmt.Errorf("更新SN状态失败: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到SN为 %s 的记录", sn)
	}

	return nil
}
