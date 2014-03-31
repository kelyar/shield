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

func copyFromComputeEngine (path string, c appengine.Context, w http.ResponseWriter, r *http.Request) ([]byte, error) {
    client := urlfetch.Client(c)
    resp, err := client.Get(ComputeEngineHost + path)

    if resp.StatusCode == 404 {
        http.NotFound(w, r)
        return []byte(""), nil
    }
    defer resp.Body.Close()
    if err == nil {
        if body, err := ioutil.ReadAll(resp.Body); err == nil {
            return []byte(body), nil
        }
    }
    return []byte(""), err
}

func ExifHandler(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    client := urlfetch.Client(c)
    imagePath := r.URL.Path[5:len(r.URL.Path)]
    url := StorageURL + "/" + Bucket + imagePath // remove '/exif'

    resp, err := client.Get(url)
    defer resp.Body.Close()

    if err != nil {
        HandleError(w, c, err)
    }

    var body []byte

    if resp.StatusCode == 404 {
        body, err = copyFromComputeEngine(imagePath, c, w, r)
        if err != nil {
            HandleError(w, c, err)
        }
    } else {
        body, _ = ioutil.ReadAll(resp.Body)
    }

    if x, err := exif.Decode(bytes.NewReader(body)); err == nil {
        if js, err := x.MarshalJSON(); err == nil {
            w.Header().Set("Content-Type", "application/json")
            fmt.Fprintf(w, string(js))
            return
        }
    }
    HandleError(w, c, err)
}
