package main

import (
	"log"
	"os"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"encoding/json"
	"strconv"
)

type State struct {
	Pin string `json:"pin"`
	Status string `json:"status"`
}

func getStates(c *gin.Context) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ARDUINO_ADDRESS, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", "059b9576-89ea-468e-81fb-564d1331055c")

	q := req.URL.Query()
	q.Add("pin", "0")
	req.URL.RawQuery = q.Encode()



	fmt.Println(req.URL.String())

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}
	defer resp.Body.Close()

	var s []State
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		log.Println(err.Error())
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	c.JSON(200, s)
}

func getStateById(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	_, err = getUserFromToken(getTokenFromRequest(c))
	if err != nil {
		log.Println(err)
		throwStatusUnauthorized(c)
		return
	}
	var action Action

	errr := db.Debug().Raw("SELECT * FROM actions WHERE id = ?", id).Scan(&action).Error
	if errr != nil {
		log.Println(errr.Error())
		throwStatusUnauthorized(c)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", ARDUINO_ADDRESS, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", "059b9576-89ea-468e-81fb-564d1331055c")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	q.Add("pin", strconv.Itoa(action.Pin))
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}
	defer resp.Body.Close()

	var s State
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		log.Println(err.Error())
		throwStatusInternalServerError(err.Error(), c)
		return
	}
	s.Pin = strconv.Itoa(action.Pin)
	c.JSON(200, s)

}