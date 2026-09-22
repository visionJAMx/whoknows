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

	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/internal/repository"
)

// setupRegisterTestRouter bygger en isoleret router med en midlertidig database.
func setupRegisterTestRouter(t *testing.T) *gin.Engine {
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

	router := gin.New()

	registerHandler := NewRegisterHandler(db)
	router.POST("/api/register", registerHandler.APIRegister)

	return router
}

// performRegisterRequest sender samme form-type, som registreringssiden bruger i browseren.
func performRegisterRequest(
	router *gin.Engine,
	form url.Values,
) *httptest.ResponseRecorder {
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/register",
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

func TestRegisterPageReturnsHTML(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.LoadHTMLGlob("../../templates/*.html")
	router.GET("/register", RegisterPage)

	request := httptest.NewRequest(
		http.MethodGet,
		"/register",
		nil,
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}
}

// TestAPIRegisterSuccess dækker acceptkriteriet: succes giver HTTP 200 med statusCode og message.
func TestAPIRegisterSuccess(t *testing.T) {
	router := setupRegisterTestRouter(t)

	form := url.Values{}
	form.Set("username", "testuser")
	form.Set("email", "test@example.com")
	form.Set("password", "secret12345")

	response := performRegisterRequest(router, form)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusOK,
			response.Code,
			response.Body.String(),
		)
	}

	var body AuthResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if body.StatusCode != http.StatusOK {
		t.Errorf("unexpected statusCode in body: %d", body.StatusCode)
	}

	if body.Message == "" {
		t.Error("expected a non-empty success message")
	}

	if strings.Contains(response.Body.String(), "secret12345") {
		t.Error("response must never contain the plaintext password")
	}
}

// TestAPIRegisterMissingFieldsReturns422 dækker: username, email og password er påkrævede,
// og valideringsfejl giver HTTP 422 med detail; hvert element indeholder loc, msg og type.
func TestAPIRegisterMissingFieldsReturns422(t *testing.T) {
	router := setupRegisterTestRouter(t)

	response := performRegisterRequest(router, url.Values{})

	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnprocessableEntity,
			response.Code,
		)
	}

	var body ValidationErrorResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(body.Detail) == 0 {
		t.Fatal("expected at least one validation error detail")
	}

	for _, detail := range body.Detail {
		if len(detail.Loc) == 0 || detail.Msg == "" || detail.Type == "" {
			t.Errorf("incomplete validation detail: %+v", detail)
		}
	}
}

// TestAPIRegisterPassword2IsOptional dækker: password2 er valgfrit som i kontrakten.
func TestAPIRegisterPassword2IsOptional(t *testing.T) {
	router := setupRegisterTestRouter(t)

	form := url.Values{}
	form.Set("username", "testuser")
	form.Set("email", "test@example.com")
	form.Set("password", "secret12345")

	response := performRegisterRequest(router, form)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected registration without password2 to succeed, got %d: %s",
			response.Code,
			response.Body.String(),
		)
	}
}

// TestAPIRegisterPassword2Mismatch dækker: når password2 sendes, skal den matche password.
func TestAPIRegisterPassword2Mismatch(t *testing.T) {
	router := setupRegisterTestRouter(t)

	form := url.Values{}
	form.Set("username", "testuser")
	form.Set("email", "test@example.com")
	form.Set("password", "secret12345")
	form.Set("password2", "does-not-match")

	response := performRegisterRequest(router, form)

	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnprocessableEntity,
			response.Code,
		)
	}
}

// TestAPIRegisterDuplicateUsername dækker: dubletter håndteres.
func TestAPIRegisterDuplicateUsername(t *testing.T) {
	router := setupRegisterTestRouter(t)

	form := url.Values{}
	form.Set("username", "testuser")
	form.Set("email", "test@example.com")
	form.Set("password", "secret12345")

	first := performRegisterRequest(router, form)
	if first.Code != http.StatusOK {
		t.Fatalf(
			"expected first registration to succeed, got %d: %s",
			first.Code,
			first.Body.String(),
		)
	}

	second := performRegisterRequest(router, form)
	if second.Code != http.StatusUnprocessableEntity {
		t.Fatalf(
			"expected duplicate registration to be rejected with %d, got %d",
			http.StatusUnprocessableEntity,
			second.Code,
		)
	}
}
