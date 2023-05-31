package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

func main() {
	telemetryConfig := appinsights.NewTelemetryConfiguration("756af033-cded-43b9-9b97-bb67dd54234d")
	telemetryConfig.MaxBatchSize = 8192
	telemetryConfig.MaxBatchInterval = 5 * time.Second
	aiClient := appinsights.NewTelemetryClientFromConfig(telemetryConfig)
	aiClient.Context().Tags.Cloud().SetRole("homeoffice")

	regions := []string{"westeurope", "northeurope", "germanywestcentral", "eastus"}
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	for _, region := range regions {
		for i := 0; i < 250; i++ {
			url := fmt.Sprintf("http://cloudexperienceday-%s.%s.azurecontainer.io:8080/singlebyte", region, region)
			req, err := http.NewRequest("POST", url, nil)
			if err != nil {
				fmt.Printf("Error creating request to %s: %v\n", url, err)
				continue
			}

			start := time.Now()
			dependency := appinsights.NewRemoteDependencyTelemetry(region, "External", url, true)

			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error sending request to %s: %v\n", url, err)
				continue
			}
			resp.Body.Close()

			dependency.Duration = time.Since(start)
			aiClient.Track(dependency)

			fmt.Printf("Successfully sent request %d to %s, response status code: %d\n", i+1, url, resp.StatusCode)
		}
	}

	type ProxyRequest struct {
		URL  string `json:"url"`
		Name string `json:"name"`
	}

	payload := ProxyRequest{
		URL: "http://cloudexperienceday-eastus.eastus.azurecontainer.io:8080/singlebyte",
		Name: "eastus",
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 250; i++ {
		url := "http://cloudexperienceday-germanywestcentral.germanywestcentral.azurecontainer.io:8080/proxy/eastus"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			fmt.Printf("Error creating request to %s: %v\n", url, err)
			continue
		}

		start := time.Now()
		dependency := appinsights.NewRemoteDependencyTelemetry("proxy", "External", url, true)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending request to %s: %v\n", url, err)
			continue
		}
		resp.Body.Close()

		dependency.Duration = time.Since(start)
		aiClient.Track(dependency)

		fmt.Printf("Successfully sent request %d to %s, response status code: %d\n", i+1, url, resp.StatusCode)
	}

	<-aiClient.Channel().Close(30 * time.Second)
}
