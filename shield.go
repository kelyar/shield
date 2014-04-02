package shield

import (
  "appengine"
  "appengine/urlfetch"
  "appengine/blobstore"
  "io/ioutil"
  "fmt"
  "net/http"
)

var (
    StorageURL = "https://storage.googleapis.com"
    Bucket  = "wixstaticdev"
    ComputeEngineHost = "http://static.gce.wixstatic.com"
)

func init() {
    http.HandleFunc("/exif/", ExifHandler)
    http.HandleFunc("/", handler)
}

func blobFileName(fileName string) string {
    return "/gs/" + Bucket + fileName
}

func FileUrl(fileName string) string {
    return StorageURL + "/" + Bucket + fileName
}

func GetAndRender(path string, c appengine.Context, w http.ResponseWriter, r *http.Request) error {
    client := urlfetch.Client(c)
    resp, err := client.Get(ComputeEngineHost + path)

    if resp.StatusCode == 404 {
        http.NotFound(w, r)
        return nil
    }
    defer resp.Body.Close()
    if err == nil {
        if body, err := ioutil.ReadAll(resp.Body); err == nil {
            w.Header().Set("Content-Type", "image/jpeg")
            fmt.Fprintf(w, "%s", body) // io.Copy?
            return nil
        }
    }
    return err
}

func HandleError(w http.ResponseWriter, c appengine.Context, err error) {
    if err.Error() == "404" {
        http.Error(w, "404 page not found", http.StatusNotFound)
    } else {
        c.Errorf(err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func RespondWithHeader(url string, c appengine.Context, w http.ResponseWriter, r *http.Request) error {
    w.Header().Set("Content-Type", "image/jpeg")
    blobKey, err := blobstore.BlobKeyForFile(c, blobFileName(url))
    if err == nil {
        w.Header().Set("X-AppEngine-BlobKey", string(blobKey))
        fmt.Fprintln(w, "")
    }
    return err
}

func handler(w http.ResponseWriter, r *http.Request) {
    // "/media/ce7c31_f0e70d3996554b4cfeff3d19aa05739b.jpg_srz_170_150_75_22_0.5_1.20_0.00_jpg_srz"
    imagePath := r.URL.Path

    if imagePath == "/favicon.ico" {
        http.NotFound(w, r)
        return
    }

    c := appengine.NewContext(r)
    client := urlfetch.Client(c)

    resp, err := client.Head(FileUrl(imagePath))
    if err != nil {
        HandleError(w, c, err)

    } else if resp.StatusCode == 200 { // image exists
        if err = RespondWithHeader(imagePath, c, w, r); err != nil {
            HandleError(w, c, err)
        }

    } else if resp.StatusCode == 404 { // no image, do request to compute engine
        if err = GetAndRender(imagePath, c, w, r); err != nil {
            HandleError(w, c, err)
        }
    }
}
