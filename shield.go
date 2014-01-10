package shield

import (
  "appengine"
  "appengine/urlfetch"
  "appengine/blobstore"
  //"appengine/user"
  //"code.google.com/p/google-api-go-client/compute/v1beta16"
  //"code.google.com/p/google-api-go-client/storage/v1beta2"

  "fmt"
  "html"
  "net/http"
)

const (
  StorageURL = "https://storage.googleapis.com"
  Bucket  = "wixstaticdev"
)

func init() {
  http.HandleFunc("/", handler)
}

func blobFileName(fileName string) string {
  return "/gs/" + Bucket + fileName
}

func fileUrl(fileName string) string {
  return StorageURL + "/" + Bucket + fileName
}

func handler(w http.ResponseWriter, r *http.Request) {
    // "/media/ce7c31_f0e70d3996554b4cfeff3d19aa05739b.jpg_srz_170_150_75_22_0.5_1.20_0.00_jpg_srz"
    imagePath := html.EscapeString(r.URL.Path)

    if imagePath == "/favicon.ico" {
      http.NotFound(w,r)
      return
    }

    c := appengine.NewContext(r)
    client := urlfetch.Client(c)

    resp, err := client.Head(fileUrl(imagePath))
    if err != nil {
      c.Errorf(err.Error())
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return

    } else if resp.StatusCode == 200 { // image exists
      w.Header().Set("Content-Type", "image/jpeg")
      blobKey, err := blobstore.BlobKeyForFile(c, blobFileName(imagePath))
      if err == nil {
        w.Header().Set("X-AppEngine-BlobKey", string(blobKey))
      }
      fmt.Fprintln(w, "")

    } else if resp.StatusCode == 404 { // no image, do request to compute engine?
    }
}
