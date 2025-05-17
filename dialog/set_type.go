package dialog

import (
	"fmt"
	"produce_tool/model"
	"reflect"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func KnownTypesString() []string {
	mptypes := model.GetDeviceTypes()
	types := make([]string, 0)
	for key, _ := range mptypes {
		types = append(types, key)
	}
	return types
}

func KnownTypes() []model.DeviceTypeInfo {
	mptypes := model.GetDeviceTypes()
	types := make([]model.DeviceTypeInfo, 0)
	for _, value := range mptypes {
		types = append(types, value)
	}
	return types
}

func CopyStruct(src interface{}, dest interface{}) {
	valSrc := reflect.ValueOf(src)
	valDest := reflect.ValueOf(dest)

	if valSrc.Kind() != reflect.Ptr || valDest.Kind() != reflect.Ptr {
		fmt.Println("Both src and dest must be pointers to structs.")
		return
	}

	valSrc = valSrc.Elem() // 取得指针指向的结构体
	valDest = valDest.Elem()

	for i := 0; i < valSrc.NumField(); i++ {
		fieldSrc := valSrc.Field(i)
		fieldDest := valDest.FieldByName(valSrc.Type().Field(i).Name)

		if fieldDest.IsValid() && fieldDest.CanSet() {
			fieldDest.Set(fieldSrc)
		}
	}
}

func RunCheckPwdDialog(owner walk.Form, selectedCb *walk.ComboBox) (int, error) {
	var dlg *walk.Dialog
	var pwd *walk.LineEdit
	return Dialog{
		AssignTo: &dlg,
		Title:    "请输入密码",
		MinSize:  Size{200, 100},
		Layout: VBox{
			Alignment: AlignHCenterVCenter,
		},
		Children: []Widget{
			LineEdit{
				AssignTo:     &pwd,
				PasswordMode: true,
				OnKeyPress: func(key walk.Key) {
					if key == walk.KeyReturn {
						if pwd.Text() == "10086" {
							dlg.Cancel()
							RunDialogAddType(owner, selectedCb)
						} else {
							dlg.Cancel()
						}
					}
				},
			},
		},
	}.Run(owner)
}

func RunDialogAddType(owner walk.Form, selectedCb *walk.ComboBox) (int, error) {

	deviceType := new(model.DeviceTypeInfo)
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton
	var selected *walk.ComboBox

	types := model.GetDeviceTypes()

	return Dialog{
		AssignTo:      &dlg,
		Title:         "添加型号",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize:       Size{750, 600},

		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "deviceType",
			DataSource:     deviceType,
			ErrorPresenter: ToolTipErrorPresenter{},
		},

		Layout: VBox{
			Alignment: AlignHNearVNear,
		},
		Children: []Widget{
			Composite{
				Layout: HBox{
					Alignment: AlignHCenterVNear,
					Margins: Margins{
						Left:   0,
						Top:    0,
						Right:  0,
						Bottom: 20,
					},
				},
				MinSize: Size{Width: 100, Height: 30},
				MaxSize: Size{Width: 700, Height: 70},
				Children: []Widget{
					GroupBox{
						Title:   "已存在设备",
						Layout:  HBox{},
						MinSize: Size{Width: 100, Height: 40},
						MaxSize: Size{Width: 120, Height: 50},
						Children: []Widget{
							ComboBox{
								AssignTo: &selected,
								//Model:         KnownTypes(),
								Model:         model.AllTypes,
								BindingMember: "DeviceType",
								DisplayMember: "DeviceType",
								MinSize:       Size{Width: 90, Height: 40},
								MaxSize:       Size{Width: 120, Height: 50},
								OnCurrentIndexChanged: func() {
									//strType := selected.Model().([]string)[selected.CurrentIndex()]
									//selectedTypeInfo := GetDeviceTypes()[strType]
									selectedTypeInfo := selected.Model().([]model.DeviceTypeInfo)[selected.CurrentIndex()]
									//StupidCopy(selectedTypeInfo, deviceType)
									CopyStruct(&selectedTypeInfo, deviceType)
									db.Reset()
								},
							},
						},
					},
				},
			},
			HSplitter{
				Children: []Widget{
					Composite{
						Alignment: AlignHCenterVNear,
						Layout: Grid{
							Columns: 2,
							Spacing: 10,
						},
						Children: []Widget{
							RadioButtonGroupBox{
								MinSize:    Size{Width: 160, Height: 50},
								MaxSize:    Size{Width: 160, Height: 50},
								Title:      "超速报警",
								Layout:     HBox{},
								DataMember: "OverSpeedAlarm",
								Buttons: []RadioButton{
									{
										Text:    "支持",
										Value:   1,
										MinSize: Size{Width: 60, Height: 30},
										MaxSize: Size{Width: 70, Height: 30},
									},
									{
										Text:    "不支持",
										Value:   0,
										MinSize: Size{Width: 60, Height: 30},
										MaxSize: Size{Width: 70, Height: 30},
									},
								},
							},
							RadioButtonGroupBox{
								MinSize:    Size{Width: 160, Height: 50},
								MaxSize:    Size{Width: 160, Height: 50},
								Title:      "监听",
								Layout:     HBox{},
								DataMember: "Listen",
								Buttons: []RadioButton{
									{
										Text:    "支持",
										Value:   1,
										MinSize: Size{Width: 60, Height: 30},
										MaxSize: Size{Width: 70, Height: 30},
									},
									{
										Text:    "不支持",
										Value:   0,
										MinSize: Size{Width: 60, Height: 30},
										MaxSize: Size{Width: 70, Height: 30},
									},
								},
							},
							RadioButtonGroupBox{
								MinSize:    Size{Width: 160, Height: 50},
								MaxSize:    Size{Width: 160, Height: 50},
								Title:      "支持灯控",
								Layout:     HBox{},
								DataMember: "LightControl",
								Buttons: []RadioButton{
									{
										Text:    "支持",
										Value:   1,
										MinSize: Size{Width: 60, Height: 30},
										MaxSize: Size{Width: 70, Height: 30},
									},
									{
										Text:    "不支持",
										Value:   0,
										MinSize: Size{Width: 60, Height: 30},
										MaxSize: Size{Width: 70, Height: 30},
									},
								},
							},
							RadioButtonGroupBox{
								MinSize:    Size{Width: 160, Height: 50},
								MaxSize:    Size{Width: 160, Height: 50},
								Title:      "防拆报警",
								Layout:     HBox{},
								DataMember: "TamperAlarm",
								Buttons: []RadioButton{
									{
										Text:    "支持",
										Value:   1,
										MinSize: Size{Width: 60, Height: 30},
										MaxSize: Size{Width: 70, Height: 30},
									},
									{
										Text:    "不支持",
										Value:   0,
										MinSize: Size{Width: 60, Height: 30},
										MaxSize: Size{Width: 70, Height: 30},
									},
								},
							},
							RadioButtonGroupBox{
								MinSize:    Size{Width: 160, Height: 50},
								MaxSize:    Size{Width: 160, Height: 50},
								Title:      "震动报警",
								Layout:     HBox{},
								DataMember: "ShakeAlarm",
								Buttons: []RadioButton{
									{
										Text:    "支持",
										Value:   1,
										MinSize: Size{Width: 60, Height: 30},
										MaxSize: Size{Width: 70, Height: 30},
									},
									{
										Text:    "不支持",
										Value:   0,
										MinSize: Size{Width: 60, Height: 30},
										MaxSize: Size{Width: 70, Height: 30},
									},
								},
							},
							RadioButtonGroupBox{
								MinSize:    Size{Width: 160, Height: 50},
								MaxSize:    Size{Width: 160, Height: 50},
								Title:      "录音",
								Layout:     HBox{},
								DataMember: "Recording",
								Buttons: []RadioButton{
									{
										Text:    "支持",
										Value:   1,
										MinSize: Size{Width: 60, Height: 30},
										MaxSize: Size{Width: 70, Height: 30},
									},
									{
										Text:    "不支持",
										Value:   0,
										MinSize: Size{Width: 60, Height: 30},
										MaxSize: Size{Width: 70, Height: 30},
									},
								},
							},
							RadioButtonGroupBox{
								MinSize:    Size{Width: 160, Height: 50},
								MaxSize:    Size{Width: 160, Height: 50},
								Title:      "低电报警",
								Layout:     HBox{},
								DataMember: "LowpowerAlarm",
								Buttons: []RadioButton{
									{
										Text:    "支持",
										Value:   1,
										MinSize: Size{Width: 60, Height: 30},
										MaxSize: Size{Width: 70, Height: 30},
									},
									{
										Text:    "不支持",
										Value:   0,
										MinSize: Size{Width: 60, Height: 30},
										MaxSize: Size{Width: 70, Height: 30},
									},
								},
							},
							RadioButtonGroupBox{
								MinSize:    Size{Width: 160, Height: 50},
								MaxSize:    Size{Width: 160, Height: 50},
								Title:      "急转弯报警",
								Layout:     HBox{},
								DataMember: "SharpTurnAlarm",
								Buttons: []RadioButton{
									{
										Text:    "支持",
										Value:   1,
										MinSize: Size{Width: 60, Height: 30},
										MaxSize: Size{Width: 70, Height: 30},
									},
									{
										Text:    "不支持",
										Value:   0,
										MinSize: Size{Width: 60, Height: 30},
										MaxSize: Size{Width: 70, Height: 30},
									},
								},
							},
							RadioButtonGroupBox{
								MinSize:    Size{Width: 160, Height: 50},
								MaxSize:    Size{Width: 160, Height: 50},
								Title:      "急加速报警",
								Layout:     HBox{},
								DataMember: "RapidAccleAlarm",
								Buttons: []RadioButton{
									{
										Text:    "支持",
										Value:   1,
										MinSize: Size{Width: 60, Height: 30},
										MaxSize: Size{Width: 70, Height: 30},
									},
									{
										Text:    "不支持",
										Value:   0,
										MinSize: Size{Width: 60, Height: 30},
										MaxSize: Size{Width: 70, Height: 30},
									},
								},
							},
							RadioButtonGroupBox{
								MinSize:    Size{Width: 160, Height: 50},
								MaxSize:    Size{Width: 160, Height: 50},
								Title:      "急减速报警",
								Layout:     HBox{},
								DataMember: "RapidDecleAlarm",
								Buttons: []RadioButton{
									{
										Text:    "支持",
										Value:   1,
										MinSize: Size{Width: 60, Height: 30},
										MaxSize: Size{Width: 70, Height: 30},
									},
									{
										Text:    "不支持",
										Value:   0,
										MinSize: Size{Width: 60, Height: 30},
										MaxSize: Size{Width: 70, Height: 30},
									},
								},
							},
							GroupBox{
								Row:     5,
								Column:  0,
								Title:   "主IP端口",
								Layout:  Grid{Columns: 2},
								MinSize: Size{Width: 140, Height: 80},
								MaxSize: Size{Width: 300, Height: 100},
								Children: []Widget{
									Label{
										Text:    "IP:",
										MinSize: Size{Width: 30, Height: 25},
										MaxSize: Size{Width: 40, Height: 25},
									},
									LineEdit{
										Text:    Bind("MainIp"),
										MinSize: Size{Width: 50, Height: 25},
										MaxSize: Size{Width: 150, Height: 25},
									},
									Label{
										Text:    "端口:",
										MinSize: Size{Width: 30, Height: 25},
										MaxSize: Size{Width: 50, Height: 25},
									},
									LineEdit{
										Text:    Bind("MainPort"),
										MinSize: Size{Width: 30, Height: 25},
										MaxSize: Size{Width: 150, Height: 25},
									},
								},
							},
							GroupBox{
								Row:     5,
								Column:  1,
								Title:   "副IP端口",
								Layout:  Grid{Columns: 2},
								MinSize: Size{Width: 140, Height: 80},
								MaxSize: Size{Width: 300, Height: 100},
								Children: []Widget{
									Label{
										Text:    "IP:",
										MinSize: Size{Width: 30, Height: 25},
										MaxSize: Size{Width: 40, Height: 25},
									},
									LineEdit{
										Text:    Bind("ViceIp"),
										MinSize: Size{Width: 50, Height: 25},
										MaxSize: Size{Width: 150, Height: 25},
									},
									Label{
										Text:    "端口:",
										MinSize: Size{Width: 30, Height: 25},
										MaxSize: Size{Width: 50, Height: 25},
									},
									LineEdit{
										Text:    Bind("VicePort"),
										MinSize: Size{Width: 30, Height: 25},
										MaxSize: Size{Width: 150, Height: 25},
									},
								},
							},
							GroupBox{
								MinSize: Size{Width: 160, Height: 50},
								MaxSize: Size{Width: 160, Height: 50},
								Title:   "APN",
								Layout:  HBox{},
								Children: []Widget{
									LineEdit{
										Text:    Bind("APN"),
										MinSize: Size{Width: 100, Height: 30},
										MaxSize: Size{Width: 120, Height: 30},
									},
								},
							},
							GroupBox{
								MinSize: Size{Width: 160, Height: 50},
								MaxSize: Size{Width: 160, Height: 50},
								Title:   "SN校验",
								Layout:  HBox{},
								Children: []Widget{
									LineEdit{
										Text:    Bind("SnLength"),
										MinSize: Size{Width: 100, Height: 30},
										MaxSize: Size{Width: 120, Height: 30},
									},
								},
							},
							GroupBox{
								MinSize: Size{Width: 160, Height: 50},
								MaxSize: Size{Width: 160, Height: 50},
								Title:   "设备型号",
								Layout:  HBox{},
								Children: []Widget{
									LineEdit{
										Text:    Bind("DeviceType", SelRequired{}),
										MinSize: Size{Width: 100, Height: 30},
										MaxSize: Size{Width: 120, Height: 30},
									},
								},
							},
							GroupBox{
								MinSize: Size{Width: 160, Height: 50},
								MaxSize: Size{Width: 160, Height: 50},
								Title:   "电量告警阈值",
								Layout:  HBox{},
								Children: []Widget{
									LineEdit{
										Text:    Bind("PowerMin"),
										MinSize: Size{Width: 100, Height: 30},
										MaxSize: Size{Width: 120, Height: 30},
									},
								},
							},
							GroupBox{
								MinSize: Size{Width: 160, Height: 50},
								MaxSize: Size{Width: 160, Height: 50},
								Title:   "协议",
								Layout:  HBox{},
								Children: []Widget{
									LineEdit{
										Text:    Bind("ProtocolValue"),
										MinSize: Size{Width: 100, Height: 30},
										MaxSize: Size{Width: 120, Height: 30},
									},
								},
							},
						},
					},
					Composite{
						Alignment: AlignHNearVNear,
						Layout: VBox{
							Alignment: AlignHNearVNear,
						},
						Children: []Widget{
							Composite{
								Layout: HBox{
									Alignment: AlignHNearVNear,
									Margins: Margins{
										Left:   0,
										Top:    0,
										Right:  0,
										Bottom: 15,
									},
								},
								Children: []Widget{
									GroupBox{
										MinSize: Size{Width: 150, Height: 100},
										MaxSize: Size{Width: 180, Height: 150},
										Title:   "信号值",
										Layout: Grid{
											Columns: 2,
											Spacing: 10,
										},
										Children: []Widget{
											Label{
												MaxSize: Size{Width: 50, Height: 40},
												Text:    "延时",
											},
											ComboBox{
												Value: Bind("SignalDelay"),
												Model: []string{"0", "1", "2", "3"},
											},
											RadioButtonGroup{
												DataMember: "SignalOpen",
												Buttons: []RadioButton{
													{
														Text:    "开启",
														Value:   1,
														MinSize: Size{Width: 30, Height: 30},
														MaxSize: Size{Width: 40, Height: 30},
													},
													{
														Text:    "关闭",
														Value:   0,
														MinSize: Size{Width: 30, Height: 30},
														MaxSize: Size{Width: 40, Height: 30},
													},
												},
											},
											Label{
												MaxSize: Size{Width: 40, Height: 40},
												Text:    "最小值",
											},
											LineEdit{
												Text: Bind("SignalMin"),
											},
											Label{
												MaxSize: Size{Width: 40, Height: 40},
												Text:    "最大值",
											},
											LineEdit{
												Text: Bind("SignalMax"),
											},
										},
									},
									GroupBox{
										MinSize: Size{Width: 150, Height: 100},
										MaxSize: Size{Width: 180, Height: 150},
										Title:   "GPS卫星值",
										Layout: Grid{
											Alignment: AlignHNearVNear,
											Columns:   2,
											Spacing:   25,
										},
										Children: []Widget{
											Label{
												MaxSize: Size{Width: 30, Height: 30},
												Text:    "延时",
											},
											ComboBox{
												MaxSize: Size{Width: 60},
												Value:   Bind("GpsDelay"),
												Model:   []string{"0", "1", "2", "3"},
											},
											RadioButtonGroup{
												DataMember: "GpsOpen",
												Buttons: []RadioButton{
													{
														Text:    "开启",
														Value:   1,
														MinSize: Size{Width: 30, Height: 30},
														MaxSize: Size{Width: 40, Height: 30},
													},
													{
														Text:    "关闭",
														Value:   0,
														MinSize: Size{Width: 30, Height: 30},
														MaxSize: Size{Width: 40, Height: 30},
													},
												},
											},
											Label{
												MaxSize: Size{Width: 60, Height: 30},
												Text:    "阈值:",
											},
											LineEdit{
												Text: Bind("GpsMin"),
											},
										},
									},
									GroupBox{
										MinSize: Size{Width: 150, Height: 100},
										MaxSize: Size{Width: 180, Height: 150},
										Title:   "WIFI",
										Layout: Grid{
											Alignment: AlignHNearVNear,
											Columns:   2,
											Spacing:   25,
										},
										Children: []Widget{
											Label{
												MaxSize: Size{Width: 30, Height: 30},
												Text:    "延时",
											},
											ComboBox{
												MaxSize: Size{Width: 60},
												Model:   []string{"0", "1", "2", "3"},
											},
											RadioButtonGroup{
												DataMember: "WifiOpen",
												Buttons: []RadioButton{
													{
														Text:    "开启",
														Value:   1,
														MinSize: Size{Width: 30, Height: 30},
														MaxSize: Size{Width: 40, Height: 30},
													},
													{
														Text:    "关闭",
														Value:   0,
														MinSize: Size{Width: 30, Height: 30},
														MaxSize: Size{Width: 40, Height: 30},
													},
												},
											},
											Label{
												MaxSize: Size{Width: 60, Height: 30},
												Text:    "阈值:",
											},
											LineEdit{
												Text: Bind("WifiMin"),
											},
										},
									},
								},
							},
							Composite{
								Layout: Grid{
									Alignment: AlignHNearVNear,
									Margins: Margins{
										Left:   0,
										Top:    0,
										Right:  0,
										Bottom: 30,
									},
									Columns: 3,
								},
								Children: []Widget{
									RadioButtonGroupBox{
										MinSize:    Size{Width: 160, Height: 50},
										MaxSize:    Size{Width: 160, Height: 50},
										Title:      "SN号",
										Layout:     HBox{},
										DataMember: "SnOpen",
										Buttons: []RadioButton{
											{
												Value:   1,
												Text:    "读取",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
											{
												Value:   0,
												Text:    "不读",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
										},
									},
									RadioButtonGroupBox{
										MinSize:    Size{Width: 160, Height: 50},
										MaxSize:    Size{Width: 160, Height: 50},
										Title:      "拨打电话",
										Layout:     HBox{},
										DataMember: "DialOpen",
										Buttons: []RadioButton{
											{
												Value:   1,
												Text:    "开启",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
											{
												Value:   0,
												Text:    "关闭",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
										},
									},
									RadioButtonGroupBox{
										MinSize:    Size{Width: 160, Height: 50},
										MaxSize:    Size{Width: 160, Height: 50},
										Title:      "SIM卡",
										Layout:     HBox{},
										DataMember: "SimOpen",
										Buttons: []RadioButton{
											{
												Value:   1,
												Text:    "开启",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
											{
												Value:   0,
												Text:    "关闭",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
										},
									},
									RadioButtonGroupBox{
										MinSize:    Size{Width: 160, Height: 50},
										MaxSize:    Size{Width: 160, Height: 50},
										Title:      "IMEI号",
										Layout:     HBox{},
										DataMember: "ImeiOpen",
										Buttons: []RadioButton{
											{
												Value:   1,
												Text:    "读取",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
											{
												Value:   0,
												Text:    "不读",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
										},
									},
									RadioButtonGroupBox{
										MinSize:    Size{Width: 160, Height: 50},
										MaxSize:    Size{Width: 160, Height: 50},
										Title:      "光感",
										Layout:     HBox{},
										DataMember: "LightOpen",
										Buttons: []RadioButton{
											{
												Value:   1,
												Text:    "开启",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
											{
												Value:   0,
												Text:    "关闭",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
										},
									},
									RadioButtonGroupBox{
										MinSize:    Size{Width: 160, Height: 50},
										MaxSize:    Size{Width: 160, Height: 50},
										Title:      "重力",
										Layout:     HBox{},
										DataMember: "GsensorOpen",
										Buttons: []RadioButton{
											{
												Value:   1,
												Text:    "开启",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
											{
												Value:   0,
												Text:    "关闭",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
										},
									},
									RadioButtonGroupBox{
										MinSize:    Size{Width: 160, Height: 50},
										MaxSize:    Size{Width: 160, Height: 50},
										Title:      "电量",
										Layout:     HBox{},
										DataMember: "PowerOpen",
										Buttons: []RadioButton{
											{
												Value:   1,
												Text:    "开启",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
											{
												Value:   0,
												Text:    "关闭",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
										},
									},
									RadioButtonGroupBox{
										MinSize:    Size{Width: 160, Height: 50},
										MaxSize:    Size{Width: 160, Height: 50},
										Title:      "挂断",
										Layout:     HBox{},
										DataMember: "EndDialOpen",
										Buttons: []RadioButton{
											{
												Value:   1,
												Text:    "开启",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
											{
												Value:   0,
												Text:    "关闭",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
										},
									},
									RadioButtonGroupBox{
										MinSize:    Size{Width: 160, Height: 50},
										MaxSize:    Size{Width: 160, Height: 50},
										Title:      "防拆",
										Layout:     HBox{},
										DataMember: "TamperOpen",
										Buttons: []RadioButton{
											{
												Value:   1,
												Text:    "开启",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
											{
												Value:   0,
												Text:    "关闭",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
										},
									},
									RadioButtonGroupBox{
										MinSize:    Size{Width: 160, Height: 50},
										MaxSize:    Size{Width: 160, Height: 50},
										Title:      "设置型号",
										Layout:     HBox{},
										DataMember: "SetTypeOpen",
										Buttons: []RadioButton{
											{
												Value:   1,
												Text:    "开启",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
											{
												Value:   0,
												Text:    "关闭",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
										},
									},
									RadioButtonGroupBox{
										MinSize:    Size{Width: 160, Height: 50},
										MaxSize:    Size{Width: 160, Height: 50},
										Title:      "IP地址",
										Layout:     HBox{},
										DataMember: "MainIpReadOpen",
										Buttons: []RadioButton{
											{
												Value:   1,
												Text:    "开启",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
											{
												Value:   0,
												Text:    "关闭",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
										},
									},
									RadioButtonGroupBox{
										MinSize:    Size{Width: 160, Height: 50},
										MaxSize:    Size{Width: 160, Height: 50},
										Title:      "APN写入",
										Layout:     HBox{},
										DataMember: "ApnWriteOpen",
										Buttons: []RadioButton{
											{
												Value:   1,
												Text:    "开启",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
											{
												Value:   0,
												Text:    "关闭",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
										},
									},
									RadioButtonGroupBox{
										MinSize:    Size{Width: 160, Height: 50},
										MaxSize:    Size{Width: 160, Height: 50},
										Title:      "副IP地址",
										Layout:     HBox{},
										DataMember: "ViceIpReadOpen",
										Buttons: []RadioButton{
											{
												Value:   1,
												Text:    "读取",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
											{
												Value:   0,
												Text:    "不读",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
										},
									},
									RadioButtonGroupBox{
										MinSize:    Size{Width: 160, Height: 50},
										MaxSize:    Size{Width: 160, Height: 50},
										Title:      "副IP写入",
										Layout:     HBox{},
										DataMember: "ViceIpWriteOpen",
										Buttons: []RadioButton{
											{
												Value:   1,
												Text:    "开启",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
											{
												Value:   0,
												Text:    "关闭",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
										},
									},
									RadioButtonGroupBox{
										MinSize:    Size{Width: 160, Height: 50},
										MaxSize:    Size{Width: 160, Height: 50},
										Title:      "读取协议",
										Layout:     HBox{},
										DataMember: "ProtocolOpen",
										Buttons: []RadioButton{
											{
												Value:   1,
												Text:    "开启",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
											{
												Value:   0,
												Text:    "关闭",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
										},
									},
									RadioButtonGroupBox{
										MinSize:    Size{Width: 160, Height: 50},
										MaxSize:    Size{Width: 160, Height: 50},
										Title:      "写入协议",
										Layout:     HBox{},
										DataMember: "ProtocolWriteOpen",
										Buttons: []RadioButton{
											{
												Value:   1,
												Text:    "开启",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
											{
												Value:   0,
												Text:    "关闭",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
										},
									},
									RadioButtonGroupBox{
										MinSize:    Size{Width: 160, Height: 50},
										MaxSize:    Size{Width: 160, Height: 50},
										Title:      "单独写入主IP",
										Layout:     HBox{},
										DataMember: "MainIpWriteOpen",
										Buttons: []RadioButton{
											{
												Value:   1,
												Text:    "开启",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
											{
												Value:   0,
												Text:    "关闭",
												MinSize: Size{Width: 60, Height: 30},
												MaxSize: Size{Width: 70, Height: 30},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			Composite{
				Layout: HBox{
					Alignment: AlignHCenterVNear,
				},
				MinSize: Size{Width: 100, Height: 30},
				MaxSize: Size{Width: 700, Height: 70},
				Children: []Widget{
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							/*
								if deviceType.DeviceType == "" {
									walk.MsgBox(nil, "Error", "请填写设备型号", walk.MsgBoxIconError)
									return
								}
							*/
							if err := db.Submit(); err != nil {
								fmt.Println(err)
								return
							} else {
								if deviceType.DeviceType != "" {
									types[deviceType.DeviceType] = *deviceType
									model.SyncDeviceTypes(types)
									types := model.GetDeviceTypes()
									model.AllTypes = make([]model.DeviceTypeInfo, 0)
									for _, value := range types {
										model.AllTypes = append(model.AllTypes, value)
									}
									selectedCb.SetModel(model.AllTypes)
									for idx, t := range model.AllTypes {
										if t.DeviceType == deviceType.DeviceType {
											selectedCb.SetCurrentIndex(idx)
											selectedCb.SetCurrentIndex(idx)
										}
									}
								}
							}

							dlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
}
