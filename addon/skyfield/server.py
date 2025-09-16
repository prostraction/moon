from http.server import HTTPServer, BaseHTTPRequestHandler
from urllib.parse import urlparse, parse_qs
import json
from moon_calc import get_moon_data_response, get_moon_position_at_time
from datetime import datetime

class QuietHTTPServer(HTTPServer):
    def handle_error(self, request, client_address):
        pass

class MoonRequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        parsed_url = urlparse(self.path)
        path = parsed_url.path
        query_params = parse_qs(parsed_url.query)

        now = datetime.now()
        
        lat = float(query_params.get('lat', [0])[0])  # Astana
        lon = float(query_params.get('lon', [0])[0])  # Astana
        
        timezone_hours = int(query_params.get('utc_hours', [5])[0])  # Astana
        timezone_minutes = int(query_params.get('utc_minutes', [0])[0])
        
        precision = int(query_params.get('precision', [2])[0])
        if not (precision <= 20):
            precision = 20    
        if not (precision >= 0):
            precision = 2   
        
        if path == '/position':
            self.handle_position_request(lat, lon, timezone_hours, timezone_minutes, precision, query_params, now)
        elif path == '/daily' or path == '/monthly':
            self.handle_moon_data_request(lat, lon, timezone_hours, timezone_minutes, precision, query_params, now, path)
    
    def handle_position_request(self, lat, lon, timezone_hours, timezone_minutes, precision, query_params, now):
        try:
            year = int(query_params.get('year', [now.year])[0])
            month = int(query_params.get('month', [now.month])[0])
            day = int(query_params.get('day', [now.day])[0])
            hour = int(query_params.get('hour', [12])[0])
            minute = int(query_params.get('minute', [0])[0])
            second = int(query_params.get('second', [0])[0])

            if not (1 <= month <= 12):
                self.send_error_response(400, "month must be between 1 and 12")
                return
            
            if not (1 <= day <= 31):
                self.send_error_response(400, "day must be between 1 and 31")
                return
            
            if not (0 <= hour <= 23):
                self.send_error_response(400, "hour must be between 0 and 23")
                return
            
            if not (0 <= minute <= 59):
                self.send_error_response(400, "minute must be between 0 and 59")
                return
            
            if not (0 <= second <= 59):
                self.send_error_response(400, "second must be between 0 and 59")
                return
            
            response_data = get_moon_position_at_time(
                lat, lon, timezone_hours, timezone_minutes, precision,
                year, month, day, hour, minute, second
            )
            
            self.send_json_response(response_data)
            
        except ValueError as e:
            self.send_error_response(400, f"Invalid parameter format: {str(e)}")
        except Exception as e:
            self.send_error_response(500, f"Internal server error: {str(e)}")
    
    def handle_moon_data_request(self, lat, lon, timezone_hours, timezone_minutes, precision, query_params, now, path):
        try:
            year = int(query_params.get('year', [now.year])[0])
            month = int(query_params.get('month', [now.month])[0])
            
            if path == '/daily':
                day = int(query_params.get('day', [now.day])[0])
                if not (1 <= day <= 31):
                    self.send_error_response(400, "day must be between 1 and 31")
                    return
            else:
                day = None
            
            if not (1 <= month <= 12):
                self.send_error_response(400, "month must be between 1 and 12")
                return
            
            response_data = get_moon_data_response(lat, lon, timezone_hours, timezone_minutes, precision, year, month, day)
            self.send_json_response(response_data)
            
        except ValueError as e:
            self.send_error_response(400, f"Invalid parameter format: {str(e)}")
        except Exception as e:
            self.send_error_response(500, f"Internal server error: {str(e)}")
    
    def send_json_response(self, response_data):
        self.send_response(200 if response_data.get('Status') == 'success' else 500)
        self.send_header('Content-type', 'application/json')
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET')
        self.end_headers()
        self.wfile.write(json.dumps(response_data, indent=2, ensure_ascii=False).encode('utf-8'))
    
    def send_error_response(self, status_code, message):
        self.send_response(status_code)
        self.send_header('Content-type', 'application/json')
        self.send_header('Access-Control-Allow-Origin', '*')
        self.end_headers()
        
        error_response = {
            'Status': 'error',
            'Message': message
        }
        self.wfile.write(json.dumps(error_response).encode('utf-8'))
    
    def do_OPTIONS(self):
        self.send_response(200)
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type')
        self.end_headers()

    def log_message(self, format, *args):
        return

def run_server(port=9997):
    server_address = ('', port)
    httpd = QuietHTTPServer(server_address, MoonRequestHandler)
    print(f'Starting moon data server on port {port}...')
    print('Available endpoints:')
    print('  GET /position - Moon position at specific time')
    print('  GET /daily    - Daily moon data (rise/set/meridian)')
    print('  GET /monthly  - Monthly moon data')
    print('')
    print('Examples:')
    print('  Position: /position?lat=51.08&lon=71.26&utc_hours=5&utc_minutes=0&year=2025&month=9&day=15&hour=20&minute=30&second=0')
    print('  Daily:    /daily?lat=51.08&lon=71.26&utc_hours=5&utc_minutes=0&year=2025&month=9&day=15')
    print('  Monthly:  /monthly?lat=51.08&lon=71.26&utc_hours=5&utc_minutes=0&year=2025&month=9')
    print('')
    print('Parameters:')
    print('  lat, lon - coordinates')
    print('  utc_hours, utc_minutes - timezone offset')
    print('  year, month, day, hour, minute, second - date and time')
    print('')
    print('Press Ctrl+C to stop the server')
    
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        httpd.shutdown()
        print('\nServer stopped')

if __name__ == "__main__":
    run_server(9997)