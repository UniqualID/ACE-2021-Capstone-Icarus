package main

import (
	"fmt"
	"math"
	"os"
	"time"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

type Base struct {
	lat  float64
	lon  float64
	name string
}

var bases = []Base{
	{
		lat:  50.5572,
		lon:  -63.4119,
		name: "SHELOB",
	},
	{
		lat:  50.5580,
		lon:  -61.0745,
		name: "SAURON",
	},
	{
		lat:  50.2306,
		lon:  -62.1124,
		name: "SCATHA",
	},
	{
		lat:  50.8405,
		lon:  -62.4290,
		name: "SMAUG",
	},
	{
		lat:  50.2223,
		lon:  -66.2656,
		name: "EDHELLOND",
	},
	{
		lat:  50.6533,
		lon:  -68.7170,
		name: "ELVENKING",
	},
	{
		lat:  49.1819,
		lon:  -68.3639,
		name: "EDORAS",
	},
	{
		lat:  49.6165,
		lon:  -67.7075,
		name: "ETHRING",
	},
	{
		lat:  50.3335,
		lon:  -60.6990,
		name: "EGLAREST",
	},
	{
		lat:  49.0181,
		lon:  -60.1556,
		name: "ANGBAND",
	},
	{
		lat:  49.0841,
		lon:  -60.9777,
		name: "RHUN",
	},
	{
		lat:  50.0000,
		lon:  -65.7590,
		name: "SAURUMAN",
	},
}

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
		fmt.Println("Unable to aunteticate to IcarusServer:", resp)
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

	// query the target vehicle for location, speed, and altitude
	statSeq = query.GetVehicleStatus(vehID)
	responseChan, _ = query.Execute()
	response = <-responseChan
	statusResponse, ok = response.Get(statSeq)
	v := statusResponse.Vehicles[0]
	if !ok {
		fmt.Println("No response")
	}

	// find nearest base
	vlat := v.Telem.Latitude
	vlon := v.Telem.Longitude
	var stDist float64 = 100
	var stBase Base
	for _, bs := range bases {
		currDist := math.Sqrt(math.Pow(vlat-bs.lat, 2) + math.Pow(vlon-bs.lon, 2))
		if currDist < stDist {
			stDist = currDist
			stBase = bs
		}
	}

	// set speed
	var vehName string
	var cruiseSpeed float32
	vehName = v.VehicleCallsign

	switch s := vehName; s[0] {
	case 'G':
		cruiseSpeed = 165
	case 'K':
		cruiseSpeed = 120
	case 'A':
		cruiseSpeed = 90
	case 'T':
		cruiseSpeed = 105
	case 'S':
		cruiseSpeed = 60
	}

	// set other params
	destLat := stBase.lat
	destLon := stBase.lon
	destAlt := v.Telem.Altitude

	// set drone to go to base

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
	fmt.Printf("Moving vehicle %s (%d) to %.4f %.0f %.4f (%s) at %.0f meters/sec then landing.\n", os.Args[1], vehID, destLat, destAlt, destLon, stBase.name, cruiseSpeed)

	// Land drone
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
	fmt.Printf("Landing vehicle %d.\n", vehID)

}
