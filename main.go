package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/kawabatas/gomock/handler"
)

var (
	portOption = flag.String("p", "8888", "port number")
)

func init() {
	flag.Parse()
}

func main() {
	fmt.Println("Start Mock Server...")
	fmt.Printf("port: %v\n", *portOption)

	http.HandleFunc("/", handler.HandleStub)
	portAddress := ":" + *portOption
	log.Fatal(http.ListenAndServe(portAddress, nil))
}
