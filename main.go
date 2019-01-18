package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
	"fmt"
	"net/url"
	"strings"
	"strconv"
)

const ARDUINO_ADDRESS = "http://epcez.myddns.rocks:3000"

func main() {

	connectToDatabase()

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "HI")
	})
	router.POST("/keep-alive", keepAliveHandler)

	router.POST("/signin", SignIn)
	router.POST("/signup", SignUp)
	router.POST("/invite", addPrivilegesToUser)

	router.GET("/action", getActionsForUser)
	router.POST("/trigger/:id", triggerAction)

	router.GET("/log")

	router.GET("/state", getStates)
	router.GET("/state/:id", getStateById)

	logsApi := router.Group("/log")

	logsApi.GET("/user", getLogsForUser)
	logsApi.GET("/action/:id", getLogsForAction)


	router.Run(":8080")

	doEvery(60*time.Second, keepAlive)

}

func keepAliveHandler(c *gin.Context){
	keepAlive(time.Time{})
}

func keepAlive(t time.Time) {
	apiUrl := "http://epcez.myddns.rocks:3000"
	fmt.Println("URL:>", apiUrl)

	data := url.Values{}

	client := &http.Client{}
	r, _ := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(r)
	fmt.Println(resp)

}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func action(c *gin.Context) {
	apiUrl := ARDUINO_ADDRESS
	fmt.Println("URL:>", apiUrl)

	data := url.Values{}
	data.Set("led_4", "1")


	client := &http.Client{}
	r, _ := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Authorization", "059b9576-89ea-468e-81fb-564d1331055c")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(r)
	fmt.Println(resp.Status)

	c.JSON(200, gin.H{
		"status": resp.Status,
	})

	return

}