package handlers

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
    "github.com/BarunW/microservices-go/productimage/files"
	//"strconv"
	"github.com/gorilla/mux"
)

// Files is a handler for reading and writing Files
type Files struct{
    l log.Logger
    store files.Storage  
}

// NewFiles creates a new file handler 
func NewFiles(s files.Storage) *Files {
    return &Files{store:s }
}

// ServeHTTP implements the http.Handler interface 
func (f *Files) ServeHTTP(rw http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    id := vars["id"]
    fn := vars["filename"]
    f.SaveFile(id, fn, rw, r.Body)
}

// MultiPart Upload Files
func (f  *Files) UploadMultipart(rw http.ResponseWriter, r *http.Request){
    fmt.Println("Multipart")
    err := r.ParseMultipartForm(128 * 1024) 

    if err != nil {
        f.l.Print("Bad request cannot parse the file")
        http.Error(rw, "Expected multipart form", http.StatusBadRequest)
        return
    }

    id := r.FormValue("id") 
    
    fl , mh, err := r.FormFile("file")
    
    if err != nil {
        http.Error(rw, "Missing File", http.StatusBadRequest)
        return
    }

    f.SaveFile(id,mh.Filename, rw, fl)    
}

func (f *Files) invalidURI(uri string, rw http.ResponseWriter){
    log.Print(fmt.Errorf("invalid %s uri", uri))
    http.Error(rw, "Invalid file path should be in the format: /[id]/[filepath]", http.StatusBadRequest)
}

// saveFile save the contents of the request to a file 
func(f *Files) SaveFile(id, path string, rw http.ResponseWriter,r io.ReadCloser ){
    read := bufio.NewReader(r)
    n := read.Size()

    fmt.Println(n)
    fp := filepath.Join(id,path)  

    err := f.store.Save(fp, r)
    if err != nil {
        fmt.Println("Unable to save file", "error", err)
        http.Error(rw, "Unable to save the file", http.StatusInternalServerError)
    }
}
