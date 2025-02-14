package handlers

import (
	"encoding/json"
	"net/http"

	"avito-internship-2025/internal/middleware"
	"avito-internship-2025/internal/service"
	"log/slog"
)

type InfoHandler struct {
	InfoService *service.InfoService
	Logger      *slog.Logger
}

func NewInfoHandler(infoService *service.InfoService, logger *slog.Logger) *InfoHandler {
	return &InfoHandler{
		InfoService: infoService,
		Logger:      logger,
	}
}

func (h *InfoHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Получен запрос /api/info")
	username, ok := r.Context().Value(middleware.ContextKeyUsername).(string)
	if !ok || username == "" {
		h.Logger.Warn("Username отсутствует в контексте")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: "Не авторизован"})
		return
	}

	info, err := h.InfoService.GetInfo(username)
	if err != nil {
		h.Logger.Error("Ошибка получения информации", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: "Ошибка получения информации"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(info)
}
