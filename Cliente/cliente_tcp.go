package main

import (
	"bufio"
	"net"
	//"net/http"
	//"sync"
	"fmt"
	"os"
	"time"

	//"bufio"
	"strings"
	//"encoding/csv"
)

var servAddr string = "0.0.0.0:9000";
//var udpAddr string = "0.0.0.0:9001"
var stop bool = false;

func udp(token string, port string) {
	// Connect to server
	connTo := "0.0.0.0:"+port

	//print(connTo)

	addr, err := net.ResolveUDPAddr("udp", connTo)
	if err != nil {
		println("No se pudo resolver la dirección UDP:", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		println("Conexión UDP fallida:", err)
	}
	defer conn.Close();

	for true {
		if stop {break;}
		message := []byte("HEARTBEAT "+token)
		_, err = conn.Write(message)
		if err != nil {
			println("Envío UDP fallido: %v", err)
			return
		}
		time.Sleep(3*time.Second);
	}
}

func listenToMsg(conn *net.TCPConn) {
	for true {
		if stop {break;}
		msg := make([]byte, 1024)
		conn.Read(msg)
		println(string(msg))
	}
}

func main() {
	fmt.Println("Ingrese sus credenciales:")

	var username,password string;

	fmt.Scan(&username, &password)

	strEcho := fmt.Sprintf("LOGIN %s %s\n", username, password);
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr fallido:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial fallido:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte(strEcho))
	if err != nil {
		println("Escritura no hecha:", err.Error())
		os.Exit(1)
	}

	//println("write to server = ", strEcho)

	reply := make([]byte, 1024)

	_, err = conn.Read(reply)
	if err != nil {
		println("Escritura al servidor fallida:", err.Error())
		os.Exit(1)
	}

	println("reply from server=", string(reply));

	rep, _ := strings.CutSuffix(string(reply[:]), "\n");

	params := strings.Split(rep, " ")

	if params[0] == "ERROR" {
		println("Login fallido con error "+params[1])
		println("Revise que sus credenciales estén ingresadas correctamente")
		os.Exit(1)
	} else {
		token := params[1]
		port, _, _ := strings.Cut(params[2], "\n")
		go udp(token, port)
		go listenToMsg(conn)
		println("Ahora que está conectado, puede mandar mensajes a través de la consola")
		println("También verá los mensajes de otros usuarios conectados al servidor en forma: INCOMING <usuario> <mensaje>")
		for true {
			//TODO: Arreglar bug que hace que solo se pueda enviar un mensaje
			reader := bufio.NewReader(os.Stdin)

			strEcho, _ := reader.ReadString('\n')
			strEcho = strings.TrimSpace(strEcho) //quitamos whitespaces rodeando al string

			query := "MSG "+token+" "+strEcho+"\n";

			_, err = conn.Write([]byte(query))
			if err != nil {
				println("Escritura no hecha:", err.Error())
				os.Exit(1)
			}
			reply := make([]byte, 4096)
			_, err = conn.Read(reply)
			if err != nil {
				println("Escritura al servidor fallida:", err.Error())
				os.Exit(1)
			}
			print(reply)
		}
	}

	stop = true;
	conn.Close()
}
