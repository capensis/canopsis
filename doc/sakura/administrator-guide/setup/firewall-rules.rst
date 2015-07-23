.. _admin-setup-firewall:

Firewall rules
==============

This document describe what ports are used by Canopsis and must be authorized.

Port Listing
------------

+-------------+-------------------------------+-----------------------------+
| Tool        | Port                          | Description                 |
+-------------+-------------------------------+-----------------------------+
| MongoDB     | 27017                         | Database (or Replica)       |
+-------------+-------------------------------+-----------------------------+
| RabbitMQ    | 5672 (5673 in cluster)        | Messaging                   |
+-------------+-------------------------------+-----------------------------+
| RabbitMQ UI | 15672 (5674 in cluster)       | RabbitMQ administration UI  |
+-------------+-------------------------------+-----------------------------+
| Webserver   | 8082                          | RESTful API                 |
+-------------+-------------------------------+-----------------------------+
| HP-RabbitMQ | 5672 & 5674 (only in cluster) | Load balancer for RabbitMQ  |
+-------------+-------------------------------+-----------------------------+
| HP-WWW      | 80 & 443 (only in cluster)    | Load balancer for Webserver |
+-------------+-------------------------------+-----------------------------+

Details
-------

MongoDB
~~~~~~~

*MongoDB* must be allowed only for Canopsis nodes, any external user must not
access it directly.

In a mono-installation, just close the port, in a cluster-installation you need
to restrict access via the firewall.

RabbitMQ & HP-RabbitMQ
~~~~~~~~~~~~~~~~~~~~~~

*RabbitMQ* must be allowed for every Canopsis nodes and connectors.

It is used to transmit events into Canopsis, leave the port open in your firewall
and restrict the authorized IP addresses.

RabbitMQ UI & HP-RabbitMQ
~~~~~~~~~~~~~~~~~~~~~~~~~

This web interface must be used only by administrators and shouldn't be open to
users.

Just like previously, leave the port open with restrictions on IP addresses.

Webserver & HP-WWW
~~~~~~~~~~~~~~~~~~

The webserver, providing the RESTful Web API, must be authorized for every-one.

Leave the port open in the firewall, with no restrictions.
