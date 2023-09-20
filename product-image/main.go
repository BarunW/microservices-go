package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	//"path"
	"productimage/files"
	"productimage/handlers"
	"time"

	"github.com/gorilla/mux"
)


func main(){
    
    // create the storage class, use local storage
    // max filesize 5MB
    stor, err := files.NewLocal("./imagestore", 1024*1000*5)
    path := os.Getenv("BASE_PATH")
    fmt.Println(path)
    if err != nil{
        log.Fatal(err)
    }
    
    // create the handlers 
    fh := handlers.NewFiles(stor)

    // gzip middleware 
    gmw := handlers.GzipHandler{}
    
    // create a new serve mux 
    sm := mux.NewRouter()
    
    // post **Upload Files **
    // filename regex : {filename:[a-zA-Z]+\\.[a-z]{3}}
    ph := sm.Methods(http.MethodPost).Subrouter()
    ph.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{4}}",fh.ServeHTTP)
    // using multipart
    ph.HandleFunc("/",fh.UploadMultipart)

    // get files
    gh := sm.Methods(http.MethodGet).Subrouter()
    gh.Handle(
        "/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{4}}",
        http.StripPrefix("/images/", http.FileServer(http.Dir("./imagestore"))), //os.Getenv("BASE_PATH")))),
    )
    
    gh.Use(gmw.GzipMiddleware)
    s := &http.Server{
        Addr: "0.0.0.0:9001",
        Handler: sm,
        IdleTimeout: 120 * time.Second,
        ReadTimeout: 5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
    go func(){
        log.Print("Starting server")
        err := s.ListenAndServe()
        if err != nil{
            log.Fatal(err)
        }

    }()
    
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt)
    signal.Notify(sigChan, os.Kill)

    sig := <- sigChan
    fmt.Println("Gracefully shutting down", sig)
    tc, cf := context.WithTimeout(context.Background(), 30*time.Second)
    cf()
    s.Shutdown(tc)
}
