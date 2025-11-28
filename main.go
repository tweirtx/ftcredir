package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type eventCodeList struct {
	EventCodes []string `json:"eventCodes"`
}

type eventResponse struct {
	EventCode   string `json:"eventCode"`
	EventName   string `json:"name"`
	EventType   string `json:"type"`
	EventStatus string `json:"status"`
	IsFinals    bool   `json:"finals"`
	Division    int    `json:"division"`
	EventStart  int64  `json:"start"`
	EventEnd    int64  `json:"end"`
	FieldCount  int    `json:"fieldCount"`
}

func main() {
	var selected_code string
	listResp, err := http.Get("http://localhost/api/v1/events")
	if err != nil {
		log.Fatal(err)
	}

	/* 	body, err := io.ReadAll(resp.Body)
	   	if err != nil {
	   		log.Fatal(err)
	   	}
	   	fmt.Println(string(body)) */

	var codes eventCodeList
	json.NewDecoder(listResp.Body).Decode(&codes)

	for index, code := range codes.EventCodes {
		resp, err := http.Get("http://localhost/api/v1/events/" + code)
		if err != nil {
			log.Fatal(err)
		}
		var parsed_response eventResponse
		json.NewDecoder(resp.Body).Decode(&parsed_response)
		// fmt.Println(parsed_response)

		start := time.Unix(parsed_response.EventStart/1000, 0)
		now := time.Now()
		end := time.Unix(parsed_response.EventEnd/1000, 0)
		one_day, _ := time.ParseDuration("24h")
		end = end.Add(one_day) // Add a day to make the timestamp normalized

		if start.Before(now) && end.After(now) {
			fmt.Println(parsed_response.EventName + " is in date range")
			if parsed_response.EventStatus != "Archive" {
				fmt.Println(parsed_response.EventName + " is not in Archive status. Therefore, it is likely the correct event!")
				if selected_code == "" {
					selected_code = parsed_response.EventCode
				} else {
					fmt.Println("Uh-oh, looks like there's multiple events in progress...We'll eventually do something about that")
				}
			}
		}

		fmt.Println("Event code " + code + " at index " + strconv.Itoa(index) + " is in the status " + parsed_response.EventStatus)
	}
	http.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(strings.Split(r.RemoteAddr, ":")[0])
		http.Redirect(w, r, "http://localhost/event/"+selected_code+"/display/", http.StatusFound)
	})
	http.ListenAndServe(":8080", nil)
	// fmt.Println(codes.EventCodes)
}
