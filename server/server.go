package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"io/ioutil"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"gopkg.in/yaml.v2"
)

type viptella_alert struct {
	Values []struct {
		SystemIP string `json:"system-ip"`
		SiteID   string `json:"site-id"`
		HostName string `json:"host-name"`
	} `json:"values"`
	Message  string `json:"message"`
	Type     string `json:"type"`
	Severity string `json:"severity"`
}

type config struct {
    Webhook string `yaml:"webhook"`
    Port string `yaml: "port"`
}

func (c *config) getConf() *config {

    yamlFile, err := ioutil.ReadFile("config.yaml")
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, c)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return c
}



func alertHandler(w http.ResponseWriter, r *http.Request) {

	// Create decoder to decode JSON body
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	// Create variable for viptella alert structure
	var t viptella_alert

	// Decode the request based on the alert structure
	err := d.Decode(&t)
	d.DisallowUnknownFields()

	// print error if it occurs during decoding.
	if err != nil {
		fmt.Println(err)
	}

	if r.URL.Path != "/alerts" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	// print the data in the request that matches our structure
	fmt.Println("Alert Message", t.Message)
	fmt.Println("Alert Severity", t.Severity)
	fmt.Println("Alert Hostname", t.Values[0].HostName)
	fmt.Println("Alert System IP", t.Values[0].SystemIP)
	fmt.Println("Alert SiteId", t.Values[0].SiteID)

	// Initialize a new Microsoft Teams client.
	mstClient := goteamsnotify.NewClient()

	// Set webhook url.
	var cfg config
	cfg.getConf()

	webhookUrl := cfg.Webhook

	// Setup message card.
	msgCard := goteamsnotify.NewMessageCard()
	msgCard.Title = t.Severity + " vManage Alert!"
	msgCard.Text = t.Message + "<br><br>" + "<h3>Details:</h3><blockquote>Hostname: " + t.Values[0].HostName + "<br>   System IP: " + t.Values[0].SystemIP + "<br>   Site ID: " + t.Values[0].SiteID + "</blockquote>"
	msgCard.ThemeColor = "#DF813D"

	// Send the message with default timeout/retry settings.
	if err := mstClient.Send(webhookUrl, msgCard); err != nil {
		log.Printf("failed to send message: %v", err)
		os.Exit(1)
	}

}

func main() {
    var cfg config
    cfg.getConf()
	http.HandleFunc("/alerts", alertHandler)

	fmt.Println("Server can be accessed on the following port: " + cfg.Port)


	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatal(err)
	}

}
