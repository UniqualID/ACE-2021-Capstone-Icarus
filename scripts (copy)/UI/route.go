package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

// const (
// 	defaultCert string = `-----BEGIN CERTIFICATE-----
// MIIC9TCCAd2gAwIBAgIRAKRQOLvrvmORBJxJKkhCpWgwDQYJKoZIhvcNAQELBQAw
// EjEQMA4GA1UEChMHQWNtZSBDbzAeFw0xOTAxMDkxNzQ1NTdaFw0yMDAxMDkxNzQ1
// NTdaMBIxEDAOBgNVBAoTB0FjbWUgQ28wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
// ggEKAoIBAQDGzdsktG2DGiQEt7ce1sWlcSc1QNbpLcRemcrGxJKw2JeWYY42R5Le
// +umtLV/xy0+ZIA47iHETj0IFFYjWdixmz5/yHnnnJbz8uKinbk3eTmaR6y+EwSAp
// gFjXFYjRgif4wPk0qnkgHaI+TJXn2dbBnpv0cX34aUKwaCa/qh0XEZ2nqmjjeowx
// mqpD4etICnaMKdJg2Z+da/YG8ExFnwYpzNS9QdfujAxHJ7DoMhPZnyc/sCmaBq+X
// eAZbMHtWFuv/24lA/KyJBmCEQGp2x9tn+HmM89SQOj1yOwOqZB87+rjENhp87rgh
// MD4vB93/Mzk48tC1LYlr7cLoBh22tskRAgMBAAGjRjBEMA4GA1UdDwEB/wQEAwIF
// oDATBgNVHSUEDDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMA8GA1UdEQQIMAaH
// BH8AAAEwDQYJKoZIhvcNAQELBQADggEBAKYVTX6fzjOU63QxdtSs9Ot6GdAqaknQ
// baTiRnAsXJuRCDZpRIEQKFECC7tdEBbyCh5FiyjpVqxn+U2/OcdNHPYxdHRevRWM
// vmNqxeOjha62Bp/JKoN0WR/NfZ7oSHFQYm5kGxTSt6n/BKovWOHI4je0PHD/YITt
// Xw5IgIEPDS2eecE+jrHBU596X1jeHSuk+XdQ9Hmo1WFYE9UK7985Oy77+zaNSKdu
// ZbBWN877w5AMdmswrC6HcDCXY0Sb3Noadfl9VcDqD69hq0vWPqAomrbjyyTUQCY3
// ru4kYaW0u8509KHlpN6ixQWmGAmVhMWtf0g1kPpzJ9HihRtAYUp5Jj4=
// -----END CERTIFICATE-----
// `
// 	defaultKey string = `-----BEGIN RSA PRIVATE KEY-----
// MIIEpgIBAAKCAQEAxs3bJLRtgxokBLe3HtbFpXEnNUDW6S3EXpnKxsSSsNiXlmGO
// NkeS3vrprS1f8ctPmSAOO4hxE49CBRWI1nYsZs+f8h555yW8/Liop25N3k5mkesv
// hMEgKYBY1xWI0YIn+MD5NKp5IB2iPkyV59nWwZ6b9HF9+GlCsGgmv6odFxGdp6po
// 43qMMZqqQ+HrSAp2jCnSYNmfnWv2BvBMRZ8GKczUvUHX7owMRyew6DIT2Z8nP7Ap
// mgavl3gGWzB7Vhbr/9uJQPysiQZghEBqdsfbZ/h5jPPUkDo9cjsDqmQfO/q4xDYa
// fO64ITA+Lwfd/zM5OPLQtS2Ja+3C6AYdtrbJEQIDAQABAoIBAQCMQh4TFkyRCzdQ
// MME8O7Bz2ZIc6yL0njqFt6EtfPA1XooMKcWom/SN5p5IdNPVBmihEtGXxNpqP08H
// wTqqe/M1kdQ5gLDmmGRuNGWgwpyjc9K/rhr3YT2sqgWDsYi2r0o+IP9w3bjZJK8b
// nvLAAZuXPKyw2AVU5gaL6N81p/IgHCoQ6GAkE19Bk2EK/vRoSKtmZsmsUY3RrezC
// HTSX4R4it2eMXBAxEF6stPK0tkv2BvOItoe5BGfWnd4VG2okVcHFkRPiaIiz60iy
// 02m6hei2dTKRe/GxaNNDNW8jNqTyJ4ZbUK8MeFVyCmFe+Zu8/LLF6bffj46HE00u
// yWovzFMxAoGBANq5gT27iGukc2MgXDc/yLKcwmXTU7hGH7TRBCu11Ms9fgInCsrz
// x9SRhI5KvysOVe80uPdtyoMx9+I/hzDDNBj0IrIXlO3tpgs26dqUN20CqRbepGwA
// aDXxy2D2lD4EHm1pDHU5SXngUECO1fyb9ErG4qeIdE4cfQ8h6ZhRblQNAoGBAOiv
// QvVNvfEm7Qz1BnCuwBTJfkScTJqLuel525c5PIj68w/B6D/cbYaKXr5kSFKroht9
// 8hr4lSHawLCAiBHc4mzXmyI1enxTH1+Sck2nzLxOWCqhExzmpNoKLI0T+Xhy0XH1
// dBQuD7QBdWsP4Y+shBdw/ehVPRqmwqY94Lgr5HQVAoGBAKWNhaJ5QK/hIKlmBAaZ
// k8qF1qqGAzdWdIdDMbn3/mH7YFY2wPeO/7EIl+Gv9/SZ/Dd7m4lEo+UbvDmWxjgF
// eHhuyZgtOz/AAk84uFcGmtE7E0tJKADLahVyt/LjkJ9ENNexjIlp3BCQ1Y2Xz6ZN
// UOIMmeAe65F4BLygeZQeBrk9AoGBAOPyxmrgBUMY+kOmSu/bEkuK9Ysrf5QrbC8A
// 9RHZvacICVQXh3oAbL/QEG7+eSecAsxh/utTOW4YCose765oMN2l/tFtiJgBKowL
// QLU4vMaBDbh9Yeb/QOJl8y0mM1A/U1YLuvMGCNY0U55VyYhh3mnEhMm1r43LbodD
// uUFTppPdAoGBAMHfXeOBXNGaCJf9wAF+qkfTG17z+rplO1ny9vX8p31i3F3wMJG0
// rDkznuYCa6PqBeav2UjbsfeqkynBggKM0wqMUcrSSCzKJEpaOjdaq1fw92Wf4tgq
// aJWMolGAClXedrb1jV2zZgl9xwi8kG1Y9EL4uVuCz2k03dFGuznMTW2v
// -----END RSA PRIVATE KEY-----
// `
// )

