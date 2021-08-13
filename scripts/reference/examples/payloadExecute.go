//payloadExecute.go
//This file contains an example of executing and action on a payload system
package main

import (
	"fmt"
	icarus "git.ironzone.ace/icarus/icarusClient"
)

func main() {

	//Create a new query pointed at the IcarusServer instance
	query := icarus.NewQuery("127.0.0.1", "9443")

	//Execute Daedalus systems
	//Final parameter is for a target ID. This parameter is not used for Camera, but is used for Antimatter Missiles and Seeker Missiles
	executeSeq := query.ExecutePayload(0, icarus.Camera, 1, icarus.EmptyParams(), 0)

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
	executeResponse, ok := response.Get(executeSeq)
	if !ok {
		fmt.Println("Payload execute response not found")
	}
	fmt.Println(executeResponse)

	//Uncomment this line to show the JSON response returned by the server
	//response.ShowResponse()

}
