package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

 type GzipHandler struct{

 }

 func (g *GzipHandler) GzipMiddleware(next http.Handler) http.Handler{
     return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) { 
         if strings.Contains(r.Header.Get("Accept-Enconding"), "gzip"){
             wrw := NewResponseWriter(rw)
             wrw.Header().Add("Content-Endcoding", "gzip")
             next.ServeHTTP(wrw, r) 
             defer wrw.Flush()
             return
         }

         next.ServeHTTP(rw, r)
        
     })
 }


 type WrapperResponseWriter struct {
    rw http.ResponseWriter
    gw *gzip.Writer
 }

 func NewResponseWriter(rw http.ResponseWriter) *WrapperResponseWriter{
    gw := gzip.NewWriter(rw) 
    return &WrapperResponseWriter{rw:rw, gw: gw}

 }

 func (wr *WrapperResponseWriter) Header() http.Header{
    return wr.rw.Header()
 }

 func (wr *WrapperResponseWriter) Write(d []byte) (int, error) {
    return wr.gw.Write(d)
 }

 func (wr *WrapperResponseWriter) WriteHeader(statusCode int ) {
    wr.rw.WriteHeader(statusCode)
 }

 func (wr *WrapperResponseWriter) Flush() {
    wr.gw.Flush()
    wr.gw.Close()
 }

