# Country Dashboard Service  
### Course: PROG2005 - Cloud Technologies  
<br>

## Overview  
In this group assignment, we are going to develop a REST web application in Golang that provides the client with the ability to configure information dashboards that are dynamically populated when requested. The dashboard configurations are saved persistently in a remote database and populated based on external services.  

It will also include a simple notification service that can listen to specific events. The application will be dockerized and deployed using an IaaS system.  

**Documentation for the external APIs used in this project:**

REST Countries API (Country information) <br>
Documentation: http://129.241.150.113:8080/

Open-Meteo APIs (For temperature and weather conditions) <br>
Documentation: https://open-meteo.com/en/features#available-apis

Currency API (For currency exchange) <br>
Documentation: http://129.241.150.113:9090/

---
## Contributions 

Mustafa:
- Firebase setup 
- Registration & Status endpoints
- README.md 

August:
- Dashboard endpoint 
- Test files 

Sethushan: 
- Webhooks 
- Notification endpoint

---
# How to run locally 

To run the Country Dashboard Service locally, follow the steps below:

### Prerequisites

- Go (version 1.18 or higher) installed.
- Your own running Firestore instance on https://firebase.google.com/  
- Internet connection to connect to external API's

### Steps:

1. Clone the repository 
2. Install necessary dependecies
3. Set up your own firebase and download the key. Name it `firebaseKey.json`
4. Make a new folder in your project files called ".env" and put the firebase key in that folder
5. Run the application with 
```bash
go run main.go
```
6. To access the application. Open a browser and go to `http://localhost:8080` with the various endpoints listed bellow.


## Final Endpoints  

- `/dashboard/v1/registrations/`  
- `/dashboard/v1/dashboards/`  
- `/dashboard/v1/notifications/`  
- `/dashboard/v1/status/`  

For detailed information and requirements, see below.

---

## Endpoint: `/dashboard/v1/registrations/`

### `(POST)` - request  
Should include:  

- Country name  
- ISO code for the country  
- Temperature: check if temperature is measured in Celsius (`true/false`)  
- Precipitation: show if it is raining, showering, or snowing? (`true/false`)  
- Capital: check if the name of the capital is shown (`true/false`)  
- Coordinates: check if coordinates are shown (`true/false`)  
- Population: check if population is shown (`true/false`)  
- Area: check if land area size is shown (`true/false`)  
- TargetCurrencies: shows all exchange rates that are displayed  

#### Request Body – `POST /dashboard/v1/registrations/`

```json
{
  "country": "Norway",
  "isoCode": "NO",
  "features": {
    "temperature": true,
    "precipitation": true,
    "capital": true,
    "coordinates": true,
    "population": true,
    "area": true,
    "targetCurrencies": ["EUR", "USD", "SEK"]
  }
}
```

Note that the POST request will be invalid if:
- Country name is not recognized from the REST countries API
- isoCode do not match it's countries iso3 code  

### (GET) - request

Returns all stored configurations (records from previous POST requests)

#### Response – `GET /dashboard/v1/registrations/`

```json
[
  {
    "id": "516dba7f015f2a68",
    "country": "Norway",
    "isoCode": "NO",
    "features": {
      "temperature": true,
      "precipitation": true,
      "capital": true,
      "coordinates": true,
      "population": true,
      "area": false,
      "targetCurrencies": ["EUR", "USD", "SEK"]
    },
    "lastChange":"2025-04-07 14:54:02 CEST"
  },
  {
    "id": "bc89adc23e27f42a",
    "country": "Denmark",
    "isoCode": "DK",
    "features": {
      "temperature": false,
      "precipitation": true,
      "capital": true,
      "coordinates": true,
      "population": false,
      "area": true,
      "targetCurrencies": ["NOK", "MYR", "JPY", "EUR"]
    },
    "lastChange": "2025-04-05 20:30:00 CEST"
  }
```

### (GET) - request with specific id

Returns a specific stored configuration.

#### Response – `GET /dashboard/v1/registrations/{id}`

Example for `id` = `516dba7f015f2a68`

