package auth

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"

	"authapi/internal/user"
)

// Handler contains http handlers for authentication
type Handler struct {
	authSvc *SFAService
}

// Auth Http Handler constructor
func NewHandler(authSvc *SFAService) *Handler {
	return &Handler{authSvc: authSvc}
}

// Hardening the server by setting security headers from owasp
func setSecurityHeaders(c echo.Context) {
	c.Response().Header().Set("Content-Type", "application/json;charset=utf-8")
	c.Response().Header().Set("X-Content-Type-Options", "nosniff")
	c.Response().Header().Set("X-Frame-Options", "DENY")
	c.Response().Header().Set("X-XSS-Protection", "0")
	c.Response().Header().Set("Cache-Control", "no-store")
	c.Response().Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; sandbox")
	c.Response().Header().Set("Server", "")
}

// Gets the status code for the error
func getStatusCode(err error) int {
	switch err {
	case user.ErrUserAlreadyExists:
		return http.StatusConflict
	case user.ErrUserNotFound:
		return http.StatusNotFound
	case user.ErrInvalidPassword:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

// Handles the login request sending back a token
func (h *Handler) LoginHandler(c echo.Context) error {
	req := &user.UserReq{}
	if err := req.DecodeValidate(c.Request().Body); err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return jsoniter.NewEncoder(c.Response().Writer).Encode(map[string]string{"message": err.Error()})
	}

	usr, token, err := h.authSvc.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		code := getStatusCode(err)
		c.Response().WriteHeader(code)
		return jsoniter.NewEncoder(c.Response().Writer).Encode(map[string]string{"message": err.Error()})
	}

	setSecurityHeaders(c)
	c.Response().Header().Set("Authorization", token)
	c.Response().Header().Set("X-Uid", usr.UserID.String())
	c.Response().Header().Set("X-Email", usr.Email)

	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})

	c.Response().WriteHeader(http.StatusOK)
	return jsoniter.NewEncoder(c.Response().Writer).Encode(map[string]string{"message": "User logged in successfully"})

}

// Handles the register request sending back a token
func (h *Handler) RegisterHandler(c echo.Context) error {
	req := &user.UserReq{}
	if err := req.DecodeValidate(c.Request().Body); err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return jsoniter.NewEncoder(c.Response().Writer).Encode(map[string]string{"message": err.Error()})
	}

	usr, token, err := h.authSvc.Register(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		code := getStatusCode(err)
		c.Response().WriteHeader(code)
		return jsoniter.NewEncoder(c.Response().Writer).Encode(map[string]string{"message": err.Error()})
	}

	setSecurityHeaders(c)
	c.Response().Header().Set("Authorization", token)
	c.Response().Header().Set("X-Uid", usr.UserID.String())
	c.Response().Header().Set("X-Email", usr.Email)

	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})

	c.Response().WriteHeader(http.StatusCreated)
	return jsoniter.NewEncoder(c.Response().Writer).Encode(map[string]string{"message": "user registered successfully"})
}
