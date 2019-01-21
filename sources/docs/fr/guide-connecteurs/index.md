# Présentation

Voici le listing des connecteurs qui peuvent fournir des événements entrants à Canopsis.

# Sommaire

- [Présentation](#présentation)
- [Sommaire](#sommaire)
	- [Base de données](#base-de-données)
		- [Mysql, PostgeSQL, Oracle, DB2, MSSQL](#mysql-postgesql-oracle-db2-mssql)
	- [Infrastructure](#infrastructure)
		- [Logstash](#logstash)
		- [SNMPtrap](#snmptrap)
		- [SNMPtrap_Custom](#snmptrapcustom)
		- [send_event](#sendevent)
	- [Supervision](#supervision)
		- [Nagios, Icinga, Centreon](#nagios-icinga-centreon)
		- [Zabbix](#zabbix)
		- [librenms](#librenms)
		- [Shinken](#shinken)
	- [Hypervision](#hypervision)
		- [Datametrie](#datametrie)



## Base de données

### Mysql, PostgeSQL, Oracle, DB2, MSSQL
- [connector-sql2canopsis](Base-de-donnees/Mysql-MariaDB-PostgreSQL-Oracle.md)  

## Infrastructure

### Logstash
- [logstash2canopsis](Infrastructure/Logstash.md)  

### SNMPtrap
- [SNMPtrap - snmp2canopsis](Infrastructure/SNMPtrap.md)  

###  SNMPtrap_Custom
**Documentation CAT (Canopsis Administration Tools)**

- [SNMPtrap_custom - snmp2canopsis](Infrastructure/SNMPtrap_custom.md)  

### send_event
- [connector-send_event2canopsis](Infrastructure/send_event.md)  


## Supervision

### Nagios, Icinga, Centreon
- [connector-neb2canopsis](Supervision/Nagios-et-Icinga.md) (Nagios & Icinga1)  
- [connector-centreon-engine](Supervision/Centreon.md) (Centreon)  
- [connector-livestatus2canopsis](Supervision/Icinga2.md) (Icinga2)

### Zabbix
- [connector-zabbix2canopsis](Supervision/Zabbix.md)  

### librenms
- [connector-librenms2canopsis](Supervision/LibreNMS.md)  

### Shinken
- [connector-shinken2canopsis](Supervision/Shinken.md)  

## Hypervision

### Datametrie
**Documentation CAT (Canopsis Administration Tools)**

- [connector-datametrie](Hypervision/datametrie.md)  
