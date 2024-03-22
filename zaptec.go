package main

import (
	"encoding/csv"
	"encoding/json"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const username = "ENTER_HERE_YOUR_USERNAME"
const password = "ENTER_HERE_YOUR_PASSWORD"
const installationId = "ENTER_HERE_YOUR_INSTALLATION_ID"
const fromDate = "2022-12-15T00%3A00%3A00.01"
const toDate = "2023-09-30T23%3A59%3A59.99"
const highFareStartHour = 7     // hour
const highFareEndHour = 19      // hour.
const lowFareNoonStartHour = 12 // hour
const lowFareNoonEndHour = 15   // hour

const highFareCost = 0.36 // CHF / kWh
const lowFareCost = 0.29656 // CHF / kWh

type JSONData struct {
	Pages int `json:"Pages"`
	Data  []struct {
		ID                string  `json:"Id"`
		DeviceID          string  `json:"DeviceId"`
		StartDateTime     string  `json:"StartDateTime"`
		EndDateTime       string  `json:"EndDateTime"`
		Energy            float64 `json:"Energy"`
		CommitMetadata    int     `json:"CommitMetadata"`
		CommitEndDateTime string  `json:"CommitEndDateTime"`
		ChargerID         string  `json:"ChargerId"`
		DeviceName        string  `json:"DeviceName"`
		ExternallyEnded   bool    `json:"ExternallyEnded"`
		EnergyDetails     []struct {
			Timestamp string  `json:"Timestamp"`
			Energy    float64 `json:"Energy"`
		} `json:"EnergyDetails"`
		ChargerFirmwareVer struct {
			Major         int `json:"Major"`
			Minor         int `json:"Minor"`
			Build         int `json:"Build"`
			Revision      int `json:"Revision"`
			MajorRevision int `json:"MajorRevision"`
			MinorRevision int `json:"MinorRevision"`
		} `json:"ChargerFirmwareVersion"`
		SignedSession string `json:"SignedSession"`
	} `json:"Data"`
}

type DeviceEnergyCost struct {
	DeviceName     string
	LowFareEnergy  float64
	HighFareEnergy float64
	LowFareCost    float64
	HighFareCost   float64
	TotalEnergy    float64
	TotalCost      float64
}

func main() {
	var jsonData = GetTimeSeriesData()
	WriteDetailAnalysisCsv(jsonData)
}

func GetTimeSeriesData() JSONData {	
	var bearerToken, err = getToken()

	url := "https://api.zaptec.com/api/chargehistory?InstallationId=" + installationId + "&From=" + fromDate + "&To=" + toDate + "&GroupBy=2&DetailLevel=1&PageSize=1000"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Bearer "+bearerToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	// Unmarshal JSON data into JSONData struct
	var data JSONData
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
	}

	return data
}

// getToken makes an OAuth request and returns the access token.
// It now accepts username and password as arguments.
func getToken() (string, error) {
	// Your API endpoint and method
	url := "https://api.zaptec.com/oauth/token"
	method := "POST"

	// Constructing the payload with the provided username and password
	payloadData := fmt.Sprintf("grant_type=password&username=%s&password=%s", username, password)
	payload := bytes.NewReader([]byte(payloadData))

	// Creating the request
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", err
	}

	// Adding necessary headers
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Making the request
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Reading the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	// A struct to match the JSON response format
	var result struct {
		AccessToken string `json:"access_token"`
	}

	// Unmarshalling the body into the struct
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	// Returning the access token
	return result.AccessToken, nil
}

