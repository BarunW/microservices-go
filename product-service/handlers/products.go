// Package classification of Product API
//
// Documentation for product API
//
//  Schemes: http
//  BasePath :/
//  Version: 0.0.1
//
//  Consumes:
//  - application/json
// 
//  Produces:
//  - application/json
//
// swagger:meta
package handlers
import (
	"context"
	"fmt"
	"log"
    "net/http"
    "github.com/BarnW/microservices-go/product-service/data"
)

// A list of products returns in the response

// A ValidationError is an error that is used when the required input fails validation.
// swagger:response validationError
type ValidationError struct {
    // The error message
    // in: body
    Body struct {
        // The validation message
        //
        // Required: true
        // Example: Expected type int
        Message string
        // An optional field name to which this validation applies
        FieldName string
    }
}

// A productsResponse is a response  of all products 
// swagger:response productsResponse
type productsResponseWrapper struct {   
    // All current products
    // in : body
    Body []struct {
        //Product Id
        //
        // Required: true
        ID int

    }
} 

// swagger:response noContent 
type productsNoContent struct{
     
}

// swagger:parameters deleteProduct
type producIDParameterWrapper struct{
    // The id of the product to delete from the database 
    // in: path
    // required: true
    ID int `json:"id"`
}

// Products is a http.Handler
type Products struct{
    l *log.Logger
}

func NewProducts(l *log.Logger) *Products{
 return &Products{l}
}

func (p *Products) unMarshalJSONBody(rw http.ResponseWriter, r *http.Request ) *data.Product{
     prod := &data.Product{} 
     err := prod.FromJSON(r.Body)
     if err != nil {
        fmt.Println("error --- unmarshalling  ",err)
        http.Error(rw, "Unable to Un-marshal json", http.StatusBadRequest)
        return nil
    }
    fmt.Println("prod about to return")
    return prod
}

type KeyProduct struct{}

func(p *Products) MiddlewareProduct(next http.Handler) http.Handler{
    fmt.Println("middle ware hit ")
    return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) { 
        prod := *p.unMarshalJSONBody(rw, r) 
        err := prod.Validate()
        if err != nil{
            p.l.Println("[ERROR] validating product", err)
            http.Error(rw,fmt.Sprintf("Error wash your hand %s", err), http.StatusBadRequest)
            return
        }
        ctx := r.Context()
        ctx =   context.WithValue(ctx,KeyProduct{},prod)
        next.ServeHTTP(rw,r.WithContext(ctx))
    })
}
