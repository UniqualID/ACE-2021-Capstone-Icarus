package main //defacto start seen in the example START

import (
	"fmt"
	"os"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

var vehicleID int = 4

//var amountFuel int = 20
var query = icarus.NewQuery("10.59.144.202", "179")

func loadFuel(vehID int) {

	//Clear the query from the queue
	query.ClearQueries()

	//Load 100 units of Fuel

	configs := icarus.AddPayloadConfig(nil, "Fuel", icarus.Fuel, 5, true)
	//payload type Fuel = 5 icarus.FUEL payloadType = 5
	//configSeq := query.EnablePayload(vehID, icarus.Fuel, true

	configSeq := query.ConfigurePayloads(vehID, configs)
	responseChan, _ := query.Execute()

	fmt.Printf("Vehicle %v fueling...\n", vehID)
	response := <-responseChan

	configResponse, ok := response.Get(configSeq)
	if !ok {
		fmt.Println("Refuel response not found")
		os.Exit(1)
	}

	if configResponse.Ok {
		fmt.Printf("Vehicle %v fueling complete\n\n", vehID)
	} else {
		fmt.Println("Error during refueling:", configResponse.Message)
	}

	//Clear the refuel query from the queue
	query.ClearQueries()
}

func main() {

	//Create a new query pointed at the IcarusServer instance
	resp, ok := query.Authenticate("valinar", "gitlabsux;(bloodybreach")

	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
		os.Exit(1)
	}

	//add fuel
	loadFuel(vehicleID)

	//Clear the refuel query from the queue
	query.ClearQueries()

}
