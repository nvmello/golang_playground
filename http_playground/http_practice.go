package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my Homepage!\n")
}
func getYo(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /yo request\n")
	io.WriteString(w, "Yooo, HTTP!\n")
}

func main() {
    fmt.Println("Setting up routes...")
    http.HandleFunc("/", getRoot)
    http.HandleFunc("/yo", getYo)
    port := 3333

    fmt.Printf("Starting server on port %d...\n", port)
    err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
    if errors.Is(err, http.ErrServerClosed) {
        fmt.Printf("server closed\n")
    } else if err != nil {
        fmt.Printf("error starting server: %s\n", err)
        os.Exit(1)
    }
}

