package network

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func DoJSONRequest(method, url string, requestBody interface{}, responseObj interface{}) (bool, string) {
	var reqBody []byte
	var err error

	// 如果有请求体，则序列化为JSON
	if requestBody != nil {
		reqBody, err = json.Marshal(requestBody)
		if err != nil {
			return false, "构建请求失败: " + err.Error()
		}
	}

	// 创建HTTP请求
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return false, "创建请求失败: " + err.Error()
	}

	// 设置请求头
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, "网络请求失败: " + err.Error()
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, "读取响应失败: " + err.Error()
	}

	// 解析响应
	if responseObj != nil {
		err = json.Unmarshal(body, responseObj)
		if err != nil {
			return false, "解析响应失败: " + err.Error()
		}
	}

	return true, ""
}

// DoFormRequest 发送表单格式的HTTP请求并解析响应
func DoFormRequest(method, urlStr string, formData map[string]string, responseObj interface{}) (bool, string) {
	// 构建表单数据
	form := url.Values{}
	for key, value := range formData {
		form.Add(key, value)
	}
	formDataStr := form.Encode()

	// 创建HTTP请求
	var req *http.Request
	var err error

	if method == "GET" {
		// 对于GET请求，将表单数据附加到URL
		if strings.Contains(urlStr, "?") {
			req, err = http.NewRequest(method, urlStr+"&"+formDataStr, nil)
		} else {
			req, err = http.NewRequest(method, urlStr+"?"+formDataStr, nil)
		}
	} else {
		// 对于POST等请求，将表单数据放在请求体中
		req, err = http.NewRequest(method, urlStr, strings.NewReader(formDataStr))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if err != nil {
		return false, "创建请求失败: " + err.Error()
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, "网络请求失败: " + err.Error()
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, "读取响应失败: " + err.Error()
	}

	// 解析响应
	if responseObj != nil {
		err = json.Unmarshal(body, responseObj)
		if err != nil {
			return false, "解析响应失败: " + err.Error()
		}
	}

	return true, ""
}
