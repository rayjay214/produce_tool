// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"produce_tool/util"
	"strconv"
	"strings"
	"time"
)

func initLog() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)
	file, err := os.OpenFile("uid_tool.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to create log file: ", err)
	} else {
		log.SetOutput(file)
	}
}

func readCsvFile(selectedFilePath string) (rows [][]string, err error) {
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

func changeFileProcess(mw *walk.MainWindow, wrote, left *walk.TextLabel, selectedFilePath string, records [][]string) {
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

func writeUid(myPort *util.MyPort, pass *util.PassParam, uid string, processEcho *walk.LineEdit) error {
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
			fmt.Println("超时")
			err = errors.New("超时")
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
			fmt.Println("写入成功")
			processEcho.SetText(fmt.Sprintf("%v:写入成功", uid))
			break
		}
		//写入uid失败
		if strings.Contains(util.GetPassParamStr(pass), "<ACK> 400 Unknown command") && bTestSuccess {
			fmt.Println("写入uid失败")
			processEcho.SetText(fmt.Sprintf("%v:写入失败", uid))
			err = errors.New("写入uid失败")
			break
		}
	}

	util.StopReader(pass)
	util.StopWriter(pass)

	return err
}

func process(portName string, uid string, processEcho *walk.LineEdit, mw *walk.MainWindow, wrote, left *walk.TextLabel, selectedFilePath string, rows [][]string) {
	err := writeUidProcess(uid, portName, processEcho)
	if err == nil {
		changeFileProcess(mw, wrote, left, selectedFilePath, rows)
	} else {
		log.Errorf("write uid %v error %v", uid, err)
	}
}

func writeUidProcess(uid string, portName string, processEcho *walk.LineEdit) error {
	fmt.Println("selected uid", uid)
	myPort := util.GetPort(portName)

	pass := new(util.PassParam)
	go util.ReadPort(myPort, pass)
	return writeUid(myPort, pass, uid, processEcho)
}

func runUidWriteWindow() {
	mw, _ := walk.NewMainWindow()

	fontFamily := "Microsoft YaHei"
	viceFontSize := 12

	var selectedCom *walk.ComboBox
	var selectedFile *walk.LineEdit
	var resultEdit *walk.LineEdit

	var dlg *walk.FileDialog

	dlg = &walk.FileDialog{
		Title:  "请选择uid文件",
		Filter: "CSV Files (*.csv)|*.csv|All Files (*.*)|*.*", // 设置过滤器
	}

	var total *walk.TextLabel
	var wrote *walk.TextLabel
	var left *walk.TextLabel

	var listUid []string
	var selectedFilePath string

	MainWindow{
		AssignTo: &mw,
		Title:    "写号工具",
		Font:     Font{PointSize: viceFontSize, Family: fontFamily},
		Size:     Size{Width: 800, Height: 300},
		Layout:   VBox{Alignment: AlignHNearVNear},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					Composite{
						Layout: Grid{
							Columns: 2,
							Spacing: 15,
						},
						Children: []Widget{
							Label{
								Text:    "选择端口:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 30},
								MaxSize: Size{Width: 80},
							},
							ComboBox{
								AssignTo:      &selectedCom,
								Font:          Font{PointSize: viceFontSize, Family: fontFamily},
								Model:         util.WholePortList,
								BindingMember: "Name",
								DisplayMember: "Name",
							},
							PushButton{
								Text:    "选择文件",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 30},
								MaxSize: Size{Width: 80},
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
										/*
											currentRow := 0
											for {
												record, err := reader.Read()
												if err != nil {
													if err.Error() == "EOF" {
														break
													}
													fmt.Println("Error reading line:", err)
													return
												}

												if currentRow == 0 {
													nWrote, _ := strconv.Atoi(record[1])
													nLeft, _ := strconv.Atoi(record[2])
													nTotal := nWrote + nLeft
													wrote.SetText(fmt.Sprintf("已写入：%v", nWrote))
													left.SetText(fmt.Sprintf("剩余：%v", nLeft))
													total.SetText(fmt.Sprintf("总数：%v", nTotal))
												}
												currentRow++
											}
										*/
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
								AssignTo: &selectedFile,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 30},
								MaxSize:  Size{Width: 200},
								ReadOnly: true,
							},
							PushButton{
								Text:       "写入",
								Font:       Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:    Size{Width: 50, Height: 40},
								MaxSize:    Size{Width: 80, Height: 40},
								ColumnSpan: 2,
								OnClicked: func() {
									rows, err := readCsvFile(selectedFilePath)
									if err != nil {
										walk.MsgBox(mw, "打开文件错误", err.Error(), walk.MsgBoxIconInformation)
										return
									}
									nWrote, _ := strconv.Atoi(rows[0][1])
									if nWrote >= len(listUid) {
										walk.MsgBox(mw, "写入错误", "uid已全部写完", walk.MsgBoxIconInformation)
										return
									}

									go process(selectedCom.Text(), listUid[nWrote], resultEdit, mw, wrote, left, selectedFilePath, rows)
									/*
										err = writeUidProcess(listUid[nWrote], selectedCom.Text(), resultEdit)
										if err == nil {
											changeFileProcess(mw, wrote, left, selectedFilePath, rows)
										} else {
											log.Errorf("write uid %v error %v", listUid[nWrote], err)
										}

									*/
								},
							},
							TextLabel{
								AssignTo:   &total,
								ColumnSpan: 2,
								Text:       "总数：0",
							},
							TextLabel{
								AssignTo:   &wrote,
								ColumnSpan: 2,
								Text:       "已写入：0",
							},
							TextLabel{
								AssignTo:   &left,
								ColumnSpan: 2,
								Text:       "剩余：0",
							},
						},
					},

					LineEdit{
						AssignTo:      &resultEdit,
						TextAlignment: AlignCenter,
						Text:          "串口输出",
						Font: Font{
							PointSize: 20,
						},
					},
				},
			},
		},
	}.Run()
}

func main() {
	initLog()
	runUidWriteWindow()
}
