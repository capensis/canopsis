location /rundeck/ {
        set $upstream {{ RUNDECK_GRAILS_URL }} ;
        proxy_pass $upstream$request_uri;
        proxy_set_header X-Forwarded-Host $host;
	proxy_redirect {{ NGINX_URL }} $scheme://$host;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header User-Agent $http_user_agent;
        proxy_set_header X-Real-IP $remote_addr;

        error_page 502 @rundeck;
}

location @rundeck {
       default_type text/html;
       return 502 "<!DOCTYPE html><h1>Error 502</h1><h2>Rundeck is not reachable.</h2><h2>Check that the service is actually started.</h2>\n";
}
