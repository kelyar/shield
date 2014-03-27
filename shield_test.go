package shield_test

import (
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
  c aetest.Context
}

const (
  MagicRequest = "/media/test.jpg"
  MagicResponse = "iamjpg"
)

var (
    wrc Wrc
)

func NewWrc(t *testing.T) Wrc {
    c, err := aetest.NewContext(nil)
    if err != nil {
        t.Fatal(err)
    }

    w := httptest.NewRecorder()

    req, err := http.NewRequest("GET", "http://example.com" + MagicRequest, nil)
    if err != nil {
        t.Fatal(err)
    }

    return Wrc{w, req, c}
}

func TestRespondWithHeader(t *testing.T) {
    wrc = NewWrc(t)

    if err := shield.RespondWithHeader(MagicRequest, wrc.c, wrc.w, wrc.r); err != nil {
        t.Fatal(err)
    }
    if "encoded_gs_file:d2l4c3RhdGljZGV2L21lZGlhL3Rlc3QuanBn" != wrc.w.Header().Get("X-AppEngine-BlobKey") {
        t.Fatal("expected encoded_gs_file, got", wrc.w.Header().Get("X-AppEngine-BlobKey"))
    }
    if "image/jpeg" != wrc.w.Header().Get("Content-Type") {
        t.Fatal("expected image/jpeg, got", wrc.w.Header().Get("Content-Type"))
    }
    defer wrc.c.Close()
}

func TestGetAndRender404(t *testing.T) {
    wrc = NewWrc(t)

    jpg404 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        http.NotFound(w, r)
    }))
    defer jpg404.Close()

    shield.ComputeEngineHost = jpg404.URL
    if err := shield.GetAndRender(MagicRequest, wrc.c, wrc.w, wrc.r); err != nil {
        t.Fatal(err)
    }
    if wrc.w.Code != 404 {
        t.Fatal("expected 404, got", wrc.w.Code)
    }
    defer wrc.c.Close()
}

func TestGetAndRender200(t *testing.T) {
    wrc = NewWrc(t)

    jpg200 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/jpeg")
        fmt.Fprint(w, MagicResponse)
    }))
    defer jpg200.Close()

    shield.ComputeEngineHost = jpg200.URL
    if err := shield.GetAndRender(MagicRequest, wrc.c, wrc.w, wrc.r); err != nil {
        t.Fatal(err)
    }
    if wrc.w.Code != 200 {
        t.Fatal("expected 200, got", wrc.w.Code)
    }
    if wrc.w.Body.String() != MagicResponse {
        t.Fatal("expected "+ MagicResponse + " , got ", wrc.w.Body.String())
    }
    defer wrc.c.Close()
}
