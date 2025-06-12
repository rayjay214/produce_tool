package db

import "time"

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

type CompareSnRecord struct {
	RecordID   uint      `gorm:"column:record_id;primaryKey;autoIncrement"`
	Sn         uint      `gorm:"column:sn"`
	Imei       uint      `gorm:"column:imei"`
	ScanSn     uint      `gorm:"column:scan_sn"`
	ScanImei   uint      `gorm:"column:scan_imei"`
	Operator   string    `gorm:"column:operator"`
	PlanId     uint      `gorm:"column:plan_id"`
	ReadStatus string    `gorm:"column:read_status"`
	CreateTime time.Time `gorm:"column:create_time;primaryKey"`
}

func (CompareSnRecord) TableName() string {
	return "compare_sn_record"
}

type Box struct {
	BoxNo    string `gorm:"column:box_no;primaryKey"`
	PlanId   int64  `gorm:"column:plan_id"`
	Count    int64  `gorm:"column:count"`
	ItemDesc string `gorm:"column:item_desc"`
	ItemCode string `gorm:"column:item_code"`
	Remark   string `gorm:"column:remark"`
}

func (Box) TableName() string {
	return "box"
}

type Device struct {
	Devno         int64  `json:"devno"`
	PlanId        int64  `json:"planId"`
	WriteStatus   string `json:"writeStatus"`
	CompareStatus string `json:"compareStatus"`
	TestStatus    string `json:"testStatus"`
	BoxStatus     string `json:"boxStatus"`
	PackingStatus string `json:"packingStatus"`
	BoxNo         string `json:"boxNo"`
}

func (Device) TableName() string {
	return "device"
}
