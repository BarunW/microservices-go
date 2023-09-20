package main

import (
	"currency/protos"
	"currency/server"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main(){
    gs := grpc.NewServer()
    cs := server.NewCurrency()
    protos.RegisterCurrencyServer(gs, cs)    
    
    reflection.Register(gs)

    l, err := net.Listen("tcp", ":9002")
    if err != nil {
        log.Fatal("Unable to listen error", err)
    }
    
    gs.Serve(l)
} 
