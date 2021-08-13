package main

import (
	"fmt"
	"os"

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

	//Output
	fmt.Printf("\n=---------------------= ALL VEHICLE STATUS =---------------------=\n\n")
	for _, v := range statusResponse.Vehicles {
		fmt.Print("=---= ")
		fmt.Print(v.VehicleCallsign, " [", v.VehicleId, "] // Team # ", v.VConfig.TeamId, " =-----------------------------=\n")
		if !v.Available {
			println("DESTROYED\n")
		} else {
			fmt.Printf("Coordinates:\t%.4f, %.0f, %.4f", v.Telem.Latitude, v.Telem.Altitude, v.Telem.Longitude)
			fmt.Printf("\t%.0f m/s", v.Telem.Velocity)
			fmt.Print("\nFuel:\t\t", v.PayStatus[5].Resources)
			fmt.Print("\tMode:\t", v.VehicleType.Mode)
			switch v.VehicleType.Mode {
			case 0:
				fmt.Print(" (NONE)\n")
			case 1:
				fmt.Print(" (TAKE_OFF)\n")
			case 2:
				fmt.Print(" (HOME)\n")
			case 3:
				fmt.Print(" (RALLY)\n")
			case 4:
				fmt.Print(" (LINGER_NOW)\n")
			case 5:
				fmt.Print(" (LINGER_WAYPOINT)\n")
			case 6:
				fmt.Print(" (FOLLOW_THE_LEADER)\n")
			case 7:
				fmt.Print(" (NAVIGATION)\n")
			case 8:
				fmt.Print(" (LAND_NOW)\n")
			case 9:
				fmt.Print(" (LAND_WAYPOINT)\n")
			case 10:
				fmt.Print(" (MANUAL_CONTROL)\n")
			}
			fmt.Println()

			// fmt.Print("\nWaypoint:\t", v.ActiveWaypoint.Type)
			// if v.ActiveWaypoint.Type == 1 {
			// 	fmt.Print(" (GOTO)\nGoing to ", v.ActiveWaypoint.Waypoint.Waypoint.Latitude, ", ",
			// 		v.ActiveWaypoint.Waypoint.Waypoint.Altitude, ", ",
			// 		v.ActiveWaypoint.Waypoint.Waypoint.Longitude, " @ ",
			// 		v.ActiveWaypoint.Waypoint.Velocity, "m/s\n\n")
			// } else if v.ActiveWaypoint.Type == 0 {
			// 	fmt.Print(" (NONE)\n\n")
			// } else if v.ActiveWaypoint.Type == 2 {
			// 	fmt.Print(" (LOITER)\n\n")
			// } else if v.ActiveWaypoint.Type == 3 {
			// 	fmt.Print(" (JUMP)\n\n")
			// }
			// if len(v.CmdList) == 1 {
			// 	fmt.Print("\nCommand:\t", v.CmdList[0].Cmd.Type)
			// 	if v.CmdList[0].Cmd.Type == 1 {
			// 		fmt.Print(" (GOTO)\nGoing to ", v.CmdList[0].Cmd.Waypoint.Waypoint.Latitude, ", ",
			// 			v.CmdList[0].Cmd.Waypoint.Waypoint.Altitude, ", ",
			// 			v.CmdList[0].Cmd.Waypoint.Waypoint.Longitude, " @ ",
			// 			v.CmdList[0].Cmd.Waypoint.Velocity, "m/s\n\n")
			// 	} else if v.CmdList[0].Cmd.Type == 0 {
			// 		fmt.Print(" (NONE)\n\n")
			// 	} else if v.CmdList[0].Cmd.Type == 2 {
			// 		fmt.Print(" (LOITER)\n\n")
			// 	} else if v.CmdList[0].Cmd.Type == 3 {
			// 		fmt.Print(" (JUMP\n\n")
			// 	}
			// } else {
			// 	fmt.Print("\nCommand:\t0 (NONE)\n\n")
			// }
		}

	}
}
