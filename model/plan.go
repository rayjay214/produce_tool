package model

import (
	"fmt"
	"produce_tool/network"
)

var AllPlans []PlanInfo

type PlanInfo struct {
	Id          int64  // id
	Name        string // 计划名称
	DeviceType  string // 设备型号
	Count       int64  // 生产设备数量
	SnType      string // 设备号类型，字典：sn_type
	Remark      string // 备注
	DisplayName string // 展示名
}

func LoadPlanNetwork() {
	netPlans, _ := network.DoGetPlanList()

	AllPlans = make([]PlanInfo, 0)
	for _, item := range netPlans {
		plan := PlanInfo{
			Id:          item.Id,
			Name:        item.Name,
			DeviceType:  item.DeviceType,
			Count:       item.Count,
			SnType:      item.SnType,
			Remark:      item.Remark,
			DisplayName: fmt.Sprintf("%v(%v_%v)", item.Name, item.DeviceType, item.Count),
		}
		AllPlans = append(AllPlans, plan)
	}
}
