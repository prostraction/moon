# API for Moon calculations

## Run

Port by default: 9998

```
git clone https://github.com/prostraction/moon
cd moon
go run cmd/main.go
```


## Methods

### GET /v1/moonPhaseCurrent

Params:
- utc: string in format `UTC+7`, `UTC+09:30`, `-3`
- lang: string. Values available: ("en", "es", "fr", "de", "ru", "jp")
- precision: int. Value between 1 and 10.

Response:

```json
{
     // by the begin of the day, 00:00 today:
    "BeginDay": {
        "MoonDays": 14.33,
        "Illumination": 96.66,
        "Phase": {
            "Name": "Full Moon",
            "NameLocalized": "Luna llena",
            "Emoji": "🌕"
        },
        "Zodiac": {
            "Name": "Aquarius",
            "NameLocalized": "Acuario",
            "Emoji": "♒"
        }
    },
    // current, at this moment (or timestamp given):
    "CurrentState": {
        "MoonDays": 15.21,
        "Illumination": 99.29,
        "Phase": {
            "Name": "Full Moon",
            "NameLocalized": "Luna llena",
            "Emoji": "🌕"
        },
        "Zodiac": {
            "Name": "Pisces",
            "NameLocalized": "Piscis",
            "Emoji": "♓"
        }
    },
    // by the end of the day 00:00 tomorrow:
    "EndDay": {
        "MoonDays": 15.33,
        "Illumination": 99.5,
        "Phase": {
            "Name": "Full Moon",
            "NameLocalized": "Luna llena",
            "Emoji": "🌕"
        },
        "Zodiac": {
            "Name": "Pisces",
            "NameLocalized": "Piscis",
            "Emoji": "♓"
        }
    },
    // all moon's day for today
    "MoonDaysDetailed": {
        "Count": 2,
        "Day": [
            {
                "Begin": "2025-09-06T16:08:06+10:00",
                "End": "2025-09-07T16:08:06+10:00"
            },
            {
                "Begin": "2025-09-07T16:08:06+10:00",
                "End": "2025-09-08T16:08:06+10:00"
            }
        ]
    },
    // all zodiacs for today
    "ZodiacDetailed": {
        "Count": 2,
        "Zodiac": [
            {
                "Name": "Aquarius",
                "NameLocalized": "Acuario",
                "Emoji": "♒",
                "Begin": "2025-09-05T04:07:06+10:00",
                "End": "2025-09-07T16:07:06+10:00"
            },
            {
                "Name": "Pisces",
                "NameLocalized": "Piscis",
                "Emoji": "♓",
                "Begin": "2025-09-07T16:07:06+10:00",
                "End": "2025-09-10T04:07:06+10:00"
            }
        ]
    }
}
```

### GET /v1/moonPhaseTimestamp

Params:
- utc: string in format `UTC+7`, `UTC+09:30`, `-3`
- t: timestamp
- lang: string. Values available: ("en", "es", "fr", "de", "ru", "jp")
- precision: int. Value between 1 and 10.

Response: as GET /v1/moonPhaseCurrent

### GET /v1/moonPhaseDate

Params:
- utc: string in format `UTC+7`, `UTC+09:30`, `-3`
- lang: string. Values available: ("en", "es", "fr", "de", "ru", "jp")
- precision: int. Value between 1 and 10.
- year: int in format YYYY (`1970`)
- month: int in format M (`1`, `12`)
- day: in in format D (`1`, `31`)
- hour: in in format h (`1`, `23`)
- minute: in in format m (`1`, `59`)
- second: in in format s (`1`, `59`)

Response: as GET /v1/moonPhaseCurrent

### GET /v1/moonTableCurrent

Params:
- utc: string in format `UTC+7`, `UTC+09:30`, `-3`

Response:

```json
[
    // first moon of the year
     {
        "NewMoon": "2024-12-31T08:27:49+10:00",
        "FirstQuarter": "2025-01-07T16:15:21+10:00",
        "FullMoon": "2025-01-14T08:27:44+10:00",
        "LastQuarter": "2025-01-22T13:25:30+10:00"
    },
    // second moon of the year
    {
        "NewMoon": "2025-01-29T22:37:18+10:00",
        "FirstQuarter": "2025-02-06T09:58:10+10:00",
        "FullMoon": "2025-02-12T23:54:26+10:00",
        "LastQuarter": "2025-02-21T20:11:33+10:00"
    },
...
    // last moon of the year
     {
        "NewMoon": "2025-12-20T11:44:25+10:00",
        "FirstQuarter": "2025-12-27T19:53:55+10:00",
        "FullMoon": "2026-01-03T20:04:15+10:00",
        "LastQuarter": "2026-01-11T02:44:20+10:00"
    }
]
```

### GET /v1/moonTableCurrent

Params:
- utc: string in format `UTC+7`, `UTC+09:30`, `-3`
- year: int in format YYYY (`1970`)

Response: as GET /v1/moonTableCurrent
