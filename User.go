package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"log"
)

type UserMicrocontrollers struct {
	ID               int               `json:"id"`
	FullName         string            `json:"full_name"`
	Email            string            `json:"email"`
	Password         string            `json:"-"`
	UUID             string            `json:"uuid"`
	Token            string            `json:"token"`
	Microcontrollers []MicroController `json:"microcontrollers"`
}

type UsersMicrocontrollers struct {
	User            User            `json:"user"`
	MicroController MicroController `json:"micro_controller"`
}

type MicroControllerWithUID struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Token  string `json:"token"`
	Domain string `json:"domain"`
	Port   int    `json:"port"`
	UserID int    `json:"user_id"`
	NumOfPins int `json:"number_of_pins"`
}

func AdminDeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	var user User
	if err := db.Debug().Where("id = ?", id).Delete(&user).Error; err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	c.JSON(200, gin.H{"id #" + strconv.Itoa(id): "deleted"})
}
func AdminUpdateUser(c *gin.Context) {
	var user User

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}

	if err := c.BindJSON(&user); err != nil {
		log.Println(err)
		throwStatusBadRequest(err.Error(), c)
		return
	}

	if err := db.Debug().Save(&user).Error; err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}
	c.JSON(200, user)
}

func AdminUserMicrocontrollers(c *gin.Context) {

	var users []UserMicrocontrollers

	query := "SELECT * FROM users u JOIN users_microcontrollers um ON u.id = um.user_id"

	if err := db.Debug().Raw(query).Scan(&users).Error; err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	var controllers []MicroControllerWithUID

	query = "SELECT m.*, um.user_id FROM microcontrollers m JOIN users_microcontrollers um ON m.id = um.controller_id"

	if err := db.Debug().Raw(query).Scan(&controllers).Error; err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	for index, user := range users {
		for _, controller := range controllers {
			if user.ID == controller.UserID {
				c := MicroController{controller.ID, controller.Name, controller.Token, controller.Domain, controller.Port, controller.NumOfPins}
				users[index].Microcontrollers = append(users[index].Microcontrollers, c)
			}
		}
	}

	c.JSON(200, users)
}

func AdminUsersMicrocontrollers(c *gin.Context) {
	query := "SELECT u.id, u.full_name, u.email, m.id, m.name, m.domain, m.port FROM users_microcontrollers um INNER JOIN users u ON um.user_id = u.id INNER JOIN microcontrollers m ON um.controller_id = m.id"
	var ums []UsersMicrocontrollers

	rows, errr := db.Raw(query).Rows() // (*sql.Rows, error)
	if errr != nil {
		log.Println(errr)
		throwStatusInternalServerError(errr.Error(), c)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var um UsersMicrocontrollers
		rows.Scan(&um.User.ID, &um.User.FullName, &um.User.Email, &um.MicroController.ID, &um.MicroController.Name, &um.MicroController.Domain, &um.MicroController.Port)
		ums = append(ums, um)
	}
	c.JSON(200, ums)
}

func AdminCreateUser(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		log.Println(err)
		throwStatusBadRequest(err.Error(), c)
		return
	}

	user.UUID = generateToken(16)

	if err := db.Debug().Create(&user).Error; err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}
	c.JSON(200, user)
}
func AdminGetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	var user User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	} else {
		c.JSON(200, user)
	}
}
func AdminGetUsers(c *gin.Context) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	} else {
		c.JSON(200, users)
	}
}
