package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

// Product represents a product in the system
// swagger:model
type Product struct{
    // The unique identifier for the produt 
    // required: false
    // min: 1
    ID int               `json:"id"`             

    // The name of the product.
    //
    // required: true
    //  max length: 255
    Name string          `json:"name" validate:"required"` 

    // Description of the product 
    //
    // required" false
    // max length : 10000
    Description string   `json:"description"` 

    // The price of the product
    //
    // required: true
    // min: 0.01
    Price float32        `json:"price" validate:"gt=0"` 

    // The SKU (Stock Keeping Unit) of the product 
    //
    // required: true
    // pattern: [a-z]+-[a-z]+-[a-z]+
    SKU string           `json:"sku" validate:"required,sku"` 
    
    // CreatedOn respresent the creation date of the product 
    CreatedOn string     `json:"-"` 

    // UpdatedOn respresent the creation date of the product 
    UpdatedOn string     `json:"-"`  


    // DeletedOn respresent the deletion date of the product 
    DeletedOn string     `json:"-"` 

}

var validate  *validator.Validate


type Products []*Product
var ErrorProductNotFound = fmt.Errorf("Product not found")

func validateSKU(fl validator.FieldLevel) bool {
    re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
    matches := re.FindAllString(fl.Field().String(), -1)
    
    if len(matches) != 1{
        return false
    }

    return true
}

func (p *Product) Validate()error{
    validate := validator.New() 
    validate.RegisterValidation("sku", validateSKU)
    return validate.Struct(p)
}


func (p *Products) ToJSON(w io.Writer) error{
    e := json.NewEncoder(w)
    return e.Encode(p)
}

func (p *Product) FromJSON(r io.Reader) error{
    d := json.NewDecoder(r)
    return d.Decode(p)    
}

func GetProducts() Products {
    return productList
}

func AddProduct(p *Product){
    p.ID = getNextId()
    productList = append(productList,p) 
}

func UpdateProduct(id int, p *Product) error{
    _, pos, err:= findProduct(id)
    if err != nil{
       return err
    }
    p.ID = id
    productList[pos] = p
    return nil
}

func findProduct(id int) (*Product, int, error){
    for i, p := range productList{
        if p.ID == id{
            return p, i, nil
        }
    }

    return nil, -1, ErrorProductNotFound
}

func getNextId() int {
    lp := productList[len(productList) - 1] 
    return lp.ID + 1  
}

var productList = []*Product{
    &Product{
        ID: 1,
        Name: "Latee",
        Description: "Frothy milky coferr",
        Price: 2.45,
        SKU: "abc323",
        CreatedOn: time.Now().UTC().String(),
        UpdatedOn: time.Now().UTC().String(),
    },
    &Product{
        ID: 2,
        Name: "Espresso",
        Description: "Short  and strong coffee without milk",
        Price: 1.99,
        SKU: "fjd34",
        CreatedOn: time.Now().UTC().String(),
        UpdatedOn: time.Now().UTC().String(),
        
    },

}
