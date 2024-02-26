package main

import (
	"log"
	"net/http"
	"time"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	muxServer := http.NewServeMux()

	muxServer.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Incoming request")
		response := "Hello World version 3"
		_, err := w.Write([]byte(response))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
		}
	})

	log.Println("Server started in localhost:8080")
	err := http.ListenAndServe("0.0.0.0:8080", muxServer)
	if err != nil {
		log.Fatalln("Error occured:", err.Error())
	}
	time.Sleep(10 * time.Second)
}