type Drone struct {
	commands [10]Command
	name     string
	launch   int // Seconds before launch
	vehID    uint32
	active   bool
}

type Waypoint struct {
	alt float32
	lat float64
	lon float64
}

type Command struct {
	gps         Waypoint
	vel         float32
	action      int
	actionParam int
}

const (
	LAND   = 1
	GOTO   = 2
	LINGER = 3
	FIRE   = 4
	BOMB   = 5
)

const DRONE_MAX = 250
const WAYPOINT_MAX = 10

var m = sync.Mutex{}

// This is not ideal
var ourAssets = []int32{158, 159, 161, 162, 163, 164, 165, 168, 169, 170, 171, 172, 173, 174, 175, 194, 195, 196, 197, 198, 199, 200, 257, 258, 259, 260, 261, 262, 263, 264, 265, 266, 267, 268, 269, 270, 271, 272, 273, 274, 275, 276, 277, 278, 279, 284, 285, 286, 288, 289, 290, 291, 306, 307, 308, 309, 310, 317, 327, 328, 332, 335, 337, 338, 339, 340, 341, 342, 343, 344, 350, 358, 359, 364, 365, 2095, 2096, 2097, 2098, 2111, 2112, 2113, 2114, 2121, 2122, 2149, 2150, 2151, 2152, 2153, 2154, 2155, 2156, 2166, 2171, 2172, 2174, 2175, 2176, 2177, 2178, 2179, 2180, 2181, 2182, 2183, 2184, 2185, 2186, 2187, 2188, 2189, 2190, 2191, 2192, 2193, 2194, 2195, 2196, 2197, 2198, 2199, 2232, 2233, 2264, 2265, 2266, 2268, 2269, 2270, 2271, 2272, 2273, 2280, 2288, 2367, 2368, 2369, 2370, 2371, 2372, 2373, 2374, 2375, 2376, 2377, 2378, 2379, 2380, 2381, 2382, 2383, 2384, 2385, 2386, 2387, 2388, 2389, 2390, 2391, 2392, 2393, 2394, 2395, 2396, 2397, 2398, 2399, 2400, 2401, 2402, 2403, 2404, 2405, 2406, 2407, 2408, 2409, 2410, 2411, 2412, 2413, 2414, 2415, 2416, 2417, 2418, 2419, 2420, 2421, 2422, 2423, 2424, 2425, 2426, 2427, 2428, 2429, 2430, 2431, 2432, 2433, 2434, 2435, 2436, 2437, 2438, 2439, 2440, 2441, 2442, 2443, 2444, 2445, 2446, 2447, 2448, 2449, 2450, 2451, 2452, 2453, 2454, 2455, 2456, 2457, 2458, 2459, 2460, 2461, 2462, 2463, 2464, 2465, 2466, 2467, 2468, 2469, 2470, 2471, 2472, 2473, 2474, 2475, 2476, 2477, 2478, 2479, 2480, 2481, 2482, 2483, 2484, 2485, 2486, 2487, 2488, 2489, 2490, 2491, 2492, 2493, 2494, 2495, 2496, 2497, 2498, 2499, 2500, 2501, 2502, 2503, 2504, 2505, 2506, 2507, 2508, 2509, 2510, 2511, 2512, 2513, 2514, 2515, 2516, 2541, 2625, 2626, 2627, 2633, 2634, 2635, 2636, 2637, 2638, 2639, 2640, 2641, 2642, 2643, 2644, 2645, 2646, 2647, 2648, 2649, 2650, 2651, 2652, 2653, 2654, 2655, 2680, 2681}

