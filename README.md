# API for Moon calculations

## Run

Port by default: 9998

```
git clone https://github.com/prostraction/moon
cd moon
go run cmd/main.go
python addon/skyfield/server.py
```


## Methods

### GET /v1/moonPhaseDate

The method returns the Moon parameters for the specified day and time. If the day or time is not specified, the current value for the unspecified fields is taken. If longitude and latitude are specified, the response will contain additional structures.

### Arguments

  | Parameter | Type | Description | Example Value |
| :--- | :--- | :--- |  :--- | 
|`utc` | `string [optional, default="UTC+0"]` | UTC in format `UTC+7`, `UTC+09:30`, `-3` | `UTC+4`
|`lang` | `string [optional, default="en"]` | Values available: ("en", "es", "fr", "de", "ru", "jp") | `es`
|`precision` | `int [optional, default=2]` | How many digits after ```.``` will be in output. Allowed range: [1, 20] | `5`
|`latitude` | `float [optional, default=none]` | Latitude of viewer's place. Used for moon position calculations: ```MoonDaysDetailed```, ```MoonRiseAndSet```, and ```MoonPosition``` object | `51.1655`
|`longitude` | `float [optional, default=none]` | Longitude of viewer's place. Used for moon position calculations: ```MoonDaysDetailed```, ```MoonRiseAndSet```, and ```MoonPosition``` object | `71.4272`
|`year` | `int [optional, default=<current year>]` | Format: YYYY Allowed range: [1, 9999] | `2025`
|`month` | `int [optional, default=<current month>]` | Format: M or MM. Allowed range: [1, 12] | `01` or `1`
|`day` | `int [optional, default=<current day>]` | Format: D or DD. Allowed range: [1, 31] | `01` or `1`
|`hour` | `int [optional, default=<current hour>]` | Format: h or hh. Allowed range: [0, 23] | `01` or `1`
|`minute` | `int [optional, default=<current minute>]` | Format: m or mm. Allowed range: [0, 59] | `01` or `1`
|`second` | `int [optional, default=<current second>]` | Format: s or ss. Allowed range: [0, 59] | `01` or `1`

----------------------------------------------------------------

### /v1/moonPhaseDate Response:

The method returns 6 objects:
- ```BeginDay```, ```CurrentState```, ```EndDay``` objects of the ```MoonStat``` structure to display the position of the moon at the beginning of the day, the specified time and the end of the day, respectively;
- ```MoonDaysDetailed```, a structure for determining the number of lunar days on a given day;
- ```ZodiacDetailed```, a structure for determining which zodiac sign the moon is in on a given time interval, when it began and ended. It contains an array for each lunar day that falls on a given Earth day;
- ```MoonRiseAndSet```, a structure for determining the moonrise, moonset and meridian on a given day.

<details>
  <summary><strong>Table</strong></summary>
  
  | Response Variable | Type | Description |
| :--- | :--- | :--- | 
|`BeginDay` | `Object of struct MoonStat [required]` | Data for the beginning of the requested day (00:00) |
|`CurrentState` | `Object of struct MoonStat [required]` | Data at specified time of requested day: hour, minute and second from request arguments. |
|`EndDay` | `Object of struct MoonStat [required]` | Data for end of requested day (00:00 next day) |
|`MoonDaysDetailed` | `Object of struct MoonDaysDetailed [optional]` | Detailed lunar day information that falls on a given Earth day. Exists only if latitude and longitude are specified. |
|`ZodiacDetailed` | `Object of struct ZodiacDetailed [required]` | Detailed zodiac transit information | 
|`MoonRiseAndSet` | `Object of struct MoonRiseAndSet [optional]` | Moon rise/set/meridian events. Exists only if latitude and longitude are specified. |
  
</details>

----------------------------------------------------------------

#### MoonStat (used as ```BeginDay```, ```CurrentState```, ```EndDay```)

```MoonStat``` objects are used to display at a  time.
In case of a method response, MoonStat will contain the values:
- BeginDay: start of the  day, 00:00AM
- CurrentState:  time
- EndDay: start of the next day, 00:00AM

<details>
  <summary><strong>Table</strong></summary>

  | Response Variable | Type | Description | Example Value |
| :--- | :--- | :--- | :--- |
|`MoonStat.MoonDays` | `Float [required]` | Lunar day number | `23.54` |
|`MoonStat.Illumination` | `Float [required]` | Percentage of Moon's disk illuminated | `38.27` |
|`MoonStat.Phase` | `Object of struct Phase [required]` | Lunar phase details | - |
|`MoonStat.Zodiac` | `Object of struct Zodiac [required]` | Zodiac sign details | - |
|`MoonStat.Position` | `Object of struct MoonPosition [optional]` | Moon position data. Exists only if latitude and longitude are specified. | - |

