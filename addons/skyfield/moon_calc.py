from skyfield import almanac
from skyfield.api import Loader, wgs84
from skyfield.almanac import meridian_transits
from datetime import datetime, timedelta

load = Loader('./skyfield-data')
eph = load('de421.bsp')
ts = load.timescale()

def get_daily_moon_data(lat, lon, timezone, year, month, day):
    location = wgs84.latlon(lat, lon)
    moon = eph['moon']
    earth = eph['earth']
    
    day_date = datetime(year, month, day)
    
    def get_meridian_time(day_date, tz_offset):
        # searching for meridian (full 24 hours)
        t0 = ts.utc(day_date.year, day_date.month, day_date.day, 0 - tz_offset)
        t1 = ts.utc(day_date.year, day_date.month, day_date.day + 1, 0 - tz_offset)
        
        f = meridian_transits(eph, moon, location)
        times, events = almanac.find_discrete(t0, t1, f)
        
        # searching for upper meridian
        for time, event in zip(times, events):
            if event == 1:  # upper meridian
                local_time = time.utc_datetime() + timedelta(hours=tz_offset)
                return local_time
        return None

    t0 = ts.utc(day_date.year, day_date.month, day_date.day, 0 - timezone)
    t1 = ts.utc(day_date.year, day_date.month, day_date.day, 24 - timezone)
    
    f_rise_set = almanac.risings_and_settings(eph, moon, location)
    times_rise_set, events_rise_set = almanac.find_discrete(t0, t1, f_rise_set)
    
    rise_time, set_time = None, None
    
    for time, event in zip(times_rise_set, events_rise_set):
        local_time = time.utc_datetime() + timedelta(hours=timezone)
        if event:  # rise
            rise_time = local_time
        else:      # set
            set_time = local_time
    
    meridian_time = get_meridian_time(day_date, timezone)
    
    noon_t = ts.utc(day_date.year, day_date.month, day_date.day, 12)
    distance_km = earth.at(noon_t).observe(moon).apparent().distance().km
    #illumination = almanac.fraction_illuminated(eph, 'moon', noon_t) * 100
    
    # timestamp for golang
    moonrise_ts = int(rise_time.timestamp()) if rise_time else None
    moonset_ts = int(set_time.timestamp()) if set_time else None
    meridian_ts = int(meridian_time.timestamp()) if meridian_time else None
    
    return {
        'date': day_date.strftime('%Y-%m-%d'),
        'moonrise': moonrise_ts,
        'moonset': moonset_ts,
        'meridian': meridian_ts,
        'isMoonRise': rise_time is not None,
        'isMoonSet': set_time is not None,
        'isMeridian': meridian_time is not None,
        'distance_km': round(distance_km, 1),
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