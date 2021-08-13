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
	VehicleID                int
	Latitude, Longitude      float64
	Altitude                 float32
	PayloadID, PayloadAction int
}

var (
	csvFileName string
)

func init() {
	flag.StringVar(&csvFileName, "csv", "mission-payloads.csv", "The file name of the CSV file that contains a mission set")
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
	tmpID, err := strconv.ParseInt(record[0], 10, 64)
	if err != nil {
		return MissionTarget{VehicleID: -1}
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
	tmpPayloadID, err := strconv.ParseInt(record[4], 10, 64)
	if err != nil {
		return MissionTarget{VehicleID: -1}
	}
	tmpPayloadAction, err := strconv.ParseInt(record[5], 10, 64)
	if err != nil {
		return MissionTarget{VehicleID: -1}
	}

	newMissionTarget.VehicleID = int(tmpID)
	newMissionTarget.Latitude = tmpLat
	newMissionTarget.Longitude = tmpLon
	newMissionTarget.Altitude = float32(tmpAlt)
	newMissionTarget.PayloadID = int(tmpPayloadID)
	newMissionTarget.PayloadAction = int(tmpPayloadAction)
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

	//Continuously monitor vehicle location and execute payloads when at target locations
	//Use Ctrl-C to exit application
	query.StartStatusStream(vehicleID)
	responseChan, _ := query.Execute()
	for {
		response := <-responseChan
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
		for _, target := range targets {
			if checkAtTarget(telem, target) {
				//Execute a payload (take picture here, but could be any payload)
				//Clear the query from the queue
				query.ClearQueries()
				//Final parameter is for a target ID. This parameter is not used for Camera, but is used for Antimatter Missiles and Seeker Missiles
				executeSeq := query.ExecutePayload(vehicleID, icarus.PayloadType(target.PayloadID), target.PayloadAction, icarus.EmptyParams(), 0)
				payloadResponseChan, _ := query.Execute()
				fmt.Println("Taking picture...")
				payloadResponseFromChan := <-payloadResponseChan
				executeResponse, ok := payloadResponseFromChan.Get(executeSeq)
				if !ok {
					fmt.Println("Payload execute response not found")
					continue
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

			}
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
		time.Sleep(time.Second)
	}
}

func checkAtTarget(telem icarus.Telemetry, target MissionTarget) bool {
	if floatCompare(telem.Latitude, target.Latitude, 0.0001) && floatCompare(telem.Longitude, target.Longitude, 0.0001) && floatCompare(float64(telem.Altitude), float64(target.Altitude), 0.1) {
		return true
	}
	return false
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
