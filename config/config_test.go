package config

import "testing"

func TestWechatConfigReadsSecretFromEnvironment(t *testing.T) {
	t.Setenv("WECHAT_APP_SECRET", "test-secret")
	cfg := WechatConfig{AppID: "wx-test"}
	cfg.ApplyEnv()
	if cfg.AppSecret != "test-secret" {
		t.Fatalf("expected environment secret, got %q", cfg.AppSecret)
	}
}
