.. _admin-setup-proxy-nginx:

.. toctree::
   :maxdepth: 2
   :titlesonly:


Reverse proxy with Nginx
========================


RHEL, CentOS / Debian
---------------------

Install nginx

.. code-block:: bash

    apt-get install nginx or yum install nginx

**Don't forget to add epel repositories for RedHat/CentOS**


Configure VHost (replace ``<DNS_NAME>``)

.. code-block:: bash


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


Restart nginx and add service on boot

.. code-block:: bash

    service nginx restart or service nginx restart
    update-rc.d nginx defaults or chkconfig nginx on 

