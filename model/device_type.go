package model

import (
	"produce_tool/network"
	"strings"
)

var AllTypes []DeviceTypeInfo

type DeviceTypeInfo struct {
	OverSpeedAlarm    int    //是否支持超速报警
	Listen            int    //是否支持监听
	LightControl      int    //是否支持灯控
	Sms               int    //是否支持短信
	TamperAlarm       int    //是否支持防拆报警
	ShakeAlarm        int    //是否支持震动报警
	Recording         int    //是否支持录音
	LowpowerAlarm     int    //是否支持低电报警
	RapidAccleAlarm   int    //是否支持急加速报警
	RapidDecleAlarm   int    //是否支持急减速报警
	SharpTurnAlarm    int    //是否支持急转弯报警
	MainIp            string //主IP地址
	MainPort          string //主IP端口
	ViceIp            string //副IP地址
	VicePort          string //副IP端口
	APN               string
	SnLength          string //SN长度校验
	DeviceType        string //设备型号
	SignalOpen        int    //是否开启信号测试
	SignalDelay       string //信号测试延时，0为关闭
	SignalMin         string //测试通过的最小值
	SignalMax         string //测试通过的最大值
	GpsOpen           int    //是否开启GPS卫星测试
	GpsDelay          string //GPS测试延时，0为关闭
	GpsMin            string //GPS测试通过的最小值
	WifiOpen          int    //是否开启WIFI测试
	WifiMin           string //WIFI测试通过的最小值
	SnOpen            int    //是否开启读取SN
	DialOpen          int    //是否开启打电话测试
	SimOpen           int    //是否开启读取SIM
	ImeiOpen          int    //是否开启读取Imei
	LightOpen         int    //是否开启光感测试
	GsensorOpen       int    //是否开启重力测试
	PowerOpen         int    //是否开启电量测试
	EndDialOpen       int    //是否开启挂断测试
	TamperOpen        int    //是否开启防拆测试
	SetTypeOpen       int    //是否开启设置型号
	MainIpReadOpen    int    //是否开启读取主IP
	ViceIpReadOpen    int    //是否开启读取副IP
	ApnWriteOpen      int    //是否开启写入APN
	MainIpWriteOpen   int    //是否开启写入主IP
	ViceIpWriteOpen   int    //是否开启写入副IP
	PowerMin          string //电量测试通过的最小值
	ProtocolOpen      int    //是否开启读取协议
	ProtocolWriteOpen int    //是否开启写入协议
	ProtocolValue     string //待设置的协议
}

var DeviceTypeInfoMap map[string]DeviceTypeInfo

func LoadDeviceType() {
	DeviceTypeInfoMap = make(map[string]DeviceTypeInfo, 0)
	types := GetDeviceTypes()
	AllTypes = make([]DeviceTypeInfo, 0)
	for _, value := range types {
		AllTypes = append(AllTypes, value)
	}
}

func getBit(num, pos int) int {
	mask := 1 << uint(pos)
	if (num & mask) != 0 {
		return 1
	}
	return 0
}

