package main

import (
	"fmt"
	"os"
	"strconv"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

func main() {

	var query icarus.QueryPackage

	if len(os.Args) < 3 {
		fmt.Println("Missing arguments")
		os.Exit(1)
	} else if len(os.Args) == 3 {
		query = icarus.NewQuery("10.59.144.32", "179")
	} else if len(os.Args) == 4 {
		query = icarus.NewQuery(os.Args[3], "179")
	} else {
		fmt.Println("Too many arguments")
		os.Exit(1)
	}
	resp, ok := query.Authenticate("valinar", "gitlabsux;(bloodybreach")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}

	vehID, _ := strconv.Atoi(os.Args[1])
	target, _ := strconv.Atoi(os.Args[2])

	for {
		statSeq := query.GetVehicleStatus(vehID)
		responseChan, _ := query.Execute()
		response := <-responseChan
		statusResponse, ok := response.Get(statSeq)
		if !ok {
			fmt.Println("No response")
			os.Exit(1)
		}

		v := statusResponse.Vehicles[0]
		radarInfo := v.PayStatus[icarus.GroundRadar]
		if len(radarInfo.Radar) > 0 {
			for asset, radarPing := range radarInfo.Radar {
				fmt.Println("Checking pings")
				if asset == int32(target) {
					cmdList := icarus.AddCmd(nil, icarus.GOTO, radarPing.Latitude, radarPing.Longitude, 200, 70, 0, 0, 0)

					gotoSeq := query.Goto(vehID, cmdList)
					responseChan, _ = query.Execute()
					response = <-responseChan
					_, ok = response.Get(gotoSeq)
					if !ok {
						fmt.Println("No response")
						os.Exit(1)
					}
					fmt.Printf("Vehicle %d at %.4f %.0f %.4f to %.4f %d %.4f at %.0f meters/sec.\n", vehID, v.Telem.Latitude, v.Telem.Altitude, v.Telem.Longitude, radarPing.Latitude, 200, radarPing.Longitude, 70)
				}
			}
		}
	}
}
