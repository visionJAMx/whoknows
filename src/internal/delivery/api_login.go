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

// APILogin validerer form-data, autentificerer brugeren og gemmer bruger-ID i sessionen.
func (handler *LoginHandler) APILogin(context *gin.Context) {
	username := strings.TrimSpace(context.PostForm("username"))
	password := context.PostForm("password")

	if username == "" || password == "" {
		context.JSON(http.StatusUnprocessableEntity, AuthResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Username and password are required",
		})
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
