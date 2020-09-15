package utils

import (
	"api/domains/auth"
	"api/domains/entry"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"io"
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

// AES-256 Encrypt a string using the ENC_KEY env variable
func Encrypt(str string) (string, error) {
	block, err := aes.NewCipher([]byte(os.Getenv("ENC_KEY")))
	if err != nil {
		return "", err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(str))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], []byte(str))

	//returns to base64 encoded string
	encmess := base64.URLEncoding.EncodeToString(cipherText)
	return encmess, nil
}

// AES-256 Decrypt a string using the ENC_KEY env variable
func Decrypt(str string) (string, error) {
	cipherText, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(os.Getenv("ENC_KEY")))
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return "", err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

// Adjust Day to only have following format "2020-09-15"
func AdjustDay(r []entry.Entry) []entry.Entry {
	for i, _ := range r {
		r[i].Day = r[i].Day[:10]
	}

	return r
}
