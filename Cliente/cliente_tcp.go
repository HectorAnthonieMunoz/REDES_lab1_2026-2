package main

import (
	"net"
	//"net/http"
	//"sync"
	//"time"
	"os"
	"fmt"
	//"bufio"
	//"strings"
	//"encoding/csv"
)


func main() {
	fmt.Println("Ingrese sus credenciales:")

	var username,password string;

	fmt.Scan(&username, &password)

	strEcho := fmt.Sprintf("LOGIN %s %s\n", username, password);
	servAddr := "0.0.0.0:9000"
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte(strEcho))
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	println("write to server = ", strEcho)

	reply := make([]byte, 4096)

	_, err = conn.Read(reply)
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	println("reply from server=", string(reply))

	conn.Close()
}
