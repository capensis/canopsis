Reverse proxy with Apache2
==========================

For use another port for webserver, you can use Apache with
``mod_proxy``

RHEL and CentOS
----------------

Install Apache2

::

    yum install httpd

Configure VHost (replace ``<DNS_NAME>``)

::

    echo "
    <VirtualHost 0.0.0.0:80>
     ProxyRequests Off
     ProxyPreservehost on
     ServerName  <DNS_NAME>
     ProxyPass / http://127.0.0.1:8082/
     ProxyPassReverse / http://127.0.0.1:8082/
    </VirtualHost>
    " > /etc/httpd/conf.d/canopsis.conf

Restart Apache and add service on boot

::

    service httpd restart
    chkconfig httpd on

