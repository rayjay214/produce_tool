// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"produce_tool/db"
	"produce_tool/util"
)

func init() {
	db.LoadCheckSnCsv()
}

func runSnCompareWindow() {
	mw, _ := walk.NewMainWindow()

	fontFamily := "Microsoft YaHei"
	viceFontSize := 12

	var selectedCom *walk.ComboBox
	var scanSn *walk.LineEdit
	var readSn *walk.LineEdit
	var readImei *walk.LineEdit
	var imeiPrefix *walk.LineEdit
	//var resultButton *walk.PushButton
	var resultEdit *walk.LineEdit

	MainWindow{
		AssignTo: &mw,
		Title:    "SN比对工具",
		Font:     Font{PointSize: viceFontSize, Family: fontFamily},
		Size:     Size{Width: 600, Height: 350},
		Layout:   VBox{Alignment: AlignHNearVNear},
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
								Text:    "选择端口:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 35},
								MaxSize: Size{Width: 80},
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
										util.DoTestOnePortCompareSn(selectedCom.Text(), scanSn, imeiPrefix.Text(), readSn, readImei, resultEdit)
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
	runSnCompareWindow()
}
