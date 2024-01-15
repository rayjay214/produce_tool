package util

import (
	"fmt"
	"produce_tool/conf"
	"produce_tool/db"
	"produce_tool/mes"
	"produce_tool/model"
	"reflect"
	"strings"
	"time"
)

type TestItem struct {
	Desc         string
	ModelColName string //对应tableview中model的属性名
	AtCmd        string
	ReturnKey    []string
	ShowKey      string
	Timeout      int //超时时间毫秒
	IsShow       bool
}

type PassParam struct {
	str        string
	stopReader bool
	stopWriter bool
}

var SelectedDeviceType model.DeviceTypeInfo

var allTestItems []TestItem
var allModifyDeviceItems []TestItem

// 用于sn对比工具
var compareSnTestItems []TestItem

// 用于写号工具
var readSnTestItems []TestItem

var CurrTestItems []TestItem

var LastPassParam map[string]*PassParam

var LastGsensorPassParam map[string]*PassParam
var LastLightPassParam map[string]*PassParam

var CompareVersion string
var CompareMainIp string
var CompareViceIp string

var PoweroffAfterTest bool

func init() {
	allTestItems = []TestItem{
		{"开启回显", "Back", "ATE1\r\n", []string{"OK"}, "", 2000, false},
		{"版本号", "Version", "AT+ATI\r\n", []string{"OK"}, "ATI\r\n", 2000, true},
		{"SIM卡", "Sim", "AT+CCID\r\n", []string{"OK", "ERROR"}, "AT+CCID", 2000, true},
		{"IMEI", "Imei", "AT+IMEI\r\n", []string{"OK", "ERROR"}, "AT+IMEI", 2000, true},
		{"SN", "Sn", "AT+SN?\r\n", []string{"OK", "ERROR"}, "SN:", 2000, true},
		{"信号", "Signal", "AT+CSQ\r\n", []string{"OK", "ERROR"}, "+CSQ:", 2000, true},
		{"卫星", "Gps", "AT+GPS\r\n", []string{"OK"}, "+GPS", 2000, true},
		{"重力", "Gsensor", "AT+GS\r\n", []string{"move", "still"}, "AT+GS", 2000, true},
		{"WIFI", "Wifi", "AT+WF\r\n", []string{"OK", "ERROR"}, "wifi test:", 6000, true},
		{"光感", "Light", "AT+PX\r\n", []string{"put on", "fall off"}, "AT+PX", 1000, true},
		{"IP地址", "MainIp", "AT+IP?\r\n", []string{"OK", "ERROR"}, "IP:", 1000, true},
		{"副IP地址", "ViceIp", "AT+IP2?\r\n", []string{"OK", "ERROR"}, "IP2:", 1000, true},
		{"设置型号", "SetType", "AT+SET=\r\n", []string{"OK", "ERROR"}, "\n\rat+set=", 2000, true},
	}

	allModifyDeviceItems = []TestItem{
		{"IMEI", "Imei", "AT+WIMEI=%v\r\n", []string{"OK", "ERROR"}, "AT+IMEI", 1000, true},
		{"SN", "Sn", "AT+SN=%v\r\n", []string{"OK", "ERROR"}, "SN:", 1000, true},
		{"设置型号", "SetType", "AT+SET=%v\r\n", []string{"OK", "ERROR"}, "\n\rat+set=", 1000, true},
		{"设置副IP", "SetViceIp", "AT+IP2=%v#%v#\r\n", []string{"OK", "ERROR"}, "IP2=", 1000, false},
	}

	compareSnTestItems = []TestItem{
		{"开启回显", "Back", "ATE1\r\n", []string{"OK"}, "", 200, false},
		{"IMEI", "Imei", "AT+IMEI\r\n", []string{"OK", "ERROR"}, "AT+IMEI", 200, true},
		{"SN", "Sn", "AT+SN?\r\n", []string{"OK", "ERROR"}, "SN:", 200, true},
	}

	readSnTestItems = []TestItem{
		//{"开启回显", "Back", "ATE1\r\n", []string{"OK"}, "", 2000, false},
		{"SN", "Sn", "AT+SN?\r\n", []string{"OK", "ERROR"}, "SN:", 2000, true},
	}

	OpenAllPorts()
	LastPassParam = make(map[string]*PassParam, 0)
	LastGsensorPassParam = make(map[string]*PassParam, 0)
	LastLightPassParam = make(map[string]*PassParam, 0)
}

func GetModifyDeviceItem(colDesc string) *TestItem {
	for _, item := range allModifyDeviceItems {
		if item.ModelColName == colDesc {
			return &item
		}
	}
	return nil
}

