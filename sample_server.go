package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Receieved request")
	fmt.Fprintf(w, "hello\n")
	fmt.Println("wrote response")
}

func headers(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Receieved request for headers")
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
	fmt.Println("Wrote response for headers")
}
func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	http.ListenAndServe(":8090", nil)
}
