package server

import (
	"context"
	"currency/protos"
	"fmt"
)

type Currency struct{}

func NewCurrency() *Currency{
    return &Currency{}
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error){
    fmt.Println("Hadle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

    return &protos.RateResponse{Rate: 0.5}, nil
}

func(Currency) mustEmbedUnimplementedCurrencyServer(){}
