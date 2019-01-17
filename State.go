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

	client := &http.Client{}
	req, err := http.NewRequest("GET", ARDUINO_ADDRESS, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	q.Add("pin", strconv.Itoa(id))
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