// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	log "github.com/sirupsen/logrus"
	"os"
	"produce_tool/db"
	"produce_tool/dialog"
	"produce_tool/model"
	"produce_tool/network"
	"produce_tool/util"
)

func init() {
	db.LoadCheckSnCsv()
	initLog()
	db.InitMysql()
}

func initLog() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)
	file, err := os.OpenFile("tool.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to create log file: ", err)
	} else {
		log.SetOutput(file)
	}
}

var version = "V2.2"
var selectedPlan *walk.ComboBox
var selectedCom *walk.ComboBox
var scanSn *walk.LineEdit
var readSn *walk.LineEdit
var readImei *walk.LineEdit
var imeiPrefix *walk.LineEdit

// var resultButton *walk.PushButton
var resultEdit *walk.LineEdit
var onlyCompareSn *walk.CheckBox

func refreshPlan() {
	if selectedPlan.CurrentIndex() == -1 {
		return
	}
	plan := selectedPlan.Model().([]model.PlanInfo)[selectedPlan.CurrentIndex()]
	network.DoGetPlan(plan.Id)
	if network.CurrentPlan.SnType == "0" {
		onlyCompareSn.SetChecked(true)
	}
}

func runSnCompareWindow() {
	mw, _ := walk.NewMainWindow()

	fontFamily := "Microsoft YaHei"
	viceFontSize := 12

	MainWindow{
		AssignTo: &mw,
		Title:    fmt.Sprintf("SN比对工具%v", version),
		Font:     Font{PointSize: viceFontSize, Family: fontFamily},
		Size:     Size{Width: 600, Height: 350},
		Layout:   VBox{Alignment: AlignHNearVNear},
		OnSizeChanged: func() {
			screenWidth := int(win.GetSystemMetrics(win.SM_CXSCREEN))
			screenHeight := int(win.GetSystemMetrics(win.SM_CYSCREEN))
			bounds := mw.Bounds()

			// 计算居中位置
			x := (screenWidth - bounds.Width) / 2
			y := (screenHeight - bounds.Height) / 2
			mw.SetBounds(walk.Rectangle{X: x, Y: y, Width: bounds.Width, Height: bounds.Height})
		},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					Composite{
						Layout: Grid{
							Columns: 2,
							Spacing: 30,
						},
						Children: []Widget{
							Label{
								Text:    "选择生产计划:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 35},
								MaxSize: Size{Width: 90},
							},
							ComboBox{
								AssignTo:              &selectedPlan,
								Font:                  Font{PointSize: viceFontSize, Family: fontFamily},
								Model:                 model.AllPlans,
								BindingMember:         "Id",
								DisplayMember:         "DisplayName",
								OnCurrentIndexChanged: refreshPlan,
								MaxSize:               Size{Width: 45},
							},
							Label{
								Text:    "选择端口:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 35},
								MaxSize: Size{Width: 90},
							},
							ComboBox{
								AssignTo:      &selectedCom,
								Font:          Font{PointSize: viceFontSize, Family: fontFamily},
								Model:         util.WholePortList,
								BindingMember: "Name",
								DisplayMember: "Name",
							},
							Label{
								Text:    "扫描SN:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 35},
								MaxSize: Size{Width: 60},
							},
							LineEdit{
								AssignTo: &scanSn,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 50},
								MaxSize:  Size{Width: 200},
								OnKeyPress: func(key walk.Key) {
									if key == walk.KeyReturn {
										if network.CurrentPlan.Id == 0 {
											walk.MsgBox(nil, "Error", "请选择生产计划", walk.MsgBoxIconError)
											return
										}
										util.DoTestOnePortCompareSn(selectedCom.Text(), scanSn, imeiPrefix.Text(), readSn, readImei, resultEdit, onlyCompareSn.Checked())
									}
								},
								OnMouseDown: func(x, y int, button walk.MouseButton) {
									scanSn.SetText("")
								},
							},
							Label{
								Text:    "读取SN:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 35},
								MaxSize: Size{Width: 60},
							},
							LineEdit{
								AssignTo: &readSn,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 35},
								MaxSize:  Size{Width: 200},
								ReadOnly: true,
							},
							Label{
								Text:    "读取IMEI:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 35},
								MaxSize: Size{Width: 80},
							},
							LineEdit{
								AssignTo: &readImei,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 35},
								MaxSize:  Size{Width: 200},
								ReadOnly: true,
							},
							Label{
								Text:    "IMEI前缀:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 35},
								MaxSize: Size{Width: 80},
							},
							LineEdit{
								AssignTo: &imeiPrefix,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 35},
								MaxSize:  Size{Width: 200},
							},
							CheckBox{
								Text:       "只比对SN",
								Font:       Font{PointSize: viceFontSize, Family: fontFamily},
								AssignTo:   &onlyCompareSn,
								MinSize:    Size{Width: 60, Height: 25},
								MaxSize:    Size{Width: 200, Height: 25},
								Enabled:    true,
								ColumnSpan: 2,
							},
						},
					},

					LineEdit{
						AssignTo:      &resultEdit,
						TextAlignment: AlignCenter,
						Font: Font{
							PointSize: 30,
						},
					},

					/*
						PushButton{
							AssignTo: &resultButton,
							Font: Font{
								PointSize: 20,
								Family:    fontFamily,
							},
							Enabled: false,
						},
					*/
				},
			},
		},
	}.Run()
}

func main() {
	loginResult := dialog.ShowLoginDialog()

	if !loginResult.Success {
		fmt.Println("登录失败:", loginResult.Message)
		os.Exit(1)
	}

	runSnCompareWindow()
}
