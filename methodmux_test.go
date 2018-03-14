package methodmux_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/acoshift/methodmux"
)

type inj struct{}

func _handler(w http.ResponseWriter, r *http.Request) {
	*r.Context().Value(inj{}).(*int) = 1
}

var (
	handler = http.HandlerFunc(_handler)
	ctxBg   = context.Background()
)

func testMethod(t *testing.T, method string, h http.Handler, pass bool) {
	var p int
	w := httptest.NewRecorder()
	h.ServeHTTP(w, (&http.Request{Method: method}).WithContext(context.WithValue(ctxBg, inj{}, &p)))
	if pass && p != 1 {
		t.Errorf("method %s not passed", method)
	} else if !pass && w.Code != http.StatusNotFound {
		t.Errorf("method %s must returns not found error", method)
	}
}

func TestMethodMux(t *testing.T) {
	testCases := []struct {
		method  string
		handler http.Handler
		pass    bool
	}{
		{http.MethodGet, methodmux.Get(handler), true},
		{http.MethodPost, methodmux.Post(handler), true},
		{http.MethodPatch, methodmux.Patch(handler), true},
		{http.MethodPut, methodmux.Put(handler), true},
		{http.MethodDelete, methodmux.Delete(handler), true},
		{http.MethodHead, methodmux.Head(handler), true},
		{http.MethodOptions, methodmux.Options(handler), true},

		{http.MethodGet, methodmux.GetPost(handler, nil), true},
		{http.MethodPost, methodmux.GetPost(nil, handler), true},

		{http.MethodPatch, methodmux.Post(handler), false},
		{http.MethodHead, methodmux.Get(handler), true},
		{http.MethodHead, methodmux.Post(handler), false},

		{http.MethodGet, methodmux.Mux{"": handler}, true},
		{http.MethodPost, methodmux.Mux{"": handler}, true},
		{http.MethodPatch, methodmux.Mux{"": handler}, true},
		{http.MethodPut, methodmux.Mux{"": handler}, true},
		{http.MethodDelete, methodmux.Mux{"": handler}, true},
		{http.MethodHead, methodmux.Mux{"": handler}, true},
		{http.MethodOptions, methodmux.Mux{"": handler}, true},
		{"PURGE", methodmux.Mux{"": handler}, true},
		{"", methodmux.Mux{"": handler}, true},
	}
	for _, c := range testCases {
		testMethod(t, c.method, c.handler, c.pass)
	}
}
