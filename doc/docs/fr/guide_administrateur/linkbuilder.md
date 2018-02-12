# Linkbuilder


## Objectif

Présenter des URLs sur un bac à alarmes pour par exemple : 

* Déclarer un ticket
* Visualiser une consigne/procédure
* Rediriger vers un screenshot


Le mécanisme de linkbuilder peut être enrichi pour des besoins spécifiques.
Voir [la section du guide de développement associée](../../en/developer_guide/canopsis_specs/link_generation/index.md)

## Utilisation du basic builder

Le basic builder permet d'enrichir les alarmes avec un lien pouvant contenir des variables.  
Ces variables sont celles contenues à la racine d'une entité. 

### Configuration basique

Soit la configuration suivante

````
$ cat basic_link_builder.json 
{
    "basic_link_builder" : {
        "base_url" : "http://www.mesconsignes.local/?q={name}+{depends}+{type}",
	    "category" : "Consigne"
    }
}
````

En considérant l'entité associée 

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

l'URL qui sera générée à la volée sur un bac à alarmes sera

```json
"links": {
    "Ticketing": ["http://www.mesconsignes.fr/?q=PING+[u'nagios/Nagios4']+resource"]
    }
```

### Mise en oeuvre backend

La configuration préalablement établie doit être postée sur l'API de Canopsis

**Phase d'authentification sur l'API**

````
curl -POST http://x.x.x.x:8082/auth -d 'username=root&password=root' -vL -c canopsis_cookie
````

**Post de la configuration**

````
curl -H "Content-Type: application/json" -X POST -d @basic_link_builder.json http://localhost:28082/api/v2/associativetable/link_builders_settings -b canopsis_cookie
````

Si une configuration existe déjà en base, remplacez **POST** par **PUT**


Notez qu'un redémarrage du moteur **context-graph** ainsi que du **webserver** est nécessaire.

### Visualisation frontend

Dans un bac à alarmes, demandez l'affichage de la colonne **links**