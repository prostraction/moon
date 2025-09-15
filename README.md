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

#### Response:


| Response Variable | Type | Description | Example Value |
| :--- | :--- | :--- | :--- |
|`BeginDay` | `Object` | Data for the beginning of the requested day (00:00) | - |
|⠀⠀`BeginDay.MoonDays` | `Float` | Lunar day number at day start | `23.54` |
|⠀⠀`BeginDay.Illumination` | `Float` | Percentage of Moon's disk illuminated | `38.27` |
|⠀⠀`BeginDay.Phase` | `Object` | Lunar phase details | - |
|⠀⠀⠀⠀⠀`BeginDay.Phase.Name` | `String` | Phase name in English | `"Waning Crescent"` |
|⠀⠀⠀⠀⠀`BeginDay.Phase.NameLocalized` | `String` | Localized phase name | `"Убывающий серп"` |
|⠀⠀⠀⠀⠀`BeginDay.Phase.Emoji` | `String` | Phase emoji | `"🌘"` |
|⠀⠀⠀⠀⠀`BeginDay.Phase.IsWaxing` | `Boolean` | True if illumination is increasing | `false` |
|⠀⠀`BeginDay.Zodiac` | `Object` | Zodiac sign details | - |
|⠀⠀⠀⠀⠀`BeginDay.Zodiac.Name` | `String` | Zodiac name in English | `"Gemini"` |
|⠀⠀⠀⠀⠀`BeginDay.Zodiac.NameLocalized` | `String` | Localized zodiac name | `"Близнецы"` |
|⠀⠀⠀⠀⠀`BeginDay.Zodiac.Emoji` | `String` | Zodiac emoji | `"♊"` |
|⠀⠀`BeginDay.Position` | `Object` | Astronomical position data | - |
|⠀⠀⠀⠀`BeginDay.Position.Timestamp` | `Integer` | Unix timestamp of calculation | `1757962800` |
|⠀⠀⠀⠀`BeginDay.Position.TimeISO` | `String` | ISO 8601 timestamp | `"2025-09-16T00:00:00+05:00"` |
|⠀⠀⠀⠀`BeginDay.Position.AzimuthDegrees` | `Float` | Compass direction (0°=North) | `57.1` |
|⠀⠀⠀⠀`BeginDay.Position.AltitudeDegrees` | `Float` | Angle above horizon (negative = below) | `8.8` |
|⠀⠀⠀⠀`BeginDay.Position.Direction` | `String` | Cardinal direction abbreviation | `"ENE"` |
|⠀⠀⠀⠀`BeginDay.Position.DistanceKm` | `Float` | Earth-Moon distance in km | `376559.9` |
|`CurrentState` | `Object` | Data at exact time of API request (same structure as BeginDay) | - |
| `EndDay` | `Object` | Data for end of requested day (00:00 next day) (same structure as BeginDay) | - |
| `MoonDaysDetailed` | `Object` | Detailed lunar day information | - |
| ⠀⠀`MoonDaysDetailed.Count` | `Integer` | Number of lunar days this calendar day | `2` |
| ⠀⠀`MoonDaysDetailed.Day` | `Array<Object>` | Array of lunar day periods | - |
| ⠀⠀⠀⠀`MoonDaysDetailed.Day[].Begin` | `String` | Start time of lunar day (ISO 8601) | `"2025-09-15T22:37:45+05:00"` |
| ⠀⠀⠀⠀`MoonDaysDetailed.Day[].IsBeginExists` | `Boolean` | True if start time is past/present | `true` |
| ⠀⠀⠀⠀`MoonDaysDetailed.Day[].End` | `String` | End time of lunar day (ISO 8601) | `"2025-09-16T23:56:10+05:00"` |
| ⠀⠀⠀⠀`MoonDaysDetailed.Day[].IsEndExists` | `Boolean` | True if end time is past | `true`, `false` |
| `ZodiacDetailed` | `Object` | Detailed zodiac transit information | - |
| ⠀⠀`ZodiacDetailed.Count` | `Integer` | Number of zodiac signs this day | `1` |
| ⠀⠀`ZodiacDetailed.Zodiac` | `Array<Object>` | Array of zodiac transit periods | - |
| ⠀⠀⠀⠀`ZodiacDetailed.Zodiac[].Name` | `String` | Zodiac sign name | `"Gemini"` |
| ⠀⠀⠀⠀`ZodiacDetailed.Zodiac[].NameLocalized` | `String` | Localized zodiac name | `"Близнецы"` |
| ⠀⠀⠀⠀`ZodiacDetailed.Zodiac[].Emoji` | `String` | Zodiac emoji | `"♊"` |
| ⠀⠀⠀⠀`ZodiacDetailed.Zodiac[].Begin` | `String` | Entry time into sign (ISO 8601) | `"2025-09-14T23:07:06+05:00"` |
| ⠀⠀⠀⠀`ZodiacDetailed.Zodiac[].End` | `String` | Exit time from sign (ISO 8601) | `"2025-09-17T11:07:06+05:00"` |
| `MoonRiseAndSet` | `Object` | Moon rise/set/meridian events | - |
| ⠀⠀`MoonRiseAndSet.Moonrise` | `Object` | Moonrise event data | - |
| ⠀⠀⠀⠀`MoonRiseAndSet.Moonrise.Timestamp` | `Integer` | Moonrise Unix timestamp | `1758048970` |
| ⠀⠀⠀⠀`MoonRiseAndSet.Moonrise.TimeISO` | `String` | Moonrise ISO time | `"2025-09-16T23:56:10+05:00"` |
| ⠀⠀⠀⠀`MoonRiseAndSet.Moonrise.AzimuthDegrees` | `Float` | Moonrise azimuth | `47.3` |
| ⠀⠀⠀⠀`MoonRiseAndSet.Moonrise.AltitudeDegrees` | `Float` | Moonrise altitude | `-0.6` |
| ⠀⠀⠀⠀`MoonRiseAndSet.Moonrise.Direction` | `String` | Moonrise direction | `"NE"` |
| `MoonRiseAndSet.Moonset` | `Object` | Moonset event data (same structure as Moonrise) | - |
| `MoonRiseAndSet.Meridian` | `Object` | Meridian event data (same structure as Moonrise) | - |
| `MoonRiseAndSet.DistanceKm` | `Float` | Approximate Earth-Moon distance | `379004.1` |
| `MoonRiseAndSet.IsMoonRise` | `Boolean` | True if moonrise occurs today | `true` |
| `MoonRiseAndSet.IsMoonSet` | `Boolean` | True if moonset occurs today | `true` |
| `MoonRiseAndSet.IsMeridian` | `Boolean` | True if meridian transit occurs today | `true` |

