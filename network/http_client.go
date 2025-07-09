package network

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	Token       string
	Username    string
	CurrentPlan ProductionPlan
	CurrentType DeviceTypeDetail
	TotalCount  int
	PassedCount int
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
	var getPageResp DeviceTypeGetPageResponse
	success, errMsg := DoFormRequest("GET", "http://factory.gps555.net/api/v1/function", nil, &getPageResp)
	if !success {
		return nil, fmt.Errorf(errMsg)
	}

	if getPageResp.Code != 200 {
		return nil, fmt.Errorf("请求失败，错误代码: %d, 错误信息: %s", getPageResp.Code, getPageResp.Msg)
	}

	return getPageResp.Data.List, nil
}

func DoGetUserInfo() error {
	var getInfoResponse GetInfoResponse
	success, errMsg := DoFormRequest("GET", "http://factory.gps555.net/api/v1/getinfo", nil, &getInfoResponse)
	if !success {
		return fmt.Errorf(errMsg)
	}

	if getInfoResponse.Code != 200 {
		return fmt.Errorf("请求失败，错误代码: %d", getInfoResponse.Code)
	}

	Username = getInfoResponse.Data.Username

	return nil
}

func DoGetPlanList() ([]ProductionPlan, error) {
	var getPageResp PlanGetPageResponse
	form := map[string]string{
		"status": "0",
	}
	success, errMsg := DoFormRequest("GET", "http://factory.gps555.net/api/v1/production-plan", form, &getPageResp)
	if !success {
		return nil, fmt.Errorf(errMsg)
	}

	if getPageResp.Code != 200 {
		return nil, fmt.Errorf("请求失败，错误代码: %d, 错误信息: %s", getPageResp.Code, getPageResp.Msg)
	}

	return getPageResp.Data.List, nil
}

func DoGetPlan(planId int64) {
	var getResp PlanGetResponse
	url := fmt.Sprintf("http://factory.gps555.net/api/v1/production-plan/%v", planId)
	success, _ := DoFormRequest("GET", url, nil, &getResp)
	if !success {
		return
	}

	if getResp.Code != 200 {
		return
	}

	CurrentPlan = getResp.Data
}

func DoGetDeviceType(deviceType string) {
	var getResp DeviceTypeGetResponse
	url := fmt.Sprintf("http://factory.gps555.net/api/v1/function/%v", deviceType)
	success, _ := DoFormRequest("GET", url, nil, &getResp)
	if !success {
		return
	}

	if getResp.Code != 200 {
		return
	}

	CurrentType = getResp.Data
}

func DoGetTotalNum() {
	var getResp DeviceGetPageResponse
	url := fmt.Sprintf("http://factory.gps555.net/api/v1/device")
	form := map[string]string{
		"planId": fmt.Sprintf("%v", CurrentPlan.Id),
	}
	success, _ := DoFormRequest("GET", url, form, &getResp)
	if !success {
		return
	}

	if getResp.Code != 200 {
		return
	}

	TotalCount = getResp.Data.Count
}

func DoGetPassedNum() {
	var getResp DeviceGetPageResponse
	url := fmt.Sprintf("http://factory.gps555.net/api/v1/device")
	form := map[string]string{
		"planId":      fmt.Sprintf("%v", CurrentPlan.Id),
		"writeStatus": "1",
		"testStatus":  "1",
	}
	success, _ := DoFormRequest("GET", url, form, &getResp)
	if !success {
		return
	}

	if getResp.Code != 200 {
		return
	}

	PassedCount = getResp.Data.Count
}

func DoGetComparePassedNum() int {
	var getResp DeviceGetPageResponse
	url := fmt.Sprintf("http://factory.gps555.net/api/v1/device")
	form := map[string]string{
		"planId":        fmt.Sprintf("%v", CurrentPlan.Id),
		"compareStatus": "1",
	}
	success, _ := DoFormRequest("GET", url, form, &getResp)
	if !success {
		return 0
	}

	if getResp.Code != 200 {
		return 0
	}

	return getResp.Data.Count
}

// 获取装箱最新的箱号
func DoGetBoxNo() string {
	var getResp BoxGetPageResponse
	url := fmt.Sprintf("http://factory.gps555.net/api/v1/box")
	form := map[string]string{
		"planId":     fmt.Sprintf("%v", CurrentPlan.Id),
		"boxNoOrder": "desc",
	}
	success, _ := DoFormRequest("GET", url, form, &getResp)
	if !success {
		return ""
	}

	if getResp.Code != 200 {
		return ""
	}

	//第一箱
	if getResp.Data.Count == 0 {
		return fmt.Sprintf("%v0001", CurrentPlan.OrderNo)
	}

	currNo := getResp.Data.List[0].BoxNo
	lastFour := strings.TrimPrefix(currNo, currNo[:len(currNo)-4])
	num, _ := strconv.Atoi(lastFour)
	num = num + 1
	boxNo := fmt.Sprintf("%v%04d", CurrentPlan.OrderNo, num)
	return boxNo
}
