package utils

import (
	"api/domains/auth"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"os"
	"strings"
	"time"
)

func ValidateStruct(obj interface{}) map[string]string {
	v := validator.New()
	err := v.Struct(obj)

	if err == nil {
		return nil
	}

	errMap := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		field := e.Namespace()
		fieldSplit := strings.Split(field, ".")
		fieldSplit = fieldSplit[1:]
		field = strings.Join(fieldSplit, ".")

		switch e.ActualTag() {
		case "required":
			errMap[field] = fmt.Sprintf("%v is a required field", field)
		default:
			continue
		}
	}

	return errMap
}

func CreateAuthToken(userID int64) (*auth.TokenWithClaims, error) {
	type JWTClaims struct {
		ID int64 `json:"id"`
		jwt.StandardClaims
	}
	customClaims := JWTClaims{
		ID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 3600,
			IssuedAt:  time.Now().Unix(),
			Issuer:    "journal_api",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return &auth.TokenWithClaims{
		Token:   tokenString,
		Expires: customClaims.StandardClaims.ExpiresAt,
	}, err
}

func ParseDate(d string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", d)
	return t, err
}

func ErrorLogger(e error) {
	if e == nil {
		return
	}
	fmt.Println(e.Error())
	f, err := os.OpenFile("error_log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	currTime := time.Now()
	str := fmt.Sprintf("Timestamp: %v \nError: %v \n", currTime, e.Error())
	_, _ = f.WriteString(str)
	f.Close()
}
