# Python send_event connector to Canopsis / AMQP

Un connecteur universel qui permet de publier des évènements dans Canopsis (utilisé par des ordonnanceurs, des scripts, des applications maison…).

## Pré requis

[Kombu](https://pypi.org/project/kombu/)

## Description

usage: `python sendevent2canopsisamqp.py [-h] -c CONFIG [-p PARAMS_JSON]`

Ce script envoie les évènements vers l'instance Canopsis via AMQP.

Il prend en paramètre un fichier de configuration INI (`-c CONFIG`), et optionnellement des valeurs JSON complémentaires (`-p PARAMS_JSON`).

## Installation

Ce script nécessite Python (compatible avec les versions 2.x et 3.x) et quelques modules supplémentaires.

Cet environnement est mis en place à l'aide de virtualenv et pip.

### Installation sous Windows

Installer Python à l'aide de l'installateur disponible ici : https://www.python.org/downloads/windows/

Puis, pour installer Pip, récupérer le script suivant : https://bootstrap.pypa.io/get-pip.py

Lancer un PowerShell, et exécuter ce script :
```powershell
PS C:\> python get-pip.py
Downloading/unpacking pip
Downloading/unpacking setuptools
Installing collected packages: pip, setuptools
Successfully installed pip setuptools
Cleaning up...
```

Une fois que pip est installé, on peut installer et mettre en place un virtualenv (en supposant que le script est stocké dans `C:\sendeventcanopsis`) :
```powershell
PS C:\> pip install virtualenv
PS C:\> virtualenv c:\sendeventcanopsis
```

On peut maintenant appeler pip pour installer les dépendances :
```powershell
PS C:\> cd sendeventcanopsis
PS C:\> .\scripts\activate.bat
(virtualenv) pip install kombu ConfigParser
```

Le connecteur est maintenant prêt à être utilisé :
```powershell
PS C:\> cd sendeventcanopsis
PS C:\> .\scripts\activate.bat
(virtualenv) python sendevent2canopsisamqp.py
```

### Installation sous Linux

On installe pip :
```shell
# Pour Debian / Ubuntu
$ sudo apt-get install python-pip
# Pour Red Hat / CentOS
$ sudo yum install python-pip
```

On installe et on met en place virtualenv :
```shell
$ pip install virtualenv
$ mkdir -p ~/venv/sendeventcanopsis
$ virtualenv ~/venv/sendeventcanopsis
```

On peut maintenant appeler pip pour installer les dépendances :
```shell
$ . ~/venv/sendeventcanopsis/bin/activate
(virtualenv) pip install ConfigParser kombu
```

Le connecteur est maintenant prêt à être utilisé :
```shell
$ . ~/venv/sendeventcanopsis/bin/activate
(virtualenv) python sendevent2canopsisamqp.py
```

## Utilisation

### Configuration du fichier INI

Ce script attend *obligatoirement* un fichier `.ini`, passé en paramètre avec l'option `-c`.

Par exemple :
```shell
(virtualenv) python sendevent2canopsisamqp.py -c sendevent2canopsisamqp.ini
```

Voici un exemple de fichier INI associé :
```ini
[amqp]
url=amqp://cpsrabbit:canopsis@localhost:5672/canopsis

[event]
connector.constant=toto
connector_name.constant=toto
event_type.constant=check
source_type.constant=resource
component.constant=localhost
resource.value=res
output.constant=output
state.constant=2
```

Il contient l'URL du serveur AMQP à interroger et la liste des évènements à y envoyer.

Dans cet exemple, l'URL est définie dans le premier bloc : on se connecte au serveur `localhost` sur le port `5672`, avec les identifiants `cpsrabbit:canopsis`, sur la ressource `/canopsis`.

Concernant les évènements (définis dans le bloc `[event]`) :

  * `connector.constant=toto` signifie qu'un évènement `connector` sera envoyé avec la valeur `toto`.
    - De la même façon, `event_type.constant=check` signifie qu'un évènement `event_type` sera envoyé avec la valeur `check`.
  * L'utilisation de **`.value`**, telle que dans `resource.value=res`, signifie que l'évènement `resource` sera renseigné à partir de la valeur `res` passée en ligne de commande (voir ci-dessous).
  * Il est aussi possible d'utiliser des regex, voir plus bas.

### Passage de valeurs supplémentaires en paramètres

Si le fichier INI contient au moins une valeur du type `.value`, alors le paramètre `-p` doit alors être obligatoirement renseigné lors de l'appel au script. Ceci permet d'injecter des valeurs dynamiques, sans devoir modifier le fichier INI pour chaque nouvelle valeur.

Ainsi, dans le cas précédent, où l'on avait une valeur `resource.value=res`, on doit renseigner une valeur `res` en paramètre :

```shell
(virtualenv) python sendevent2canopsisamqp.py -c sendevent2canopsisamqp.ini -p "{\"res\": \"test\"}"
```

Note : le paramètre `-p` attend un tableau JSON en argument. Par exemple, si plusieurs évènements doivent être précisés :

```shell
(virtualenv) python sendevent2canopsisamqp.py -c sendevent2canopsisamqp.ini -p "{\"param1\": \"valeur\", \"param2\": 42}"
```

### Utilisation des regex

Des valeurs **`.regex`** peuvent aussi être utilisées dans le fichier INI.

L'évènement aura pour valeur toute chaîne correspondant à la regex donnée.

De cette façon, l'évènement suivant :
```ini
[event]
component.regex=my_(.*)
```

transformera `component="my_first_value"` en `component="first_value"`.
