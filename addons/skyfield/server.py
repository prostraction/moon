from http.server import HTTPServer, BaseHTTPRequestHandler
from urllib.parse import urlparse, parse_qs
import json
from moon_calc import get_moon_data_response
from datetime import datetime

class QuietHTTPServer(HTTPServer):
    def handle_error(self, request, client_address):
        pass

class MoonRequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        parsed_url = urlparse(self.path)
        query_params = parse_qs(parsed_url.query)

        now = datetime.now()
        
        lat = float(query_params.get('lat', [51.08])[0])  # Astana
        lon = float(query_params.get('lon', [71.26])[0])  # Astana
        timezone = int(query_params.get('timezone', [5])[0]) # Astana
        year = int(query_params.get('year', [now.year])[0])
        month = int(query_params.get('month', [now.month])[0])
        
        
        day_param = query_params.get('day', [None])[0]
        day = int(day_param) if day_param is not None and day_param.isdigit() else None
        
        if not (1 <= month <= 12):
            self.send_error_response(400, "Month must be between 1 and 12")
            return
        
        if day is not None and not (1 <= day <= 31):
            self.send_error_response(400, "Day must be between 1 and 31")
            return
        
        response_data = get_moon_data_response(lat, lon, timezone, year, month, day)
        
        self.send_response(200 if response_data['status'] == 'success' else 500)
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
            'status': 'error',
            'message': message
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
    print('  GET /?lat=51.08&lon=71.26&timezone=5&year=2025&month=9')
    print('  GET /?lat=51.08&lon=71.26&timezone=5&year=2025&month=9&day=15')
    print('Parameters: lat, lon, timezone, year, month, day (day is optional)')
    print('Press Ctrl+C to stop the server')
    
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        httpd.shutdown()

if __name__ == "__main__":
    run_server(9997)