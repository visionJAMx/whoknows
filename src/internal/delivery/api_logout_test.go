package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/internal/repository"
	"github.com/visionJAMx/whoknows/src/service"
)

func setupLogoutTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	databasePath := filepath.Join(t.TempDir(), "test.db")

	db, err := repository.Open(databasePath)
	if err != nil {
		t.Fatalf("could not open test database: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	ctx := context.Background()
	if err := repository.Initialize(ctx, db); err != nil {
		t.Fatalf("could not initialize test database: %v", err)
	}

	passwordHash, err := service.HashPassword("correct-password")
	if err != nil {
		t.Fatalf("could not hash password: %v", err)
	}

	_, err = repository.CreateUser(
		ctx,
		db,
		"osman",
		"osman@example.com",
		passwordHash,
	)
	if err != nil {
		t.Fatalf("could not create test user: %v", err)
	}

	router := gin.New()

	store := cookie.NewStore(
		[]byte("test-session-secret-with-at-least-32-characters"),
	)
	router.Use(sessions.Sessions("whoknows_session", store))

	loginHandler := NewLoginHandler(db)
	router.POST("/api/login", loginHandler.APILogin)

	router.POST("/api/logout", LogoutHandler)

	return router
}

func loginAndGetSessionCookie(router *gin.Engine, t *testing.T) string {
	t.Helper()

	form := url.Values{}
	form.Set("username", "osman")
	form.Set("password", "correct-password")

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/login",
		strings.NewReader(form.Encode()),
	)
	request.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	cookie := response.Header().Get("Set-Cookie")
	if cookie == "" {
		t.Fatalf("expected login to set a session cookie")
	}

	return cookie
}

func performLogoutRequest(router *gin.Engine, sessionCookie string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
	if sessionCookie != "" {
		request.Header.Set("Cookie", sessionCookie)
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	return response
}

func TestAPILogoutClearsActiveSession(t *testing.T) {
	router := setupLogoutTestRouter(t)
	sessionCookie := loginAndGetSessionCookie(router, t)

	response := performLogoutRequest(router, sessionCookie)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}

	var body AuthResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if body.Message != "You were logged out" {
		t.Errorf("unexpected message: %s", body.Message)
	}

	if response.Header().Get("Set-Cookie") == "" {
		t.Error("expected logout to update the session cookie")
	}
}

func TestAPILogoutWithoutActiveSession(t *testing.T) {
	router := setupLogoutTestRouter(t)

	response := performLogoutRequest(router, "")

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}
}

func TestLogoutRouteRejectsGet(t *testing.T) {
	router := setupLogoutTestRouter(t)

	request := httptest.NewRequest(http.MethodGet, "/api/logout", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code == http.StatusOK {
		t.Error("expected GET /api/logout to NOT log the user out (CSRF risk)")
	}
}
