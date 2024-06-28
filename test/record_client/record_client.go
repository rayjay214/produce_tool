package main

import (
	"bufio"
	"fmt"
	"github.com/tarm/serial"
	"io"
	"os"
	"time"
)

func doCheckResult(port *serial.Port, result *string) {
	if *result == "ATE1\r\n" {
		_, err := port.Write([]byte("ATE1\r\n\r\n\r\n\r\nOK\r\n"))
		if err != nil {
			fmt.Println("send err", err)
		}
		*result = ""
	} else if *result == "AT+IMEI\r\n" {
		_, err := port.Write([]byte("AT+IMEI\r\n\r\n\r\n861223063584557\r\n\r\nOK\r\n"))
		if err != nil {
			fmt.Println("send err", err)
		}
		*result = ""
	} else if *result == "AT+CREC\r\n" {
		fmt.Println("write +CREC:0")
		_, err := port.Write([]byte("+CREC:0\r\nOK\r\n"))
		if err != nil {
			fmt.Println("send err", err)
		}
		*result = ""
		time.Sleep(time.Second * 5)
		//录音完成
		fileInfo, _ := os.Stat("sample-15s.wav")
		strNotify := fmt.Sprintf("+CFTRANTX:DATA,%v\r\n", fileInfo.Size())
		fmt.Println("send notify", strNotify)
		_, err = port.Write([]byte(strNotify))
		if err != nil {
			fmt.Println("send err", err)
		}
		time.Sleep(time.Millisecond * 200)
		//发送wav文件

		file, _ := os.Open("sample-15s.wav")
		defer file.Close()
		bufReader := bufio.NewReader(file)
		chunkSize := 1024
		for {
			buf := make([]byte, chunkSize)
			n, err := bufReader.Read(buf)
			if n == 0 {
				break
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			_, err = port.Write(buf[:n])
			if err != nil {
				fmt.Println("send err", err)
			}
		}
		time.Sleep(time.Millisecond * 200)
		//成功通知
		_, err = port.Write([]byte("\r\n+CFTRANTX:0\r\n"))
		if err != nil {
			fmt.Println("send err", err)
		}

	}
}

func reader(port *serial.Port) {
	buf := make([]byte, 128)
	result := ""
	for {
		n, err := port.Read(buf)
		if err != nil {
			fmt.Println("Failed to read from COM21:", err)
			return
		}
		if n > 0 {
			data := buf[:n]
			result += string(data)
		}
		fmt.Println("check", result)
		doCheckResult(port, &result)
	}
}

func main() {
	configCOM21 := &serial.Config{Name: "COM21", Baud: 9600}
	portCOM21, err := serial.OpenPort(configCOM21)
	if err != nil {
		fmt.Println("Failed to open COM21:", err)
		return
	}
	defer portCOM21.Close()

	go reader(portCOM21)
	select {}
}
