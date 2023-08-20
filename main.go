package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func writeJson(rw http.ResponseWriter, status int, v any) error {

	rw.Header().Add("Content-type", "application/json")
	rw.WriteHeader(status)
	return json.NewEncoder(rw).Encode(v)
}

func main() {

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		scanner := bufio.NewReader(r.Body)
		header := r.Header
		fmt.Println(header)

		lineByte, _, err := scanner.ReadLine()

		if err != nil {
			log.Fatal("err", err)
		}
		line := string(lineByte)

		fmt.Fprintf(rw, "Hello %s", string(line))
		writeJson(rw, http.StatusAccepted, "message:accepted")
	})

	http.ListenAndServe(":4000", nil)
}
