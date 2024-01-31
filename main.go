package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// read in config as global variable
var mc = readMikuservConfig("./mikuservConfig.json")

func main() {

	router := gin.Default()
	router.GET("/mikuserv/ping", ping)
	router.GET("/mikuserv/stravaToken", stravaToken)
	router.POST("/mikuserv/contact", contact)
	router.Run(":" + mc["servicePort"])
}

type CommentForm struct {
	Email    string
	Text     string
	Time     string
	Password string
}

func ping(c *gin.Context) {
	responseStruct := struct {
		MainMessage   string
		ConfigMessage string
	}{
		MainMessage:   "Mikuserv is running!",
		ConfigMessage: mc["pingMessage"],
	}
	fmt.Println(responseStruct)
	c.IndentedJSON(http.StatusOK, responseStruct)
}

func contact(c *gin.Context) {
	jsonBytes, err := io.ReadAll(c.Request.Body)

	if err != nil {
		panic(err)
	}
	var form CommentForm
	json.Unmarshal(jsonBytes, &form)

	if form.Password != "" {
		c.String(http.StatusForbidden, "Error in form")
		return
	}

	fmt.Printf("%T", form.Password)

	var body bytes.Buffer
	body.WriteString(`{
		"from": {"email":"` + mc["contactFormFromEmail"] + `"},
		"personalizations":[
			{
			"to":[{"email":"` + mc["contactFormToEmail"] + `"}],
			"dynamic_template_data":{"email":"` + form.Email + `","time":"` + form.Time + `","message":"` + form.Text + `"}
			}
		],
		"template_id":"d-62476eb8af63450a8498b9a439d9016d"}`)

	req, err := http.NewRequest(http.MethodPost, "https://api.sendgrid.com/v3/mail/send", &body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+mc["sendgridApiKey"])
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusAccepted {
		c.String(http.StatusInternalServerError, "Unexpected response from Sendgrid: Code %d, Body %s", resp.StatusCode, resp.Body)
	} else {
		c.String(http.StatusAccepted, "")
	}
	fmt.Println(resp)
}

func stravaToken(c *gin.Context) {
	code := c.Query("code")

	jsonBody := []byte(`{"client_id": "102789", "client_secret": "` + mc["stravaClientSecret"] + `", "code": "` + code + `", "grant_type":"authorization_code"}`)
	//bodyReader := bytes.NewReader(jsonBody)

	fmt.Println(string(jsonBody))

	req, err := http.NewRequest(http.MethodPost, "https://www.strava.com/api/v3/oauth/token?client_id=102789&client_secret="+mc["stravaClientSecret"]+"&code="+code+"&grant_type=authorization_code", nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error encountered while receiving response from Strava")
		c.IndentedJSON(http.StatusInternalServerError, "error from http request to Strava")
	}
	defer resp.Body.Close()

	fmt.Println("Status: "+resp.Status, "Headers : ", resp.Header, "Transfer Encoding: ", resp.TransferEncoding)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "error reading response body")
	}

	c.IndentedJSON(http.StatusOK, body)
}

func readMikuservConfig(configPath string) map[string]string {
	dat, err := os.ReadFile(configPath)
	if err != nil {
		panic("Failed to read mikuserv config file at" + configPath)
	}
	configMap := make(map[string]string)
	json.Unmarshal(dat, &configMap)
	return configMap
}
