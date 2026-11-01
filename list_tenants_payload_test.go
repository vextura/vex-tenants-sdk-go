package tenants

// Hand-written (not generated): pins the @httpPayload decode of
// GET /tenants, whose body is a bare JSON array, and JSON null for an
// empty registry.

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func serve(t *testing.T, body string) *Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return New(srv.URL, "")
}

func TestListTenants_BareArray(t *testing.T) {
	out, err := serve(t, `[{"id":"yc-demo","name":"YC","regions":[{"region":"kz-north-1","status":"active"}],"created_at":"2026-11-01T00:00:00Z"}]`).ListTenants(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Tenants) != 1 || out.Tenants[0].Id != "yc-demo" || out.Tenants[0].Regions[0].Region != "kz-north-1" {
		t.Fatalf("unexpected decode: %+v", out.Tenants)
	}
}

func TestListTenants_NullIsEmpty(t *testing.T) {
	out, err := serve(t, `null`).ListTenants(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Tenants) != 0 {
		t.Fatalf("null must decode as empty, got %d", len(out.Tenants))
	}
	for range out.Tenants { // ranging a nil slice is safe: null behaves as []
	}
}
