package server

import (
	"context"
	"github.com/BarunW/microservices-go/currency-service/protos"
	"fmt"
)

type Currency struct{
    protos.UnimplementedCurrencyServer
}


func (Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error){
    fmt.Println("Hadle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

    return &protos.RateResponse{Rate: 0.5}, nil
}


