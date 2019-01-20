package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"log"
)

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
