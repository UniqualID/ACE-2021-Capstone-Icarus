package main //defacto start seen in the example START

import (
	"fmt"
	"os"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

var query = icarus.NewQuery("10.59.144.202", "179")

//5 waypoints (these are correct according to current sortie plan - 10 Jun 21)
const landLocationLat = 50.3335
const landLocationLon = -60.6999

func nav(vehID int) {

	//Clear the take off query from the queue
	query.ClearQueries()

	//Change mode to Navigate
	navSeq := query.SetNavMode(vehID, icarus.NAVIGATION)
	responseChan, _ := query.Execute()
	//CHECK if this prints correctly
	fmt.Printf("Vehicle %v entering navigation mode...\n", vehID)
	response := <-responseChan
	navResponse, ok := response.Get(navSeq)
	if !ok {
		fmt.Println("NAVIGATE response not found")
		os.Exit(1)
	}

	if navResponse.Ok {
		fmt.Printf("Vehicle %v ready to navigate...\n", vehID)
	} else {
		fmt.Println("Error during mode change:", navResponse.Message)
	}
	//Clear the mode change query from the queue
	query.ClearQueries()

	//We haven't decided what the altitude should be for each waypoint yet, but this was the one used in the example
	var cruiseAlt float32 = 1000.0
	//var cruiseSpeed float32 = 120.0

	//AddCmd(command list, command type, latitude, longitude, altitude, velocity, turn radius, linger time, transit heading)
	cmdList := icarus.AddCmd(nil, icarus.GOTO, landLocationLat, landLocationLon, cruiseAlt/5, 10, 0, 0, 0)
	//Iterate this process???

	//collate entire list of commands for vehicle 0 and declare it as gotoSeq
	gotoSeq := query.Goto(vehID, cmdList)
	responseChan, _ = query.Execute()
	response = <-responseChan
	gotoResponse, ok := response.Get(gotoSeq)
	if !ok {
		fmt.Println("Go to response not found")
		os.Exit(1)
	}

	if gotoResponse.Ok {
		fmt.Printf("Vehicle %v Navigating to waypoints\n\n", vehID)
	} else {
		fmt.Println("Error during navigation:", gotoResponse.Message)
	}
	//Clear the navigation query from the queue
	query.ClearQueries()

}

func main() {

	//Create a new query pointed at the IcarusServer instance
	resp, ok := query.Authenticate("valinar", "gitlabsux;(bloodybreach")

	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
		os.Exit(1)
	}

	//Clear the refuel query from the queue
	query.ClearQueries()

	//give list of waypoints to drone, instruct to go to all of them
	nav(4)

}
