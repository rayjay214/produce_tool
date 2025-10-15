// Copyright 2023 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"produce_tool/util"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var version = "V1.1"

type SerialPortItem struct {
	Name     string
	Selected bool
	Status   string
}

func initLog() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)
	file, err := os.OpenFile("multi_uid_tool.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to create log file: ", err)
	} else {
		log.SetOutput(file)
	}
}

func readCsvFile2(selectedFilePath string) (rows [][]string, err error) {
	file, err := os.Open(selectedFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func changeFileProcess2(mw *walk.MainWindow, wrote, left *walk.TextLabel, selectedFilePath string, records [][]string) {
	var nWrote, nLeft int
	if len(records) > 0 {
		nWrote, _ = strconv.Atoi(records[0][1])
		nLeft, _ = strconv.Atoi(records[0][2])
		nWrote += 1
		nLeft -= 1
		records[0][1] = fmt.Sprintf("%v", nWrote)
		records[0][2] = fmt.Sprintf("%v", nLeft)
	}

	file, err := os.Create(selectedFilePath)
	if err != nil {
		walk.MsgBox(mw, "打开文件错误", "请先关闭其他打开该文件的程序", walk.MsgBoxIconInformation)
		return
	}
	defer file.Close()

	fmt.Println("before", records)

	writer := csv.NewWriter(file)
	err = writer.WriteAll(records)
	if err != nil {
		fmt.Println("Error writing CSV:", err)
		return
	}
	writer.Flush()
	wrote.SetText(fmt.Sprintf("已写入：%v", nWrote))
	left.SetText(fmt.Sprintf("剩余：%v", nLeft))
	fmt.Println("write success")
}

func writeUid2(myPort *util.MyPort, pass *util.PassParam, uid string, processEcho *walk.LineEdit) error {
	var err error
	processEcho.SetText("sysinfo test")
	_, err = myPort.Port.Write([]byte("sysinfo test\r\n"))
	if err != nil {
		log.Errorf("err is %v, port is %v", err, myPort.Name)
		myPort.Vaild = false
		return err
	}
	bTestSuccess := false

	timeout := 3 * time.Second
	startTime := time.Now()
	for {
		if time.Since(startTime) >= timeout {
			err = errors.New("timeout")
			processEcho.SetText(fmt.Sprintf("%v:超时", uid))
			break
		}
		//测试串口通过
		if strings.Contains(util.GetPassParamStr(pass), "<ACK> 200 OK") && !bTestSuccess {
			bTestSuccess = true
			util.SetPassParamStr(pass, "")
			processEcho.SetText(fmt.Sprintf("sysinfo setuid uid=%v", uid))
			_, err = myPort.Port.Write([]byte(fmt.Sprintf("sysinfo setuid uid=%v\r\n", uid)))
			if err != nil {
				log.Errorf("err is %v, port is %v", err, myPort.Name)
				myPort.Vaild = false
				break
			}
		}
		//写入uid成功
		if strings.Contains(util.GetPassParamStr(pass), "<ACK> 200 OK") && bTestSuccess {
			log.Infof("%v write uid success", uid)
			processEcho.SetText(fmt.Sprintf("%v:写入成功", uid))
			break
		}
		//写入uid失败
		if strings.Contains(util.GetPassParamStr(pass), "<ACK> 400 Unknown command") && bTestSuccess {
			log.Infof("%v write uid failed", uid)
			processEcho.SetText(fmt.Sprintf("%v:写入失败", uid))
			err = errors.New("write uid failed")
			break
		}
	}

	util.StopReaderSafe(pass)
	util.StopWriter(pass)

	return err
}

func process2(portName string, uid string, processEcho *walk.LineEdit, mw *walk.MainWindow, wrote, left *walk.TextLabel, selectedFilePath string, rows [][]string) {
	err := writeUidProcess(uid, portName, processEcho)
	if err == nil {
		changeFileProcess2(mw, wrote, left, selectedFilePath, rows)
	} else {
		log.Errorf("write uid %v error %v", uid, err)
	}
}

func writeUidProcess(uid string, portName string, processEcho *walk.LineEdit) error {
	myPort := util.GetPort(portName)

	pass := util.NewPassParam()
	go util.ReadPortSafe(myPort, pass)
	return writeUid2(myPort, pass, uid, processEcho)
}

func main() {
	initLog()
	//util.OpenAllPorts()
	runSerialDisplayWindow()
}

func runSerialDisplayWindow() {
	mw, _ := walk.NewMainWindow()

	fontFamily := "Microsoft YaHei"
	viceFontSize := 12

	// 创建串口项列表
	serialItems := make([]*SerialPortItem, 0)
	for _, port := range util.WholePortList {
		serialItems = append(serialItems, &SerialPortItem{
			Name:     port.Name,
			Selected: false,
			Status:   "",
		})
	}

	// 创建UI组件
	var selectedFile *walk.LineEdit
	var dlg *walk.FileDialog
	var total *walk.TextLabel
	var wrote *walk.TextLabel
	var left *walk.TextLabel

	var listUid []string
	var selectedFilePath string

	// 创建串口复选框和状态文本框
	comCheckBoxes := make([]*walk.CheckBox, len(serialItems))
	statusLineEdits := make([]*walk.LineEdit, len(serialItems))

	dlg = &walk.FileDialog{
		Title:  "请选择文件",
		Filter: "CSV Files (*.csv)|*.csv|All Files (*.*)|*.*",
	}

	MainWindow{
		AssignTo: &mw,
		Title:    fmt.Sprintf("多串口写UID工具%v", version),
		//Title:  fmt.Sprintf("多串口写UID工具"),
		Font:   Font{PointSize: viceFontSize, Family: fontFamily},
		Size:   Size{Width: 600, Height: 400},
		Layout: VBox{Alignment: AlignHNearVNear},
		Children: []Widget{
			// 串口列表区域
			Composite{
				Layout: VBox{MarginsZero: true},
				Children: func() []Widget {
					var widgets []Widget

					// 为每个串口创建一行
					// todo 限定不超过8个
					for i, item := range serialItems {
						// 创建局部变量来固定当前循环的值
						currentIndex := i
						currentItem := item

						widgets = append(widgets, Composite{
							Layout: HBox{MarginsZero: true},
							Children: []Widget{
								Label{
									Text:    currentItem.Name,
									MinSize: Size{Width: 60},
									Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								},
								CheckBox{
									AssignTo: &comCheckBoxes[currentIndex],
									MinSize:  Size{Width: 20},
									Checked:  currentItem.Selected,
									OnCheckedChanged: func() {
										// 更新选中状态
										serialItems[currentIndex].Selected = comCheckBoxes[currentIndex].Checked()
									},
									Text: "选中写入",
								},
								LineEdit{
									AssignTo:   &statusLineEdits[currentIndex],
									ReadOnly:   true,
									Text:       currentItem.Status,
									MinSize:    Size{Width: 400},
									Font:       Font{PointSize: viceFontSize, Family: fontFamily},
									Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
								},
							},
						})
					}

					return widgets
				}(),
			},

			// 文件选择区域
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					PushButton{
						Text:    "选择文件",
						Font:    Font{PointSize: viceFontSize, Family: fontFamily},
						MinSize: Size{Width: 80, Height: 30},
						OnClicked: func() {
							accepted, err := dlg.ShowOpen(mw)
							if err != nil {
								fmt.Println("Error or no file selected:", err)
								return
							}

							if accepted {
								selectedFilePath = dlg.FilePath
								selectedFile.SetText(filepath.Base(dlg.FilePath))
								file, err := os.Open(dlg.FilePath)
								if err != nil {
									fmt.Println("Error opening file:", err)
									return
								}
								defer file.Close()
								listUid = []string{}

								reader := csv.NewReader(file)

								//小文件，一次性读取
								rows, err := reader.ReadAll()
								if err != nil {
									fmt.Println("Error reading CSV:", err)
									return
								}
								for idx, row := range rows {
									if idx == 0 {
										nWrote, _ := strconv.Atoi(row[1])
										nLeft, _ := strconv.Atoi(row[2])
										nTotal := nWrote + nLeft

										wrote.SetText(fmt.Sprintf("已写入：%v", nWrote))
										left.SetText(fmt.Sprintf("剩余：%v", nLeft))
										total.SetText(fmt.Sprintf("总数：%v", nTotal))
									} else {
										listUid = append(listUid, row[0])
									}
								}
								fmt.Println(listUid)
							}
						},
					},
					LineEdit{
						AssignTo:   &selectedFile,
						Font:       Font{PointSize: viceFontSize, Family: fontFamily},
						MinSize:    Size{Width: 400},
						ReadOnly:   true,
						Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
					},
				},
			},

			// 写入按钮
			PushButton{
				Text:    "写入",
				Font:    Font{PointSize: viceFontSize, Family: fontFamily},
				MinSize: Size{Width: 80, Height: 40},
				MaxSize: Size{Width: 120, Height: 40},
				OnClicked: func() {
					go func() {
						for i, item := range serialItems {
							if !item.Selected {
								continue
							}
							rows, err := readCsvFile2(selectedFilePath)
							if err != nil {
								walk.MsgBox(mw, "打开文件错误", err.Error(), walk.MsgBoxIconInformation)
								return
							}
							nWrote, _ := strconv.Atoi(rows[0][1])
							if nWrote >= len(listUid) {
								walk.MsgBox(mw, "写入错误", "uid已全部写完", walk.MsgBoxIconInformation)
								return
							}

							process2(item.Name, listUid[nWrote], statusLineEdits[i], mw, wrote, left, selectedFilePath, rows)
						}
					}()
				},
			},

			// 统计信息区域
			Composite{
				Layout: VBox{MarginsZero: true},
				Children: []Widget{
					TextLabel{
						AssignTo: &total,
						Text:     "总数：0",
						Font:     Font{PointSize: viceFontSize, Family: fontFamily},
					},
					TextLabel{
						AssignTo: &wrote,
						Text:     "已写入：0",
						Font:     Font{PointSize: viceFontSize, Family: fontFamily},
					},
					TextLabel{
						AssignTo: &left,
						Text:     "剩余：0",
						Font:     Font{PointSize: viceFontSize, Family: fontFamily},
					},
				},
			},
		},
	}.Run()
}
