package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/resend/resend-go/v2"
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
	Date     string
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
		c.String(http.StatusInternalServerError, "Bad JSON from form request")
		return
	}
	var form CommentForm
	json.Unmarshal(jsonBytes, &form)

	if form.Password != "" {
		c.String(http.StatusForbidden, "Error in form")
		return
	}
	fmt.Println("JSON received: " + string(jsonBytes))

	client := resend.NewClient(mc["resendApiKey"])

	// Send
	params := &resend.SendEmailRequest{
		From:    "tmiku.net Notifications <" + mc["contactFormFromEmail"] + ">",
		To:      []string{mc["contactFormToEmail"]},
		Subject: "Contact form submitted (<" + form.Email + ">)",
		Html: `
		A reader has submitted a contact form. Their address is <` + form.Email + `> and they submitted the form at ` + form.Date + `. Message below.<br/>
		<br/>
		------------------------------------<br/>
		` + form.Text + `<br/>
		------------------------------------<br/>
		`,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		c.String(http.StatusInternalServerError, "API rejected email request")
		return
	}

	fmt.Println("Email sent successfully:", sent.Id, "at", time.Now().Format(time.RFC3339))
	c.String(http.StatusAccepted, "")
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
