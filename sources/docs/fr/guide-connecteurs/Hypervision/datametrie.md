# /!\ Documentation CAT /!\

## DatametrieManager

 * Project : `canopsis/canopsis` -> `sources/canopsis/datametrie/manager.py`

### Description

Accède à l'API de Datametrie, et transforme les informations recueillis en event Canopsis.

### Utilisation

Ligne d'import :

```python
from canopsis.datametrie.manager import DatametrieManager
```

Utilisation minimaliste :

```python
import logging
from canopsis.datametrie.manager import DatametrieManager

logger = logging.getLogger('datametrie')
conf = {
    Datametrie.CONF_CAT: {
        'user': '',
        'password': ''
    }
}

manager = DatametrieManager(config=conf, logger=logger)

manager.process_metrics()
```

Par défaut, les métriques sont processées au beat de l'engine Datametrie.

### Configuration

 * Ini : -> `etc/datametrie/datametrie.conf`

 - *user* : datametrie user account
 - *password* : datametrie password for previous user account
 - *prefix* : begin of the datametrie url (before verb)
 - *monitor_cache_length* : time in seconds before reloading the list of all known alarms
 - *excluded_alarm_types* (optional) : comma separated list with AlarmType to ignore (in upper case)
 - *http_proxy_url* (optional) : proxy url
 - *https_proxy_url* (optional) : proxy url (for https)

### Fonctions

 - **process_metrics()**: récupère les infos par l'API et génère les events correspondants


### Datametrie API

#### Alarm_status example

La description d'une alarme. Voila ce que peut donner l'api sur le service `/Get_Current_Alarms_All_Monitors` :

```json
{
	"MONITOR_ID": 123456,
	"ALARM_ID": 67108585,
	"ALARM_TYPE": "ORANGE",
	"ACK": 0,
	"ALARM_START_DATE": "05\/10\/2017 11:18:11",
	"ALARM_START_DATE_GMT": "05\/10\/2017 09:18:11",
	"RESULT_DATE": "05\/10\/2017 11:15:00",
	"RESULT_DATE_GMT": "05\/10\/2017 09:15:00",
	"ALARM_END_DATE": "",
	"ALARM_END_DATE_GMT": ""
}
```

#### Monitor example

La description d'un moniteur. Voila ce que peut donner l'api sur le service `/Get_Monitors` :

```json
{
	"MONITOR_ID": 123456,
	"MONITOR_NAME": "Bateaux",
	"CLIENT_ID": 987654,
	"MONITORING_DESCRIPTION": "Transaction",
	"IS_SHARED": 0,
	"MONITOR_STATUS": "STOPPED",
	"PERIODICITY": 10,
	"LOCATIONS": [
		{"SITE_ID": 1},
		{"SITE_ID": 61}
	],
	"PURCHASE_ORDER": 0
}
```
