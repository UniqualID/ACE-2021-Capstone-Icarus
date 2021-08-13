package icarusClient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/net/publicsuffix"
)

type QueryType uint32

const (
	SingleVehicle QueryType = iota
	AllVehicles
	WaypointList
	NewVehicle
	GotoLocation
	PayloadQuery
	RemoveVehicle
	StartStream
	EndStream
	Error
	SetMode
	Update
)

func (qt QueryType) String() string {
	switch qt {
	case SingleVehicle:
		return "SingleVehicle"
	case AllVehicles:
		return "AllVehicles"
	case WaypointList:
		return "WaypointList"
	case NewVehicle:
		return "NewVehicle"
	case GotoLocation:
		return "Goto"
	case PayloadQuery:
		return "PayloadQuery"
	case RemoveVehicle:
		return "RemoveVehicle"
	case StartStream:
		return "StartStream"
	case EndStream:
		return "EndStream"
	case Error:
		return "Error"
	case SetMode:
		return "SetMode"
	case Update:
		return "Update"
	default:
		return "UnkownQuery"
	}
}

type PayloadType int

const (
	InvalidPayload      PayloadType = 0
	AllPayloads                     = 0
	ThermalLance                    = 3
	Camera                          = 4
	Fuel                            = 5
	Phosphex                        = 7
	PhosphexRemediation             = 8
	AirRadar                        = 9
	AntiMatterMissile               = 10
	AllRadar                        = 11
	GroundRadar                     = 12
	SAM                             = 13
	Cargo                           = 14
	SeekerMissile                   = 15
)

func (pt PayloadType) String() string {
	switch pt {
	case InvalidPayload:
		return "InvalidPayload"
	case ThermalLance:
		return "ThermalLance"
	case Camera:
		return "Camera"
	case Fuel:
		return "Fuel"
	case Phosphex:
		return "Phosphex"
	case PhosphexRemediation:
		return "PhosphexRemediation"
	case AirRadar:
		return "AirRadar"
	case AntiMatterMissile:
		return "AntiMatterMissile"
	case AllRadar:
		return "AllRadar"
	case GroundRadar:
		return "GroundRadar"
	default:
		return "InvalidPayload"
	}
}

type CmdType int32

const (
	GOTO   CmdType = 1
	LOITER CmdType = 2
	JUMP   CmdType = 3
)

type AutopilotType int32

const (
	FIXED_WING    AutopilotType = 1
	ROTOR_CRAFT   AutopilotType = 2
	SURFACE_CRAFT AutopilotType = 3
)

type NavMode int32

const (
	TAKE_OFF          NavMode = 1
	HOME              NavMode = 2
	RALLY             NavMode = 3
	LINGER_NOW        NavMode = 4
	LINGER_WAYPOINT   NavMode = 5
	FOLLOW_THE_LEADER NavMode = 6
	NAVIGATION        NavMode = 7
	LAND_NOW          NavMode = 8
	LAND_WAYPOINT     NavMode = 9
	MANUAL_CONTROL    NavMode = 10
)

//The authentication information the Icarus Server is expecting to initiate a session
type AuthPacket struct {
	Username string
	Password string
}

type Command struct {
	Type     CmdType `protobuf:"varint,1,req,name=type,enum=c3po.pb.Command_Type" json:"type,omitempty"`
	Waypoint Goto    `protobuf:"bytes,2,opt,name=waypoint" json:"waypoint,omitempty"`
	Jump     *Jump   `protobuf:"bytes,3,opt,name=jump" json:"jump,omitempty"`
}

type Cmd struct {
	Cmd Command
}

