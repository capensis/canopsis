# CORS configuration file for Canopsis
# https://developer.mozilla.org/docs/Web/HTTP/CORS

if ($request_method = 'OPTIONS') {
	add_header 'Access-Control-Allow-Origin' "$scheme://$canopsis_server_name";
	add_header 'Access-Control-Allow-Credentials' 'true';
	add_header 'Access-Control-Allow-Headers' 'Authorization,Cache-Control,Content-Type';
	add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, OPTIONS';
	add_header 'Access-Control-Max-Age' '31536000';
	add_header 'Content-Type' 'text/plain; charset=UTF-8';
	add_header 'Content-Length' 0;

	return 204;
}

if ($request_method = 'GET') {
	add_header 'Access-Control-Allow-Origin' "$scheme://$canopsis_server_name";
	add_header 'Access-Control-Allow-Credentials' 'true';
	add_header 'Access-Control-Allow-Headers' 'Authorization,Cache-Control,Content-Type';
}

if ($request_method = 'POST') {
	add_header 'Access-Control-Allow-Origin' "$scheme://$canopsis_server_name";
	add_header 'Access-Control-Allow-Credentials' 'true';
	add_header 'Access-Control-Allow-Headers' 'Authorization,Cache-Control,Content-Type';
}

if ($request_method = 'PUT') {
	add_header 'Access-Control-Allow-Origin' "$scheme://$canopsis_server_name";
	add_header 'Access-Control-Allow-Credentials' 'true';
	add_header 'Access-Control-Allow-Headers' 'Authorization,Cache-Control,Content-Type';
}
