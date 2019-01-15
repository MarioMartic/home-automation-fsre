package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type MicroController struct {
	ID	int `json:"id"`
	Name string `json:"name"`
	Token string `json:"token"`
	Domain string `json:"domain"`
	Port int `json:"port"`
}

func AdminGetMicroControllers(c *gin.Context){
	var microcontrollers []MicroController
	if err := db.Raw("SELECT * FROM microcontrollers").Scan(&microcontrollers).Error; err != nil {
		throwStatusInternalServerError(err.Error(), c)
		return
	}
	throwStatusOk(microcontrollers, c)
	return
}

func AdminGetMicroControllerByID(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	microcontroller, err := getMicroControllerByID(id)
	if err != nil {
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	throwStatusOk(microcontroller, c)
	return
}

func AdminUpdateMicroControllerByID(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	microcontroller, err := getMicroControllerByID(id)
	if err != nil {
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	if err := db.Exec("UPDATE microcontollers SET name=?, token=?, domain=?, port=? WHERE id =?", microcontroller.Name, microcontroller.Token, microcontroller.Domain, microcontroller.Port, id).Error; err != nil {
		throwStatusInternalServerError(err.Error(), c)
		return
	}
	throwStatusOk(microcontroller, c)
	return
}

func AdminDeleteMicroControllerByID(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	microcontroller, err := getMicroControllerByID(id)
	if err != nil {
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	if err := db.Exec("DELETE FROM microcontollers WHERE id =?", microcontroller.ID).Error; err != nil {
		throwStatusInternalServerError(err.Error(), c)
		return
	}
	throwStatusOk(RESULT_SUCCESS, c)
	return
}


func getMicroControllerByID(id int)(MicroController, error){
	var microcontroller MicroController
	if err := db.Raw("SELECT * FROM microcontrollers WHERE id =?", id).Scan(&microcontroller).Error; err != nil {
		return MicroController{}, err
	}
	return microcontroller, nil
}
