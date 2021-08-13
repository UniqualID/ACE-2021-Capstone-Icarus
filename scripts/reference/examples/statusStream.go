//statusStream.go
//This file contains an example of starting a status stream for a vehicle with ID 0.
package main

import (
	"fmt"
	icarus "git.ironzone.ace/icarus/icarusClient"
)

func main() {
	//Create a new query pointed at the IcarusServer instance
	query := icarus.NewQuery("127.0.0.1", "9443")

	//Start a status stream for vehicle (ID: 0). This gives a continuous status update to the response channel
	statSeq := query.StartStatusStream(0)

	//Authenticate to the server
	resp, ok := query.Authenticate("test1", "testing")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}
	fmt.Println("Authenticated")

	//Uncomment this line to show the JSON query being sent to the server
	//query.ShowQuery()

	//Execute the query and read the responses
	responseChan, _ := query.Execute()

	for i := 0; i < 100000; i++ {
		response := <-responseChan
		statusResponse, ok := response.Get(statSeq)
		if !ok {
			fmt.Println("Status response not found")
		}
		fmt.Println(statusResponse)
		//Uncomment this line to show the JSON response returned by the server
		//response.ShowResponse()

	}

}
