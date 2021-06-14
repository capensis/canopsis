[![Build Status](https://travis-ci.org/bashrc666/ansible-influxdb.svg?branch=master)](https://travis-ci.org/bashrc666/ansible-influxdb)

InfluxDB
=========

Install of influxdb package


Compatibility
-------------
 - CentOS6
 - CentOS7

Requirements
------------

ansible >= 2.1.2

Role Variables
--------------

    influxdb_database_name: test
    influxdb_username: test
    influxdb_password: test
    influxdb_ip_address: 127.0.0.1

Dependencies
------------

None

Everything will be installed by roles if it's not installed

Example Playbook
----------------

```
- hosts: localhost
  remote_user: root
  vars:
    influxdb_database_name: test
    influxdb_username: test
    influxdb_password: test
    influxdb_ip_address: 127.0.0.1
  roles:
    - ansible-influxdb
```

License
-------

BSD

Author Information
------------------

gotoole/CAPENSIS - 2016
