package db

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"produce_tool/network"
	"strconv"
	"time"
)

var MysqlConn *gorm.DB

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

	//修改设备测试状态
	var nSn int64
	if network.CurrentPlan.SnType == "1" {
		nSn, _ = strconv.ParseInt(record.Imei, 10, 64)
	} else {
		nSn, _ = strconv.ParseInt(record.Sn, 10, 64)
	}
	result = MysqlConn.Table("device").Where("devno = ?", nSn).
		Update("test_status", "1")
	if result.Error != nil {
		log.Errorf("%v change status failed:%v", nSn, result.Error)
	}
}

func baseCheck(sn int64, planId int64, device *Device) error {
	if MysqlConn == nil {
		log.Error("mysql conn invalid")
		return fmt.Errorf("数据库连接失败")
	}

	result := MysqlConn.Table("device").Select("*").Where("devno = ?", sn).First(device)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("SN不存在")
		}
		return fmt.Errorf("查询SN失败: %v", result.Error)
	}

	if device.PlanId != planId {
		return fmt.Errorf("SN不属于当前计划")
	}
	return nil
}

func CheckSn(sn, planId int64) error {
	var device Device
	err := baseCheck(sn, planId, &device)
	if err != nil {
		return err
	}

	if device.WriteStatus != "0" {
		return fmt.Errorf("SN已被使用")
	}

	return nil
}

func CheckCompareSn(sn string, planId int64) error {
	nSn, _ := strconv.ParseInt(sn, 10, 64)
	var device Device
	err := baseCheck(nSn, planId, &device)
	if err != nil {
		return err
	}

	//todo 上个流程通过才可以

	return nil
}

func CheckBoxSn(sn string, planId int64) error {
	nSn, _ := strconv.ParseInt(sn, 10, 64)
	var device Device
	err := baseCheck(nSn, planId, &device)
	if err != nil {
		return err
	}

	//todo 上个流程通过才可以

	return nil
}

func CheckPackingSn(sn string, planId int64) error {
	nSn, _ := strconv.ParseInt(sn, 10, 64)
	var device Device
	err := baseCheck(nSn, planId, &device)
	if err != nil {
		return err
	}

	if device.BoxStatus == "1" {
		return fmt.Errorf("该SN已装箱")
	}
	
	return nil
}

func WriteStatus(sn string) error {
	if MysqlConn == nil {
		log.Error("mysql conn invalid")
		return fmt.Errorf("数据库连接失败")
	}

	nSn, _ := strconv.ParseInt(sn, 10, 64)

	result := MysqlConn.Table("device").Where("devno = ?", nSn).
		Update("write_status", "1")

	if result.Error != nil {
		log.Errorf("更新SN状态失败: %v", result.Error)
		return fmt.Errorf("更新SN状态失败: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Errorf("未找到SN为 %s 的记录", sn)
		return fmt.Errorf("未找到SN为 %s 的记录", sn)
	}

	return nil
}

func CompareStatus(sn string) error {
	if MysqlConn == nil {
		log.Error("mysql conn invalid")
		return fmt.Errorf("数据库连接失败")
	}

	nSn, _ := strconv.ParseInt(sn, 10, 64)

	result := MysqlConn.Table("device").Where("devno = ?", nSn).
		Update("compare_status", "1")

	if result.Error != nil {
		log.Errorf("更新SN状态失败: %v", result.Error)
		return fmt.Errorf("更新SN状态失败: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Errorf("未找到SN为 %s 的记录", sn)
		return fmt.Errorf("未找到SN为 %s 的记录", sn)
	}

	return nil
}

func BoxStatus(sn string) error {
	if MysqlConn == nil {
		log.Error("mysql conn invalid")
		return fmt.Errorf("数据库连接失败")
	}

	nSn, _ := strconv.ParseInt(sn, 10, 64)

	result := MysqlConn.Table("device").Where("devno = ?", nSn).
		Update("box_status", "1")

	if result.Error != nil {
		log.Errorf("更新SN状态失败: %v", result.Error)
		return fmt.Errorf("更新SN状态失败: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Errorf("未找到SN为 %s 的记录", sn)
		return fmt.Errorf("未找到SN为 %s 的记录", sn)
	}

	return nil
}

func InsertCompareFailedRecord(sn, scanSn, imei, scanImei string) error {
	if MysqlConn == nil {
		log.Error("mysql conn invalid")
		return fmt.Errorf("数据库连接失败")
	}

	nSn, _ := strconv.ParseInt(sn, 10, 64)
	nScanSn, _ := strconv.ParseInt(scanSn, 10, 64)
	nImei, _ := strconv.ParseInt(imei, 10, 64)
	nScanImei, _ := strconv.ParseInt(scanImei, 10, 64)

	var record CompareSnRecord
	record.Imei = uint(nImei)
	record.Sn = uint(nSn)
	record.ScanSn = uint(nScanSn)
	record.ScanImei = uint(nScanImei)
	record.PlanId = uint(network.CurrentPlan.Id)
	record.ReadStatus = "0"
	record.Operator = network.Username
	record.CreateTime = time.Now()

	result := MysqlConn.Create(&record)
	if result.Error != nil {
		log.Errorf("insert failed:%v", result.Error)
	} else {
		log.Info("insert success")
	}

	return nil
}

func InsertBoxRecord(box Box, snList []string) {
	if MysqlConn == nil {
		log.Error("mysql conn invalid")
		return
	}
	result := MysqlConn.Create(&box)
	if result.Error != nil {
		log.Errorf("insert failed:%v", result.Error)
	} else {
		log.Info("insert success")
	}

	result = MysqlConn.Table("device").Where("devno in ?", snList).UpdateColumns(map[string]interface{}{
		"box_no":         box.BoxNo,
		"packing_status": "1",
	})
	if result.Error != nil {
		log.Errorf("update failed:%v", result.Error)
	} else {
		log.Info("update success")
	}
}
