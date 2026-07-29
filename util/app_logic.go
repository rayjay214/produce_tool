package util

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/lxn/walk"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"produce_tool/db"
	"produce_tool/network"
	"strconv"
	"strings"
	"sync"
	"time"
)

func readSnImei(myport *MyPort, items []TestItem, pass *PassParam) (string, string, string) {
	var sn, imei, mode string
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

		if item.Desc == "MODE" {
			mode = showValue
		}
	}
	pass.stopReader = true
	pass.stopWriter = true

	return sn, imei, mode
}

func compareSn(p CompareSnParams, sn, imei, scanSn string, record db.CheckSnRecord) {
	if p.OnlyCompareSn {
		if sn == scanSn {
			err := db.CheckCompareSn(sn, network.CurrentPlan.Id)
			if err != nil {
				brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
				p.ResultEdit.SetBackground(brush)
				p.ResultEdit.SetText(err.Error())
			} else {
				if p.SetMode {
					if p.ModeResult != "0" {
						brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
						p.ResultEdit.SetBackground(brush)
						p.ResultEdit.SetText("模式设置失败")
						p.ScanSnEdit.SetText("")
						return
					}
				}
				brush, _ := walk.NewSolidColorBrush(walk.RGB(0, 255, 0))
				p.ResultEdit.SetBackground(brush)
				p.ResultEdit.SetText("PASS")
				p.ScanSnEdit.SetText("")
				db.WriteCheckSnLog(record)
				db.CompareStatus(sn)
			}
		} else if sn != scanSn {
			brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
			p.ResultEdit.SetBackground(brush)
			p.ResultEdit.SetText("比对失败")
			db.InsertCompareFailedRecord(sn, scanSn, "0", "0")
		} else {
			brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
			p.ResultEdit.SetBackground(brush)
			p.ResultEdit.SetText("FAIL")
		}
	} else {
		if sn == scanSn && imei == (p.ImeiPrefix+scanSn) {
			err := db.CheckCompareSn(sn, network.CurrentPlan.Id)
			if err != nil {
				brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
				p.ResultEdit.SetBackground(brush)
				p.ResultEdit.SetText(err.Error())
			} else {
				brush, _ := walk.NewSolidColorBrush(walk.RGB(0, 255, 0))
				p.ResultEdit.SetBackground(brush)
				p.ResultEdit.SetText("PASS")
				p.ScanSnEdit.SetText("")
				db.WriteCheckSnLog(record)
				if network.CurrentPlan.SnType == "1" {
					db.CompareStatus(imei)
				} else {
					db.CompareStatus(sn)
				}
			}
		} else if sn != scanSn {
			brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
			p.ResultEdit.SetBackground(brush)
			p.ResultEdit.SetText("比对失败")
			db.InsertCompareFailedRecord(sn, scanSn, imei, p.ImeiPrefix+scanSn)
		} else if imei != (p.ImeiPrefix + scanSn) {
			brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
			p.ResultEdit.SetBackground(brush)
			p.ResultEdit.SetText("IMEI前缀错误")
			db.InsertCompareFailedRecord(sn, scanSn, imei, p.ImeiPrefix+scanSn)
		} else {
			brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
			p.ResultEdit.SetBackground(brush)
			p.ResultEdit.SetText("FAIL")
		}
	}
}

func compareImei(p CompareSnParams, imei, scanSn string, record db.CheckSnRecord) {
	if imei == scanSn {
		err := db.CheckCompareSn(imei, network.CurrentPlan.Id)
		if err != nil {
			brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
			p.ResultEdit.SetBackground(brush)
			p.ResultEdit.SetText(err.Error())
		} else {
			brush, _ := walk.NewSolidColorBrush(walk.RGB(0, 255, 0))
			p.ResultEdit.SetBackground(brush)
			p.ResultEdit.SetText("PASS")
			p.ScanSnEdit.SetText("")
			db.WriteCheckSnLog(record)
			db.CompareStatus(imei)
		}
	} else if imei != scanSn {
		brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
		p.ResultEdit.SetBackground(brush)
		p.ResultEdit.SetText("比对失败")
		db.InsertCompareFailedRecord(imei, scanSn, imei, p.ImeiPrefix+scanSn)
	} else {
		brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
		p.ResultEdit.SetBackground(brush)
		p.ResultEdit.SetText("FAIL")
	}
}

