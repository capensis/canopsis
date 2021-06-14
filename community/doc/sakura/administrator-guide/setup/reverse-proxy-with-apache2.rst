.. _admin-setup-proxy-apache:

.. toctree::
   :maxdepth: 2
   :titlesonly:


Reverse proxy with Apache2
==========================

For use another port for webserver, you can use Apache with
``mod_proxy``

RHEL, CentOS / Debian
---------------------

Install Apache2

.. code-block:: bash

    apt-get install apache2 or yum install httpd


Configure VHost (replace ``<DNS_NAME>``)

.. code-block:: bash

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



Restart Apache and add service on boot

.. code-block:: bash

    service apache2 restart or service httpd restart
    update-rc.d apache2 defaults or chkconfig httpd on 