```json
{
  "id": "516dba7f015f2a68",
  "country": "Norway",
  "isoCode": "NO",
  "features": {
    "temperature": true,
    "precipitation": true,
    "capital": true,
    "coordinates": true,
    "population": true,
    "area": false,
    "targetCurrencies": ["EUR", "USD", "SEK"]
  },
  "lastChange": "2025-04-07 14:54:02 CEST"
}
```

### Time format management
Originally time was shown in `unix.time` timestamp, 
this would work and looks fine and readable on the database, but since the client would read the time in JSON format it would look like this
`2025-04-09T10:07:57.248557Z`, we therefore have altered it to always show the time in CEST as shown above. 
This is not flexible from where the client is located but rather will always be set to the current time in CEST.  

### (PUT) - request

Update an individual configuration identified by its ID. This will update configuration as requested and will also update the timestamp of `lastChange`.

#### Request – `PUT /dashboard/v1/registrations/{id}`

Where `{id}` is the id of the registration to be changed

```json
{
  "country": "Norway",
  "isoCode": "NO",
  "features": {
    "temperature": false,  // value to be changed
    "precipitation": true,
    "capital": true,
    "coordinates": true,
    "population": true,
    "area": false,
    "targetCurrencies": ["EUR", "SEK"]  // value to be changed
  }
}
```

### (DELETE) - Request

Delete an individual configuration identified by its ID. This request will result in the deletion of the configuration from the server.

#### Request – `DELETE /dashboard/v1/registrations/{id}`

Example request for `id` = `516dba7f015f2a68`:

```json
{
  "Method": "DELETE",
  "Path": "/dashboard/v1/registrations/516dba7f015f2a68"
}
```



---

## Endpoint: `/dashboard/v1/dashboards/`

### (GET) - request

Retrieve the populated dashboard for a given country configuration. The dashboard includes features like: 
- Temperature
- Precipitation
- Capital
- Coordinates (latitude, longitude)
- Population
- Area
- TargetCurrencies
- Last retrieval

If any of these features are checked off as false in the `registrations` document, it will not show up in the dashboards page. 

**Request Parameters**: (brought from the registrations)
- Country name
- ISO code <br>

Validators will make sure that these match when trying to POST a registration

Example request: 
```
Request: GET
Path: /dashboard/v1/dashboards/{id}
```
Where {id} is an id from a register. <br>
**Response**:

```json
{
"country": "Norway",
"isoCode": "NO",
"features": {
        "temperature": -1.2,         // Mean temperature across all forecasted temperature values for country's coordinates
        "precipitation": 0.80,       // Mean precipitation across all returned precipitation values
        "capital": "Oslo",           // Capital: Where multiple values exist, take the first
        "coordinates": {             // Those are the country geocoordinates
                "latitude": 62.0,
                "longitude": 10.0
                        },
        "population": 5379475,
        "area": 323802.0,
        "targetCurrencies": {
                "EUR": 0.087701435,  // this is the current NOK to EUR exchange rate (where multiple currencies exist for a given country, take the first)
                "USD": 0.095184741, 
                "SEK": 0.97827275
                }
         },
"lastRetrieval":"2025-04-09 14:54:02 CEST" // this should be the current time (i.e., the time of retrieval)
}
```
## Endpoint: `/dashboard/v1/notifications/`

Users can register webhooks that are triggered by the service based on specified events.
These events consists of: 
- If a new configuration is created
- If a configuration is changed or deleted
- Invocation event (when dashboard for a given country is invoked)
- Users can register multiple webhooks

These webhooks are stored persistently in the database and will stay until manually deleted with a `DELETE` request 

### (POST) - Register Webhook

#### Request – `POST /dashboard/v1/notifications/`


```
Method: POST
Path: /dashboard/v1/notifications/
Content type: application/json
```

## Webhook Registration Body

When registering a webhook, the body should contain the following information:

- `url`: The URL that will be triggered when the specified event occurs. This is the service endpoint to be invoked.
- `country`: The country associated with the event trigger. If this field is left empty, the webhook will be triggered for any country.
- `event`: The event type that triggers the webhook. The following event types are supported:
  - `REGISTER`: The webhook is invoked when a new configuration is registered.
  - `CHANGE`: The webhook is invoked when a configuration is modified.
  - `DELETE`: The webhook is invoked when a configuration is deleted.
  - `INVOKE`: The webhook is invoked when a dashboard for a given country is retrieved (i.e., a GET request on the `/dashboard/v1/dashboards/` endpoint).


