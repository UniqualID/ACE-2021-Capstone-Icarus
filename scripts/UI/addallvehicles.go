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
MIIDWjCCAkKgAwIBAgIQWucCvQ7YKNkspkqB5e47ITANBgkqhkiG9w0BAQsFADAW
MRQwEgYDVQQDDAtFYXN5LVJTQSBDQTAeFw0yMTA4MTExMzAwNTZaFw0yMzExMTQx
MzAwNTZaMBcxFTATBgNVBAMMDHZlaGljbGUtY2VydDCCASIwDQYJKoZIhvcNAQEB
BQADggEPADCCAQoCggEBAMuWNDkuD3rjD8LsFMGX3114L4rBg1HE5NqxR99t5BJ5
2wHMQ5nM3D40l/Azv7HGMgV+KcBe1BDmqyFwNWHm56eGlkMqj8HvIPSFw71cfDlt
jH+exNA8jNW2Mh1TzbKxpCsY3nCMvYNqVD3BB8dBTvQFX2qbZkfvoG8kyr+XQcIZ
HdMbWP/+fw39LRyIFi9ZZO5Njl0rEWQumzht79mnTyoGnFynHPhcVuKfNDbLh+TA
falmtoA+AGonu+8bG22IEvO8JscgbVu4sA/7Oh32gdWxZRUx9LDetFluTRHDbt1p
YEh6mr4IL2LU3YYC36mco7X+KCdG7QDJUo86QZ8wssUCAwEAAaOBojCBnzAJBgNV
HRMEAjAAMB0GA1UdDgQWBBRqsaM4A8SFLGkYL/35TjQVaMajdzBRBgNVHSMESjBI
gBTuVAS4OHJ4qep+BDi/95tCjTFA3aEapBgwFjEUMBIGA1UEAwwLRWFzeS1SU0Eg
Q0GCFEyGVKe6891GwxhCApXhKYLOE1q3MBMGA1UdJQQMMAoGCCsGAQUFBwMCMAsG
A1UdDwQEAwIHgDANBgkqhkiG9w0BAQsFAAOCAQEACB5tVqvWmDLzD971H0Nr+g5k
Q0GWk0XJjUWV3VdfE6HHL1qOKnjg1RNRKSWJcdxogVZ99ApEC5QFlb0az9J/yXIK
WumgfTl+IA2/U/MjpgvkcjQarHExj/whOO2GkITR4qG5ki66xahMQU86Qkz4kdlE
3VzJv/sYHvHafoqT1daLeMDIyams32XZecs058yj3SHtfudkv9qSHM69pUCKQc+y
7W0BxSMW3A8NLSaa/sW68gGeDXAMCc/CUDfr/DrmiH74q9Jw36zUeRuyiXEoGa6B
mHfkIrWmKfQl4ivd4eEjl4y3Riw0TIM2+rK9nIok+I3HMeZAFc+6chzpy3Vw4w==
-----END CERTIFICATE-----
`
	defaultKey string = `-----BEGIN RSA PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDLljQ5Lg964w/C
