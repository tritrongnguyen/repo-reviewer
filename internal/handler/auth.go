package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/tritrongnguyen/repo-reviewer.git/internal/domain"
	"github.com/tritrongnguyen/repo-reviewer.git/internal/service"
	"github.com/tritrongnguyen/repo-reviewer.git/pkg/helper"
	"github.com/tritrongnguyen/repo-reviewer.git/pkg/validator"
)

type AuthHandler struct {
	auth service.AuthService
}

func NewAuthHandler(a service.AuthService) *AuthHandler {
	return &AuthHandler{
		auth: a,
	}
}

type authRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	json.NewDecoder(r.Body).Decode(&req)

	if err := validator.Validate.Struct(req); err != nil {
		helper.RespondWithJson(w, http.StatusCreated,
			helper.APIResponse[domain.User]{
				Code:    http.StatusBadRequest,
				Message: "Invalid email or password",
			},
		)
		return
	}

	err := h.auth.SignUp(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			helper.RespondWithJson(w, http.StatusConflict,
				helper.APIResponse[domain.User]{
					Code:    http.StatusConflict,
					Message: "Email already exists",
				},
			)

			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	helper.RespondWithJson(w, http.StatusCreated,
		helper.APIResponse[domain.User]{
			Code:    http.StatusCreated,
			Message: "Created",
		},
	)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	json.NewDecoder(r.Body).Decode(&req)

	if err := validator.Validate.Struct(req); err != nil {
		helper.RespondWithJson(w, http.StatusBadRequest,
			helper.APIResponse[domain.User]{
				Code:    http.StatusBadRequest,
				Message: "Invalid email or password",
			},
		)
		return
	}

	sessionID, err := h.auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		helper.RespondWithJson(w, http.StatusUnauthorized,
			helper.APIResponse[domain.User]{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			},
		)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		Secure:   false, // true if https
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	helper.RespondWithJson(w, http.StatusOK,
		helper.APIResponse[domain.User]{
			Code:    http.StatusOK,
			Message: "Login Successfully",
		},
	)
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {

}
