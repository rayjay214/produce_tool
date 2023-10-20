package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func GetDeviceTypes() map[string]DeviceTypeInfo {
	filePath := "deivce_types.json"

	_, err := os.Stat(filePath)
	if err != nil {
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return nil
		}
		jsonData, _ := json.Marshal(DeviceTypeInfoMap)
		_, err = file.Write(jsonData)
		return DeviceTypeInfoMap
	}

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Failed to open or create file:", err)
		return nil
	}
	defer file.Close()
	dataBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}
	var types map[string]DeviceTypeInfo
	json.Unmarshal(dataBytes, &types)
	return types
}

func SyncDeviceTypes(types map[string]DeviceTypeInfo) {
	filePath := "deivce_types.json"
	file, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Failed to open or create file:", err)
	}
	defer file.Close()
	jsonData, _ := json.Marshal(types)
	fmt.Println(string(jsonData))
	file.Truncate(0)
	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println(err)
	}
}

func GetCurrType() string {
	filePath := "type.dat"

	_, err := os.Stat(filePath)
	if err != nil {
		_, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return ""
		}
		return ""
	}

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Failed to open or create file:", err)
		return ""
	}
	defer file.Close()
	dataBytes, err := ioutil.ReadAll(file)
	return string(dataBytes)
}

func SyncCurrType(strType string) {
	filePath := "type.dat"
	file, err := os.OpenFile(filePath, os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Failed to open or create file:", err)
	}
	defer file.Close()
	_, err = file.Write([]byte(strType))
}