**Body** (Example for an `INVOKE` event):

```json
{
   "url": "https://localhost:8080/client/",  // URL to be invoked when event occurs
   "country": "NO",                          // Country that is registered, or empty if all countries
   "event": "INVOKE"                         // Event on which it is invoked
}
```
This will respond with the ID for the registration that can used to see detail information or to delete the webhook registration.
```
{
    "id": {id}
}
```
Where `{id}` is the ID of the registation.
### (DELETE) - Delete webhook

```
Method: DELETE
Path: /dashboard/v1/notifications/{id}
```
Where {id} is the ID returned during the webhook registration


### (GET) - View registered webhook
```
Method: GET
Path: /dashboard/v1/notifications/{id}
```
Where {id} is the ID FOR THE webhook registration
## Webhook Registration Response

The response to a webhook registration request will include the ID assigned by the server to the registered webhook. This ID can be used to view or delete the webhook.

### Response Body

```json
{
   "id": "OIdksUDwveiwe",  // Unique ID assigned by the server
   "url": "https://localhost:8080/client/",  // URL to be invoked when event occurs
   "country": "NO",                        // Country that is registered, or empty if all countries
   "event": "INVOKE"                       // Event on which it is invoked
}
````

### (GET) - View all registered webhooks

**Method:** GET  
**Path:** /dashboard/v1/notifications/

## Webhook Registrations Response

The response to a GET request for all registered webhooks will return a list of all currently registered webhooks. Each entry includes the details of the webhook, such as the ID, URL, country, and event type.

### Response Body

```json
{
   {
      "id": "OIdksUDwveiwe",  // Unique ID assigned by the server
      "url": "https://localhost:8080/client/",  // URL to be invoked when event occurs
      "country": "NO",                         // Country that is registered, or empty if all countries
      "event": "INVOKE"                        // Event on which it is invoked
   },
   {
      "id": "DiSoisivucios",  // Another unique ID for a different webhook
      "url": "https://localhost:8081/anotherClient/",  // URL for a different service to invoke
      "country": "",  // Empty or omitted if registered for all countries
      "event": "REGISTER"  // Event type for this webhook
   }
   // Additional webhook registrations...
}
```


### Webhook Invocation (upon trigger)

When a webhook is triggered, it should send information as follows. Where multiple webhooks are triggered, the information should be sent separately (i.e., one notification per triggered webhook). 

For testing purposes, you may want to set up another service that is able to receive the invocation. You can use [Webhook.site](https://webhook.site/) for development testing.

**Method:** POST  
**Path:** `<url specified in the corresponding webhook registration>`  
**Content type:** application/json  

### Request Body (Exemplary message based on schema):

```json
{
   "id": "OIdksUDwveiwe",       // Unique ID assigned to the webhook
   "country": "NO",             // Country code related to the event (can be empty for all countries)
   "event": "INVOKE",           // Event type that triggered the webhook
   "time": "20240223 06:23"     // Time at which the event occurred
}
```

### Endpoint 'Status': Monitoring service availability

The status interface indicates the availability of all individual services this service depends on. These can include more services than the ones specified. If additional services are included, you can specify them with the suffix `api`. The reporting occurs based on the status codes returned by the dependent services. The status interface also provides information about the number of registered webhooks and the uptime of the service.

#### Request (GET)

**Method:** GET  
**Path:** /dashboard/v1/status/

#### Response

**Content type:** application/json

**Status code:** 200 if everything is OK, appropriate error code otherwise.

### Response Body:

```json
{
   "countries_api": <http status code for *REST Countries API*>,
   "meteo_api": <http status code for *Meteo API*>, 
   "currency_api": <http status code for *Currency API*>,
   "notification_db": <http status code for *Notification database*>,
   ...
   "webhooks": <number of registered webhooks>,
   "version": "v1",
   "uptime": <time in seconds from the last service restart>
}