type Goto struct {
	Waypoint Position `protobuf:"bytes,1,req,name=waypoint" json:"waypoint,omitempty"`
	Velocity float32  `protobuf:"fixed32,2,opt,name=velocity" json:"velocity,omitempty"`
	Radius   float32  `protobuf:"fixed32,3,opt,name=radius" json:"radius,omitempty"`
	Seconds  uint32   `protobuf:"varint,4,opt,name=seconds" json:"seconds,omitempty"`
	Heading  float32  `protobuf:"fixed32,5,opt,name=heading" json:"heading,omitempty"`
}

type Jump struct {
	Target uint32 `protobuf:"varint,1,req,name=target" json:"target,omitempty"`
	Count  uint32 `protobuf:"varint,2,opt,name=count" json:"count,omitempty"`
}

type Position struct {
	Latitude  float64 `protobuf:"fixed64,1,req,name=latitude" json:"latitude,omitempty"`
	Longitude float64 `protobuf:"fixed64,2,req,name=longitude" json:"longitude,omitempty"`
	Altitude  float32 `protobuf:"fixed32,3,opt,name=altitude" json:"altitude,omitempty"`
}

//The information the Icarus Server is expecting in a query
type IcarusQuery struct {
	QueryId         uint32
	Type            QueryType
	VehicleId       uint32
	TeamId          uint32
	VehicleCallsign string
	Daedalus        DaedalusQuery
	CmdList         []Cmd
	VConfig         VehicleConfig
}

func (q IcarusQuery) String() string {
	return fmt.Sprintf("{ QueryId: %v, Type: %v, VehicleId: %v, TeamId: %v, DaedalusQuery: %v, CmdList: %v, VConfig: %v}", q.QueryId, q.Type, q.VehicleId, q.TeamId, q.Daedalus, q.CmdList, q.VConfig)
}

//The information expected from an Icarus Server response
type IcarusResponse struct {
	ResponseId      uint32
	Type            QueryType
	Vehicles        []VehicleStatus
	PayloadResponse DaedalusResponse
	Ok              bool
	Message         string
}

type DaedalusType uint8

const (
	Configure DaedalusType = iota
	Enable
	Status
	Execute
	Unknown
	CargoStatus
)

const (
	DefaultC3poTime     int32 = 1000
	DefaultDaedalusTime int32 = 10000
)

//The information used to interact with the Daedalus Server
type DaedalusQuery struct {
	Type           DaedalusType
	Configurations []PayloadStatus
	PayloadID      PayloadType
	Enabled        bool
	Action         int32
	Parameters     map[string]string
	Target         int32
}

type DaedalusResponse struct {
	Type       DaedalusType
	ErrorCode  int32
	Error      string
	Status     []PayloadStatus
	File       []string
	Parameters map[string]string
	Radar      map[int32]RadarPing
}

/*RadarPing represents an object found by one of the radar payload systems.

StructureType will be 0 for vehicles or 1 for infrastructure (such as airbases).

The definition of Type depends on if the RadarPing represents a vehicle (type is vehicle role) or infrastructure (type is infrastructure type)

Vehicle Roles:
  1: FIGHTER
  2: BOMBER
  4: ISR
  5: ROVER
  6: MULTI
  7: WMD
  8: SAM
  9: RADAR
 10: LOGISTICS
 11: SPEC-OPS
Infrastructure Types:
  1: Airbase
  2: Depot
  3: Datacenter
  4: Fort
  5: WMD-PRODUCTION
  6: ENCAMPMENT
  7: AnimalHospital
  8: EMBASSY
  9: NETWORK-NODE
 10: INSTALLATION
*/
type RadarPing struct {
	Type          int32
	Latitude      float64
	Longitude     float64
	Altitude      float32
	Heading       float32
	StructureType int32
}

//Connection information about a single vehicle
type VehicleConfig struct {
	Ip, C3poPort, DaedalusPort         string
	NavMode                            NavMode
	ProxyMode                          int
	ProxyConfig                        []string
	DaedalusCert, DaedalusKey          []byte
	TeamId                             uint32
	C3poUpdateTime, DaedalusUpdateTime int32
}

