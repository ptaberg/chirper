package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-vk-api/vk"
	"golang.org/x/oauth2"
	vkAuth "golang.org/x/oauth2/vk"
)

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Photo     string `json:"photo_400_orig"`
	City      City   `json:"city"`
}

type City struct {
	Title string `json:"title"`
}

type Status struct {
	Text string `json:"text"`
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	conf := &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:8090/auth",
		Scopes:       []string{"photos", "status"},
		Endpoint:     vkAuth.Endpoint,
	}

	r.GET("/", func(c *gin.Context) {
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"authUrl": url,
		})
	})

	r.GET("/auth", func(c *gin.Context) {
		ctx := context.Background()
		authCode := c.Request.URL.Query()["code"]
		tok, err := conf.Exchange(ctx, authCode[0])
		if err != nil {
			log.Fatal(err)
		}

		client, err := vk.NewClientWithOptions(vk.WithToken(tok.AccessToken))
		if err != nil {
			log.Fatal(err)
		}

		user := getCurrentUser(client)
		status := getUserStatus(client, user.ID)

		c.HTML(http.StatusOK, "auth.html", gin.H{
			"user":   user,
			"status": status,
		})

		fmt.Println(user.ID)
		fmt.Println(status)
	})

	r.Run(":8090")
}

func getUserStatus(api *vk.Client, id int64) Status {
	var status Status

	_ = api.CallMethod("status.get", vk.RequestParams{
		"user_id": strconv.Itoa(int(id)),
	}, &status)

	return status
}

func getCurrentUser(api *vk.Client) User {
	var users []User

	_ = api.CallMethod("users.get", vk.RequestParams{
		"fields": "city,photo_400_orig",
	}, &users)

	fmt.Println(users[0])
	return users[0]
}
