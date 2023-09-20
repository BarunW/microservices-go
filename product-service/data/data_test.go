package data

import "testing"


func  TestValidation(t *testing.T){
    p := &Product{
        Name: "alu",
        Price: 1.99,
    }

    err := p.Validate()
    
    if err !=nil{
        t.Fatal(err)
    }
}


