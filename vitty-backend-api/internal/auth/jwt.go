package auth

import (
	"errors"
	"strings"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm/clause"
)

func CreateJWTToken(username string, email string, role string, jwtKey string) (string, error) {
	claims := jwt.MapClaims{}
	claims["username"] = username
	claims["email"] = email
	claims["role"] = role

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}

func getJWTClaims(tokenString string, jwtKey string) (jwt.MapClaims, error) {
	var err error
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, err
		}
		return []byte(jwtKey), nil
	})

	if err != nil {
		return nil, errors.New("invalid JWT Token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		return claims, nil
	}
	return claims, errors.New("invalid claims")
}

func GetUserFromJWTToken(tokenString string, jwtKey string) (models.User, error) {
	authorizationString := strings.Split(tokenString, " ")

	if len(authorizationString) != 2 {
		return models.User{}, errors.New("invalid Authorization Header")
	}

	token := authorizationString[1]

	if token == "" {
		return models.User{}, errors.New("not logged in")
	}

	claims, err := getJWTClaims(token, jwtKey)

	var user models.User

	if err != nil {
		return user, err
	}
	if database.DB.Where("username = ?", claims["username"]).First(&user).RowsAffected == 0 {
		return user, errors.New("request made by invalid user")
	}

	database.DB.Where("username = ?", claims["username"]).Preload(clause.Associations).First(&user)
	return user, nil
}
