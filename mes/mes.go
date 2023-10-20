package mes

import (
	"bytes"
	"encoding/xml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"produce_tool/conf"
)

type SoapHeader struct {
	XMLName xml.Name `xml:"soap:Header"`
	Content MySoapHeader
}

type MySoapHeader struct {
	XMLName  xml.Name `xml:"http://tempuri.org/ MySoapHeader"`
	UserName string   `xml:"UserName"`
	Password string   `xml:"Password"`
}

type OverStations struct {
	XMLName xml.Name `xml:"soap:Envelope"`
	SoapNS  string   `xml:"xmlns:soap,attr"`
	XsiNS   string   `xml:"xmlns:xsi,attr"`
	XsdNS   string   `xml:"xmlns:xsd,attr"`
	Header  SoapHeader
	Body    SoapBody
}

type SoapBody struct {
	XMLName xml.Name `xml:"soap:Body"`
	Content OverStationsBody
}

type OverStationsBody struct {
	XMLName       xml.Name `xml:"http://tempuri.org/ OverStations"`
	WorkProcedure string   `xml:"workprocedure"`
	IsnCode       string   `xml:"IsnCode"`
	OuterCode     string   `xml:"OuterCode"`
	OuterCode1    string   `xml:"OuterCode1"`
	OuterCode2    string   `xml:"OuterCode2"`
	OuterCode3    string   `xml:"OuterCode3"`
	User          string   `xml:"User"`
	Detail        string   `xml:"Detail"`
}

type OverStationsResponse struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    struct {
		OverStationsResponse struct {
			XMLName            xml.Name `xml:"http://tempuri.org/ OverStationsResponse"`
			OverStationsResult string   `xml:"OverStationsResult"`
			Msg                string   `xml:"Msg"`
		} `xml:"http://tempuri.org/ OverStationsResponse"`
	} `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
}

func constructBody(sn string, detail string, workprocedure string) string {
	header := MySoapHeader{
		UserName: "admin",
		Password: "admin123456",
	}

	body := OverStationsBody{
		//WorkProcedure: "DIGNWEIQICESHI",
		WorkProcedure: workprocedure,
		IsnCode:       sn,
		OuterCode:     "",
		OuterCode1:    "",
		OuterCode2:    "",
		OuterCode3:    "",
		User:          "ADMIN",
		Detail:        detail,
	}

	xmlData := OverStations{
		SoapNS: "http://schemas.xmlsoap.org/soap/envelope/",
		XsiNS:  "http://www.w3.org/2001/XMLSchema-instance",
		XsdNS:  "http://www.w3.org/2001/XMLSchema",
		Header: SoapHeader{
			Content: header,
		},
		Body: SoapBody{
			Content: body,
		},
	}

	xmlBytes, err := xml.MarshalIndent(xmlData, "", "    ")
	if err != nil {
		return ""
	}

	xmlString := xml.Header + string(xmlBytes)
	return xmlString
}

func constructCheckBody(sn string, workprocedure string) string {
	header := CheckMySoapHeader{
		UserName: "admin",
		Password: "admin123456",
	}

	body := IsOverStationBody{
		//WorkProcedure: "DIGNWEIQICESHI",
		WorkProcedure: workprocedure,
		IsnCode:       sn,
	}

	xmlData := IsOverStations{
		SoapNS: "http://schemas.xmlsoap.org/soap/envelope/",
		XsiNS:  "http://www.w3.org/2001/XMLSchema-instance",
		XsdNS:  "http://www.w3.org/2001/XMLSchema",
		Header: CheckSoapHeader{
			Content: header,
		},
		Body: CheckSoapBody{
			Content: body,
		},
	}

	xmlBytes, err := xml.MarshalIndent(xmlData, "", "  ")
	if err != nil {
		log.Errorf("Error:%v", err)
		return ""
	}

	xmlString := xml.Header + string(xmlBytes)
	return xmlString
}

func SetMesReq(sn string, detail string, workprocedure string) bool {
	url := conf.MesUrl

	// 构造POST请求的body数据
	payload := []byte(constructBody(sn, detail, workprocedure))
	log.Infof("req body is %v", string(payload))

	if url == "" {
		url = "http://172.18.2.242:83/MESWeb/DataConnection.asmx"
	}

	// 发送POST请求
	resp, err := http.Post(url, "text/xml;charset=utf-8", bytes.NewBuffer(payload))
	if err != nil {
		log.Infof("Error sending request:%v", err)
		return false
	}
	defer resp.Body.Close()

	// 处理响应
	var bRst bool
	if resp.StatusCode == http.StatusOK {
		log.Infof("Request successful!")
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Error reading response:%v", err)
			return false
		}
		log.Infof("resp body is %v", string(respBody))
		var response OverStationsResponse
		err = xml.Unmarshal([]byte(respBody), &response)
		if err != nil {
			log.Errorf("Error:%v", err)
			return false
		}
		if response.Body.OverStationsResponse.OverStationsResult == "true" {
			bRst = true
		} else {
			bRst = false
		}

	} else {
		log.Infof("Request failed with status code:%v", resp.StatusCode)
	}
	return bRst
}

type CheckSoapHeader struct {
	XMLName xml.Name `xml:"soap:Header"`
	Content CheckMySoapHeader
}

type CheckMySoapHeader struct {
	XMLName  xml.Name `xml:"http://tempuri.org/ MySoapHeader"`
	UserName string   `xml:"UserName"`
	Password string   `xml:"Password"`
}

type IsOverStations struct {
	XMLName xml.Name `xml:"soap:Envelope"`
	SoapNS  string   `xml:"xmlns:soap,attr"`
	XsiNS   string   `xml:"xmlns:xsi,attr"`
	XsdNS   string   `xml:"xmlns:xsd,attr"`
	Header  CheckSoapHeader
	Body    CheckSoapBody
}

type CheckSoapBody struct {
	XMLName xml.Name `xml:"soap:Body"`
	Content IsOverStationBody
}

type IsOverStationBody struct {
	XMLName       xml.Name `xml:"http://tempuri.org/ IsOverStation"`
	WorkProcedure string   `xml:"workprocedure"`
	IsnCode       string   `xml:"ISNCode"`
}

type IsOverStationResponse struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    IsOverStationResponseBody
}

type IsOverStationResponseBody struct {
	XMLName               xml.Name                  `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	IsOverStationResponse IsOverStationResponseData `xml:"IsOverStationResponse"`
}

