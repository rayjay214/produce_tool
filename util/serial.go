package util

import (
	"fmt"
	"github.com/lxn/walk"
	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
	bs "go.bug.st/serial"
	"produce_tool/db"
	"produce_tool/mes"
	//"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SortByName []*MyPort

func (a SortByName) Len() int           { return len(a) }
func (a SortByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

const (
	Init = iota
	Testing
	Waiting
	Passed
	Failed
)

type MyPort struct {
	Name  string
	Port  *serial.Port
	Vaild bool
}

var PortList []*MyPort
var WholePortList []*MyPort

func GetPorts() []*MyPort {
	return PortList
}

func GetPort(portName string) *MyPort {
	for _, port := range PortList {
		if port.Name == portName {
			return port
		}
	}
	return nil
}

func CheckPorts() {
	for _, myport := range PortList {
		c := &serial.Config{Name: myport.Name,
			Baud:        115200,
			ReadTimeout: time.Millisecond * 500,
			Size:        8,
			Parity:      'N',
			StopBits:    0,
		}

		myport.Port.Close()
		time.Sleep(10 * time.Millisecond)

		s, err := serial.OpenPort(c)
		if err != nil {
			fmt.Printf("open port err %v\n", err)
			log.Errorf("open port err %v", err)
			myport.Vaild = false
			continue
		}
		myport.Vaild = true
		myport.Port = s
	}
}

func OpenAllPorts() {
	ports, err := bs.GetPortsList()
	if err != nil {
		fmt.Println("Error getting serial ports:", err)
		log.Errorf("Error getting serial ports:%v", err)
		return
	}

	for _, port := range ports {
		c := &serial.Config{Name: port,
			Baud:        115200,
			ReadTimeout: time.Millisecond * 500,
			Size:        8,
			Parity:      'N',
			StopBits:    0,
		}
		s, err := serial.OpenPort(c)
		if err != nil {
			fmt.Printf("open port err %v\n", err)
			log.Errorf("open port err %v", err)
			continue
		}
		myPort := new(MyPort)
		myPort.Name = port
		myPort.Port = s
		myPort.Vaild = true

		WholePortList = append(WholePortList, myPort)
	}

	sort.Sort(SortByName(WholePortList))
	for _, myport := range WholePortList {
		PortList = append(PortList, myport)
	}
}

func RefreshPorts(blockedPortList []string) {
	PortList = nil
	for _, myport := range WholePortList {
		bAdd := true
		for _, blockedPort := range blockedPortList {
			if blockedPort == myport.Name {
				bAdd = false
				break
			}
		}
		if bAdd {
			PortList = append(PortList, myport)
		}
	}
}

func setDevType(myport *MyPort, pass *PassParam) {
	param := GetWriteTypeParam()
	fmt.Printf("param is %v\n", param)
	modifyDevice(myport, pass, "SetType", param, false)
}

func setViceIp(myport *MyPort, pass *PassParam) {
	time.Sleep(1 * time.Second)
	value := fmt.Sprintf("%v#%v#", SelectedDeviceType.ViceIp, SelectedDeviceType.VicePort)
	modifyDevice(myport, pass, "ViceIpWrite", value, false)
}

func writeItems(myport *MyPort, items []TestItem, pass *PassParam) {
	var wg sync.WaitGroup
	model := GetTableModel()
	tableItem := model.items[PortNameRowidx[myport.Name]]
	bForceStop := false
	for _, item := range items {
		if item.ModelColName == "SetType" {
			setDevType(myport, pass)
			continue
		}

		if item.ModelColName == "ViceIpWrite" {
			setViceIp(myport, pass)
			continue
		}

		b := writeComm(myport, item, pass)
		if pass.stopWriter {
			bForceStop = true
			break
		}

		_, respValue := getValue(pass.str, item.ShowKey)
		var showValue string
		if b && !strings.Contains(pass.str, "ERROR") {
			showValue = respValue
		} else if !b {
			showValue = "获取超时"
		} else {
			showValue = "失败"
		}

		if b && respValue != "" && (item.ModelColName == "Gsensor" || item.ModelColName == "Light") {
			showValue = "等待中"
			wg.Add(1)
			go doubleCheck(&wg, myport, respValue, item)
		}

		if showValue == respValue && item.ModelColName == "Signal" {
			t := strings.Split(showValue, ",")
			csq, _ := strconv.Atoi(t[0])
			nSignalMin, _ := strconv.Atoi(SelectedDeviceType.SignalMin)
			nSignalMax, _ := strconv.Atoi(SelectedDeviceType.SignalMax)
			if csq >= nSignalMin && csq <= nSignalMax {
				showValue = fmt.Sprintf("通过(%v)", csq)
			} else {
				showValue = fmt.Sprintf("失败(%v)", csq)
			}
		}

		if showValue == respValue && item.ModelColName == "Gps" {
			gps, _ := strconv.Atoi(showValue)
			nGpsMin, _ := strconv.Atoi(SelectedDeviceType.GpsMin)
			if gps >= nGpsMin {
				showValue = fmt.Sprintf("通过(%v)", gps)
			} else {
				showValue = fmt.Sprintf("失败(%v)", gps)
			}
		}

		if showValue == respValue && item.ModelColName == "Wifi" {
			wifi, _ := strconv.Atoi(showValue)
			nWifiMin, _ := strconv.Atoi(SelectedDeviceType.WifiMin)
			if wifi >= nWifiMin {
				showValue = fmt.Sprintf("通过(%v)", wifi)
			} else {
				showValue = fmt.Sprintf("失败(%v)", wifi)
			}
		}

		if item.ModelColName == "Version" && CompareVersion != "" {

			if showValue != CompareVersion {
				fmt.Printf("show value %v, compare value %v\n", showValue, CompareVersion)
				showValue = fmt.Sprintf("%v(匹配失败)", showValue)
			}
		}

		if item.ModelColName == "MainIp" && CompareMainIp != "" {
			if showValue != CompareMainIp {
				showValue = fmt.Sprintf("%v(匹配失败)", showValue)
			}
		}

		if item.ModelColName == "ViceIp" && CompareViceIp != "" {
			if showValue != CompareViceIp {
				showValue = fmt.Sprintf("%v(匹配失败)", showValue)
			}
		}

		if item.IsShow {
			v := reflect.ValueOf(tableItem)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			v.FieldByName(item.ModelColName).SetString(showValue)

			model.PublishRowChanged(PortNameRowidx[myport.Name])
		}
	}
	pass.stopReader = true
	pass.stopWriter = true

	wg.Wait() //等待重力和光感流程结束
	if !bForceStop {
		DoFinish(myport, tableItem)
	}
}

func modifyDevice(myport *MyPort, pass *PassParam, colDesc string, writeValue string, stopWhenFinish bool) {
	model := GetTableModel()
	tableItem := model.items[PortNameRowidx[myport.Name]]
	modifyDeviceItem := GetModifyDeviceItem(colDesc)
	if modifyDeviceItem == nil {
		return
	}

	writeSuccess := false
	rstSuccess := false
	for i := 0; i < 1; i++ {
		strCmd := fmt.Sprintf(modifyDeviceItem.AtCmd, writeValue)
		fmt.Printf("write %v\n", strCmd)
		_, err := myport.Port.Write([]byte(strCmd))
		if err != nil {
			writeSuccess = false
			break
		}
		time.Sleep(10 * time.Millisecond)

		//等待设备返回结果
		timeout := time.Duration(modifyDeviceItem.Timeout) * time.Millisecond
		startTime := time.Now()
		for {
			if time.Since(startTime) >= timeout {
				writeSuccess = false
				break
			}
			time.Sleep(10 * time.Millisecond)
			if strings.Contains(pass.str, "OK") || strings.Contains(pass.str, "ok") {
				writeSuccess = true
				rstSuccess = true
				break
			}
			if strings.Contains(pass.str, "ERROR") || strings.Contains(pass.str, "error") {
				writeSuccess = true
				rstSuccess = false
				break
			}
		}
	}
	var showValue string
	if writeSuccess && rstSuccess {
		showValue = fmt.Sprintf("写入成功(%s)", writeValue)
	} else if writeSuccess && !rstSuccess {
		showValue = "写入失败"
	} else {
		showValue = "超时"
	}

	v := reflect.ValueOf(tableItem)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	v.FieldByName(modifyDeviceItem.ModelColName).SetString(showValue)
	model.PublishRowChanged(PortNameRowidx[myport.Name])

	if stopWhenFinish {
		pass.stopReader = true
		pass.stopWriter = true
	}
}

func closeDevice(myport *MyPort, pass *PassParam, sn string) {
	bSuccess := true
	for i := 0; i < 2; i++ {
		strCmd := "AT+POWEROFF\r\n"
		_, err := myport.Port.Write([]byte(strCmd))
		if err != nil {
			break
		}
		time.Sleep(10 * time.Millisecond)

		//等待设备返回结果
		timeout := time.Duration(1000) * time.Millisecond
		startTime := time.Now()
		for {
			if time.Since(startTime) >= timeout {
				bSuccess = false
				break
			}
			time.Sleep(10 * time.Millisecond)
			if strings.Contains(pass.str, "OK") || strings.Contains(pass.str, "ok") {
				break
			}
			if strings.Contains(pass.str, "ERROR") || strings.Contains(pass.str, "error") {
				bSuccess = false
				break
			}
		}
	}

	fmt.Printf("close device %v rst %v\n", sn, bSuccess)

	pass.stopReader = true
	pass.stopWriter = true
}

func readPort(myPort *MyPort, pass *PassParam) {
	buf := make([]byte, 128)
	for {
		if pass.stopReader {
			log.Infoln("reader stop")
			break
		}
		if !myPort.Vaild {
			log.Errorf("port %v invalid", myPort.Name)
			return
		}
		n, err := myPort.Port.Read(buf)
		if err != nil {
			log.Errorf("read err %v, port %v", err, myPort)
			myPort.Vaild = false
			return
		}
		if n > 0 {
			data := buf[:n]
			pass.str += string(data)
		}
	}
}

func DoCloseDevice(myPort *MyPort, Sn string) {
	pass := new(PassParam)
	go readPort(myPort, pass)
	go closeDevice(myPort, pass, Sn)
}

func DoTestAllPortsAllItems() {
	model := GetTableModel()
	model.ResetRows()

	items := GetTestItems()
	myPorts := GetPorts()

	log.Infof("begin onekey test, ports %v, item %v", myPorts, items)

	for _, myPort := range myPorts {
		if !myPort.Vaild {
			log.Errorf("port %v valid", myPort.Name)
			continue
		}

		lastPass, ok := LastPassParam[myPort.Name]
		if ok {
			if lastPass.stopWriter == false {
				lastPass.stopWriter = true
				time.Sleep(100 * time.Millisecond)
			}
		}
		lastGsensorPass, ok := LastGsensorPassParam[myPort.Name]
		if ok {
			if lastGsensorPass.stopWriter == false {
				lastGsensorPass.stopWriter = true
				time.Sleep(100 * time.Millisecond)
			}
		}
		lastLightPass, ok := LastLightPassParam[myPort.Name]
		if ok {
			if lastLightPass.stopWriter == false {
				lastLightPass.stopWriter = true
				time.Sleep(100 * time.Millisecond)
			}
		}

		pass := new(PassParam)
		LastPassParam[myPort.Name] = pass

		go readPort(myPort, pass)
		go writeItems(myPort, items, pass)
	}
}

func DoTestOnePortAllItems(portName string, idx int) {
	model := GetTableModel()
	model.ClearRow(idx)

	items := GetTestItems()
	myPort := GetPort(portName)
	if myPort.Name == portName {
		lastPass, ok := LastPassParam[myPort.Name]
		if ok {
			if lastPass.stopWriter == false {
				lastPass.stopWriter = true
				time.Sleep(100 * time.Millisecond)
			}
		}
		pass := new(PassParam)
		go readPort(myPort, pass)
		go writeItems(myPort, items, pass)
		LastPassParam[myPort.Name] = pass
	}
}

func DoOnePortWriteSn(portName string, sn string) {
	myPort := GetPort(portName)
	//bSuccess := checkMesSn(sn)
	bSuccess := true
	if bSuccess {
		if myPort.Name == portName {
			pass := new(PassParam)
			go readPort(myPort, pass)
			go modifyDevice(myPort, pass, "Sn", sn, true)
		}
	} else {
		model := GetTableModel()
		tableItem := model.items[PortNameRowidx[portName]]
		tableItem.Sn = "已过站"
		model.PublishRowChanged(PortNameRowidx[portName])
	}
}

func DoOnePortWriteImei(portName string, imei string) {
	myPort := GetPort(portName)
	if myPort.Name == portName {
		pass := new(PassParam)
		go readPort(myPort, pass)
		go modifyDevice(myPort, pass, "Imei", imei, true)
	}
}

func DoTestAllPortsOneItem(itemDesc string, label *walk.Label) {
	model := GetTableModel()
	model.ResetRows()

	items := GetTestItem(itemDesc)
	myPorts := GetPorts()

	for _, myPort := range myPorts {
		pass := new(PassParam)
		go readPort(myPort, pass)
		go writeItems(myPort, items, pass)
	}
}

func writeComm(myport *MyPort, item TestItem, pass *PassParam) bool {
	pass.str = ""
	for i := 0; i < 2; i++ {
		if pass.stopWriter {
			break
		}
		log.Infof("write to port %v, %v", myport.Name, string(item.AtCmd))
		_, err := myport.Port.Write([]byte(item.AtCmd))
		if err != nil {
			log.Errorf("err is %v, port is %v", err, myport.Name)
			myport.Vaild = false
			return false
		}
		time.Sleep(10 * time.Millisecond)

		//等待设备返回结果
		timeout := time.Duration(item.Timeout) * time.Millisecond
		startTime := time.Now()
		for {
			if time.Since(startTime) >= timeout {
				log.Println("Timeout!")
				break
			}
			time.Sleep(10 * time.Millisecond)
			for _, retkey := range item.ReturnKey {
				contains := strings.Contains(pass.str, retkey)
				if contains {
					log.Infof("get response %v", pass.str)
					return true
				}
			}
		}
	}
	return false
}

func doubleCheck(wg *sync.WaitGroup, myPort *MyPort, lastValue string, item TestItem) {
	defer wg.Done()
	if SelectedDeviceType.WifiOpen > 0 { //2G的wifi测试，需要3s+才能返回结果，避免跟重力同时测试发生冲突
		time.Sleep(4000 * time.Millisecond)
	} else {
		time.Sleep(1000 * time.Millisecond)
	}
	model := GetTableModel()
	tableItem := model.items[PortNameRowidx[myPort.Name]]
	pass := new(PassParam)
	showValue := "失败"
	if item.ModelColName == "Gsensor" {
		LastGsensorPassParam[myPort.Name] = pass
	}
	if item.ModelColName == "Light" {
		LastLightPassParam[myPort.Name] = pass
	}

	go readPort(myPort, pass)
	for i := 0; i < 45; i++ {
		if pass.stopWriter {
			break
		}
		b := writeComm(myPort, item, pass)
		_, respValue := getValue(pass.str, item.ShowKey)
		fmt.Printf("respValue:%v, lastValue:%v\n", respValue, lastValue)
		if b && respValue != "" && respValue != lastValue {
			showValue = "通过"
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
	pass.stopReader = true
	v := reflect.ValueOf(tableItem)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	v.FieldByName(item.ModelColName).SetString(showValue)
	model.PublishRowChanged(PortNameRowidx[myPort.Name])
}

func checkMesSn(sn string) bool {
	return mes.CheckMesReq(sn, "DIGNWEIQICESHI")
}

func readSnImei(myport *MyPort, items []TestItem, pass *PassParam) (string, string) {
	var sn, imei string
	for _, item := range items {
		b := writeComm(myport, item, pass)
		_, respValue := getValue(pass.str, item.ShowKey)
		var showValue string
		if b && !strings.Contains(pass.str, "ERROR") {
			showValue = respValue
		} else if !b {
			showValue = "获取超时"
		} else {
			showValue = "失败"
		}

		if item.Desc == "IMEI" {
			imei = showValue
		}

		if item.Desc == "SN" {
			sn = showValue
		}

	}
	pass.stopReader = true
	pass.stopWriter = true

	return sn, imei
}

// 用于SN比较工具
func DoTestOnePortCompareSn(portName string, scanSnEdit *walk.LineEdit, prefix string, readSn *walk.LineEdit, readImei *walk.LineEdit, resultEdit *walk.LineEdit) {
	var sn, imei string

	items := GetCompareSnTestItems()
	myPort := GetPort(portName)
	if myPort.Name == portName {
		pass := new(PassParam)
		go readPort(myPort, pass)
		sn, imei = readSnImei(myPort, items, pass)
	}
	readSn.SetText(sn)
	readImei.SetText(imei)

	scanSn := scanSnEdit.Text()

	record := db.CheckSnRecord{
		RST:        "通过",
		Imei:       imei,
		Sn:         sn,
		CreateTime: time.Now(),
	}

	if sn == scanSn && imei == (prefix+scanSn) {
		//匹配通过，去上报MES
		result := mes.CheckMesReq(scanSn, "CHAHAO")
		if result {
			result = mes.SetMesReq(scanSn, "", "CHAHAO")
			if result {
				brush, _ := walk.NewSolidColorBrush(walk.RGB(0, 255, 0))
				resultEdit.SetBackground(brush)
				resultEdit.SetText("PASS")
				scanSnEdit.SetText("")
				db.InsertCheckSnCsv(record)
			} else {
				brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
				resultEdit.SetBackground(brush)
				resultEdit.SetText("MES过站失败")
			}
		} else {
			brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
			resultEdit.SetBackground(brush)
			resultEdit.SetText("MES查号失败")
		}
	} else if sn != scanSn {
		brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
		resultEdit.SetBackground(brush)
		resultEdit.SetText("比对失败")
	} else if imei != (prefix + scanSn) {
		brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
		resultEdit.SetBackground(brush)
		resultEdit.SetText("IMEI前缀错误")
	} else {
		brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
		resultEdit.SetBackground(brush)
		resultEdit.SetText("FAIL")
	}

}

func readCommSn(myport *MyPort, items []TestItem, pass *PassParam) string {
	var sn string
	for _, item := range items {
		b := writeComm(myport, item, pass)
		_, respValue := getValue(pass.str, item.ShowKey)
		var showValue string
		if b && !strings.Contains(pass.str, "ERROR") {
			showValue = respValue
		} else if !b {
			showValue = "获取超时"
		} else {
			showValue = "失败"
		}

		if item.Desc == "SN" {
			sn = showValue
		}
	}
	pass.stopReader = true
	pass.stopWriter = true

	return sn
}

func writeCommSn(myport *MyPort, pass *PassParam, writeValue string) string {
	modifyDeviceItem := GetModifyDeviceItem("Sn")
	if modifyDeviceItem == nil {
		return ""
	}

	writeSuccess := false
	rstSuccess := false
	for i := 0; i < 1; i++ {
		strCmd := fmt.Sprintf(modifyDeviceItem.AtCmd, writeValue)
		_, err := myport.Port.Write([]byte(strCmd))
		if err != nil {
			writeSuccess = false
			break
		}
		time.Sleep(100 * time.Millisecond)

		//等待设备返回结果
		timeout := time.Duration(modifyDeviceItem.Timeout) * time.Millisecond
		startTime := time.Now()
		for {
			if time.Since(startTime) >= timeout {
				writeSuccess = false
				break
			}
			time.Sleep(10 * time.Millisecond)
			if strings.Contains(pass.str, "OK") || strings.Contains(pass.str, "ok") {
				writeSuccess = true
				rstSuccess = true
				break
			}
			if strings.Contains(pass.str, "ERROR") || strings.Contains(pass.str, "error") {
				writeSuccess = true
				rstSuccess = false
				break
			}
		}
	}

	var showValue string
	if writeSuccess && rstSuccess {
		showValue = fmt.Sprintf("写入成功(%s)", writeValue)
	} else if writeSuccess && !rstSuccess {
		showValue = "写入失败"
	} else {
		showValue = "超时"
	}

	pass.stopReader = true
	pass.stopWriter = true

	return showValue
}

// 用于写号工具
func DoTestOnePortWriteSn(portName string, SnValue string, readSn *walk.LineEdit, resultEdit *walk.LineEdit, scanSn *walk.LineEdit) {
	myPort := GetPort(portName)
	writeRst := ""
	if myPort.Name == portName {
		pass := new(PassParam)
		go readPort(myPort, pass)
		writeRst = writeCommSn(myPort, pass, SnValue)
	}

	time.Sleep(1000 * time.Millisecond)
	var sn string
	items := GetReadSnTestItems()

	if myPort.Name == portName {
		pass := new(PassParam)
		go readPort(myPort, pass)
		sn = readCommSn(myPort, items, pass)
	}
	readSn.SetText(sn)

	log.Infof("rayjay rst:%v, sn:%v", writeRst, sn)

	if strings.Contains(writeRst, "成功") && sn == SnValue {
		brush, _ := walk.NewSolidColorBrush(walk.RGB(0, 255, 0))
		resultEdit.SetBackground(brush)
		resultEdit.SetText("SN写入成功")
		scanSn.SetText("")
	} else {
		brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
		resultEdit.SetBackground(brush)
		resultEdit.SetText("SN写入失败")
	}

}