Phase structure:

  | Response Variable | Type | Description | Example Value |
| :--- | :--- | :--- | :--- |
|`Phase.Name` | `String [required]` | Phase name in English | `"Waning Crescent"` |
|`Phase.NameLocalized` | `String [required]` | Localized phase name | `"Убывающий серп"` |
|`Phase.Emoji` | `String [required]` | Phase emoji | `"🌘"` |
|`Phase.IsWaxing` | `Boolean [required]` | True if Moon is waxing / illumination is increasing | `false` |

Zodiac structure:

  | Response Variable | Type | Description | Example Value |
| :--- | :--- | :--- | :--- |
|`Zodiac.Name` | `String [required]` | Zodiac name in English | `"Gemini"` |
|`Zodiac.NameLocalized` | `String [required]` | Localized zodiac name | `"Близнецы"` |
|`Zodiac.Emoji` | `String [required]` | Zodiac emoji | `"♊"` |

Position structure (Exists only if latitude and longitude are specified):

  | Response Variable | Type | Description | Example Value |
| :--- | :--- | :--- | :--- |
|`Position.Timestamp` | `Integer [required]` | Unix timestamp of calculation | `1757962800` |
|`Position.TimeISO` | `String [required]` | ISO 8601 timestamp | `"2025-09-16T00:00:00+05:00"` |
|`Position.AzimuthDegrees` | `Float [required]` | Compass direction (0°=North) | `57.1` |
|`Position.AltitudeDegrees` | `Float [required]` | Angle above horizon (negative = below) | `8.8` |
|`Position.Direction` | `String [required]` | Cardinal direction abbreviation | `"ENE"` |
|`Position.DistanceKm` | `Float [required]` | Earth-Moon distance in km | `376559.9` |

</details>

</details>

----------------------------------------------------------------

#### MoonDaysDetailed

```MoonDaysDetailed``` is a structure that contains an array for each lunar day that falls on a given Earth day. Exists only if latitude and longitude are specified.

<details>
  <summary><strong>Table</strong></summary>

  | Response Variable | Type | Description | Example Value |
| :--- | :--- | :--- | :--- |
|`MoonDaysDetailed.Count` | `Integer [required]` | Number of lunar days this calendar day | `2` |
|`MoonDaysDetailed.Day` | `Array<Object> [required]` | Array of lunar day periods | - |
|`MoonDaysDetailed.Day[].Begin [optional]` | `String` | Start time of lunar day (ISO 8601) | `"2025-09-15T22:37:45+05:00"` |
|`MoonDaysDetailed.Day[].IsBeginExists [required]` | `Boolean` | True if start time is past/present | `true` |
|`MoonDaysDetailed.Day[].End [optional]` | `String` | End time of lunar day (ISO 8601) | `"2025-09-16T23:56:10+05:00"` |
|`MoonDaysDetailed.Day[].IsEndExists [required]` | `Boolean` | True if end time is past | `true`, `false` |

</details>

----------------------------------------------------------------

#### ZodiacDetailed

```ZodiacDetailed``` is a structure for determining which zodiac sign the moon is in on a given time interval, when it began and ended. It contains an array for each lunar day that falls on a given Earth day.

<details>
  <summary><strong>Table</strong></summary>

  | Response Variable | Type | Description | Example Value |
| :--- | :--- | :--- | :--- |
|`ZodiacDetailed.Count` | `Integer` | Number of zodiac signs this day | `1` |
|`ZodiacDetailed.Zodiac` | `Array<Object>` | Array of zodiac transit periods | - |
|`ZodiacDetailed.Zodiac[].Name` | `String` | Zodiac sign name | `"Gemini"` |
|`ZodiacDetailed.Zodiac[].NameLocalized` | `String` | Localized zodiac name | `"Близнецы"` |
|`ZodiacDetailed.Zodiac[].Emoji` | `String` | Zodiac emoji | `"♊"` |
|`ZodiacDetailed.Zodiac[].Begin` | `String` | Entry time into sign (ISO 8601) | `"2025-09-14T23:07:06+05:00"` |
|`ZodiacDetailed.Zodiac[].End` | `String` | Exit time from sign (ISO 8601) | `"2025-09-17T11:07:06+05:00"` |

</details>

----------------------------------------------------------------

#### MoonRiseAndSet

