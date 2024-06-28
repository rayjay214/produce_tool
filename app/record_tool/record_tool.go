// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"os"
	"produce_tool/util"
)

func runRecordTestWindow() {
	mw, _ := walk.NewMainWindow()

	fontFamily := "Microsoft YaHei"
	viceFontSize := 12

	var selectedCom *walk.ComboBox
	var readImei *walk.LineEdit
	var resultEdit *walk.LineEdit

	MainWindow{
		AssignTo: &mw,
		Title:    "录音测试工具",
		Font:     Font{PointSize: viceFontSize, Family: fontFamily},
		Size:     Size{Width: 600, Height: 190},
		Layout:   VBox{Alignment: AlignHNearVNear},
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
								Text:    "设备IMEI:",
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
							PushButton{
								Text:       "开始测试",
								Font:       Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:    Size{Width: 50, Height: 40},
								MaxSize:    Size{Width: 80, Height: 40},
								ColumnSpan: 2,
								OnClicked: func() {
									//需要实时显示测试流程，不能阻塞渲染主线程
									go util.DoTestRecord(mw, selectedCom.Text(), readImei, resultEdit)
								},
							},
						},
					},

					LineEdit{
						AssignTo:      &resultEdit,
						TextAlignment: AlignCenter,
						Text:          "测试开始",
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
	dirPath := "record"
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.Mkdir(dirPath, os.ModePerm)
	}
	runRecordTestWindow()
}
