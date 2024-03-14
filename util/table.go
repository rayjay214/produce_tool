package util

import (
	"github.com/lxn/walk"
	"reflect"
)

type MyTableRow struct {
	Com         string
	Pass        string
	Version     string
	Signal      string
	Wifi        string
	Sim         string
	Sn          string
	Imei        string
	Gps         string
	Light       string
	MainIp      string
	ViceIp      string
	Gsensor     string
	SetType     string
	ViceIpWrite string
	Power       string
}

type MyTableModel struct {
	walk.TableModelBase
	items []*MyTableRow
}

var ColumnIdxNames map[int]string //列对应的Model结构体字段名
var ColumnNamesIdx map[string]int //列对应的Model结构体字段名

var PortNameRowidx map[string]int //串口名所在的行
var RowidxPortName map[int]string //行所对应的串口名

var tableModel *MyTableModel

var tv *walk.TableView

func GetTableView() *walk.TableView {
	tv = new(walk.TableView)
	return tv
}

func GetTableModel() *MyTableModel {
	return tableModel
}

func GetModelItems() []*MyTableRow {
	return tableModel.items
}

func init() {
	ColumnIdxNames = make(map[int]string, 0)
	ColumnNamesIdx = make(map[string]int, 0)

	ArrangeCols()
	tableModel = NewMyTableModel()
	tableModel.Refresh()
}

func RefreshTableModel() {
	ArrangeCols()
	tableModel.Refresh()
}

func ArrangeCols() {
	ColumnIdxNames[0] = "Com"
	ColumnIdxNames[1] = "Pass"
	ColumnNamesIdx["Com"] = 0
	ColumnNamesIdx["Pass"] = 1
	items := GetAllTestItems()
	for i, item := range items {
		if item.IsShow {
			ColumnIdxNames[i+1] = item.ModelColName
			ColumnNamesIdx[item.ModelColName] = i + 1
		}
	}
}

func (m *MyTableModel) Refresh() {
	tableModel.items = make([]*MyTableRow, 0)
	PortNameRowidx = make(map[string]int, 0)
	RowidxPortName = make(map[int]string, 0)

	myPorts := GetPorts()
	for i, myPort := range myPorts {
		row := MyTableRow{Com: myPort.Name}
		tableModel.items = append(tableModel.items, &row)
		PortNameRowidx[myPort.Name] = i
		RowidxPortName[i] = myPort.Name
	}
	m.PublishRowsReset()
}

func (m *MyTableModel) ClearRow(idx int) {
	if idx < m.RowCount() {
		rv := reflect.ValueOf(m.items[idx])
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		rt := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			fieldType := rt.Field(i)
			if fieldType.Name != "Com" {
				fieldValue := rv.Field(i)
				fieldValue.SetString("")
			}
		}
		m.PublishRowChanged(idx)
	}
}

func (m *MyTableModel) ResetRows() {
	for i := 0; i < m.RowCount(); i++ {
		m.ClearRow(i)
	}
}

func NewMyTableModel() *MyTableModel {
	m := new(MyTableModel)
	m.ResetRows()
	return m
}

func (m *MyTableModel) RowCount() int {
	return len(m.items)
}

func (m *MyTableModel) Value(row, col int) interface{} {
	item := m.items[row]
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.IsValid() {
		return v.FieldByName(ColumnIdxNames[col]).Interface()
	}
	return ""
}
