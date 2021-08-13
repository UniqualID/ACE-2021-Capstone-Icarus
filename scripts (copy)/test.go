package main

import (
	"fmt"
	"os"
	"strconv"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

func main() {

	var query icarus.QueryPackage

	if len(os.Args) < 2 {
		fmt.Println("Missing arguments")
		os.Exit(1)
	} else if len(os.Args) == 2 {
		query = icarus.NewQuery("10.59.144.32", "179")
	} else if len(os.Args) == 3 {
		query = icarus.NewQuery(os.Args[2], "179")
	} else {
		fmt.Println("Too many arguments")
		os.Exit(1)
	}
	resp, ok := query.Authenticate("valinar", "gitlabsux;(bloodybreach")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}

	vehID, _ := strconv.Atoi(os.Args[1])

	// configs := query.ExecutePayload(vehID, 3, 1, nil, 2111)
	// responseChan, _ := query.Execute()
	// response := <-responseChan
	// configResponse, ok := response.Get(configs)
	// if !ok {
	// 	fmt.Println("No response")
	// 	os.Exit(1)
	// }

	// if configResponse.Ok {
	// 	fmt.Printf("Vehicle %v fire.\n", vehID)
	// } else {
	// 	fmt.Println("Error during fire: \n", configResponse.Message)
	// }

	statSeq := query.GetVehicleStatus(vehID)
	responseChan, _ := query.Execute()
	response := <-responseChan
	statusResponse, ok := response.Get(statSeq)
	if !ok {
		fmt.Println("No response")
		os.Exit(1)
	}

	v := statusResponse.Vehicles[0]

	fmt.Println("Enabled:", v.PayStatus[3].Enabled)
	fmt.Println("Name:", v.PayStatus[3].Name)
	fmt.Println("ID:", v.PayStatus[3].Id)
	fmt.Println("Parameters:", v.PayStatus[3].Parameters)
	fmt.Println("Resources:", v.PayStatus[3].Resources)

	fmt.Println("Payload All:", v.PayStatus[0].Resources)
	fmt.Println("Thermal Lance:", v.PayStatus[3].Resources)
	fmt.Println("Camera:", v.PayStatus[4].Resources)
	fmt.Println("Fuel:", v.PayStatus[5].Resources)
	fmt.Println("Phosphex:", v.PayStatus[7].Resources)
	fmt.Println("PhosphexRemediation:", v.PayStatus[8].Resources)
	fmt.Println("AirRadar:", v.PayStatus[9].Resources)
	fmt.Println("AntiMatterMissile:", v.PayStatus[10].Resources)
	fmt.Println("AllRadar:", v.PayStatus[11].Resources)
	fmt.Println("GroundRadar:", v.PayStatus[12].Resources)
	fmt.Println("SAM:", v.PayStatus[13].Resources)
	fmt.Println("Cargo:", v.PayStatus[14].Resources)
	fmt.Println("SeekerMissile:", v.PayStatus[15].Resources)
}
