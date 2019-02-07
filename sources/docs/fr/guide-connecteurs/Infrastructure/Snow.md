# Snow vers Canopsis

!!! attention
    Ce connecteur n'est disponible que dans l'édition CAT de Canopsis.

## Introduction

Cette documentation détaille l'enrichissement des événements Canopsis depuis un serveur Service Now via le connecteur `snow2canopsis`.

## Fonctionnement

Le connecteur snow2canopsis se connecte une fois jour avec le serveur Service Now pour synchroniser Canopsis avec la CMDB de Service Now.

### Configuration

Dans le cadre d'une installation avec Docker de Canopsis, le fichier de configuration du connecteur est celui de ses variables d'environnement : `snow2canopsis.env`.

```ini
HTTPS_PROXY=HTTPS Proxy to connect to Service Now server
https_proxy=HTTPS Proxy to connect to Service Now server
CPS_SNOW_URL= base url of the snow server
CPS_SNOW_USERNAME= snow username
CPS_SNOW_PASSWORD= snow password
CPS_SNOW_FIELDS= list of parameters to update in canopsis from snow, separated with commas, example: "name,sys_id,usedfor"
CPS_HOST= Canopsis address (http://...)
CPS_PORT= Canopsis port
CPS_AUTH_KEY= Canopsis auth_key
CPS_CONNECTOR= The connector identifier in canopsis
CPS_CONNECTOR_NAME= The connector name in canopsis
```

#### Proxy HTTPS

Les 2 lignes `HTTPS_PROXY` et `https_proxy` ne doivent être utilisées que lorsqu'un Proxy HTTPS est sollicité pour se connecter au serveur Service Now.

Les 2 URL des Proxy doivent être indiquées juste aprs le `=` et sans guillemets.

Exemple :

```ini
HTTPS_PROXY=192.168.0.0.:8080
https_proxy=192.168.0.0.:8080
```

#### Authkey

La `auth_key` à rentrer est celle de l'utilisateur `root` de Canopsis. Elle sert au connecteur `snow2canopsis` à se connecter à l'instance de Canopsis.

#### CPS_HOST

Dans le cadre d'une installation avec Docker de Canopsis, la variable `CPS_HOST` doit avoir comme valeur `http://webserver`, le nom du service Web de Canopsis au sein du fichier `docker-compose.yml`.

```yaml
(...)
webserver:
  image: canopsis/canopsis-cat:${CANOPSIS_IMAGE_TAG}
(...)
  environment:
    - CPS_WEBSERVER=1
(...)
```

#### CPS_PORT

Dans le cadre d'une installation avec Docker de Canopsis, la variable `CPS_PORT` doit avoir comme valeur `8082`, le port du service Web de Canopsis au sein du fichier `docker-compose.yml`.

```yaml
(...)
webserver:
  image: canopsis/canopsis-cat:${CANOPSIS_IMAGE_TAG}
(...)
  environment:
    - CPS_WEBSERVER=1
  ports:
    - 80:8082
(...)
```

## Dépannage

### Vérification des logs

Dans le cadre d'une installation avec Docker de Canopsis, on peut consulter les logs du connecteur `snow2canopsis` :

```sh
$  docker-compose logs connector-snow2canopsis
```

#### Exemple de log positif

```sh
connector-snow2canopsis_1   | Running snow2canopsis
connector-snow2canopsis_1   | # Waiting for Canopsis API to be ready...
connector-snow2canopsis_1   | # Building snow hierarchy
connector-snow2canopsis_1   |
connector-snow2canopsis_1   | # Building entities from the hierarchy
connector-snow2canopsis_1   |
connector-snow2canopsis_1   | # Sending entities to Canopsis
connector-snow2canopsis_1   | 603 entities uploaded, waiting for completion
connector-snow2canopsis_1   | import c76aacac-yyz8-h7zi-7he6-810293029092 status: ongoing
connector-snow2canopsis_1   | import c76aacac-yyz8-h7zi-7he6-810293029092 status: ongoing
connector-snow2canopsis_1   | import c76aacac-yyz8-h7zi-7he6-810293029092 status: ongoing
connector-snow2canopsis_1   | import c76aacac-yyz8-h7zi-7he6-810293029092 status: ongoing
connector-snow2canopsis_1   | import c76aacac-yyz8-h7zi-7he6-810293029092 status: ongoing
connector-snow2canopsis_1   | import c76aacac-yyz8-h7zi-7he6-810293029092 status: ongoing
connector-snow2canopsis_1   | import c76aacac-yyz8-h7zi-7he6-810293029092 status: ongoing
connector-snow2canopsis_1   | import c76aacac-yyz8-h7zi-7he6-810293029092 status: ongoing
connector-snow2canopsis_1   | import c76aacac-yyz8-h7zi-7he6-810293029092 status: ongoing
connector-snow2canopsis_1   | import c76aacac-yyz8-h7zi-7he6-810293029092 status: ongoing
connector-snow2canopsis_1   | import c76aacac-yyz8-h7zi-7he6-810293029092 status: ongoing
connector-snow2canopsis_1   | import c76aacac-yyz8-h7zi-7he6-810293029092 status: ongoing
connector-snow2canopsis_1   | import c76aacac-yyz8-h7zi-7he6-810293029092 status: ongoing
connector-snow2canopsis_1   | import c76aacac-yyz8-h7zi-7he6-810293029092 status: ongoing
connector-snow2canopsis_1   | import c76aacac-yyz8-h7zi-7he6-810293029092 status: ongoing
connector-snow2canopsis_1   | import c76aacac-yyz8-h7zi-7he6-810293029092 status: ongoing
connector-snow2canopsis_1   | import c76aacac-yyz8-h7zi-7he6-810293029092 status: ongoing
connector-snow2canopsis_1   | import c76aacac-yyz8-h7zi-7he6-810293029092 status: done
```

#### Problème d'URL de Canopsis

Exemple d'erreur due à une mauvaise valeur pour `CPS_HOST`.

```sh
connector-snow2canopsis_1   |     raise MaxRetryError(_pool, url, error or ResponseError(cause))
connector-snow2canopsis_1   | urllib3.exceptions.MaxRetryError: HTTPConnectionPool(host='localhost', port=28082): Max retries exceeded with url: /?authkey=427c1756-24b0-11e9-a237-0242ac15000c (Caused by NewConnectionError('<urllib3.connection.HTTPConnection object at 0x7faa442555c0>: Failed to establish a new connection: [Errno 111] Connection refused',))
```

#### Problème de authkey

Exemple d'erreur due à une mauvaise valeur pour `CPS_AUTH_KEY` :

```sh
connector-snow2canopsis_1   | requests.exceptions.ConnectionError: HTTPConnectionPool(host='webserver', port=8082): Max retries exceeded with url: /?authkey=c6871fc4-2478-11e9-a647-0242ac140010 (Caused by NewConnectionError('<urllib3.connection.HTTPConnection object at 0x7ff78649f5c0>: Failed to establish a new connection: [Errno 111] Connection refused',))
```
