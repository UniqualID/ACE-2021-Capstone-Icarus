//This file contains an example ISR mission
package main

import (
	"fmt"
	icarus "git.ironzone.ace/icarus/icarusClient"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type RadarPing struct {
	VehicleID           int
	VehicleType         int
	Latitude, Longitude float64
	Altitude, Heading   float32
}

func main() {
	//Create a new query pointed at the IcarusServer instance
	query := icarus.NewQuery("127.0.0.1", "9443")

	resp, ok := query.Authenticate("test1", "testing")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
		os.Exit(1)
	}

	//This mission assumes the vehicle is waiting for takeoff at an appropriate airbase
	statusSeq := query.GetVehicleStatus(0)
	responseChan, _ := query.Execute()
	response := <-responseChan
	statusResponse, ok := response.Get(statusSeq)
	if !ok {
		fmt.Println("Unable to get initial telemetry")
		os.Exit(1)
	}
	vehicles := statusResponse.Vehicles
	if len(vehicles) < 1 {
		fmt.Println("Unable to get initial telemetry")
		os.Exit(1)
	}
	telem := vehicles[0].Telem
	initialLat := telem.Latitude
	initialLon := telem.Longitude
	//Clear the refuel query from the queue
	query.ClearQueries()

	//Load 50 units of Fuel
	configs := icarus.AddPayloadConfig(nil, "Fuel", icarus.Fuel, 50, true)

	//Assumes ISR vehicle has vehicle ID 0
	configSeq := query.ConfigurePayloads(0, configs)
	//Need to create a new response channel for every execute. The Execute function closes the previous channel
	responseChan, _ = query.Execute()
	fmt.Println("Loading fuel...")
	response = <-responseChan

	configResponse, ok := response.Get(configSeq)
	if !ok {
		fmt.Println("Refuel response not found")
		os.Exit(1)
	}

	if configResponse.Ok {
		fmt.Println("Refueling complete")
	} else {
		fmt.Println("Error during refueling:", configResponse.Message)
	}

	//Clear the refuel query from the queue
	query.ClearQueries()

	//Change mode to Take off
	takeOffSeq := query.SetNavMode(0, icarus.TAKE_OFF)
	responseChan, _ = query.Execute()
	fmt.Println("Taking off...")
	response = <-responseChan
	takeOffResponse, ok := response.Get(takeOffSeq)
	if !ok {
		fmt.Println("TAKE_OFF response not found")
		os.Exit(1)
	}

	if takeOffResponse.Ok {
		fmt.Println("UAV launch complete")
	} else {
		fmt.Println("Error during take off:", takeOffResponse.Message)
	}
	//Clear the take off query from the queue
	query.ClearQueries()

	//Sleep to ensure UAV is in air before navigating
	time.Sleep(5 * time.Second)

	//Change mode to Navigate
	navSeq := query.SetNavMode(0, icarus.NAVIGATION)
	responseChan, _ = query.Execute()
	fmt.Println("Entering Navigate mode...")
	response = <-responseChan
	navResponse, ok := response.Get(navSeq)
	if !ok {
		fmt.Println("NAVIGATE response not found")
		os.Exit(1)
	}

	if navResponse.Ok {
		fmt.Println("UAV ready to navigate")
	} else {
		fmt.Println("Error during mode change:", navResponse.Message)
	}
	//Clear the mode change query from the queue
	query.ClearQueries()

	//Navigate to point of interest and linger for 60 seconds then navigate to secondary location
	firstLocationLat := 43.0344
	firstLocationLon := -75.6543
	secondLocationLat := 43.0354
	secondLocationLon := -75.6573
	var secondLocationAlt float32 = 30.0
	cmdList := icarus.AddCmd(nil, icarus.LOITER, firstLocationLat, firstLocationLon, 25, 1, 0, 60, 0)
	cmdList = icarus.AddCmd(cmdList, icarus.GOTO, secondLocationLat, secondLocationLon, secondLocationAlt, 1, 0, 0, 0)
	gotoSeq := query.Goto(0, cmdList)
	responseChan, _ = query.Execute()
	fmt.Println("Navigating to waypoint...")
	response = <-responseChan
	gotoResponse, ok := response.Get(gotoSeq)
	if !ok {
		fmt.Println("Go to response not found")
		os.Exit(1)
	}

	if gotoResponse.Ok {
		fmt.Println("Navigating to waypoint")
	} else {
		fmt.Println("Error during navigation:", gotoResponse.Message)
	}
	//Clear the navigation query from the queue
	query.ClearQueries()

	//Wait until at final location and stay at second location until done taking picture
	statusSeq = query.StartStatusStream(0)
	responseChan, stopChan := query.Execute()
	for {
		response = <-responseChan
		//All status stream responses are sequence number 0
		statusResponse, ok := response.Get(0)
		if !ok {
			fmt.Println("No status received")
			continue
		}
		vehicles := statusResponse.Vehicles
		if len(vehicles) < 1 {
			fmt.Println("No vehicles in status")
			continue
		}
		telem := vehicles[0].Telem
		if floatCompare(telem.Latitude, secondLocationLat, 0.0001) && floatCompare(telem.Longitude, secondLocationLon, 0.0001) && floatCompare(float64(telem.Altitude), float64(secondLocationAlt), 0.01) {
			stopChan <- true
			break
		}
		//Check radar information
		radarInfo := vehicles[0].PayStatus[icarus.AirRadar]
		//Radar's current range is found in radarInfo.Resources parameter
		if len(radarInfo.Radar) > 0 {
			for vehicleID, radarPing := range radarInfo.Radar {

				//Depreciated with introduction of radar parameter
				//radarPing := parseRadar(int(vehicleID), info)

				//ID 0 means nothing was returned by the radar payload
				if vehicleID != 0 {
					//Printing here for usability, but could save for something else or respond with Air-to-Air munitions if this was a fighter
					fmt.Printf("Radar found vehicle(%d) at (%f,%f,%f)\n", vehicleID, radarPing.Latitude, radarPing.Longitude, radarPing.Altitude)
				}
			}
		}

	}
	//Clear the navigation query from the queue
	query.ClearQueries()

	cmdList = icarus.AddCmd(nil, icarus.GOTO, secondLocationLat, secondLocationLon, secondLocationAlt, 1, 0, 0, 0)

	gotoSeq = query.Goto(0, cmdList)
	responseChan, _ = query.Execute()
	fmt.Println("Stopping at second location...")
	response = <-responseChan
	gotoResponse, ok = response.Get(gotoSeq)
	if !ok {
		fmt.Println("Go to response not found")
		os.Exit(1)
	}

	if gotoResponse.Ok {
		fmt.Println("Stopping at second location")
	} else {
		fmt.Println("Error during navigation:", gotoResponse.Message)
	}
	//Clear the navigation query from the queue
	query.ClearQueries()

	//Execute a payload (take picture here, but could be any payload)
	//Final parameter is for a target ID. This parameter is not used for Camera, but is used for Antimatter Missiles and Seeker Missiles
	executeSeq := query.ExecutePayload(0, icarus.Camera, 1, icarus.EmptyParams(), 0)
	responseChan, _ = query.Execute()
	fmt.Println("Taking picture...")
	response = <-responseChan
	executeResponse, ok := response.Get(executeSeq)
	if !ok {
		fmt.Println("Payload execute response not found")
		os.Exit(1)
	}
	payloadResponse := executeResponse.PayloadResponse
	if payloadResponse.ErrorCode != 0 {
		fmt.Println("Error taking picture:", payloadResponse.Error)
	}
	if len(payloadResponse.File) > 0 {
		fmt.Println("Picture taken. Saved on server as")
		fmt.Println(payloadResponse.File[0])
	}
	//Clear the navigation query from the queue
	query.ClearQueries()

	//Return to home base
	cmdList = icarus.AddCmd(nil, icarus.GOTO, initialLat, initialLon, 20, 1, 0, 0, 0)

	gotoSeq = query.Goto(0, cmdList)
	responseChan, _ = query.Execute()
	fmt.Println("Returning to base...")
	response = <-responseChan
	gotoResponse, ok = response.Get(gotoSeq)
	if !ok {
		fmt.Println("Go to response not found")
		os.Exit(1)
	}

	if gotoResponse.Ok {
		fmt.Println("Navigating to base")
	} else {
		fmt.Println("Error during navigation:", gotoResponse.Message)
	}
	//Clear the navigation query from the queue
	query.ClearQueries()

	//Wait until at initial base then land
	statusSeq = query.StartStatusStream(0)
	responseChan, stopChan = query.Execute()
	for {
		response = <-responseChan
		statusResponse, ok := response.Get(0)
		if !ok {
			continue
		}
		vehicles := statusResponse.Vehicles
		if len(vehicles) < 1 {
			continue
		}
		telem := vehicles[0].Telem
		if floatCompare(telem.Latitude, initialLat, 0.0001) && floatCompare(telem.Longitude, initialLon, 0.0001) && floatCompare(float64(telem.Altitude), 20.0, 0.1) {
			stopChan <- true
			break
		}

	}
	fmt.Println("At home base")

	//Clear the navigation query from the queue
	query.ClearQueries()

	landSeq := query.SetNavMode(0, icarus.LAND_NOW)
	responseChan, _ = query.Execute()

	fmt.Println("Landing...")
	response = <-responseChan
	landResponse, ok := response.Get(landSeq)
	if !ok {
		fmt.Println("Land response not found")
		os.Exit(1)
	}
	if landResponse.Ok {
		fmt.Println("UAV landing initialized")
	} else {
		fmt.Println("Unable to land")
	}
}
func floatCompare(a, b float64, diff float64) bool {
	if math.Abs(a-b) < diff {
		return true
	}
	return false
}

