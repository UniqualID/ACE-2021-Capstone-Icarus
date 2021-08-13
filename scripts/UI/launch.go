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
		fmt.Println("Missing arguments")
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

	navSeq := query.SetNavMode(vehID, icarus.NAVIGATION) // Drones can go straight into navigation mode; may be issue in future
	responseChan, _ = query.Execute()
	response = <-responseChan
	_, ok = response.Get(navSeq)
	if !ok {
		fmt.Println("No response")
	}
	fmt.Print((time.Now()).Format("15:04:05"), " ")
	fmt.Printf("Launching vehicle %s (%d)\n", os.Args[1], vehID)
}
