package main

import (
	"fmt"
	"github.com/epimorphics/prometheus-sns-webhook/pkg/server"
	"net/http"
)

func main() {
	router := server.NewRouter()
	fmt.Println(http.ListenAndServe(":3000", router))
}
