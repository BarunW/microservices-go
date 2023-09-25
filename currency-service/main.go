package main

import (
	"log"
	"log/slog"
	"net"
	"github.com/BarunW/microservices-go/currency-service/data"
	protos "github.com/BarunW/microservices-go/currency-service/protos"
	"github.com/BarunW/microservices-go/currency-service/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main(){
    var sl *slog.Logger
    // grpc server
    gs := grpc.NewServer()
    
    //rate 
    er, err := data.NewRates(sl) 
    if err != nil{
        log.Fatal(err)
    }

    // server[handle the methods from the client]
    cs := server.NewCurrency(er)
    
    //register grpc 
    protos.RegisterCurrencyServer(gs, cs) 

    //register the reflection service which allows clients to determine the methods
    // sever reflection
    reflection.Register(gs)

    l, err := net.Listen("tcp", ":9002")
    if err != nil {
        log.Fatal("Unable to listen error", err)
    }
        
    slog.Info("[Currency Service Serving on :9002]")
    gs.Serve(l)
}

