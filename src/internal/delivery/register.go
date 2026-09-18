package delivery

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/visionJAMx/whoknows/src/internal/repository"
)

// RegisterPage viser registreringsformularen.
func RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register", gin.H{
		"error": "",
	})
}

// ValidationErrorDetail følger API-kontraktens format for et enkelt valideringsproblem.
type ValidationErrorDetail struct {
	Loc  []string `json:"loc"`
	Msg  string   `json:"msg"`
	Type string   `json:"type"`
}

// ValidationErrorResponse pakker en eller flere valideringsfejl til et 422-svar efter kontrakten.
type ValidationErrorResponse struct {
	Detail []ValidationErrorDetail `json:"detail"`
}

// RegisterHandler forbinder HTTP-laget med databasen via service-laget.
type RegisterHandler struct {
	db *sql.DB
}

// NewRegisterHandler opretter en handler med den delte databaseforbindelse.
func NewRegisterHandler(db *sql.DB) *RegisterHandler {
	return &RegisterHandler{db: db}
}

// APIRegister validerer form-data, opretter brugeren og returnerer JSON efter kontrakten i issue #48.
// Password hashes altid sikkert, og hverken password eller password-hash indgår nogensinde i svaret.
func (handler *RegisterHandler) APIRegister(c *gin.Context) {
	username := strings.TrimSpace(c.PostForm("username"))
	email := strings.TrimSpace(c.PostForm("email"))
	password := c.PostForm("password")
	password2, password2Provided := c.GetPostForm("password2")

	var details []ValidationErrorDetail

	if username == "" {
		details = append(details, ValidationErrorDetail{
			Loc:  []string{"body", "username"},
			Msg:  "Username is required",
			Type: "value_error.missing",
		})
	} else if len(username) < 3 || len(username) > 50 {
		details = append(details, ValidationErrorDetail{
			Loc:  []string{"body", "username"},
			Msg:  "Username must be between 3 and 50 characters",
			Type: "value_error.any_str.length",
		})
	}

	if email == "" {
		details = append(details, ValidationErrorDetail{
			Loc:  []string{"body", "email"},
			Msg:  "Email is required",
			Type: "value_error.missing",
		})
	} else if !strings.Contains(email, "@") {
		details = append(details, ValidationErrorDetail{
			Loc:  []string{"body", "email"},
			Msg:  "Email must be a valid email address",
			Type: "value_error.email",
		})
	}

	if password == "" {
		details = append(details, ValidationErrorDetail{
			Loc:  []string{"body", "password"},
			Msg:  "Password is required",
			Type: "value_error.missing",
		})
	} else if len(password) < 10 {
		details = append(details, ValidationErrorDetail{
			Loc:  []string{"body", "password"},
			Msg:  "Password must be at least 10 characters",
			Type: "value_error.any_str.length",
		})
	}

	// password2 er valgfrit jf. kontrakten; når den sendes, skal den matche password.
	if password2Provided && password2 != password {
		details = append(details, ValidationErrorDetail{
			Loc:  []string{"body", "password2"},
			Msg:  "Passwords do not match",
			Type: "value_error.mismatch",
		})
	}

	if len(details) > 0 {
		c.JSON(http.StatusUnprocessableEntity, ValidationErrorResponse{Detail: details})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusInternalServerError, AuthResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Could not create account",
		})
		return
	}

	_, err = repository.CreateUser(
		c.Request.Context(),
		handler.db,
		username,
		email,
		string(passwordHash),
	)
	if err != nil {
		if repository.IsDuplicateUserError(err) {
			c.JSON(http.StatusUnprocessableEntity, ValidationErrorResponse{
				Detail: []ValidationErrorDetail{{
					Loc:  []string{"body", "username"},
					Msg:  "Username or email is already taken",
					Type: "value_error.duplicate",
				}},
			})
			return
		}

		_ = c.Error(err)
		c.JSON(http.StatusInternalServerError, AuthResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Could not create account",
		})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		StatusCode: http.StatusOK,
		Message:    "You were successfully registered",
	})
}