type Vehicle struct {
	AutopilotType AutopilotType
	Mode          NavMode
	Valid         bool
}

//Status information about a single vehicle
type VehicleStatus struct {
	VehicleId       uint32
	VehicleCallsign string
	VehicleType     Vehicle
	Telem           Telemetry
	PayStatus       map[int]PayloadStatus
	CmdList         []Cmd
	ActiveWaypoint  Command
	LastComms       uint64
	Available       bool
	VConfig         VehicleConfig
}

//Payload status information from a vehicle
type PayloadStatus struct {
	Id         PayloadType
	Resources  int32
	Name       string
	Enabled    bool
	Parameters map[string]string
	Radar      map[int32]RadarPing
}

type Telemetry struct {
	Latitude, Longitude         float64
	Altitude, Heading, Velocity float32
}

//A list of queries that can be sent to the included IP:Port in a single packet
type QueryPackage struct {
	currentQueryId uint32
	jar            *cookiejar.Jar
	tls            *tls.Config
	Ip, Port       string
	Queries        []IcarusQuery
}

//Remove all current queries from the package. This allows a single QueryPackage to be used multiple times.
func (q *QueryPackage) ClearQueries() bool {
	q.Queries = make([]IcarusQuery, 0)
	if len(q.Queries) > 0 {
		return false
	}
	return true
}

//GetVehicleStatus adds a single status request to the query. It will return the status information of the vehicle with the given ID.
func (q *QueryPackage) GetVehicleStatus(vehicleID int) uint32 {
	id := q.currentQueryId
	q.currentQueryId++
	query := IcarusQuery{
		QueryId:   id,
		Type:      SingleVehicle,
		VehicleId: uint32(vehicleID),
	}
	q.Queries = append(q.Queries, query)
	return id
}

//GetAllVehicleStatus adds a status request for all vehicles to the query. It will return the status information of all connected vehicles.
func (q *QueryPackage) GetAllVehicleStatus() uint32 {
	id := q.currentQueryId
	q.currentQueryId++
	query := IcarusQuery{
		QueryId: id,
		Type:    AllVehicles,
	}
	q.Queries = append(q.Queries, query)
	return id
}

//StartStatusStream initiates a status update stream for the given vehicle. Once received, the IcarusServer application will send status responses for the given vehicle at regular intervals without the need to send the request again. This can be stopped using the StopStatusStream query.
//
func (q *QueryPackage) StartStatusStream(vehicleID int) uint32 {
	id := q.currentQueryId
	q.currentQueryId++
	query := IcarusQuery{
		QueryId:   id,
		Type:      StartStream,
		VehicleId: uint32(vehicleID),
	}
	q.Queries = append(q.Queries, query)
	return id
}

//StopStatusStream stops a status update stream for the given vehicle. Once received, the IcarusServer application will stop sending status responses for the given vehicle. This can be restarted by using the StartStatusStream query.
func (q *QueryPackage) StopStatusStream(vehicleID int) uint32 {
	id := q.currentQueryId
	q.currentQueryId++
	query := IcarusQuery{
		QueryId:   id,
		Type:      EndStream,
		VehicleId: uint32(vehicleID),
	}
	q.Queries = append(q.Queries, query)
	return id
}

//GetWaypointList adds a 'retrieve waypoints' request to the query. The response will contain a list of waypoints retrieved from the given vehicle.
func (q *QueryPackage) GetWaypointList(vehicleID int) uint32 {
	id := q.currentQueryId
	q.currentQueryId++
	query := IcarusQuery{
		QueryId:   id,
		Type:      WaypointList,
		VehicleId: uint32(vehicleID),
	}
	q.Queries = append(q.Queries, query)
	return id
}

