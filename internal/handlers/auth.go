package handlers

import (
	"encoding/json"
	"net/http"

	"avito-internship-2025/internal/service"
	"log/slog"
)

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type AuthHandler struct {
	AuthService *service.AuthService
	Logger      *slog.Logger
}

func NewAuthHandler(authService *service.AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
		Logger:      logger,
	}
}

func (h *AuthHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Получен запрос /api/auth")
	var req service.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("Ошибка декодирования запроса", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: "Неверный запрос"})
		return
	}

	res, err := h.AuthService.Authenticate(req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			h.Logger.Warn("Неверные учетные данные", slog.String("username", req.Username))
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Errors: "Неверный пароль"})
			return
		}
		h.Logger.Error("Ошибка аутентификации", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: "Внутренняя ошибка сервера"})
		return
	}

	h.Logger.Info("Аутентификация успешна", slog.String("username", req.Username))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
