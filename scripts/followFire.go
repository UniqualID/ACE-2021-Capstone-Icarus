// TODO test the follow and kill script

package main

import (
	"fmt"
	"os"
	"strconv"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

func main() {

	var query icarus.QueryPackage

	if len(os.Args) < 5 {
		fmt.Println("Missing arguments")
		os.Exit(1)
	} else if len(os.Args) == 5 {
		query = icarus.NewQuery("10.59.144.32", "179")
	} else if len(os.Args) == 6 {
		query = icarus.NewQuery(os.Args[5], "179")
	} else {
		fmt.Println("Too many arguments")
		os.Exit(1)
	}
	resp, ok := query.Authenticate("valinar", "gitlabsux;(bloodybreach")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}

	vehID, _ := strconv.Atoi(os.Args[1])
	payID, _ := strconv.Atoi(os.Args[2])
	action, _ := strconv.Atoi(os.Args[3])
	target, _ := strconv.Atoi(os.Args[4])

	for {
		statSeq := query.GetVehicleStatus(vehID)
		responseChan, _ := query.Execute()
		response := <-responseChan
		statusResponse, ok := response.Get(statSeq)
		if !ok {
			fmt.Println("No response")
			os.Exit(1)
		}

		// fuse radar to find stuff
		v := statusResponse.Vehicles[0]
		radarInfo := v.PayStatus[icarus.AllRadar]
		if len(radarInfo.Radar) > 0 {
			for asset, radarPing := range radarInfo.Radar {
				fmt.Println("Checking pings")
				// locate our target
				if asset == int32(target) {
					cmdList := icarus.AddCmd(nil, icarus.GOTO, radarPing.Latitude, radarPing.Longitude, 200, 70, 0, 0, 0)
					// go to the target
					gotoSeq := query.Goto(vehID, cmdList)
					responseChan, _ = query.Execute()
					response = <-responseChan
					_, ok = response.Get(gotoSeq)
					if !ok {
						fmt.Println("No response")
						os.Exit(1)
					}
					fmt.Printf("Vehicle %d at %.4f %.0f %.4f to %.4f %d %.4f at %d meters/sec.\n", vehID, v.Telem.Latitude, v.Telem.Altitude, v.Telem.Longitude, radarPing.Latitude, 200, radarPing.Longitude, 70)

					executePayload(vehID, payID, action, target)
				}
			}
		}
	}
}

// pew pew bitch
func executePayload(vehID int, payID int, action int, target int) {
	if len(os.Args) != 5 {
		fmt.Println("Missing/Invalid arguments")
		os.Exit(1)
	}

	var query = icarus.NewQuery("10.59.144.32", "179")
	resp, ok := query.Authenticate("valinar", "gitlabsux;(bloodybreach")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}

	// execute the shot
	configs := query.ExecutePayload(vehID, icarus.PayloadType(payID), action, nil, target)
	responseChan, _ := query.Execute()
	response := <-responseChan
	_, ok = response.Get(configs)
	if !ok {
		fmt.Println("No response")
		os.Exit(1)
	}

	// message to operator with what happened
	if action == 0 {
		fmt.Printf("Vehicle %v jettison.\n", vehID)
	} else if action == 1 {
		fmt.Printf("Vehicle %v fired.\n", vehID)
	} else {

		fmt.Println("Incorrect action")
	}
}
