package server

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/BarunW/microservices-go/currency-service/data"
	"github.com/BarunW/microservices-go/currency-service/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Currency struct{
    rates *data.ExchangeRates 
    subscriptions map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
    protos.UnimplementedCurrencyServer
}

func (c *Currency) handleUpdates(){
    ru := c.rates.MonitorRates(5 * time.Second)

    outer:
    for range ru{
        // loop over subscribe client 
        for k, v := range c.subscriptions{
            // loop over subscribed rates
            for _, rr := range v{
                rate , err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
                if err != nil{
                    slog.Error("[Unable to update the get  rate]", rr.GetBase().String(), rr.GetDestination().String(), err)
                    break outer           
                }
                errs := k.Send(&protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: rate})            
                if errs != nil{
                    slog.Error("[Unable to update the rate]", rr.GetBase().String(), rr.GetDestination().String(), err)
                    break outer
                }
            }
        }
    }
    defer c.handleUpdates()
}


func NewCurrency(r *data.ExchangeRates) *Currency{
    c := &Currency{
        rates : r, 
        subscriptions:make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest)}

    go c.handleUpdates()
    return c
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error){
    fmt.Println("Hadle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())
    if rr.Base == rr.Destination {
        resp := status.Newf(
            codes.InvalidArgument,
            "Base currency %s can not be the same as the destination currency %s",
            rr.Base.String(),
            rr.Destination.String(),
        ) 
        
        st, err := resp.WithDetails(rr)
        if err != nil{
            return  nil, err
        }
        return nil, st.Err()
    }

    rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
    if err != nil {
        return nil, err
    }

    return &protos.RateResponse{Base: rr.GetBase(),Destination: rr.GetDestination(),Rate: rate}, nil
}

func (c *Currency) SubscribeRates( src protos.Currency_SubscribeRatesServer) error{
    for {
        rr, err := src.Recv()
        if err == io.EOF{
            slog.Info("Unable to read from client", "error", err)
            break
        }
        if err != nil {
            slog.Error("Unable to read from client", "error", err)
            break
        }
        slog.Info("Handle client request",rr)

        rrs, ok := c.subscriptions[src]
        if !ok {
            rrs = []*protos.RateRequest{}
        }

        rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
        if err != nil{
            slog.Error("[Unable to get rate SR]",rate, rr.GetBase().String(), rr.GetDestination().String(), err)
        }else{
            rrs = append(rrs, rr)
        }
        c.subscriptions[src] = rrs
    }
    return nil
}



func (c *Currency) SubscribeRatesV2( src protos.Currency_SubscribeRatesServer) error{
    for {
        rr, err := src.Recv()
        if err == io.EOF{
            slog.Info("Unable to read from client", "error", err)
            break
        }
        if err != nil {
            slog.Error("Unable to read from client", "error", err)
            break
        }
        slog.Info("Handle client request",rr)

        rrs, ok := c.subscriptions[src]
        if !ok {
            rrs = []*protos.RateRequest{}
        }
        // check that subscriptions does not exists
      // check if already in the subscribe list and return a custom gRPC error
		for _, r := range rrs {
			// if we already have subscribe to this currency return an error
			if r.Base == rr.Base && r.Destination == rr.Destination {
				slog.Error("Subscription already active", "base", rr.Base.String(), "dest", rr.Destination.String())

				grpcError := status.New(codes.InvalidArgument, "Subscription already active for rate")
				grpcError, err = grpcError.WithDetails(rr)
				if err != nil {
					slog.Error("Unable to add metadata to error message", "error", err)
					continue
				}

				// Can't return error as that will terminate the connection, instead must send an error which
				// can be handled by the client Recv stream.
				rrs := &protos.StreamingRateResponse_Error{Error: grpcError.Proto()}
				src.Send(&protos.StreamingRateResponse{Message: rrs})
			}
		}
        //rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
        rrs = append(rrs, rr)
        c.subscriptions[src] = rrs
    }
    return nil
}
