package main

import (
	protos "github.com/BarunW/microservices-go/currency-service/protos"
	"github.com/BarunW/microservices-go/currency-service/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
    "log"
	"net"
)

func main(){
    gs := grpc.NewServer()
    cs := server.Currency{}
    protos.RegisterCurrencyServer(gs, cs) 
    reflection.Register(gs)

    l, err := net.Listen("tcp", ":9002")
    if err != nil {
        log.Fatal("Unable to listen error", err)
    }
    
    gs.Serve(l)
}

