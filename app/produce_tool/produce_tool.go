// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/lxn/win"
	"os"
	"produce_tool/conf"
	"produce_tool/db"
	"produce_tool/dialog"
	"produce_tool/model"
	"produce_tool/network"
	"produce_tool/util"
	"reflect"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	log "github.com/sirupsen/logrus"
)

var version = "V2.8"

var tv *walk.TableView
var tableColumns []TableViewColumn
var tableModel *util.MyTableModel
var singleFunctionButtons []*walk.PushButton
var singleFunctionMapping map[int]string
var selectedCb *walk.ComboBox
var onnKeyTestBtn *walk.PushButton
var planCb *walk.ComboBox

// 写号相关
var controlBtn *walk.PushButton
var checkSn *walk.CheckBox
var checkImei *walk.CheckBox
var textHeader *walk.TextEdit
var textSn *walk.TextEdit

// 屏蔽com口
var blockedCom *walk.TextEdit

// 测完是否关机
var checkPowerOff *walk.CheckBox

// usb模式
var usbMode *walk.CheckBox

// 阈值校验
var checkSignalMin *walk.TextEdit
var checkSignalMax *walk.TextEdit
var checkGpsValue *walk.TextEdit
var checkWifiValue *walk.TextEdit
var checkPowerValue *walk.TextEdit

// 比较校验
var compareVersion *walk.TextEdit
var compareMainIp *walk.TextEdit
var compareViceIp *walk.TextEdit

// 待修改值
var modifyIp *walk.TextEdit
var modifyApn *walk.TextEdit
var modifyProtocol *walk.TextEdit

// 通过数量
var totalCnt *walk.LineEdit
var passedCnt *walk.LineEdit

func init() {
	initLog()
	initTableColumns()
	initSingleFunctionButtons()
	initRefreshTimer()
	initSyncConfTimer()
	initRefreshCountTimer()
	initConf()
	db.InitMysql()
	//model.LoadDeviceType()
	//model.LoadDeviceTypeNetwork()
	db.LoadTestRstCsv()
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

func initConf() {
	conf.LoadConf()
}

func initSingleFunctionButtons() {
	singleFunctionMapping = make(map[int]string, 0)
	singleFunctionMapping[0] = "SIM"
	singleFunctionMapping[1] = "信号"
	singleFunctionMapping[2] = "重力"
	singleFunctionMapping[3] = "WIFI"
	singleFunctionMapping[4] = "IMEI"
	singleFunctionMapping[5] = "光感"
	singleFunctionMapping[6] = "回音"
	singleFunctionMapping[7] = "防拆"
	singleFunctionMapping[8] = "型号"
	singleFunctionMapping[9] = "GPS"

	singleFunctionButtons = make([]*walk.PushButton, 0)
	for i := 0; i < len(singleFunctionMapping); i++ {
		button := new(walk.PushButton)
		singleFunctionButtons = append(singleFunctionButtons, button)
	}
}

func initTableColumns() {
	items := util.GetAllTestItems()
	tableColumns = []TableViewColumn{
		TableViewColumn{Title: "COM", Width: 70, Name: "Com"},
		TableViewColumn{Title: "通过", Width: 90, Name: "Pass"},
	}

	for _, item := range items {
		if item.IsShow {
			column := TableViewColumn{Title: item.Desc, Name: item.ModelColName}
			tableColumns = append(tableColumns, column)
		}
	}
}

func initSyncConfTimer() {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for {
			select {
			case <-ticker.C:
				conf.SyncConf()
			}
		}
	}()
}

func initRefreshTimer() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				for i := 0; i < util.GetTableModel().RowCount(); i++ {
					util.GetTableModel().PublishRowChanged(i)
				}
			}
		}
	}()
}

func initRefreshCountTimer() {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ticker.C:
				if network.CurrentPlan.Id != 0 {
					network.DoGetPassedNum()
					passedCnt.SetText(fmt.Sprintf("通过数量：%v", network.PassedCount))
				}
			}
		}
	}()
}

