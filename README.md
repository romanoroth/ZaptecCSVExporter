# Zaptec CSV Exporter

## Zaptec Energy Cost Analysis Tool

This tool is specifically designed for owners of Zaptec charging devices. It utilizes the Zaptec API to fetch energy consumption data and performs an analysis distinguishing between high fare and low fare periods. The primary purpose of this tool is to provide a detailed analysis of energy usage during these differing fare periods, helping users to understand and optimize their energy costs.

## Features

- Fetch energy consumption data specifically for Zaptec devices via the Zaptec API.
- Calculate energy costs for high fare and low fare periods.
- Output detailed analysis of high and low fare energy usage into a CSV file.
- Handle time zone conversions and adjustments for daylight saving time.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Go programming environment (Go 1.14 or higher recommended).
- Internet access to make requests to the Zaptec API.
- A Zaptec account with valid credentials.

## Configuration

This tool requires you to enter your Zaptec API credentials and other specific information:

```go
const username = "ENTER_HERE_YOUR_USERNAME"          // Your Zaptec API username.
const password = "ENTER_HERE_YOUR_PASSWORD"          // Your Zaptec API password.
const installationId = "ENTER_HERE_YOUR_INSTALLATION_ID" // The installation ID for your Zaptec device.
```

Replace ENTER_HERE_YOUR_USERNAME, ENTER_HERE_YOUR_PASSWORD, and ENTER_HERE_YOUR_INSTALLATION_ID with your actual Zaptec API credentials and installation ID.

Additionally, configure the following constants based on your tariff and timing preferences:

```go
highFareStartHour, highFareEndHour, lowFareNoonStartHour, lowFareNoonEndHour
highFareCost, lowFareCost
```

## Usage
To run this program, execute the following command in your terminal:

```bash
go run main.go
```

This command initiates the application, which will retrieve energy consumption data for your specified date range, calculate the costs for high and low fare periods, and generate a detailData.csv file containing the detailed analysis.

## Main Use Case
The primary use case of this tool is to analyze and understand the distribution of energy consumption between high fare and low fare periods for Zaptec charging devices. This analysis helps users make informed decisions to optimize their energy usage and reduce costs.

Contributing
Contributions to this project are welcome. Please adhere to the standard fork-and-pull request workflow on GitHub to submit your changes.

## License
This project is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
