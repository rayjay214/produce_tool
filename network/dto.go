package network

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 登录响应结构体
type LoginResponse struct {
	Code   int    `json:"code"`
	Expire string `json:"expire"`
	Token  string `json:"token"`
}

// LoginResult 登录结果
type LoginResult struct {
	Success bool
	Token   string
	Message string
}

type DeviceTypeDetail struct {
	DeviceType        string      `json:"deviceType"`
	FuncBit           int         `json:"funcBit"`
	MainIp            string      `json:"mainIp"`
	ViceIp            string      `json:"viceIp"`
	SignalOpen        string      `json:"signalOpen"`
	GpsOpen           string      `json:"gpsOpen"`
	WifiOpen          string      `json:"wifiOpen"`
	SnOpen            string      `json:"snOpen"`
	SimOpen           string      `json:"simOpen"`
	ImeiOpen          string      `json:"imeiOpen"`
	LightOpen         string      `json:"lightOpen"`
	GsensorOpen       string      `json:"gsensorOpen"`
	PowerOpen         string      `json:"powerOpen"`
	MainIpReadOpen    string      `json:"mainIpReadOpen"`
	ViceIpReadOpen    string      `json:"viceIpReadOpen"`
	MainIpWriteOpen   string      `json:"mainIpWriteOpen"`
	ViceIpWriteOpen   string      `json:"viceIpWriteOpen"`
	SetTypeOpen       string      `json:"setTypeOpen"`
	ProtocolOpen      string      `json:"protocolOpen"`
	ProtocolWriteOpen string      `json:"protocolWriteOpen"`
	SignalMin         string      `json:"signalMin"`
	SignalMax         string      `json:"signalMax"`
	GpsMin            string      `json:"gpsMin"`
	WifiMin           string      `json:"wifiMin"`
	PowerMin          string      `json:"powerMin"`
	ProtocolValue     string      `json:"protocolValue"`
	Prefix            string      `json:"prefix"`
	SnType            string      `json:"snType"`
	ImeiPrefix        string      `json:"imeiPrefix"`
	WriteSn           string      `json:"writeSn"`
	WriteImei         string      `json:"writeImei"`
	WriteApn          string      `json:"writeApn"`
	Apn               string      `json:"apn"`
	WriteWifiPrior    string      `json:"writeWifiPrior"`
	IsUltraLong       string      `json:"IsUltraLong"`
	CreatedAt         string      `json:"createdAt"`
	UpdatedAt         string      `json:"updatedAt"`
	DeletedAt         interface{} `json:"deletedAt"`
	CreateBy          int         `json:"createBy"`
	UpdateBy          int         `json:"updateBy"`
}

type DeviceTypeGetPageResponse struct {
	RequestId string `json:"requestId"`
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Data      struct {
		Count     int                `json:"count"`
		PageIndex int                `json:"pageIndex"`
		PageSize  int                `json:"pageSize"`
		List      []DeviceTypeDetail `json:"list"`
	} `json:"data"`
}

type DeviceTypeGetResponse struct {
	RequestId string           `json:"requestId"`
	Code      int              `json:"code"`
	Msg       string           `json:"msg"`
	Data      DeviceTypeDetail `json:"data"`
}

type GetInfoResponse struct {
	RequestId string `json:"requestId"`
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Data      struct {
		Username string `json:"username"`
		UserId   int64  `json:"userId"`
		Nickname string `json:"nickname"`
		RoleId   int64  `json:"roleId"`
		FamilyId int64  `json:"familyId"`
	}
}
