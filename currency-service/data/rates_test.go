package data_test

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/BarunW/microservices-go/currency-service/data"
)

func TestNewRates(t *testing.T){
    tr, err := data.NewRates(&slog.Logger{})  
    if err != nil {
        t.Fatal(err)
    }

    fmt.Printf("%+v",tr)
}
