package files_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"productimage/files"
	"testing"
)


func setupLocal(t *testing.T) (*files.Local, string, func()){
    dir := os.TempDir()    
    l, err := files.NewLocal(dir,5 )

    if err != nil {
        t.Fatal(err)
    }

    return l, dir, func() {
        os.Remove(dir)
    }
}

func TestLocal( t *testing.T){
    savePath := "/1/test.jpg"
    fileContents := "Hello world"

    l, dir, cleanup := setupLocal(t)
    defer cleanup()
    
    err := l.Save(savePath, bytes.NewBuffer([]byte(fileContents)))

    if err != nil {
        t.Fatal(err)
    }

    // check the file has been corectly written
    f, err := os.Open(filepath.Join(dir,savePath))

    //check the contents of the file
    d, err := io.ReadAll(f)
    if err != nil {
        t.Fatal(err)
    }
    
    fmt.Println(string(d))

}