func main() {
	var query icarus.QueryPackage

	if len(os.Args) < 2 {
		fmt.Println("Missing arguments")
		os.Exit(1)
	} else if len(os.Args) == 2 {
		query = icarus.NewQuery("127.0.0.1", "9443")
	} else if len(os.Args) == 3 {
		query = icarus.NewQuery(os.Args[2], "9443")
	} else {
		fmt.Println("Too many arguments")
		os.Exit(1)
	}

	query = icarus.NewQuery("127.0.0.1", "9443")
	resp, ok := query.Authenticate("test1", "testing")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}

	// Delay start
	// start := false
	// for !start {
	// 	time.Sleep(10 * time.Second)
	// 	fmt.Println("Hour:", time.Now().Hour(), "Minute:", time.Now().Minute())
	// 	if time.Now().Hour() == 2 && time.Now().Minute() == 24 {
	// 		start = true
	// 	}
	// }

	// openFile, err := os.Open("assets.csv") // IF IT NEEDS TO ADD VEHICLES
	// if err != nil {
	// 	log.Fatal("Unable to open CSV:", err.Error())
	// }

	// reader := csv.NewReader(openFile)
	// reader.Comment = '#'
	// reader.Read()
	// for {
	// 	record, err := reader.Read()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatal("Error reading CSV:", err.Error())
	// 	}

	// 	secMode, _ := strconv.Atoi(record[6])
	// 	secMap := make([]string, 8)
	// 	secMap[0] = ""
	// 	secMap[1] = record[7]
	// 	secMap[2] = record[8]
	// 	secMap[3] = record[9]
	// 	secMap[4] = record[10]
	// 	secMap[5] = record[11]
	// 	secMap[6] = record[12]
	// 	secMap[7] = record[13]

	// 	// Add new vehicle to IcarusServer
	// 	addSeq := query.AddNewVehicle(record[2], record[3], record[1], record[4], secMode, secMap, []byte(defaultCert), []byte(defaultKey), 1, icarus.DefaultC3poTime, icarus.DefaultDaedalusTime)
	// 	responseChan, _ := query.Execute()
	// 	response := <-responseChan
	// 	fmt.Println(response.Get(addSeq))
	// }
	// fmt.Println("All vehicles added.")

	openFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal("Unable to open CSV:", err.Error())
	}

	var drone [DRONE_MAX]Drone

	reader := csv.NewReader(openFile)
	reader.Comment = '#'
	reader.Read()
	reader.Read()

	for i := 0; i < DRONE_MAX; i++ {

		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error reading CSV:", err.Error())
		}

		for j, value := range record {

			if value == "" {
				continue
			}

			if j == 0 {
				drone[i].name = value
				continue
			}

			if j == 1 {
				if value == "" {
					drone[i].launch = 0
				} else {
					drone[i].launch, _ = strconv.Atoi(value)
				}
				continue
			}

			if (j-2)%6 == 0 && record[j] != "" && record[j+1] != "" {

				// Has at least one waypoint
				drone[i].active = true
				cmd := &drone[i].commands[(j-2)/6]

				// Altitude
				if record[j+2] != "" {
					flt64, _ := strconv.ParseFloat(record[j+2], 32)
					cmd.gps.alt = float32(flt64)
				} else {
					cmd.gps.alt = float32(500 + i*10)
				}

				// Velocity
				if record[j+3] != "" {
					flt64, _ := strconv.ParseFloat(record[j+3], 32)
					cmd.vel = float32(flt64)

					switch s := drone[i].name; s[0] {
					case 'G':
						if cmd.vel > 165 {
							cmd.vel = 165
						}
					case 'K':
						if cmd.vel > 120 {
							cmd.vel = 120
						}
					case 'A':
						if cmd.vel > 90 {
							cmd.vel = 90
						}
					case 'T':
						if cmd.vel > 105 {
							cmd.vel = 105
						}
					case 'S':
						if cmd.vel > 60 {
							cmd.vel = 60
						}
					}
				} else {
					switch s := drone[i].name; s[0] {
					case 'G':
						cmd.vel = 165
					case 'K':
						cmd.vel = 120
					case 'A':
						cmd.vel = 90
					case 'T':
						cmd.vel = 105
					case 'S':
						cmd.vel = 60
					}
				}

				cmd.gps.lat, _ = strconv.ParseFloat(record[j], 64)
				cmd.gps.lon, _ = strconv.ParseFloat(record[j+1], 64)

				// Action
				if record[j+4] == "LAND" {
					cmd.action = LAND
					fmt.Println("Added", drone[i].name, "waypoint to land at", cmd.gps.lat, cmd.gps.alt, cmd.gps.lon, "at speed", cmd.vel)
					break
				} else if record[j+4] == "" || record[j+4] == "GOTO" {
					cmd.action = GOTO
					cmd.actionParam = 0
					fmt.Println("Added", drone[i].name, "waypoint to", cmd.gps.lat, cmd.gps.alt, cmd.gps.lon, "at speed", cmd.vel)
				} else if record[j+4] == "FIRE" {
					cmd.action = FIRE
					cmd.actionParam, _ = strconv.Atoi(record[j+5])
					fmt.Println("Added", drone[i].name, "waypoint to", cmd.gps.lat, cmd.gps.alt, cmd.gps.lon, "at speed", cmd.vel, "to fire and linger for", cmd.actionParam, "seconds")
				} else if record[j+4] == "BOMB" {
					cmd.action = BOMB
					cmd.actionParam, _ = strconv.Atoi(record[j+5])
					fmt.Println("Added", drone[i].name, "waypoint to", cmd.gps.lat, cmd.gps.alt, cmd.gps.lon, "at speed", cmd.vel, "to bomb and linger for", cmd.actionParam, "seconds")
				} else {
					cmd.action = LINGER
					cmd.actionParam, _ = strconv.Atoi(record[j+5])
					fmt.Println("Added", drone[i].name, "waypoint to", cmd.gps.lat, cmd.gps.alt, cmd.gps.lon, "at speed", cmd.vel, "to linger for", cmd.actionParam, "seconds")

				}
			}
		}
	}

	// Get active IDs
	statSeq := query.GetAllVehicleStatus()
	responseChan, _ := query.Execute()
	response := <-responseChan
	statusResponse, _ := response.Get(statSeq)

	for _, v := range statusResponse.Vehicles {
		for i := 0; i < DRONE_MAX; i++ {
			if drone[i].active && drone[i].name == v.VehicleCallsign {
				drone[i].vehID = v.VehicleId
				fmt.Println(v.VehicleCallsign, "with", v.PayStatus[5].Resources, "fuel gets ID", drone[i].vehID)
				break
			}
		}
	}

	var done [DRONE_MAX]bool

	for i := 0; i < DRONE_MAX; i++ {
		done[i] = false
	}

	for i := 0; i < DRONE_MAX; i++ {

		if drone[i].active {
			go func(drone Drone, index int) {
				fmt.Printf("Sending %s after %d seconds\n", drone.name, drone.launch)
				time.Sleep(time.Duration(drone.launch) * time.Second)
				fmt.Printf("Sending %s now\n", drone.name)

				m.Lock()
				launchSeq := query.SetNavMode(int(drone.vehID), icarus.NAVIGATION)
				responseChan, _ = query.Execute()
				response := <-responseChan
				_, ok = response.Get(launchSeq)
				if !ok {
					fmt.Println("Error:", response)
				}

				query.ClearQueries()
				m.Unlock()

				// Loop through each waypoint
				for j := 0; j < WAYPOINT_MAX && drone.commands[j].action != 0; j++ {

					m.Lock()
					cmd := icarus.AddCmd(nil, icarus.GOTO, drone.commands[j].gps.lat, drone.commands[j].gps.lon, drone.commands[j].gps.alt, drone.commands[j].vel, 0, 0, 0)

					gotoSeq := query.Goto(int(drone.vehID), cmd)
					responseChan, _ = query.Execute()
					response = <-responseChan
					_, ok = response.Get(gotoSeq)
					if !ok {
						fmt.Println("Error:", response)
					}

					query.ClearQueries()
					m.Unlock()

					var inRange bool = false
					// Perform vehicle action
					switch drone.commands[j].action {
					case LAND:
						for !inRange {
							m.Lock()
							statSeq := query.GetVehicleStatus(int(drone.vehID))
							responseChan, _ := query.Execute()
							response := <-responseChan
							statusResponse, ok := response.Get(statSeq)
							if !ok {
								fmt.Println("Error:", response)
							}
							v := statusResponse.Vehicles[0]
							latDist := math.Pow(drone.commands[j].gps.lat-v.Telem.Latitude, 2)
							lonDist := math.Pow(drone.commands[j].gps.lon-v.Telem.Longitude, 2)
							// altDist := math.Pow(float64(drone.commands[j].gps.alt-v.Telem.Altitude), 2)
							distance := math.Sqrt(latDist+lonDist) * 50000

							query.ClearQueries()
							m.Unlock()

							if distance > 3000 {
								time.Sleep(5 * time.Second)
							} else if distance > 1000 {
								m.Lock()
								homeCmd := icarus.AddCmd(nil, icarus.GOTO, drone.commands[j].gps.lat, drone.commands[j].gps.lon, drone.commands[j].gps.alt, 60, 0, 0, 0)
								gotoSeq = query.Goto(int(drone.vehID), homeCmd)
								responseChan, _ = query.Execute()
								response = <-responseChan
								_, ok = response.Get(gotoSeq)
								if !ok {
									fmt.Println("Error:", response)
								}
								query.ClearQueries()
								m.Unlock()
								time.Sleep(1 * time.Second)
							} else if distance > 500 {
								m.Lock()
								homeCmd := icarus.AddCmd(nil, icarus.GOTO, drone.commands[j].gps.lat, drone.commands[j].gps.lon, drone.commands[j].gps.alt, 30, 0, 0, 0)
								gotoSeq = query.Goto(int(drone.vehID), homeCmd)
								responseChan, _ = query.Execute()
								response = <-responseChan
								_, ok = response.Get(gotoSeq)
								if !ok {
									fmt.Println("Error:", response)
								}
								query.ClearQueries()
								m.Unlock()
								time.Sleep(1 * time.Second)
							} else if distance > 20 {
								m.Lock()
								homeCmd := icarus.AddCmd(nil, icarus.GOTO, drone.commands[j].gps.lat, drone.commands[j].gps.lon, drone.commands[j].gps.alt, 10, 0, 0, 0)
								gotoSeq = query.Goto(int(drone.vehID), homeCmd)
								responseChan, _ = query.Execute()
								response = <-responseChan
								_, ok = response.Get(gotoSeq)
								if !ok {
									fmt.Println("Error:", response)
								}
								query.ClearQueries()
								m.Unlock()
								time.Sleep(100 * time.Millisecond)
							} else {
								inRange = true
							}
						}

						// Land
						m.Lock()
						landSeq := query.SetNavMode(int(drone.vehID), icarus.LAND_NOW)
						responseChan, _ = query.Execute()
						response = <-responseChan
						_, ok = response.Get(landSeq)
						if !ok {
							fmt.Println("No response")
						}

						query.ClearQueries()
						m.Unlock()
						fmt.Print((time.Now()).Format("15:04:05"), " ")
						fmt.Printf("Landing vehicle %s (%d).\n", drone.name, drone.vehID)
						break
					case GOTO:
						fallthrough
					case LINGER:
						for !inRange {
							m.Lock()
							statSeq := query.GetVehicleStatus(int(drone.vehID))
							responseChan, _ := query.Execute()
							response := <-responseChan
							statusResponse, ok := response.Get(statSeq)
							if !ok {
								fmt.Println("Error:", response)
							}
							v := statusResponse.Vehicles[0]

							query.ClearQueries()
							m.Unlock()

							latDist := math.Pow(drone.commands[j].gps.lat-v.Telem.Latitude, 2)
							lonDist := math.Pow(drone.commands[j].gps.lon-v.Telem.Longitude, 2)
							altDist := math.Pow(float64(drone.commands[j].gps.alt-v.Telem.Altitude), 2)
							distance := math.Sqrt(latDist+lonDist+altDist) * 50000

							if distance > 2000 {
								time.Sleep(5 * time.Second)
							} else if distance > 1000 {
								time.Sleep(1 * time.Second)
							} else if distance > 200 {
								time.Sleep(100 * time.Millisecond)
							} else {
								inRange = true
							}
						}

						if drone.commands[j].action == LINGER {
							time.Sleep(time.Duration(drone.commands[j].actionParam) * time.Second)
						}
						break
					case BOMB:
						fallthrough
					case FIRE:
						for !inRange {
							m.Lock()
							statSeq := query.GetVehicleStatus(int(drone.vehID))
							responseChan, _ := query.Execute()
							response := <-responseChan
							statusResponse, ok := response.Get(statSeq)
							if !ok {
								fmt.Println("Error:", response)
							}
							v := statusResponse.Vehicles[0]
							query.ClearQueries()
							m.Unlock()

							foundTarget := false
							if v.Available && (v.PayStatus[9].Resources > 0 || v.PayStatus[12].Resources > 0 || v.PayStatus[11].Resources > 0) {

								var radarInfo icarus.PayloadStatus

								if v.PayStatus[9].Resources > 0 { //AirRadar
									radarInfo = v.PayStatus[icarus.AirRadar]
								} else if v.PayStatus[12].Resources > 0 { //GroundRadar
									radarInfo = v.PayStatus[icarus.GroundRadar]
								} else if v.PayStatus[11].Resources > 0 { //AllRadar
									radarInfo = v.PayStatus[icarus.AllRadar]
								}

								if len(radarInfo.Radar) > 0 {
									for asset, radarPing := range radarInfo.Radar {
										//ID 0 means nothing was returned by the radar payload
										if asset != 0 {
											// var foundTarget bool = false
											for _, ourAsset := range ourAssets {
												if ourAsset == asset {
													foundTarget = true
													// fmt.Printf("Radar found target (%d) type (%d) at (%f,%f,%f)\n", asset, radarPing.Type, radarPing.Latitude, radarPing.Longitude, radarPing.Altitude)
													distance := math.Abs(v.Telem.Latitude-radarPing.Latitude) + math.Abs(v.Telem.Longitude-radarPing.Longitude) + math.Abs(float64(v.Telem.Altitude)-float64(radarPing.Altitude))

													if v.PayStatus[10].Resources > 0 && distance <= 7500 {
														m.Lock()
														configs := query.ExecutePayload(int(drone.vehID), icarus.PayloadType(10), 1, nil, int(asset))
														responseChan, _ = query.Execute()
														response = <-responseChan
														_, ok = response.Get(configs)
														if !ok {
															fmt.Print((time.Now()).Format("15:04:05"), " ")
															fmt.Println("No response")
															os.Exit(1)
														}
														query.ClearQueries()
														m.Unlock()

														fmt.Print((time.Now()).Format("15:04:05"), " ")
														fmt.Printf("Vehicle %s (%d) fired payload %s\n", os.Args[1], drone.vehID, icarus.PayloadType(10).String())

													}
												}
											}
										}
									}
								}
							}

							latDist := math.Pow(drone.commands[j].gps.lat-v.Telem.Latitude, 2)
							lonDist := math.Pow(drone.commands[j].gps.lon-v.Telem.Longitude, 2)
							altDist := math.Pow(float64(drone.commands[j].gps.alt-v.Telem.Altitude), 2)
							distance := math.Sqrt(latDist+lonDist+altDist) * 50000

							if distance > 2000 && !foundTarget {
								time.Sleep(5 * time.Second)
							} else if distance > 1000 && !foundTarget {
								time.Sleep(1 * time.Second)
							} else if distance > 100 {
								time.Sleep(100 * time.Millisecond)
							} else {
								inRange = true
							}
						}

						// CAMERA
						if drone.name[0] == 'K' {
							configs := query.ExecutePayload(int(drone.vehID), icarus.PayloadType(4), 1, nil, 0)
							responseChan, _ = query.Execute()
							response = <-responseChan
							_, ok = response.Get(configs)
							if !ok {
								fmt.Print((time.Now()).Format("15:04:05"), " ")
								fmt.Println("No response")
								os.Exit(1)
							}

							fmt.Printf("Vehicle %s (%d) fired payload %s\n", os.Args[1], drone.vehID, icarus.PayloadType(4).String())
						}

						// BOMB
						if drone.commands[j].action == BOMB {
							m.Lock()
							statSeq := query.GetVehicleStatus(int(drone.vehID))
							responseChan, _ := query.Execute()
							response := <-responseChan
							statusResponse, ok := response.Get(statSeq)
							if !ok {
								fmt.Println("Error:", response)
							}
							v := statusResponse.Vehicles[0]
							query.ClearQueries()

							bombCount := v.PayStatus[3].Resources
							for bombCount > 0 {

								configs := query.ExecutePayload(int(drone.vehID), icarus.PayloadType(3), 1, nil, 0)
								responseChan, _ = query.Execute()
								response = <-responseChan
								_, ok = response.Get(configs)
								if !ok {
									fmt.Print((time.Now()).Format("15:04:05"), " ")
									fmt.Println("No response")
									os.Exit(1)
								}
								query.ClearQueries()
								m.Unlock()

								bombCount--
								fmt.Print((time.Now()).Format("15:04:05"), " ")
								fmt.Printf("Vehicle %s (%d) fired payload %s\n", os.Args[1], drone.vehID, icarus.PayloadType(10).String())
							}
						}

						// LINGER
						for s := 0; s < drone.commands[j].actionParam; s++ {
							time.Sleep(1 * time.Second)

							m.Lock()
							statSeq := query.GetVehicleStatus(int(drone.vehID))
							responseChan, _ := query.Execute()
							response := <-responseChan
							statusResponse, ok := response.Get(statSeq)
							if !ok {
								fmt.Println("Error:", response)
							}
							v := statusResponse.Vehicles[0]
							query.ClearQueries()
							m.Unlock()

							if v.Available && (v.PayStatus[9].Resources > 0 || v.PayStatus[12].Resources > 0 || v.PayStatus[11].Resources > 0) {

								var radarInfo icarus.PayloadStatus

								if v.PayStatus[9].Resources > 0 { //AirRadar
									radarInfo = v.PayStatus[icarus.AirRadar]
								} else if v.PayStatus[12].Resources > 0 { //GroundRadar
									radarInfo = v.PayStatus[icarus.GroundRadar]
								} else if v.PayStatus[11].Resources > 0 { //AllRadar
									radarInfo = v.PayStatus[icarus.AllRadar]
								}

								if len(radarInfo.Radar) > 0 {
									for asset, radarPing := range radarInfo.Radar {
										//ID 0 means nothing was returned by the radar payload
										if asset != 0 {
											// var foundTarget bool = false
											for _, ourAsset := range ourAssets {
												if ourAsset == asset {
													fmt.Printf("Radar found target (%d) type (%d) at (%f,%f,%f)\n", asset, radarPing.Type, radarPing.Latitude, radarPing.Longitude, radarPing.Altitude)
													distance := math.Abs(v.Telem.Latitude-radarPing.Latitude) + math.Abs(v.Telem.Longitude-radarPing.Longitude) + math.Abs(float64(v.Telem.Altitude)-float64(radarPing.Altitude))
													if v.PayStatus[10].Resources > 0 && distance <= 7500 {
														m.Lock()
														configs := query.ExecutePayload(int(drone.vehID), icarus.PayloadType(10), 1, nil, int(asset))
														responseChan, _ = query.Execute()
														response = <-responseChan
														_, ok = response.Get(configs)
														if !ok {
															fmt.Print((time.Now()).Format("15:04:05"), " ")
															fmt.Println("No response")
															os.Exit(1)
														}
														query.ClearQueries()
														m.Unlock()

														fmt.Print((time.Now()).Format("15:04:05"), " ")
														fmt.Printf("Vehicle %s (%d) fired payload %s\n", os.Args[1], drone.vehID, icarus.PayloadType(10).String())
													}
												}
											}
										}
									}
								}
							}

						}
						break

					default:
						break
					}

				}
				fmt.Print((time.Now()).Format("15:04:05"), " ")

				fmt.Println("Drone", drone.name, "done")
				done[index] = true

			}(drone[i], i)
		}
	}

	// Keep looping until all drones are done
	allDone := false
	for !allDone {
		allDone = true
		time.Sleep(1 * time.Minute)
		for j := 0; j < DRONE_MAX; j++ {
			if drone[j].active && !done[j] {
				allDone = false
			}
		}
	}
	fmt.Print((time.Now()).Format("15:04:05"), " ")
	fmt.Println("Route completed")
}