func styleFunc(style *walk.CellStyle) {
	font, _ := walk.NewFont("Microsoft YaHei", 12, 0)
	style.Font = font
	items := util.GetModelItems()
	item := items[style.Row()]
	rv := reflect.ValueOf(item)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	rt := rv.Type()
	propertyName := util.ColumnIdxNames[style.Col()]
	if util.ContainsOne(rv.FieldByName(propertyName).String(), "失败", "超时", "已被使用", "不存在", "不属于") {
		style.BackgroundColor = walk.RGB(255, 0, 0)
	}
	if util.ContainsOne(rv.FieldByName(propertyName).String(), "写入") {
		style.BackgroundColor = walk.RGB(66, 239, 245)
	}

	if propertyName == "Pass" {
		passTest := true
		allColumnFilled := true
		waiting := false
		for i := 0; i < rv.NumField(); i++ {
			fieldValue := rv.Field(i)
			fieldType := rt.Field(i)
			_, ok := util.ColumnNamesIdx[fieldType.Name]
			if i == 2 { //固定为MES
				if util.ContainsOne(fieldValue.String(), "失败") {
					passTest = false
				}
				continue
			}
			if ok && i != 1 {
				if (fieldValue.String() == "" && tv.Columns().ByName(fieldType.Name).Visible()) || util.ContainsOne(fieldValue.String(), "失败", "超时", "等待", "已过站") {
					passTest = false
				}
				if fieldValue.String() == "" && tv.Columns().ByName(fieldType.Name).Visible() {
					allColumnFilled = false
				}
				if util.ContainsOne(fieldValue.String(), "等待") {
					waiting = true
				}
			}
		}
		if passTest {
			item.Pass = "测试通过"
			style.BackgroundColor = walk.RGB(0, 255, 0)
		} else if waiting {
			item.Pass = "等待中"
			style.BackgroundColor = walk.RGB(0, 0, 255)
		} else if allColumnFilled {
			item.Pass = "测试失败"
			style.BackgroundColor = walk.RGB(255, 0, 0)
		}
	}
}

func refreshType() {
	if selectedCb.CurrentIndex() == -1 {
		return
	}
	selectedType := selectedCb.Model().([]model.DeviceTypeInfo)[selectedCb.CurrentIndex()]
	modifyIp.SetText(fmt.Sprintf("%v:%v", selectedType.MainIp, selectedType.MainPort))
	modifyApn.SetText(selectedType.APN)
	network.DoGetDeviceType(selectedType.DeviceType)
	/*
		if network.CurrentType.SnType == "1" {
			textHeader.SetText(network.CurrentType.ImeiPrefix)
		}
	*/
	if network.CurrentType.WriteSn == "1" {
		checkSn.SetChecked(true)
	}
	if network.CurrentType.WriteImei == "1" {
		checkImei.SetChecked(true)
	}

	if selectedType.ProtocolValue == "0" {
		modifyProtocol.SetText("GT06")
	} else {
		modifyProtocol.SetText("JT808")
	}
	checkSignalMax.SetText(selectedType.SignalMax)
	checkSignalMin.SetText(selectedType.SignalMin)
	checkGpsValue.SetText(selectedType.GpsMin)
	checkWifiValue.SetText(selectedType.WifiMin)
	checkPowerValue.SetText(selectedType.PowerMin)
	util.SelectedDeviceType = selectedType
	util.SyncTestItems()
	util.RefreshTableModel()

	for i := 0; i < tv.Columns().Len(); i++ {
		tv.Columns().At(i).SetVisible(true)
	}
	util.FilterTableColumn(tv, selectedType)

	conf.SelectedType = selectedType.DeviceType
	conf.SyncConf()
}

func refreshPlan() {
	if planCb.CurrentIndex() == -1 {
		return
	}
	plan := planCb.Model().([]model.PlanInfo)[planCb.CurrentIndex()]
	network.DoGetPlan(plan.Id)
	network.DoGetTotalNum()
	totalCnt.SetText(fmt.Sprintf("总体数量：%v", network.TotalCount))
	listModel, _ := selectedCb.Model().([]model.DeviceTypeInfo)
	for i := 0; i < len(listModel); i++ {
		item := listModel[i]
		if item.DeviceType == plan.DeviceType {
			selectedCb.SetCurrentIndex(i)
		}
	}
}

