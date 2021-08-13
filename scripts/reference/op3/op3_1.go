package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

const (
	defaultCert string = `-----BEGIN CERTIFICATE-----
MIIC9TCCAd2gAwIBAgIRAKRQOLvrvmORBJxJKkhCpWgwDQYJKoZIhvcNAQELBQAw
EjEQMA4GA1UEChMHQWNtZSBDbzAeFw0xOTAxMDkxNzQ1NTdaFw0yMDAxMDkxNzQ1
NTdaMBIxEDAOBgNVBAoTB0FjbWUgQ28wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQDGzdsktG2DGiQEt7ce1sWlcSc1QNbpLcRemcrGxJKw2JeWYY42R5Le
+umtLV/xy0+ZIA47iHETj0IFFYjWdixmz5/yHnnnJbz8uKinbk3eTmaR6y+EwSAp
gFjXFYjRgif4wPk0qnkgHaI+TJXn2dbBnpv0cX34aUKwaCa/qh0XEZ2nqmjjeowx
mqpD4etICnaMKdJg2Z+da/YG8ExFnwYpzNS9QdfujAxHJ7DoMhPZnyc/sCmaBq+X
eAZbMHtWFuv/24lA/KyJBmCEQGp2x9tn+HmM89SQOj1yOwOqZB87+rjENhp87rgh
MD4vB93/Mzk48tC1LYlr7cLoBh22tskRAgMBAAGjRjBEMA4GA1UdDwEB/wQEAwIF
oDATBgNVHSUEDDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMA8GA1UdEQQIMAaH
BH8AAAEwDQYJKoZIhvcNAQELBQADggEBAKYVTX6fzjOU63QxdtSs9Ot6GdAqaknQ
baTiRnAsXJuRCDZpRIEQKFECC7tdEBbyCh5FiyjpVqxn+U2/OcdNHPYxdHRevRWM
vmNqxeOjha62Bp/JKoN0WR/NfZ7oSHFQYm5kGxTSt6n/BKovWOHI4je0PHD/YITt
Xw5IgIEPDS2eecE+jrHBU596X1jeHSuk+XdQ9Hmo1WFYE9UK7985Oy77+zaNSKdu
ZbBWN877w5AMdmswrC6HcDCXY0Sb3Noadfl9VcDqD69hq0vWPqAomrbjyyTUQCY3
ru4kYaW0u8509KHlpN6ixQWmGAmVhMWtf0g1kPpzJ9HihRtAYUp5Jj4=
-----END CERTIFICATE-----
`
	defaultKey string = `-----BEGIN RSA PRIVATE KEY-----
MIIEpgIBAAKCAQEAxs3bJLRtgxokBLe3HtbFpXEnNUDW6S3EXpnKxsSSsNiXlmGO
NkeS3vrprS1f8ctPmSAOO4hxE49CBRWI1nYsZs+f8h555yW8/Liop25N3k5mkesv
hMEgKYBY1xWI0YIn+MD5NKp5IB2iPkyV59nWwZ6b9HF9+GlCsGgmv6odFxGdp6po
43qMMZqqQ+HrSAp2jCnSYNmfnWv2BvBMRZ8GKczUvUHX7owMRyew6DIT2Z8nP7Ap
mgavl3gGWzB7Vhbr/9uJQPysiQZghEBqdsfbZ/h5jPPUkDo9cjsDqmQfO/q4xDYa
fO64ITA+Lwfd/zM5OPLQtS2Ja+3C6AYdtrbJEQIDAQABAoIBAQCMQh4TFkyRCzdQ
MME8O7Bz2ZIc6yL0njqFt6EtfPA1XooMKcWom/SN5p5IdNPVBmihEtGXxNpqP08H
wTqqe/M1kdQ5gLDmmGRuNGWgwpyjc9K/rhr3YT2sqgWDsYi2r0o+IP9w3bjZJK8b
nvLAAZuXPKyw2AVU5gaL6N81p/IgHCoQ6GAkE19Bk2EK/vRoSKtmZsmsUY3RrezC
HTSX4R4it2eMXBAxEF6stPK0tkv2BvOItoe5BGfWnd4VG2okVcHFkRPiaIiz60iy
02m6hei2dTKRe/GxaNNDNW8jNqTyJ4ZbUK8MeFVyCmFe+Zu8/LLF6bffj46HE00u
yWovzFMxAoGBANq5gT27iGukc2MgXDc/yLKcwmXTU7hGH7TRBCu11Ms9fgInCsrz
x9SRhI5KvysOVe80uPdtyoMx9+I/hzDDNBj0IrIXlO3tpgs26dqUN20CqRbepGwA
aDXxy2D2lD4EHm1pDHU5SXngUECO1fyb9ErG4qeIdE4cfQ8h6ZhRblQNAoGBAOiv
QvVNvfEm7Qz1BnCuwBTJfkScTJqLuel525c5PIj68w/B6D/cbYaKXr5kSFKroht9
8hr4lSHawLCAiBHc4mzXmyI1enxTH1+Sck2nzLxOWCqhExzmpNoKLI0T+Xhy0XH1
dBQuD7QBdWsP4Y+shBdw/ehVPRqmwqY94Lgr5HQVAoGBAKWNhaJ5QK/hIKlmBAaZ
k8qF1qqGAzdWdIdDMbn3/mH7YFY2wPeO/7EIl+Gv9/SZ/Dd7m4lEo+UbvDmWxjgF
eHhuyZgtOz/AAk84uFcGmtE7E0tJKADLahVyt/LjkJ9ENNexjIlp3BCQ1Y2Xz6ZN
UOIMmeAe65F4BLygeZQeBrk9AoGBAOPyxmrgBUMY+kOmSu/bEkuK9Ysrf5QrbC8A
9RHZvacICVQXh3oAbL/QEG7+eSecAsxh/utTOW4YCose765oMN2l/tFtiJgBKowL
QLU4vMaBDbh9Yeb/QOJl8y0mM1A/U1YLuvMGCNY0U55VyYhh3mnEhMm1r43LbodD
uUFTppPdAoGBAMHfXeOBXNGaCJf9wAF+qkfTG17z+rplO1ny9vX8p31i3F3wMJG0
rDkznuYCa6PqBeav2UjbsfeqkynBggKM0wqMUcrSSCzKJEpaOjdaq1fw92Wf4tgq
aJWMolGAClXedrb1jV2zZgl9xwi8kG1Y9EL4uVuCz2k03dFGuznMTW2v
-----END RSA PRIVATE KEY-----
`
)

