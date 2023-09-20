package main

import (
	//	"encoding/json"
	"context"
	"fmt"
	"log"
	"microservice/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"
	"github.com/gorilla/mux"
    "github.com/go-openapi/runtime/middleware"
)

func main() {

    l := log.New(os.Stdout, "product-api", log.LstdFlags)

    // create the handler
    ph:= handlers.NewProducts(l)
    
    // gRPC clients server

    // create a new servemux from gorilla/mux package
    sm:= mux.NewRouter()

    getRouter := sm.Methods("GET").Subrouter()    
    getRouter.HandleFunc("/",ph.GetProducts)
    
    putRouter := sm.Methods("PUT").Subrouter()
    putRouter.HandleFunc("/{id:[0-9]+}",ph.UpdateProduct)
    putRouter.Use(ph.MiddlewareProduct)

    postRouter := sm.Methods(http.MethodPost).Subrouter()
    postRouter.HandleFunc("/",ph.AddProduct)
    postRouter.Use(ph.MiddlewareProduct)
    
    deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
    deleteRouter.HandleFunc("/{id:[0-9]+}",ph.DeleteProduct)
    deleteRouter.Use(ph.MiddlewareProduct)

    opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
    sh := middleware.Redoc(opts, nil)
    getRouter.Handle("/docs",sh)
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
