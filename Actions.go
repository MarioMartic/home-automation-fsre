package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"fmt"
	"net/url"
	"net/http"
	"strings"
)

type Action struct {
	ID int `json:"id"`
	Pin int `json:"pin"`
	Type int `json:"type"`
	Name string `json:"name"`
	ControllerID int `json:"controller_id"`
}

type ActionResponse struct {
	ID int `json:"id"`
	Pin int `json:"pin"`
	Type int `json:"type"`
	Name string `json:"name"`
	MicroController MicroController `json:"micro_controller" gorm:"foreignkey:controller_id;association_foreignkey:id"`
}

type UserMicroController struct {
	UserID int `json:"user_id"`
	MicroControllerID int `json:"micro_controller_id"`
}

func getActionsForUser(c *gin.Context){
	token := getTokenFromRequest(c)
	user, err := getUserFromToken(token)
	if err != nil {
		log.Println(err)
		throwStatusUnauthorized(c)
		return
	}

	query := `SELECT * FROM actions a 
	JOIN microcontrollers m ON a.controller_id = m.id 
	JOIN users_microcontrollers um ON m.id = um.controller_id 
	JOIN users u ON um.user_id = u.id WHERE u.id=?`

	var actions []Action

	if err := db.Debug().Raw(query, user.ID).Scan(&actions).Error; err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	throwStatusOk(actions, c)
	return
}

func triggerAction(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	user, err := getUserFromToken(getTokenFromRequest(c))
	if err != nil {
		log.Println(err)
		throwStatusUnauthorized(c)
		return
	}

	log.Println("USER;", user)
	var action Action

	errr := db.Debug().Raw("SELECT * FROM actions WHERE id = ?", id).Scan(&action).Error
	if errr != nil {
		log.Println(errr.Error())
		throwStatusUnauthorized(c)
		return
	}

	var umc []UserMicroController
	errMC := db.Debug().Raw("SELECT * FROM users_microcontrollers WHERE user_id = ? AND controller_id=?", user.ID, action.ControllerID).Scan(&umc).Error
	if errMC != nil {
		log.Println("Count = 0 asdasd")
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	if len(umc) == 0 {
		throwStatusUnauthorized(c)
		return
	}

	var controller MicroController

	if err := db.Debug().Raw("SELECT * FROM microcontrollers WHERE id = ?", action.ControllerID).Scan(&controller).Error; err != nil{
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	if err := sendRequest(controller, action); err != nil {
		throwStatusInternalServerError(err.Error(), c)
		log.Println(err)
		return
	}

	if err := logAction(action, user); err != nil {
		throwStatusInternalServerError(err.Error(), c)
		log.Println(err)
		return
	}

	throwStatusOk("OK", c)

}

func sendRequest(controller MicroController, action Action) error {
	apiUrl := "http://" + controller.Domain + ":" + strconv.Itoa(controller.Port)
	fmt.Println("URL:>", apiUrl)

	data := url.Values{}
	data.Set("pin_" + strconv.Itoa(action.Pin), strconv.Itoa(action.Type))


	client := &http.Client{}
	r, _ := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Authorization", controller.Token)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)

	return nil

}

func logAction(action Action, user User) error {

	var log = "User " + user.FullName + " with id = " + strconv.Itoa(user.ID) + " triggered action named '" + action.Name + "' on microcontroller with id = " + strconv.Itoa(action.ControllerID)
	if err := db.Debug().Exec("INSERT INTO logs (user_id, action_id, log) VALUES (?, ?, ?)", user.ID, action.ID, log).Error; err != nil {
		return err
	}

	return nil
}