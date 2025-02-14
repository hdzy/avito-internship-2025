package handlers

import (
	"encoding/json"
	"net/http"

	"avito-internship-2025/internal/middleware"
	"avito-internship-2025/internal/service"
	"log/slog"
)

type BuyHandler struct {
	CoinService *service.CoinService
	Logger      *slog.Logger
}

func NewBuyHandler(coinService *service.CoinService, logger *slog.Logger) *BuyHandler {
	return &BuyHandler{
		CoinService: coinService,
		Logger:      logger,
	}
}

func (h *BuyHandler) BuyMerch(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Получен запрос на покупку мерча")

	item := r.URL.Path[len("/api/buy/"):]
	if item == "" {
		h.Logger.Warn("Не указан параметр item")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: "Не указан товар для покупки"})
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

	if err := h.CoinService.BuyMerch(username, item); err != nil {
		h.Logger.Error("Ошибка при покупке мерча", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: "Ошибка покупки мерча"})
		return
	}

	h.Logger.Info("Покупка мерча успешно выполнена", slog.String("item", item))
	w.WriteHeader(http.StatusOK)
}
