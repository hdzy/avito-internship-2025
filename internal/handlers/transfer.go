package handlers

import (
	"encoding/json"
	"net/http"

	"avito-internship-2025/internal/middleware"
	"avito-internship-2025/internal/service"
	"log/slog"
)

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type TransferHandler struct {
	CoinService *service.CoinService
	Logger      *slog.Logger
}

func NewTransferHandler(coinService *service.CoinService, logger *slog.Logger) *TransferHandler {
	return &TransferHandler{
		CoinService: coinService,
		Logger:      logger,
	}
}

func (h *TransferHandler) SendCoin(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Получен запрос на перевод монет")
	var req SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("Ошибка декодирования запроса", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: "Неверный запрос"})
		return
	}

	// Извлекаем username из контекста
	username, ok := r.Context().Value(middleware.ContextKeyUsername).(string)
	if !ok || username == "" {
		h.Logger.Warn("Username отсутствует в контексте")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: "Не авторизован"})
		return
	}

	if err := h.CoinService.TransferCoins(username, req.ToUser, req.Amount); err != nil {
		h.Logger.Error("Ошибка при переводе монет", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: "Ошибка перевода монет"})
		return
	}

	h.Logger.Info("Перевод монет успешно выполнен")
	w.WriteHeader(http.StatusOK)
}
