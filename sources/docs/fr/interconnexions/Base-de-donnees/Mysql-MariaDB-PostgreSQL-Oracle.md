# connecteur de base de données SQL vers Canopsis / AMQP

Permet de transformer des résultats de requêtes SQL en évènements Canopsis via AMQP.

## Sommaire

[MySQL/MariaDB](#mysql-mariadb)  
[PostGreSQL](#postgresql)  
[Oracle](#oracle)  
[DB2](#db2)  
[MSSQL](#mssql)  

# Pré requis :

[kombu](https://pypi.python.org/pypi/kombu)  
[sqlalchemy](https://pypi.python.org/pypi/sqlalchemy)  

## Description

usage: `python connector-sql2canopsis.py [-h] -c CONFIG [-l LOGLEVEL]`

Il prend obligatoirement en paramètre un fichier de configuration INI (`-c CONFIG`) et optionnellement un degré d'alerte (`-l LOGLEVEL`).

## Installation

Python ≥ 2.7.10 doit être présent (attention, le script n'est *pas* compatible avec Python 3) et doit avoir été compilé avec l'option `--enable-unicode=ucs4`.

Lancer les commandes suivantes pour s'en assurer :
```shell
$ python -V
Python 2.7.14+
$ python -c 'import sys; print sys.maxunicode'
1114111
```

Si ces conditions ne sont pas réunies (c'est-à-dire si l'on a un autre résultat que `1114111`, ou une plus ancienne version de Python ou bien Python 3), cela signifie que l'environnement n'est pas compatible avec ce script.

Mais si tout est conforme, on peut alors installer pip et les fichiers de développement Python (note : sur CentOS, le dépôt EPEL doit être activé) :
```shell
# Pour Debian / Ubuntu
$ sudo apt-get install python-pip python-dev
# Pour Red Hat / CentOS
$ sudo yum install python-pip python-devel
```

On installe et on met en place virtualenv :
```shell
$ sudo pip install virtualenv
$ mkdir -p ~/venv/sql2canopsis
$ virtualenv ~/venv/sql2canopsis
```

On peut maintenant appeler pip pour installer les dépendances :
```shell
$ . ~/venv/sql2canopsis/bin/activate
(virtualenv) pip install -r requirements.txt
```

Vérifier le bon chargement des dépendances du script :
```shell
(virtualenv) python connector-sql2canopsis.py -h
usage: connector-sql2canopsis.py [-h] -c CONFIG [-l LOGLEVEL]
...
```

## Installation d'un « dialect »

On a aussi besoin d'installer un `dialect` sqlachemy (<https://docs.sqlalchemy.org/en/latest/dialects/index.html>) en fonction de la base de données ciblée.

### MySQL / MariaDB

Pour MySQL / MariaDB, il faut tout d'abord installer ses fichiers de développement :
```shell
# Pour Debian / Ubuntu
$ sudo apt-get install libmariadbclient-dev
# Ou, pour les versions plus anciennes (ex : Debian < 9)
$ sudo apt-get install libmysqlclient-dev

# Pour Red Hat / CentOS
$ sudo yum install mariadb-devel mariadb-libs
# Ou, pour les versions plus anciennes (ex : CentOS < 7)
$ sudo yum install mysql-devel mysql-libs
```

Installation du module Python nécessaire pour MySQL / MariaDB :
```shell
(virtualenv) pip install mysql-python
```

L'URL de connexion à ajouter au fichier de configuration sera de ce type (voir plus bas) :
```ini
[database]
url=mysql://user:password@mysql_host/database
```

### PostgreSQL

Pour PostgreSQL, certaines configurations nécessitent l'installation de ses fichiers de développement :
```shell
# Pour Debian / Ubuntu
$ sudo apt-get install postgresql-server-dev-all libpq-dev
# Pour Red Hat / CentOS
$ sudo yum install postgresql-libs postgresql-devel
```

```shell
(virtualenv) pip install psycopg2
```

L'URL de connexion à ajouter au fichier de configuration sera de ce type (voir plus bas) :
```ini
[database]
url=postgresql://user:password@postgresql_host/database
```

### Oracle

Prérequis : récupérer et installer Oracle Instant Client et son SDK (`instantclient-basic-linux` et `instantclient-sdk-linux`). Disponibles sur <https://www.oracle.com/technetwork/database/database-technologies/instant-client/overview/index.html> (attention, un compte est nécessaire).

Définir le répertoire d'installation d'Oracle Instant Client :
```shell
$ export ORACLE_HOME=/chemin/absolu/vers/oracle_instant_client
```

Trouver un fichier `libclntsh.so*` (éventuellement dans `$ORACLE_HOME` ou `$ORACLE_HOME/lib`) et y ajouter un lien symbolique de compatibilité :
```shell
$ cd /dossier/contenant/libclntsh
$ sudo ln -sf libclntsh.so.XXXXXX libclntsh.so
```

On peut alors installer le module pip :
```shell
(virtualenv) pip install cx_Oracle
```

L'URL de connexion à ajouter au fichier de configuration sera de ce type (voir plus bas) :
```ini
[database]
url=oracle://user:password@oracle_host:1521/sid
# Variante
#url=oracle+cx_oracle://user:password@tnsname
```

### Plus de détails sur les URL de connexion

Voir le lien suivant afin d'en savoir plus sur le paramètre `url`, si nécessaire : <https://docs.sqlalchemy.org/en/latest/core/engines.html#database-urls>

## Utilisation

### Configuration du fichier INI

Ce script attend *obligatoirement* un fichier `.ini`, passé en paramètre avec l'option `-c`.

Par exemple :
```shell
(virtualenv) python connector-sql2canopsis.py -c sql2canopsis.ini
```

Voici un exemple de fichier INI associé :
```ini
[database]
url=mysql://user:password@mysql_host/tickets
query=SELECT id, title, description, etablissement FROM tickets
use_last_value=false
fetch_size=1000
encoding=utf8

[amqp]
url=amqp://cpsrabbit:canopsis@localhost:5672/canopsis

[event]
connector.constant=edc2canopsis
connector_name.constant=instance1
event_type.constant=check
source_type.constant=resource
component.value=id
resource.value=etablissement
output.value=description
```

### Injection de valeurs avec `{sql_last_value}`

Le mot-clé `{sql_last_value}` peut être ajouté à la clause `WHERE` de la requête SQL afin d'y incorporer des résultats d'une requête précédente.

Pour cela :

 * Il faut tout d'abord activer cette fonctionnalité à l'aide de l'option `use_last_value=true`. Mettre cette valeur à `false` désactive cette option.
 * Il faut ensuite ajouter une clause `WHERE` à la requête SQL, en injectant la valeur de substitution `{sql_last_value}` à l'endroit souhaité (attention : sans espace supplémentaire, ceci est très important). Cette valeur de substitution sera remplacée par le contenu du fichier `last_value_retention_file` qui doit aussi être défini.
 * Définir, à l'aide de l'option `last_value_column`, le nom de la colonne dont la dernière valeur sera enregistrée dans le fichier `last_value_retention_file`.
 * Modifier le fichier défini dans `last_value_retention_file` afin d'y ajouter une valeur de départ. Par exemple, `1970-01-01 00:00:01` pour une date. Si le fichier défini dans `last_value_retention_file` n'existe pas, la valeur `0` sera utilisée.

Dans l'exemple précédent, la nouvelle configuration du bloc `[database]` deviendrait alors :

```ini
[database]
url=mysql://user:password@mysql_host/tickets
query=SELECT id, title, description, etablissement FROM tickets WHERE unix_timestamp(date) > unix_timestamp('{sql_last_value}')
use_last_value=true
last_value_column=date
last_value_retention_file=/un/chemin/absolu/est/recommandé/sql_last_value.ret
fetch_size=1000
encoding=utf8
```

## Métriques

Si des métriques sont nécessaires, celles-ci doivent être décrites dans le bloc `[event]` :
```ini
[event]
connector.constant=edc2canopsis
connector_name.constant=instance1
event_type.constant=check
source_type.constant=resource
component.value=id
resource.value=etablissement
output.value=description
metrique1.metrictype=GAUGE
metrique1.metric.value=nb
metrique1.metric.min=50
metrique2.metric.value=age
metrique2.metric.crit=5:
```

Dans cet exemple, la criticité « min » sera atteinte lorsque la valeur associée à `metrique1` (ici `nb`) dépassera 50.

De la même façon, la criticité « crit » sera atteinte lorsque la valeur associée à `metrique2` (ici `age`) sera en dessous de 5. Les « : » à la suite du nombre (`5:`) permettent en effet de préciser que la criticité sera atteinte lorsque la valeur sera *inférieure* au nombre indiqué.

Lorsque plusieurs métriques sont utilisées pour générer des évènements, le niveau de criticité le plus important sera retenu.

### DB2

Pour ce qui est de DB2, l'un des systèmes de gestion de base de données propriétaire d'IBM, veuillez vous rendre [ici](https://docs.sqlalchemy.org/en/latest/dialects/#production-ready), ou encore [ici](https://github.com/ibmdb/python-ibmdb).

### MSSQL

Pour ce qui est de DB2, le système de gestion de base de données propriétaire de Microsoft, veuillez vous rendre [ici](https://docs.sqlalchemy.org/en/latest/dialects/mssql.html), ou encore [ici](https://docs.sqlalchemy.org/en/latest/core/engines.html#microsoft-sql-server).
