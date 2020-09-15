package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)

	yamlHandler, err := YAMLHandler(mux)
	if err != nil {
		panic(err)
	}

	fmt.Print("Listening port :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
}
