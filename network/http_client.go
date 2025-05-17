package network

import (
	"fmt"
)

func DoLogin(username, password string) (bool, string, string) {
	// 使用JSON方式登录
	loginReq := LoginRequest{
		Username: username,
		Password: password,
	}

	var loginResp LoginResponse
	success, errMsg := DoJSONRequest("POST", "http://factory.gps555.net/api/v1/login", loginReq, &loginResp)
	if !success {
		return false, "", errMsg
	}

	// 判断登录是否成功
	if loginResp.Code == 200 {
		return true, loginResp.Token, "登录成功"
	} else {
		return false, "", "登录失败，错误代码: " + fmt.Sprintf("%d", loginResp.Code)
	}
}
