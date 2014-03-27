package shield

import (
  "appengine"
  "appengine/urlfetch"
  "bytes"
  "io/ioutil"
  "encoding/json"
  "github.com/rwcarlsen/goexif/exif"
  "fmt"
  "net/http"
)

type exifResponse struct {
  Make string     `json:"make"`
  Model string    `json:"model"`
  Focal string    `json:"focal"`
  Iso string      `json:"iso"`
  Aperture string `json:"aperture"`
  Shutter string  `json:"speed"`
  Size string
}

func ExifHandler(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    client := urlfetch.Client(c)
    url := StorageURL + "/" + Bucket + r.URL.Path[5:len(r.URL.Path)] // remove '/exif'
    resp, err := client.Get(url)

    if resp.StatusCode == 404 {
        http.NotFound(w, r)
        return
    }
    defer resp.Body.Close()
    if err == nil {
        body, _ := ioutil.ReadAll(resp.Body)
        x, err := exif.Decode(bytes.NewReader(body))
        if err != nil {
          c.Debugf("DECODE FAILED")
          HandleError(w, c, err)
          return
        }

        maker,err := x.Get(exif.Make)
        model,err := x.Get(exif.Model)
        iso,  err := x.Get(exif.ISOSpeedRatings)
        focal,err := x.Get(exif.FocalLength)
        apert,err := x.Get(exif.ApertureValue)
        //speed,err := x.Get(exif.FNumber)
        //if err == nil { speedStr = speed.StringVal() }

        js, err := json.Marshal(&exifResponse{
            Make: maker.String(),
            Model: model.String(),
            Iso: iso.String(),
            Focal: focal.String(),
            Aperture: apert.String(),
            //Shutter: speedStr,
        })
        if err != nil {
          HandleError(w, c, err)
          return
        }
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprintf(w, string(js))
    }
}
