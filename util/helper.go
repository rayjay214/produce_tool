package util

import (
	"github.com/lxn/walk"
	"produce_tool/model"
	"reflect"
	"strings"
)

func getValue(str, findS string) (int, string) {
	sLeft := ""
	if str == "" {
		return 0, sLeft
	}

	nPos := 0
	if findS != "" {
		nPos = strings.Index(str, findS)
	}
	if nPos < 0 {
		return 0, sLeft
	}

	sLeft = str[nPos+len(findS):]
	sLeft = strings.TrimLeft(sLeft, " ")
	sLeft = strings.TrimLeft(sLeft, "\r\n")

	nPos = strings.Index(sLeft, "\r\n")
	if nPos >= 0 {
		sLeft = sLeft[:nPos]
	}

	sLeft = strings.TrimRight(sLeft, " ")
	sLeft = strings.TrimRight(sLeft, "\r\n")

	return len(sLeft), sLeft
}

func ContainsOne(str string, substrs ...string) bool {
	for _, substr := range substrs {
		if strings.Contains(str, substr) {
			return true
		}
	}
	return false
}

func GetFromStatus(data string) map[string]string {
	fields := strings.Split(data, ",")
	result := make(map[string]string)
	for _, field := range fields {
		keyValue := strings.SplitN(field, ":", 2)
		if len(keyValue) != 2 {
			continue
		}
		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])
		result[key] = value
	}
	return result
}

func StopReader(param *PassParam) {
	param.stopReader = true
}

func StopWriter(param *PassParam) {
	param.stopWriter = true
}

func GetPassParamStr(param *PassParam) string {
	return param.str
}

func SetPassParamStr(param *PassParam, value string) {
	param.str = value
}

func FilterTableColumn(tv *walk.TableView, selectedType model.DeviceTypeInfo) {
	v := reflect.ValueOf(selectedType)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name

		if ContainsOne(fieldName, "DialOpen", "EndDialOpen", "TamperOpen", "ApnWriteOpen") {
			continue
		}

		if strings.Contains(fieldName, "Open") {
			value := v.Field(i).Interface()
			if intValue, ok := value.(int); ok {
				if intValue <= 0 {
					colName := fieldName[:len(fieldName)-4]
					if colName == "MainIpRead" {
						colName = "MainIp"
					}
					if colName == "ViceIpRead" {
						colName = "ViceIp"
					}
					tv.Columns().ByName(colName).SetVisible(false)
				}
			}
		}
	}
}
