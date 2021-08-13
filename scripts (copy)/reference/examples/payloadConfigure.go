//payloadConfigure.go
//This file contains an example of configuring a payload system
package main

import (
	"fmt"
	icarus "git.ironzone.ace/icarus/icarusClient"
)

func main() {
	//Create a new query pointed at the IcarusServer instance
	query := icarus.NewQuery("127.0.0.1", "9443")

	//Create payload configurations
	configs := icarus.AddPayloadConfig(nil, "payload1", icarus.ThermalLance, 1, true)
	configs = icarus.AddPayloadConfig(configs, "payload2", icarus.AntiMatterMissile, 1, true)

	//Configure Daedalus systems
	configSeq := query.ConfigurePayloads(0, configs)

	//Authenticate to the server
	resp, ok := query.Authenticate("test1", "testing")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}

	//Uncomment this line to show the JSON query being sent to the server
	//query.ShowQuery()

	//Execute the query and read the responses
	responseChan, _ := query.Execute()

	fmt.Println("Waiting for responses:")
	response := <-responseChan
	configResponse, ok := response.Get(configSeq)
	if !ok {
		fmt.Println("Payload configure response not found")
	}
	fmt.Println(configResponse)

	//Uncomment this line to show the JSON response returned by the server
	//response.ShowResponse()

}