func parseRadar(id int, info string) RadarPing {
	newPing := RadarPing{VehicleID: id}
	if id == 0 {
		return newPing
	}
	//Radar information format: "vehicleType Latitude Longitude Altitude Heading"
	splitInfo := strings.Split(info, " ")
	if len(splitInfo) != 5 {
		return newPing
	}
	var err error
	vType, err := strconv.ParseInt(splitInfo[0], 10, 64)
	if err != nil {
		return RadarPing{VehicleID: id}
	}
	newPing.VehicleType = int(vType)
	newPing.Latitude, err = strconv.ParseFloat(splitInfo[1], 64)
	if err != nil {
		return RadarPing{VehicleID: id}
	}
	newPing.Longitude, err = strconv.ParseFloat(splitInfo[2], 64)
	if err != nil {
		return RadarPing{VehicleID: id}
	}
	var tmpFloat float64
	tmpFloat, err = strconv.ParseFloat(splitInfo[3], 32)
	if err != nil {
		return RadarPing{VehicleID: id}
	}
	newPing.Altitude = float32(tmpFloat)
	tmpFloat, err = strconv.ParseFloat(splitInfo[4], 32)
	if err != nil {
		return RadarPing{VehicleID: id}
	}
	newPing.Heading = float32(tmpFloat)

	return newPing
}
