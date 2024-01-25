/*
CISCO SAMPLE CODE LICENSE Version 1.1 Copyright (c) 2020 Cisco and/or its affiliates

These terms govern this Cisco Systems, Inc. ("Cisco"), example or demo source code and its associated documentation (together, the "Sample Code"). By downloading, copying, modifying, compiling, or redistributing the Sample Code, you accept and agree to be bound by the following terms and conditions (the "License"). If you are accepting the License on behalf of an entity, you represent that you have the authority to do so (either you or the entity, "you"). Sample Code is not supported by Cisco TAC and is not tested for quality or performance. This is your only license to the Sample Code and all rights not expressly granted are reserved.

LICENSE GRANT: Subject to the terms and conditions of this License, Cisco hereby grants to you a perpetual, worldwide, non-exclusive, non- transferable, non-sublicensable, royalty-free license to copy and modify the Sample Code in source code form, and compile and redistribute the Sample Code in binary/object code or other executable forms, in whole or in part, solely for use with Cisco products and services. For interpreted languages like Java and Python, the executable form of the software may include source code and compilation is not required.

CONDITIONS: You shall not use the Sample Code independent of, or to replicate or compete with, a Cisco product or service. Cisco products and services are licensed under their own separate terms and you shall not use the Sample Code in any way that violates or is inconsistent with those terms (for more information, please visit: www.cisco.com/go/terms).

OWNERSHIP: Cisco retains sole and exclusive ownership of the Sample Code, including all intellectual property rights therein, except with respect to any third-party material that may be used in or by the Sample Code. Any such third-party material is licensed under its own separate terms (such as an open source license) and all use must be in full accordance with the applicable license. This License does not grant you permission to use any trade names, trademarks, service marks, or product names of Cisco. If you provide any feedback to Cisco regarding the Sample Code, you agree that Cisco, its partners, and its customers shall be free to use and incorporate such feedback into the Sample Code, and Cisco products and services, for any purpose, and without restriction, payment, or additional consideration of any kind. If you initiate or participate in any litigation against Cisco, its partners, or its customers (including cross-claims and counter-claims) alleging that the Sample Code and/or its use infringe any patent, copyright, or other intellectual property right, then all rights granted to you under this License shall terminate immediately without notice.

LIMITATION OF LIABILITY: CISCO SHALL HAVE NO LIABILITY IN CONNECTION WITH OR RELATING TO THIS LICENSE OR USE OF THE SAMPLE CODE, FOR DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO DIRECT, INCIDENTAL, AND CONSEQUENTIAL DAMAGES, OR FOR ANY LOSS OF USE, DATA, INFORMATION, PROFITS, BUSINESS, OR GOODWILL, HOWEVER CAUSED, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGES.

DISCLAIMER OF WARRANTY: SAMPLE CODE IS INTENDED FOR EXAMPLE PURPOSES ONLY AND IS PROVIDED BY CISCO "AS IS" WITH ALL FAULTS AND WITHOUT WARRANTY OR SUPPORT OF ANY KIND. TO THE MAXIMUM EXTENT PERMITTED BY LAW, ALL EXPRESS AND IMPLIED CONDITIONS, REPRESENTATIONS, AND WARRANTIES INCLUDING, WITHOUT LIMITATION, ANY IMPLIED WARRANTY OR CONDITION OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON- INFRINGEMENT, SATISFACTORY QUALITY, NON-INTERFERENCE, AND ACCURACY, ARE HEREBY EXCLUDED AND EXPRESSLY DISCLAIMED BY CISCO. CISCO DOES NOT WARRANT THAT THE SAMPLE CODE IS SUITABLE FOR PRODUCTION OR COMMERCIAL USE, WILL OPERATE PROPERLY, IS ACCURATE OR COMPLETE, OR IS WITHOUT ERROR OR DEFECT.

GENERAL: This License shall be governed by and interpreted in accordance with the laws of the State of California, excluding its conflict of laws provisions. You agree to comply with all applicable United States export laws, rules, and regulations. If any provision of this License is judged illegal, invalid, or otherwise unenforceable, that provision shall be severed and the rest of the License shall remain in full force and effect. No failure by Cisco to enforce any of its rights related to the Sample Code or to a breach of this License in a particular situation will act as a waiver of such rights. In the event of any inconsistencies with any other terms, this License shall take precedence.
*/
package main

import (
        "encoding/json"
        "fmt"
        "log"
        "net/http"
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
        // Create variable for viptella alert structure
        var t viptella_alert

        // Decode the request based on the alert structure
        err := d.Decode(&t)

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
                return
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
