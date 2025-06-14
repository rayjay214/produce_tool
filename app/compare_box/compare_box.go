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

var version = "V2.1"
var selectedPlan *walk.ComboBox
var boxSn *walk.LineEdit
var deviceSn *walk.LineEdit

var resultEdit *walk.TextEdit

func refreshPlan() {
	if selectedPlan.CurrentIndex() == -1 {
		return
	}
	plan := selectedPlan.Model().([]model.PlanInfo)[selectedPlan.CurrentIndex()]
	network.DoGetPlan(plan.Id)
}

func compareBoxSn() bool {
	strBoxSn := boxSn.Text()
	strDeviceSn := deviceSn.Text()

	if strBoxSn == strDeviceSn {
		err := db.CheckBoxSn(strBoxSn, network.CurrentPlan.Id)
		if err != nil {
			brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
			resultEdit.SetBackground(brush)
			errMsg := fmt.Sprintf("%s, 彩盒SN:%s, 机身SN:%s", err.Error(), strBoxSn, strDeviceSn)
			resultEdit.SetText(errMsg)
			return false
		} else {
			brush, _ := walk.NewSolidColorBrush(walk.RGB(0, 255, 0))
			resultEdit.SetBackground(brush)
			resultEdit.SetText("PASS")
			db.BoxStatus(strBoxSn)
			return true
		}
	} else {
		brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
		resultEdit.SetBackground(brush)
		errMsg := fmt.Sprintf("FAIL, 彩盒SN:%s, 机身SN:%s", strBoxSn, strDeviceSn)
		resultEdit.SetText(errMsg)
		return false
	}
}

func runSnCompareWindow() {
	mw, _ := walk.NewMainWindow()

	fontFamily := "Microsoft YaHei"
	viceFontSize := 12

	MainWindow{
		AssignTo: &mw,
		Title:    fmt.Sprintf("彩盒标机身标比对工具%v", version),
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
								Text:    "彩盒SN:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 35},
								MaxSize: Size{Width: 60},
							},
							LineEdit{
								AssignTo: &boxSn,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 50},
								MaxSize:  Size{Width: 200},
								OnMouseDown: func(x, y int, button walk.MouseButton) {
									boxSn.SetText("")
								},
								OnKeyPress: func(key walk.Key) {
									if key == walk.KeyReturn {
										deviceSn.SetFocus()
									}
								},
							},
							Label{
								Text:    "机身SN:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 35},
								MaxSize: Size{Width: 60},
							},
							LineEdit{
								AssignTo: &deviceSn,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 35},
								MaxSize:  Size{Width: 200},
								OnMouseDown: func(x, y int, button walk.MouseButton) {
									deviceSn.SetText("")
								},
								OnKeyPress: func(key walk.Key) {
									if key == walk.KeyReturn {
										if network.CurrentPlan.Id == 0 {
											walk.MsgBox(nil, "Error", "请选择生产计划", walk.MsgBoxIconError)
											return
										}
										bSuccess := compareBoxSn()
										if bSuccess {
											boxSn.SetText("")
											deviceSn.SetText("")
											boxSn.SetFocus()
										} else {
											boxSn.SetText("")
											deviceSn.SetText("")
											boxSn.SetFocus()
										}
									}
								},
							},
						},
					},
					TextEdit{
						AssignTo:      &resultEdit,
						TextAlignment: AlignNear,
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
	loginResult := dialog.ShowLoginDialog()

	if !loginResult.Success {
		fmt.Println("登录失败:", loginResult.Message)
		os.Exit(1)
	}

	runSnCompareWindow()
}
