package main //defacto start seen in the example START

import (
	"fmt"
	"os"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

var vehicleID uint32 = 4
var query = icarus.NewQuery("10.59.144.202", "179")

func land_now(vehID uint32) {
	//Clear the navigation query from the queue
	query.ClearQueries()

	landSeq := query.SetNavMode(int(vehID), icarus.LAND_NOW)
	responseChan, _ := query.Execute()

	fmt.Printf("Vehicle %v preparing to land...\n", vehID)
	response := <-responseChan
	landResponse, ok := response.Get(landSeq)
	if !ok {
		fmt.Println("Land response not found")
		os.Exit(1)
	}
	if landResponse.Ok {
		fmt.Printf("Vehicle %v landed\n\n", vehID)
	} else {
		fmt.Println("Unable to land")
	}

	//Clear the navigation query from the queue
	query.ClearQueries()
}

func main() {

	//Create a new query pointed at the IcarusServer instance
	resp, ok := query.Authenticate("valinar", "gitlabsux;(bloodybreach")

	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
		os.Exit(1)
	}

	//land
	land_now(vehicleID)

	//Clear the refuel query from the queue
	query.ClearQueries()

}

// func atFinalWaypoint(vehID uint32) {

// 	//Wait until at initial base then land
// 	//statusSeq := query.StartStatusStream(int(vehID))
// 	responseChan, stopChan := query.Execute()
// 	for {
// 		response := <-responseChan
// 		statusResponse, ok := response.Get(vehID)
// 		if !ok {
// 			continue
// 		}
// 		vehicles := statusResponse.Vehicles
// 		if len(vehicles) < 1 {
// 			continue
// 		}
// 		telem := vehicles[0].Telem
// 		fmt.Printf("Vehicle %v Lat: %v, Lon: %v\n", vehID, telem.Latitude, telem.Longitude)
// 		if floatCompare(telem.Latitude, landLocationLat, 0.01) && floatCompare(telem.Longitude, landLocationLon, 0.01) && floatCompare(float64(telem.Altitude), 1000, 500) {
// 			stopChan <- true
// 			break
// 		}

// 	}
// 	fmt.Printf("Vehicle %v at final waypoint...\n", vehID)

// 	//Clear the navigation query from the queue
// 	query.ClearQueries()
// }