func LoadDeviceTypeNetwork() {
	DeviceTypeInfoMap = make(map[string]DeviceTypeInfo, 0)
	netTypes, _ := network.DoGetDeviceTypes()
	deviceTypes := make(map[string]DeviceTypeInfo)
	for _, item := range netTypes {
		signalOpen := 0
		if item.SignalOpen == "1" {
			signalOpen = 1
		}
		gpsOpen := 0
		if item.GpsOpen == "1" {
			gpsOpen = 1
		}
		wifiOpen := 0
		if item.WifiOpen == "1" {
			wifiOpen = 1
		}
		snOpen := 0
		if item.SnOpen == "1" {
			snOpen = 1
		}
		simOpen := 0
		if item.SimOpen == "1" {
			simOpen = 1
		}
		imeiOpen := 0
		if item.ImeiOpen == "1" {
			imeiOpen = 1
		}
		lightOpen := 0
		if item.LightOpen == "1" {
			lightOpen = 1
		}
		gsensorOpen := 0
		if item.GsensorOpen == "1" {
			gsensorOpen = 1
		}
		powerOpen := 0
		if item.PowerOpen == "1" {
			powerOpen = 1
		}
		setTypeOpen := 0
		if item.SetTypeOpen == "1" {
			setTypeOpen = 1
		}
		mainIpReadOpen := 0
		if item.MainIpReadOpen == "1" {
			mainIpReadOpen = 1
		}
		viceIpReadOpen := 0
		if item.ViceIpReadOpen == "1" {
			viceIpReadOpen = 1
		}
		viceIpWriteOpen := 0
		if item.ViceIpWriteOpen == "1" {
			viceIpWriteOpen = 1
		}

		mainIpWriteOpen := 0
		if item.MainIpWriteOpen == "1" {
			mainIpWriteOpen = 1
		}
		protocolOpen := 0
		if item.ProtocolOpen == "1" {
			protocolOpen = 1
		}
		protocolWriteOpen := 0
		if item.ProtocolWriteOpen == "1" {
			protocolWriteOpen = 1
		}
		apnWriteOpen := 0
		if item.WriteApn == "1" {
			apnWriteOpen = 1
		}

		var mainIp, mainPort string
		if item.MainIp != "" {
			listMainIp := strings.Split(item.MainIp, ":")
			if len(listMainIp) == 2 {
				mainIp = listMainIp[0]
				mainPort = listMainIp[1]
			}
		}
		var viceIp, vicePort string
		if item.ViceIp != "" {
			listViceIp := strings.Split(item.ViceIp, ":")
			if len(listViceIp) == 2 {
				viceIp = listViceIp[0]
				vicePort = listViceIp[1]
			}
		}

		overSpeedAlarm := 0
		if getBit(item.FuncBit, 6) == 1 {
			overSpeedAlarm = 1
		}
		tamperAlarm := 0
		if getBit(item.FuncBit, 3) == 1 {
			tamperAlarm = 1
		}
		shakeAlarm := 0
		if getBit(item.FuncBit, 4) == 1 {
			shakeAlarm = 1
		}
		lowpowerAlarm := 0
		if getBit(item.FuncBit, 5) == 1 {
			lowpowerAlarm = 1
		}
		sharpTurnAlarm := 0
		if getBit(item.FuncBit, 14) == 1 {
			sharpTurnAlarm = 1
		}
		rapidAccleAlarm := 0
		if getBit(item.FuncBit, 12) == 1 {
			rapidAccleAlarm = 1
		}
		rapidDecleAlarm := 0
		if getBit(item.FuncBit, 13) == 1 {
			rapidDecleAlarm = 1
		}

		deviceInfo := DeviceTypeInfo{
			DeviceType:        item.DeviceType,
			MainIp:            mainIp,
			MainPort:          mainPort,
			ViceIp:            viceIp,
			VicePort:          vicePort,
			SignalOpen:        signalOpen,
			GpsOpen:           gpsOpen,
			WifiOpen:          wifiOpen,
			SnOpen:            snOpen,
			SimOpen:           simOpen,
			ImeiOpen:          imeiOpen,
			LightOpen:         lightOpen,
			GsensorOpen:       gsensorOpen,
			PowerOpen:         powerOpen,
			SetTypeOpen:       setTypeOpen,
			MainIpReadOpen:    mainIpReadOpen,
			MainIpWriteOpen:   mainIpWriteOpen,
			ViceIpReadOpen:    viceIpReadOpen,
			ViceIpWriteOpen:   viceIpWriteOpen,
			ProtocolOpen:      protocolOpen,
			ProtocolWriteOpen: protocolWriteOpen,
			SignalMin:         item.SignalMin,
			SignalMax:         item.SignalMax,
			GpsMin:            item.GpsMin,
			WifiMin:           item.WifiMin,
			PowerMin:          item.PowerMin,
			ProtocolValue:     item.ProtocolValue,
			OverSpeedAlarm:    overSpeedAlarm,
			TamperAlarm:       tamperAlarm,
			ShakeAlarm:        shakeAlarm,
			LowpowerAlarm:     lowpowerAlarm,
			SharpTurnAlarm:    sharpTurnAlarm,
			RapidDecleAlarm:   rapidDecleAlarm,
			RapidAccleAlarm:   rapidAccleAlarm,
			ApnWriteOpen:      apnWriteOpen,
			APN:               item.Apn,
		}

		deviceTypes[item.DeviceType] = deviceInfo
	}
	AllTypes = make([]DeviceTypeInfo, 0)
	for _, value := range deviceTypes {
		AllTypes = append(AllTypes, value)
	}
}
