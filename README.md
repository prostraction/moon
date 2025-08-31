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
        "MoonDays": 8.75,
        "Illumination": 52.96,
        "Phase": {
            "Name": "First quarter",
            "Emoji": "🌓"
        },
        "Zodiac": {
            "Name": "Sagittarius",
            "Emoji": "♐"
        }
    },
    // current, at this moment (or timestamp given):
    "CurrentState": {
        "MoonDays": 8.65,
        "Illumination": 55.89,
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
        "MoonDays": 7.75,
        "Illumination": 47.24,
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
        "NewMoon": "2024-12-30T22:27:49Z",
        "FirstQuarter": "2025-01-07T06:15:49Z",
        "FullMoon": "2025-01-13T22:27:44Z",
        "LastQuarter": "2025-01-22T03:25:49Z"
    },
    {
        "NewMoon": "2025-01-29T12:37:18Z",
        "FirstQuarter": "2025-02-05T23:58:25Z",
        "FullMoon": "2025-02-12T13:54:26Z",
        "LastQuarter": "2025-02-21T10:12:41Z"
    },
    // second moon of the year
...
    // last moon of the year
    {
        "NewMoon": "2025-12-20T01:44:25Z",
        "FirstQuarter": "2025-12-27T09:54:01Z",
        "FullMoon": "2026-01-03T10:04:15Z",
        "LastQuarter": "2026-01-10T16:46:01Z"
    }
]
```

### GET /v1/moonTableCurrent

Params:
- utc: string in format `UTC+7`, `UTC+09:30`, `-3`
- year: int in format YYYY (`1970`)

Response: as GET /v1/moonTableCurrent
