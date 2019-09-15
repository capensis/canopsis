# Connecteur Cisco IoT FND `iotfnd2canopsis`

!!! attention
    Ce connecteur n'est disponible que dans l'édition CAT de Canopsis.

## Description

Convertit des évènements de supervision Cisco IoT FND en évènements Canopsis (via AMQP).

Il prend en paramètre un fichier de configuration INI (`-c CONFIG`).

## Installation

Ce script nécessite Python (compatible avec les versions 3.x) et quelques modules supplémentaires.

Cet environnement est mis en place à l'aide de `virtualenv` et `pip`.

### Installation sous Windows

Installer Python à l'aide de l'installateur disponible ici : <https://www.python.org/downloads/windows/>.

Puis, pour installer Pip, récupérer le script suivant : <https://bootstrap.pypa.io/get-pip.py>

Lancer un PowerShell, et exécuter ce script :
```powershell
PS C:\> python get-pip.py
Downloading/unpacking pip
Downloading/unpacking setuptools
Installing collected packages: pip, setuptools
Successfully installed pip setuptools
Cleaning up...
```

Une fois que pip est installé, on peut installer et mettre en place un virtualenv (en supposant que le script est stocké dans `C:\iotfnd2canopsis`) :
```powershell
PS C:\> pip install virtualenv
PS C:\> virtualenv c:\iotfnd2canopsis
```

On peut maintenant appeler pip pour installer les dépendances :
```powershell
PS C:\> cd iotfnd2canopsis
PS C:\> .\scripts\activate.bat
(virtualenv) pip install kombu ConfigParser zeep
```

Le connecteur est maintenant prêt à être utilisé :
```powershell
PS C:\> cd iotfnd2canopsis
PS C:\> .\scripts\activate.bat
(virtualenv) python iotfnd2canopsis.py -c iotfnd2canopsis.ini
```

### Installation sous Linux

**Note :** sous CentOS, EPEL est nécessaire.

On installe pip :
```sh
# Pour Debian
apt install python3-pip
# Pour CentOS
yum install python3-pip
```

On installe et on met en place virtualenv :
```sh
pip3 install virtualenv
mkdir -p ~/venv/iotfnd2canopsis
virtualenv ~/venv/iotfnd2canopsis
```

On peut maintenant appeler pip pour installer les dépendances :
```sh
. ~/venv/iotfnd2canopsis/bin/activate
(virtualenv) pip install ConfigParser kombu zeep
```

Le connecteur est maintenant prêt à être utilisé :
```sh
. ~/venv/iotfnd2canopsis/bin/activate
(virtualenv) python iotfnd2canopsis.py -c iotfnd2canopsis.ini
```

## Utilisation

### Configuration du fichier INI

Ce script attend *obligatoirement* un fichier `.ini`, passé en paramètre avec l'option `-c`.

Par exemple :
```sh
(virtualenv) python iotfnd2canopsis.py -c iotfnd2canopsis.ini
```

Voici un exemple de fichier INI associé :
```ini
[amqp]
url = amqp://cpsrabbit:canopsis@localhost:5672/canopsis

[iotfnd]
url = https://localhost/nbapi/issue?wsdl
username = USER
password = PASS
timestamp_file_path = /tmp/iotfnd2canopsis.timestamp
count = 1000
query = issue:down issueStatus:OPEN

[event]
connector.constant = iotfnd
connector_name.constant = iotfnd
event_type.constant = check
source_type.constant = component
component.value = eid
output.value = issueMessage
timestamp.value = issueLastUpdateTimestamp
state.value = state
```

Il contient l'URL du serveur AMQP à interroger et la liste des évènements à y envoyer. Ces événements sont créés à partir des *issues* que l'on va récupérer via l'API IOT FND.

Dans cet exemple, l'URL est définie dans le premier bloc : on se connecte au serveur `localhost` sur le port `5672`, avec les identifiants `cpsrabbit:canopsis`, sur la ressource `/canopsis`.

Dans le deuxième bloc, on peut définir différents filtres pour le paramètre `query`. Dans tous les cas, le `lastUpdateTime` est ajouté automatiquement à la requête et il est mis à jour après chaque appel à l'API. `lastUpdateTime` correspond au datetime défini dans le fichier `timestamp_file_path` ou à `time.now()`. Pour plus d'informations sur les filtres, [vous pouvez consulter la documentation de l'API IOT FND](https://www.cisco.com/c/en/us/td/docs/routers/connectedgrid/iot_fnd/api_guide/3_0/IoT-FND_NB_API/issue.html#87674 "North Bound API User Guide for the Cisco IoT Field Network Director").

Voici un exemple d'issue récupéré directement de l'API IOT FND :
```
{
    'eid': 'IERZwd-ZXJ0aWZpY2F0-57F3CC7+9429C098',
    'issueLastUpdateTime': datetime.datetime(2018, 1, 9, 16, 20, 32, 913000, tzinfo=<isodate.tzinfo.Utc object at 0x7f0697849518>),
    'issueMessage': 'Has been triggered.',
    'issueOccurTime': datetime.datetime(2017, 12, 30, 4, 3, 11, 366000, tzinfo=<isodate.tzinfo.Utc object at 0x7f0697849518>),
    'issueStatus': 'CLOSED',
    'issueTypeName': 'ContactAsserted',
    'meterId': None
}
```

### Exécution du script

Après avoir configuré le fichier INI, vérifiez que vous êtes bien dans le virtualenv et que vous utilisez une version de Python 3.X (``python --version``). Si c'est le cas, vous pouvez lancer le script avec la commande : `python iotfnd2canopsis.py -c iotfnd2canopsis.ini`. Après la fin du script, vous pouvez consulter Canopsis et vérifier que les issues de l'API IOT FND se retrouvent dans Canopsis.

Le script Python s'arrête après chaque utilisation, c'est-à-dire qu'il doit lancé à chaque fois que l'on souhaite exécuter à nouveau le traitement. Pour envoyer les issues vers Canopsis de manière régulière, il vous faudra utiliser un ordonnanceur (par exemple, `crontab` ou `Schtasks.exe`).
