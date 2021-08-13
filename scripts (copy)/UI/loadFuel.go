package main

import (
	"fmt"
	"os"
	"strconv"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Missing/Invalid arguments")
		os.Exit(1)
	}

	var query = icarus.NewQuery("127.0.0.1", "9443")
	resp, ok := query.Authenticate("test1", "testing")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}

	var vehID = -1
	statSeq := query.GetAllVehicleStatus()
	responseChan, _ := query.Execute()
	response := <-responseChan
	statusResponse, _ := response.Get(statSeq)

	for _, v := range statusResponse.Vehicles {
		if os.Args[1] == v.VehicleCallsign {
			vehID = int(v.VehicleId)
			break
		}
	}
	if vehID == -1 {
		fmt.Println("Vehicle not found.")
		os.Exit(1)
	}
	query.ClearQueries()

	// vehID, _ := strconv.Atoi(os.Args[1])
	fuelAmt, _ := strconv.Atoi(os.Args[2])
	configs := icarus.AddPayloadConfig(nil, "Fuel", icarus.Fuel, fuelAmt, true)
	configSeq := query.ConfigurePayloads(vehID, configs)
	responseChan, _ = query.Execute()
	response = <-responseChan
	_, ok = response.Get(configSeq)
	if !ok {
		fmt.Println("No response")
		os.Exit(1)
	}
	fmt.Printf("Fueling vehicle %d.\n", vehID)
}
