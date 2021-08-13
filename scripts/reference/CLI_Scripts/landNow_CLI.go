package main //defacto start seen in the example START

import (
	"fmt"
	"os"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

var vehicleID int

var query = icarus.NewQuery("10.59.144.202", "179")

func land_now(vehicleID int) {
	//Clear the navigation query from the queue
	query.ClearQueries()

	landSeq := query.SetNavMode(int(vehicleID), icarus.LAND_NOW)
	responseChan, _ := query.Execute()

	fmt.Printf("Vehicle %v preparing to land...\n", vehicleID)
	response := <-responseChan
	landResponse, ok := response.Get(landSeq)
	if !ok {
		fmt.Println("Land response not found")
		os.Exit(1)
	}
	if landResponse.Ok {
		fmt.Printf("Vehicle %v landed\n\n", vehicleID)
	} else {
		fmt.Println("Unable to land")
	}

	//Clear the navigation query from the queue
	query.ClearQueries()
}

func main() {

	// User enters vehicle ID (vehID) of drone
	fmt.Println("Enter the vehID of the drone to land:")
	fmt.Scanln(&vehicleID)

	//Create a new query pointed at the IcarusServer instance
	resp, ok := query.Authenticate("valinar", "gitlabsux;(bloodybreach")

	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
		os.Exit(1)
	}

	//land
	land_now(vehicleID)

}
