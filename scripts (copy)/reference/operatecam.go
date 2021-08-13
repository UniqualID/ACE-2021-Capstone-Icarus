package main //defacto start seen in the example START

import (
	"fmt"
	"os"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

var vehicleID int = 0
var query = icarus.NewQuery("10.59.144.202", "179")

func operatecam(vehID int) {

	//Clear the query from the queue
	query.ClearQueries()

	//payload type 4 = camera
	// var payloadType int = 4

	configs := query.ExecutePayload(vehID, 4, 1, nil, 0)

	//payload type Fuel = 5 icarus.FUEL payloadType = 5
	//configSeq := query.EnablePayload(vehID, icarus.Fuel, true

	// configSeq := query.ConfigurePayloads(vehID, configs) // I AM NOT SURE WHAT THIS LINE DOES BUT I THINK IT RETURNS THE SAME THING AS EXECUTE PAYLOAD
	responseChan, _ := query.Execute()

	fmt.Printf("Vehicle %v taking picture...\n", vehID)
	response := <-responseChan

	configResponse, ok := response.Get(configs) // changed from configSeq
	if !ok {
		fmt.Println("picture not taken")
		os.Exit(1)
	}

	var images []string = configResponse.PayloadResponse.File

	if configResponse.Ok {
		fmt.Printf("Vehicle %v picture complete\n\n", vehID)
		fmt.Printf("Images returned by vehicle %v: %v \n", vehID, images)
	} else {
		fmt.Println("Error during picture taken: \n", configResponse.Message)
	}

	for i, s := range images {
		fmt.Println(i, s)
	}

	//Clear the take picture query from the queue
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
	operatecam(vehicleID)

	//Clear the refuel query from the queue
	query.ClearQueries()

}
