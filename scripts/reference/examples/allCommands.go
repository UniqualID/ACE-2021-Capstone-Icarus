//allCommands.go
//This file contains an example of all available commands. This shows how multiple commands can be sent in a single query.
package main

import (
	"fmt"
	icarus "git.ironzone.ace/icarus/icarusClient"
	"time"
)

//This file shows the usage of all available icarus commands
//The produced JSON query of all commands can be seen in the allCommandExample.json file
func main() {
	//Create a new query pointed at the IcarusServer instance
	query := icarus.NewQuery("127.0.0.1", "9443")

	//Add new vehicle to IcarusServer
	addSeq := query.AddNewVehicle("127.0.0.1", "5001", "Test Vehicle", "44444", 0, nil, nil, nil, 0, icarus.DefaultC3poTime, icarus.DefaultDaedalusTime)

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

	//Wait 2 seconds before continueing
	time.Sleep(2 * time.Second)

	//Clear current queries
	query.ClearQueries()

	//Add new vehicle
	addSeq = query.AddNewVehicle("127.0.0.1", "5001", "Test Vehicle 2", "44444", 0, nil, nil, nil, 0, icarus.DefaultC3poTime, icarus.DefaultDaedalusTime)

	//Get status of first vehicle (ID: 0)
	query.GetVehicleStatus(0)

	//Get status of all vehicles
	query.GetAllVehicleStatus()

	//Start status stream for vehicle (ID: 0)
	query.StartStatusStream(0)

	//Stop the status stream
	query.StopStatusStream(0)

	//Get waypoints from vehicle (ID: 0)
	query.GetWaypointList(0)

	//Remove vehicle (ID: 1) from IcarusServer
	query.RemoveVehicle(1)

	//Create the command list to upload to the vehicle
	cmdList := icarus.AddCmd(nil, icarus.GOTO, 43.0343799, -75.6469088, 25, 10, 0, 0, 3.14)
	cmdList = icarus.AddCmd(cmdList, icarus.GOTO, 43.0340694, -75.6487492, 35, 20, 0, 0, 4.7)

	//Set the waypoint list on vehicle 0
	query.Goto(0, cmdList)

	//Set navigation mode to NAVIGATION
	query.SetNavMode(0, icarus.NAVIGATION)

	//Uncomment this line to show the JSON query being sent to the server
	//query.ShowQuery()

	//Execute the query and read the responses
	responseChan2, _ := query.Execute()

	fmt.Println("Waiting for responses:")
	response = <-responseChan2
	fmt.Println(response)

	//Uncomment this line to show the JSON response returned by the server
	//response.ShowResponse()

}