func GetTestItems() []TestItem {
	return CurrTestItems
}

func GetAllTestItems() []TestItem {
	return allTestItems
}

func GetCompareSnTestItems() []TestItem {
	return compareSnTestItems
}

func GetReadSnTestItems() []TestItem {
	return readSnTestItems
}

func GetTestItem(desc string) []TestItem {
	//回显是无论如何都要返回的
	items := make([]TestItem, 0)
	items = append(items, CurrTestItems[0])
	for _, item := range CurrTestItems {
		if item.Desc == desc {
			items = append(items, item)
			return items
		}
	}
	return nil
}

func setBit(num, pos int) int {
	mask := 1 << uint(pos)
	return num | mask
}

/*
 * 写入设备功能位序说明
 * ----------------------------------------------------------------------------------------------------------------
 * |   0     |   1    |   2    |   3    |   4    |   5    |   6    |   7    |   8    |   9    |   10    |   11    |
 * ----------------------------------------------------------------------------------------------------------------
 * | 周期定位 | 监听   |短信设置  |防拆报警 |震动报警 |低电报警 |超速报警 |支持灯控  | ACC    | 继电器  |   录音  | 单片机   |
 * ----------------------------------------------------------------------------------------------------------------
 *
 * at+set=功能,版本,IP:端口
 * eg：at+set=00000000,SK,192.168.1.1:9000
 *
 */
func GetWriteTypeParam() string {
	devFunc := 0
	if SelectedDeviceType.Listen > 0 {
		devFunc = setBit(devFunc, 1)
	}
	if SelectedDeviceType.Sms > 0 {
		devFunc = setBit(devFunc, 2)
	}
	if SelectedDeviceType.TamperAlarm > 0 {
		devFunc = setBit(devFunc, 3)
	}
	if SelectedDeviceType.ShakeAlarm > 0 {
		devFunc = setBit(devFunc, 4)
	}
	if SelectedDeviceType.LowpowerAlarm > 0 {
		devFunc = setBit(devFunc, 5)
	}
	if SelectedDeviceType.OverSpeedAlarm > 0 {
		devFunc = setBit(devFunc, 6)
	}
	if SelectedDeviceType.LightControl > 0 {
		devFunc = setBit(devFunc, 7)
	}
	if SelectedDeviceType.Recording > 0 {
		devFunc = setBit(devFunc, 10)
	}
	if SelectedDeviceType.DeviceType == "" || SelectedDeviceType.MainIp == "" || SelectedDeviceType.MainPort == "" {
		return ""
	}

	param := fmt.Sprintf("%08x,%s,%s:%s", devFunc, SelectedDeviceType.DeviceType, SelectedDeviceType.MainIp, SelectedDeviceType.MainPort)
	return param
}

func SyncTestItems() {
	CurrTestItems = make([]TestItem, 0)
	for _, item := range allTestItems {
		switch item.ModelColName {
		case "Back":
			CurrTestItems = append(CurrTestItems, item)
		case "Version":
			CurrTestItems = append(CurrTestItems, item)
		case "Sim":
			if SelectedDeviceType.SimOpen > 0 {
				CurrTestItems = append(CurrTestItems, item)
			}
		case "Imei":
			if SelectedDeviceType.ImeiOpen > 0 {
				CurrTestItems = append(CurrTestItems, item)
			}
		case "Sn":
			if SelectedDeviceType.SnOpen > 0 {
				CurrTestItems = append(CurrTestItems, item)
			}
		case "Signal":
			if SelectedDeviceType.SignalOpen > 0 {
				CurrTestItems = append(CurrTestItems, item)
			}
		case "Gps":
			if SelectedDeviceType.GpsOpen > 0 {
				CurrTestItems = append(CurrTestItems, item)
			}
		case "Gsensor":
			if SelectedDeviceType.GsensorOpen > 0 {
				CurrTestItems = append(CurrTestItems, item)
			}
		case "Wifi":
			if SelectedDeviceType.WifiOpen > 0 {
				CurrTestItems = append(CurrTestItems, item)
			}
		case "Light":
			if SelectedDeviceType.LightOpen > 0 {
				CurrTestItems = append(CurrTestItems, item)
			}
		case "SetType":
			if SelectedDeviceType.SetTypeOpen > 0 {
				CurrTestItems = append(CurrTestItems, item)
			}
		case "MainIp":
			if SelectedDeviceType.MainIpReadOpen > 0 {
				CurrTestItems = append(CurrTestItems, item)
			}
		case "ViceIp":
			if SelectedDeviceType.ViceIpReadOpen > 0 {
				CurrTestItems = append(CurrTestItems, item)
			}
		}
	}
}

