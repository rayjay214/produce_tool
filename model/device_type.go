package model

var AllTypes []DeviceTypeInfo

type DeviceTypeInfo struct {
	OverSpeedAlarm  int    //是否支持超速报警
	Listen          int    //是否支持监听
	LightControl    int    //是否支持灯控
	Sms             int    //是否支持短信
	TamperAlarm     int    //是否支持防拆报警
	ShakeAlarm      int    //是否支持震动报警
	Recording       int    //是否支持录音
	LowpowerAlarm   int    //是否支持低电报警
	MainIp          string //主IP地址
	MainPort        string //主IP端口
	ViceIp          string //副IP地址
	VicePort        string //副IP端口
	APN             string
	SnLength        string //SN长度校验
	DeviceType      string //设备型号
	SignalOpen      int    //是否开启信号测试
	SignalDelay     string //信号测试延时，0为关闭
	SignalMin       string //测试通过的最小值
	SignalMax       string //测试通过的最大值
	GpsOpen         int    //是否开启GPS卫星测试
	GpsDelay        string //GPS测试延时，0为关闭
	GpsMin          string //GPS测试通过的最小值
	WifiOpen        int    //是否开启WIFI测试
	WifiMin         string //WIFI测试通过的最小值
	SnOpen          int    //是否开启读取SN
	DialOpen        int    //是否开启打电话测试
	SimOpen         int    //是否开启读取SIM
	ImeiOpen        int    //是否开启读取Imei
	LightOpen       int    //是否开启光感测试
	GsensorOpen     int    //是否开启重力测试
	EchoOpen        int    //是否开启回音测试
	EndDialOpen     int    //是否开启挂断测试
	TamperOpen      int    //是否开启防拆测试
	SetTypeOpen     int    //是否开启设置型号
	MainIpReadOpen  int    //是否开启读取主IP
	ViceIpReadOpen  int    //是否开启读取副IP
	ApnWriteOpen    int    //是否开启写入APN
	ViceIpWriteOpen int    //是否开启写入副IP
}

var DeviceTypeInfoMap map[string]DeviceTypeInfo

func init() {
	DeviceTypeInfoMap = make(map[string]DeviceTypeInfo, 0)
	types := GetDeviceTypes()
	AllTypes = make([]DeviceTypeInfo, 0)
	for _, value := range types {
		AllTypes = append(AllTypes, value)
	}
}
