package main

import (
	"fmt"
	"os"
	"strconv"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

func main() {
	if len(os.Args) != 6 {
		fmt.Println("Missing/Invalid arguments")
		os.Exit(1)
	}

	var query = icarus.NewQuery("10.59.144.202", "179")
	resp, ok := query.Authenticate("valinar", "gitlabsux;(bloodybreach")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
		os.Exit(1)
	}

	vehID, _ := strconv.Atoi(os.Args[1])
	navSeq := query.SetNavMode(vehID, icarus.NAVIGATION)
	responseChan, _ := query.Execute()
	response := <-responseChan
	_, ok = response.Get(navSeq)
	if !ok {
		fmt.Println("No response")
		os.Exit(1)
	}

	destLat, _ := strconv.ParseFloat(os.Args[2], 64)
	destLon, _ := strconv.ParseFloat(os.Args[3], 64)
	alt, _ := strconv.ParseFloat(os.Args[4], 32)
	destAlt := float32(alt)
	speed, _ := strconv.ParseFloat(os.Args[5], 32)
	cruiseSpeed := float32(speed)

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
	fmt.Printf("Moving vehicle %d to %.4f %.0f %.4f at %.0f meters/sec.\n", vehID, destLat, destAlt, destLon, cruiseSpeed)

	// Trigger when in range
	var inRange bool = false
	for !inRange {
		statSeq := query.GetVehicleStatus(vehID)
		responseChan, _ := query.Execute()
		response := <-responseChan
		statusResponse, ok := response.Get(statSeq)
		if !ok {
			fmt.Println("No response")
			os.Exit(1)
		}
		v := statusResponse.Vehicles[0]
		if v.Telem.Latitude > destLat-.0005 && v.Telem.Latitude < destLat+.0005 && v.Telem.Longitude > destLon-.0005 && v.Telem.Longitude < destLon+.0005 {
			inRange = true
		}
	}

	// Go to linger
	// navSeq = query.SetNavMode(vehID, icarus.LINGER_NOW)
	// responseChan, _ = query.Execute()
	// response = <-responseChan
	// _, ok = response.Get(navSeq)
	// if !ok {
	// 	fmt.Println("No response")
	// }
	// fmt.Printf("Moving vehicle %d to linger at %.4f %.0f %.4f.\n", vehID, destLat, destAlt, destLon)

	//picture
	configs := query.ExecutePayload(vehID, 4, 1, nil, 0)

	// configSeq := query.ConfigurePayloads(vehID, configs) // I AM NOT SURE WHAT THIS LINE DOES BUT I THINK IT RETURNS THE SAME THING AS EXECUTE PAYLOAD
	responseChan, _ = query.Execute()

	fmt.Printf("Vehicle %v taking picture...\n", vehID)
	response = <-responseChan

	configResponse, ok := response.Get(configs) // changed from configSeq
	if !ok {
		fmt.Println("picture not taken")
		os.Exit(1)
	}

	if configResponse.Ok {
		fmt.Printf("Vehicle %v picture complete\n\n", vehID)
		fmt.Printf("Images returned by vehicle %v: %v \n", vehID, configResponse.PayloadResponse.File)
	} else {
		fmt.Println("Error during picture taken: \n", configResponse.Message)
	}
}
