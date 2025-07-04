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
	defer formats.Release()
	cwd, _ := os.Getwd()
	templatePath := filepath.Join(cwd, sourceFilename)
	format := oleutil.MustCallMethod(formats, "Open", templatePath, false, "").ToIDispatch()
	defer format.Release()

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

	//先默认打印一次
	_, err = oleutil.CallMethod(format, "PrintOut", false, false)
	if err != nil {
		return fmt.Errorf("打印失败: %v", err)
	}

	//同时保存模板文件
	newPath := filepath.Join(cwd, "template", dstFilename)
	oleutil.MustCallMethod(format, "SaveAs", newPath, true)

	oleutil.MustCallMethod(format, "Close", 2)

	oleutil.MustCallMethod(btAppDispatch, "Quit", 1)

	return nil
}

func PrintFileCmd(templateFilePath string) error {

	cmd := exec.Command(conf.BartendPath, fmt.Sprintf("/AF=%v", templateFilePath), "/P", "/X")

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("打印失败: %v", err)
	}

	return nil
}

func PrintFileOle(templateFilePath string) error {
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

	formats := oleutil.MustGetProperty(btAppDispatch, "Formats").ToIDispatch()
	format := oleutil.MustCallMethod(formats, "Open", templateFilePath, false, "").ToIDispatch()

	_, err = oleutil.CallMethod(format, "PrintOut", false, false)
	if err != nil {
		return fmt.Errorf("打印失败: %v", err)
	}

	_, err = oleutil.CallMethod(format, "Close", 0) // 0 = don't save changes
	if err != nil {
		return fmt.Errorf("关闭出错: %v", err)
	}
	format.Release()

	_, _ = oleutil.CallMethod(btAppDispatch, "Quit")

	return nil
}

func PrintInMemory(templateFilename string, sn string, cnt int) error {
	btApp, err := oleutil.CreateObject("BarTender.Application")
	if err != nil {
		return fmt.Errorf("BarTender 启动失败: %v", err)
	}
	btAppDispatch, err := btApp.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return fmt.Errorf("接口失败: %v", err)
	}
	defer btAppDispatch.Release()

	formats := oleutil.MustGetProperty(btAppDispatch, "Formats").ToIDispatch()
	defer formats.Release()
	cwd, _ := os.Getwd()
	templatePath := filepath.Join(cwd, "template", templateFilename)
	format := oleutil.MustCallMethod(formats, "Open", templatePath, false, "").ToIDispatch()
	defer format.Release()

	oleutil.MustCallMethod(format, "SetNamedSubStringValue", "sn", sn)

	for i := 0; i < cnt; i++ {
		_, err = oleutil.CallMethod(format, "PrintOut", false, false)
		if err != nil {
			return fmt.Errorf("打印失败: %v", err)
		}
	}

	_, err = oleutil.CallMethod(format, "Close", 0) // 0 = don't save changes
	if err != nil {
		return fmt.Errorf("关闭出错: %v", err)
	}
	format.Release()

	_, _ = oleutil.CallMethod(btAppDispatch, "Quit")

	return nil
}
