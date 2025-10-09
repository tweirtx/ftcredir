package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	EventStart  int    `json:"start"`
	EventEnd    int    `json:"end"`
	FieldCount  int    `json:"fieldCount"`
}

func main() {
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

		fmt.Println("Event code " + code + " at index " + strconv.Itoa(index) + " is in the status " + parsed_response.EventStatus)
	}
	fmt.Println(codes.EventCodes)
}