func DoFinish(myport *MyPort, item *MyTableRow) {
	rvTableRow := reflect.ValueOf(item)
	if rvTableRow.Kind() == reflect.Ptr {
		rvTableRow = rvTableRow.Elem()
	}
	rtTableRow := rvTableRow.Type()

	rvDeviceType := reflect.ValueOf(SelectedDeviceType)
	if rvDeviceType.Kind() == reflect.Ptr {
		rvDeviceType = rvDeviceType.Elem()
	}
	//rtDeviceType := rvDeviceType.Type()

	bPass := true
	for i := 0; i < rvTableRow.NumField(); i++ {
		fieldValueTableRow := rvTableRow.Field(i)
		fieldTypeTableRow := rtTableRow.Field(i)

		if fieldTypeTableRow.Name == "Pass" || fieldTypeTableRow.Name == "Com" {
			continue
		}

		if fieldTypeTableRow.Name == "Mes" {
			continue
		}

		if fieldTypeTableRow.Name == "Version" {
			if ContainsOne(fieldValueTableRow.String(), "失败", "超时", "等待") {
				bPass = false
				break
			}
			continue
		}

		switchName := fieldTypeTableRow.Name + "Open"
		if fieldTypeTableRow.Name == "MainIp" || fieldTypeTableRow.Name == "ViceIp" {
			switchName = fieldTypeTableRow.Name + "ReadOpen"
		}

		if rvDeviceType.FieldByName(switchName).Int() <= 0 {
			continue
		}

		if ContainsOne(fieldValueTableRow.String(), "失败", "超时", "等待") {
			bPass = false
		}

		if fieldValueTableRow.String() == "" {
			bPass = false
		}
	}

	if bPass && item.Sn != "" && item.Sn != "13100018888" {
		conf.CntMutex.Lock()
		conf.PassedCnt += 1
		conf.CntMutex.Unlock()

		//go SaveResultToMysql(*item, "通过")
		//go SaveResultToExcel(*item, "通过")
		//go SaveResultToCsv(*item, "通过")
		go SaveResultToMes(item, myport)
	}

	if bPass && item.Sn != "13100018888" && CompareVersion != "" && PoweroffAfterTest {
		fmt.Printf("begin to close %v\n", item.Sn)
		DoCloseDevice(myport, item.Sn)
	}

}

// todo use reflect
func makeRecord(item MyTableRow, result string, mes string) db.TestRecord {
	record := db.TestRecord{}
	record.Pass = result
	record.Mes = mes
	record.Version = item.Version
	record.Sim = item.Sim
	record.Sn = strings.Trim(strings.Trim(item.Sn, "写入成功("), ")")
	record.Imei = strings.Trim(strings.Trim(item.Imei, "写入成功("), ")")
	record.Signal = item.Signal
	record.Gps = item.Gps
	record.Gsensor = item.Gsensor
	record.Wifi = item.Wifi
	record.MainIp = item.MainIp
	record.ViceIp = item.ViceIp
	record.SetType = item.SetType
	record.CreateTime = time.Now()
	return record
}

func SaveResultToMysql(item MyTableRow, result string, mes string) {
	record := makeRecord(item, result, mes)
	db.InsertRecordMysql(record)
}

func SaveResultToExcel(item MyTableRow, result string, mes string) {
	record := makeRecord(item, result, mes)
	db.InsertRecordExcel(record)
}

func SaveResultToCsv(item MyTableRow, result string, mes string) {
	record := makeRecord(item, result, mes)
	db.InsertRecordCsv(record)
}

func SaveResultToMes(item *MyTableRow, myport *MyPort) {
	detail := ""
	rv := reflect.ValueOf(item)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		fieldValue := rv.Field(i)
		fieldType := rt.Field(i)
		if fieldType.Name == "Com" || fieldType.Name == "Pass" {
			continue
		}
		s := fmt.Sprintf("%v|%v;", fieldType.Name, fieldValue.String())
		detail += s
	}

	result := mes.SetMesReq(item.Sn, detail, "DIGNWEIQICESHI")
	if result {
		item.Mes = "成功"
		go SaveResultToCsv(*item, "通过", "通过")
		go SaveResultToMysql(*item, "通过", "通过")
	} else {
		item.Mes = "失败"
		go SaveResultToCsv(*item, "通过", "失败")
		go SaveResultToMysql(*item, "通过", "失败")
	}
	model := GetTableModel()
	model.PublishRowChanged(PortNameRowidx[myport.Name])
}
