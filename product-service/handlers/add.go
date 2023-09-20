package handlers

import(
    "net/http"
    "microservice/data"
)

func (p *Products) AddProduct(rw http.ResponseWriter, r*http.Request){
    prod := r.Context().Value(KeyProduct{}).(data.Product)
    data.AddProduct(&prod)
}
