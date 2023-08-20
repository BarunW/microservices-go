package handlers

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
)

type hello struct {
	l *log.Logger
}

func (*hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	scanner := bufio.NewReader(r.Body)
	header := r.Header
	fmt.Println(header)

	lineByte, _, err := scanner.ReadLine()

	if err != nil {
		http.Error(rw, "Opps", http.StatusBadRequest)
	}
	line := string(lineByte)

	fmt.Fprintf(rw, "Hello %s", string(line))
}
