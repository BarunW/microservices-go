package data

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
    "math/rand"
)

type ExchangeRates struct {
    log *slog.Logger
    rates map[string]float64
}

func NewRates(l *slog.Logger) (*ExchangeRates, error){
    exrate := &ExchangeRates{log: l, rates: make(map[string]float64)}
    err := exrate.getRates()
    return exrate, err
}


func (e *ExchangeRates) GetRate(base, dest string) (float64, error){
    br, ok := e.rates[base]
    if !ok{
        return 0, fmt.Errorf("Rate not found for currency %s",base)
    }

    dr, ok := e.rates[dest]
    if !ok {
        return 0, fmt.Errorf("Rate not found for currency %s",dest)
    }

    return dr/br, nil
}

func (e *ExchangeRates) MonitorRates(interval time.Duration) chan struct{}{
    ret := make(chan struct{})
    
    go func(){
        ticker := time.NewTicker(interval)
        for {
            select {
                case <-ticker.C:
                    for k, v := range e.rates{
                        change := (rand.Float64()/10)

                        direction := rand.Intn(1)

                        if direction == 0{
                            change = 1 - change
                        } else {
                            change = 1 + change
                        }
                        e.rates[k] = v * change 
                    }

            }
            ret <-struct{}{}
        }
    }()
    return ret 
}

func(e *ExchangeRates) getRates() error{
    resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
    if err != nil {
        return nil
    }
    
    if resp.StatusCode != http.StatusOK{
        return fmt.Errorf("Expected error code 200 got %d", resp.StatusCode)
    }
    defer resp.Body.Close()

    md := &Cubes{}
    xml.NewDecoder(resp.Body).Decode(&md)

    for _, v := range md.CubeData{
        r, err := strconv.ParseFloat(v.Rate,64)
        if err != nil {
            return err
        }
        e.rates[v.Currency] = r
    }
    return nil
}

type Cubes struct {
    CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
    Currency string `xml:"currency,attr"`
    Rate string `xml:"rate,attr"`
}


