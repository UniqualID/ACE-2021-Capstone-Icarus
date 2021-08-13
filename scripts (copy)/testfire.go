package main

import (
	"fmt"
	"os"
	"strconv"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

func main() {

	var query icarus.QueryPackage

	if len(os.Args) < 3 {
		fmt.Println("Missing arguments")
		os.Exit(1)
	} else if len(os.Args) == 3 {
		query = icarus.NewQuery("10.59.144.32", "179")
	} else if len(os.Args) == 4 {
		query = icarus.NewQuery(os.Args[3], "179")
	} else {
		fmt.Println("Too many arguments")
		os.Exit(1)
	}
	resp, ok := query.Authenticate("valinar", "gitlabsux;(bloodybreach")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}

	vehID, _ := strconv.Atoi(os.Args[1])
	iffID, _ := strconv.Atoi(os.Args[2])

	configs := query.ExecutePayload(vehID, 10, 1, nil, iffID)
	responseChan, _ := query.Execute()
	response := <-responseChan
	configResponse, ok := response.Get(configs)
	if !ok {
		fmt.Println("No response")
		os.Exit(1)
	}

	if configResponse.Ok {
		fmt.Printf("Vehicle %v fire.\n", vehID)
	} else {
		fmt.Println("Error during fire: \n", configResponse.Message)
	}

}
