from skyfield import almanac
from skyfield.api import Loader, wgs84
from skyfield.almanac import meridian_transits
from datetime import datetime, timedelta

load = Loader('./skyfield-data')
eph = load('de421.bsp')
ts = load.timescale()

def get_azimuth_and_altitude(time, observer, target):
    astrometric = observer.at(time).observe(target)
    apparent = astrometric.apparent()
    alt, az, distance = apparent.altaz()
    
    AzimuthDegrees = az.degrees
    AltitudeDegrees = alt.degrees
    
    directions = ['N', 'NE', 'E', 'SE', 'S', 'SW', 'W', 'NW']
    direction_index = int((AzimuthDegrees + 22.5) % 360 / 45)
    direction = directions[direction_index]
    
    return round(AzimuthDegrees, 1), round(AltitudeDegrees, 1), direction

def get_daily_moon_data(lat, lon, timezone, year, month, day):
    location = wgs84.latlon(lat, lon)
    moon = eph['moon']
    earth = eph['earth']
    observer = earth + location
    
    day_date = datetime(year, month, day)
    
    def get_meridian_time_and_direction(day_date, tz_offset):
        # searching for meridian (full 24 hours)
        t0 = ts.utc(day_date.year, day_date.month, day_date.day, 0 - tz_offset)
        t1 = ts.utc(day_date.year, day_date.month, day_date.day + 1, 0 - tz_offset)
        
        f = meridian_transits(eph, moon, location)
        times, events = almanac.find_discrete(t0, t1, f)
        
        # searching for upper meridian
        for time, event in zip(times, events):
            if event == 1:  # upper meridian
                local_time = time.utc_datetime() + timedelta(hours=tz_offset)
                alt, az, distance = (observer.at(time).observe(moon).apparent().altaz())
                AltitudeDegrees = round(alt.degrees, 1)
                return local_time, 180.0, AltitudeDegrees, 'S'
        return None, None, None, None

    t0 = ts.utc(day_date.year, day_date.month, day_date.day, 0 - timezone)
    t1 = ts.utc(day_date.year, day_date.month, day_date.day, 24 - timezone)
    
    # horizon_degrees=0 for more pricise calculatuion when altitude = 0°
    f_rise_set = almanac.risings_and_settings(eph, moon, location, horizon_degrees=0)
    times_rise_set, events_rise_set = almanac.find_discrete(t0, t1, f_rise_set)
    
    rise_time, rise_azimuth, rise_altitude, rise_direction = None, None, None, None
    set_time, set_azimuth, set_altitude, set_direction = None, None, None, None
    
    for time, event in zip(times_rise_set, events_rise_set):
        local_time = time.utc_datetime() + timedelta(hours=timezone)
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
    
    meridian_time, meridian_azimuth, meridian_altitude, meridian_direction = get_meridian_time_and_direction(day_date, timezone)
    
    noon_t = ts.utc(day_date.year, day_date.month, day_date.day, 12)
    DistanceKm = earth.at(noon_t).observe(moon).apparent().distance().km
    
    # timestamp for golang
    moonrise_ts = int(rise_time.timestamp()) if rise_time else None
    moonset_ts = int(set_time.timestamp()) if set_time else None
    meridian_ts = int(meridian_time.timestamp()) if meridian_time else None
    
    return {
        'Moonrise': {
            'Timestamp': moonrise_ts,
            'AzimuthDegrees': rise_azimuth,
            'AltitudeDegrees': rise_altitude,
            'Direction': rise_direction,
        } if rise_time is not None else None,
        'Moonset': {
            'Timestamp': moonset_ts,
            'AzimuthDegrees': set_azimuth,
            'AltitudeDegrees': set_altitude,
            'Direction': set_direction,
        } if set_time is not None else None,
        'Meridian': {
            'Timestamp': meridian_ts,
            'AzimuthDegrees': meridian_azimuth,
            'AltitudeDegrees': meridian_altitude,
            'Direction': meridian_direction,
        } if meridian_time is not None else None,
        'DistanceKm': round(DistanceKm, 1),
        'IsMoonRise': rise_time is not None,
        'IsMoonSet': set_time is not None,
        'IsMeridian': meridian_time is not None,
    }

def calculate_moon_data(lat, lon, timezone, year, month, day=None):
    if day is not None:
        # request for selected day:
        return get_daily_moon_data(lat, lon, timezone, year, month, day)
    else:
        # request for selected month:        
        if month == 12:
            next_month = datetime(year + 1, 1, 1)
        else:
            next_month = datetime(year, month + 1, 1)
        
        last_day = (next_month - timedelta(days=1)).day
        
        moon_data = []
        
        for day_num in range(1, last_day + 1):
            daily_data = get_daily_moon_data(lat, lon, timezone, year, month, day_num)
            moon_data.append(daily_data)
        
        return moon_data

def get_moon_data_response(lat, lon, timezone, year, month, day=None):
    try:
        data = calculate_moon_data(lat, lon, timezone, year, month, day)
        
        if timezone >= 0:
            UtcOffset = f"UTC+{timezone}"
        else:
            UtcOffset = f"UTC{timezone}"
        
        response = {
            'Status': 'success',
            'Parameters': {
                'Latitude': lat,
                'Longitude': lon,
                'Timezone': timezone,
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
        
    except Exception as e:
        if timezone >= 0:
            UtcOffset = f"UTC+{timezone}"
        else:
            UtcOffset = f"UTC{timezone}"
            
        error_response = {
            'Status': 'error',
            'Message': str(e),
            'Parameters': {
                'Latitude': lat,
                'Longitude': lon,
                'Timezone': timezone,
                'UtcOffset': UtcOffset,
                'Year': year,
                'Month': month
            }
        }
        
        if day is not None:
            error_response['Parameters']['Day'] = day
        
        return error_response