package router

import (
	"blog/config"
	"blog/internal/app"
	"blog/internal/handler"
	"blog/internal/service"
	"blog/pkg/code"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCurrentUserProfileRouteRejectsMissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	originalConfig := config.Cfg
	config.Cfg = &config.Config{CORS: config.CORSConfig{
		AllowOrigins: []string{"http://localhost"},
		AllowMethods: []string{http.MethodGet},
		AllowHeaders: []string{"Authorization"},
	}}
	t.Cleanup(func() { config.Cfg = originalConfig })
	container := &app.Container{Handler: handler.New(service.New(nil, nil, nil))}
	router := Register(container)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/users/me", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	expectedCode := fmt.Sprintf(`"code":%d`, code.Unauthorized.BizCode)
	if !strings.Contains(recorder.Body.String(), expectedCode) {
		t.Fatalf("expected unauthorized business code: %s", recorder.Body.String())
	}
}
