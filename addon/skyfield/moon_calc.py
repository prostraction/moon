from skyfield import almanac
from skyfield.api import Loader, wgs84
from skyfield.almanac import meridian_transits
from datetime import datetime, timezone, timedelta

load = Loader('./skyfield-data')
eph = load('de421.bsp')
ts = load.timescale()

def get_azimuth_and_altitude(time, observer, target):
    astrometric = observer.at(time).observe(target)
    apparent = astrometric.apparent()
    alt, az, distance = apparent.altaz()
    
    AzimuthDegrees = az.degrees
    AltitudeDegrees = alt.degrees
    
    directions = [
        'N', 'NNE', 'NE', 'ENE', 
        'E', 'ESE', 'SE', 'SSE',
        'S', 'SSW', 'SW', 'WSW', 
        'W', 'WNW', 'NW', 'NNW'
    ]
    
    direction_index = int((AzimuthDegrees + 11.25) % 360 / 22.5)
    direction = directions[direction_index]
    
    return AzimuthDegrees, AltitudeDegrees, direction

def get_moon_position_at_time(lat, lon, timezone_hours, timezone_minutes, precision, year, month, day, hour, minute, second):
    total_timezone_offset = timezone_hours + timezone_minutes / 60.0
    offset = timedelta(hours=total_timezone_offset)
    tz = timezone(offset)

    local_dt = datetime(year, month, day, hour, minute, second, tzinfo=tz)
    utc_dt = local_dt - timedelta(hours=total_timezone_offset)
    
    location = wgs84.latlon(lat, lon)
    moon = eph['moon']
    earth = eph['earth']
    observer = earth + location
    
    time = ts.utc(utc_dt.year, utc_dt.month, utc_dt.day, 
                 utc_dt.hour, utc_dt.minute, utc_dt.second)
    
    azimuth, altitude, direction = get_azimuth_and_altitude(time, observer, moon)
    distance_km = earth.at(time).observe(moon).apparent().distance().km
    
    return {
        'Status': 'success',
        'Timestamp': int(local_dt.timestamp()),
        'AzimuthDegrees': round(azimuth, precision),
        'AltitudeDegrees': round(altitude, precision),
        'Direction': direction,
        'DistanceKm': round(distance_km, precision),
        'DateTime': {
            'Year': year,
            'Month': month,
            'Day': day,
            'Hour': hour,
            'Minute': minute,
            'Second': second
        },
        'TimezoneOffset': {
            'Hours': timezone_hours,
            'Minutes': timezone_minutes,
            'TotalHours': total_timezone_offset
        }
    }

