package main

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
	"strings"
	"errors"
)

type User struct {
	ID 		 int 	`json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	UUID 	 string `json:"uuid"`
}

type Credentials struct {
	FullName string `json:"full_name"`
	Email    string `json:"email",	db:"email"`
	Password string `json:"password", db:"password"`
}

const SigningKey = "dfajlkfjqwopie"
const UserJWTExpirationTime  = time.Hour * 24 * 7

func SignUp(c *gin.Context) {
	user := &User{}
	if bindErr := c.BindJSON(&user); bindErr != nil {
		log.Println(bindErr)
		throwStatusBadRequest(bindErr.Error(), c)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.UUID = generateToken(16)
	user.Password = string(hashedPassword)

	if err = db.Save(user).Error; err != nil {
		log.Println(err.Error())
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	throwStatusOk("OK", c)
	return
}

func SignIn(c *gin.Context) {
	creds := &Credentials{}
	if bindErr := c.BindJSON(&creds); bindErr != nil {
		log.Println(bindErr)
		throwStatusBadRequest(bindErr.Error(), c)
		return
	}

	user := User{}

	if err := db.Debug().Raw("select * from users where email=?", creds.Email).Scan(&user).Error; err != nil {
		log.Println(err, "useriimejl")
		throwStatusUnauthorized(c)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		log.Println(err.Error(), "userpasvord")
		throwStatusUnauthorized(c)
		return
	}

	tokenString := generateToken(64)

	if err := db.Exec("UPDATE users SET login_token = ? WHERE id =?", tokenString, user.ID).Error; err != nil {
		log.Println(err)
		throwStatusInternalServerError("DB_ERR", c)
		return
	}

	tokenMap := map[string]string{"token":tokenString}

	throwStatusOk(tokenMap, c)
	return
}

func getUserFromToken(token string) (User, error){
	var user User

	query := `SELECT * FROM users WHERE login_token = ?`

	count := db.Debug().Raw(query, token).Scan(&user).RowsAffected
	if count == 0 {
		return User{}, errors.New("DB_ERR")
	}
	return user, nil
}

func getTokenFromRequest(c *gin.Context) string {
	token := c.Request.Header.Get("Authorization")
	if len(token)>0{
		token = strings.TrimPrefix(token, "Bearer ")
		log.Println("TOKEN: ", token)
		return token
	}
	return ""
}

