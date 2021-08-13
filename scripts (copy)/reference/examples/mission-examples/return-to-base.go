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
	csvFileName string
)

func init() {
	flag.StringVar(&csvFileName, "csv", "return-to-base.csv", "The file name of the CSV file that contains a mission set")
	flag.Parse()
}

func main() {
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
func parseRecord(record []string) MissionTarget {
	if len(record) != 6 {
		return MissionTarget{VehicleID: -1}
	}
	var err error
	newMissionTarget := MissionTarget{}
	tmpID, err := strconv.ParseInt(record[1], 10, 64)
	if err != nil {
		return MissionTarget{VehicleID: -1}
	}
	tmpLat, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return MissionTarget{VehicleID: -1}
	}
	tmpLon, err := strconv.ParseFloat(record[3], 64)
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
	newMissionTarget.VehicleID = int(tmpID)
	newMissionTarget.Latitude = tmpLat
	newMissionTarget.Longitude = tmpLon
	newMissionTarget.Altitude = 25.0
	newMissionTarget.Heading = float32(tmpHead)
	newMissionTarget.Velocity = float32(tmpVel)
	return newMissionTarget
}
func runMission(vehicleID int, targets []MissionTarget, wg *sync.WaitGroup) {
	defer wg.Done()
	//Only have one base listed for this vehicle
	if len(targets) > 1 {
		return
	}
	target := targets[0]
	//Create a new query pointed at the IcarusServer instance
	query := icarus.NewQuery("127.0.0.1", "9443")

	resp, ok := query.Authenticate("test1", "testing")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
		return
	}

	//Change mode to Navigate
	navSeq := query.SetNavMode(vehicleID, icarus.NAVIGATION)
	responseChan, _ := query.Execute()
	fmt.Println("Entering Navigate mode...")
	response := <-responseChan
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

	//Navigate to home base
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

	//Wait until at initial base then land
	query.StartStatusStream(vehicleID)
	responseChan, stopChan := query.Execute()
	for {
		response = <-responseChan
		//streamed status is always sequence 0
		statusResponse, ok := response.Get(0)
		if !ok {
			continue
		}
		vehicles := statusResponse.Vehicles
		if len(vehicles) < 1 {
			continue
		}
		telem := vehicles[0].Telem
		if floatCompare(telem.Latitude, target.Latitude, 0.0001) && floatCompare(telem.Longitude, target.Longitude, 0.0001) && floatCompare(float64(telem.Altitude), float64(target.Altitude), 0.1) {
			//Once at desired location (the home base) stop checking status and land
			stopChan <- true
			break
		}

	}
	fmt.Printf("Vehicle %d at home base\n", vehicleID)

	//Clear the navigation query from the queue
	query.ClearQueries()

	landSeq := query.SetNavMode(vehicleID, icarus.LAND_NOW)
	responseChan, _ = query.Execute()

	fmt.Println("Landing...")
	response = <-responseChan
	landResponse, ok := response.Get(landSeq)
	if !ok {
		fmt.Println("Land response not found")
		return
	}
	if landResponse.Ok {
		fmt.Println("UAV landing initialized")
	} else {
		fmt.Println("Unable to land")
	}
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
