# Connecteur Zabbix vers Canopsis (connector-zabbix2canopsis)

Convertit des alertes issues de triggers Zabbix en évènements Canopsis.

## Prérequis

- Zabbix 4.2.6

## Introduction

Le connecteur Zabbix est un [script shell](#mise-en-place-du-script-zabbix2filebeatsh) configuré comme `media type` et déclenché par une `action` permettant de tracer dans un fichier de log les changements d'état des `triggers` configurés dans Zabbix. Le fichier de log est ensuite lu par Filebeat qui transfère le contenu à Logstash.

## Mapping

Canopsis utilise 4 états (0 : Info, 1 : Mineur , 2 : Majeur, 3 : Critique).

Zabbix utilise 6 états (0 : Not classified, 1 : Information, 2 : Warning, 3 : Average, 4 : High, 5 : Disaster).

Du fait de ces différences le mapping suivant a été décidé :

| Zabbix             | Canopsis     |
| ------------------ | ------------ |
| 0 - Not classified | 0 - Info     |
| 1 - Information    | 0 - Info     |
| 2 - Warning        | 1 - Mineur   |
| 3 - Average        | 1 - Mineur   |
| 4 - High           | 2 - Majeur   |
| 5  - Disaster      | 3 - Critical |

## Préparation à la mise en place de "zabbix2filebeat.sh"

Cette préparation est réalisée dans l'interface web de Zabbix.

### Création du media type "zabbix2filebeat"

Dans l'onglet Administration, cliquez sur "Media types" puis "Create media type".

Renseignez ensuite les informations ci-dessous :

* Name : `zabbix2filebeat`
* Type : `Script`
* Script name : `zabbix2filebeat.sh`
* Script parameters : `{ALERT.MESSAGE}`
* Enabled : `X`

### Ajouter le media à l'utilisateur admin

Toujours dans l'onglet Administration, cliquez sur "Users" puis "Admin".

Dans les propriétés de l'utilisateur cliquez sur l'onglet "Media" puis sur "Add" et renseignez le formulaire comme suit :

* Type : `zabbix2filebeat`
* Send to : `admin@admin`
* When active : `1-7,00:00-24:00`
* Use if severity :
    - Not classified : `X`
    - Information : `X`
    - Warning : `X`
    - Average : `X`
    - High : `X`
    - Disaster : `X`
* Enabled : `X`

### Création de l'action "zabbix2filebeat"

Dans l'onglet Configuration, cliquez sur "Actions" puis "Create action" et renseignez les différents onglets :

#### Onglet Action

* Name : `zabbix2filebeat`
* Enabled : `X`

#### Onglet Operations

* Default operation step duration : `1h`
* Default subject : `{TRIGGER.STATUS}: {TRIGGER.NAME}`
* Default message :
```
    ACTION.ID={ACTION.ID}
    ACTION.NAME={ACTION.NAME}
    DATE={DATE}
    ESC.HISTORY={ESC.HISTORY}
    EVENT.ACK.HISTORY={EVENT.ACK.HISTORY}
    EVENT.ACK.STATUS={EVENT.ACK.STATUS}
    EVENT.AGE={EVENT.AGE}
    EVENT.DATE={EVENT.DATE}
    EVENT.ID={EVENT.ID}
    EVENT.RECOVERY.DATE={EVENT.RECOVERY.DATE}
    EVENT.RECOVERY.ID={EVENT.RECOVERY.ID}
    EVENT.RECOVERY.STATUS={EVENT.RECOVERY.STATUS}
    EVENT.RECOVERY.TAGS={EVENT.RECOVERY.TAGS}
    EVENT.RECOVERY.TIME={EVENT.RECOVERY.TIME}
    EVENT.RECOVERY.VALUE={EVENT.RECOVERY.VALUE}
    EVENT.STATUS={EVENT.STATUS}
    EVENT.TAGS={EVENT.TAGS}
    EVENT.TIME={EVENT.TIME}
    EVENT.VALUE={EVENT.VALUE}
    HOST.CONN={HOST.CONN}
    HOST.DESCRIPTION={HOST.DESCRIPTION}
    HOST.DNS={HOST.DNS}
    HOST.HOST={HOST.HOST}
    HOSTNAME={HOSTNAME}
    HOST.IP={HOST.IP}
    HOST.NAME={HOST.NAME}
    HOST.PORT={HOST.PORT}
    INVENTORY.ALIAS={INVENTORY.ALIAS}
    INVENTORY.ASSET.TAG={INVENTORY.ASSET.TAG}
    INVENTORY.CHASSIS={INVENTORY.CHASSIS}
    INVENTORY.CONTACT={INVENTORY.CONTACT}
    PROFILE.CONTACT={PROFILE.CONTACT}
    INVENTORY.CONTRACT.NUMBER={INVENTORY.CONTRACT.NUMBER}
    INVENTORY.DEPLOYMENT.STATUS={INVENTORY.DEPLOYMENT.STATUS}
    INVENTORY.HARDWARE={INVENTORY.HARDWARE}
    PROFILE.HARDWARE={PROFILE.HARDWARE}
    INVENTORY.HARDWARE.FULL={INVENTORY.HARDWARE.FULL}
    INVENTORY.HOST.NETMASK={INVENTORY.HOST.NETMASK}
    INVENTORY.HOST.NETWORKS={INVENTORY.HOST.NETWORKS}
    INVENTORY.HOST.ROUTER={INVENTORY.HOST.ROUTER}
    INVENTORY.HW.ARCH={INVENTORY.HW.ARCH}
    INVENTORY.HW.DATE.DECOMM={INVENTORY.HW.DATE.DECOMM}
    INVENTORY.HW.DATE.EXPIRY={INVENTORY.HW.DATE.EXPIRY}
    INVENTORY.HW.DATE.INSTALL={INVENTORY.HW.DATE.INSTALL}
    INVENTORY.HW.DATE.PURCHASE={INVENTORY.HW.DATE.PURCHASE}
    INVENTORY.INSTALLER.NAME={INVENTORY.INSTALLER.NAME}
    INVENTORY.LOCATION={INVENTORY.LOCATION}
    PROFILE.LOCATION={PROFILE.LOCATION}
    INVENTORY.LOCATION.LAT={INVENTORY.LOCATION.LAT}
    INVENTORY.LOCATION.LON={INVENTORY.LOCATION.LON}
    INVENTORY.MACADDRESS.A={INVENTORY.MACADDRESS.A}
    PROFILE.MACADDRESS={PROFILE.MACADDRESS}
    INVENTORY.MACADDRESS.B={INVENTORY.MACADDRESS.B}
    INVENTORY.MODEL={INVENTORY.MODEL}
    INVENTORY.NAME={INVENTORY.NAME}
    INVENTORY.NOTES={INVENTORY.NOTES}
    INVENTORY.OOB.IP={INVENTORY.OOB.IP}
    INVENTORY.OOB.NETMASK={INVENTORY.OOB.NETMASK}
    INVENTORY.OOB.ROUTER={INVENTORY.OOB.ROUTER}
    INVENTORY.OS={INVENTORY.OS}
    PROFILE.OS={PROFILE.OS}
    INVENTORY.OS.FULL={INVENTORY.OS.FULL}
    INVENTORY.OS.SHORT={INVENTORY.OS.SHORT}
    INVENTORY.POC.PRIMARY.CELL={INVENTORY.POC.PRIMARY.CELL}
    INVENTORY.POC.PRIMARY.EMAIL={INVENTORY.POC.PRIMARY.EMAIL}
    INVENTORY.POC.PRIMARY.NAME={INVENTORY.POC.PRIMARY.NAME}
    INVENTORY.POC.PRIMARY.NOTES={INVENTORY.POC.PRIMARY.NOTES}
    INVENTORY.POC.PRIMARY.PHONE.A={INVENTORY.POC.PRIMARY.PHONE.A}
    INVENTORY.POC.PRIMARY.PHONE.B={INVENTORY.POC.PRIMARY.PHONE.B}
    INVENTORY.POC.PRIMARY.SCREEN={INVENTORY.POC.PRIMARY.SCREEN}
    INVENTORY.POC.SECONDARY.CELL={INVENTORY.POC.SECONDARY.CELL}
    INVENTORY.POC.SECONDARY.EMAIL={INVENTORY.POC.SECONDARY.EMAIL}
    INVENTORY.POC.SECONDARY.NAME={INVENTORY.POC.SECONDARY.NAME}
    INVENTORY.POC.SECONDARY.NOTES={INVENTORY.POC.SECONDARY.NOTES}
    INVENTORY.POC.SECONDARY.PHONE.A={INVENTORY.POC.SECONDARY.PHONE.A}
    INVENTORY.POC.SECONDARY.PHONE.B={INVENTORY.POC.SECONDARY.PHONE.B}
    INVENTORY.POC.SECONDARY.SCREEN={INVENTORY.POC.SECONDARY.SCREEN}
    INVENTORY.SERIALNO.A={INVENTORY.SERIALNO.A}
    PROFILE.SERIALNO={PROFILE.SERIALNO}
    INVENTORY.SERIALNO.B={INVENTORY.SERIALNO.B}
    INVENTORY.SITE.ADDRESS.A={INVENTORY.SITE.ADDRESS.A}
    INVENTORY.SITE.ADDRESS.B={INVENTORY.SITE.ADDRESS.B}
    INVENTORY.SITE.ADDRESS.C={INVENTORY.SITE.ADDRESS.C}
    INVENTORY.SITE.CITY={INVENTORY.SITE.CITY}
    INVENTORY.SITE.COUNTRY={INVENTORY.SITE.COUNTRY}
    INVENTORY.SITE.NOTES={INVENTORY.SITE.NOTES}
    INVENTORY.SITE.RACK={INVENTORY.SITE.RACK}
    INVENTORY.SITE.STATE={INVENTORY.SITE.STATE}
    INVENTORY.SITE.ZIP={INVENTORY.SITE.ZIP}
    INVENTORY.SOFTWARE={INVENTORY.SOFTWARE}
    INVENTORY.SOFTWARE.APP.A={INVENTORY.SOFTWARE.APP.A}
    INVENTORY.SOFTWARE.APP.B={INVENTORY.SOFTWARE.APP.B}
    INVENTORY.SOFTWARE.APP.C={INVENTORY.SOFTWARE.APP.C}
    INVENTORY.SOFTWARE.APP.D={INVENTORY.SOFTWARE.APP.D}
    INVENTORY.SOFTWARE.APP.E={INVENTORY.SOFTWARE.APP.E}
    INVENTORY.SOFTWARE.FULL={INVENTORY.SOFTWARE.FULL}
    INVENTORY.TAG={INVENTORY.TAG}
    INVENTORY.TYPE={INVENTORY.TYPE}
    INVENTORY.TYPE.FULL={INVENTORY.TYPE.FULL}
    INVENTORY.URL.A={INVENTORY.URL.A}
    INVENTORY.URL.B={INVENTORY.URL.B}
    INVENTORY.URL.C={INVENTORY.URL.C}
    INVENTORY.VENDOR={INVENTORY.VENDOR}
    ITEM.DESCRIPTION={ITEM.DESCRIPTION}
    ITEM.ID={ITEM.ID}
    ITEM.KEY={ITEM.KEY}
    TRIGGER.KEY={TRIGGER.KEY}
    ITEM.KEY.ORIG={ITEM.KEY.ORIG}
    ITEM.LASTVALUE={ITEM.LASTVALUE}
    ITEM.LOG.AGE={ITEM.LOG.AGE}
    ITEM.LOG.DATE={ITEM.LOG.DATE}
    ITEM.LOG.EVENTID={ITEM.LOG.EVENTID}
    ITEM.LOG.NSEVERITY={ITEM.LOG.NSEVERITY}
    ITEM.LOG.SEVERITY={ITEM.LOG.SEVERITY}
    ITEM.LOG.SOURCE={ITEM.LOG.SOURCE}
    ITEM.LOG.TIME={ITEM.LOG.TIME}
    ITEM.NAME={ITEM.NAME}
    ITEM.NAME.ORIG={ITEM.NAME.ORIG}
    ITEM.VALUE={ITEM.VALUE}
    PROXY.DESCRIPTION={PROXY.DESCRIPTION}
    PROXY.NAME={PROXY.NAME}
    TIME={TIME}
    TRIGGER.DESCRIPTION={TRIGGER.DESCRIPTION}
    TRIGGER.COMMENT={TRIGGER.COMMENT}
    TRIGGER.EVENTS.ACK={TRIGGER.EVENTS.ACK}
    TRIGGER.EVENTS.PROBLEM.ACK={TRIGGER.EVENTS.PROBLEM.ACK}
    TRIGGER.EVENTS.PROBLEM.UNACK={TRIGGER.EVENTS.PROBLEM.UNACK}
    TRIGGER.EVENTS.UNACK={TRIGGER.EVENTS.UNACK}
    TRIGGER.HOSTGROUP.NAME={TRIGGER.HOSTGROUP.NAME}
    TRIGGER.EXPRESSION={TRIGGER.EXPRESSION}
    TRIGGER.EXPRESSION.RECOVERY={TRIGGER.EXPRESSION.RECOVERY}
    TRIGGER.ID={TRIGGER.ID}
    TRIGGER.NAME={TRIGGER.NAME}
    TRIGGER.NAME.ORIG={TRIGGER.NAME.ORIG}
    TRIGGER.NSEVERITY={TRIGGER.NSEVERITY}
    TRIGGER.SEVERITY={TRIGGER.SEVERITY}
    TRIGGER.STATUS={TRIGGER.STATUS}
    STATUS={STATUS}
    TRIGGER.TEMPLATE.NAME={TRIGGER.TEMPLATE.NAME}
    TRIGGER.URL={TRIGGER.URL}
    TRIGGER.VALUE={TRIGGER.VALUE}
```
* Pause operations for suppressed problems : `X`
* Operations :
    -  Cliquez sur "New" et renseignez :
        - Operation details :
        - Steps : `1` - `0`
            - Step duration : `0`
            - Operation type : `Send message`
            - Send to User groups : `Zabbix administrators`
            - Send only to : `zabbix2filebeat`
            - Default message : `X`
            - Cliquez sur `Add`

#### Onglet Recovery operations

* Default subject  : `{TRIGGER.STATUS}: {TRIGGER.NAME}`
* Default message :
```
  ACTION.ID={ACTION.ID}
  ACTION.NAME={ACTION.NAME}
  DATE={DATE}
  ESC.HISTORY={ESC.HISTORY}
  EVENT.ACK.HISTORY={EVENT.ACK.HISTORY}
  EVENT.ACK.STATUS={EVENT.ACK.STATUS}
  EVENT.AGE={EVENT.AGE}
  EVENT.DATE={EVENT.DATE}
  EVENT.ID={EVENT.ID}
  EVENT.RECOVERY.DATE={EVENT.RECOVERY.DATE}
  EVENT.RECOVERY.ID={EVENT.RECOVERY.ID}
  EVENT.RECOVERY.STATUS={EVENT.RECOVERY.STATUS}
  EVENT.RECOVERY.TAGS={EVENT.RECOVERY.TAGS}
  EVENT.RECOVERY.TIME={EVENT.RECOVERY.TIME}
  EVENT.RECOVERY.VALUE={EVENT.RECOVERY.VALUE}
  EVENT.STATUS={EVENT.STATUS}
  EVENT.TAGS={EVENT.TAGS}
  EVENT.TIME={EVENT.TIME}
  EVENT.VALUE={EVENT.VALUE}
  HOST.CONN={HOST.CONN}
  HOST.DESCRIPTION={HOST.DESCRIPTION}
  HOST.DNS={HOST.DNS}
  HOST.HOST={HOST.HOST}
  HOSTNAME={HOSTNAME}
  HOST.IP={HOST.IP}
  HOST.NAME={HOST.NAME}
  HOST.PORT={HOST.PORT}
  INVENTORY.ALIAS={INVENTORY.ALIAS}
  INVENTORY.ASSET.TAG={INVENTORY.ASSET.TAG}
  INVENTORY.CHASSIS={INVENTORY.CHASSIS}
  INVENTORY.CONTACT={INVENTORY.CONTACT}
  PROFILE.CONTACT={PROFILE.CONTACT}
  INVENTORY.CONTRACT.NUMBER={INVENTORY.CONTRACT.NUMBER}
  INVENTORY.DEPLOYMENT.STATUS={INVENTORY.DEPLOYMENT.STATUS}
  INVENTORY.HARDWARE={INVENTORY.HARDWARE}
  PROFILE.HARDWARE={PROFILE.HARDWARE}
  INVENTORY.HARDWARE.FULL={INVENTORY.HARDWARE.FULL}
  INVENTORY.HOST.NETMASK={INVENTORY.HOST.NETMASK}
  INVENTORY.HOST.NETWORKS={INVENTORY.HOST.NETWORKS}
  INVENTORY.HOST.ROUTER={INVENTORY.HOST.ROUTER}
  INVENTORY.HW.ARCH={INVENTORY.HW.ARCH}
  INVENTORY.HW.DATE.DECOMM={INVENTORY.HW.DATE.DECOMM}
  INVENTORY.HW.DATE.EXPIRY={INVENTORY.HW.DATE.EXPIRY}
  INVENTORY.HW.DATE.INSTALL={INVENTORY.HW.DATE.INSTALL}
  INVENTORY.HW.DATE.PURCHASE={INVENTORY.HW.DATE.PURCHASE}
  INVENTORY.INSTALLER.NAME={INVENTORY.INSTALLER.NAME}
  INVENTORY.LOCATION={INVENTORY.LOCATION}
  PROFILE.LOCATION={PROFILE.LOCATION}
  INVENTORY.LOCATION.LAT={INVENTORY.LOCATION.LAT}
  INVENTORY.LOCATION.LON={INVENTORY.LOCATION.LON}
  INVENTORY.MACADDRESS.A={INVENTORY.MACADDRESS.A}
  PROFILE.MACADDRESS={PROFILE.MACADDRESS}
  INVENTORY.MACADDRESS.B={INVENTORY.MACADDRESS.B}
  INVENTORY.MODEL={INVENTORY.MODEL}
  INVENTORY.NAME={INVENTORY.NAME}
  INVENTORY.NOTES={INVENTORY.NOTES}
  INVENTORY.OOB.IP={INVENTORY.OOB.IP}
  INVENTORY.OOB.NETMASK={INVENTORY.OOB.NETMASK}
  INVENTORY.OOB.ROUTER={INVENTORY.OOB.ROUTER}
  INVENTORY.OS={INVENTORY.OS}
  PROFILE.OS={PROFILE.OS}
  INVENTORY.OS.FULL={INVENTORY.OS.FULL}
  INVENTORY.OS.SHORT={INVENTORY.OS.SHORT}
  INVENTORY.POC.PRIMARY.CELL={INVENTORY.POC.PRIMARY.CELL}
  INVENTORY.POC.PRIMARY.EMAIL={INVENTORY.POC.PRIMARY.EMAIL}
  INVENTORY.POC.PRIMARY.NAME={INVENTORY.POC.PRIMARY.NAME}
  INVENTORY.POC.PRIMARY.NOTES={INVENTORY.POC.PRIMARY.NOTES}
  INVENTORY.POC.PRIMARY.PHONE.A={INVENTORY.POC.PRIMARY.PHONE.A}
  INVENTORY.POC.PRIMARY.PHONE.B={INVENTORY.POC.PRIMARY.PHONE.B}
  INVENTORY.POC.PRIMARY.SCREEN={INVENTORY.POC.PRIMARY.SCREEN}
  INVENTORY.POC.SECONDARY.CELL={INVENTORY.POC.SECONDARY.CELL}
  INVENTORY.POC.SECONDARY.EMAIL={INVENTORY.POC.SECONDARY.EMAIL}
  INVENTORY.POC.SECONDARY.NAME={INVENTORY.POC.SECONDARY.NAME}
  INVENTORY.POC.SECONDARY.NOTES={INVENTORY.POC.SECONDARY.NOTES}
  INVENTORY.POC.SECONDARY.PHONE.A={INVENTORY.POC.SECONDARY.PHONE.A}
  INVENTORY.POC.SECONDARY.PHONE.B={INVENTORY.POC.SECONDARY.PHONE.B}
  INVENTORY.POC.SECONDARY.SCREEN={INVENTORY.POC.SECONDARY.SCREEN}
  INVENTORY.SERIALNO.A={INVENTORY.SERIALNO.A}
  PROFILE.SERIALNO={PROFILE.SERIALNO}
  INVENTORY.SERIALNO.B={INVENTORY.SERIALNO.B}
  INVENTORY.SITE.ADDRESS.A={INVENTORY.SITE.ADDRESS.A}
  INVENTORY.SITE.ADDRESS.B={INVENTORY.SITE.ADDRESS.B}
  INVENTORY.SITE.ADDRESS.C={INVENTORY.SITE.ADDRESS.C}
  INVENTORY.SITE.CITY={INVENTORY.SITE.CITY}
  INVENTORY.SITE.COUNTRY={INVENTORY.SITE.COUNTRY}
  INVENTORY.SITE.NOTES={INVENTORY.SITE.NOTES}
  INVENTORY.SITE.RACK={INVENTORY.SITE.RACK}
  INVENTORY.SITE.STATE={INVENTORY.SITE.STATE}
  INVENTORY.SITE.ZIP={INVENTORY.SITE.ZIP}
  INVENTORY.SOFTWARE={INVENTORY.SOFTWARE}
  INVENTORY.SOFTWARE.APP.A={INVENTORY.SOFTWARE.APP.A}
  INVENTORY.SOFTWARE.APP.B={INVENTORY.SOFTWARE.APP.B}
  INVENTORY.SOFTWARE.APP.C={INVENTORY.SOFTWARE.APP.C}
  INVENTORY.SOFTWARE.APP.D={INVENTORY.SOFTWARE.APP.D}
  INVENTORY.SOFTWARE.APP.E={INVENTORY.SOFTWARE.APP.E}
  INVENTORY.SOFTWARE.FULL={INVENTORY.SOFTWARE.FULL}
  INVENTORY.TAG={INVENTORY.TAG}
  INVENTORY.TYPE={INVENTORY.TYPE}
  INVENTORY.TYPE.FULL={INVENTORY.TYPE.FULL}
  INVENTORY.URL.A={INVENTORY.URL.A}
  INVENTORY.URL.B={INVENTORY.URL.B}
  INVENTORY.URL.C={INVENTORY.URL.C}
  INVENTORY.VENDOR={INVENTORY.VENDOR}
  ITEM.DESCRIPTION={ITEM.DESCRIPTION}
  ITEM.ID={ITEM.ID}
  ITEM.KEY={ITEM.KEY}
  TRIGGER.KEY={TRIGGER.KEY}
  ITEM.KEY.ORIG={ITEM.KEY.ORIG}
  ITEM.LASTVALUE={ITEM.LASTVALUE}
  ITEM.LOG.AGE={ITEM.LOG.AGE}
  ITEM.LOG.DATE={ITEM.LOG.DATE}
  ITEM.LOG.EVENTID={ITEM.LOG.EVENTID}
  ITEM.LOG.NSEVERITY={ITEM.LOG.NSEVERITY}
  ITEM.LOG.SEVERITY={ITEM.LOG.SEVERITY}
  ITEM.LOG.SOURCE={ITEM.LOG.SOURCE}
  ITEM.LOG.TIME={ITEM.LOG.TIME}
  ITEM.NAME={ITEM.NAME}
  ITEM.NAME.ORIG={ITEM.NAME.ORIG}
  ITEM.VALUE={ITEM.VALUE}
  PROXY.DESCRIPTION={PROXY.DESCRIPTION}
  PROXY.NAME={PROXY.NAME}
  TIME={TIME}
  TRIGGER.DESCRIPTION={TRIGGER.DESCRIPTION}
  TRIGGER.COMMENT={TRIGGER.COMMENT}
  TRIGGER.EVENTS.ACK={TRIGGER.EVENTS.ACK}
  TRIGGER.EVENTS.PROBLEM.ACK={TRIGGER.EVENTS.PROBLEM.ACK}
  TRIGGER.EVENTS.PROBLEM.UNACK={TRIGGER.EVENTS.PROBLEM.UNACK}
  TRIGGER.EVENTS.UNACK={TRIGGER.EVENTS.UNACK}
  TRIGGER.HOSTGROUP.NAME={TRIGGER.HOSTGROUP.NAME}
  TRIGGER.EXPRESSION={TRIGGER.EXPRESSION}
  TRIGGER.EXPRESSION.RECOVERY={TRIGGER.EXPRESSION.RECOVERY}
  TRIGGER.ID={TRIGGER.ID}
  TRIGGER.NAME={TRIGGER.NAME}
  TRIGGER.NAME.ORIG={TRIGGER.NAME.ORIG}
  TRIGGER.NSEVERITY={TRIGGER.NSEVERITY}
  TRIGGER.SEVERITY={TRIGGER.SEVERITY}
  TRIGGER.STATUS={TRIGGER.STATUS}
  STATUS={STATUS}
  TRIGGER.TEMPLATE.NAME={TRIGGER.TEMPLATE.NAME}
  TRIGGER.URL={TRIGGER.URL}
  TRIGGER.VALUE={TRIGGER.VALUE}
```
* Operations :
    - Cliquez sur "New" puis renseignez :
    - Operation details :
        - Operation type : `Send message`
        - Send to User groups : `Zabbix administrators`
        - Send only to : `zabbix2filebeat`
        - Default message : `X`

### Mise en place du script "zabbix2filebeat.sh"

#### Création du script

```bash
cat > /usr/lib/zabbix/alertscripts/zabbix2filebeat.sh << EOF
#!/bin/bash
LOG_PATH=/var/log/zabbix/notifications.log
ALERT_MESSAGE="${@%Q}"
OLD_IFS=${IFS}
MSG=''
IS_RECOVERY=0
STATE=0

# Parse alert_message as key value
######

while IFS=$'\n' read args;do
	IFS="${OLD_IFS}"
	set -- $(echo ${args%Q} | sed 's/=/ /;s/\r//g')
	key="${1//./_}"
	shift
	value="${@//\\/\\\\\\}"
	value="${value//\"/\\\\\"}"
	value="${value//\'/\\\'}"
	eval "${key%Q}=$'${value%Q}'"
done < <(echo "${ALERT_MESSAGE%Q}")

###################
# Translate/Create some values
###################

# test if we are in recovery
if [[ "${EVENT_RECOVERY_VALUE}" != '{EVENT.RECOVERY.VALUE}' ]];then
	IS_RECOVERY=1
fi

# declare TIMESTAMP depending of recovery or not
if [[ "${IS_RECOVERY}" -eq 1 ]];then
	TIMESTAMP="$(date --date="${EVENT_RECOVERY_DATE//.//} ${EVENT_RECOVERY_TIME}" +"%s")"
else
	TIMESTAMP="$(date --date="${EVENT_DATE//.//} ${EVENT_TIME}" +"%s")"
	case ${TRIGGER_NSEVERITY} in
		0|1)
			STATE=0
			;;
		2|3)
			STATE=1
			;;
		4)
			STATE=2
			;;
		5)
			STATE=3
			;;
	esac
fi

CONNECTOR_NAME=$(cat /etc/hostname)

##################
# Create JSON message
##################

MSG+='{'
# Canopsis
MSG+=$'"timestamp":"'${TIMESTAMP}'"'
MSG+=$',"connector":"zabbix"'
MSG+=$',"connector_name":"'${CONNECTOR_NAME}'"'
MSG+=$',"event_type":"check"'
MSG+=$',"source_type":"resource"'
MSG+=$',"component":"'${HOST_CONN}'"'
MSG+=$',"resource":"'${TRIGGER_NAME_ORIG}'"'
MSG+=$',"output":"'${TRIGGER_NAME}'"'
MSG+=$',"state":"'${STATE}'"'
# Original zabbix values asked by the final customer
MSG+=$',"ITEM.ID":"'${ITEM_ID}'"'
MSG+=$',"ITEM.DESCRIPTION":"'${ITEM_DESCRIPTION}'"'
MSG+=$',"ITEM.NAME":"'${ITEM_NAME}'"'
MSG+=$',"ITEM.VALUE":"'${ITEM_VALUE}'"'
MSG+=$',"HOST.DNS":"'${HOST_DNS}'"'
MSG+=$',"HOST.HOST":"'${HOST_HOST}'"'
MSG+=$',"EVENT.ACK.STATUS":"'${EVENT_ACK_STATUS}'"'
MSG+=$',"EVENT.STATUS":"'${EVENT_STATUS}'"'
MSG+=$'}'

# Print message
echo "${MSG%Q}" >> ${LOG_PATH}
EOF
```

#### Rendre exécutable le script

```shell
chmod +x /usr/lib/zabbix/alertscripts/zabbix2filebeat.sh
```

#### Rotation du log via logrotate

```shell
cat > /etc/logrotate.d/zabbix-notifications << EOF
/var/log/zabbix/notifications.log {
    daily
    rotate 30
    missingok
    notifempty
    compress
    delaycompress
}
EOF
```

#### Output du script

Le script crée et alimente le fichier de log suivant : `/var/log/zabbix/notifications.log`

Le contenu non minifié d'une ligne de log est de la forme suivante :
```json
{
	"timestamp": "1538486892",
	"connector": "zabbix",
	"connector_name": "zabbix1.plateforme-dev",
	"event_type": "check",
	"source_type": "resource",
	"component": "127.0.0.1",
	"resource": "Zabbix agent on {HOST.NAME} is unreachable for 5 minutes",
	"output": "Zabbix agent on Zabbix server is unreachable for 5 minutes",
	"state": "0",
	"ITEM.ID": "23287",
	"ITEM.DESCRIPTION": "The agent always returns 1 for this item. It could be used in combination with nodata() for availability check.",
	"ITEM.NAME": "Agent ping",
	"ITEM.VALUE": "Up (1)",
	"HOST.DNS": "",
	"HOST.HOST": "Zabbix server",
	"EVENT.ACK.STATUS": "No",
	"EVENT.STATUS": "PROBLEM"
}
```

### Filebeat

#### Installation

Pour Red Hat et CentOS :

```shell
curl -L -O https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-6.4.2-x86_64.rpm
rpm -ivh filebeat-6.4.2-x86_64.rpm
```

Pour Debian consultez la [page de téléchargement](https://www.elastic.co/fr/downloads/beats/filebeat) de Filebeat et la documentation de votre distribution.

#### Configuration

```shell
cat > /etc/filebeat/filebeat.yml << EOF
filebeat.prospectors:
- type: log
  paths:
    - /var/log/zabbix/notifications.log
output.logstash:
  hosts: ["X.X.X.X:5044"]
EOF
```

!!! attention
    Sur la dernière ligne, n'oubliez pas de remplacer `X.X.X.X` par l'adresse du serveur Logstash de la plateforme ciblée.

#### Activation du démarrage automatique du service

```shell
systemctl enable filebeat
```

#### Rechargement la configuration

```shell
systemctl restart filebeat
```

### Logstash

On utilisera Logstash pour traiter les informations en provenance de Zabbix et les transmettre à Canopsis.

#### Installation

L'installation et la configuration initiale de Logstash restent à votre charge et ne seront pas détaillées dans cette documentation.

#### Configuration

Vous trouverez ci-dessous le contenu du fichier de configuration du pipeline.

```sh
cat > /etc/logstash/conf.d/filebeat-pipeline.conf << EOF
input {
    beats {
        port => "5044"
        codec => "json"
    }
}

filter {
    # Conversion des champs en integer
        mutate {
          convert => ["[state]", "integer"]
          convert => ["[state_type]", "integer"]
          convert => ["[timestamp]", "integer"]
        }


        # construction de la clef de routage (routing key) nécessaire à Canopsis
        mutate {
            add_field => {"[@metadata][canopsis_rk_]" => "%{connector}.%{connector_name}.%{event_type}.%{source_type}.%{component}" }
        }

        if [source_type] == "resource" {
            mutate {
                add_field => {"[@metadata][canopsis_rk]" => "%{[@metadata][canopsis_rk_]}.%{resource}" }
            }
        } else {
            mutate {
                add_field => {"[@metadata][canopsis_rk]" => "%{[@metadata][canopsis_rk_]" }
            }
        }

        # nettoyage: suppression des champs et métadonnées maintenant inutiles
        mutate {
            remove_field => ["[@metadata][canopsis_rk_]", "message", "@timestamp", "input"]
        }


        # Ajouter un timestamp à l'event (convertir une date en timestamp) :

        #ruby{
        #  code =>"event.set('timestamp', event.get('@timestamp').to_i)"
        #}

}
output {
  stdout { codec => rubydebug }

  rabbitmq {
      host => "rabbitmq"
      vhost => canopsis
      user => "cpsrabbit"
      password => "canopsis"
      exchange => "canopsis.events"
      exchange_type => "topic"
      message_properties => { "content_type" => "application/json" }
      key => "%{[@metadata][canopsis_rk]}"
  }
}
EOF
```

Les évènements seront envoyés à Canopsis par le biais du bus AMQP de RabbitMQ.
