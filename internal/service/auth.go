package service

import (
	"log/slog"
	"time"

	"avito-internship-2025/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	EmployeeRepo *repository.EmployeeRepository
	JWTSecret    []byte
	Logger       *slog.Logger
}

func NewAuthService(repo *repository.EmployeeRepository, secret string, logger *slog.Logger) *AuthService {
	return &AuthService{
		EmployeeRepo: repo,
		JWTSecret:    []byte(secret),
		Logger:       logger,
	}
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

// ErrInvalidCredentials ошибка несовпадения пароля
var ErrInvalidCredentials = jwt.NewValidationError("Неавторизован", jwt.ValidationErrorSignatureInvalid)

// Authenticate ищет сотрудника по username или создает его
// и генерирует JWT-token
func (s *AuthService) Authenticate(req AuthRequest) (*AuthResponse, error) {
	employee, err := s.EmployeeRepo.GetEmployeeByUsername(req.Username)
	if err != nil {
		return nil, err
	}

	if employee != nil {
		if err := bcrypt.CompareHashAndPassword([]byte(employee.PasswordHash), []byte(req.Password)); err != nil {
			return nil, ErrInvalidCredentials
		}
		s.Logger.Info("Пользователь аутентифицирован", slog.String("username", req.Username))
	} else {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		employee, err = s.EmployeeRepo.CreateEmployee(req.Username, string(hashedPassword))
		if err != nil {
			return nil, err
		}
		s.Logger.Info("Создан новый сотрудник", slog.String("username", req.Username), slog.Int("employee_id", employee.ID))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"employee_id": employee.ID,
		"username":    employee.Username,
		"exp":         time.Now().Add(1 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(s.JWTSecret)
	if err != nil {
		return nil, err
	}

	s.Logger.Debug("JWT токен сгенерирован", slog.String("username", req.Username))
	return &AuthResponse{Token: tokenString}, nil
}
