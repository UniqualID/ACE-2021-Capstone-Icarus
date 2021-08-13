package main

import (
	"fmt"
	"os"
	"strconv"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

func main() {

	var query icarus.QueryPackage

	if len(os.Args) < 4 {
		fmt.Println("Missing arguments")
		os.Exit(1)
	} else if len(os.Args) == 4 {
		query = icarus.NewQuery("127.0.0.1", "9443")
	} else if len(os.Args) == 5 {
		query = icarus.NewQuery(os.Args[4], "9443")
	} else {
		fmt.Println("Too many arguments")
		os.Exit(1)
	}

	resp, ok := query.Authenticate("test1", "testing")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}

	var vehID = -1

	// Get vehicle ID
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

	tmp, _ := strconv.Atoi(os.Args[2])
	payID := icarus.PayloadType(tmp)
	amt, _ := strconv.Atoi(os.Args[3])

	configs := icarus.AddPayloadConfig(nil, payID.String(), payID, amt, true)
	configSeq := query.ConfigurePayloads(vehID, configs)
	responseChan, _ = query.Execute()
	response = <-responseChan
	_, ok = response.Get(configSeq)
	if !ok {
		fmt.Println("No response")
		os.Exit(1)
	}
	fmt.Printf("Loading payload %s on vehicle %s (%d)\n", payID.String(), os.Args[1], vehID)
}