type GPS struct {
	lat float64
	lon float64
}

type Sortie struct {
	gps []GPS
}

const MAX_VEL = 120.0
const ALT_DIF = 200

const JOHNNY = 2
const FRANCIS = 1
const BABY = 0

var query = icarus.NewQuery("10.59.144.202", "179")

func main() {
	// Authentication
	resp, ok := query.Authenticate("valinar", "gitlabsux;(bloodybreach")
	if !ok {
		fmt.Print((time.Now()).Format("15:04:05"), " ")
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
		os.Exit(1)
	}

	// Add vehicles from assets.csv
	addVehicles()

	// Print status
	status_all()

	// Get active IDs
	statSeq := query.GetAllVehicleStatus()
	responseChan, _ := query.Execute()
	response := <-responseChan
	statusResponse, ok := response.Get(statSeq)
	if !ok {
		fmt.Print((time.Now()).Format("15:04:05"), " ")
		fmt.Println("No response")
	}
	activeIDs := make([]uint32, 4)
	for i, v := range statusResponse.Vehicles {
		activeIDs[i] = v.VehicleId
		fmt.Print((time.Now()).Format("15:04:05"), " ")
		fmt.Println("Active ID:", v.VehicleId)
	}

	// Get sorties
	var sortie [3]Sortie

	openFile, err := os.Open("sorties.csv")
	if err != nil {
		log.Fatal("Unable to open CSV:", err.Error())
	}

	reader := csv.NewReader(openFile)
	reader.Comment = '#'
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error reading CSV:", err.Error())
		}

		sortieNum, _ := strconv.Atoi(record[0])
		gpsNum, _ := strconv.Atoi(record[1])
		var tempGPS GPS
		tempGPS.lat, _ = strconv.ParseFloat(record[2], 64)
		tempGPS.lon, _ = strconv.ParseFloat(record[3], 64)
		sortie[sortieNum].gps = append(sortie[sortieNum].gps, tempGPS)
		fmt.Print((time.Now()).Format("15:04:05"), " ")
		fmt.Printf("Drone %d - Waypoint %d: %.4f, %.4f\n", sortieNum, gpsNum, sortie[sortieNum].gps[gpsNum].lat, sortie[sortieNum].gps[gpsNum].lon)
	}

	// ASSUMING NO FUEL UP

	// Set bases
	var BASE [3]GPS
	BASE[BABY].lat = 50.3335
	BASE[BABY].lon = -60.6990
	BASE[FRANCIS].lat = 50.5580
	BASE[FRANCIS].lon = -61.0745
	BASE[JOHNNY].lat = 49.6165
	BASE[JOHNNY].lon = -67.7075

	// Launch Drones

	// Wait time for launch
	//fmt.Println((time.Now()).Format("15:04:05"), "Waiting to launch.")
	//time.Sleep(10 * time.Minute)

	// Go BABY
	fmt.Println((time.Now()).Format("15:04:05"), "GO BABY")
	launch(activeIDs[BABY])
	for _, dest := range sortie[BABY].gps {
		//AddCmd(command list, command type, latitude, longitude, altitude, velocity, turn radius, linger time, transit heading)
		cmdList := icarus.AddCmd(nil, icarus.GOTO, dest.lat, dest.lon, ALT_DIF+ALT_DIF*BABY, MAX_VEL, 0, 0, 0)
		gotoSeq := query.Goto(int(activeIDs[BABY]), cmdList)
		responseChan, _ = query.Execute()
		response = <-responseChan
		_, ok = response.Get(gotoSeq)
		if !ok {
			fmt.Print((time.Now()).Format("15:04:05"), " ")
			fmt.Println("No response")
			os.Exit(1)
		}

		fmt.Print((time.Now()).Format("15:04:05"), " ")
		fmt.Printf("Moving vehicle %d to %.4f %d %.4f at %.0f meters/sec.\n", activeIDs[BABY], dest.lat, ALT_DIF+ALT_DIF*BABY, dest.lon, MAX_VEL)

		// Trigger when in range
		var inRange bool = false
		for !inRange {
			time.Sleep(10 * time.Second)
			statSeq := query.GetVehicleStatus(int(activeIDs[BABY]))
			responseChan, _ := query.Execute()
			response := <-responseChan
			statusResponse, ok := response.Get(statSeq)
			if !ok {
				fmt.Print((time.Now()).Format("15:04:05"), " ")
				fmt.Println("No response")
				os.Exit(1)
			}
			v := statusResponse.Vehicles[0]
			if v.Telem.Latitude > dest.lat-.0020 && v.Telem.Latitude < dest.lat+.0020 && v.Telem.Longitude > dest.lon-.0020 && v.Telem.Longitude < dest.lon+.0020 {
				inRange = true
			}
		}

		// Take picture
		configs := query.ExecutePayload(int(activeIDs[BABY]), 4, 1, nil, 0)
		responseChan, _ = query.Execute()
		fmt.Print((time.Now()).Format("15:04:05"), " ")
		fmt.Printf("Vehicle %v taking picture.\n", activeIDs[BABY])
		response = <-responseChan

		configResponse, ok := response.Get(configs)
		if !ok {
			fmt.Print((time.Now()).Format("15:04:05"), " ")
			fmt.Println("Error during picture taken: \n", configResponse.Message)
		} else {
			fmt.Print((time.Now()).Format("15:04:05"), " ")
			fmt.Printf("Image returned by vehicle %v: %v\n", activeIDs[BABY], configResponse.PayloadResponse.File)
		}

		// Check radar
		statSeq := query.GetVehicleStatus(int(activeIDs[BABY]))
		responseChan, _ := query.Execute()
		response := <-responseChan
		statusResponse, ok := response.Get(statSeq)
		if !ok {
			fmt.Println("No response")
			os.Exit(1)
		}
		v := statusResponse.Vehicles[0]

		// Output radar
		fmt.Println("\nPrint radar info:")
		radarInfo := v.PayStatus[icarus.AllRadar]
		//Radar's current range is found in radarInfo.Resources parameter
		if len(radarInfo.Radar) > 0 {
			for asset, radarPing := range radarInfo.Radar {
				//ID 0 means nothing was returned by the radar payload
				if asset != 0 {
					if radarPing.StructureType == 0 {
						fmt.Printf("Radar found vehicle(%d) at (%f,%f,%f)\n", asset, radarPing.Latitude, radarPing.Longitude, radarPing.Altitude)
					} else {
						fmt.Printf("Radar found infrastructure(%d) at (%f,%f,%f)\n", asset, radarPing.Latitude, radarPing.Longitude, radarPing.Altitude)
					}
				}
			}
		}

	}
	// Go back to base
	cmdList := icarus.AddCmd(nil, icarus.GOTO, BASE[BABY].lat, BASE[BABY].lon, ALT_DIF+ALT_DIF*BABY, MAX_VEL, 0, 0, 0)
	gotoSeq := query.Goto(int(activeIDs[BABY]), cmdList)
	responseChan, _ = query.Execute()
	response = <-responseChan
	_, ok = response.Get(gotoSeq)
	if !ok {
		fmt.Print((time.Now()).Format("15:04:05"), " ")
		fmt.Println("No response")
		os.Exit(1)
	}

	fmt.Print((time.Now()).Format("15:04:05"), " ")
	fmt.Printf("Moving vehicle %d to base at %.4f %d %.4f at %.0f meters/sec.\n", activeIDs[BABY], BASE[BABY].lat, ALT_DIF+ALT_DIF*BABY, BASE[BABY].lon, MAX_VEL)

	// Trigger when in range
	var inRange bool = false
	for !inRange {
		time.Sleep(10 * time.Second)
		statSeq := query.GetVehicleStatus(int(activeIDs[BABY]))
		responseChan, _ := query.Execute()
		response := <-responseChan
		statusResponse, ok := response.Get(statSeq)
		if !ok {
			fmt.Print((time.Now()).Format("15:04:05"), " ")
			fmt.Println("No response")
			os.Exit(1)
		}
		v := statusResponse.Vehicles[0]
		if v.Telem.Latitude > BASE[BABY].lat-.0100 && v.Telem.Latitude < BASE[BABY].lat+.0100 && v.Telem.Longitude > BASE[BABY].lon-.0100 && v.Telem.Longitude < BASE[BABY].lon+.0100 {
			inRange = true
		}
	}

	// Land
	landSeq := query.SetNavMode(int(activeIDs[BABY]), icarus.LAND_NOW)
	responseChan, _ = query.Execute()
	response = <-responseChan
	_, ok = response.Get(landSeq)
	if !ok {
		fmt.Print((time.Now()).Format("15:04:05"), " ")
		fmt.Println("No response")
		os.Exit(1)
	}
	fmt.Print((time.Now()).Format("15:04:05"), " ")
	fmt.Printf("Landing vehicle %d.\n", activeIDs[BABY])

	// Print status
	status_all()

	fmt.Println((time.Now()).Format("15:04:05"), "End script.")

}

