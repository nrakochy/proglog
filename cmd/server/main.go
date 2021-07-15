package main

import (
	"log"
	"strings"

	"github.com/nrakochy/proglog/internal/server"
)

func main(){
	const port = ":9000"
	srv := server.HTTPServer(port)
	var str strings.Builder
	str.WriteString("Starting server on port ")
	str.WriteString(port)
	log.Println(str.String())
	log.Fatal(srv.ListenAndServe())
}
