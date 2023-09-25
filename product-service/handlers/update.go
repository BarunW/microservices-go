package handlers

import (
	"fmt"
	"github.com/BarnW/microservices-go/product-service/data"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request){
    fmt.Println("update Product")
    // extract id
    sid := mux.Vars(r)["id"]

    id, err_conv := strconv.Atoi(sid)
    if err_conv  != nil{
        http.Error(rw,"Server Error", http.StatusInternalServerError)
        return 
    }
    
    prod := r.Context().Value(KeyProduct{}).(data.Product) 

    err := p.productDB.UpdateProduct(id, &prod)
    
    if err == data.ErrorProductNotFound {
        http.Error(rw, "Product Not Found", http.StatusNotFound)
        return
    }
    
    if err != nil {
        fmt.Println("IError",err)
        http.Error(rw, "Server Error", http.StatusInternalServerError)
        return  
    }
    return
}
