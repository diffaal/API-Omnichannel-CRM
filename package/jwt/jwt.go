package jwt

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type AccessToken struct {
	AccessToken  string
	ExpiredToken int64
}

type AccessTokenNodes struct {
	AccessToken string
}

func CreateToken(user_id string, username string, role int) (*AccessToken, error) {
	at := &AccessToken{}
	at.ExpiredToken = time.Now().Add(time.Minute * 480).Unix()
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = user_id
	atClaims["username"] = username
	atClaims["role"] = role
	atClaims["exp"] = at.ExpiredToken

	atTemp := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	at.AccessToken, err = atTemp.SignedString([]byte(viper.GetString("Jwt_secret")))
	if err != nil {
		return nil, err
	}

	return at, nil
}

func ExtractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	strArr := strings.Split(token, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Wrong signature method")
		}
		return []byte(viper.GetString("Jwt_secret")), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}

	return nil
}

func GetDataFromToken(atn *AccessTokenNodes) map[string]interface{} {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(atn.AccessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("Jwt_secret")), nil
	})

	if err != nil {
		responseData := map[string]interface{}{
			"error":    err.Error(),
			"user_id":  "",
			"username": "",
			"role":     0,
		}
		return responseData
	}

	if token.Valid {
		responseData := map[string]interface{}{
			"error":           nil,
			"user_id":         claims["user_id"],
			"username":        claims["username"].(string),
			"role":            claims["role"].(float64),
			"channel_account": claims["channel_account"],
		}
		return responseData
	}

	responseData := map[string]interface{}{}

	return responseData
}
