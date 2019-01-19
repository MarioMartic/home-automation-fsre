package main

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
	"errors"
	"strconv"
)

type User struct {
	ID 		 	int 	`json:"id"`
	FullName 	string `json:"full_name"`
	Email    	string `json:"email"`
	Password 	string `json:"password"`
	UUID 	 	string `json:"uuid"`
	LoginToken	string `json:"login_token"`
}

type UserV2 struct {
	ID 		 int 	`json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	UUID 	 string `json:"uuid"`
	Token	 string `json:"token"`
	MicrocontrollersCount int `json:"microcontrollers_count"`
}

type Credentials struct {
	FullName string `json:"full_name"`
	Email    string `json:"email",	db:"email"`
	Password string `json:"password", db:"password"`
}

type UUID struct {
	Text string `json:"uuid"`
}

type UserID struct {
	Text string `json:"id"`
}


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

	user := UserV2{}

	if err := db.Debug().Raw("select u.*, count(um.user_id) as microcontrollers_count from users u left join users_microcontrollers um on um.user_id = u.id where u.email=? GROUP BY u.id", creds.Email).Scan(&user).Error; err != nil {
		log.Println(err, "useriimejl")
		throwStatusUnauthorized(c)
		return
	}

	log.Println(user)

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

	user.Password = ""
	user.Token = tokenString

	throwStatusOk(user, c)
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

func addPrivilegesToUser(c *gin.Context) {
	uuid := &UUID{}
	if bindErr := c.BindJSON(&uuid); bindErr != nil {
		log.Println(bindErr)
		throwStatusBadRequest(bindErr.Error(), c)
		return
	}

	user, err := getUserFromToken(getTokenFromRequest(c))
	if err != nil {
		log.Println(err)
		throwStatusUnauthorized(c)
		return
	}

	controllers, err := getMicroControllerByUserID(user.ID)
	if err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}


	for _, controller := range controllers {
		query := "INSERT INTO users_microcontrollers(user_id, controller_id) VALUES ((SELECT id FROM users WHERE uuid = ?), ?)"
		if err := db.Exec(query, uuid.Text, controller.ID).Error; err != nil {
			log.Println(err)
			continue
		}
	}

	throwStatusOk("OK", c )

}

func deletePrivilegesToUser(c *gin.Context) {
	var userId = UserID{}
	if bindErr := c.BindJSON(&userId); bindErr != nil {
		log.Println(bindErr)
		throwStatusBadRequest(bindErr.Error(), c)
		return
	}
	log.Println(userId)

	user, err := getUserFromToken(getTokenFromRequest(c))
	if err != nil {
		log.Println(err)
		throwStatusUnauthorized(c)
		return
	}
	log.Println(user)

	if(strings.Compare(userId.Text, strconv.Itoa(user.ID)) == 0){
		throwStatusBadRequest("Can't delete yourself", c)
	}

	controllers, err := getMicroControllerByUserID(user.ID)
	if err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}


	for _, controller := range controllers {
		query := "DELETE FROM users_microcontrollers WHERE user_id = ? AND controller_id = ?"
		print(query, userId.Text, controller.ID)
		if err := db.Debug().Exec(query, userId.Text, controller.ID).Error; err != nil {
			log.Println(err)
			throwStatusInternalServerError("Error while deleting", c)
		}
	}

	throwStatusOk("OK", c )

}

func getConnectedUsers(c *gin.Context) {
	creds, err := getUserFromToken(getTokenFromRequest(c))
	if err != nil {
		log.Println(err)
		throwStatusUnauthorized(c)
		return
	}

	var users []UserV2

	query := "select u.* from users u left join users_microcontrollers um on um.user_id = u.id where um.controller_id IN (select controller_id from users_microcontrollers where user_id = ?) GROUP BY u.id"

	if err := db.Debug().Raw(query, creds.ID).Scan(&users).Error; err != nil {
		log.Println(err, "useriimejl")
		throwStatusUnauthorized(c)
		return
	}
	print(users)
	throwStatusOk(users, c)
	return

}



func AdminMiddlerware(c *gin.Context){
	if true {
		c.Next()
	}else {
		throwStatusUnauthorized(c)
	}
}