```json
{
  "BeginDay": {
    "MoonDays": 23.54,
    "Illumination": 38.27,
    "Phase": {
      "Name": "Waning Crescent",
      "NameLocalized": "Убывающий серп",
      "Emoji": "🌘",
      "IsWaxing": false
    },
    "Zodiac": {
      "Name": "Gemini",
      "NameLocalized": "Близнецы",
      "Emoji": "♊"
    },
    "Position": {
      "Timestamp": 1757962800,
      "TimeISO": "2025-09-16T00:00:00+05:00",
      "AzimuthDegrees": 57.1,
      "AltitudeDegrees": 8.8,
      "Direction": "ENE",
      "DistanceKm": 376559.9
    }
  },
  "CurrentState": {
    "MoonDays": 24.43,
    "Illumination": 29,
    "Phase": {
      "Name": "Waning Crescent",
      "NameLocalized": "Убывающий серп",
      "Emoji": "🌘",
      "IsWaxing": false
    },
    "Zodiac": {
      "Name": "Gemini",
      "NameLocalized": "Близнецы",
      "Emoji": "♊"
    },
    "Position": {
      "Timestamp": 1758039653,
      "TimeISO": "2025-09-16T21:20:53+05:00",
      "AzimuthDegrees": 15.7,
      "AltitudeDegrees": -12.5,
      "Direction": "NNE",
      "DistanceKm": 379635
    }
  },
  "EndDay": {
    "MoonDays": 24.54,
    "Illumination": 27.91,
    "Phase": {
      "Name": "Waning Crescent",
      "NameLocalized": "Убывающий серп",
      "Emoji": "🌘",
      "IsWaxing": false
    },
    "Zodiac": {
      "Name": "Gemini",
      "NameLocalized": "Близнецы",
      "Emoji": "♊"
    },
    "Position": {
      "Timestamp": 1758049200,
      "TimeISO": "2025-09-17T00:00:00+05:00",
      "AzimuthDegrees": 48,
      "AltitudeDegrees": -0.1,
      "Direction": "NE",
      "DistanceKm": 380020.6
    }
  },
  "MoonDaysDetailed": {
    "Count": 2,
    "Day": [
      {
        "Begin": "2025-09-15T22:37:45+05:00",
        "IsBeginExists": true,
        "End": "2025-09-16T23:56:10+05:00",
        "IsEndExists": true
      },
      {
        "Begin": "2025-09-16T23:56:10+05:00",
        "IsBeginExists": true,
        "End": "0001-01-01T00:00:00Z",
        "IsEndExists": false
      }
    ]
  },
  "ZodiacDetailed": {
    "Count": 1,
    "Zodiac": [
      {
        "Name": "Gemini",
        "NameLocalized": "Близнецы",
        "Emoji": "♊",
        "Begin": "2025-09-14T23:07:06+05:00",
        "End": "2025-09-17T11:07:06+05:00"
      }
    ]
  },
  "MoonRiseAndSet": {
    "Moonrise": {
      "Timestamp": 1758048970,
      "TimeISO": "2025-09-16T23:56:10+05:00",
      "AzimuthDegrees": 47.3,
      "AltitudeDegrees": -0.6,
      "Direction": "NE"
    },
    "Moonset": {
      "Timestamp": 1758022192,
      "TimeISO": "2025-09-16T16:29:52+05:00",
      "AzimuthDegrees": 314.5,
      "AltitudeDegrees": -0.6,
      "Direction": "NW"
    },
    "Meridian": {
      "Timestamp": 1757990436,
      "TimeISO": "2025-09-16T07:40:36+05:00",
      "AzimuthDegrees": 180,
      "AltitudeDegrees": 65.7,
      "Direction": "S"
    },
    "DistanceKm": 379004.1,
    "IsMoonRise": true,
    "IsMoonSet": true,
    "IsMeridian": true
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
