//this file contains nav information
//We took the file "isr_example.go" from the example missions directory and cut out the parts relevant to NAVIGATION
//We modified the code to fit our sortie
//on battlespace VPN type godoc -http localhost:8080 in terminal then visit http://localhost:8080 in web browser

package main //defacto start seen in the example START

import (
	"fmt"
	"os"
	"time"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

var numDrones = 4

//CHANGE
var vehiclesIDs = [4]int{4, 5, 6, 7}

var query = icarus.NewQuery("10.59.144.202", "179")

//5 waypoints (these are correct according to current sortie plan - 10 Jun 21)
const firstLocationLat = 49.841
const firstLocationLon = -65.9895
const secondLocationLat = 50.1927
const secondLocationLon = -64.6313
const thirdLocationLat = 50.0105
const thirdLocationLon = -62.7558
const fourthLocationLat = 49.7292
const fourthLocationLon = -61.9451
const fifthLocationLat = 50.3228
const fifthLocationLon = -60.7200
const sixthLocationLat = 50.3292
const sixthLocationLon = -60.7079
const landLocationLat = 50.3335
const landLocationLon = -60.6999

func takeOff(vehID int) {
	//Clear the current query from the queue
	query.ClearQueries()

	//Change mode to Take off
	takeOffSeq := query.SetNavMode(vehID, icarus.TAKE_OFF)
	responseChan, _ := query.Execute()
	fmt.Printf("Vehicle %v taking off...\n", vehID)

	response := <-responseChan
	takeOffResponse, ok := response.Get(takeOffSeq)

	if !ok {
		fmt.Println("TAKE_OFF response not found")
		os.Exit(1)
	}

	if takeOffResponse.Ok {
		fmt.Printf("Vehicle %v take off complete\n\n", vehID)
	} else {
		fmt.Println("Error during take off:", takeOffResponse.Message)
	}

	//Clear the current query from the queue
	query.ClearQueries()

	//Sleep to ensure UAV is in air before navigating
	time.Sleep(5 * time.Second)
}

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
	var cruiseSpeed float32 = 120.0

	//AddCmd(command list, command type, latitude, longitude, altitude, velocity, turn radius, linger time, transit heading)
	cmdList := icarus.AddCmd(nil, icarus.GOTO, firstLocationLat, firstLocationLon, cruiseAlt, cruiseSpeed, 0, 0, 0)
	cmdList = icarus.AddCmd(cmdList, icarus.GOTO, secondLocationLat, secondLocationLon, cruiseAlt, cruiseSpeed, 0, 0, 0)
	cmdList = icarus.AddCmd(cmdList, icarus.GOTO, thirdLocationLat, thirdLocationLon, cruiseAlt, cruiseSpeed, 0, 0, 0)
	cmdList = icarus.AddCmd(cmdList, icarus.GOTO, fourthLocationLat, fourthLocationLon, cruiseAlt, cruiseSpeed, 0, 0, 0)
	cmdList = icarus.AddCmd(cmdList, icarus.GOTO, fifthLocationLat, fifthLocationLon, cruiseAlt, cruiseSpeed, 0, 0, 0)
	cmdList = icarus.AddCmd(cmdList, icarus.GOTO, sixthLocationLat, fifthLocationLon, 100, 40, 0, 0, 0)
	cmdList = icarus.AddCmd(cmdList, icarus.GOTO, landLocationLat, landLocationLon, 50, 10, 0, 0, 0)
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

	for i := 0; i < numDrones; i++ {

		//take off
		takeOff(vehiclesIDs[i])

		//give list of waypoints to drone, instruct to go to all of them
		nav(vehiclesIDs[i])

		//wait before next drone takes off
		time.Sleep(120 * time.Second)
	}

}
