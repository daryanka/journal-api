package services

import (
	"api/clients"
	"api/domains/auth"
	"api/utils"
	"api/utils/xerror"
	"database/sql"
	"fmt"
	"github.com/daryanka/myorm"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceI interface {
	Register(req *auth.CreateUser) (*auth.TokenWithClaims, *xerror.XerrorT)
	Login(req *auth.LoginRequest) (*auth.TokenWithClaims, *xerror.XerrorT)
}

type authService struct{}

var AuthService AuthServiceI = &authService{}

func (i *authService) Register(req *auth.CreateUser) (*auth.TokenWithClaims, *xerror.XerrorT) {
	// Check that email is not already in use
	var emailID int64
	err := clients.ClientOrm.Table("users").
		Select("id").
		Where("email", "=", req.Email).
		First(&emailID)

	if err != sql.ErrNoRows {
		return nil, xerror.NewBadRequest("email already in use", "EMAIL_USED")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 11)

	if err != nil {
		return nil, xerror.NewInternalError("error creating user")
	}

	id, err := clients.ClientOrm.Table("users").Insert(myorm.H{
		"email":    req.Email,
		"password": hashedPassword,
		"2fa_url":  "test",
	})

	if err != nil {
		return nil, xerror.NewInternalError("unable to create user")
	}

	token, err := utils.CreateAuthToken(id)

	if err != nil {
		fmt.Print(err.Error())
		return nil, xerror.NewInternalError("error unable to log user in", "REDIRECT_LOGIN")
	}

	return token, nil
}

func (i *authService) Login(req *auth.LoginRequest) (*auth.TokenWithClaims, *xerror.XerrorT) {
	// Find user with email
	var User struct {
		ID       int64  `db:"id"`
		Email    string `db:"email"`
		Password string `db:"password"`
	}

	err := clients.ClientOrm.Table("users").
		Select("id", "email", "password").
		Where("email", "=", req.Email).
		First(&User)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, xerror.NewBadRequest("invalid credentials")
		}
		return nil, xerror.NewInternalError("server error")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(req.Password)); err != nil {
		return nil, xerror.NewBadRequest("invalid credentials")
	}

	// Generate Token
	token, err := utils.CreateAuthToken(User.ID)

	if err != nil {
		return nil, xerror.NewInternalError("server error")
	}

	return token, nil
}

func GetAuthUser(c *gin.Context) auth.User {
	u, _ := c.Get("user")

	user, ok := u.(auth.User)
	if !ok {
		return auth.User{}
	}
	return user
}
