package handlers

import (
	"net/http"

	"github.com/BarnW/microservices-go/product-service/data"
)

// swagger:route GET /products listProducts
// Returns a list of products
// response:
// 200: productsResponse


type GenericErrorMessage struct{
    Message error
}


func(p *Products)  GetProducts(rw http.ResponseWriter, r *http.Request){
    
    rw.Header().Add("Content-Type","application/json")
    lp := data.Products{}
    cur := r.URL.Query().Get("cur")
    // getRate
    prods, err := p.productDB.GetProducts(cur)
    if err != nil{
        rw.WriteHeader(http.StatusInternalServerError)
        lp.ToJSON(&GenericErrorMessage{Message: err}, rw)
        return
    }
    // Serialize to JSON
    err = lp.ToJSON(prods,rw)
    if err != nil { 
        p.l.Error("Unable to serialize the prod", err)
        return
    }
}




