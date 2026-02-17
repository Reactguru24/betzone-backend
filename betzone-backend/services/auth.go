package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/betzone/backend/config"
	"github.com/betzone/backend/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	jwtSecret string
	tokenKey  string // Betkraft token key for callback validation
	db        *gorm.DB
}

type CustomClaims struct {
	UserID string `json:"user_id"`
	Phone  string `json:"phone"`
	jwt.RegisteredClaims
}

// NewAuthService creates a new authentication service with database
func NewAuthService(jwtSecret string, db *gorm.DB) *AuthService {
	return &AuthService{
		jwtSecret: jwtSecret,
		tokenKey:  "", // Will be set later
		db:        db,
	}
}

// NewAuthServiceWithConfig creates auth service with full config
func NewAuthServiceWithConfig(cfg *config.Config, db *gorm.DB) *AuthService {
	return &AuthService{
		jwtSecret: cfg.JWTSecret,
		tokenKey:  cfg.BetkraftAppKey, // Use AppKey as tokenKey for callback validation
		db:        db,
	}
}

// GetTokenKey returns the Betkraft token key for callback validation
func (as *AuthService) GetTokenKey() string {
	if as.tokenKey == "" {
		// Fallback: you could set this via config or environment
		return ""
	}
	return as.tokenKey
}

// HashPassword hashes a password using bcrypt
func (as *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(hash), nil
}

// CheckPassword verifies a password against its hash
func (as *AuthService) CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Signup creates a new user account in the database
func (as *AuthService) Signup(req *models.SignupRequest) (*models.User, string, error) {
	// Check if user already exists
	var existingUser models.User
	if result := as.db.Where("phone = ?", req.Phone).First(&existingUser); result.Error == nil {
		return nil, "", errors.New("user with this phone number already exists")
	} else if result.Error != gorm.ErrRecordNotFound {
		return nil, "", fmt.Errorf("database error: %v", result.Error)
	}

	// Hash password
	hashedPassword, err := as.HashPassword(req.Password)
	if err != nil {
		return nil, "", err
	}

	// Create user
	user := models.User{
		ID:        generateUserID(),
		Phone:     req.Phone,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Balance:   0,
		Status:    "active",
	}

	// Save user to database
	if result := as.db.Create(&user); result.Error != nil {
		return nil, "", fmt.Errorf("failed to create user: %v", result.Error)
	}

	// Generate JWT token
	token, err := as.GenerateToken(&user)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

// Signin authenticates a user and returns a JWT token
func (as *AuthService) Signin(req *models.SigninRequest) (*models.User, string, error) {
	// Find user by phone
	user := &models.User{}
	if result := as.db.Where("phone = ?", req.Phone).First(user); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, "", errors.New("invalid phone number or password")
		}
		return nil, "", fmt.Errorf("database error: %v", result.Error)
	}

	// Check password
	if err := as.CheckPassword(user.Password, req.Password); err != nil {
		return nil, "", errors.New("invalid phone number or password")
	}

	// Generate JWT token
	token, err := as.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// GenerateToken creates a JWT token for a user
func (as *AuthService) GenerateToken(user *models.User) (string, error) {
	claims := CustomClaims{
		UserID: user.ID,
		Phone:  user.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "betzone-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(as.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

// VerifyToken validates and parses a JWT token
func (as *AuthService) VerifyToken(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(as.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	return claims, nil
}

// GetUserByID retrieves a user by their ID from the database
func (as *AuthService) GetUserByID(userID string) (*models.User, error) {
	user := &models.User{}
	if result := as.db.Where("id = ?", userID).First(user); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %v", result.Error)
	}
	return user, nil
}

// GetUserByPhone retrieves a user by their phone number from the database
func (as *AuthService) GetUserByPhone(phone string) (*models.User, error) {
	user := &models.User{}
	if result := as.db.Where("phone = ?", phone).First(user); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %v", result.Error)
	}
	return user, nil
}

// UpdateUserBalance updates a user's balance
func (as *AuthService) UpdateUserBalance(userID string, amount float64) error {
	if result := as.db.Model(&models.User{}).Where("id = ?", userID).Update("balance", gorm.Expr("balance + ?", amount)); result.Error != nil {
		return fmt.Errorf("failed to update balance: %v", result.Error)
	}
	return nil
}

// generateUserID generates a UUID-like user ID
func generateUserID() string {
	return fmt.Sprintf("user_%d", time.Now().UnixNano())
}
