package main

import (
	"fmt"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

var vehicleID int

func main() {

	// User enters vehicle ID (vehID) of drone
	fmt.Println("Enter the vehicle ID of the drone to launch:")
	fmt.Scanln(&vehicleID)

	var query = icarus.NewQuery("10.59.144.202", "179")
	resp, ok := query.Authenticate("valinar", "gitlabsux;(bloodybreach")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}

	takeoffSeq := query.SetNavMode(vehicleID, icarus.TAKE_OFF)
	query.SetNavMode(vehicleID, icarus.NAVIGATION)
	responseChan2, _ := query.Execute()

	response := <-responseChan2
	_, ok = response.Get(takeoffSeq)
	if !ok {
		fmt.Println("No response")
	}
	fmt.Printf("Launching vehicle %d.\n", vehicleID)
}
