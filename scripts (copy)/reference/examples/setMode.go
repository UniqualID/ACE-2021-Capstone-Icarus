//setMode.go
//This file contains an example of setting the navigation mode of a vehicle.
package main

import (
	"fmt"
	icarus "git.ironzone.ace/icarus/icarusClient"
	"time"
)

func main() {
	//Create a new query pointed at the IcarusServer instance
	query := icarus.NewQuery("127.0.0.1", "9443")

	//Set vehicle mode to Land
	landSeq := query.SetNavMode(1, icarus.LAND_NOW)

	//Authenticate to the server
	resp, ok := query.Authenticate("test1", "testing")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}

	//Uncomment this line to show the JSON query being sent to the server
	//query.ShowQuery()

	//Execute the query and read the responses
	responseChan, _ := query.Execute()

	fmt.Println("Waiting for responses:")
	response := <-responseChan
	landResponse, ok := response.Get(landSeq)
	if !ok {
		fmt.Println("Status response not found")
	}
	fmt.Println(landResponse)

	//Uncomment this line to show the JSON response returned by the server
	//response.ShowResponse()

	//Wait for five seconds before taking off
	time.Sleep(5 * time.Second)

	//Clear queries
	query.ClearQueries()

	//Set vehicle mode to take off
	takeoffSeq := query.SetNavMode(1, icarus.TAKE_OFF)
	//After takeoff, start navigating
	query.SetNavMode(1, icarus.NAVIGATION)

	//Uncomment this line to show the JSON query being sent to the server
	//query.ShowQuery()

	//Execute the query and read the responses
	responseChan2, _ := query.Execute()

	fmt.Println("Waiting for responses:")
	response = <-responseChan2
	takeoffResponse, ok := response.Get(takeoffSeq)
	if !ok {
		fmt.Println("Status response not found")
	}
	fmt.Println(takeoffResponse)

	//Uncomment this line to show the JSON response returned by the server
	//response.ShowResponse()

}
