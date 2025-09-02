# API for Moon calculations

## Run

Port by default: 3000

```
git clone https://github.com/prostraction/moon
cd moon
go run cmd/main.go
```


## Methods

### GET /v1/moonPhaseCurrent

Params:
- utc: string in format `UTC+7`, `UTC+09:30`, `-3`

Response:

```json
{
    // by the end of the day 00:00 tomorrow:
     "EndDay": {
        "MoonDays": 9.16,
        "Illumination": 56.93,
        "Phase": {
            "Name": "Waxing Gibbous",
            "Emoji": "🌔"
        },
        "Zodiac": {
            "Name": "Sagittarius",
            "Emoji": "♐"
        }
    },
    // current, at this moment (or timestamp given):
     "CurrentState": {
        "MoonDays": 9.12,
        "Illumination": 56.53,
        "Phase": {
            "Name": "Waxing Gibbous",
            "Emoji": "🌔"
        },
        "Zodiac": {
            "Name": "Sagittarius",
            "Emoji": "♐"
        }
    },
    // by the begin of the day, 00:00 today:
    "BeginDay": {
        "MoonDays": 8.16,
        "Illumination": 51.21,
        "Phase": {
            "Name": "First quarter",
            "Emoji": "🌓"
        },
        "Zodiac": {
            "Name": "Sagittarius",
            "Emoji": "♐"
        }
    },
    "MoonDaysDetailed": null,
    "ZodiacDetailed": {
        "Count": 0,
        "Zodiac": null
    }
}
```

### GET /v1/moonPhaseTimestamp

Params:
- utc: string in format `UTC+7`, `UTC+09:30`, `-3`
- t: timestamp

Response: as GET /v1/moonPhaseCurrent

### GET /v1/moonPhaseDate

Params:
- utc: string in format `UTC+7`, `UTC+09:30`, `-3`
- year: int in format YYYY (`1970`)
- month: int in format M (`1`, `12`) // change later to MM
- day: in in format D (`1`, `31`) // change later to DD
- hour: in in format h (`1`, `23`) // change later to hh
- minute: in in format m (`1`, `59`) // change later to mm
- second: in in format s (`1`, `59`) // change later to ss

Response: as GET /v1/moonPhaseCurrent



### GET /v1/moonTableCurrent

Params:
- utc: string in format `UTC+7`, `UTC+09:30`, `-3`

Response:

```json
[
    // first moon of the year
    {
        "NewMoon": "2022-12-23T00:17:56-10:00",
        "FirstQuarter": "2022-12-29T00:09:17-10:00",
        "FullMoon": "2023-01-06T13:09:55-10:00",
        "LastQuarter": "2023-01-14T10:56:00-10:00"
    },
    // second moon of the year
    {
        "NewMoon": "2023-01-21T10:55:30-10:00",
        "FirstQuarter": "2023-01-27T23:25:44-10:00",
        "FullMoon": "2023-02-05T08:30:44-10:00",
        "LastQuarter": "2023-02-13T10:04:20-10:00"
    },
...
    // last moon of the year
    {
        "NewMoon": "2023-12-12T13:32:07-10:00",
        "FirstQuarter": "2023-12-18T11:36:55-10:00",
        "FullMoon": "2023-12-26T14:33:43-10:00",
        "LastQuarter": "2024-01-03T06:28:26-10:00"
    }
]
```

### GET /v1/moonTableCurrent

Params:
- utc: string in format `UTC+7`, `UTC+09:30`, `-3`
- year: int in format YYYY (`1970`)

Response: as GET /v1/moonTableCurrent