def get_daily_moon_data(lat, lon, timezone_hours, timezone_minutes, precision, year, month, day):
    location = wgs84.latlon(lat, lon)
    moon = eph['moon']
    earth = eph['earth']
    observer = earth + location
    
    day_date = datetime(year, month, day)
    total_timezone_hours = timezone_hours + timezone_minutes / 60.0
    
    def get_meridian_time_and_direction(day_date, total_tz_hours):
        # searching for meridian (full 24 hours)
        t0 = ts.utc(day_date.year, day_date.month, day_date.day, 0 - total_tz_hours)
        t1 = ts.utc(day_date.year, day_date.month, day_date.day + 1, 0 - total_tz_hours)
        
        f = meridian_transits(eph, moon, location)
        times, events = almanac.find_discrete(t0, t1, f)
        
        # searching for upper meridian
        for time, event in zip(times, events):
            if event == 1:  # upper meridian
                local_time = time.utc_datetime()
                alt, az, distance = (observer.at(time).observe(moon).apparent().altaz())
                AltitudeDegrees = round(alt.degrees, 1)
                return local_time, 180.0, AltitudeDegrees, 'S'
        return None, None, None, None

    t0 = ts.utc(day_date.year, day_date.month, day_date.day, 0 - total_timezone_hours)
    t1 = ts.utc(day_date.year, day_date.month, day_date.day, 24 - total_timezone_hours)
    
    # horizon_degrees=0 make less pricise calculatuion, but it calculated to altitude = 0°
    f_rise_set = almanac.risings_and_settings(eph, moon, location)#, horizon_degrees=0)
    times_rise_set, events_rise_set = almanac.find_discrete(t0, t1, f_rise_set)
    
    rise_time, rise_azimuth, rise_altitude, rise_direction = None, None, None, None
    set_time, set_azimuth, set_altitude, set_direction = None, None, None, None
    
    for time, event in zip(times_rise_set, events_rise_set):
        local_time = time.utc_datetime()
        AzimuthDegrees, AltitudeDegrees, direction = get_azimuth_and_altitude(time, observer, moon)
        
        if event:  # rise
            rise_time = local_time
            rise_azimuth = AzimuthDegrees
            rise_altitude = AltitudeDegrees
            rise_direction = direction
        else:      # set
            set_time = local_time
            set_azimuth = AzimuthDegrees
            set_altitude = AltitudeDegrees
            set_direction = direction
    
    meridian_time, meridian_azimuth, meridian_altitude, meridian_direction = get_meridian_time_and_direction(day_date, total_timezone_hours)
    
    noon_t = ts.utc(day_date.year, day_date.month, day_date.day, 12)
    DistanceKm = earth.at(noon_t).observe(moon).apparent().distance().km
    
    # timestamp for golang
    moonrise_ts = int(rise_time.timestamp()) if rise_time else None
    moonset_ts = int(set_time.timestamp()) if set_time else None
    meridian_ts = int(meridian_time.timestamp()) if meridian_time else None

    moonrise_t = ts.utc(rise_time) if rise_time else None
    moonset_t  = ts.utc(set_time) if set_time else None
    meridian_t = ts.utc(meridian_time) if meridian_time else None

    distance_at_moonrise = earth.at(moonrise_t).observe(moon).apparent().distance().km if rise_time else None
    distance_at_moonset = earth.at(moonset_t).observe(moon).apparent().distance().km if set_time else None
    distance_at_meridian = earth.at(meridian_t).observe(moon).apparent().distance().km if meridian_time else None

    return {
        'Moonrise': {
            'Timestamp': moonrise_ts,
            'AzimuthDegrees': round(rise_azimuth, precision),
            'AltitudeDegrees': round(rise_altitude, precision),
            'Direction': rise_direction,
            'DistanceKm': round(distance_at_moonrise, precision),
        } if rise_time is not None else None,
        'Moonset': {
            'Timestamp': moonset_ts,
            'AzimuthDegrees': round(set_azimuth, precision),
            'AltitudeDegrees': round(set_altitude, precision),
            'Direction': set_direction,
            'DistanceKm': round(distance_at_moonset, precision),
        } if set_time is not None else None,
        'Meridian': {
            'Timestamp': meridian_ts,
            'AzimuthDegrees': round(meridian_azimuth, precision),
            'AltitudeDegrees': round(meridian_altitude, precision),
            'Direction': meridian_direction,
            'DistanceKm': round(distance_at_meridian, precision),
        } if meridian_time is not None else None,
        'IsMoonRise': rise_time is not None,
        'IsMoonSet': set_time is not None,
        'IsMeridian': meridian_time is not None,
    }

def calculate_moon_data(lat, lon, timezone_hours, timezone_minutes, precision, year, month, day=None):
    if day is not None:
        return get_daily_moon_data(lat, lon, timezone_hours, timezone_minutes, precision, year, month, day)
    else:     
        if month == 12:
            next_month = datetime(year + 1, 1, 1)
        else:
            next_month = datetime(year, month + 1, 1)
        
        last_day = (next_month - timedelta(days=1)).day 
        moon_data = []

        for day_num in range(1, last_day + 1):
            daily_data = get_daily_moon_data(lat, lon, timezone_hours, timezone_minutes, precision, year, month, day_num)
            moon_data.append(daily_data)
        return moon_data

def get_moon_data_response(lat, lon, timezone_hours, timezone_minutes, precision, year, month, day=None):
    data = calculate_moon_data(lat, lon, timezone_hours, timezone_minutes, precision, year, month, day)
    total_hours = timezone_hours + timezone_minutes / 60.0
    if total_hours >= 0:
        UtcOffset = f"UTC+{abs(timezone_hours):02d}:{timezone_minutes:02d}"
    else:
        UtcOffset = f"UTC-{abs(timezone_hours):02d}:{timezone_minutes:02d}"
    
    response = {
        'Status': 'success',
        'Parameters': {
            'Latitude': lat,
            'Longitude': lon,
            'TimezoneHours': timezone_hours,
            'TimezoneMinutes': timezone_minutes,
            'TotalTimezoneHours': round(total_hours, precision),
            'UtcOffset': UtcOffset,
            'Year': year,
            'Month': month
        },
        'Data': data
    }
    
    if day is not None:
        response['Parameters']['Day'] = day
        response['Range'] = 'single_day'
    else:
        response['Range'] = 'Full_month'
        response['DaysCount'] = len(data)
    return response
