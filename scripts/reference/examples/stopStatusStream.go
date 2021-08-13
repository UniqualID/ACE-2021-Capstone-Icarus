//stopStatusStream.go
//This file contains an example of stopping an active status stream.
package main

import (
	"fmt"
	icarus "git.ironzone.ace/icarus/icarusClient"
	"time"
)

func main() {
	//Create a new query pointed at the IcarusServer instance
	query := icarus.NewQuery("127.0.0.1", "9443")

	//Start a status stream for vehicle (ID: 0). This gives a continuous status update to the response channel
	statSeq := query.StartStatusStream(0)

	//Authenticate to the server
	resp, ok := query.Authenticate("test1", "testing")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}
	fmt.Println("Authenticated")

	//Uncomment this line to show the JSON query being sent to the server
	//query.ShowQuery()

	//Execute the query and read the responses
	responseChan, stopChan := query.Execute()

	go func() {
		for {
			response := <-responseChan
			statusResponse, ok := response.Get(statSeq)
			if !ok {
				if statusResponse.Type != icarus.StartStream {
					break
				}
				fmt.Println("Status response not found")
			}
			fmt.Println(statusResponse)

			//Uncomment this line to show the JSON response returned by the server
			//response.ShowResponse()

		}
	}()

	//Wait here for 20 seconds to see stream results, then stop the stream
	time.Sleep(20 * time.Second)
	ok = query.ClearQueries()
	if !ok {
		fmt.Println("Unable to clear queries")
		return
	}

	//Uncomment this line to show the JSON query being sent to the server
	//query.ShowQuery()

	query.Execute()
	fmt.Println("Initiating stream stop")
	stopChan <- true
	fmt.Println("Waiting 20 more seconds")
	time.Sleep(20 * time.Second)
}
