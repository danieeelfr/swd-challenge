package login

import (
	"context"
	"time"

	"github.com/danieeelfr/swd-challenge/internal/config"
	"github.com/danieeelfr/swd-challenge/internal/models"

	srv "github.com/danieeelfr/swd-challenge/internal/service/user"
	"github.com/dgrijalva/jwt-go"
)

// JWTCustomClaims holds the token claims values
type JWTCustomClaims struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
	Role string `json:"role"`
	jwt.StandardClaims
}

// Login holds the login implementation itself
type Login struct {
	cfg         *config.MySQLRepositoryConfig
	userService *srv.User
}

// NewLoginService returns the login implementation
func NewLoginService(cfg *config.Config) (*Login, error) {
	us, err := srv.NewUserService(cfg)
	if err != nil {
		return nil, err
	}
	return &Login{
		cfg:         cfg.MySQLRepositoryConfig,
		userService: us,
	}, nil
}

// GetToken ...
func (l *Login) GetToken(user *models.User) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = user.User
	claims["name"] = user.Name
	claims["id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return t, nil

}

// GetUser get the user to validate on the login process
func (l *Login) GetUser(ctx context.Context, user, password string) (*models.User, error) {
	filter := models.UserFilter{User: user, Password: password}

	u, err := l.userService.GetUser(ctx, &filter)
	if err != nil {
		return nil, err
	}

	return u, nil
}