type CompareSnParams struct {
	PortName      string
	ScanSnEdit    *walk.LineEdit
	ImeiPrefix    string
	ReadSnEdit    *walk.LineEdit
	ReadImeiEdit  *walk.LineEdit
	ResultEdit    *walk.LineEdit
	ModeEdit      *walk.LineEdit
	OnlyCompareSn bool
	SetMode       bool
	ModeResult    string
}

// 用于SN比较工具
func DoTestOnePortCompareSn(p CompareSnParams) {
	var sn, imei, mode string

	var items []TestItem
	if p.SetMode {
		items = GetCompareSnWithModeTestItems()
	} else {
		items = GetCompareSnTestItems()
	}

	myPort := GetPort(p.PortName)
	if myPort.Name == p.PortName {
		pass := new(PassParam)
		go readPort(myPort, pass)
		sn, imei, mode = readSnImei(myPort, items, pass)
	}
	p.ReadSnEdit.SetText(sn)
	p.ReadImeiEdit.SetText(imei)
	p.ModeEdit.SetText(mode)
	p.ModeResult = mode

	scanSn := p.ScanSnEdit.Text()

	record := db.CheckSnRecord{
		RST:        "通过",
		Imei:       imei,
		Sn:         sn,
		CreateTime: time.Now(),
	}

	if network.CurrentType.SnType == "1" {
		compareImei(p, imei, scanSn, record)
	} else {
		compareSn(p, sn, imei, scanSn, record)
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

// 用于写IMEI工具
func writeCommImei(myport *MyPort, pass *PassParam, writeValue string) string {
	modifyDeviceItem := GetModifyDeviceItem("Imei")
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
		//timeout := time.Duration(modifyDeviceItem.Timeout) * time.Millisecond
		timeout := time.Duration(4000) * time.Millisecond
		startTime := time.Now()
		for {
			if time.Since(startTime) >= timeout {
				writeSuccess = false
				break
			}
			time.Sleep(10 * time.Millisecond)
			if strings.Contains(pass.str, "OK") || strings.Contains(pass.str, "ok") {
				//返回两次，第一次表示开始执行，第二次表示执行成功
				for {
					pass.str = ""
					if time.Since(startTime) >= timeout {
						writeSuccess = false
						break
					}
					time.Sleep(10 * time.Millisecond)
					if strings.Contains(pass.str, "OK") && strings.Contains(pass.str, "IMEI:") {
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
				//writeSuccess = true
				//rstSuccess = true
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

func readCommImei(myport *MyPort, items []TestItem, pass *PassParam) string {
	var imei string
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
	}
	pass.stopReader = true
	pass.stopWriter = true

	return imei
}

// 用于写IMEI工具
func DoTestOnePortWriteImei(portName string, ImeiValue string, readImei *walk.LineEdit, resultEdit *walk.LineEdit, scanImei *walk.LineEdit) {
	myPort := GetPort(portName)
	writeRst := ""
	if myPort.Name == portName {
		pass := new(PassParam)
		go readPort(myPort, pass)
		writeRst = writeCommImei(myPort, pass, ImeiValue)
	}

	//time.Sleep(1000 * time.Millisecond)
	var imei string
	items := GetReadImeiTestItems()

	if myPort.Name == portName {
		pass := new(PassParam)
		go readPort(myPort, pass)
		imei = readCommImei(myPort, items, pass)
	}
	readImei.SetText(imei)

	log.Infof("rayjay rst:%v, imei:%v", writeRst, imei)

	if strings.Contains(writeRst, "成功") && imei == ImeiValue {
		brush, _ := walk.NewSolidColorBrush(walk.RGB(0, 255, 0))
		resultEdit.SetBackground(brush)
		resultEdit.SetText("IMEI写入成功")
		scanImei.SetText("")
	} else {
		brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
		resultEdit.SetBackground(brush)
		resultEdit.SetText("IMEI写入失败")
	}
}

func checkRecordTest(mw *walk.MainWindow, myport *MyPort, items []TestItem, pass *PassParam, imei string, resultEdit *walk.LineEdit) string {
	var result string
	item := items[0]
	_, err := myport.Port.Write([]byte(item.AtCmd))
	if err != nil {
		log.Errorf("err is %v, port is %v", err, myport.Name)
		myport.Vaild = false
	}
	wavLength := -1
	cnt := 0
	for {
		if cnt > 100*30 {
			result = "超时"
			break
		}
		if strings.Contains(pass.str, "+CREC:0") {
			if strings.Contains(pass.str, "ERROR") {
				result = "录音下发失败"
				break
			}
			pass.str = ""
			mw.Synchronize(func() {
				resultEdit.SetText("设备录音中")
			})
		}
		if strings.Contains(pass.str, "+CFTRANTX:DATA") {
			strList := strings.Split(pass.str, ",")
			strLength := strings.Trim(strList[1], "\r\n")
			wavLength, _ = strconv.Atoi(strLength)
			fmt.Println(wavLength)
			pass.str = ""
			mw.Synchronize(func() {
				resultEdit.SetText("接收录音中")
			})
		}
		if len(pass.str) == wavLength { //录音文件接收完成
			os.WriteFile(fmt.Sprintf("record/%v.amr", imei), []byte(pass.str), 0644)
			pass.str = ""
			mw.Synchronize(func() {
				resultEdit.SetText("录音接收完成")
			})
		}
		if strings.Contains(pass.str, "+CFTRANTX:0") {
			mw.Synchronize(func() {
				resultEdit.SetText("播放录音中")
			})
			//转码amr->wav
			src := fmt.Sprintf("record/%v.amr", imei)
			dst := fmt.Sprintf("record/%v.wav", imei)
			cmd := exec.Command("./ffmpeg", "-y", "-i", src, dst)
			err = cmd.Run()
			if err != nil {
				fmt.Println("failed to convert AMR to WAV: %w", err)
			}
			playWav(dst)
			result = "测试完成"
			break
		}
		time.Sleep(time.Millisecond * 100)
		cnt++
	}

	pass.stopReader = true
	pass.stopWriter = true

	return result
}

// 用于录音测试工具
func DoTestRecord(mw *walk.MainWindow, portName string, readImei *walk.LineEdit, resultEdit *walk.LineEdit) {
	myPort := GetPort(portName)

	//首先清空窗口颜色
	brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 255, 255))
	resultEdit.SetBackground(brush)

	//读取IMEI
	var imei string
	if myPort.Name == portName {
		items := GetReadImeiTestItems()
		pass := new(PassParam)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			readPort(myPort, pass)
		}()
		imei = readCommImei(myPort, items, pass)
		mw.Synchronize(func() {
			readImei.SetText(imei)
			resultEdit.SetText("获取IMEI成功")
		})
		wg.Wait() //要等reader结束才能开启下一个reader，不然下一个reader的第一条消息可能被上一个读走
	}

	if imei == "失败" || imei == "获取超时" {
		brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
		resultEdit.SetBackground(brush)
		resultEdit.SetText("读取IMEI失败")
		return
	}

	//开始测试录音
	recordItems := GetRecodTestItems()
	pass := new(PassParam)
	go readPort(myPort, pass)
	result := checkRecordTest(mw, myPort, recordItems, pass, imei, resultEdit)
	brush, _ = walk.NewSolidColorBrush(walk.RGB(0, 255, 0))
	resultEdit.SetBackground(brush)
	resultEdit.SetText(result)
}

func playWav(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Failed to open audio file:", err)
		return
	}
	defer f.Close()

	streamer, format, err := wav.Decode(f)
	if err != nil {
		fmt.Println("Failed to decode audio file:", err)
		return
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}