7BTBl99deC+KwYNRxOTasUffbeQSedsBzEOZzNw+NJfwM7+xxjIFfinAXtQQ5qsh
cDVh5uenhpZDKo/B7yD0hcO9XHw5bYx/nsTQPIzVtjIdU82ysaQrGN5wjL2DalQ9
wQfHQU70BV9qm2ZH76BvJMq/l0HCGR3TG1j//n8N/S0ciBYvWWTuTY5dKxFkLps4
be/Zp08qBpxcpxz4XFbinzQ2y4fkwH2pZraAPgBqJ7vvGxttiBLzvCbHIG1buLAP
+zod9oHVsWUVMfSw3rRZbk0Rw27daWBIepq+CC9i1N2GAt+pnKO1/ignRu0AyVKP
OkGfMLLFAgMBAAECggEAE/CzTsJvK7cfrRTYd5m7e/kwluObTB//1le01XJ6+0BK
LiPmnyLMk58tHo7DANOLxLP1vOzM0pb1PgIyiFNIb0UkZJG/oNweGzUa1epAFJPh
RgKz/j1gAfKld+/kXtM9ZCc8akOusDdy5tWIQIDqDcaU8OklnHgg+6Hb5eYXv73v
JaX6faA+pAyfrRXPUv29savCrOfC9T+/XQFhadsxv7NJ39k3xSUpRj33xhnKHJew
TErv6eik6AWCIwBmjEd7e030yc9lxj6AQCBhGLNLnqTAOEs+7iq/Z+o2dTUoo+Is
s2Hv1nNWD4fEfYjlzop/G7RtAP+/YpQfemT5EnjCZQKBgQD43p8XBAVsyCnIqLWy
m/f6SH2FOZAe6bpUlUCB/BDofp4UxehK0Vqz00KMlheW88SWtRwOzRqh8BpMf5jY
w1VVwYOq9OuuroYEvqLSL3OJC/bI2Y1XoYSc1jiyERJiEyZscQwE706tn5BTs65o
xRsmC4FFESnYU9eme+HMpxBeewKBgQDRa3J7MLT7m7hdFgA8iiAaY7Xwdzg0rm6C
MMRxr0z0c+54EcbFxkvXAnhznpETgOKAf8m/38iGsNVQyULJcAJ96Z5goZU8KIpt
y5jiuVV0Q41uDF1gQ5HdBQoKNzmcUoEk/kqF/7tw0YMm0ssE0FBKfWZ8NvxnBJ+X
9GEh6wMPvwKBgQCnGqxOAvg1i81qm8WtEUcXujb2Dqmz6BKiGrl+zib4RZSVtDF+
k0ZI+rBzv5BFXpcA7sjRM20PkS3HceHKopYZB+AGGYBrAWqhovOnGA+G1Q623EUo
sx5dRQY6onXqfptNMzbz1U/KCxsY6MxSMw9Ao3cATj7/r4RcmFAMX50BIwKBgFUr
2+2aS6EhHgRN1F7K1m+lKfPYqPVUFAHDD8Ikx6mMFOqkYDFDisixaoerb9l8y6Hz
VvxAaW/OL5OhpxYCBiFriExIq0dCPbqV2WIsekLEzpp5UOi70nEPgZvSBg47x+Zh
vbFt5q/lpe9+P69/gYgqJCpUuqazTK0iaPbAfAhFAoGBAO2wqAAtAFILyN1Pt1CX
TmIfJL9wp9Ss2XuwLrv/1XpLI27/S9rzVbdGsnAiXrspCTLbwno+nnUjcg1yZd4t
e3ffcmJig1KwTrClKMiS+o/0zAm+5ndjOJ8QjiljcCR7qo6j9pAPucgYqpq1D/tW
y4brRIrbe7Blj8yuFjLT744/
-----END RSA PRIVATE KEY-----
`
)

func main() {
	var query icarus.QueryPackage

	if len(os.Args) == 1 {
		fmt.Println("Missing arguments")
		fmt.Println("./addallvehicles [realassets.csv/testassets.csv]")
		os.Exit(1)
	} else if len(os.Args) == 2 {
		query = icarus.NewQuery("10.59.144.207", "22")
	} else if len(os.Args) == 3 {
		query = icarus.NewQuery(os.Args[2], "22")
	} else {
		fmt.Println("Too many arguments")
		os.Exit(1)
	}

	openFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal("Unable to open CSV:", err.Error())
	}

	resp, ok := query.Authenticate("valinar", "thisPasswordNeedsToWorkPLZ")
	if !ok {
		fmt.Print((time.Now()).Format("15:04:05"), " ")
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}

	fmt.Printf("Adding files from %s.\n", os.Args[1])
	reader := csv.NewReader(openFile)
	reader.Comment = '#'
	reader.Read()
	// for i := 0; i < 20; i++ {

	// counter := 1
	var allResponses []uint32
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
		fmt.Println(response.Get(addSeq))
		fmt.Println()
		query.ClearQueries()
		// allResponses = append(allResponses, query.AddNewVehicle(record[2], record[3], record[1], record[4], secMode, secMap, []byte(defaultCert), []byte(defaultKey), 1, icarus.DefaultC3poTime, icarus.DefaultDaedalusTime))
		// fmt.Println("Adding: ", counter)
		// counter += 1
	}

	responseChan, _ := query.Execute()
	response := <-responseChan
	for _, v := range allResponses {

		shit, _ := response.Get(v)
		fmt.Print(shit, "\n")
	}
	fmt.Print((time.Now()).Format("15:04:05"), " ")
	fmt.Println("All vehicles added.")
}
