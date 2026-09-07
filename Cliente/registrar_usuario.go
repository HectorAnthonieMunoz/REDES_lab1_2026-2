package main

import (
	//"net"
	"net/http"
	//"sync"
	//"time"
	//"os"
	"fmt"
	//"bufio"
	"strings"
	//"encoding/csv"
)

var clienteHTTP = http.DefaultClient;

func main() {
	url := "http://localhost:8888/register";

	var username,password string;

	println("Ingrese su nombre y contraseña:")
	
	fmt.Scan(&username, &password);

	body := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password)

	//fmt.Println(body)

	r, err := http.NewRequest("POST", url, strings.NewReader(body))
	if (err != nil) {
		panic(err)
	}
	r.Header.Add("Content-Type", "application/json")

	resp, err := clienteHTTP.Do(r);
	if (err != nil) {
		panic(err)
	}

	code := resp.Status;

	fmt.Println(code)

}
