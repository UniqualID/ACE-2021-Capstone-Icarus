package main //defacto start seen in the example START

import (
	"fmt"
	"os"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

var vehicleID int
var fuelAmount int

//var amountFuel int = 20
var query = icarus.NewQuery("10.59.144.202", "179")

func loadFuel(vehicleID int, fuelAmount int) {

	//Clear the query from the queue
	query.ClearQueries()

	//Load 100 units of Fuel

	configs := icarus.AddPayloadConfig(nil, "Fuel", icarus.Fuel, fuelAmount, true)
	//payload type Fuel = 5 icarus.FUEL payloadType = 5
	//configSeq := query.EnablePayload(vehID, icarus.Fuel, true

	configSeq := query.ConfigurePayloads(vehicleID, configs)
	responseChan, _ := query.Execute()

	fmt.Printf("Vehicle %v fueling...\n", vehicleID)
	response := <-responseChan

	configResponse, ok := response.Get(configSeq)
	if !ok {
		fmt.Println("Refuel response not found")
		os.Exit(1)
	}

	if configResponse.Ok {
		fmt.Printf("Vehicle %v fueling complete\n\n", vehicleID)
	} else {
		fmt.Println("Error during refueling:", configResponse.Message)
	}

	//Clear the refuel query from the queue
	query.ClearQueries()
}

func main() {

	// User enters vehicle ID (vehID) of drone
	fmt.Println("Enter the vehicle ID of the drone to fuel:")
	fmt.Scanln(&vehicleID)

	// User enters amount of fuel
	fmt.Println("Enter the amount of fuel:")
	fmt.Scanln(&fuelAmount)

	//Create a new query pointed at the IcarusServer instance
	resp, ok := query.Authenticate("valinar", "gitlabsux;(bloodybreach")

	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
		os.Exit(1)
	}

	//add fuel
	loadFuel(vehicleID, fuelAmount)

	//Clear the refuel query from the queue
	query.ClearQueries()

}