//AddNewVehicle adds a 'connect to vehicle' command to the query. Once received, the IcarusServer application will add a vehicle with the provided configuration to they connection list and attempt to communicate with that vehicle.
func (q *QueryPackage) AddNewVehicle(ip, c3poPort, callsign, daedalusPort string, proxyMode int, proxyConfig []string, daedalusCert, daedalusKey []byte, teamID int, c3poUpdate, daedalusUpdate int32) uint32 {
	if proxyConfig == nil {
		proxyConfig = make([]string, 1)
	}
	if daedalusCert == nil {
		daedalusCert = make([]byte, 0)
	}
	if daedalusKey == nil {
		daedalusKey = make([]byte, 0)
	}

	id := q.currentQueryId
	q.currentQueryId++
	config := VehicleConfig{
		Ip:                 ip,
		C3poPort:           c3poPort,
		DaedalusPort:       daedalusPort,
		ProxyMode:          proxyMode,
		ProxyConfig:        proxyConfig,
		DaedalusCert:       daedalusCert,
		DaedalusKey:        daedalusKey,
		TeamId:             uint32(teamID),
		C3poUpdateTime:     c3poUpdate,
		DaedalusUpdateTime: daedalusUpdate,
	}
	query := IcarusQuery{
		QueryId:         id,
		Type:            NewVehicle,
		VehicleCallsign: callsign,
		VConfig:         config,
	}
	q.Queries = append(q.Queries, query)
	return id
}

//RemoveVehicle add a remove request to the query. This will stop the Icarus server from communicating with the given vehicle and remove it from all internal datastructures.
func (q *QueryPackage) RemoveVehicle(vehicleID int) uint32 {
	id := q.currentQueryId
	q.currentQueryId++
	query := IcarusQuery{
		QueryId:   id,
		Type:      RemoveVehicle,
		VehicleId: uint32(vehicleID),
	}
	q.Queries = append(q.Queries, query)
	return id
}

func newWaypoint(lat, lon float64, alt float32) Position {
	return Position{Latitude: lat, Longitude: lon, Altitude: alt}
}

func newCmd(cmdType CmdType, position Position, velocityToTarget, turnRadius float32, lingerSeconds uint32, transitHeading float32) Cmd {
	target := Goto{Waypoint: position, Velocity: velocityToTarget, Radius: turnRadius, Seconds: lingerSeconds, Heading: transitHeading}
	cmd := Command{Type: cmdType, Waypoint: target}
	command := Cmd{Cmd: cmd}
	return command

}

//AddCmd adds a C3PO command to an existing command list. This is used to send a vehicle to a set of waypoint locations using the provided information. To create the list, pass nil as the cmdList argument to AddCmd. By default, cmdType should be set to GOTO unless there is a linger at a given waypoint. If a linger is required, use LOITER as the cmdType.
func AddCmd(cmdList []Cmd, cmdType CmdType, lat, lon float64, alt, velocityToTarget, turnRadius float32, lingerSeconds uint32, transitHeading float32) []Cmd {
	if cmdList == nil {
		cmdList = make([]Cmd, 0)
	}
	waypoint := newWaypoint(lat, lon, alt)
	cmd := newCmd(cmdType, waypoint, velocityToTarget, turnRadius, lingerSeconds, transitHeading)
	cmd.Cmd.Jump = nil
	cmdList = append(cmdList, cmd)
	return cmdList
}

//Goto sets a vehicle's waypoint list to the provided cmdList. This will cause the vehicle to abandon any active waypoints and start navigating through the new sat of waypoints.
func (q *QueryPackage) Goto(vehicleID int, cmdList []Cmd) uint32 {
	id := q.currentQueryId
	q.currentQueryId++
	query := IcarusQuery{
		QueryId:   id,
		Type:      GotoLocation,
		VehicleId: uint32(vehicleID),
		CmdList:   cmdList,
	}
	q.Queries = append(q.Queries, query)
	return id
}