```MoonRiseAndSet``` is a structure for determining the moonrise, moonset and meridian on a given day.

<details>
  <summary><strong>Table</strong></summary>

| Response Variable | Type | Description | Example Value |
| :--- | :--- | :--- | :--- |
|`MoonRiseAndSet.Moonrise` | `Object of struct MoonPosition [optional]` | Moonrise position data. Exists only if IsMoonRise = true | - |
|`MoonRiseAndSet.Moonset` | `Object of struct MoonPosition [optional]` | Moonset position data. Exists only if IsMoonSet = true | - |
|`MoonRiseAndSet.Meridian` | `Object of struct MoonPosition [optional]` | Meridian position data, Exists only if IsMeridian = true | - |
|`MoonRiseAndSet.DistanceKm` | `Float [required]` | Approximate Earth-Moon distance | `379004.1` |
|`MoonRiseAndSet.IsMoonRise` | `Boolean [required]` | True if moonrise occurs at given day | `true` |
|`MoonRiseAndSet.IsMoonSet` | `Boolean [required]` | True if moonset occurs at given day | `true` |
|`MoonRiseAndSet.IsMeridian` | `Boolean [required]` | True if meridian transit occurs at given day | `true` |

MoonPosition structure:

<details>
  <summary><strong>Table</strong></summary>

| Response Variable | Type | Description | Example Value |
| :--- | :--- | :--- | :--- |
|`MoonPosition.Timestamp` | `Integer [required]` | Moonrise Unix timestamp | `1758048970` |
|`MoonPosition.TimeISO` | `String [required]` | Moonrise ISO time | `"2025-09-16T23:56:10+05:00"` |
|`MoonPosition.AzimuthDegrees` | `Float [required]` | Moonrise azimuth | `47.3` |
|`MoonPosition.AltitudeDegrees` | `Float [required]` | Moonrise altitude | `-0.6` |
|`MoonPosition.Direction` | `String [required]` | Moonrise direction | `"ENE"` |

</details>

</details>

----------------------------------------------------------------

#### Full response example:

Response of:
```GET /v1/moonPhaseDate?lang=ru&utc=5&latitude=51.1655&longitude=71.4272&year=2025&month=09&day=15&precision=5&hour=12&minute=0&second=0```

<details>
  <summary><strong>JSON</strong></summary>

```json
{
  "BeginDay": {
    "MoonDays": 22.53674,
    "Illumination": 49.42951,
    "Phase": {
      "Name": "Third quarter",
      "NameLocalized": "Последняя четверть",
      "Emoji": "🌗",
      "IsWaxing": false
    },
    "Zodiac": {
      "Name": "Gemini",
      "NameLocalized": "Близнецы",
      "Emoji": "♊"
    },
    "Position": {
      "Timestamp": 1757876400,
      "TimeISO": "2025-09-15T00:00:00+05:00",
      "AzimuthDegrees": 66.96877,
      "AltitudeDegrees": 17.50543,
      "Direction": "ENE",
      "DistanceKm": 373227.05417
    }
  },
  "CurrentState": {
    "MoonDays": 23.03674,
    "Illumination": 43.78443,
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
      "Timestamp": 1757919600,
      "TimeISO": "2025-09-15T12:00:00+05:00",
      "AzimuthDegrees": 279.37345,
      "AltitudeDegrees": 29.04833,
      "Direction": "W",
      "DistanceKm": 374869.93889
    }
  },
  "EndDay": {
    "MoonDays": 23.53674,
    "Illumination": 38.2726,
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
      "AzimuthDegrees": 57.08873,
      "AltitudeDegrees": 8.79646,
      "Direction": "ENE",
      "DistanceKm": 376559.88437
    }
  },
  "MoonDaysDetailed": {
    "Count": 2,
    "Day": [
      {
        "Begin": "2025-09-14T21:31:05+05:00",
        "IsBeginExists": true,
        "End": "2025-09-15T22:37:45+05:00",
        "IsEndExists": true
      },
      {
        "Begin": "2025-09-15T22:37:45+05:00",
        "IsBeginExists": true,
        "End": "2025-09-16T23:56:10+05:00",
        "IsEndExists": true
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
      "Timestamp": 1757957865,
      "TimeISO": "2025-09-15T22:37:45+05:00",
      "AzimuthDegrees": 42.31555,
      "AltitudeDegrees": -0.56667,
      "Direction": "NE"
    },
    "Moonset": {
      "Timestamp": 1757933213,
      "TimeISO": "2025-09-15T15:46:53+05:00",
      "AzimuthDegrees": 318.44328,
      "AltitudeDegrees": -0.56667,
      "Direction": "NW"
    },
    "Meridian": {
      "Timestamp": 1757900413,
      "TimeISO": "2025-09-15T06:40:13+05:00",
      "AzimuthDegrees": 180,
      "AltitudeDegrees": 67.1,
      "Direction": "S"
    },
    "DistanceKm": 375569.36571,
    "IsMoonRise": true,
    "IsMoonSet": true,
    "IsMeridian": true
  }
}
```

