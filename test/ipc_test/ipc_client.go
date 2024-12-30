package main

import (
	"fmt"
	"github.com/tarm/serial"
	"strings"
)

func doCheckResult(port *serial.Port, result *string) {
	fmt.Println(*result)
	if *result == "sysinfo test" {
		_, err := port.Write([]byte("<ACK> 200 OK"))
		if err != nil {
			fmt.Println("send err", err)
		}
		*result = ""
	} else if strings.Contains(*result, "sysinfo setuid") {
		//_, err := port.Write([]byte("<ACK> 400 Unknown command"))
		_, err := port.Write([]byte("<ACK> 200 OK"))
		fmt.Println("send <ACK> 400 Unknown command")
		if err != nil {
			fmt.Println("send err", err)
		}
		*result = ""
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
