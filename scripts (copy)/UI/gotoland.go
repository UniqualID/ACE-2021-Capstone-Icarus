package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

func main() {
	var query icarus.QueryPackage

	if len(os.Args) < 6 {
		fmt.Println("Missing arguments")
		os.Exit(1)
	} else if len(os.Args) == 6 {
		query = icarus.NewQuery("127.0.0.1", "9443")
	} else if len(os.Args) == 7 {
		query = icarus.NewQuery(os.Args[6], "9443")
	} else {
		fmt.Println("Too many arguments")
		os.Exit(1)
	}
	resp, ok := query.Authenticate("test1", "testing")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
		os.Exit(1)
	}

	var vehID = -1
	var vehName string

	// Get vehicle ID
	statSeq := query.GetAllVehicleStatus()
	responseChan, _ := query.Execute()
	response := <-responseChan
	statusResponse, _ := response.Get(statSeq)

	for _, v := range statusResponse.Vehicles {
		if os.Args[1] == v.VehicleCallsign {
			vehID = int(v.VehicleId)
			vehName = v.VehicleCallsign
			break
		}
	}
	if vehID == -1 {
		fmt.Println("Vehicle not found.")
		os.Exit(1)
	}

	navSeq := query.SetNavMode(vehID, icarus.NAVIGATION)
	responseChan, _ = query.Execute()
	response = <-responseChan
	_, ok = response.Get(navSeq)
	if !ok {
		fmt.Println("No response")
		os.Exit(1)
	}

	destLat, _ := strconv.ParseFloat(os.Args[2], 64)
	destLon, _ := strconv.ParseFloat(os.Args[3], 64)
	alt, _ := strconv.ParseFloat(os.Args[4], 32)
	destAlt := float32(alt)
	tmp, _ := strconv.ParseFloat(os.Args[5], 32)
	cruiseSpeed := float32(tmp)

	// Velocity check
	switch s := vehName; s[0] {
	case 'G':
		if cruiseSpeed > 165 {
			cruiseSpeed = 165
		}
	case 'K':
		if cruiseSpeed > 120 {
			cruiseSpeed = 120
		}
	case 'A':
		if cruiseSpeed > 90 {
			cruiseSpeed = 90
		}
	case 'T':
		if cruiseSpeed > 105 {
			cruiseSpeed = 105
		}
	case 'S':
		if cruiseSpeed > 60 {
			cruiseSpeed = 60
		}
	}

	//AddCmd(command list, command type, latitude, longitude, altitude, velocity, turn radius, linger time, transit heading)
	cmdList := icarus.AddCmd(nil, icarus.GOTO, destLat, destLon, destAlt, cruiseSpeed, 0, 0, 0)

	gotoSeq := query.Goto(vehID, cmdList)
	responseChan, _ = query.Execute()
	response = <-responseChan
	_, ok = response.Get(gotoSeq)
	if !ok {
		fmt.Println("No response")
		os.Exit(1)
	}
	fmt.Printf("Moving vehicle %s (%d) to %.4f %.0f %.4f at %.0f meters/sec then landing.\n", os.Args[1], vehID, destLat, destAlt, destLon, cruiseSpeed)

	// Land when in range
	var inRange bool = false
	for !inRange {
		statSeq := query.GetVehicleStatus(vehID)
		responseChan, _ := query.Execute()
		response := <-responseChan
		statusResponse, ok := response.Get(statSeq)
		if !ok {
			fmt.Print((time.Now()).Format("15:04:05"), " ")
			fmt.Println("No response")
			os.Exit(1)
		}
		v := statusResponse.Vehicles[0]
		latDist := math.Pow(destLat-v.Telem.Latitude, 2)
		longDist := math.Pow(destLon-v.Telem.Longitude, 2)
		distance := math.Sqrt(latDist+longDist) * 50000

		if distance > 3000 {
			time.Sleep(5 * time.Second)
		} else if distance > 1000 {
			homeCmd := icarus.AddCmd(nil, icarus.GOTO, destLat, destLon, destAlt, 60, 0, 0, 0)
			gotoSeq = query.Goto(int(vehID), homeCmd)
			responseChan, _ = query.Execute()
			response = <-responseChan
			_, ok = response.Get(gotoSeq)
			if !ok {
				fmt.Println("Error:", response)
			}
			query.ClearQueries()
			time.Sleep(1 * time.Second)
		} else if distance > 500 {
			homeCmd := icarus.AddCmd(nil, icarus.GOTO, destLat, destLon, destAlt, 30, 0, 0, 0)
			gotoSeq = query.Goto(int(vehID), homeCmd)
			responseChan, _ = query.Execute()
			response = <-responseChan
			_, ok = response.Get(gotoSeq)
			if !ok {
				fmt.Println("Error:", response)
			}
			query.ClearQueries()
			time.Sleep(1 * time.Second)
		} else if distance > 20 {
			homeCmd := icarus.AddCmd(nil, icarus.GOTO, destLat, destLon, destAlt, 10, 0, 0, 0)
			gotoSeq = query.Goto(int(vehID), homeCmd)
			responseChan, _ = query.Execute()
			response = <-responseChan
			_, ok = response.Get(gotoSeq)
			if !ok {
				fmt.Println("Error:", response)
			}
			query.ClearQueries()
			time.Sleep(100 * time.Millisecond)
		} else {
			inRange = true
		}

	}

	// Land
	landSeq := query.SetNavMode(vehID, icarus.LAND_NOW)
	responseChan, _ = query.Execute()
	response = <-responseChan
	_, ok = response.Get(landSeq)
	if !ok {
		fmt.Print((time.Now()).Format("15:04:05"), " ")
		fmt.Println("No response")
		os.Exit(1)
	}
	fmt.Print((time.Now()).Format("15:04:05"), " ")
	fmt.Printf("Landing vehicle %s (%d)\n", os.Args[1], vehID)
}
