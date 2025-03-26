# Country Dashboard Service 

Final endpoints:

/dashboard/v1/registrations/
/dashboard/v1/dashboards/
/dashboard/v1/notifications/
/dashboard/v1/status/

For information and requirements see bellow

## Endpoint: Registrations
(POST)
Should include:

-country name
-iso code for said country 
-check if temprature is measured in Celsius or not (true/false)
-precipitation -> is it raining, showering or snowing? (t/f)
-Capital -> check if the name of the capital is shown (t/f)
-Check if coordinates are shown or not (t/f)
-Check if population is shown or not (t/f)
-Check if land area size is shown or not (t/f)
-TargetCurrencies shows all exhange rates that are shown.


(GET)

Should include:

A record of post requests for all stored configurations. 
i.e info from post on every country. 


(PUT) 

Request for an updated individual configuration identified by its ID.
Will also update the associated timestamp (lastChange)


(DELETE)
Delete individual configuration identified by its ID. Should lead to deletion of the configuration on the server


## Endpoint - Dashboards
(GET)
Should include:
-Country name
-isoCode 
-features:
	-temprature
	-precipipitation
	-capital
	-coordinates:{
		-latitude
		-longditude
		}
	-population
	-area
	-targetCurrencies
-last retrival (should be current time/ time of retrival)

## Endpoint - Notification
(POST)
-URL will be triggered upon an event (service should be invoked)
-Diffrent triggers:
    REGISTER - If a new config is registred
    CHANGE - If a config is modified
    DELETE - If a config is deleted 
    INVOKE - If a dashboard is retriedved (i.e, GET request on dashboard endpoint) 

(DELETE)
-Deletes a webhook given its id

(GET)
Similar to POST request body, but check with ID to a spesific webhook
-If no ID is supplied, it shall show all registred webhooks


WEBHOOK INVOCATION (upon a trigger)

When a webhook is triggered it shall send information as follows. For multiple webhooks, each hook should have their own information and notification


## Endpoint- Status
(GET)
-Should show status on all API's in use(countries_api, meteo_api, currency_api, notification_db)
-Show the number of registred webhooks
-version 
-uptime 