func runMainWindow() {
	mw, _ := walk.NewMainWindow()

	tableModel = util.GetTableModel()
	btnHeight := 55
	fontSize := 12
	fontFamily := "Microsoft YaHei"
	viceFontSize := 10

	//addDeviceType := new(model.DeviceTypeInfo)

	MainWindow{
		AssignTo: &mw,
		Title:    fmt.Sprintf("生产测试工具%v", version),
		Font:     Font{PointSize: viceFontSize, Family: fontFamily},
		Size:     Size{Width: 1200, Height: 650},
		Layout:   VBox{Alignment: AlignHNearVNear},
		//设置居中
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
			Composite{
				Layout: HBox{
					Alignment: AlignHCenterVCenter,
					Margins:   Margins{Left: 0, Top: 0, Right: 0, Bottom: 0},
				},
				MaxSize: Size{Width: 1100, Height: 80},
				Children: []Widget{
					Composite{
						MinSize: Size{Width: 100, Height: 80},
						MaxSize: Size{Width: 150, Height: 80},
						Layout: Grid{
							Columns:   1,
							Spacing:   0,
							Alignment: AlignHNearVNear,
						},
						Children: []Widget{
							LineEdit{
								AssignTo:  &totalCnt,
								Text:      fmt.Sprintf("总体数量：%v", 0),
								Font:      Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:   Size{Width: 100, Height: 25},
								MaxSize:   Size{Width: 150, Height: 25},
								Alignment: AlignHNearVCenter,
								TextColor: walk.RGB(255, 0, 0),
								ReadOnly:  true,
							},
							LineEdit{
								AssignTo:  &passedCnt,
								Text:      fmt.Sprintf("通过数量：%v", 0),
								Font:      Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:   Size{Width: 100, Height: 25},
								MaxSize:   Size{Width: 150, Height: 25},
								Alignment: AlignHNearVCenter,
								TextColor: walk.RGB(255, 0, 0),
								ReadOnly:  true,
							},
						},
					},
					/*
						PushButton{
							Text:      "管理型号",
							Font:      Font{PointSize: 14, Family: fontFamily},
							Alignment: AlignHNearVCenter,
							MinSize:   Size{Width: 60, Height: 100},
							MaxSize:   Size{Width: 100, Height: 100},
							OnClicked: func() {
								dialog.RunCheckPwdDialog(mw, selectedCb)
							},
						},

					*/
					GroupBox{
						MinSize: Size{Width: 80, Height: 100},
						MaxSize: Size{Width: 350, Height: 100},
						Title:   "请选择生产计划:",
						Font:    Font{PointSize: 12, Family: fontFamily},
						Layout:  HBox{},
						Children: []Widget{
							ComboBox{
								AssignTo:      &planCb,
								Font:          Font{PointSize: viceFontSize, Family: fontFamily},
								Model:         model.AllPlans,
								BindingMember: "Id",
								DisplayMember: "DisplayName",
								MaxSize:       Size{Width: 320, Height: btnHeight},

								OnCurrentIndexChanged: refreshPlan,
							},
						},
					},
					GroupBox{
						MinSize: Size{Width: 80, Height: 100},
						MaxSize: Size{Width: 150, Height: 100},
						Title:   "请选择型号:",
						Font:    Font{PointSize: 12, Family: fontFamily},
						Layout:  HBox{},
						Children: []Widget{
							ComboBox{
								AssignTo:              &selectedCb,
								Font:                  Font{PointSize: viceFontSize, Family: fontFamily},
								Model:                 model.AllTypes,
								BindingMember:         "DeviceType",
								DisplayMember:         "DeviceType",
								MaxSize:               Size{Width: 100, Height: btnHeight},
								OnCurrentIndexChanged: refreshType,
								Enabled:               false,
							},
						},
					},
					GroupBox{
						MinSize: Size{Width: 80, Height: 100},
						MaxSize: Size{Width: 200, Height: 100},
						Font:    Font{PointSize: 12, Family: fontFamily},
						Title:   "额外选项:",
						Layout:  HBox{},
						Children: []Widget{
							CheckBox{
								AssignTo: &checkPowerOff,
								Text:     "测完关机",
								Font:     Font{PointSize: 10, Family: fontFamily},
								MinSize:  Size{Width: 30, Height: 25},
								MaxSize:  Size{Width: 80, Height: 25},
								OnCheckedChanged: func() {
									util.PoweroffAfterTest = checkPowerOff.Checked()
								},
							},
							CheckBox{
								AssignTo: &usbMode,
								Text:     "usb模式",
								Font:     Font{PointSize: 10, Family: fontFamily},
								MinSize:  Size{Width: 30, Height: 25},
								MaxSize:  Size{Width: 80, Height: 25},
								OnCheckedChanged: func() {
									util.IsUsbMode = usbMode.Checked()
								},
							},
						},
					},
				},
			},
			Composite{
				Layout: HBox{
					Alignment: AlignHNearVNear,
					Margins:   Margins{Left: 0, Top: 0, Right: 0, Bottom: 0},
				},
				Children: []Widget{
					GroupBox{
						Title:  "COM信息",
						Font:   Font{PointSize: fontSize, Family: fontFamily},
						Layout: VBox{Alignment: AlignHNearVNear},
						Children: []Widget{
							ScrollView{
								Layout:  VBox{},
								MinSize: Size{Width: 1000, Height: 450},
								Children: []Widget{
									TableView{
										AssignTo: &tv,
										//AlternatingRowBG: true,
										Columns: tableColumns,
										Model:   tableModel,
										OnItemActivated: func() {
											for portName, idx := range util.PortNameRowidx {
												if idx == tv.CurrentIndex() {
													if selectedCb.CurrentIndex() < 0 {
														walk.MsgBox(nil, "Error", "请选择型号", walk.MsgBoxIconError)
														return
													}
													util.DoTestOnePortAllItems(portName, idx)
												}
											}
										},
										StyleCell: styleFunc,
									},
								},
							},
						},
					},
				},
			},
			HSplitter{
				Children: []Widget{
					Composite{
						Layout: Flow{Alignment: AlignHCenterVCenter},
						Children: []Widget{
							GroupBox{
								Alignment: AlignHCenterVCenter,
								Title:     "修改IMEI或SN",
								Font:      Font{PointSize: viceFontSize, Family: fontFamily},
								Layout:    Grid{Columns: 3},
								MinSize:   Size{Width: 250, Height: 110},
								MaxSize:   Size{Width: 500, Height: 110},
								Children: []Widget{
									PushButton{
										AssignTo: &controlBtn,
										Text:     "开启写号",
										Font:     Font{PointSize: viceFontSize, Family: fontFamily},
										MinSize:  Size{Width: 50, Height: 25},
										MaxSize:  Size{Width: 80, Height: 25},
										OnClicked: func() {
											enabled := textHeader.Enabled()
											//checkSn.SetEnabled(!enabled)
											//checkImei.SetEnabled(!enabled)
											textHeader.SetEnabled(!enabled)
											textSn.SetEnabled(!enabled)
											if enabled {
												controlBtn.SetText("开启写号")
											} else {
												controlBtn.SetText("关闭写号")
											}
										},
									},
									CheckBox{
										Text:     "SN",
										Font:     Font{PointSize: viceFontSize, Family: fontFamily},
										AssignTo: &checkSn,
										MinSize:  Size{Width: 40, Height: 25},
										MaxSize:  Size{Width: 60, Height: 25},
										Enabled:  false,
									},
									CheckBox{
										Text:     "IMEI",
										Font:     Font{PointSize: viceFontSize, Family: fontFamily},
										AssignTo: &checkImei,
										MinSize:  Size{Width: 40, Height: 25},
										MaxSize:  Size{Width: 60, Height: 25},
										Enabled:  false,
									},
									TextEdit{
										Text:     "",
										Font:     Font{PointSize: viceFontSize, Family: fontFamily},
										AssignTo: &textHeader,
										MinSize:  Size{Width: 80, Height: 25},
										MaxSize:  Size{Width: 100, Height: 25},
										Enabled:  false,
									},
									TextEdit{
										Text:       "",
										Font:       Font{PointSize: viceFontSize, Family: fontFamily},
										AssignTo:   &textSn,
										MinSize:    Size{Width: 100, Height: 25},
										MaxSize:    Size{Width: 150, Height: 25},
										ColumnSpan: 2,
										Enabled:    false,
										OnKeyPress: func(key walk.Key) {
											if util.IsUsbMode {
												util.CheckPorts() //USB的需要重新打开端口，串口的不需要，可以不调用此函数
											}
											if key == walk.KeyReturn {
												if tv.CurrentIndex() < 0 || tv.CurrentIndex() > tv.Model().(*util.MyTableModel).RowCount() {
													log.Infof("row %v invalid", tv.CurrentIndex())
													tv.SetCurrentIndex(0)
												}
												if planCb.CurrentIndex() < 0 {
													walk.MsgBox(nil, "Error", "请选择生产计划", walk.MsgBoxIconError)
													return
												}
												go func(idx int, sn string, imei string) {
													if checkSn.Checked() {
														util.DoOnePortWriteSn(util.RowidxPortName[idx], sn)
													}
													if checkImei.Checked() {
														time.Sleep(time.Second)
														util.DoOnePortWriteImei(util.RowidxPortName[idx], imei)
													}
												}(tv.CurrentIndex(), textSn.Text(), textHeader.Text()+textSn.Text())
												textSn.SetText("")
												if tv.CurrentIndex()+1 >= tv.Model().(*util.MyTableModel).RowCount() {
													tv.SetCurrentIndex(-1)
												} else {
													tv.SetCurrentIndex(tv.CurrentIndex() + 1)
												}
											}
										},
									},
								},
							},
						},
					},
					Composite{
						Layout: Flow{Alignment: AlignHCenterVCenter},
						Children: []Widget{
							PushButton{
								Alignment: AlignHCenterVCenter,
								Text:      "一键测试",
								Font:      Font{PointSize: 20, Family: fontFamily},
								MinSize:   Size{Width: 250, Height: 120},
								MaxSize:   Size{Width: 500, Height: 120},
								OnClicked: func() {
									if selectedCb.CurrentIndex() < 0 {
										walk.MsgBox(nil, "Error", "请选择型号", walk.MsgBoxIconError)
										return
									}
									if planCb.CurrentIndex() < 0 {
										walk.MsgBox(nil, "Error", "请选择生产计划", walk.MsgBoxIconError)
										return
									}
									if util.IsUsbMode {
										util.CheckPorts() //USB的需要重新打开端口，串口的不需要，可以不调用此函数
									}
									util.DoTestAllPortsAllItems()
								},
							},
						},
					},
				},
			},
			Composite{
				Layout: HBox{
					Alignment: AlignHNearVNear,
					Margins:   Margins{Left: 0, Top: 0, Right: 0, Bottom: 0},
					Spacing:   5,
				},
				Children: []Widget{
					GroupBox{
						Title:   "屏蔽COM口",
						Font:    Font{PointSize: viceFontSize, Family: fontFamily},
						Layout:  Grid{Columns: 1},
						MinSize: Size{Width: 150, Height: 100},
						MaxSize: Size{Width: 300, Height: 100},
						Children: []Widget{
							TextEdit{
								AssignTo: &blockedCom,
								Text:     conf.BlockedCom,
								MinSize:  Size{Width: 150, Height: 25},
								MaxSize:  Size{Width: 300, Height: 25},
							},
							PushButton{
								Text:    "确定",
								Font:    Font{PointSize: 9, Family: fontFamily},
								MinSize: Size{Width: 80, Height: 25},
								MaxSize: Size{Width: 300, Height: 25},
								OnClicked: func() {
									var sliPortIdx []string
									var sliPortName []string
									if blockedCom.Text() != "" {
										sliPortIdx = strings.Split(blockedCom.Text(), ",")
									}
									for _, idx := range sliPortIdx {
										portName := fmt.Sprintf("COM%v", idx)
										sliPortName = append(sliPortName, portName)
									}
									util.RefreshPorts(sliPortName)
									util.RefreshTableModel()
									conf.BlockedCom = blockedCom.Text()
									conf.SyncConf()
								},
							},
							PushButton{
								Text:    "刷新端口",
								Font:    Font{PointSize: 9, Family: fontFamily},
								MinSize: Size{Width: 80, Height: 25},
								MaxSize: Size{Width: 300, Height: 25},
								OnClicked: func() {
									util.OpenAllPorts()
									util.RefreshTableModel()
								},
							},
						},
					},
					GroupBox{
						Title:   "GSM信号值",
						Font:    Font{PointSize: viceFontSize, Family: fontFamily},
						Layout:  Grid{Columns: 2},
						MinSize: Size{Width: 70, Height: 100},
						MaxSize: Size{Width: 150, Height: 100},
						Children: []Widget{
							Label{
								Text:    "最大值:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 35, Height: 25},
								MaxSize: Size{Width: 50, Height: 25},
							},
							TextEdit{
								Text:     "32",
								AssignTo: &checkSignalMax,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 35, Height: 25},
								MaxSize:  Size{Width: 150, Height: 25},
								ReadOnly: true,
							},
							Label{
								Text:    "最小值:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 35, Height: 25},
								MaxSize: Size{Width: 50, Height: 25},
							},
							TextEdit{
								Text:     "14",
								AssignTo: &checkSignalMin,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 35, Height: 25},
								MaxSize:  Size{Width: 150, Height: 25},
								ReadOnly: true,
							},
						},
					},
					GroupBox{
						Title:   "测试阈值",
						Font:    Font{PointSize: viceFontSize, Family: fontFamily},
						Layout:  Grid{Columns: 2},
						MinSize: Size{Width: 70, Height: 100},
						MaxSize: Size{Width: 150, Height: 120},
						Children: []Widget{
							Label{
								Text:    "GPS:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 35, Height: 25},
								MaxSize: Size{Width: 40, Height: 25},
							},
							TextEdit{
								Text:     "4",
								AssignTo: &checkGpsValue,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 35, Height: 25},
								MaxSize:  Size{Width: 150, Height: 25},
								ReadOnly: true,
							},
							Label{
								Text:    "WIFI:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 35, Height: 25},
								MaxSize: Size{Width: 40, Height: 25},
							},
							TextEdit{
								Text:     "1",
								AssignTo: &checkWifiValue,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 35, Height: 25},
								MaxSize:  Size{Width: 150, Height: 25},
								ReadOnly: true,
							},
							Label{
								Text:    "电量:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 35, Height: 25},
								MaxSize: Size{Width: 40, Height: 25},
							},
							TextEdit{
								Text:     "1",
								AssignTo: &checkPowerValue,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 35, Height: 25},
								MaxSize:  Size{Width: 150, Height: 25},
								ReadOnly: true,
							},
						},
					},
					GroupBox{
						Title:   "IP地址, 端口",
						Font:    Font{PointSize: viceFontSize, Family: fontFamily},
						Layout:  Grid{Columns: 2},
						MinSize: Size{Width: 180, Height: 100},
						MaxSize: Size{Width: 300, Height: 100},
						Children: []Widget{
							Label{
								Text:    "IP:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 30, Height: 25},
								MaxSize: Size{Width: 40, Height: 25},
							},
							TextEdit{
								Text:     Bind("MainIp"),
								AssignTo: &modifyIp,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 120, Height: 25},
								MaxSize:  Size{Width: 250, Height: 25},
								ReadOnly: true,
							},
							Label{
								Text:    "APN:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 30, Height: 25},
								MaxSize: Size{Width: 50, Height: 25},
							},
							TextEdit{
								Text:     Bind("ApnValue"),
								AssignTo: &modifyApn,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 30, Height: 25},
								MaxSize:  Size{Width: 150, Height: 25},
								ReadOnly: true,
							},
							Label{
								Text:    "协议:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 30, Height: 25},
								MaxSize: Size{Width: 50, Height: 25},
							},
							TextEdit{
								Text:     Bind("ProtocolValue"),
								AssignTo: &modifyProtocol,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 30, Height: 25},
								MaxSize:  Size{Width: 150, Height: 25},
								ReadOnly: true,
							},
						},
					},
					Composite{
						Layout:  VBox{},
						MinSize: Size{Width: 250, Height: 100},
						MaxSize: Size{Width: 500, Height: 100},
						Children: []Widget{
							Composite{
								Layout: Grid{Columns: 2},
								Children: []Widget{
									TextEdit{
										AssignTo:  &compareVersion,
										Text:      "",
										Font:      Font{PointSize: viceFontSize, Family: fontFamily},
										MinSize:   Size{Width: 180, Height: 23},
										MaxSize:   Size{Width: 300, Height: 23},
										MaxLength: 100,
										OnTextChanged: func() {
											util.CompareVersion = compareVersion.Text()
										},
									},
									PushButton{
										Text:    "添加比对版本",
										Font:    Font{PointSize: viceFontSize, Family: fontFamily},
										MinSize: Size{Width: 100, Height: 23},
										MaxSize: Size{Width: 100, Height: 23},
										OnClicked: func() {
											if tv.CurrentIndex() < 0 || tv.CurrentIndex() > tv.Model().(*util.MyTableModel).RowCount() {
												log.Infof("row %v invalid", tv.CurrentIndex())
												return
											}
											item := util.GetModelItems()[tv.CurrentIndex()]
											compareVersion.SetText(strings.Trim(item.Version, "(匹配失败)"))
										},
									},
									TextEdit{
										AssignTo: &compareMainIp,
										Text:     "",
										Font:     Font{PointSize: viceFontSize, Family: fontFamily},
										MinSize:  Size{Width: 180, Height: 23},
										MaxSize:  Size{Width: 300, Height: 23},
										OnTextChanged: func() {
											util.CompareMainIp = compareMainIp.Text()
										},
									},
									PushButton{
										Text:    "添加IP比对",
										Font:    Font{PointSize: viceFontSize, Family: fontFamily},
										MinSize: Size{Width: 100, Height: 23},
										MaxSize: Size{Width: 100, Height: 23},
										OnClicked: func() {
											if tv.CurrentIndex() < 0 || tv.CurrentIndex() > tv.Model().(*util.MyTableModel).RowCount() {
												log.Infof("row %v invalid", tv.CurrentIndex())
												return
											}
											item := util.GetModelItems()[tv.CurrentIndex()]
											if tv.Columns().ByName("MainIp").Visible() {
												compareMainIp.SetText(strings.Trim(item.MainIp, "(匹配失败)"))
											}
										},
									},
									TextEdit{
										AssignTo: &compareViceIp,
										Text:     "",
										Font:     Font{PointSize: viceFontSize, Family: fontFamily},
										MinSize:  Size{Width: 180, Height: 23},
										MaxSize:  Size{Width: 300, Height: 23},
										OnTextChanged: func() {
											util.CompareViceIp = compareViceIp.Text()
										},
									},
									PushButton{
										Text:    "添加副IP比对",
										Font:    Font{PointSize: viceFontSize, Family: fontFamily},
										MinSize: Size{Width: 100, Height: 23},
										MaxSize: Size{Width: 100, Height: 23},
										OnClicked: func() {
											if tv.CurrentIndex() < 0 || tv.CurrentIndex() > tv.Model().(*util.MyTableModel).RowCount() {
												log.Infof("row %v invalid", tv.CurrentIndex())
												return
											}
											item := util.GetModelItems()[tv.CurrentIndex()]
											if tv.Columns().ByName("ViceIp").Visible() {
												compareViceIp.SetText(strings.Trim(item.ViceIp, "(匹配失败)"))
											}
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}.Run()
}

func main() {
	log.Info("start")
	walk.AppendToWalkInit(func() {
		walk.FocusEffect, _ = walk.NewBorderGlowEffect(walk.RGB(0, 63, 255))
		walk.InteractionEffect, _ = walk.NewDropShadowEffect(walk.RGB(63, 63, 63))
		walk.ValidationErrorEffect, _ = walk.NewBorderGlowEffect(walk.RGB(255, 0, 0))
	})

	// 显示登录对话框
	loginResult := dialog.ShowLoginDialog()

	// 检查登录结果
	if !loginResult.Success {
		fmt.Println("登录失败:", loginResult.Message)
		os.Exit(1)
	}

	runMainWindow()
}
