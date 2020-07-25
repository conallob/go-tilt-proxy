package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"github.com/alexhowarth/go-tilt"
)


// TODO(conall): Convert http_client and request out of the loops using pointers

// TODO(conall): Convert json_payload into a struct and allocate one struct per tilt

func main() {
	sensors := tilt.NewScanner()
	sensors.Scan(5 * time.Second)
	brewfatherURL := "https://log.brewfather.net/stream?id=EzKG3ra2Yp0oMM"
	http_client := &http.Client{}

	for {
	  for _, t := range sensors.Tilts() {
            // API Docs: https://docs.brewfather.app/integrations/custom-stream
	    json_payload, _ := json.Marshal(map[string]interface{}{
               "name":      	 t.Colour(),
	       // TODO(conall): Add a flag to control temperature units
	       "temp":       	 t.Celsius(),
	       "temp_unit":  	 "C",
	       "gravity":    	 t.Gravity(),
	       "gravity_unit":  "G",
	       "device_source": "github.com/conallob/go-tilt-proxy",
	       // TODO(conall): Set "ext_temp" using data from a Ruuvi
	     })

	    request, err := http.NewRequest(
		"POST", brewfatherURL,
	        bytes.NewBuffer(json_payload),
	    )
	    request.Header.Set("Content-Type", "application/json")

	    response, err := http_client.Do(request)
            if err != nil {
                log.Printf("%v", err)
            }
	    defer response.Body.Close()

	    log.Printf("%v", t.Colour())
	    log.Printf("Temp: %vC", t.Celsius())
	    log.Printf("Gravity: %vG", t.Gravity())
          }

	  time.Sleep(15 * time.Minute)

        }
}
