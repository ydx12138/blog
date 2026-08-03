package wechat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const defaultSessionURL = "https://api.weixin.qq.com/sns/jscode2session"

var (
	ErrInvalidCode   = errors.New("wechat login code is required")
	ErrWechatLogin   = errors.New("wechat code exchange failed")
	ErrMissingOpenID = errors.New("wechat response missing openid")
)

type Session struct {
	OpenID string
}

type Exchanger interface {
	ExchangeCode(context.Context, string) (Session, error)
}

type Client struct {
	appID      string
	appSecret  string
	baseURL    string
	httpClient *http.Client
}

func NewClient(appID, appSecret, baseURL string) *Client {
	if baseURL == "" {
		baseURL = defaultSessionURL
	}
	return &Client{
		appID:      appID,
		appSecret:  appSecret,
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
	}
}

func (c *Client) ExchangeCode(ctx context.Context, code string) (Session, error) {
	if strings.TrimSpace(code) == "" {
		return Session{}, ErrInvalidCode
	}

	endpoint, err := url.Parse(c.baseURL)
	if err != nil {
		return Session{}, fmt.Errorf("parse wechat session endpoint: %w", err)
	}
	query := endpoint.Query()
	query.Set("appid", c.appID)
	query.Set("secret", c.appSecret)
	query.Set("js_code", code)
	query.Set("grant_type", "authorization_code")
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return Session{}, fmt.Errorf("build wechat request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Session{}, fmt.Errorf("exchange wechat code: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return Session{}, fmt.Errorf("%w: http status %d", ErrWechatLogin, resp.StatusCode)
	}

	var result struct {
		OpenID  string `json:"openid"`
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Session{}, fmt.Errorf("decode wechat response: %w", err)
	}
	if result.ErrCode != 0 {
		return Session{}, fmt.Errorf("%w: %d %s", ErrWechatLogin, result.ErrCode, result.ErrMsg)
	}
	if result.OpenID == "" {
		return Session{}, ErrMissingOpenID
	}
	return Session{OpenID: result.OpenID}, nil
}
