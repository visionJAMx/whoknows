package delivery

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/service"
)

// LoginRequest beskriver formularfelterne i OpenAPI.
// Selve login-handleren læser fortsat felterne med PostForm.
type LoginRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// AuthResponse følger API-kontrakten for login- og registreringssvar.
type AuthResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// LoginHandler forbinder HTTP-laget med databasen via service-laget.
type LoginHandler struct {
	db *sql.DB
}

// NewLoginHandler opretter en handler med den delte databaseforbindelse.
func NewLoginHandler(db *sql.DB) *LoginHandler {
	return &LoginHandler{db: db}
}

// APILogin validerer loginoplysninger og gemmer bruger-ID i sessionen.
// @Summary Log brugeren ind
// @Description Autentificerer brugeren og opretter en sessionscookie.
// @Tags Authentication
// @Produce json
// @Param credentials formData LoginRequest true "Loginoplysninger"
// @Success 200 {object} AuthResponse
// @Failure 401 {object} AuthResponse
// @Failure 422 {object} ValidationErrorResponse
// @Failure 500 {object} AuthResponse
// @Router /api/login [post]
func (handler *LoginHandler) APILogin(context *gin.Context) {
	username := strings.TrimSpace(context.PostForm("username"))
	password := context.PostForm("password")

	// Brug samme valideringsformat som registrering og OpenAPI-kontrakten.
	var details []ValidationErrorDetail

	if username == "" {
		details = append(details, ValidationErrorDetail{
			Loc:  []string{"body", "username"},
			Msg:  "Username is required",
			Type: "value_error.missing",
		})
	}

	if password == "" {
		details = append(details, ValidationErrorDetail{
			Loc:  []string{"body", "password"},
			Msg:  "Password is required",
			Type: "value_error.missing",
		})
	}

	// Returnér alle manglende felter, før vi forsøger at logge ind.
	if len(details) > 0 {
		context.JSON(
			http.StatusUnprocessableEntity,
			ValidationErrorResponse{Detail: details},
		)
		return
	}

	user, err := service.Authenticate(
		context.Request.Context(),
		handler.db,
		username,
		password,
	)
	if err != nil {
		// Samme svar bruges ved ukendt bruger og forkert password for ikke at afsløre brugernavne.
		if errors.Is(err, service.ErrInvalidCredentials) {
			context.JSON(http.StatusUnauthorized, AuthResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "Invalid username or password",
			})
			return
		}

		_ = context.Error(err)
		context.JSON(http.StatusInternalServerError, AuthResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Could not complete login",
		})
		return
	}

	// Kun bruger-ID gemmes i den signerede session-cookie; aldrig password eller password-hash.
	session := sessions.Default(context)
	session.Set("user_id", user.ID)

	if err := session.Save(); err != nil {
		_ = context.Error(err)
		context.JSON(http.StatusInternalServerError, AuthResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Could not save login session",
		})
		return
	}

	context.JSON(http.StatusOK, AuthResponse{
		StatusCode: http.StatusOK,
		Message:    "You were logged in",
	})
}
