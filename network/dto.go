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
	FuccBit           int         `json:"fuccBit"`
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
	CreatedAt         string      `json:"createdAt"`
	UpdatedAt         string      `json:"updatedAt"`
	DeletedAt         interface{} `json:"deletedAt"`
	CreateBy          int         `json:"createBy"`
	UpdateBy          int         `json:"updateBy"`
}

type GetPageResponse struct {
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
