//This file contains an example ISR mission
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	icarus "git.ironzone.ace/icarus/icarusClient"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RadarPing struct {
	VehicleID           int
	VehicleType         int
	Latitude, Longitude float64
	Altitude, Heading   float32
}

type MissionTarget struct {
	VehicleID                   int
	Latitude, Longitude         float64
	Altitude, Heading, Velocity float32
	Linger                      int
}

var (
	csvFileName       string
	availableVehicles map[string]int
)

func init() {
	flag.StringVar(&csvFileName, "csv", "isr-deploy.csv", "The file name of the CSV file that contains a mission set")
	flag.Parse()
}

func main() {
	availableVehicles = getAvailableVehicles()
	openFile, err := os.Open(csvFileName)
	if err != nil {
		log.Fatal("Unable to open CSV:", err.Error())
	}
	reader := csv.NewReader(openFile)
	reader.Comment = '#'
	vehicleMap := make(map[int][]MissionTarget)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error reading CSV:", err.Error())
		}
		var newMissionTarget MissionTarget
		newMissionTarget = parseRecord(record)
		if newMissionTarget.VehicleID > -1 {
			targetSlice := vehicleMap[newMissionTarget.VehicleID]
			vehicleMap[newMissionTarget.VehicleID] = append(targetSlice, newMissionTarget)

		}

	}
	var wg sync.WaitGroup
	for vehicleID, targets := range vehicleMap {
		wg.Add(1)
		go runMission(vehicleID, targets, &wg)
		time.Sleep(3 * time.Second)
	}
	wg.Wait()
}

func getAvailableVehicles() map[string]int {
	query := icarus.NewQuery("127.0.0.1", "9443")

	_, ok := query.Authenticate("test1", "testing")
	if !ok {
		return make(map[string]int)
	}
	_, vehicles := query.GetVehicleID("")
	return vehicles
}

func parseRecord(record []string) MissionTarget {
	if len(record) != 7 {
		return MissionTarget{VehicleID: -1}
	}

	var err error
	newMissionTarget := MissionTarget{}
	var tmpID int64
	tmpID, err = strconv.ParseInt(record[0], 10, 64)
	if err != nil {
		newID, ok := availableVehicles[record[0]]
		if !ok {
			return MissionTarget{VehicleID: -1}
		}
		tmpID = int64(newID)
	}
	tmpLat, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return MissionTarget{VehicleID: -1}
	}
	tmpLon, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return MissionTarget{VehicleID: -1}
	}
	tmpAlt, err := strconv.ParseFloat(record[3], 32)
	if err != nil {
		return MissionTarget{VehicleID: -1}
	}
	tmpHead, err := strconv.ParseFloat(record[4], 32)
	if err != nil {
		return MissionTarget{VehicleID: -1}
	}
	tmpVel, err := strconv.ParseFloat(record[5], 32)
	if err != nil {
		return MissionTarget{VehicleID: -1}
	}
	tmpLinger, err := strconv.ParseInt(record[6], 10, 64)
	if err != nil {
		return MissionTarget{VehicleID: -1}
	}
	newMissionTarget.VehicleID = int(tmpID)
	newMissionTarget.Latitude = tmpLat
	newMissionTarget.Longitude = tmpLon
	newMissionTarget.Altitude = float32(tmpAlt)
	newMissionTarget.Heading = float32(tmpHead)
	newMissionTarget.Velocity = float32(tmpVel)
	newMissionTarget.Linger = int(tmpLinger)
	return newMissionTarget
}
func runMission(vehicleID int, targets []MissionTarget, wg *sync.WaitGroup) {
	defer wg.Done()
	//Create a new query pointed at the IcarusServer instance
	query := icarus.NewQuery("127.0.0.1", "9443")

	resp, ok := query.Authenticate("test1", "testing")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
		return
	}

	//Load 50 units of Fuel
	configs := icarus.AddPayloadConfig(nil, "Fuel", icarus.Fuel, 50, true)

	//Assumes ISR vehicle, may add other payloads here using above syntax if munitions are needed
	configSeq := query.ConfigurePayloads(vehicleID, configs)
	//Create the chanels to get execute responses

	responseChan, _ := query.Execute()
	fmt.Println("Loading fuel...")
	response := <-responseChan

	configResponse, ok := response.Get(configSeq)
	if !ok {
		fmt.Println("Refuel response not found")
		return
	}

	if configResponse.Ok {
		fmt.Println("Refueling complete")
	} else {
		fmt.Println("Error during refueling:", configResponse.Message)
	}

	//Clear the refuel query from the queue
	query.ClearQueries()

	//Change mode to Take off
	takeOffSeq := query.SetNavMode(vehicleID, icarus.TAKE_OFF)
	responseChan, _ = query.Execute()
	fmt.Println("Taking off...")
	response = <-responseChan
	takeOffResponse, ok := response.Get(takeOffSeq)
	if !ok {
		fmt.Println("TAKE_OFF response not found")
		return
	}

	if takeOffResponse.Ok {
		fmt.Println("UAV launch complete")
	} else {
		fmt.Println("Error during take off:", takeOffResponse.Message)
	}
	//Clear the take off query from the queue
	query.ClearQueries()

	//Sleep to ensure UAV is in air before navigating
	time.Sleep(1 * time.Second)

	//Change mode to Navigate
	navSeq := query.SetNavMode(vehicleID, icarus.NAVIGATION)

	responseChan, _ = query.Execute()
	fmt.Println("Entering Navigate mode...")
	response = <-responseChan
	navResponse, ok := response.Get(navSeq)
	if !ok {
		fmt.Println("NAVIGATE response not found")
		return
	}

	if navResponse.Ok {
		fmt.Println("UAV ready to navigate")
	} else {
		fmt.Println("Error during mode change:", navResponse.Message)
	}
	//Clear the mode change query from the queue
	query.ClearQueries()

	//Navigate to point of interest and linger for 60 seconds then navigate to secondary location
	//The vehicle will loop through the given targets until it is told to do something else
	cmdList := getCmdList(targets)

	gotoSeq := query.Goto(vehicleID, cmdList)
	responseChan, _ = query.Execute()
	fmt.Println("Navigating to waypoint...")
	response = <-responseChan
	gotoResponse, ok := response.Get(gotoSeq)
	if !ok {
		fmt.Println("Go to response not found")
		return
	}

	if gotoResponse.Ok {
		fmt.Println("Navigating to waypoint")
	} else {
		fmt.Println("Error during navigation:", gotoResponse.Message)
	}
	//Clear the navigation query from the queue
	query.ClearQueries()
}

func getCmdList(targets []MissionTarget) []icarus.Cmd {
	var cmdList []icarus.Cmd = nil
	for _, target := range targets {
		var cmdType icarus.CmdType
		if target.Linger > 0 {
			cmdType = icarus.LOITER
		} else {
			cmdType = icarus.GOTO
		}
		cmdList = icarus.AddCmd(cmdList, cmdType, target.Latitude, target.Longitude, target.Altitude, target.Velocity, 0, uint32(target.Linger), target.Heading)
	}
	return cmdList

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
