package shield

import (
  "appengine"
  "appengine/urlfetch"
  "bytes"
  "io/ioutil"
  "github.com/rwcarlsen/goexif/exif"
  "fmt"
  "net/http"
)

type Wrc struct {
  w *http.ResponseWriter
  r *http.Request
  c appengine.Context
}

func copyFromComputeEngine (path string, wrcobj Wrc) ([]byte, error) {
    client := urlfetch.Client(wrcobj.c)
    resp, err := client.Get(ComputeEngineHost + path)

    if resp.StatusCode == 404 {
        return []byte(""), fmt.Errorf("404")
    }
    defer resp.Body.Close()
    if err == nil {
        if body, err := ioutil.ReadAll(resp.Body); err == nil {
            return []byte(body), nil
        }
    }
    return nil, err
}

func getImageContent(path string, wrcobj Wrc) ([]byte, error) {
    var body []byte
    client := urlfetch.Client(wrcobj.c)

    resp, err := client.Get(FileUrl(path))
    defer resp.Body.Close()
    if err != nil {
        return nil, err
    }

    if resp.StatusCode == 404 {
        body, err = copyFromComputeEngine(path, wrcobj)
        if err != nil {
            return nil, err
        }
    } else {
        body, _ = ioutil.ReadAll(resp.Body)
    }
    return body, nil
}

func ExifHandler(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    imagePath := r.URL.Path[5:len(r.URL.Path)] // remove '/exif' segment

    body, err := getImageContent(imagePath, Wrc{&w, r, c})
    if err == nil {
        x, err := exif.Decode(bytes.NewReader(body))
        if err == nil {
            js, err := x.MarshalJSON()
            if err == nil {
                w.Header().Set("Content-Type", "application/json")
                fmt.Fprintf(w, string(js))
                return
            }
        }
    }
    HandleError(w, c, err)
}