func launch(vehID uint32) {

	navSeq := query.SetNavMode(int(vehID), icarus.NAVIGATION) // Drones can go straight into navigation mode; may be issue in future
	responseChan, _ := query.Execute()
	response := <-responseChan
	_, ok := response.Get(navSeq)
	if !ok {
		fmt.Print((time.Now()).Format("15:04:05"), " ")
		fmt.Println("No response")
	}
	fmt.Print((time.Now()).Format("15:04:05"), " ")
	fmt.Printf("Launching vehicle %d.\n", vehID)
}

func addVehicles() {
	openFile, err := os.Open("assets.csv")
	if err != nil {
		log.Fatal("Unable to open CSV:", err.Error())
	}

	reader := csv.NewReader(openFile)
	reader.Comment = '#'
	reader.Read()
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error reading CSV:", err.Error())
		}

		secMode, _ := strconv.Atoi(record[6])
		secMap := make([]string, 8)
		secMap[0] = ""
		secMap[1] = record[7]
		secMap[2] = record[8]
		secMap[3] = record[9]
		secMap[4] = record[10]
		secMap[5] = record[11]
		secMap[6] = record[12]
		secMap[7] = record[13]

		// Add new vehicle to IcarusServer
		addSeq := query.AddNewVehicle(record[2], record[3], record[1], record[4], secMode, secMap, []byte(defaultCert), []byte(defaultKey), 1, icarus.DefaultC3poTime, icarus.DefaultDaedalusTime)
		responseChan, _ := query.Execute()
		response := <-responseChan
		fmt.Print((time.Now()).Format("15:04:05"), " ")
		fmt.Println(response.Get(addSeq))
	}
	fmt.Print((time.Now()).Format("15:04:05"), " ")
	fmt.Println("All vehicles added.")
}

