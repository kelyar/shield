package shield

import (
  "appengine"
  "appengine/urlfetch"
  //"appengine/user"
  //"code.google.com/p/google-api-go-client/compute/v1beta16"
  //"code.google.com/p/google-api-go-client/storage/v1beta2"

  "fmt"
  "net/http"
)

const (
  GCS_URL = "http://static.wix.com"
)

func init() {
  http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
    // "/media/ce7c31_f0e70d3996554b4cfeff3d19aa05739b.jpg_srz_170_150_75_22_0.5_1.20_0.00_jpg_srz"
    image_url := r.URL.String()

    c := appengine.NewContext(r)
    client := urlfetch.Client(c)
    resp, err := client.Head(GCS_URL + image_url)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return

    } else if resp.StatusCode == 200 { // image exists
      w.Header().Set("Content-Type", "image/jpeg")
      w.Header().Set("X-AppEngine-BlobKey", "here will be blobkey")
      fmt.Fprintln(w, "")

    } else if resp.StatusCode == 404 { // no image, do request to compute engine?
    }
}
