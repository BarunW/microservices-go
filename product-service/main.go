package main

import (
	//	"encoding/json"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/BarnW/microservices-go/product-service/handlers"
	"github.com/BarunW/microservices-go/currency-service/protos"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

    l := log.New(os.Stdout, "product-api", log.LstdFlags)
    
    // gRPC coonect to currency server
    conn, err := grpc.Dial("localhost:9002", grpc.WithTransportCredentials(insecure.NewCredentials()))
    defer conn.Close()
    if err != nil{
        panic(err)
    }

    
    // gRPC clients
    cc := protos.NewCurrencyClient(conn)
    
    // create the handler
    ph:= handlers.NewProducts(&slog.Logger{},cc)
    
    // create a new servemux from gorilla/mux package
    sm:= mux.NewRouter()

    // GET product
    getRouter := sm.Methods("GET").Subrouter()    
    getRouter.HandleFunc("/products",ph.GetProducts).Queries("cur","{[A-Z]{3}}") 
    getRouter.HandleFunc("/products",ph.GetProducts)
    // UPDATE product
    putRouter := sm.Methods("PUT").Subrouter()
    putRouter.HandleFunc("/products/{id:[0-9]+}",ph.UpdateProduct)
    putRouter.Use(ph.MiddlewareProduct)

    // POST product
    postRouter := sm.Methods(http.MethodPost).Subrouter()
    postRouter.HandleFunc("/products/",ph.AddProduct)
    postRouter.Use(ph.MiddlewareProduct)
    
    // DELETE product
    deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
    deleteRouter.HandleFunc("/products/{id:[0-9]+}",ph.DeleteProduct)
    deleteRouter.Use(ph.MiddlewareProduct)

    opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
    sh := middleware.Redoc(opts, nil)
    getRouter.Handle("/products/docs",sh)
    getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./"))) 


    s := &http.Server{
        Addr: ":9000",                   // configure bind adress
        Handler: sm,                     // set the default handler 
        ErrorLog: l,                     // set the logger for the server
        IdleTimeout : 120 * time.Second, // max time for connections using TCP keep-Alive
        ReadTimeout:  5 * time.Second,   // max time for read request from the client
        WriteTimeout: 10 * time.Second,  // max time for to write a respond to the client
    }
    
    // start the server
    go func() {
        slog.Info("[Product-Service running on port 9000]")
        err := s.ListenAndServe()
        if err != nil {
            l.Fatal(err)
        }
    }()
    

    sigChan := make(chan os.Signal)
    signal.Notify(sigChan, os.Interrupt)
    signal.Notify(sigChan,os.Kill)
    
    sig := <- sigChan
    fmt.Println("gracefully shutting down", sig)
    tc, f :=context.WithTimeout(context.Background(), 30*time.Second)
    f()
    s.Shutdown(tc)
}
