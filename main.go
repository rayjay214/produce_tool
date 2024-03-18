// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	log "github.com/sirupsen/logrus"
	"os"
	"produce_tool/conf"
	"produce_tool/db"
	"produce_tool/model"
	"produce_tool/util"
	"reflect"
	"strings"
	"time"
)

var tv *walk.TableView
var tableColumns []TableViewColumn
var tableModel *util.MyTableModel
var singleFunctionButtons []*walk.PushButton
var singleFunctionMapping map[int]string
var selectedCb *walk.ComboBox
var onnKeyTestBtn *walk.PushButton

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

// 阈值校验
var checkSignalMin *walk.TextEdit
var checkSignalMax *walk.TextEdit
var checkGpsValue *walk.TextEdit
var checkWifiValue *walk.TextEdit

// 比较校验
var compareVersion *walk.TextEdit
var compareMainIp *walk.TextEdit
var compareViceIp *walk.TextEdit

// 待修改值
var modifyIp *walk.TextEdit
var modifyPort *walk.TextEdit

// 通过数量
var passedCnt *walk.LineEdit

func init() {
	initLog()
	initTableColumns()
	initSingleFunctionButtons()
	initRefreshTimer()
	initSyncConfTimer()
	initConf()
	db.InitMysql()
	model.LoadDeviceType()
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
				if passedCnt != nil {
					passedCnt.SetText(fmt.Sprintf("当前测试通过数量:%v", conf.PassedCnt))
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
	if util.ContainsOne(rv.FieldByName(propertyName).String(), "失败", "超时", "已过站") {
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
	modifyIp.SetText(selectedType.MainIp)
	modifyPort.SetText(selectedType.MainPort)
	checkSignalMax.SetText(selectedType.SignalMax)
	checkSignalMin.SetText(selectedType.SignalMin)
	checkGpsValue.SetText(selectedType.GpsMin)
	checkWifiValue.SetText(selectedType.WifiMin)
	util.SelectedDeviceType = selectedType
	util.SyncTestItems()
	util.RefreshTableModel()

	for i := 0; i < tv.Columns().Len(); i++ {
		tv.Columns().At(i).SetVisible(true)
	}

	//stupid method for the moment
	if selectedType.SignalOpen <= 0 {
		tv.Columns().ByName("Signal").SetVisible(false)
	}
	if selectedType.GpsOpen <= 0 {
		tv.Columns().ByName("Gps").SetVisible(false)
	}
	if selectedType.WifiOpen <= 0 {
		tv.Columns().ByName("Wifi").SetVisible(false)
	}
	if selectedType.SnOpen <= 0 {
		tv.Columns().ByName("Sn").SetVisible(false)
	}
	if selectedType.SimOpen <= 0 {
		tv.Columns().ByName("Sim").SetVisible(false)
	}
	if selectedType.ImeiOpen <= 0 {
		tv.Columns().ByName("Imei").SetVisible(false)
	}
	if selectedType.LightOpen <= 0 {
		tv.Columns().ByName("Light").SetVisible(false)
	}
	if selectedType.GsensorOpen <= 0 {
		tv.Columns().ByName("Gsensor").SetVisible(false)
	}
	if selectedType.SetTypeOpen <= 0 {
		tv.Columns().ByName("SetType").SetVisible(false)
	}
	if selectedType.MainIpReadOpen <= 0 {
		tv.Columns().ByName("MainIp").SetVisible(false)
	}
	if selectedType.ViceIpReadOpen <= 0 {
		tv.Columns().ByName("ViceIp").SetVisible(false)
	}
	if selectedType.ViceIpWriteOpen <= 0 {
		tv.Columns().ByName("ViceIpWrite").SetVisible(false)
	}
	if selectedType.PowerOpen <= 0 {
		tv.Columns().ByName("Power").SetVisible(false)
	}
	/*
		if selectedType.DialOpen <= 0 {
			tv.Columns().ByName("Dial").SetVisible(false)
		}

		if selectedType.EndDialOpen <= 0 {
			tv.Columns().ByName("EndDial").SetVisible(false)
		}
		if selectedType.TamperOpen <= 0 {
			tv.Columns().ByName("Tamper").SetVisible(false)
		}
		if selectedType.ApnWriteOpen <= 0 {
			tv.Columns().ByName("Apn").SetVisible(false)
		}
	*/

	conf.SelectedType = selectedType.DeviceType
	conf.SyncConf()
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
		Title:    "生产测试工具",
		Font:     Font{PointSize: viceFontSize, Family: fontFamily},
		Size:     Size{Width: 900, Height: 650},
		Layout:   VBox{Alignment: AlignHNearVNear},
		OnSizeChanged: func() {

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
								AssignTo:  &passedCnt,
								Text:      fmt.Sprintf("当前测试通过数量:%v", conf.PassedCnt),
								Font:      Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:   Size{Width: 100, Height: 25},
								MaxSize:   Size{Width: 150, Height: 25},
								Alignment: AlignHNearVCenter,
								TextColor: walk.RGB(255, 0, 0),
								ReadOnly:  true,
							},
							PushButton{
								Text:    "清零",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 20, Height: 25},
								MaxSize: Size{Width: 80, Height: 25},
								OnClicked: func() {
									conf.CntMutex.Lock()
									conf.PassedCnt = 0
									conf.CntMutex.Unlock()
									passedCnt.SetText(fmt.Sprintf("当前测试通过数量:%v", conf.PassedCnt))
								},
							},
						},
					},
					PushButton{
						Text:      "管理型号",
						Font:      Font{PointSize: 14, Family: fontFamily},
						Alignment: AlignHNearVCenter,
						MinSize:   Size{Width: 60, Height: 100},
						MaxSize:   Size{Width: 100, Height: 100},
						OnClicked: func() {
							model.RunCheckPwdDialog(mw, selectedCb)
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
							},
						},
					},
					GroupBox{
						MinSize: Size{Width: 80, Height: 100},
						MaxSize: Size{Width: 150, Height: 100},
						Font:    Font{PointSize: 12, Family: fontFamily},
						Title:   "测完是否关机:",
						Layout:  HBox{},
						Children: []Widget{
							CheckBox{
								AssignTo: &checkPowerOff,
								Text:     "关机",
								Font:     Font{PointSize: 10, Family: fontFamily},
								MinSize:  Size{Width: 30, Height: 25},
								MaxSize:  Size{Width: 50, Height: 25},
								OnCheckedChanged: func() {
									util.PoweroffAfterTest = checkPowerOff.Checked()
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
											enabled := checkSn.Enabled()
											checkSn.SetEnabled(!enabled)
											checkImei.SetEnabled(!enabled)
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
											if key == walk.KeyReturn {
												if tv.CurrentIndex() < 0 || tv.CurrentIndex() > tv.Model().(*util.MyTableModel).RowCount() {
													log.Infof("row %v invalid", tv.CurrentIndex())
													tv.SetCurrentIndex(0)
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
									//util.CheckPorts() //USB的需要重新打开端口，串口的不需要，可以不调用此函数
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
						MinSize: Size{Width: 150, Height: 80},
						MaxSize: Size{Width: 300, Height: 80},
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
						},
					},
					GroupBox{
						Title:   "GSM信号值",
						Font:    Font{PointSize: viceFontSize, Family: fontFamily},
						Layout:  Grid{Columns: 2},
						MinSize: Size{Width: 70, Height: 80},
						MaxSize: Size{Width: 150, Height: 80},
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
						Title:   "定位信号值",
						Font:    Font{PointSize: viceFontSize, Family: fontFamily},
						Layout:  Grid{Columns: 2},
						MinSize: Size{Width: 70, Height: 80},
						MaxSize: Size{Width: 150, Height: 80},
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
						},
					},
					GroupBox{
						Title:   "IP地址, 端口",
						Font:    Font{PointSize: viceFontSize, Family: fontFamily},
						Layout:  Grid{Columns: 2},
						MinSize: Size{Width: 180, Height: 80},
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
								Text:    "端口:",
								Font:    Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize: Size{Width: 30, Height: 25},
								MaxSize: Size{Width: 50, Height: 25},
							},
							TextEdit{
								Text:     Bind("MainPort"),
								AssignTo: &modifyPort,
								Font:     Font{PointSize: viceFontSize, Family: fontFamily},
								MinSize:  Size{Width: 30, Height: 25},
								MaxSize:  Size{Width: 150, Height: 25},
								ReadOnly: true,
							},
						},
					},
					Composite{
						Layout:  VBox{},
						MaxSize: Size{Width: 500, Height: 100},
						MinSize: Size{Width: 250, Height: 80},
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

	runMainWindow()
}