</details>

----------------------------------------------------------------

### GET /v1/moonPhaseCurrent

The method returns the Moon parameters for the current day and time. If the day or time is not specified, the current value for the unspecified fields is taken. If longitude and latitude are specified, the response will contain additional structures.

This is a synonym for the moonPhaseDate method without day and time arguments.

### Arguments

  | Parameter | Type | Description | Example Value |
| :--- | :--- | :--- |  :--- | 
|`utc` | `string [optional, default="UTC+0"]` | UTC in format `UTC+7`, `UTC+09:30`, `-3` | `UTC+4`
|`lang` | `string [optional, default="en"]` | Values available: ("en", "es", "fr", "de", "ru", "jp") | `es`
|`precision` | `int [optional, default=2]` | How many digits after ```.``` will be in output. Allowed range: [1, 20] | `5`
|`latitude` | `float [optional, default=none]` | Latitude of viewer's place. Used for moon position calculations: ```MoonDaysDetailed```, ```MoonRiseAndSet```, and ```MoonPosition``` object | `51.1655`
|`longitude` | `float [optional, default=none]` | Longitude of viewer's place. Used for moon position calculations: ```MoonDaysDetailed```, ```MoonRiseAndSet```, and ```MoonPosition``` object | `71.4272`

### Response

Response as [GET /v1/moonPhaseDate](https://github.com/prostraction/moon/#v1moonphasedate-response)

----------------------------------------------------------------

### GET /v1/moonPhaseTimestamp

The method returns the Moon parameters for the given timestamp. If it is not specified, the current value for the timestamp is taken. If longitude and latitude are specified, the response will contain additional structures.

This is a synonym for the moonPhaseDate method but with timestamp instead of date.

### Arguments

  | Parameter | Type | Description | Example Value |
| :--- | :--- | :--- |  :--- | 
|`utc` | `string [optional, default="UTC+0"]` | UTC in format `UTC+7`, `UTC+09:30`, `-3` | `UTC+4`
|`lang` | `string [optional, default="en"]` | Values available: ("en", "es", "fr", "de", "ru", "jp") | `es`
|`precision` | `int [optional, default=2]` | How many digits after ```.``` will be in output. Allowed range: [1, 20] | `5`
|`latitude` | `float [optional, default=none]` | Latitude of viewer's place. Used for moon position calculations: ```MoonDaysDetailed```, ```MoonRiseAndSet```, and ```MoonPosition``` object | `51.1655`
|`longitude` | `float [optional, default=none]` | Longitude of viewer's place. Used for moon position calculations: ```MoonDaysDetailed```, ```MoonRiseAndSet```, and ```MoonPosition``` object | `71.4272`
|`timestamp` | `int [optional, default=<current>]` | Timestamp for calculations | `1758045697`

### Response

Response as [GET /v1/moonPhaseDate](https://github.com/prostraction/moon/#v1moonphasedate-response)

----------------------------------------------------------------

### GET /v1/moonTableYear

The method returns the moon phases for the given year. The response contains an array for each month, each element of which contains the time of the new moon, first quarter, full moon, last quarter.

### Arguments

  | Parameter | Type | Description | Example Value |
| :--- | :--- | :--- |  :--- | 
|`utc` | `string [optional, default="UTC+0"]` | UTC in format `UTC+7`, `UTC+09:30`, `-3` | `UTC+4`
|`year` | `int [optional, default=<current year>]` | Format: YYYY Allowed range: [1, 9999] | `2025`

### /v1/moonTableYear Response

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

----------------------------------------------------------------

### GET /v1/moonTableCurrent

The method returns the moon phases for the current year. The response contains an array for each month, each element of which contains the time of the new moon, first quarter, full moon, last quarter.

### Arguments

  | Parameter | Type | Description | Example Value |
| :--- | :--- | :--- |  :--- | 
|`utc` | `string [optional, default="UTC+0"]` | UTC in format `UTC+7`, `UTC+09:30`, `-3` | `UTC+4`

### Response:

Response: as GET [/v1/moonTableYear](https://github.com/prostraction/moon#v1moontableyear-response)
