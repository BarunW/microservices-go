package handlers

import(
    "net/http"
    "github.com/BarnW/microservices-go/product-service/data"
)

func (p *Products) AddProduct(rw http.ResponseWriter, r*http.Request){
    prod := r.Context().Value(KeyProduct{}).(data.Product)
    data.AddProduct(&prod)
}
