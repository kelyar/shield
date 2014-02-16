package shield_test

import (
  "appengine"
  "appengine/aetest"
  "fmt"
  "net/http"
  "net/http/httptest"
  "shield"
  "testing"
)

type Wrc struct {
  w http.ResponseWriter
  r *http.Request
  c appengine.Context
}

func NewWrc(t *testing.T) Wrc {
    c, err := aetest.NewContext(nil)
    if err != nil {
        t.Fatal(err)
    }

    w := httptest.NewRecorder()

    req, err := http.NewRequest("GET", "http://example.com/media/test.jpg", nil)
    if err != nil {
        t.Fatal(err)
    }

    return Wrc{w, req, c}
}

func TestRespondWithHeader(t *testing.T) {
    wrc := NewWrc(t)

    if err := shield.RespondWithHeader("/media/test.jpg", wrc.c, wrc.w, wrc.r); err != nil {
        t.Fatal(err)
    }
    if "encoded_gs_file:d2l4c3RhdGljZGV2L21lZGlhL3Rlc3QuanBn" != wrc.w.Header().Get("X-AppEngine-BlobKey") {
        t.Fatal("expected encoded_gs_file, got %v", wrc.w.Header().Get("X-AppEngine-BlobKey"))
    }
    if "image/jpeg" != wrc.w.Header().Get("Content-Type") {
        t.Fatal("expected image/jpeg, got %v", wrc.w.Header().Get("Content-Type"))
    }
}

func TestGetAndRender(t *testing.T) {
    wrc := NewWrc(t)

    jpg200 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "HELLO")
        w.Header().Set("Content-Type", "application/jpeg")
        fmt.Fprintln(w, "iamjpg")
    }))
    defer jpg200.Close()

    jpg404 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        http.NotFound(w, r)
    }))
    defer jpg404.Close()

    shield.ComputeEngineHost = jpg404.URL
    if err := shield.GetAndRender("/media/test.jpg", wrc.c, wrc.w, wrc.r); err != nil {
        t.Fatal(err)
    }

    shield.ComputeEngineHost = jpg200.URL
    if err := shield.GetAndRender("/media/test.jpg", wrc.c, wrc.w, wrc.r); err != nil {
        t.Fatal(err)
    }
}
