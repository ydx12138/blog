package wechat

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientExchangeCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("appid") != "app-id" {
			t.Fatal("missing appid")
		}
		if r.URL.Query().Get("secret") != "app-secret" {
			t.Fatal("missing secret")
		}
		if r.URL.Query().Get("js_code") != "login-code" {
			t.Fatal("missing code")
		}
		if r.URL.Query().Get("grant_type") != "authorization_code" {
			t.Fatal("missing grant type")
		}
		_, _ = w.Write([]byte(`{"openid":"openid-1","session_key":"secret"}`))
	}))
	defer server.Close()

	got, err := NewClient("app-id", "app-secret", server.URL).ExchangeCode(context.Background(), "login-code")
	if err != nil {
		t.Fatalf("ExchangeCode returned error: %v", err)
	}
	if got.OpenID != "openid-1" {
		t.Fatalf("expected openid-1, got %q", got.OpenID)
	}
}

func TestClientExchangeCodeRejectsWechatError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"errcode":40029,"errmsg":"invalid code"}`))
	}))
	defer server.Close()

	_, err := NewClient("app-id", "app-secret", server.URL).ExchangeCode(context.Background(), "login-code")
	if err == nil {
		t.Fatal("expected WeChat error")
	}
}
