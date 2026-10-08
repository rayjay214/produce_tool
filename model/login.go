package model

import (
	"encoding/json"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"io/ioutil"
	"os"
	"path/filepath"
	"produce_tool/network"
)

type LoginResult struct {
	Success bool
	Token   string
	Message string
}

type SavedCredentials struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	RememberPass bool   `json:"remember_pass"`
}

func getCredentialsFilePath() string {
	exePath, err := os.Executable()
	if err != nil {
		return "credentials.json"
	}

	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, "credentials.json")
}

func saveCredentials(username, password string, remember bool) error {
	if !remember {
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

	savedUsername, savedPassword, savedRemember := loadCredentials()

	Dialog{
		AssignTo:      &dlg,
		Title:         "系统登录",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize:       Size{Width: 400, Height: 280},
		Layout:        VBox{},
		Font:          Font{PointSize: fontSize, Family: fontFamily},
		OnSizeChanged: func() {
			screenWidth := int(win.GetSystemMetrics(win.SM_CXSCREEN))
			screenHeight := int(win.GetSystemMetrics(win.SM_CYSCREEN))

			bounds := dlg.Bounds()

			x := (screenWidth - bounds.Width) / 2
			y := (screenHeight - bounds.Height) / 2

			dlg.SetBounds(walk.Rectangle{X: x, Y: y, Width: bounds.Width, Height: bounds.Height})
		},
		Children: []Widget{
			Composite{
				Layout: VBox{MarginsZero: false, Margins: Margins{Top: 20, Left: 50, Right: 50}},
				Children: []Widget{
					TextLabel{
						Text:      "系统登录",
						Alignment: AlignHCenterVCenter,
						Font:      Font{PointSize: fontSize + 4, Family: fontFamily, Bold: true},
						MinSize:   Size{Height: 40},
					},
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
								Text:     savedUsername,
							},
						},
					},
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
								Text:         savedPassword,
							},
						},
					},
					CheckBox{
						Alignment: AlignHNearVNear,
						AssignTo:  &rememberPassCB,
						Text:      "记住密码",
						Checked:   savedRemember,
					},
					TextLabel{
						AssignTo:  &loginStatus,
						Text:      "",
						Alignment: AlignHCenterVCenter,
						MinSize:   Size{Height: 30},
						Font:      Font{PointSize: fontSize, Family: fontFamily},
					},
					Composite{
						Layout: HBox{MarginsZero: false, Margins: Margins{Top: 20}, Alignment: AlignHCenterVCenter},
						Children: []Widget{
							PushButton{
								AssignTo: &acceptPB,
								Text:     "登录",
								MinSize:  Size{Width: 100, Height: 30},
								OnClicked: func() {
									username := usernameEdit.Text()
									password := passwordEdit.Text()

									if username == "" || password == "" {
										loginStatus.SetText("用户名和密码不能为空")
										return
									}

									loginStatus.SetText("正在登录...")
									success, token, message := network.DoLogin(username, password)

									if success {
										saveCredentials(username, password, rememberPassCB.Checked())

										result.Success = true
										result.Token = token
										result.Message = "登录成功"
										dlg.Synchronize(func() {
											LoadDeviceTypeNetwork()
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
