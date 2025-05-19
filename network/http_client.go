package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	Token string
)

func DoLogin(username, password string) (bool, string, string) {
	loginReq := LoginRequest{
		Username: username,
		Password: password,
	}

	var loginResp LoginResponse
	success, errMsg := DoJSONRequest("POST", "http://factory.gps555.net/api/v1/login", loginReq, &loginResp)
	if !success {
		return false, "", errMsg
	}

	if loginResp.Code == 200 {
		Token = loginResp.Token
		return true, loginResp.Token, "登录成功"
	} else {
		return false, "", "登录失败，错误代码: " + fmt.Sprintf("%d", loginResp.Code)
	}
}

func DoGetDeviceTypes() ([]DeviceTypeDetail, error) {
	// 创建请求
	req, err := http.NewRequest("GET", "http://factory.gps555.net/api/v1/function", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头，添加token
	req.Header.Set("Authorization", "Bearer "+Token)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("网络请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析响应
	var getPageResp GetPageResponse
	err = json.Unmarshal(body, &getPageResp)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查响应状态
	if getPageResp.Code != 200 {
		return nil, fmt.Errorf("请求失败，错误代码: %d, 错误信息: %s", getPageResp.Code, getPageResp.Msg)
	}

	return getPageResp.Data.List, nil
}
