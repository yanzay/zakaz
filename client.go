package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const apiURL = "https://metro.zakaz.ua/api/query.json"

type Args struct {
	StoreIds      []string `json:"store_ids"`
	OnlyAvailable bool     `json:"only_available"`
	ZoneID        string   `json:"zone_id"`
	DeliveryType  string   `json:"delivery_type"`
}

type SubRequest struct {
	Args Args   `json:"args"`
	V    string `json:"v"`
	Type string `json:"type"`
	ID   string `json:"id"`
}

type Request struct {
	Request []SubRequest `json:"request"`
}

type DayWindows struct {
	Date    string `json:"date"`
	Windows []struct {
		Status         int    `json:"status"`
		ContractZoneID string `json:"contract_zone_id"`
		Price          struct {
			Num0 int `json:"0"`
		} `json:"price"`
		EndOrderingTime   string  `json:"end_ordering_time"`
		TzOffset          float64 `json:"tz_offset"`
		IsAvailable       bool    `json:"is_available"`
		StartOrderingTime float64 `json:"start_ordering_time"`
		ID                string  `json:"id"`
		IsInPast          bool    `json:"is_in_past"`
		RangeTime         string  `json:"range_time"`
		Default           bool    `json:"default"`
		Title             string  `json:"title"`
		Ts                float64 `json:"ts"`
		EndOrderingTimeTs float64 `json:"end_ordering_time_ts"`
		Time              string  `json:"time"`
		StartTime         float64 `json:"start_time"`
	} `json:"windows"`
	IsDayoff bool   `json:"is_dayoff"`
	Title    string `json:"title"`
}

type Response struct {
	Meta struct {
	} `json:"meta"`
	Responses []struct {
		Data struct {
			Items []struct {
				Windows []*DayWindows `json:"windows"`
				ID      string        `json:"id"`
				ZoneID  string        `json:"zone_id"`
			} `json:"items"`
		} `json:"data"`
		Error bool `json:"error"`
	} `json:"responses"`
}

func GetWindows(storeID, zoneID string) ([]*DayWindows, error) {
	req := &Request{
		Request: []SubRequest{{
			Args: Args{
				StoreIds:      []string{storeID},
				OnlyAvailable: true,
				ZoneID:        zoneID,
				DeliveryType:  "plan",
			},
			V:    "0.1",
			Type: "timewindows.list",
			ID:   "timewindows_list",
		}},
	}
	buf := &bytes.Buffer{}
	json.NewEncoder(buf).Encode(req)
	resp, err := http.Post(apiURL, "application/json", buf)
	if err != nil {
		return nil, err
	}
	response := &Response{}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return nil, err
	}
	if len(response.Responses) == 0 || len(response.Responses[0].Data.Items) == 0 {
		return nil, fmt.Errorf("empty response")
	}
	return response.Responses[0].Data.Items[0].Windows, nil
}