func status_all() {
	statSeq := query.GetAllVehicleStatus()
	responseChan, _ := query.Execute()
	response := <-responseChan
	statusResponse, ok := response.Get(statSeq)
	if !ok {
		fmt.Print((time.Now()).Format("15:04:05"), " ")
		fmt.Println("No response")
	}

	//Output
	fmt.Print((time.Now()).Format("15:04:05"), " ")
	fmt.Println("All Vehicles Status:")
	for _, v := range statusResponse.Vehicles {
		fmt.Print("=---= ")
		fmt.Print(v.VehicleCallsign, " [", v.VehicleId, "] // Team # ", v.VConfig.TeamId, " =-----------------------=\n")
		if !v.Available {
			println("DESTROYED\n")
		} else {
			fmt.Printf("Coordinates:\t%.4f, %.0f, %.4f", v.Telem.Latitude, v.Telem.Altitude, v.Telem.Longitude)
			fmt.Print("\nVector:\t\t", v.Telem.Heading, "°") // This is broken
			switch d := v.Telem.Heading; {
			case d < 22.5:
				fmt.Print("N")
			case d < 67.5:
				fmt.Print("NE")
			case d < 112.5:
				fmt.Print("E")
			case d < 157.5:
				fmt.Print("SE")
			case d < 202.5:
				fmt.Print("S")
			case d < 227.5:
				fmt.Print("SW")
			case d < 292.5:
				fmt.Print("W")
			case d < 327.5:
				fmt.Print("NW")
			case d >= 327.5:
				fmt.Print("N")
			default:
			}
			fmt.Printf("\t%.0f m/s", v.Telem.Velocity)
			fmt.Print("\nFuel:\t\t", v.PayStatus[5].Resources)
			// fmt.Print("\nValid:\t\t", v.VehicleType.Valid) Find out what this is
			fmt.Print("\nNav Mode:\t", v.VehicleType.Mode)
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
		}
	}
}
