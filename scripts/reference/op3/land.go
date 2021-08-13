package main

import (
	"fmt"
	"os"
	"strconv"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Missing/Invalid arguments")
		os.Exit(1)
	}

	var query = icarus.NewQuery("10.59.144.202", "179")
	resp, ok := query.Authenticate("valinar", "gitlabsux;(bloodybreach")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}

	vehID, _ := strconv.Atoi(os.Args[1])
	landSeq := query.SetNavMode(vehID, icarus.LAND_NOW)
	responseChan, _ := query.Execute()
	response := <-responseChan
	_, ok = response.Get(landSeq)
	if !ok {
		fmt.Println("No response")
		os.Exit(1)
	}
	fmt.Printf("Landing vehicle %d.\n", vehID)
}
