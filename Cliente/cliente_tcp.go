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
var lock bool = false;

func udp(token string, port string) {
	//Este es el código para el hilo encargado de hacer la conexión UDP
	//y mandar latidos al endpoint periódicamente.
	//Toma como parámetros el token y el port obtenidos del servidor
	//al momento del login.
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

func listenToMsg(conn *net.TCPConn, c chan string) {
	//Hilo que escucha perpetuamente al servidor UDP hasta terminar la conexión
	//para recibir mensajes de otros usuarios y también pasar las respuestas
	//de las queries del cliente de vuelta al hilo principal a través del canal c
	//que es tomado como parámetro. También se toma como parámetro la conexión TCP
	//propiamente tal.
	for true {
		if stop {break;}
		//if lock {continue;}
		msg := make([]byte, 1024)
		conn.Read(msg)

		strMsg := string(msg[:])
		strMsg, _, _ = strings.Cut(strMsg, "\n")

		splitMsg := strings.SplitN(strMsg, " ", 3)

		if (splitMsg[0] != "INCOMING") {
			c <- string(strMsg)
		} else {
			println(splitMsg[1]+": "+splitMsg[2])
		}
	}
}

func main() {
	//Función principal del programa. Pedimos las credenciales de login al usuario
	//antes de intentar conectar.
	fmt.Println("Ingrese sus credenciales:")

	var username,password string;

	fmt.Scan(&username, &password)

	// Intento de login al servidor
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

	// Se hace el login request
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

	//println("reply from server=", string(reply));

	rep, _ := strings.CutSuffix(string(reply[:]), "\n");

	params := strings.Split(rep, " ")

	//Si el servidor devolvió error, paramos. Caso contrario, pasamos al programa principal del cliente
	if params[0] == "ERROR" {
		println("Login fallido con error "+params[1])
		println("Revise que sus credenciales estén ingresadas correctamente")
		os.Exit(1)
	} else {
		//Código principal del cliente mientras está conectado al servidor
		canal := make(chan string)

		token := params[1]
		port, _, _ := strings.Cut(params[2], "\n")
		go udp(token, port) //Iniciamos los goroutines (en esencia, hilos hijos)
		go listenToMsg(conn, canal)
		println("Ahora que está conectado, puede mandar mensajes a través de la consola")
		println("También verá los mensajes de otros usuarios conectados al servidor en forma <usuario>: <mensaje>")
		println("Para cerrar sesión, simplemente escriba LOGOUT (todo mayúscula)")
		for !stop {
			//print(username+": ")
			reader := bufio.NewReader(os.Stdin)

			strEcho, _ := reader.ReadString('\n')

			strEcho = strings.TrimSpace(strEcho) //quitamos whitespaces rodeando al string

			//Si la query es logout, hacerlo en este bloque
			if (strEcho == "LOGOUT") {
				query := "LOGOUT\n";
				_, err = conn.Write([]byte(query))
				if err != nil {
					println("Escritura no hecha:", err.Error())
					os.Exit(1)
				}
				stop = true;
				break;
			}

			//Construcción del mensaje TCP a enviar
			query := "MSG "+token+" "+strEcho+"\n";

			_, err = conn.Write([]byte(query))
			if err != nil {
				println("Escritura no hecha:", err.Error())
				os.Exit(1)
			}
			reply := <- canal; //Obtenemos la respuesta a través del canal de escucha
			replyParams := strings.Split(reply, " ");
			if (replyParams[0] == "ACK") { //Leemos el reply, si empieza con ACK estamos bien...
				println("Servidor acusa recibo de mensaje exitosamente")
			} else if (replyParams[0] == "ERROR"){ //...caso contrario, hay error y debemos manejarlo de acuerdo al caso
				errMsg := replyParams[1];
				if (errMsg == "UNKOWN_COMMAND") {
					println("Se envió un comando inválido al servidor")
				} else { //Sesión inválida o expirada, en ambos casos hay que detener la conexión
					println("La sesión ha expirado o no es válida")
					stop = true;
				}
			}
		}
	}

	stop = true;
	conn.Close()
}
