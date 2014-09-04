Reverse proxy with nginx
========================

For use another port for webserver, you can use Nginx as a reverse proxy

Debian/Ubuntu
--------------

Install Nginx

::

    apt-get install nginx

Configure VHost (replace ``<DNS_NAME>``)

::

    cat >/etc/nginx/sites-available/canopsis<<EOF
    server {
        listen 80;
        server_name <DNS_NAME>;

        access_log  /var/log/nginx/canopsis.access.log;
        error_log   /var/log/nginx/canopsis.error.log;

        location / { 
            proxy_redirect     off;

            # you need to change this to "https", if you set "ssl" directive to "on"
            proxy_set_header   X-FORWARDED_PROTO http;
            proxy_set_header   Host              $http_host;
            proxy_set_header   X-Real-IP         $remote_addr;

            proxy_read_timeout 300;
            proxy_connect_timeout 300;

            proxy_pass http://127.0.0.1:8082;
        }
    }
    EOF

Then enable the newly created vhost:

::

    ln -s /etc/nginx/sites-available/canopsis /etc/nginx/sites-enabled/

Finally, reload Nginx

::

    /etc/init.d/nginx reload

Now you can access canopsis using the ``<DNS_NAME>`` on port 80
