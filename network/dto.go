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
	CreatedAt         string      `json:"createdAt"`
	UpdatedAt         string      `json:"updatedAt"`
	DeletedAt         interface{} `json:"deletedAt"`
	CreateBy          int         `json:"createBy"`
	UpdateBy          int         `json:"updateBy"`
}

type ProductionPlan struct {
	Id         int64       `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`
	Name       string      `json:"name" gorm:"type:varchar(64);comment:计划名称"`              // 计划名称
	DeviceType string      `json:"deviceType" gorm:"type:varchar(32);comment:设备型号"`        // 设备型号
	Count      int64       `json:"count" gorm:"type:bigint unsigned;comment:生产设备数量"`       // 生产设备数量
	SnType     string      `json:"snType" gorm:"type:varchar(4);comment:设备号类型，字典：sn_type"` // 设备号类型，字典：sn_type
	Remark     string      `json:"remark" gorm:"type:varchar(1024);comment:备注"`            // 备注
	CreatedAt  string      `json:"createdAt"`
	UpdatedAt  string      `json:"updatedAt"`
	DeletedAt  interface{} `json:"deletedAt"`
	CreateBy   int         `json:"createBy"`
	UpdateBy   int         `json:"updateBy"`
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

type PlanGetPageResponse struct {
	RequestId string `json:"requestId"`
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Data      struct {
		Count     int              `json:"count"`
		PageIndex int              `json:"pageIndex"`
		PageSize  int              `json:"pageSize"`
		List      []ProductionPlan `json:"list"`
	} `json:"data"`
}

type PlanGetResponse struct {
	RequestId string         `json:"requestId"`
	Code      int            `json:"code"`
	Msg       string         `json:"msg"`
	Data      ProductionPlan `json:"data"`
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
