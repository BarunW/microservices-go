package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"time"

	"github.com/BarunW/microservices-go/currency-service/protos"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
    Price float64        `json:"price" validate:"gt=0"` 

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

type ProductsDB struct{
    currency protos.CurrencyClient
    log *slog.Logger
    rates map[string]float64
    client protos.Currency_SubscribeRatesClient
}

func NewProductDB(c protos.CurrencyClient, l *slog.Logger) *ProductsDB{
    pb :=  &ProductsDB{c,l,make(map[string]float64),nil}
    go pb.handleUpdates()

    return pb
}


func (p *ProductsDB) handleUpdates(){
    sub, err := p.currency.SubscribeRates(context.Background())
    if err != nil{
        slog.Error("Unable to subscribe for rates", "error", err)
    }
    p.client = sub
    for {
        rr, err := sub.Recv()
        if err != nil {
            slog.Error("Error recieving message", "error",err)
            return
        }
        p.rates[rr.Destination.String()] = rr.Rate
    }
}

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

// Object to JSON convertor
func (p *Products) ToJSON(i interface{},  w io.Writer) error{
    e := json.NewEncoder(w)
    return e.Encode(i)
}

// JSON DECODER
func (p *Product) FromJSON(r io.Reader) error{
    d := json.NewDecoder(r)
    return d.Decode(p)    
}

func (p *ProductsDB) getRequestRate(dest string) (float64,error){ 
    // if cached return 
    if r, ok := p.rates[dest]; ok {
        slog.Info("[Cache hit]")
        return r, nil
    }

    rr := &protos.RateRequest{
        Base: protos.Currencies(protos.Currencies_value["INR"]),
        Destination: protos.Currencies(protos.Currencies_value[dest]),
    }

    // get initial rate 
    resp, err := p.currency.GetRate(context.Background(), rr)
    if err != nil{
        if s, ok := status.FromError(err); ok {
            md := s.Details()[0].(*protos.RateRequest)
            if s.Code() == codes.InvalidArgument{
                return -1, err//fmt.Errorf("Base %s and Destination %s are same[Unable to get the Data]",md.Base.String(), md.Destination.String())
            }
            return -1, fmt.Errorf("Unable to get rate from currency server, Base %s  Destination %s",md.Base.String(), md.Destination.String())
        }
        
    }
    p.rates[dest] = resp.Rate

    // subscribe for updates 
    p.client.Send(rr)

    return resp.Rate, err
}

// Return the list of product 
func(p *ProductsDB) GetProducts(currency string) (Products, error){
    if currency == ""{
        return productList, nil
    }
    r, err := p.getRequestRate(currency) 
    if err != nil{
        slog.Error("[Error while getting RATE] ", err)
        return nil, err
    } 
    pr := Products{}

    for _, p := range productList{
        np := *p
        np.Price = np.Price * r
        pr = append(pr, &np)
    }
    return pr, nil

}

// Add Product 
func AddProduct(p *Product){
    p.ID = getNextId()
    productList = append(productList,p) 
}


// Update Product 
func(p *ProductsDB)UpdateProduct(id int, prod  *Product) error{
    _, pos, err:= p.findProductById(id)
    if err != nil{
       return err
    }
    prod.ID = id
    productList[pos] = prod
    return nil
}


func(p *ProductsDB) findProductById(id int) (*Product, int, error){
    if id < 1 {
        return  nil, id, ErrorProductNotFound
    } 
    for i, p := range productList{
        if p.ID == id{
            return p, i, nil
        }
    }
    return nil, -1, ErrorProductNotFound

}

func(p *ProductsDB) GetProductById(id int, currency string) (*Product, int, error){
    pr, _, err := p.findProductById(id)
    if err != nil{ 
        return nil, -1, ErrorProductNotFound
    }

    r, err := p.getRequestRate(currency)
    if err != nil {
        return nil, -1, fmt.Errorf("Currency Not Found")
    }
    
    prod := Product{}
    prod = *pr
    prod.Price = pr.Price * r
    return &prod,0, nil
    
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
        Price: 399,
        SKU: "abc323",
        CreatedOn: time.Now().UTC().String(),
        UpdatedOn: time.Now().UTC().String(),
    },
    &Product{
        ID: 2,
        Name: "Espresso",
        Description: "Short  and strong coffee without milk",
        Price: 499,
        SKU: "fjd34",
        CreatedOn: time.Now().UTC().String(),
        UpdatedOn: time.Now().UTC().String(),
        
    },

}
