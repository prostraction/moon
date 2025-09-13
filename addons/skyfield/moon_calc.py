from skyfield import almanac
from skyfield.api import Loader, wgs84
from skyfield.almanac import meridian_transits
from datetime import datetime, timedelta

load = Loader('./skyfield-data')
eph = load('de421.bsp')
ts = load.timescale()

def get_azimuth(time, observer, target):
    """Вычисляет азимут объекта в градусах и направлении (N, NE, E, SE, S, SW, W, NW)"""
    astrometric = observer.at(time).observe(target)
    apparent = astrometric.apparent()
    alt, az, distance = apparent.altaz()
    
    azimuth_degrees = az.degrees
    directions = ['N', 'NE', 'E', 'SE', 'S', 'SW', 'W', 'NW']
    direction_index = int((azimuth_degrees + 22.5) % 360 / 45)
    direction = directions[direction_index]
    
    return round(azimuth_degrees, 1), direction

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
                # always (180°) for meridian
                return local_time, 180.0, 'S'
        return None, None, None

    t0 = ts.utc(day_date.year, day_date.month, day_date.day, 0 - timezone)
    t1 = ts.utc(day_date.year, day_date.month, day_date.day, 24 - timezone)
    
    f_rise_set = almanac.risings_and_settings(eph, moon, location)
    times_rise_set, events_rise_set = almanac.find_discrete(t0, t1, f_rise_set)
    
    rise_time, rise_azimuth, rise_direction = None, None, None
    set_time, set_azimuth, set_direction = None, None, None
    
    for time, event in zip(times_rise_set, events_rise_set):
        local_time = time.utc_datetime() + timedelta(hours=timezone)
        azimuth_degrees, direction = get_azimuth(time, observer, moon)
        
        if event:  # rise
            rise_time = local_time
            rise_azimuth = azimuth_degrees
            rise_direction = direction
        else:      # set
            set_time = local_time
            set_azimuth = azimuth_degrees
            set_direction = direction
    
    meridian_time, meridian_azimuth, meridian_direction = get_meridian_time_and_direction(day_date, timezone)
    
    noon_t = ts.utc(day_date.year, day_date.month, day_date.day, 12)
    distance_km = earth.at(noon_t).observe(moon).apparent().distance().km
    
    # timestamp for golang
    moonrise_ts = int(rise_time.timestamp()) if rise_time else None
    moonset_ts = int(set_time.timestamp()) if set_time else None
    meridian_ts = int(meridian_time.timestamp()) if meridian_time else None
    
    return {
        'date': day_date.strftime('%Y-%m-%d'),
        'moonrise': {
            'timestamp': moonrise_ts,
            'time': rise_time.strftime('%H:%M:%S') if rise_time else None,
            'azimuth_degrees': rise_azimuth,
            'direction': rise_direction,
        } if rise_time is not None else None,
        'moonset': {
            'timestamp': moonset_ts,
            'time': set_time.strftime('%H:%M:%S') if set_time else None,
            'azimuth_degrees': set_azimuth,
            'direction': set_direction,
        } if set_time is not None else None,
        'meridian': {
            'timestamp': meridian_ts,
            'time': meridian_time.strftime('%H:%M:%S') if meridian_time else None,
            'azimuth_degrees': meridian_azimuth,
            'direction': meridian_direction,
        } if meridian_time is not None else None,
        'distance_km': round(distance_km, 1),
        'isMoonRise': rise_time is not None,
        'isMoonSet': set_time is not None,
        'isMeridian': meridian_time is not None,
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
            utc_offset = f"UTC+{timezone}"
        else:
            utc_offset = f"UTC{timezone}"
        
        response = {
            'status': 'success',
            'parameters': {
                'latitude': lat,
                'longitude': lon,
                'timezone': timezone,
                'utc_offset': utc_offset,
                'year': year,
                'month': month
            },
            'data': data
        }
        
        if day is not None:
            response['parameters']['day'] = day
            response['range'] = 'single_day'
        else:
            response['range'] = 'full_month'
            response['days_count'] = len(data)
        
        return response
        
    except Exception as e:
        if timezone >= 0:
            utc_offset = f"UTC+{timezone}"
        else:
            utc_offset = f"UTC{timezone}"
            
        error_response = {
            'status': 'error',
            'message': str(e),
            'parameters': {
                'latitude': lat,
                'longitude': lon,
                'timezone': timezone,
                'utc_offset': utc_offset,
                'year': year,
                'month': month
            }
        }
        
        if day is not None:
            error_response['parameters']['day'] = day
        
        return error_response