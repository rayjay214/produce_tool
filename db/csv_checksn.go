package db

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

type CheckSnRecord struct {
	RST        string
	Imei       string
	Sn         string
	CreateTime time.Time
}

var CheckSnCsv *os.File

func initCheckSnCsv() (*os.File, error) {
	file, err := os.Create("查号记录.csv")
	if err != nil {
		return nil, err
	}
	writer := csv.NewWriter(file)
	headers := []string{
		"测试结果",
		"IMEI",
		"SN",
		"创建时间",
	}
	writer.Write(headers)
	writer.Flush()

	return file, nil
}

func LoadCheckSnCsv() {
	file, err := os.OpenFile("查号记录.csv", os.O_APPEND, 0644)
	if err != nil {
		newFile, err := initCheckSnCsv()
		if err != nil {
			CheckSnCsv = newFile
		}
		return
	}
	CheckSnCsv = file
}

func InsertCheckSnCsv(record CheckSnRecord) {
	if CheckSnCsv == nil {
		LoadCheckSnCsv()
	}

	data := []string{}
	data = append(data, record.RST)
	data = append(data, record.Imei)
	data = append(data, record.Sn)
	data = append(data, record.CreateTime.Format("2006-01-02 15:04:05"))

	writer := csv.NewWriter(CheckSnCsv)
	writer.Write(data)
	writer.Flush()
}

func WriteCheckSnLog(record CheckSnRecord) {
	file, err := os.Create(fmt.Sprintf("log/%v.txt", record.Sn))
	if err != nil {
		return
	}
	defer file.Close()

	data := []byte(fmt.Sprintf("%v,%v,%v", record.Sn, record.Imei, record.CreateTime))
	_, err = file.Write(data)
	if err != nil {
		return
	}
}
