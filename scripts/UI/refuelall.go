package main

import (
	"fmt"
	"os"
	"time"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

func main() {

	var query icarus.QueryPackage

	if len(os.Args) < 2 {
		fmt.Print((time.Now()).Format("15:04:05"), " ")
		fmt.Println("Missing/Invalid arguments")
		os.Exit(1)
	} else if len(os.Args) == 2 {
		query = icarus.NewQuery("10.59.144.207", "22")
	} else if len(os.Args) == 3 {
		query = icarus.NewQuery(os.Args[2], "22")
	} else {
		fmt.Println("Too many arguments")
		os.Exit(1)
	}

	resp, ok := query.Authenticate("valinar", "thisPasswordNeedsToWorkPLZ")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}

	// Get vehicle ID
	statSeq := query.GetAllVehicleStatus()
	responseChan, _ := query.Execute()
	response := <-responseChan
	statusResponse, _ := response.Get(statSeq)

	var configSeq uint32
	var configs []icarus.PayloadStatus

	for _, v := range statusResponse.Vehicles {
		if v.Available && v.VehicleType.Mode == icarus.LAND_NOW && v.PayStatus[5].Resources < 100 {

			// fmt.Printf("Loading fuel on vehicle %s (%d)\n", v.VehicleCallsign, v.VehicleId)
			configs = icarus.AddPayloadConfig(configs, "Fuel", icarus.Fuel, 100, true)
			configSeq = query.ConfigurePayloads(int(v.VehicleId), configs)
		}
	}
	responseChan, _ = query.Execute()
	response = <-responseChan
	_, ok = response.Get(configSeq)
	if !ok {
		fmt.Println("No response")
		os.Exit(1)
	}
	fmt.Print((time.Now()).Format("15:04:05"), " ")
	fmt.Println("All vehicles fueled")
}
