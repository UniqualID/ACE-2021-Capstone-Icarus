//removeVehicle.go
//This file contains an example of removing a vehicle from Icarus Server
package main

import (
	"fmt"
	icarus "git.ironzone.ace/icarus/icarusClient"
	"time"
)

func main() {
	//Create a new query pointed at the IcarusServer instance
	query := icarus.NewQuery("127.0.0.1", "9443")

	//Add new vehicle to IcarusServer
	addSeq := query.AddNewVehicle("127.0.0.1", "5001", "Test Vehicle", "44444", 0, make([]string, 0), make([]byte, 0), make([]byte, 0), 0, icarus.DefaultC3poTime, icarus.DefaultDaedalusTime)

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
	fmt.Println(response.Get(addSeq))

	//Uncomment this line to show the JSON response returned by the server
	//response.ShowResponse()

	//Wait 5 seconds before removing vehicle
	time.Sleep(5 * time.Second)

	//Clear current queries
	query.ClearQueries()

	//Remove vehicle (ID: 0) from IcarusServer
	removeSeq := query.RemoveVehicle(0)

	//Uncomment this line to show the JSON query being sent to the server
	//query.ShowQuery()

	//Execute the query and read the responses
	responseChan2, _ := query.Execute()

	fmt.Println("Waiting for responses:")
	response = <-responseChan2
	fmt.Println(response.Get(removeSeq))

	//Uncomment this line to show the JSON response returned by the server
	//response.ShowResponse()

}
