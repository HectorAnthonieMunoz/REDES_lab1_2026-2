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
	url := "http://0.0.0.0:9000";

	var username,nombre string;

	fmt.Scan(&username, &nombre);

	body := fmt.Sprintf(`%s %s`, username, nombre)

	//fmt.Println(body)

	r, err := http.NewRequest("LOGIN", url, strings.NewReader(body))
	if (err != nil) {
		panic(err)
	}
	r.Header.Add("Content-Type", "text")

	resp, err := clienteHTTP.Do(r);
	if (err != nil) {
		panic(err)
	}

	code := resp.Status;

	fmt.Println(code)

}
