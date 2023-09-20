package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// swagger:route DELETE/{id} products deleteProduct
// Returns a list of products
// response:
// 201: No content

// Delete a product form the  database
func(p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request){
    // this will always convert because of the router 
    sid := mux.Vars(r)["id"]

    id, err := strconv.Atoi(sid)
    if err != nil {
        fmt.Println("Error in id", err)
        http.Error(rw,"there is error in your id", http.StatusInternalServerError)

    }

    fmt.Println(id) 
    return
}
