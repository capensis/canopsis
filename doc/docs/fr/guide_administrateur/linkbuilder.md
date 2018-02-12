# Linkbuilder


## Objectif

Présenter des URLs sur un bac à alarmes pour par exemple : 

* Déclarer un ticket
* Visualiser une consigne/procédure
* Rediriger vers un screenshot


Le mécanisme de linkbuilder peut êter enrichi pour des besoins spécifiques.

Voir [dev guide](../../en/developer_guide/canopsis_specs/link_generation/index.md)

## Utilisation du basic builder

Le basic builder permet d'enrichir les alarmes avec un lien pouvant contenir des variables.  
Ces variables sont celles contenues à la racine d'une entité. 

### Configuration simple

Construisez un fichier qui servira de payload à une requête POST.

````
$ cat basic_link_builder.json 
{
    "basic_link_builder" : {
        "base_url" : "http://www.mesconsignes.local/?q={name}+{depends}+{type}",
	    "category" : "Consigne"
    }
}
````

En considérant l'entité suivante 

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