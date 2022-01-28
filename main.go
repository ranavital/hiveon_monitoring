package main

import (
	"encoding/json"
	"fmt"
    "net/http"
)



func main() {
    resp, err := http.Get("https://hiveon.net/api/v1/stats/miner/72fc9a5770cd96f2686c816fd3672840feb96364/ETH/workers")
    if err != nil {
	    panic(err)
    }

    var responseJson = map[string]map[string]map[string]interface{}{}

	if err := json.NewDecoder(resp.Body).Decode(&responseJson); err != nil {
		panic(err)
	}


	for k, v := range responseJson["workers"] {
		if v["online"] != true{
			fmt.Printf("Worker %s is offline\n", k)
		}
	}
}
