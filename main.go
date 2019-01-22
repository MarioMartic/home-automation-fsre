package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
	"fmt"
	"net/url"
	"strings"
	"strconv"
	"github.com/gin-contrib/cors"
)

const ARDUINO_ADDRESS = "http://epcez.myddns.rocks:3000"
var handler http.Handler

func main() {

	connectToDatabase()

	router := gin.New()
	router.Use(gin.Logger())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"*"}
	config.AllowMethods = []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"}

	router.Use(cors.New(config))

	/*
	router.Use(func(c *gin.Context) {
			// Run this on all requests   
			// Should be moved to a proper middleware 
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Token")
	c.Next()
	})

	router.OPTIONS("/*cors", func(c *gin.Context) {
		// Empty 200 response
	}) */

	
	router.POST("/keep-alive", keepAliveHandler)

	router.POST("/signin", SignIn)
	router.POST("/signup", SignUp)
	router.POST("/invite", addPrivilegesToUser)
	router.GET("/user/connected", getConnectedUsers)
	router.POST("/user/delete", deletePrivilegesToUser)

	router.GET("/action", getActionsForUser)
	router.POST("/trigger/:id", triggerAction)

	router.GET("/log")

	router.GET("/state", getStates)
	router.GET("/state/:id", getStateById)

	logsApi := router.Group("/log")

	logsApi.GET("/user", getLogsForUser)
	logsApi.GET("/action/:id", getLogsForAction)

	router.POST("/admin/signin", AdminSignIn)

	adminApi := router.Group("/admin", AdminMiddleware)
	{

		adminApi.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, "HI")
		})

		controllerApi := adminApi.Group("/controllers")
		{
			controllerApi.POST("/", AdminCreateMicroController)
			controllerApi.GET("/", AdminGetMicroControllers)
			controllerApi.GET("/:id", AdminGetMicroControllerByID)
			controllerApi.PUT("/:id", AdminUpdateMicroControllerByID)
			controllerApi.DELETE("/:id", AdminDeleteMicroControllerByID)
			controllerApi.POST("/bind", bindUserWithController)
		}

		userApi := adminApi.Group("/users")
		{
			userApi.POST("/", AdminCreateUser)
			userApi.GET("/", AdminGetUsers)
			userApi.GET("/:id", AdminGetUser)
			userApi.PUT("/:id", AdminUpdateUser)
			userApi.DELETE("/:id", AdminDeleteUser)
		}

		actionApi := adminApi.Group("/actions")
		{
			actionApi.POST("/", AdminCreateAction)
			actionApi.GET("/", AdminGetActions)
			actionApi.GET("/:id", AdminGetAction)
			actionApi.PUT("/:id", AdminUpdateAction)
			actionApi.DELETE("/:id", AdminDeleteAction)
		}

		logsAdminApi := adminApi.Group("/logs")
		{
			logsAdminApi.GET("/", AdminGetLogs)
			logsAdminApi.DELETE("/:id", AdminDeleteLog)
		}
	}

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