func WriteDetailAnalysisCsv(jsonData JSONData) {
	// Open CSV file for writing
	csvFile, err := os.Create("detailData.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer csvFile.Close()

	// Create CSV writer
	writer := csv.NewWriter(csvFile)

	// Write header row
	header := []string{"Timestamp", "DeviceName", "Energy", "Cost", "highFare"}
	err = writer.Write(header)
	if err != nil {
		fmt.Println(err)
		return
	}

	totalHighFareEnergy := 0.0
	totalLowFareEnergy := 0.0
	totalHighFareCost := 0.0
	totalLowFareCost := 0.0

	// Create a map to keep track of the energy and cost for each device
	deviceEnergyCostMap := make(map[string]DeviceEnergyCost)

	// write the data to the CSV file
	for _, d := range jsonData.Data {
		for _, ed := range d.EnergyDetails {
			var convertedTime = ConvertAndFixTimeZone(ed.Timestamp)
			var highFare = CheckIfHighFare(convertedTime)
			var energy = ed.Energy
			var fareCost = CalculateFare(highFare, energy)
			if highFare == false {
				totalLowFareEnergy += energy
				totalLowFareCost += fareCost
			} else {
				totalHighFareEnergy += energy
				totalHighFareCost += fareCost
			}

			row := []string{convertedTime.Format(time.RFC3339), d.DeviceName, fmt.Sprintf("%v", ed.Energy), fmt.Sprintf("%v", fareCost), strconv.FormatBool(highFare)}

			// Update device energy and cost
			deviceEnergyCost, ok := deviceEnergyCostMap[d.DeviceName]
			if !ok {
				deviceEnergyCost = DeviceEnergyCost{}
			}

			if highFare {
				deviceEnergyCost.HighFareEnergy += energy
				deviceEnergyCost.HighFareCost += fareCost
			} else {
				deviceEnergyCost.LowFareEnergy += energy
				deviceEnergyCost.LowFareCost += fareCost
			}

			deviceEnergyCost.TotalEnergy += energy
			deviceEnergyCost.TotalCost += fareCost

			deviceEnergyCostMap[d.DeviceName] = deviceEnergyCost
			// write the row to the CSV file
			err := writer.Write(row)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	// flush the writer to ensure all data is written
	writer.Flush()

	// Print out the device energy and cost information
	fmt.Printf("From Date: \t%s\n", fromDate)
	fmt.Printf("To Date:   \t%s\n", toDate)
	fmt.Printf("Hochtarif Kosten: \t%8.2f CHF\n", highFareCost)
	fmt.Printf("Niedertarif Kosten::  \t%8.2f CHF\n", lowFareCost)
	fmt.Printf("Total Hochtarif Energie:\t%8.2f kWh\n", totalHighFareEnergy)
	fmt.Printf("Total Niedertarif Energie:\t%8.2f kWh\n", totalLowFareEnergy)
	fmt.Printf("Total Energie:   \t%8.2f kWh\n", totalLowFareEnergy+totalHighFareEnergy)
	fmt.Printf("Total Hochtarif Kosten:\t%8.2f CHF\n", totalHighFareCost)
	fmt.Printf("Total Niedertarif Kosten::\t%8.2f CHF\n", totalLowFareCost)
	fmt.Printf("Total Kosten:     \t%8.2f CHF\n", totalLowFareCost+totalHighFareCost)
	for deviceName, deviceEnergyCost := range deviceEnergyCostMap {
		fmt.Printf("Device Name: %s\t", fmt.Sprintf("%-15s", deviceName)[:15])
		fmt.Printf("Total Energie: %8.2f kWh\t", deviceEnergyCost.TotalEnergy)
		fmt.Printf("Hochtarif Energie: %8.2f kWh\t", deviceEnergyCost.HighFareEnergy)
		fmt.Printf("Niedertarif Energie: %8.2f kWh\t", deviceEnergyCost.LowFareEnergy)
		fmt.Printf("Total Kosten: %8.2f CHF\t", deviceEnergyCost.TotalCost)
		fmt.Printf("Hochtarif Kosten: %8.2f CHF\t", deviceEnergyCost.HighFareCost)
		fmt.Printf("Niedertarif Kosten:: %8.2f CHF\n", deviceEnergyCost.LowFareCost)
	}
}

func CalculateFare(highFare bool, energy float64) float64 {
	if highFare {
		return energy * highFareCost
	}
	return energy * lowFareCost
}

func CheckIfHighFare(timeToCheck time.Time) bool {
	if timeToCheck.Weekday() == time.Sunday {
		return false // Sunday is always low fare
	}

	highFareStart := time.Date(timeToCheck.Year(), timeToCheck.Month(), timeToCheck.Day(), highFareStartHour, 0, 0, 0, timeToCheck.Location())
	highFareEnd := time.Date(timeToCheck.Year(), timeToCheck.Month(), timeToCheck.Day(), highFareEndHour, 0, 0, 0, timeToCheck.Location())
	lowFareNoonStart := time.Date(timeToCheck.Year(), timeToCheck.Month(), timeToCheck.Day(), lowFareNoonStartHour, 0, 0, 0, timeToCheck.Location())
	lowFareNoonEnd := time.Date(timeToCheck.Year(), timeToCheck.Month(), timeToCheck.Day(), lowFareNoonEndHour, 0, 0, 0, timeToCheck.Location())
	if timeToCheck.IsDST() { // Summer fare
		if timeToCheck.After(lowFareNoonStart) && timeToCheck.Before(lowFareNoonEnd) {
			return false
		}
	}
	if timeToCheck.After(highFareStart) && timeToCheck.Before(highFareEnd) {
		return true
	} else {
		return false
	}
}

func SetHour(t time.Time, hour int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), hour, t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

func ConvertAndFixTimeZone(timestamp string) time.Time {
	// Parse input string into a time.Time object
	var inputTime = ConvertStringToTime(timestamp)

	// Get the Zurich timezone, accounting for DST
	loc, err := time.LoadLocation("Europe/Zurich")
	if err != nil {
		panic(err)
	}

	// Convert the input time to Zurich time
	correctedTime := inputTime.In(loc)

	return correctedTime
}

func ConvertStringToTime(timeToCheck string) time.Time {
	convertedTimeToCheck, err := time.Parse(time.RFC3339Nano, timeToCheck)
	if err != nil {
		err = nil
		convertedTimeToCheck, err := time.Parse(time.RFC3339, timeToCheck)
		if err != nil {
			err = nil
			var layout = "2006-01-02T15:04:05.99"
			convertedTimeToCheck, err := time.Parse(layout, "net/url")
			if err != nil {
				err = nil
				s, err := url.PathUnescape(timeToCheck)
				convertedTimeToCheck, err := time.Parse(layout, s)
				if err != nil {
					fmt.Println(err)
					return convertedTimeToCheck
				}
			}
			return convertedTimeToCheck
		}
		return convertedTimeToCheck
	}
	return convertedTimeToCheck
}