type IsOverStationResponseData struct {
	XMLName             xml.Name `xml:"http://tempuri.org/ IsOverStationResponse"`
	IsOverStationResult string   `xml:"IsOverStationResult"`
	Msg                 string   `xml:"Msg"`
}

func CheckMesReq(sn string, workprocedure string) bool {
	url := conf.MesUrl

	// 构造POST请求的body数据
	payload := []byte(constructCheckBody(sn, workprocedure))
	log.Infof("check req body is %v", string(payload))

	if url == "" {
		url = "http://172.18.2.242:83/MESWeb/DataConnection.asmx"
	}

	// 发送POST请求
	resp, err := http.Post(url, "text/xml;charset=utf-8", bytes.NewBuffer(payload))
	if err != nil {
		log.Infof("Error sending request:%v", err)
		return false
	}
	defer resp.Body.Close()

	// 处理响应
	var bRst bool
	if resp.StatusCode == http.StatusOK {
		log.Infof("Request successful!")
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Error reading response:%v", err)
			return false
		}
		log.Infof("check resp body is %v", string(respBody))
		var response IsOverStationResponse
		err = xml.Unmarshal([]byte(respBody), &response)
		if err != nil {
			log.Errorf("Error:%v", err)
			return false
		}
		if response.Body.IsOverStationResponse.IsOverStationResult == "true" {
			bRst = true
		} else {
			bRst = false
		}
	} else {
		log.Infof("Request failed with status code:%v", resp.StatusCode)
	}
	return bRst
}
