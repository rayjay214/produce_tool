// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
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
	"strconv"
)

func init() {
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

var version = "V2.6"
var selectedPlan *walk.ComboBox
var scanSn *walk.LineEdit
var resultEdit *walk.LineEdit
var printCnt *walk.LineEdit

func refreshPlan() {
	if selectedPlan.CurrentIndex() == -1 {
		return
	}
	plan := selectedPlan.Model().([]model.PlanInfo)[selectedPlan.CurrentIndex()]
	network.DoGetPlan(plan.Id)
}

func runPrintWindow(btAppDispatch *ole.IDispatch) {
	mw, _ := walk.NewMainWindow()

	fontFamily := "Microsoft YaHei"
	viceFontSize := 12

	MainWindow{
		AssignTo: &mw,
		Title:    fmt.Sprintf("标签打印工具%v", version),
		Font:     Font{PointSize: viceFontSize, Family: fontFamily},
		Size:     Size{Width: 600, Height: 200},
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
							Spacing: 25,
						},
						Children: []Widget{
							Label{
								Text:    "选择计划:",
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
								Text:    "扫描打印:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 35},
								MaxSize: Size{Width: 90},
							},
							LineEdit{
								AssignTo: &scanSn,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 35},
								MaxSize:  Size{Width: 200},
								OnMouseDown: func(x, y int, button walk.MouseButton) {
									scanSn.SetText("")
								},
								OnKeyPress: func(key walk.Key) {
									if key == walk.KeyReturn {
										err := db.CheckPrintSn(scanSn.Text(), network.CurrentPlan.Id)
										if err != nil {
											brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
											resultEdit.SetBackground(brush)
											resultEdit.SetText(err.Error())
											return
										} else {
											brush, _ := walk.NewSolidColorBrush(walk.RGB(0, 255, 0))
											resultEdit.SetBackground(brush)
											resultEdit.SetText("打印中")
										}
										//开始打印
										cnt := 1
										if printCnt.Text() != "" {
											cnt, _ = strconv.Atoi(printCnt.Text())
										}
										err = util.PrintInMemory(btAppDispatch, fmt.Sprintf("%v.btw", network.CurrentPlan.DeviceType), scanSn.Text(), cnt)
										if err != nil {
											resultEdit.SetText("打印失败")
											return
										}
										db.PrintStatus(scanSn.Text())
										resultEdit.SetText("打印完成")
										scanSn.SetText("")
									}
								},
							},
							Label{
								Text:    "打印次数:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 35},
								MaxSize: Size{Width: 90},
							},
							LineEdit{
								AssignTo: &printCnt,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 35},
								MaxSize:  Size{Width: 200},
								Text:     "1",
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

	ole.CoInitialize(0)
	defer ole.CoUninitialize()
	btApp, err := oleutil.CreateObject("BarTender.Application")
	if err != nil {
		return
	}
	btAppDispatch, err := btApp.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return
	}
	defer btAppDispatch.Release()

	runPrintWindow(btAppDispatch)
	_, _ = oleutil.CallMethod(btAppDispatch, "Quit")
}
