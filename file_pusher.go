package main

import (
  "bytes"
  "fmt"
  "io"
  "io/ioutil"
  "log"
  "mime/multipart"
  "net/http"
  "os"
  "path/filepath"
)

func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
  file, err := os.Open(path)
  if err != nil {
      return nil, err
  }
  defer file.Close()

  body := &bytes.Buffer{}
  writer := multipart.NewWriter(body)
  part, err := writer.CreateFormFile(paramName, filepath.Base(path))
  if err != nil {
      return nil, err
  }
  _, err = io.Copy(part, file)

  for key, val := range params {
      _ = writer.WriteField(key, val)
  }
  err = writer.Close()
  if err != nil {
      return nil, err
  }

  req, err := http.NewRequest("POST", uri, body)
  req.Header.Set("Authorization", "Token " + os.Getenv("TOKEN"))
  req.Header.Set("Content-Type", writer.FormDataContentType())
  return req, err
}

func main() {
  full_path := os.Args[1]

  files, _ := ioutil.ReadDir(full_path)

  rf, err := os.Create(os.Args[2])
  if err != nil {
    panic(err)
  }

  defer rf.Close()

  for _, f := range files {
    fmt.Println(f.Name())

    path := full_path + f.Name()

    request, err := newfileUploadRequest(os.Args[3], map[string]string{}, "upload", path)

    if err != nil {
      log.Fatal(err)
    }

    client := &http.Client{}
    resp, err := client.Do(request)

    if err != nil {
      log.Fatal(err)
    } else {
      if _, err = rf.WriteString(f.Name() + ";"); err != nil {
        panic(err)
      }
      io.Copy(rf, resp.Body)

      resp.Body.Close()
      fmt.Println(resp.StatusCode)
      fmt.Println(resp.Header)
    }
  }
}
