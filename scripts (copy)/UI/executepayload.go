package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

func main() {

	var query icarus.QueryPackage

	if len(os.Args) < 5 {
		fmt.Println("Missing arguments")
		os.Exit(1)
	} else if len(os.Args) == 5 {
		query = icarus.NewQuery("127.0.0.1", "9443")
	} else if len(os.Args) == 6 {
		query = icarus.NewQuery(os.Args[5], "9443")
	} else {
		fmt.Println("Too many arguments")
		os.Exit(1)
	}

	resp, ok := query.Authenticate("test1", "testing")
	if !ok {
		fmt.Print((time.Now()).Format("15:04:05"), " ")
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}

	payID, _ := strconv.Atoi(os.Args[2])
	action, _ := strconv.Atoi(os.Args[3])
	target, _ := strconv.Atoi(os.Args[4])

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

	// /*
	// 	For a fighter that executes an air-to-air missile
	// 	Check that the target vehicle is in range
	// */
	// if icarus.PayloadType(payID) == 10 {
	// 	// i am a fighter

	// 	// find target location
	// 	statSeq := query.GetVehicleStatus(vehID)
	// 	responseChan, _ := query.Execute()
	// 	response := <-responseChan
	// 	statusResponse, ok := response.Get(statSeq)
	// 	v := statusResponse.Vehicles[0]

	// 	if !ok {
	// 		fmt.Print((time.Now()).Format("15:04:05"), " ")
	// 		fmt.Println("No response")
	// 		os.Exit(1)
	// 	}

	// 	var radarInfo icarus.PayloadStatus
	// 	if v.PayStatus[9].Resources > 0 {
	// 		radarInfo = v.PayStatus[icarus.AirRadar]
	// 	} else if v.PayStatus[12].Resources > 0 {
	// 		radarInfo = v.PayStatus[icarus.GroundRadar]
	// 	} else if v.PayStatus[11].Resources > 0 {
	// 		radarInfo = v.PayStatus[icarus.AllRadar]
	// 	}

	// 	// find our location

	// 	var targetLat float64
	// 	var targetLong float64
	// 	var targetAlt float64
	// 	if len(radarInfo.Radar) > 0 {
	// 		for asset, radarPing := range radarInfo.Radar {
	// 			if asset == int32(target) {
	// 				targetLat = radarPing.Latitude
	// 				targetLong = radarPing.Longitude
	// 				targetAlt = float64(radarPing.Altitude)
	// 			}
	// 		}
	// 	}

	// 	// calculate our separation from the target

	// 	ourLat := v.Telem.Latitude
	// 	ourLong := v.Telem.Longitude
	// 	ourAlt := v.Telem.Altitude

	// 	latSep := math.Abs(ourLat - targetLat)
	// 	longSep := math.Abs(ourLong - targetLong)
	// 	altSep := math.Abs(float64(ourAlt) - targetAlt)

	// 	totalSep := math.Sqrt(latSep*latSep + longSep*longSep + altSep*altSep)

	// 	// check if the separation is within our allowed range
	// 	const MAXAIRMISSILERANGE = 7500

	// 	if totalSep > MAXAIRMISSILERANGE {
	// 		fmt.Print((time.Now()).Format("15:04:05"), " ")
	// 		fmt.Printf("Target out of air-to-air missile range")

	// 	} else {
	// 		// execute payload
	// 		configs := query.ExecutePayload(vehID, icarus.PayloadType(payID), action, nil, target)
	// 		responseChan, _ := query.Execute()
	// 		response := <-responseChan
	// 		_, ok = response.Get(configs)
	// 		if !ok {
	// 			fmt.Print((time.Now()).Format("15:04:05"), " ")
	// 			fmt.Println("No response")
	// 			os.Exit(1)
	// 		}

	// 		fmt.Print((time.Now()).Format("15:04:05"), " ")
	// 		fmt.Printf("Varicle %v fired.\n", vehID)

	// 	}

	// } else {
	// 	// non-air to air payload
	configs := query.ExecutePayload(vehID, icarus.PayloadType(payID), action, nil, target)
	responseChan, _ = query.Execute()
	response = <-responseChan
	_, ok = response.Get(configs)
	if !ok {
		fmt.Print((time.Now()).Format("15:04:05"), " ")
		fmt.Println("No response")
		os.Exit(1)
	}

	fmt.Print((time.Now()).Format("15:04:05"), " ")
	if action == 0 {
		fmt.Printf("Vehicle %s (%d) jettisoned payload %s\n", os.Args[1], vehID, icarus.PayloadType(payID).String())
	} else if action == 1 {
		fmt.Printf("Vehicle %s (%d) fired payload %s\n", os.Args[1], vehID, icarus.PayloadType(payID).String())
	} else {
		fmt.Println("Incorrect action")
	}
}