//AddPayloadConfig creates and appends to a list of payload configurations that can be used to configure an instance of a Daedalus Server.
func AddPayloadConfig(configList []PayloadStatus, name string, payloadID PayloadType, resources int, enabled bool) []PayloadStatus {
	if configList == nil {
		configList = make([]PayloadStatus, 0)
	}
	config := PayloadStatus{
		Id:        payloadID,
		Resources: int32(resources),
		Name:      name,
		Enabled:   enabled,
	}
	configList = append(configList, config)
	return configList
}

//ConfigurePayloads sets the loaded payloads on the specified vehicle to the given configurations
func (q *QueryPackage) ConfigurePayloads(vehicleID int, configs []PayloadStatus) uint32 {
	query := DaedalusQuery{
		Type:           Configure,
		Configurations: configs,
	}
	return q.PayloadQuery(vehicleID, query)
}

//EnablePayload enables or disables the given payload on the specified vehicle
func (q *QueryPackage) EnablePayload(vehicleID int, payloadID PayloadType, enabled bool) uint32 {
	query := DaedalusQuery{
		Type:      Enable,
		PayloadID: payloadID,
		Enabled:   enabled,
	}
	return q.PayloadQuery(vehicleID, query)
}

//StatusPayload returns the payload status of a specific payload (identified by the payload ID value), or for all loaded payloads (if the payloadID value is 0).
func (q *QueryPackage) StatusPayload(vehicleID int, payloadID PayloadType) uint32 {
	query := DaedalusQuery{
		Type:      Status,
		PayloadID: payloadID,
	}
	return q.PayloadQuery(vehicleID, query)
}

//CargoStatus returns the status of the cargo payload.
func (q *QueryPackage) CargoStatus(vehicleID int) uint32 {
	query := DaedalusQuery{
		Type: CargoStatus,
	}
	return q.PayloadQuery(vehicleID, query)
}

//ExecutePayload performs the given payload action on the specified vehicle.
func (q *QueryPackage) ExecutePayload(vehicleID int, payloadID PayloadType, action int, params map[string]string, target int) uint32 {
	query := DaedalusQuery{
		Type:       Execute,
		PayloadID:  payloadID,
		Action:     int32(action),
		Parameters: params,
		Target:     int32(target),
	}
	return q.PayloadQuery(vehicleID, query)

}

//PayloadQuery attaches the given Daedalus query to the next Icarus query execute
func (q *QueryPackage) PayloadQuery(vehicleID int, daedalusQuery DaedalusQuery) uint32 {
	id := q.currentQueryId
	q.currentQueryId++
	query := IcarusQuery{
		QueryId:   id,
		Type:      PayloadQuery,
		VehicleId: uint32(vehicleID),
		Daedalus:  daedalusQuery,
	}
	q.Queries = append(q.Queries, query)
	return id
}

//SetNavMode sets the navigation mode of the vehicle. This can be set to any valid C3PO navigation mode including TAKE_OFF, NAVIGATE, HOME, LAND_NOW or any other mode found in the C3PO library.
func (q *QueryPackage) SetNavMode(vehicleID int, mode NavMode) uint32 {
	id := q.currentQueryId
	q.currentQueryId++
	config := VehicleConfig{
		NavMode: mode,
	}
	query := IcarusQuery{
		QueryId:   id,
		Type:      SetMode,
		VehicleId: uint32(vehicleID),
		VConfig:   config,
	}
	q.Queries = append(q.Queries, query)
	return id
}

//Update initiates an application update by downloading a new DEB file and initiating 'apt install' on the DEB file
func (q *QueryPackage) Update(versionNumber string) uint32 {
	id := q.currentQueryId
	q.currentQueryId++
	query := IcarusQuery{
		QueryId:         id,
		Type:            Update,
		VehicleCallsign: versionNumber,
	}
	q.Queries = append(q.Queries, query)
	return id
}

