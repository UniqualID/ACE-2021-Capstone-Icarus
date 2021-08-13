package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

var infraTrans = map[string]string{
	"1":  "Airbase",
	"2":  "Depot",
	"3":  "Datacenter",
	"4":  "Fort",
	"5":  "WMD-PRODUCTION",
	"6":  "ENCAMPMENT",
	"7":  "AnimalHospital",
	"8":  "EMBASSY",
	"9":  "NETWORK-NODE",
	"10": "INSTALLATION",
	"11": "MYSTERY",
}

var vehTrans = map[string]string{
	"1":  "FIGHTER",
	"2":  "BOMBER",
	"4":  "ISR",
	"5":  "ROVER",
	"6":  "MULTI",
	"7":  "WMD",
	"8":  "SAM",
	"9":  "RADAR",
	"10": "LOGISTICS",
	"11": "SPEC-OPS",
}

func BinarySearch(a []int, x int) int {
	r := -1 // not found
	start := 0
	end := len(a) - 1
	for start <= end {
		mid := (start + end) / 2
		if a[mid] == x {
			r = mid // found
			break
		} else if a[mid] < x {
			start = mid + 1
		} else if a[mid] > x {
			end = mid - 1
		}
	}
	return r
}

type entityInfo struct {
	TypeEnt   string
	TypeName  string
	Loc       string
	TimeStamp string
	timeobj   time.Time
}

type jsonInfo struct {
	Vehicle entityInfo
	Infra   entityInfo
}

