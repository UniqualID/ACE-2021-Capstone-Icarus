//goto.go
//This file contains an example of controlling the waypoint list on a vehicle.
package main

import (
	"fmt"
	icarus "git.ironzone.ace/icarus/icarusClient"
)

func main() {
	//Create a new query pointed at the IcarusServer instance
	query := icarus.NewQuery("127.0.0.1", "9443")

	//Create the command list to upload to the vehicle

	cmdList := icarus.AddCmd(nil, icarus.GOTO, 43.0343799, -75.6469088, 25, 10, 0, 0, 3.14)
	cmdList = icarus.AddCmd(cmdList, icarus.GOTO, 43.0340694, -75.6487492, 35, 20, 0, 0, 4.7)

	//Set the waypoint list on vehicle 0
	gotoSeq := query.Goto(0, cmdList)

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
	gotoResponse, ok := response.Get(gotoSeq)
	if !ok {
		fmt.Println("Goto response not found")
	}
	fmt.Println(gotoResponse)

	//Uncomment this line to show the JSON response returned by the server
	//response.ShowResponse()

}
