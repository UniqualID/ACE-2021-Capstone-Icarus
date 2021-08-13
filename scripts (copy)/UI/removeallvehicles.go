//removeVehicle.go
//This file contains an example of removing a vehicle from Icarus Server
package main

import (
	"fmt"
	"os"
	"time"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

func main() {

	var query icarus.QueryPackage

	if len(os.Args) == 1 {
		query = icarus.NewQuery("127.0.0.1", "9443")
	} else if len(os.Args) == 2 {
		query = icarus.NewQuery(os.Args[1], "9443")
	} else {
		fmt.Println("Too many arguments")
		os.Exit(1)
	}

	resp, ok := query.Authenticate("test1", "testing")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}

	statSeq := query.GetAllVehicleStatus()
	responseChan, _ := query.Execute()
	response := <-responseChan
	statusResponse, ok := response.Get(statSeq)
	if !ok {
		fmt.Println("No response")
	}

	for _, v := range statusResponse.Vehicles {
		removeSeq := query.RemoveVehicle(int(v.VehicleId))
		responseChan2, _ := query.Execute()
		response = <-responseChan2
		_, ok = response.Get(removeSeq)
		if !ok {
			fmt.Println("No response")
		}
	}
	fmt.Print((time.Now()).Format("15:04:05"), " ")
	fmt.Println("All vehicles removed.")
}
