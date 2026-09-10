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

// setupLoginTestRouter bygger en isoleret router med midlertidig database og testbruger.
func setupLoginTestRouter(t *testing.T) *gin.Engine {
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

	return router
}

// performLoginRequest sender samme form-type, som login-siden bruger i browseren.
func performLoginRequest(
	router *gin.Engine,
	username string,
	password string,
) *httptest.ResponseRecorder {
	form := url.Values{}
	form.Set("username", username)
	form.Set("password", password)

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

	return response
}

// TestAPILoginSuccess kontrollerer både HTTP-svar og at en session-cookie bliver oprettet.
func TestAPILoginSuccess(t *testing.T) {
	router := setupLoginTestRouter(t)

	response := performLoginRequest(
		router,
		"osman",
		"correct-password",
	)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}

	if response.Header().Get("Set-Cookie") == "" {
		t.Error("expected successful login to create a session cookie")
	}
}

// TestAPILoginRejectsWrongPassword sikrer, at ugyldige oplysninger afvises uden brugerlæk.
func TestAPILoginRejectsWrongPassword(t *testing.T) {
	router := setupLoginTestRouter(t)

	response := performLoginRequest(
		router,
		"osman",
		"wrong-password",
	)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			response.Code,
		)
	}

	var body AuthResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if body.Message != "Invalid username or password" {
		t.Errorf("unexpected error message: %s", body.Message)
	}
}

// TestAPILoginRejectsMissingFields sikrer, at obligatoriske formfelter valideres.
func TestAPILoginRejectsMissingFields(t *testing.T) {
	router := setupLoginTestRouter(t)

	response := performLoginRequest(router, "", "")

	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnprocessableEntity,
			response.Code,
		)
	}
}
