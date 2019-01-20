package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"log"
	"net/url"
	"net/http"
)

type MicroController struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Token  string `json:"token"`
	Domain string `json:"domain"`
	Port   int    `json:"port"`
}

func AdminGetMicroControllers(c *gin.Context) {
	var microcontrollers []MicroController
	if err := db.Raw("SELECT * FROM microcontrollers").Scan(&microcontrollers).Error; err != nil {
		throwStatusInternalServerError(err.Error(), c)
		return
	}
	throwStatusOk(microcontrollers, c)
	return
}

func AdminGetMicroControllerByID(c *gin.Context) {
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

func AdminUpdateMicroControllerByID(c *gin.Context) {
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

	if err := c.BindJSON(&microcontroller); err != nil {
		log.Println(err)
		return
	}

	validationErrors := microcontroller.validate()
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, validationErrors)
		return
	}

	if err := db.Exec("UPDATE microcontollers SET name=?, token=?, domain=?, port=? WHERE id =?", microcontroller.Name, microcontroller.Token, microcontroller.Domain, microcontroller.Port, id).Error; err != nil {
		throwStatusInternalServerError(err.Error(), c)
		return
	}
	throwStatusOk(microcontroller, c)
	return
}

func AdminDeleteMicroControllerByID(c *gin.Context) {
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

func AdminCreateMicroController(c *gin.Context) {
	var controller MicroController
	if err := c.BindJSON(&controller); err != nil {
		log.Println(err)
		return
	}

	count := db.Raw("SELECT * FROM microcontrollers WHERE domain = ? AND port = ?", controller.Domain, controller.Port).RowsAffected

	if count != 0 {
		throwStatusBadRequest("Controller already exists", c)
	}

	validationErrors := controller.validate()
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, validationErrors)
		return
	}

	if err := db.Create(&controller); err != nil {
		log.Println(err)
		return
	}

	c.JSON(200, controller)

}

func (m *MicroController) validate() url.Values {
	errs := url.Values{}

	if m.Token == "" {
		errs.Add("token", "Token is required!")
	}

	if m.Port == 0 || m.Port < 1000 || m.Port > 15000 {
		errs.Add("port", "Invalid port number!")
	}

	if m.Domain == ""  {
		errs.Add("domain", "Domain is required!")
	}

	return errs
}

func getMicroControllerByID(id int) (MicroController, error) {
	var microcontroller MicroController
	if err := db.Raw("SELECT * FROM microcontrollers WHERE id =?", id).Scan(&microcontroller).Error; err != nil {
		return MicroController{}, err
	}
	return microcontroller, nil
}

func getMicroControllerByUserID(id int) ([]MicroController, error) {
	var microcontrollers []MicroController
	if err := db.Raw("SELECT * FROM microcontrollers m JOIN users_microcontrollers um ON m.id = um.controller_id WHERE um.user_id =?", id).Scan(&microcontrollers).Error; err != nil {
		return nil, err
	}
	return microcontrollers, nil
}
