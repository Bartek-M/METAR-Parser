# METAR Parser
Aviation weather parser - `METAR`, with simple runway selection tool created in **Go**. Based on standard [VATSIM Metar Service](https://metar.vatsim.net/) - text API. Helpful for quick look into current weather situation and active runways.

**Example result:**
```txt
METAR Parser - 010000Z

 MVFR | 29/33 | EPWA 010000Z 25005KT 9999 FEW021 BKN042 02/M01 Q1028 NOSIG
```

> If you find any bugs, feel free to create a new issue in this repository.

## New QNH
When QNH change is detected during loop interval, additional information is added
```txt
Q1024 -> Q1025
```

## Setup
### API format
```
url/ STATION
METAR string \n
```
> **NOTE:** 
> - Station ID will be added right after URL
> - API has to be text format with METAR strings only
> - METAR strings have to be separated with newline - `\n`. 

### Configuration file
- `api` - API url link

- `interval` - request interval time in seconds; minimum: 20

- `stations` - list of stations for API fetch

- `excludeNoConfig` - hide fetched data without airport configuration

- `minimums` - weather minimums (from best to worst); list of **4** categories
    - `category` - category names
    - `visibility` - horizontal visibility in **meters** 
    - `ceiling` - ceiling height in **feet**
    > **NOTE:** `third+` category prefers runway with **ILS** system

- `windLimits` - wind limits in **knots**
    - `preferred wind` - importance: other conditions > wind
    - `max wind` - importance: other conditions < wind

- `airports` - airport configuration list; **ICAO**: *config*
    - `runways` - runway list: `{"id": string, "hdg": int, "ils": bool}`
    - `preference` - departure / arrival runway preference order
    - `lvp` - LVP operation [visibility, ceiling]
    > **NOTE:** if no runway preference specified, runways will be used in order


### JSON structure:
```json
{
    "api": "https://metar.vatsim.net/",
    "interval": -1,
    "stations": [
        "EP",
        ...
    ],
    "excludeNoConfig": false,
    "minimums": {
        "category": ["VFR", "MVFR", "IFR", "LIFR"],
        "visibility": [8000, 5000, 1500, 0],
        "ceiling": [3000, 1000, 500, 0]
    },
    "windLimits": [7, 18],
    "airports": {
        "EPWA": {
            "runways": [
                {"id": "29", "hdg": 295},
                {"id": "15", "hdg": 146},
                {"id": "33", "hdg": 326, "ils": true},
                {"id": "11", "hdg": 108, "ils": true},
                ...
            ],
            "preference": {
                "dep": ["29", "15", "33", "11", ...],
                "arr": ["33", "11", "15", "29", ...]
            },
            "lvp": [550, 200]
        },
        ...
    }
}
```

## Manual setup
```bash
go build -o METAR-Parser.exe ./cmd/
```

Compiled application will land in current directory

> **NOTE:** You need [Golang](https://go.dev/dl/) installed on your machine.
> 
> Remove `.exe` if you're using macOS / Linux