package files

import (
	
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Local struct{
    maxFileSize int
    basePath string
}

func NewLocal(basePath string, maxSize int) (*Local, error){
    p, err := filepath.Abs(basePath)
    if err != nil{
        return nil, err
    }
    return &Local{basePath: p}, nil
}


func (l *Local) Save(path string, contents io.Reader) error{
    fmt.Println("uploading files")
    fp := l.fullPath(path)
    
    // get the directory path and make sure it exists
    d :=  filepath.Dir(fp)
    err := os.MkdirAll(d, os.ModePerm)
    if err != nil{
        return err
    }
    
    //  if the file exists delete it 
    _, err = os.Stat(fp)
    if err == nil{
        err := os.Remove(fp)
        if err != nil{
            return err
        }
    } else if !os.IsNotExist(err){
        return fmt.Errorf("Unable to get the file info: %w", err) 
    }

    f, err := os.Create(fp)
    if err != nil{
        return fmt.Errorf("%w",err)
    }

    defer f.Close()

    // write the contents to the new file
    // ensure that we are not writing greater than max bytes 
    a, err := io.Copy(f, contents) 

    if err != nil {
        return fmt.Errorf("Unable to write to file: %w", err)
    }
    fmt.Printf("%d bytes has written",a)
    return nil

}

func (l *Local) Get(path string) (*os.File, error){
    // get the full path for the file 
    fp :=l.fullPath(path)
    
    // open the file 
    f, err := os.Open(fp)

    if err != nil {
        return nil, err
    }

    return f, nil
}

func (l *Local) fullPath(path string) string{
    // append the given path to the base fullPath
    return filepath.Join(l.basePath,path)
}
