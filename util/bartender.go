package util

import (
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"os"
	"os/exec"
	"path/filepath"
	"produce_tool/conf"
)

type BtwParam struct {
	KTXSN    string //箱号
	PO       string //订单号
	QTY      string //数量
	ITEMNAME string //备注
	ITEMCODE string //产品编码
	ITEMDESC string //品名规格
	SNLIST   []string
}

func GenBtwFile(sourceFilename, dstFilename string, param BtwParam) error {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	btApp, err := oleutil.CreateObject("BarTender.Application")
	if err != nil {
		return fmt.Errorf("BarTender 启动失败: %v", err)
	}
	btAppDispatch, err := btApp.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return fmt.Errorf("接口失败: %v", err)
	}
	defer btAppDispatch.Release()

	oleutil.PutProperty(btAppDispatch, "Visible", false)

	formats := oleutil.MustGetProperty(btAppDispatch, "Formats").ToIDispatch()
	cwd, _ := os.Getwd()
	templatePath := filepath.Join(cwd, sourceFilename)
	format := oleutil.MustCallMethod(formats, "Open", templatePath, false, "").ToIDispatch()

	oleutil.MustCallMethod(format, "SetNamedSubStringValue", "KTXSN", param.KTXSN)
	oleutil.MustCallMethod(format, "SetNamedSubStringValue", "PO", param.PO)
	oleutil.MustCallMethod(format, "SetNamedSubStringValue", "QTY", param.QTY)
	oleutil.MustCallMethod(format, "SetNamedSubStringValue", "ITEMNAME", param.ITEMNAME)
	oleutil.MustCallMethod(format, "SetNamedSubStringValue", "ITEMCODE", param.ITEMCODE)
	oleutil.MustCallMethod(format, "SetNamedSubStringValue", "ITEMDESC", param.ITEMDESC)

	for i, sn := range param.SNLIST {
		key := fmt.Sprintf("S/N%v", i)
		oleutil.MustCallMethod(format, "SetNamedSubStringValue", key, sn)
	}

	newPath := filepath.Join(cwd, dstFilename)
	oleutil.MustCallMethod(format, "SaveAs", newPath, true)

	oleutil.MustCallMethod(format, "Close", 2)
	oleutil.MustCallMethod(btAppDispatch, "Quit", 1)

	return nil
}

func Print(templateFilePath string) error {

	cmd := exec.Command(conf.BartendPath, fmt.Sprintf("/AF=%v", templateFilePath), "/P", "/X")

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("打印失败: %v", err)
	}

	return nil
}
