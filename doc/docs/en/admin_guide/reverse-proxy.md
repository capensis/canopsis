Reverse proxy
=============

Reverse proxifying your connections to Canopsis can be useful to :

-   use standards ports for clients (http(s))
-   add ssl between clients and Canopsis
-   compress traffic


Reverse proxy with Apache2
==========================

For use another port for webserver, you can use Apache with `mod_proxy`

RHEL, CentOS / Debian
---------------------

Install Apache2

```bash
apt-get install apache2 or yum install httpd
```

Configure VHost (replace `<DNS_NAME>`)

```bash
<VirtualHost *:80>
    ProxyRequests Off
    ProxyPreserveHost On

    <Location />
        ProxyPass http://127.0.0.1:8082/
        ProxyPassReverse http://127.0.0.1:8082/

        SetOutputFilter DEFLATE
        SetEnvIfNoCase Request_URI \.(?:gif|jpe?g|png)$ no-gzip dont-vary
        SetEnvIfNoCase Request_URI \.(?:exe|t?gz|zip|gz2|sit|rar)$ no-gzip dont-vary
    </Location>
</VirtualHost>
```

Restart Apache and add service on boot

```bash
service apache2 restart or service httpd restart
update-rc.d apache2 defaults or chkconfig httpd on 
```

Reverse proxy with Nginx
========================

RHEL, CentOS / Debian
---------------------

Install nginx

```bash
apt-get install nginx or yum install nginx
```

**Don't forget to add epel repositories for RedHat/CentOS**

Configure VHost (replace `<DNS_NAME>`)

```bash
server {
    listen 80;
    server_name <DNS_NAME>;

    gzip on;
    gzip_disable "msie6";

    gzip_comp_level 6;
    gzip_min_length 1100;
    gzip_buffers 16 8k;
    gzip_proxied any;
    gzip_types
        text/plain
        text/css
        text/js
        text/xml
        text/javascript
        application/javascript
        application/x-javascript
        application/json
        application/xml
        application/xml+rss;

    location / {
        proxy_set_header X-Forwarded-Host $host;
        proxy_set_header X-Forwarded-Server $host;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_pass http://127.0.0.1:8082;
    }
}
```

Restart nginx and add service on boot

```bash
service nginx restart or service nginx restart
update-rc.d nginx defaults or chkconfig nginx on 
```
