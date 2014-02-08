package shield

import (
  "appengine"
  "appengine/urlfetch"
  "appengine/blobstore"
  "io/ioutil"
  "fmt"
  "html"
  "net/http"
)

const (
    StorageURL = "https://storage.googleapis.com"
    Bucket  = "wixstaticdev"
    ComputeEngineHost = "http://static.gce.wixstatic.com"
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

func getAndRender(w http.ResponseWriter, c appengine.Context, path string) error {
    client := urlfetch.Client(c)
    resp, err := client.Get(ComputeEngineHost + path)
    defer resp.Body.Close()
    if err == nil {
        if body, err := ioutil.ReadAll(resp.Body); err == nil {
            w.Header().Set("Content-Type", "image/jpeg")
            fmt.Fprintf(w, "%s", body) //serveContent ?
            return nil
        }
    }
    return err
}

func handleError(w http.ResponseWriter, c appengine.Context, err error) {
    c.Errorf(err.Error())
    http.Error(w, err.Error(), http.StatusInternalServerError)
}

func handler(w http.ResponseWriter, r *http.Request) {
    // "/media/ce7c31_f0e70d3996554b4cfeff3d19aa05739b.jpg_srz_170_150_75_22_0.5_1.20_0.00_jpg_srz"
    imagePath := html.EscapeString(r.URL.Path)

    if imagePath == "/favicon.ico" {
        http.NotFound(w, r)
        return
    }

    c := appengine.NewContext(r)
    client := urlfetch.Client(c)

    resp, err := client.Head(fileUrl(imagePath))
    if err != nil {
        handleError(w, c, err)
        return

    } else if resp.StatusCode == 200 { // image exists
        w.Header().Set("Content-Type", "image/jpeg")
        blobKey, err := blobstore.BlobKeyForFile(c, blobFileName(imagePath))
        if err == nil {
          w.Header().Set("X-AppEngine-BlobKey", string(blobKey))
        }
        fmt.Fprintln(w, "")

    } else if resp.StatusCode == 404 { // no image, do request to compute engine
        if err = getAndRender(w, c, imagePath); err != nil {
            handleError(w, c, err)
        }
    }
}
