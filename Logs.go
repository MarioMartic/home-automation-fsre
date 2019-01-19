package main

import (
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"log"
)

type Log struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	ActionID int `json:"action_id"`
	Log string `json:"log"`
	CreatedAt time.Time `json:"created_at"`
}

func getLogsForUser(c *gin.Context){
	user, err := getUserFromToken(getTokenFromRequest(c))
	if err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	var logs []Log
	if err := db.Raw("SELECT * FROM logs WHERE user_id=?", user.ID).Scan(&logs).Error; err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	throwStatusOk(logs, c)
	return
}

func getLogsForAction(c *gin.Context) {
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

	log.Println(user);

	query := `SELECT * FROM logs WHERE action_id = ?`

	var logs []Log
	if err := db.Raw(query, id).Scan(&logs).Error; err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	throwStatusOk(logs, c)
	return
}

func logData(user User, action Action) error {
	message := "User " + user.FullName + " with ID: " + strconv.Itoa(user.ID) + "triggered action " + action.Name
	if err := db.Exec("INSERT INTO logs(user_id, action_id, log) VALUES(?,?,?)", user.ID, action.ID, message).Error; err != nil {
		return err
	}
	return nil
}

func AdminGetLogs(c *gin.Context){
	var logs []Log
	if err := db.Raw("SELECT * FROM logs").Scan(&logs).Error; err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	throwStatusOk(logs, c)
	return
}

func AdminDeleteLog(c *gin.Context){
	var logData Log

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	if err := db.Raw("SELECT * FROM logs WHERE id = ?", id).Scan(&logData).Error; err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	if err := db.Where("id = ?", id).Delete(&logData).Error; err != nil {
		log.Println(err)
		throwStatusInternalServerError(err.Error(), c)
		return
	}

	throwStatusOk(logData, c)
	return
}