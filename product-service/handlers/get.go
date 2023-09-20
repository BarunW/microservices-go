package handlers

import (
	"microservice/data"
	"net/http"
)

// swagger:route GET /products listProducts
// Returns a list of products
// response:
// 200: productsResponse

func(p *Products)  GetProducts(rw http.ResponseWriter, r *http.Request){
   p.l.Println("Handle Get Products") 
   
   lp := data.GetProducts()

   // serialize the list to JSON
    err := lp.ToJSON(rw)
    if err != nil{
        http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
    }
}
