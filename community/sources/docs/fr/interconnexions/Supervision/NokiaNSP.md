# Connecteur Nokia NSP `nokiansp2canopsis`

!!! info
    Ce connecteur n'est disponible que dans l'édition Pro de Canopsis.

## Description

Convertit des évènements de supervision [Nokia NSP](https://www.nokia.com/networks/products/network-services-platform/) (*Network Services Platform*) en évènements Canopsis.

Il envoie les alarmes du composant *Fault Management* de la solution Nokia NSP vers Canopsis via AMQP.

## Installation

Le script nécessite Python 3 (version 3.5 ou supérieure), ainsi que l'outil `virtualenv`.

### Sur CentOS 7

Sur CentOS, Python 3.6 peut être installé depuis le dépôt [EPEL](https://fedoraproject.org/wiki/EPEL).

```sh
yum install epel-release
yum install python36 python36-virtualenv
```

Dans le dossier `connector-nokiansp2canopsis`, mettre en place le *virtualenv* :

```sh
virtualenv-3.6 venv
```

### Sur Debian 9

Sur Debian 9, la version de Python 3 distribuée est une version 3.5.

```sh
apt install python3 virtualenv
```

Dans le dossier `connector-nokiansp2canopsis`, mettre en place le *virtualenv* :

```sh
virtualenv -p python3 venv
```

### Suite et fin de mise en œuvre (toutes distributions)

Activer le *virtualenv* et y installer les dépendances grâce à `pip` :

```sh
. venv/bin/activate
(venv) $ pip install pika requests
```

Le connecteur peut à présent être exécuté dans cet environnement.

## Utilisation

### Configuration

Le script du connecteur doit obligatoirement être invoqué avec un fichier de configuration au format INI, qu'il convient d'indiquer avec l'option `-c` / `--config`.

Exemple :

```sh
(venv) $ ./nokiansp2canopsis.py -c nokiansp2canopsis.ini
```

Exemple de fichier INI (livré dans le dépôt) :

```ini
[amqp]
url = amqp://cpsrabbit:canopsis@canopsis:5672/canopsis

[nokiansp]
host = nokiansp
username = USER
password = PASS
ssl_verify = True
timeout = 15
timestamp_file_path = ./nokiansp2canopsis.timestamp
alarms_file_path = ./nokiansp2canopsis.alarms

[event]
connector.constant = NOKIA_NSP
connector_name.constant = NSP_NFMT_HYP_API_ALRM
event_type.constant = check
source_type.constant = resource
; Direct lookup in Nokia NSP API fields
resource.value = objectFullName
output.value = probableCause
timestamp.value = lastTimeDetected
; Calulated fields
component.value = canopsisComponent
state.value = canopsisState

[states]
warning = 1
minor = 1
major = 2
critical = 3
```

Les valeurs de configuration peuvent ainsi être adaptées à l'environnement où le connecteur est intégré, pour les différentes sections suivantes.

#### Section `[amqp]`

Cette section contient une unique clé, `url`, pour indiquer les coordonnées du serveur AMQP auquel envoyer les évènements.

La syntaxe d'URI `amqp://` fait l'objet d'une [spécification](https://www.rabbitmq.com/uri-spec.html).

#### Section `[nokiansp]`

Les paramètres de cette section concernent la connexion à l'outil source, ici l'API Nokia NSP.

*  `host` : nom d'hôte ou adresse IP de l'instance Nokia NSP
*  `username`, `password` : authentification sur l'API Nokia NSP
*  `ssl_verify` : (`True` ou `False`) active ou non la vérification de validité du certificat SSL lors de la communication avec l'API
*  `timeout` : temps en secondes au bout duquel le connecteur interrompt la requête à l'API NSP en l'absence de réponse

Afin de filtrer les alarmes à envoyer à l'hyperviseur et pour détecter la résolution des alarmes NSP, le connecteur enregistre l'horodatage de la dernière alarme rencontrée ainsi que la liste des alarmes en cours.

*  le paramètre `timestamp_file_path` définit le chemin du fichier où sauvegarder le dernier horodatage ;
*  le paramètre `alarms_file_path` définit le chemin du fichier où sauvegarder les alarmes.

Sauf cas particulier, ces deux chemins n'ont pas besoin d'être changés.

#### Section `[event]`

Cette section définit l'ensemble des clés et leur contenu pour renseigner les évènements envoyés à Canopsis.

Chaque clé notée `<key>.constant` sera renseignée avec la valeur notée dans le fichier de configuration.

Chaque clé notée `<key>.value` sera renseignée dynamiquement d'après une information présente dans la structure de l'alarme obtenue de NSP.

Exemple :

```ini
output.value = probableCause
```

Traduction :

> Le champ `output` de l'évènement Canopsis sera rempli avec le contenu du champ `probableCause` présent dans la structure de l'alarme Nokia NSP.

Cas particulier : les champs calculés `canopsisState` et `canopsisComponent` sont construits dans le code du connecteur d'après des règles plus complexes.

#### Section `[states]`

Cette section permet de configurer la correspondance entre les noms des niveaux de criticité dans l'outil Nokia NSP et les valeurs numériques attendues par Canopsis pour représenter la criticité du check.

Niveaux Canopsis :

*  0 - INFO
*  1 - MINOR
*  2 - MAJOR
*  3 - CRITICAL.

Si un nom de criticité Nokia NSP n'est pas défini dans la configuration, le connecteur lui associera par défaut la valeur 0 (niveau `INFO`, c'est-à-dire OK).

### Exécution

Une fois la configuration renseignée, le connecteur peut être invoqué depuis le *virtualenv*.

```sh
(venv) $ ./nokiansp2canopsis.py -c nokiansp2canopsis.ini
```

### Planification

Le script est un programme dont l'exécution est ponctuelle. Son lancement doit être assuré par un planificateur.

Par exemple, pour un lancement par `cron`, un script comme ci-dessous peut être mis en place pour exécuter le programme dans l'environnement adéquat :

```sh
#!/bin/sh

cd /opt/connector-nokiansp2canopsis
. venv/bin/activate
python nokiansp2canopsis.py -c nokiansp2canopsis.ini
```

L'entrée en `crontab` pour déclenchement toutes les 5 minutes peut alors être exprimée ainsi, dans un fichier `/etc/cron.d/nokiansp2canopsis` :

```sh
# crontab for nokiansp2canopsis

*/5 * * * * user /opt/connector-nokiansp2canopsis/nokiansp2canopsis.sh
```

### Exemples d'exécution

Le script affiche en sortie les évènements envoyés à l'hyperviseur.

S'il arrive à joindre le service AMQP, il envoie au moins un évènement intitulé *heartbeat* indiquant la bonne exécution du processus ou, le cas échéant, l'erreur rencontrée lors des requêtes à l'API Nokia NSP.

Exemple de sortie avec *heartbeat* OK :

```
Event NOKIA_NSP.NSP_NFMT_HYP_API_ALRM.check.resource.NSP_NFMT_HYP_API_ALRM.heartbeat sent
{
  "connector": "NOKIA_NSP",
  "connector_name": "NSP_NFMT_HYP_API_ALRM",
  "event_type": "check",
  "source_type": "resource",
  "component": "NSP_NFMT_HYP_API_ALRM",
  "resource": "heartbeat",
  "state": 0,
  "output": "Connector excecution successful"
}
```

Exemple de sortie avec *heartbeat* KO :

```
Event NOKIA_NSP.NSP_NFMT_HYP_API_ALRM.check.resource.NSP_NFMT_HYP_API_ALRM.heartbeat sent
{
  "connector": "NOKIA_NSP",
  "connector_name": "NSP_NFMT_HYP_API_ALRM",
  "event_type": "check",
  "source_type": "resource",
  "component": "NSP_NFMT_HYP_API_ALRM",
  "resource": "heartbeat",
  "state": 3,
  "output": "Caught exception: HTTPSConnectionPool(host='139.178.72.122', port=443):
             Max retries exceeded with url: /rest-gateway/rest/api/v1/auth/token
             (Caused by ConnectTimeoutError(<urllib3.connection.VerifiedHTTPSConnection object at 0x7f2c58098240>,
                                            'Connection to 139.178.72.122 timed out. (connect timeout=15)'))"
}
```
