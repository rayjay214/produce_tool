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
	"path/filepath"
	"produce_tool/conf"
	"produce_tool/db"
	"produce_tool/dialog"
	"produce_tool/model"
	"produce_tool/network"
	"produce_tool/util"
	"strconv"
	"strings"
)

func init() {
	initLog()
	conf.LoadConf()
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
var selectedCount *walk.ComboBox
var itemCode *walk.LineEdit
var itemDesc *walk.LineEdit
var remark *walk.LineEdit
var scanSn *walk.LineEdit
var snList *walk.TextEdit
var scanCount *walk.Label
var printTemplate *walk.LineEdit

var nScanCount int

var resultEdit *walk.LineEdit

var mpSnList = make(map[string]struct{})

func refreshPlan() {
	if selectedPlan.CurrentIndex() == -1 {
		return
	}
	plan := selectedPlan.Model().([]model.PlanInfo)[selectedPlan.CurrentIndex()]
	network.DoGetPlan(plan.Id)
}

func runSnCompareWindow() {
	mw, _ := walk.NewMainWindow()

	fontFamily := "Microsoft YaHei"
	viceFontSize := 12

	var dlg *walk.FileDialog

	dlg = &walk.FileDialog{
		Title:  "请选择打印模板文件",
		Filter: "BTW Files (*.btw)", // 设置过滤器
	}

	MainWindow{
		AssignTo: &mw,
		Title:    fmt.Sprintf("装箱工具%v", version),
		Font:     Font{PointSize: viceFontSize, Family: fontFamily},
		Size:     Size{Width: 800, Height: 350},
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
							Spacing: 10,
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
								Text:    "选择数量:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 35},
								MaxSize: Size{Width: 90},
							},
							ComboBox{
								AssignTo: &selectedCount,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								Model:    []string{"50", "70", "100", "5"},
								MaxSize:  Size{Width: 45},
							},
							Label{
								Text:    "品名规格:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 35},
								MaxSize: Size{Width: 70},
							},
							LineEdit{
								AssignTo: &itemDesc,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 50},
								MaxSize:  Size{Width: 200},
							},
							Label{
								Text:    "产品编码:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 35},
								MaxSize: Size{Width: 70},
							},
							LineEdit{
								AssignTo: &itemCode,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 50},
								MaxSize:  Size{Width: 200},
							},
							Label{
								Text:    "备注:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 35},
								MaxSize: Size{Width: 70},
							},
							LineEdit{
								AssignTo: &remark,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 50},
								MaxSize:  Size{Width: 200},
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
								MinSize:  Size{Width: 35},
								MaxSize:  Size{Width: 200},
								OnMouseDown: func(x, y int, button walk.MouseButton) {
									scanSn.SetText("")
								},
								OnKeyPress: func(key walk.Key) {
									if key == walk.KeyReturn {
										strCount := selectedCount.Text()
										nCount, _ := strconv.Atoi(strCount)
										if nScanCount >= nCount {
											walk.MsgBox(mw, "装箱超限", "装箱超限，请生成打印标签后重新装箱", walk.MsgBoxIconInformation)
											return
										}

										if network.CurrentPlan.Id == 0 {
											walk.MsgBox(nil, "Error", "请选择生产计划", walk.MsgBoxIconError)
											return
										}
										err := db.CheckPackingSn(scanSn.Text(), network.CurrentPlan.Id)
										if err != nil {
											brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
											resultEdit.SetBackground(brush)
											resultEdit.SetText(err.Error())
											return
										} else {
											brush, _ := walk.NewSolidColorBrush(walk.RGB(0, 255, 0))
											resultEdit.SetBackground(brush)
											resultEdit.SetText("PASS")
										}
										if _, exists := mpSnList[scanSn.Text()]; exists {
											walk.MsgBox(mw, "重复扫描", "该SN已被扫描，请勿重复扫描", walk.MsgBoxIconInformation)
											scanSn.SetFocus()
											return
										}
										mpSnList[scanSn.Text()] = struct{}{}
										strSnList := snList.Text()
										strSnList = strSnList + scanSn.Text() + "\r\n"
										snList.SetText(strSnList)
										nScanCount++
										scanCount.SetText(fmt.Sprintf("已扫描数量:%v", nScanCount))
										if nScanCount >= nCount {
											walk.MsgBox(mw, "装箱完成", "装箱完成，请生成打印标签模板", walk.MsgBoxIconInformation)
											return
										}
										scanSn.SetText("")
										scanSn.SetFocus()
									}
								},
							},
							Label{
								AssignTo:   &scanCount,
								Text:       fmt.Sprintf("已扫描数量: %v", nScanCount),
								Font:       Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:    Size{Width: 50},
								MaxSize:    Size{Width: 150},
								ColumnSpan: 2,
							},
							TextEdit{
								AssignTo:   &snList,
								Text:       "",
								VScroll:    true, // 启用垂直滚动条
								HScroll:    false,
								ColumnSpan: 2,
								MinSize:    Size{Height: 100},
								MaxSize:    Size{Height: 500},
								ReadOnly:   true,
							},
							PushButton{
								Text:    "生成装箱标签",
								Font:    Font{PointSize: 9, Family: fontFamily},
								MinSize: Size{Width: 80, Height: 30},
								MaxSize: Size{Width: 300, Height: 30},
								OnClicked: func() {
									nSelectCnt, _ := strconv.Atoi(selectedCount.Text())
									if nScanCount != nSelectCnt {
										walk.MsgBox(mw, "装箱错误", fmt.Sprintf("已扫描数量(%v)与选择数量(%v)不一致，请检查", nScanCount, nSelectCnt), walk.MsgBoxIconInformation)
										return
									}
									if itemCode.Text() == "" || itemDesc.Text() == "" || remark.Text() == "" {
										walk.MsgBox(mw, "装箱错误", fmt.Sprintf("请先补全信息"), walk.MsgBoxIconInformation)
										return
									}

									boxNo := network.DoGetBoxNo()
									if boxNo == "" {
										walk.MsgBox(mw, "网络异常", "获取箱号失败", walk.MsgBoxIconInformation)
										return
									}
									param := util.BtwParam{
										KTXSN:    boxNo,
										PO:       network.CurrentPlan.OrderNo,
										QTY:      selectedCount.Text(),
										ITEMNAME: remark.Text(),
										ITEMCODE: itemCode.Text(),
										ITEMDESC: itemDesc.Text(),
									}
									param.SNLIST = make([]string, 0)
									lines := strings.Split(snList.Text(), "\r\n")

									for _, line := range lines {
										if line != "" {
											param.SNLIST = append(param.SNLIST, line)
										}
									}
									dstFilename := fmt.Sprintf("%v.btw", boxNo)
									srcFilename := fmt.Sprintf("ktx%v.btw", selectedCount.Text())

									go func() {
										err := util.GenBtwFile(srcFilename, dstFilename, param)
										if err != nil {
											mw.Synchronize(func() {
												walk.MsgBox(mw, "生成失败", "箱码标签生成失败："+err.Error(), walk.MsgBoxIconError)
											})
											return
										}

										// 成功后回主线程继续执行后续操作
										mw.Synchronize(func() {
											walk.MsgBox(mw, "成功", "箱码标签生成成功", walk.MsgBoxIconInformation)

											box := db.Box{
												BoxNo:    boxNo,
												PlanId:   network.CurrentPlan.Id,
												Count:    int64(len(param.SNLIST)),
												ItemDesc: param.ITEMDESC,
												ItemCode: param.ITEMCODE,
												Remark:   param.ITEMNAME,
											}

											db.InsertBoxRecord(box, param.SNLIST)
											mpSnList = make(map[string]struct{})
											snList.SetText("")
											nScanCount = 0
											scanCount.SetText(fmt.Sprintf("已扫描数量:%v", nScanCount))

											// 设置打印模板路径
											cwd, _ := os.Getwd()
											templatePath := filepath.Join(cwd, dstFilename)
											printTemplate.SetText(templatePath)
										})
									}()
								},
								ColumnSpan: 2,
							},
							PushButton{
								Text:    "选择打印模板",
								Font:    Font{PointSize: 9, Family: fontFamily},
								MinSize: Size{Width: 80, Height: 30},
								MaxSize: Size{Width: 300, Height: 30},
								OnClicked: func() {
									accepted, err := dlg.ShowOpen(mw)
									if err != nil {
										fmt.Println("Error or no file selected:", err)
										return
									}
									if accepted {
										printTemplate.SetText(dlg.FilePath)
									}
								},
							},
							LineEdit{
								AssignTo: &printTemplate,
								Font:     Font{PointSize: 9, Family: fontFamily},
								MinSize:  Size{Width: 80, Height: 30},
								MaxSize:  Size{Width: 300, Height: 30},
							},
							PushButton{
								Text:       "打印箱码",
								Font:       Font{PointSize: 9, Family: fontFamily},
								MinSize:    Size{Width: 80, Height: 30},
								MaxSize:    Size{Width: 300, Height: 30},
								ColumnSpan: 2,
								OnClicked: func() {
									go func() {
										err := util.PrintFileOle(printTemplate.Text())
										// 确保回到主线程更新 UI
										mw.Synchronize(func() {
											if err != nil {
												walk.MsgBox(mw, "打印失败", err.Error(), walk.MsgBoxIconError)
											} else {
												walk.MsgBox(mw, "打印成功", "箱码打印成功", walk.MsgBoxIconInformation)
											}
										})
									}()
								},
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

	runSnCompareWindow()
}
