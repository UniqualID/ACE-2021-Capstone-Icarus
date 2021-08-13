//payloadStatus.go
//This file contains an example of retrieving payload status of a payload system
package main

import (
	"fmt"
	icarus "git.ironzone.ace/icarus/icarusClient"
)

func main() {
	//Create a new query pointed at the IcarusServer instance
	query := icarus.NewQuery("127.0.0.1", "9443")

	//Get status of all payloads on vehicle 0
	statusSeq := query.StatusPayload(1, icarus.AllPayloads)

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
	statusResponse, ok := response.Get(statusSeq)
	if !ok {
		fmt.Println("Status response not found")
	}
	fmt.Println(statusResponse)

	//Uncomment this line to show the JSON response returned by the server
	//response.ShowResponse()

}
