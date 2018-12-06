# Utilisation du linkbuilder

## Objectif

Le linkbuilder permet de présenter des URL sur un bac à alarmes pour, par exemple :

*  déclarer un ticket ;
*  visualiser une consigne ou une procédure ;
*  rediriger vers un screenshot.

## Utilisation du basic builder

Le basic builder permet d'enrichir les alarmes avec un lien pouvant contenir des variables.

Ces variables sont celles contenues à la racine d'une entité.

### Configuration basique

Soit la configuration suivante :

```sh
$ cat basic_link_builder.json
{
    "basic_link_builder" : {
        "base_url" : "http://www.mesconsignes.local/?q={name}+{depends}+{type}",
	    "category" : "Consigne"
    }
}
```

(« category » est le nom sous lequel apparaitra le lien généré ; à défaut, « links » apparaitra)

En considérant l'entité associée :

```json
"entity": {
    "impact": ["srv-mail"],
    "name": "PING",
    "enable_history": [1518429148],
    "measurements": {},
    "enabled": true,
    "depends": ["nagios/Nagios4"],
    "infos": {},
    "_id": "PING/srv-mail",
    "type": "resource"
}
```

l'URL qui sera générée à la volée sur un bac à alarmes sera :

```json
"links": {
    "Ticketing": ["http://www.mesconsignes.fr/?q=PING+['nagios/Nagios4']+resource"]
}
```

### Basic alarm link builder

Une classe permettant de construire des liens à partir d'informations de l'alarme liées à l'entité ciblée est disponible. Par exemple, avec la configuration suivante :

```json
{
    "basic_alarm_link_builder" : {
        "base_url" : "http://www.mesconsignes.local/?q={alarm.v.component}"
    }
}
```

On va rechercher la valeur `component`, dans `v` de l'`alarm`. Si l'on ne précise pas **alarm**, la valeur sera recherchée dans l'entité.

### Mise en œuvre backend

La configuration préalablement établie doit être postée sur l'API de Canopsis.

**Phase d'authentification sur l'API :**

```sh
curl -POST http://x.x.x.x:8082/auth -d 'username=root&password=root' -vL -c canopsis_cookie
```

**Envoi de la configuration :**

```sh
curl -H "Content-Type: application/json" -X POST -d @basic_link_builder.json http://localhost:28082/api/v2/associativetable/link_builders_settings -b canopsis_cookie
```

Si une configuration existe déjà en base, remplacez `POST` par `PUT`.

Notez qu'un redémarrage du moteur `context-graph` ainsi que du `webserver` est nécessaire.

### Visualisation frontend

Dans un bac à alarmes, demandez l'affichage de la colonne `links`.
