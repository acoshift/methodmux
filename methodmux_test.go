package methodmux_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/acoshift/methodmux"
)

type injecter struct{}

func _handler(w http.ResponseWriter, r *http.Request) {
	*r.Context().Value(injecter{}).(*int) = 1
}

var (
	handler = http.HandlerFunc(_handler)
	ctxBg   = context.Background()
)

func testMethod(t *testing.T, method string, h http.Handler, pass bool) {
	var p int
	w := httptest.NewRecorder()
	h.ServeHTTP(w, (&http.Request{Method: method}).WithContext(context.WithValue(ctxBg, injecter{}, &p)))
	if pass && p != 1 {
		t.Errorf("method %s not passed", method)
	} else if !pass {
		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("method %s must returns method not allowed error", method)
		}
	}
}

func TestMethodMux(t *testing.T) {
	testMethod(t, http.MethodGet, methodmux.Get(handler), true)
	testMethod(t, http.MethodPost, methodmux.Post(handler), true)
	testMethod(t, http.MethodPatch, methodmux.Patch(handler), true)
	testMethod(t, http.MethodPut, methodmux.Put(handler), true)
	testMethod(t, http.MethodDelete, methodmux.Delete(handler), true)
	testMethod(t, http.MethodHead, methodmux.Head(handler), true)
	testMethod(t, http.MethodOptions, methodmux.Options(handler), true)

	testMethod(t, http.MethodGet, methodmux.GetPost(handler, nil), true)
	testMethod(t, http.MethodPost, methodmux.GetPost(nil, handler), true)

	testMethod(t, http.MethodPatch, methodmux.Post(handler), false)
	testMethod(t, http.MethodHead, methodmux.Get(handler), true)
	testMethod(t, http.MethodHead, methodmux.Post(handler), false)

	testMethod(t, http.MethodGet, methodmux.Mux{"": handler}, true)
	testMethod(t, http.MethodPost, methodmux.Mux{"": handler}, true)
	testMethod(t, http.MethodPatch, methodmux.Mux{"": handler}, true)
	testMethod(t, http.MethodPut, methodmux.Mux{"": handler}, true)
	testMethod(t, http.MethodDelete, methodmux.Mux{"": handler}, true)
	testMethod(t, http.MethodHead, methodmux.Mux{"": handler}, true)
	testMethod(t, http.MethodOptions, methodmux.Mux{"": handler}, true)
	testMethod(t, "PURGE", methodmux.Mux{"": handler}, true)
	testMethod(t, "", methodmux.Mux{"": handler}, true)
}