//GetVehicleID returns a vehicle ID along with a map of all callsign->vids when given a vehicle callsign
//This query package must be authenticated before calling GetVehicleID
func (q *QueryPackage) GetVehicleID(callsign string) (int, map[string]int) {
	allStatusSeq := q.GetAllVehicleStatus()
	defer q.ClearQueries()
	responseChan, _ := q.Execute()
	response := <-responseChan
	statusResponse, ok := response.Get(allStatusSeq)
	if !ok {
		return -1, make(map[string]int)
	}
	returnMap := make(map[string]int)
	for _, veh := range statusResponse.Vehicles {
		name := veh.VehicleCallsign
		id := veh.VehicleId
		returnMap[name] = int(id)
	}
	requestedID, ok := returnMap[callsign]
	if !ok {
		return -1, returnMap
	}
	return requestedID, returnMap
}

//Execute uploads the list of queries to the IcarusServer application and returns true on success. All responses will be returned through the provided response channel. The Authenticate function must be called before this function or Execute will fail.
func (q *QueryPackage) Execute() (chan QueryResponse, chan bool) {
	responseChan := make(chan QueryResponse, 1)
	stopChan := make(chan bool)

	go func() {
		defer close(responseChan)
		url := url.URL{Scheme: "wss", Host: q.Ip + ":" + q.Port, Path: "/api/v1"}
		var dialer = &websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 5 * time.Second,
			TLSClientConfig:  q.tls,
			Jar:              q.jar,
		}

		conn, _, err := dialer.Dial(url.String(), nil)
		if err != nil {
			fmt.Println("Execute dial error:", err)
			return
		}
		defer conn.Close()

		//TODO: set conn read/write timeout?
		readWriteTimeout := time.Duration(20)
		conn.SetReadDeadline(time.Now().Add(time.Second * readWriteTimeout))
		conn.SetWriteDeadline(time.Now().Add(time.Second * readWriteTimeout))

		streaming := false
		for _, v := range q.Queries {
			if v.Type == StartStream {
				streaming = true
			}
		}
		packet, err := json.Marshal(q.Queries)
		if err != nil {
			fmt.Println("Error marshaling query:", err)
			return
		}

		err = conn.WriteMessage(websocket.BinaryMessage, packet)
		if err != nil {
			fmt.Println("Error writing to socket:", err)
			return
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading from socket:", err.Error())
			return
		}

		response := QueryResponse{}
		err = json.Unmarshal(message, &response.Responses)
		if err != nil {
			fmt.Println("Error parsing response:", err)
			return
		}

		if streaming {
		Loop:
			for {
				conn.SetReadDeadline(time.Now().Add(time.Second * readWriteTimeout))
				_, message, err := conn.ReadMessage()
				if err != nil {
					fmt.Println("Error reading from socket:", err)
					break
				}

				response := QueryResponse{}
				err = json.Unmarshal(message, &response.Responses)
				if err != nil {
					fmt.Println("Error parsing response:", err)
					continue
				}
				responseChan <- response

				select {
				case <-stopChan:
					break Loop
				default:
				}

			}
		} else {
			conn.SetWriteDeadline(time.Now().Add(time.Second * readWriteTimeout))
			err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "Query Complete"))
			if err != nil {
				fmt.Println("Error closing socket:", err)
			}

			responseChan <- response
			return
		}
		conn.SetWriteDeadline(time.Now().Add(time.Second * readWriteTimeout))
		err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "Query Complete"))
		if err != nil {
			fmt.Println("Error closing socket:", err)
		}
		return
	}()
	return responseChan, stopChan
}

