package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

func init() {
	viper.SetEnvPrefix("sms")

	viper.BindEnv("mattermost_webhook_url")
	viper.BindEnv("host")
	viper.BindEnv("port")

	viper.SetDefault("addr", "0.0.0.0")
	viper.SetDefault("port", "1323")

	if viper.GetString("mattermost_webhook_url") == "" {
		log.Fatalf("SMS_MATTERMOST_WEBHOOK_URL environment variable must be set!")
	}
}

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/:channel", func(c *gin.Context) {
		channel := c.Param("channel")

		jsonByteData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Fatalf("Error reading body: %v", err)
		}
		jsonStringData := string(jsonByteData)
		// log.Printf("The value of JSON is: %s", jsonStringData)
		title := gjson.Get(jsonStringData, "event.title").String()

		var customer_info string
		source := gjson.Get(jsonStringData, "event.contexts.Research Source.research_source").String()
		scheme := gjson.Get(jsonStringData, "event.contexts.Customer Scheme.customer_scheme").String()
		name := gjson.Get(jsonStringData, "event.contexts.Customer Name.customer_name").String()
		if source == "unknown" || source == "" {
			customer_info = strings.Join([]string{scheme, name}, "\n")
		} else {
			customer_info = source
		}
		headers := gjson.Get(jsonStringData, "event.request.headers").Array()
		var study_uid string
		for _, header := range headers {
			headerArray := header.Array()
			if headerArray[0].String() == "Study-Instance-Uid" {
				study_uid = headerArray[1].String()
				break
			}
		}

		project_name := gjson.Get(jsonStringData, "project").String()
		var alert_text string 
		if project_name == "breastcancer" {
			alert_text = ":beda::beda::beda::beda::beda::beda::beda::beda::beda::beda:"
		} else {
			alert_text = ":grozny_ebaka:"
		}

		postBody, err := json.Marshal(map[string]interface{}{
			"channel": channel,
			"text": alert_text,
			"attachments": []interface{}{
				map[string]interface{}{
					"title":       title,
					"fallback":    title,
					"color":       "#FF0000",
					"author_name": "Leshtry",
					"author_icon": "https://i.ibb.co/QPz299J/photo-2023-12-17-00-34-13.jpg",
					"title_link":  gjson.Get(jsonStringData, "url").String(),
					"fields": []interface{}{
						map[string]interface{}{
							"short": false,
							"title": "Environment",
							"value": gjson.Get(jsonStringData, "event.environment").String(),
						},
						map[string]interface{}{
							"short": false,
							"title": "Customer Info",
							"value": customer_info,
						},
						map[string]interface{}{
							"short": false,
							"title": "StudyInstanceUID",
							"value": study_uid,
						},
					},
				},
			},
		})
		if err != nil {
			log.Fatalf("Error during json marshal: %v", err)
		}

		resp, err := http.Post(
			viper.GetString("mattermost_webhook_url"),
			"application/json",
			bytes.NewBuffer(postBody),
		)
		if err != nil {
			log.Fatalf("Error when performing webhook call: %v", err)
		}
		defer resp.Body.Close()
	})

	r.Run(fmt.Sprintf(
		"%s:%s",
		viper.GetString("host"),
		viper.GetString("port"),
	))
}
