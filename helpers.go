package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"fmt"
	"crypto/rand"
)

var RESULT_SUCCESS = map[string]string{"result:":"success"}

func respondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}

	c.JSON(code, resp)
	c.AbortWithStatus(code)
}

func throwStatusOk(i interface{}, c *gin.Context) {
	if i != nil {
		c.JSON(http.StatusOK, i)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "OK",
	})
}

func throwStatusBadRequest(msg string, c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"message": msg,
	})
}

func throwStatusInternalServerError(msg string, c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"message": msg,
	})
}

func throwStatusUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"message": "ERR_USER_UNAUTHORIZED",
	})
}

func generateToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
