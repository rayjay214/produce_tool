package db

import (
	"encoding/csv"
	"os"
)

var TestRstCsv *os.File

func init() {
	LoadCsv()
}

func initNewCsv() (*os.File, error) {
	file, err := os.Create("测试记录.csv")
	if err != nil {
		return nil, err
	}
	writer := csv.NewWriter(file)
	headers := []string{
		"是否通过",
		"MES",
		"版本号",
		"SIM卡",
		"IMEI",
		"SN",
		"信号值",
		"卫星",
		"重力",
		"WIFI",
		"光感",
		"IP地址",
		"副IP地址",
		"设置型号",
		"创建时间",
	}
	writer.Write(headers)
	writer.Flush()

	return file, nil
}

func LoadCsv() {
	file, err := os.OpenFile("测试记录.csv", os.O_APPEND, 0644)
	if err != nil {
		newFile, err := initNewCsv()
		if err != nil {
			TestRstCsv = newFile
		}
		return
	}
	TestRstCsv = file
}

func InsertRecordCsv(record TestRecord) {
	if TestRstCsv == nil {
		LoadCsv()
	}

	data := []string{}
	data = append(data, record.Pass)
	data = append(data, record.Mes)
	data = append(data, record.Version)
	data = append(data, record.Sim)
	data = append(data, record.Imei)
	data = append(data, record.Sn)
	data = append(data, record.Signal)
	data = append(data, record.Gps)
	data = append(data, record.Gsensor)
	data = append(data, record.Wifi)
	data = append(data, record.Light)
	data = append(data, record.MainIp)
	data = append(data, record.ViceIp)
	data = append(data, record.SetType)
	data = append(data, record.CreateTime.Format("2006-01-02 15:04:05"))

	writer := csv.NewWriter(TestRstCsv)
	writer.Write(data)
	writer.Flush()
}
