package db

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

var TestRstExcel *excelize.File

func init() {
	//LoadExcel()
}

func initNewExcel() (*excelize.File, error) {
	f := excelize.NewFile()
	index, err := f.NewSheet("测试记录")
	if err != nil {
		return nil, err
	}
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

	currLetter := 'A'
	for _, header := range headers {
		pos := fmt.Sprintf("%c1", currLetter)
		f.SetCellValue("测试记录", pos, header)
		currLetter++
	}
	f.SetActiveSheet(index)
	if err := f.SaveAs("测试记录.xlsx"); err != nil {
		return nil, err
	}
	return f, nil
}

func LoadExcel() {
	f, err := excelize.OpenFile("测试记录.xlsx")
	if err != nil {
		newFile, err := initNewExcel()
		if err != nil {
			TestRstExcel = newFile
		}
		return
	}
	TestRstExcel = f
}

func InsertRecordExcel(record TestRecord) {
	if TestRstExcel == nil {
		LoadExcel()
	}
	rows, err := TestRstExcel.GetRows("测试记录")
	if err != nil {
		log.Errorf("get sheet failed %v", err)
		return
	}
	data := []string{}
	/*
		rv := reflect.ValueOf(record)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		for i := 0; i < rv.NumField(); i++ {
			data = append(data, rv.Field(i).String())
		}
	*/
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

	rowIndex := len(rows) + 1
	for colIndex, cellValue := range data {
		cellName, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex)
		TestRstExcel.SetCellValue("测试记录", cellName, cellValue)
	}
	if err := TestRstExcel.Save(); err != nil {
		log.Errorf("save excel failed %v", err)
		return
	}
}
