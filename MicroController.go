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
	NumOfPins int `json:"number_of_pins"`
}

func (MicroController) TableName() string {
	return "microcontrollers"
}

type UM struct {
	UserID int `json:"user_id"`
	ControllerID int `json:"controller_id"`
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
	var mc []MicroController
	query := "SELECT * FROM microcontrollers WHERE domain = ? AND port = ? AND id != ?"

	count := db.Debug().Raw(query, microcontroller.Domain, microcontroller.Port, microcontroller.ID).Scan(&mc).RowsAffected

	if count != 0 {
		throwStatusBadRequest("ERR_DOMAIN_PORT_DUPLICATION", c)
		return
	}

	if err := db.Exec("UPDATE microcontrollers SET name=?, token=?, domain=?, port=? WHERE id =?", microcontroller.Name, microcontroller.Token, microcontroller.Domain, microcontroller.Port, id).Error; err != nil {
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

	if err := db.Exec("DELETE FROM microcontrollers WHERE id =?", microcontroller.ID).Error; err != nil {
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

	var mc []MicroController
	count := db.Debug().Raw("SELECT * FROM microcontrollers WHERE domain = ? AND port = ?", controller.Domain, controller.Port).Scan(&mc).RowsAffected

	if count != 0 {
		throwStatusBadRequest("ERR_DOMAIN_PORT_DUPLICATION", c)
		return
	}

	validationErrors := controller.validate()
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, validationErrors)
		return
	}


	if err := db.Debug().Exec("INSERT INTO `microcontrollers` (`name`,`token`,`domain`,`port`,`number_of_pins`) " +
		"VALUES (?,?,?,?,?)", controller.Name, controller.Token, controller.Domain, controller.Port, controller.NumOfPins).Error; err != nil {
			throwStatusBadRequest(err.Error(), c)
		log.Println(err)
		return
	}

	throwStatusOk(controller, c)
	return
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

	if m.Name == "" {
		errs.Add("name", "Name is required!")
	}

	if m.NumOfPins != 13 {
		errs.Add("number_of_pins", "Invalid number of pins")
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

func bindUserWithController(c *gin.Context) {
	var um UM


	if err := c.BindJSON(&um); err != nil {
		log.Println(err)
		return
	}

	if um.UserID == 0 || um.ControllerID == 0 {
		throwStatusBadRequest("Nemoj nula", c)
		return
	}

	query := "INSERT INTO users_microcontrollers (user_id, controller_id) VALUES (?, ?)"


	if err := db.Exec(query, um.UserID, um.ControllerID).Error; err != nil {
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	throwStatusOk("OK", c)

}

func unbindUserWithController(c *gin.Context) {
	var um UM
	if err := c.BindJSON(&um); err != nil {
		log.Println(err)
		return
	}

	if um.UserID == 0 || um.ControllerID == 0 {
		throwStatusBadRequest("Nemoj nula", c)
		return
	}

	query := "DELETE FROM users_microcontrollers WHERE user_id = ? AND controller_id = ?"

	if err := db.Exec(query, um.UserID, um.ControllerID).Error; err != nil {
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	throwStatusOk("OK", c)

}
