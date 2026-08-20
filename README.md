# railway-go-api
This contains the code necessary to bootstrap a simple REST api in go with minimal 3rd party dependencies.


### Build
go build -o railway-garden-server.exe ./cmd/api

# Features 
## Photo Logic. 
Need somewhere to store a photo of the plant. And the API logic to upload a photo. 

## Needs Watering flag
Set at plant level. 
This logic needs to run daily? Twice daily? Should it be run on the back of rainfall updates?

deficit = sum over lookback_days of:
    (zone.transpiration * plant.et0_multiplier) - (zone.rainfall * plant.rainfall_effectiveness)

Set flag needs_watering

if deficit > deficit_threshold → needs_watering = true
* Watering plants counts as deficit_threshold of rainfall. 
* Setting needs_watering should make a history log entry.
* This logic removes the need for "next watering date"

## Zone API: Plants need attention.
Return list of plant ID's where needs_watering = true.

This logic should live in service/weather.go
Or at Cron.go

Need to re-write this. 
cron should call a water management service. 
This service needs to reference plants, zone & locations/weather data. 
It should not directly talk to the DB. 