func main() {
	// ourAssets := []int{158, 159, 161, 162, 163, 164, 165, 168, 169, 170, 171, 172, 173, 174, 175, 194, 195, 196, 197, 198, 199, 200, 257, 258, 284, 285, 286, 288, 289, 290, 291, 306, 307, 308, 309, 310, 317, 327, 328, 332, 335, 337, 338, 339, 340, 341, 342, 343, 344, 350, 358, 359, 364, 365, 2095, 2096, 2097, 2098, 2111, 2112, 2113, 2114, 2121, 2122, 2149, 2150, 2151, 2152, 2153, 2154, 2155, 2156, 2171, 2172, 2174, 2175, 2176, 2177, 2178, 2179, 2180, 2181, 2182, 2183, 2184, 2185, 2186, 2187, 2188, 2189, 2190, 2191, 2192, 2193, 2194, 2195, 2196, 2197, 2198, 2199, 2232, 2233, 2264, 2265, 2266, 2268, 2269, 2270, 2271, 2272, 2273, 2280, 2288, 2541}
	ourAssets := []int{158, 159, 161, 162, 163, 164, 165, 168, 169, 170, 171, 172, 173, 174, 175, 194, 195, 196, 197, 198, 199, 200, 257, 258, 259, 260, 261, 262, 263, 264, 265, 266, 267, 268, 269, 270, 271, 272, 273, 274, 275, 276, 277, 278, 279, 284, 285, 286, 288, 289, 290, 291, 306, 307, 308, 309, 310, 317, 327, 328, 332, 335, 337, 338, 339, 340, 341, 342, 343, 344, 350, 358, 359, 364, 365, 2095, 2096, 2097, 2098, 2111, 2112, 2113, 2114, 2121, 2122, 2149, 2150, 2151, 2152, 2153, 2154, 2155, 2156, 2166, 2171, 2172, 2174, 2175, 2176, 2177, 2178, 2179, 2180, 2181, 2182, 2183, 2184, 2185, 2186, 2187, 2188, 2189, 2190, 2191, 2192, 2193, 2194, 2195, 2196, 2197, 2198, 2199, 2232, 2233, 2264, 2265, 2266, 2268, 2269, 2270, 2271, 2272, 2273, 2280, 2288, 2367, 2368, 2369, 2370, 2371, 2372, 2373, 2374, 2375, 2376, 2377, 2378, 2379, 2380, 2381, 2382, 2383, 2384, 2385, 2386, 2387, 2388, 2389, 2390, 2391, 2392, 2393, 2394, 2395, 2396, 2397, 2398, 2399, 2400, 2401, 2402, 2403, 2404, 2405, 2406, 2407, 2408, 2409, 2410, 2411, 2412, 2413, 2414, 2415, 2416, 2417, 2418, 2419, 2420, 2421, 2422, 2423, 2424, 2425, 2426, 2427, 2428, 2429, 2430, 2431, 2432, 2433, 2434, 2435, 2436, 2437, 2438, 2439, 2440, 2441, 2442, 2443, 2444, 2445, 2446, 2447, 2448, 2449, 2450, 2451, 2452, 2453, 2454, 2455, 2456, 2457, 2458, 2459, 2460, 2461, 2462, 2463, 2464, 2465, 2466, 2467, 2468, 2469, 2470, 2471, 2472, 2473, 2474, 2475, 2476, 2477, 2478, 2479, 2480, 2481, 2482, 2483, 2484, 2485, 2486, 2487, 2488, 2489, 2490, 2491, 2492, 2493, 2494, 2495, 2496, 2497, 2498, 2499, 2500, 2501, 2502, 2503, 2504, 2505, 2506, 2507, 2508, 2509, 2510, 2511, 2512, 2513, 2514, 2515, 2516, 2541, 2625, 2626, 2627, 2633, 2634, 2635, 2636, 2637, 2638, 2639, 2640, 2641, 2642, 2643, 2644, 2645, 2646, 2647, 2648, 2649, 2650, 2651, 2652, 2653, 2654, 2655}
	// ourAssets := []int{0}
	var vehicleMap = make(map[int32]entityInfo)
	var infraMap = make(map[int32]entityInfo)

	var query = icarus.NewQuery("10.59.144.207", "22")
	resp, ok := query.Authenticate("valinar", "thisPasswordNeedsToWorkPLZ")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}
	statSeq := query.GetAllVehicleStatus()
	// statSeq := query.GetVehicleStatus(83)

	for true {

		// statSeq := query.GetVehicleStatus(138)
		responseChan, _ := query.Execute()
		response := <-responseChan
		statusResponse, ok := response.Get(statSeq)
		if !ok {
			fmt.Println("No response")
		}

		//Output
		fmt.Println("=--------------------------------------=ALL RADAR RESPONSES:=----------------------------------------=")
		for _, v := range statusResponse.Vehicles {
			if v.Available && (v.PayStatus[9].Resources > 0 || v.PayStatus[12].Resources > 0 || v.PayStatus[11].Resources > 0) {

				// fmt.Print("=---= ")
				// fmt.Print(v.VehicleCallsign, " [", v.VehicleId, "] // Team # ", v.VConfig.TeamId, " =-----------------------=\n")

				// fmt.Printf("Coordinates:\t%.4f, %.0f, %.4f\n", v.Telem.Latitude, v.Telem.Altitude, v.Telem.Longitude)
				// fmt.Println("\tPrint radar info:")
				var radarInfo icarus.PayloadStatus

				if v.PayStatus[9].Resources > 0 { //AirRadar
					radarInfo = v.PayStatus[icarus.AirRadar]
				} else if v.PayStatus[12].Resources > 0 { //GroundRadar
					radarInfo = v.PayStatus[icarus.GroundRadar]

				} else if v.PayStatus[11].Resources > 0 { //AllRadar
					radarInfo = v.PayStatus[icarus.AllRadar]
				}

				curTime := time.Now()
				timestamp := curTime.Format(time.RFC3339)
				if len(radarInfo.Radar) > 0 {
					for asset, radarPing := range radarInfo.Radar {
						//ID 0 means nothing was returned by the radar payload
						if asset != 0 {
							if radarPing.StructureType == 0 {
								if BinarySearch(ourAssets, int(asset)) == -1 {
									str := fmt.Sprintf("%f,%f,%f", radarPing.Latitude, radarPing.Longitude, radarPing.Altitude)
									vehicleMap[asset] = entityInfo{strconv.Itoa(int(radarPing.Type)), vehTrans[strconv.Itoa(int(radarPing.Type))], str, timestamp, curTime}
									// fmt.Printf("Radar found vehicle (%d) type (%d) at (%f,%f,%f)\n", asset, radarPing.Type, radarPing.Latitude, radarPing.Longitude, radarPing.Altitude)
								}
							} else {
								if BinarySearch(ourAssets, int(asset)) == -1 {
									str := fmt.Sprintf("%f,%f,%f", radarPing.Latitude, radarPing.Longitude, radarPing.Altitude)
									infraMap[asset] = entityInfo{strconv.Itoa(int(radarPing.Type)), infraTrans[strconv.Itoa(int(radarPing.Type))], str, timestamp, curTime}
									// fmt.Printf("Radar found infrastructure (%d) type (%d) at (%f,%f,%f)\n", asset, radarPing.Type, radarPing.Latitude, radarPing.Longitude, radarPing.Altitude)
								}
							}
						}
					}
				}
			}
		}

		// fmt.Println(len(infraMap), len(vehicleMap))
		for iff, info := range infraMap {
			duration := time.Now().Sub(info.timeobj)
			fmt.Printf("Infrastructure: (%d) type %s(%s)\t\tat\t(%s)\t%s\t\n", iff, info.TypeEnt, infraTrans[info.TypeEnt], info.Loc, duration)
		}

		for iff, info := range vehicleMap {
			duration := time.Now().Sub(info.timeobj)
			fmt.Printf("Vehicle: (%d) type %s(%s)\t\tat\t(%s)\t%s\t\n", iff, info.TypeEnt, vehTrans[info.TypeEnt], info.Loc, duration)
		}

		// data := jsonInfo{vehicle:vehicleMap, infraMap}

		var outputMap = make(map[int32]entityInfo)

		for k, v := range vehicleMap {
			outputMap[k] = v
		}
		for k, v := range infraMap {
			outputMap[k] = v
		}
		file, _ := json.Marshal(outputMap)
		_ = ioutil.WriteFile("radar.json", file, 0644)

		time.Sleep(1 * time.Second)
	}
}
