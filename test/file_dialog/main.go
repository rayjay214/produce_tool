package main

import (
	"fmt"
	"github.com/lxn/walk"
)

func main() {
	var dlg *walk.FileDialog

	// 创建一个文件选择对话框
	dlg = &walk.FileDialog{
		Title:  "请选择uid文件",
		Filter: "CSV Files (*.csv)|*.csv|All Files (*.*)|*.*", // 设置过滤器
	}

	// 创建一个简单的窗体作为 "owner"
	var mw *walk.MainWindow
	mw, err := walk.NewMainWindow()
	if err != nil {
		fmt.Println("Error creating main window:", err)
		return
	}

	// 显示文件选择对话框并传入主窗口作为 owner
	accepted, err := dlg.ShowOpen(mw)
	if err != nil {
		fmt.Println("Error or no file selected:", err)
		return
	}

	// 如果文件被选择，输出文件路径
	if accepted {
		fmt.Println("File selected:", dlg.FilePath)
	} else {
		fmt.Println("No file selected.")
	}

	// 关闭主窗口
	mw.Close()
}
