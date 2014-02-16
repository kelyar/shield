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
  w *httptest.ResponseRecorder
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
    if wrc.w.Code != 404 {
        t.Fatal("expected 404, got %v", wrc.w.Code)
    }

    wrc.w = httptest.NewRecorder()
    shield.ComputeEngineHost = jpg200.URL
    if err := shield.GetAndRender("/media/test.jpg", wrc.c, wrc.w, wrc.r); err != nil {
        t.Fatal(err)
    }
    if wrc.w.Code != 200 {
        t.Fatal("expected 200, got %v", wrc.w.Code)
    }
    if wrc.w.Body.String() != "iamjpg\n" {
        t.Fatal("expected iamjpg, got ", wrc.w.Body.String())
    }
}