//Authenticate sends the authentication packet to the server and saves the resulting cookie in the QueryPackage for later use. This function must be called before Execute is called or Execute will fail. It returns true if the connection was successful and false if it failed.
func (q *QueryPackage) Authenticate(user, pass string) (string, bool) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return err.Error(), false
	}

	q.jar = jar
	url := url.URL{Scheme: "https", Host: q.Ip + ":" + q.Port, Path: "/api/v1/login"}
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	q.tls = tlsConfig
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{
		Jar:       q.jar,
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	authReq := AuthPacket{
		Username: user,
		Password: pass,
	}

	packet, err := json.Marshal(authReq)
	if err != nil {
		return err.Error(), false
	}

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(packet))
	if err != nil {
		return err.Error(), false
	}

	resp, err := client.Do(req)
	if err != nil {
		return err.Error(), false
	}

	defer resp.Body.Close()
	status := resp.StatusCode
	if status > 299 {
		return "HTTP Error: " + resp.Status, false
	}
	if status > 199 && status < 300 {
		return "", true
	}
	return "", true
}

//ShowQuery prints the current query in text JSON format
func (q *QueryPackage) ShowQuery() {
	json, err := json.MarshalIndent(q.Queries, "", "\t")
	if err != nil {
		fmt.Println("Error marshaling query")
		return
	}
	fmt.Println(string(json))
}

//A list of responses from the Icarus Server
type QueryResponse struct {
	Responses []IcarusResponse
}

//ShowResponse prints the current response in text JSON format
func (q *QueryResponse) ShowResponse() {
	json, err := json.MarshalIndent(q.Responses, "", "\t")
	if err != nil {
		fmt.Println("Error marshaling response")
		return
	}
	fmt.Println(string(json))

}

//Get returns the IcarusResponse corresponding to the IcarusQuery with the given ID. This can be used to get a single response out of the list of responses returned by the server.
func (r *QueryResponse) Get(queryID uint32) (IcarusResponse, bool) {
	for _, v := range r.Responses {
		if v.ResponseId == uint32(queryID) {
			return v, true
		}
	}
	return IcarusResponse{}, false
}

//NewQuery creates the datastructures necessary to send quiries to IcarusServer
func NewQuery(ip, port string) QueryPackage {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		fmt.Println("Cookie jar unavailable")
		return QueryPackage{}
	}
	return QueryPackage{
		Ip:      ip,
		Port:    port,
		Queries: make([]IcarusQuery, 0),
		jar:     jar,
	}
}

//EmptyParams returns an empty map correctly formed for payload operations
func EmptyParams() map[string]string {
	return make(map[string]string)
}

//LoadCargo returns a map correctly formatted for loading the given amount of cargo
func LoadCargo(payload PayloadType, amount int) (map[string]string, error) {
	returnMap := make(map[string]string)
	jsonMap := make(map[PayloadType]int)
	jsonMap[payload] = amount
	jsonData, err := json.Marshal(jsonMap)
	if err != nil {
		return nil, err
	}
	returnMap["load"] = string(jsonData)
	return returnMap, nil
}

//UnloadCargo returns a map correctly formatted for unloading the given amount of cargo
func UnloadCargo(payload PayloadType, amount int) (map[string]string, error) {
	returnMap := make(map[string]string)
	jsonMap := make(map[PayloadType]int)
	jsonMap[payload] = amount
	jsonData, err := json.Marshal(jsonMap)
	if err != nil {
		return nil, err
	}
	returnMap["unload"] = string(jsonData)
	return returnMap, nil
}

//LoadMultiCargo returns a map correctly formatted for loading the given multiple cargo types
func LoadMultiCargo(cargo map[PayloadType]int) (map[string]string, error) {
	returnMap := make(map[string]string)
	jsonData, err := json.Marshal(cargo)
	if err != nil {
		return make(map[string]string), err
	}
	returnMap["load"] = string(jsonData)
	return returnMap, nil

}

//UnloadMultiCargo returns a map correctly formatted for loading the given multiple cargo types
func UnloadMultiCargo(cargo map[PayloadType]int) (map[string]string, error) {
	returnMap := make(map[string]string)
	jsonData, err := json.Marshal(cargo)
	if err != nil {
		return make(map[string]string), err
	}
	returnMap["unload"] = string(jsonData)
	return returnMap, nil

}
