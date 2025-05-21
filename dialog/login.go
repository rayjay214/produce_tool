package dialog

import (
	"encoding/json"
	"github.com/lxn/win"
	"io/ioutil"
	"os"
	"path/filepath"
	"produce_tool/model"
	"produce_tool/network"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// LoginResult 登录结果
type LoginResult struct {
	Success bool
	Token   string
	Message string
}

// SavedCredentials 保存的凭据
type SavedCredentials struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	RememberPass bool   `json:"remember_pass"`
}

// 获取凭据文件路径
func getCredentialsFilePath() string {
	// 获取应用程序所在目录
	exePath, err := os.Executable()
	if err != nil {
		return "credentials.json"
	}

	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, "credentials.json")
}

// 保存凭据到文件
func saveCredentials(username, password string, remember bool) error {
	if !remember {
		// 如果不记住密码，则删除凭据文件
		os.Remove(getCredentialsFilePath())
		return nil
	}

	creds := SavedCredentials{
		Username:     username,
		Password:     password,
		RememberPass: remember,
	}

	data, err := json.Marshal(creds)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(getCredentialsFilePath(), data, 0600)
}

// 从文件加载凭据
func loadCredentials() (string, string, bool) {
	data, err := ioutil.ReadFile(getCredentialsFilePath())
	if err != nil {
		return "", "", false
	}

	var creds SavedCredentials
	err = json.Unmarshal(data, &creds)
	if err != nil {
		return "", "", false
	}

	return creds.Username, creds.Password, creds.RememberPass
}

// ShowLoginDialog 显示登录对话框
func ShowLoginDialog() (result LoginResult) {
	var dlg *walk.Dialog
	var usernameEdit, passwordEdit *walk.LineEdit
	var loginStatus *walk.TextLabel
	var acceptPB, cancelPB *walk.PushButton
	var rememberPassCB *walk.CheckBox

	fontFamily := "Microsoft YaHei"
	fontSize := 12

	result = LoginResult{
		Success: false,
		Token:   "",
		Message: "用户取消登录",
	}

	// 加载保存的凭据
	savedUsername, savedPassword, savedRemember := loadCredentials()

	Dialog{
		AssignTo:      &dlg,
		Title:         "系统登录",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize:       Size{Width: 400, Height: 280}, // 增加高度以容纳记住密码选项
		Layout:        VBox{},
		Font:          Font{PointSize: fontSize, Family: fontFamily},
		OnSizeChanged: func() {
			// 获取屏幕尺寸
			screenWidth := int(win.GetSystemMetrics(win.SM_CXSCREEN))
			screenHeight := int(win.GetSystemMetrics(win.SM_CYSCREEN))

			// 获取对话框尺寸
			bounds := dlg.Bounds()

			// 计算居中位置
			x := (screenWidth - bounds.Width) / 2
			y := (screenHeight - bounds.Height) / 2

			// 设置对话框位置
			dlg.SetBounds(walk.Rectangle{X: x, Y: y, Width: bounds.Width, Height: bounds.Height})
		},
		Children: []Widget{
			Composite{
				Layout: VBox{MarginsZero: false, Margins: Margins{Top: 20, Left: 50, Right: 50}},
				Children: []Widget{
					// 标题
					TextLabel{
						Text:      "系统登录",
						Alignment: AlignHCenterVCenter,
						Font:      Font{PointSize: fontSize + 4, Family: fontFamily, Bold: true},
						MinSize:   Size{Height: 40},
					},
					// 用户名输入框
					Composite{
						Layout: HBox{MarginsZero: false, Margins: Margins{Top: 10}},
						Children: []Widget{
							TextLabel{
								Text:    "用户名：",
								MinSize: Size{Width: 80},
							},
							LineEdit{
								AssignTo: &usernameEdit,
								MinSize:  Size{Width: 200},
								Text:     savedUsername, // 填充保存的用户名
							},
						},
					},
					// 密码输入框
					Composite{
						Layout: HBox{MarginsZero: false, Margins: Margins{Top: 10}},
						Children: []Widget{
							TextLabel{
								Text:    "密码：",
								MinSize: Size{Width: 80},
							},
							LineEdit{
								AssignTo:     &passwordEdit,
								MinSize:      Size{Width: 200},
								PasswordMode: true,
								Text:         savedPassword, // 填充保存的密码
							},
						},
					},
					CheckBox{
						Alignment: AlignHNearVNear,
						AssignTo:  &rememberPassCB,
						Text:      "记住密码",
						Checked:   savedRemember,
					},
					// 状态显示
					TextLabel{
						AssignTo:  &loginStatus,
						Text:      "",
						Alignment: AlignHCenterVCenter,
						MinSize:   Size{Height: 30},
						Font:      Font{PointSize: fontSize, Family: fontFamily},
					},
					// 按钮区域
					Composite{
						Layout: HBox{MarginsZero: false, Margins: Margins{Top: 20}, Alignment: AlignHCenterVCenter},
						Children: []Widget{
							PushButton{
								AssignTo: &acceptPB,
								Text:     "登录",
								MinSize:  Size{Width: 100, Height: 30},
								OnClicked: func() {
									// 获取用户名和密码
									username := usernameEdit.Text()
									password := passwordEdit.Text()

									// 验证输入
									if username == "" || password == "" {
										loginStatus.SetText("用户名和密码不能为空")
										return
									}

									// 调用登录API
									loginStatus.SetText("正在登录...")
									success, token, message := network.DoLogin(username, password)

									if success {
										// 保存凭据（如果勾选了记住密码）
										saveCredentials(username, password, rememberPassCB.Checked())

										result.Success = true
										result.Token = token
										result.Message = "登录成功"
										dlg.Synchronize(func() {
											model.LoadDeviceTypeNetwork()
											model.LoadPlanNetwork()
											network.DoGetUserInfo()
											dlg.Accept()
										})
									} else {
										loginStatus.SetText("登录失败: " + message)
									}
								},
							},
							PushButton{
								AssignTo: &cancelPB,
								Text:     "取消",
								MinSize:  Size{Width: 100, Height: 30},
								OnClicked: func() {
									dlg.Cancel()
								},
							},
						},
					},
				},
			},
		},
	}.Run(nil)

	return result
}
