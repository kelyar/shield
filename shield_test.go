package shield_test

import (
  "appengine/aetest"
  "fmt"
  "net/http"
  "net/http/httptest"
  "shield"
  "testing"
)

func TestRespondWithHeader(t *testing.T) {

    jpgok := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "HELLO")
        w.Header().Set("Content-Type", "application/jpeg")
        fmt.Fprintln(w, "iamjpg")
    }))
    defer jpgok.Close()

    shield.StorageURL = jpgok.URL

    c, err := aetest.NewContext(nil)
    if err != nil {
        t.Fatal(err)
    }
    defer c.Close()

    w := httptest.NewRecorder()
    req, err := http.NewRequest("GET", "http://example.com/media/test.jpg", nil)
    if err != nil {
        t.Fatal(err)
    }

    if err = shield.RespondWithHeader("/media/test.jpg", c, w, req); err != nil {
        t.Fatal(err)
    }
    if "encoded_gs_file:d2l4c3RhdGljZGV2L21lZGlhL3Rlc3QuanBn" != w.Header().Get("X-AppEngine-BlobKey") {
        t.Fatal("expected encoded_gs_file, got %v", w.Header().Get("X-AppEngine-BlobKey"))
    }
    if err = shield.GetAndRender(w, req, "/media/test.jpg"); err != nil {
        t.Fatal(err)
    }
}
