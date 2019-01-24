
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"log"
	"net/url"
	"net/http"
)

func AdminDeleteAction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	var action Action
	if err := db.Debug().Where("id = ?", id).Delete(&action).Error; err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	c.JSON(200, gin.H{"id #" + strconv.Itoa(id): "deleted"})
}
func AdminUpdateAction(c *gin.Context) {
	var action Action

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	if err := db.Where("id = ?", id).First(&action).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}

	if err := c.BindJSON(&action); err != nil {
		log.Println(err)
		throwStatusBadRequest(err.Error(), c)
		return
	}

	validationErrors := action.validate()
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, validationErrors)
		return
	}
	query := "SELECT * FROM actions WHERE controller_id = ? AND pin = ?"
	var actions []Action
	count := db.Debug().Raw(query, action.ControllerID, action.Pin).Scan(&actions).RowsAffected
	if count != 0 {
		throwStatusBadRequest("ERR_PIN_DUPLICATION", c)
		return
	}
	if err := db.Debug().Save(&action).Error; err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}
	c.JSON(200, action)
}
func AdminCreateAction(c *gin.Context) {
	var action Action
	if err := c.BindJSON(&action); err != nil {
		log.Println(err)
		throwStatusBadRequest(err.Error(), c)
		return
	}

	validationErrors := action.validate()
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, validationErrors)
		return
	}

	query := "SELECT * FROM actions WHERE controller_id = ? AND pin = ?"
	var actions []Action

	count := db.Debug().Raw(query, action.ControllerID, action.Pin).Scan(&actions).RowsAffected

	if count != 0 {
		throwStatusBadRequest("ERR_PIN_DUPLICATION", c)
		return
	}

	if err := db.Debug().Create(&action).Error; err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}
	c.JSON(200, action)
}
func AdminGetAction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	var action Action
	if err := db.Where("id = ?", id).First(&action).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	} else {
	c.JSON(200, action)
	}
}
func AdminGetActions(c *gin.Context) {
	var actions []Action
	if err := db.Find(&actions).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	} else {
		c.JSON(200, actions)
	}
}

func (a *Action) validate() url.Values {
	errs := url.Values{}

	if a.ControllerID == 0 {
		errs.Add("controller_id", "Controller ID is required!")
	}

	if a.Pin == 0 || a.Pin < 2 || a.Pin > 13 {
		errs.Add("pin", "Invalid pin number!")
	}

	if a.Name == ""  {
		errs.Add("name", "Action name is required!")
	}

	return errs
}
