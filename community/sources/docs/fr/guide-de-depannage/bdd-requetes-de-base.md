# Requêtes en base

## Healthcheck

### 1) Connexion à la BDD

Connexion en ligne de commande :
```
# mongo canopsis -u cpsmongo -p MOT_DE_PASSE --host XXXX
MongoDB shell version v3.4.13
MongoDB server version: 3.4.13
rs0:PRIMARY> 
```

Ou sinon, avec [Robo3T](https://robomongo.org) (anciennement RoboMongo).
Dans Robo3T, les requêtes doivent être validées avec Ctrl + Entrée afin d'être exécutées.

**Attention :** MongoDB, et plus particulièrement Robo3T, sont sensibles à la véritable casse du hostname des nœuds MongoDB.

### 2) Identifier un évènement suivant le nom de son composant, sa ressource et son statut.

**Rappels**

Pour comprendre le fonctionnement interne de Canopsis, il est important de comprendre le modèle de données associé.  
Pour rappel, les entités sont stockées dans la collection `default_entities` tandis que les alarmes sont stockées dans `periodical_alarm`.

Prenons l'exemple de l'entité suivante :

* component : `NOM_MACHINE`
* resource : `Traffic - eth0`
* connector : `centreon`
* connector_name : `NOM_CONN`

Pour rechercher cette entité en base de données, vous pouvez exécuter la requête suivante :

```js
db.default_entities.find({"_id": "Traffic - eth0/NOM_MACHINE"}).pretty()
{
        "_id" : "Traffic - eth0/NOM_MACHINE",
        "name" : "Traffic - eth0",
        "impact" : [
                "NOM_MACHINE"
        ],
        "depends" : [
                "centreon/NOM_CONN"
        ],
        "enable_history" : [
                NumberLong(1527085105)
        ],
        "measurements" : null,
        "enabled" : true,
        "infos" : {
                "hostgroups" : {
                        "name" : "hostgroups",
                        "description" : "",
                        "value" : "XXXX"
                },
                "component_alias" : {
                        "name" : "component_alias",
                        "description" : "",
                        "value" : "XXXX"
                }
        },
        "type" : "resource"
}
```

**Recherche**

Si vous souhaitez rechercher les **alarmes** associées à cette entité, il faut requêter la collection `periodical_alarm`.

```js
db.periodical_alarm.find({"v.component" : "NOM_MACHINE", "v.resource" : "Traffic - eth0"}).pretty()
{
        "_id" : "XXXX",
        "t" : NumberLong(1524667917),
        "d" : "Traffic - eth0/NOM_MACHINE",
        "v" : {
                "state" : {
                        "_t" : "stateinc",
                        "t" : NumberLong(1524667892),
                        "a" : "NOM_CONN",
                        "m" : "ERROR: Interface Status Request : No response from remote host \"IP_MACHINE\"",
                        "val" : NumberLong(1)
                },
                "status" : {
                        "_t" : "statusinc",
                        "t" : NumberLong(1524667892),
                        "a" : "NOM_CONN",
                        "m" : "ERROR: Interface Status Request : No response from remote host \"IP_MACHINE\"",
                        "val" : NumberLong(1)
                },
                "steps" : [
                        {
                                "_t" : "stateinc",
                                "t" : NumberLong(1524667892),
                                "a" : "NOM_CONN",
                                "m" : "ERROR: Interface Status Request : No response from remote host \"IP_MACHINE\"",
                                "val" : NumberLong(1)
                        },
                        {
                                "_t" : "statusinc",
                                "t" : NumberLong(1524667892),
                                "a" : "NOM_CONN",
                                "m" : "ERROR: Interface Status Request : No response from remote host \"IP_MACHINE\"",
                                "val" : NumberLong(1)
                        }
                ],
                "component" : "NOM_MACHINE",
                "connector" : "centreon",
                "connector_name" : "NOM_CONN",
                "creation_date" : NumberLong(1524667917),
                "display_name" : "XXXX",
                "extra" : {

                },
                "initial_output" : "",
                "last_update_date" : NumberLong(1524667917),
                "last_event_date" : NumberLong(1524667917),
                "resource" : "Traffic - eth0",
                "tags" : [ ]
        }
}
```

Pour « mesurer » la criticité courante de l'alarme, il faut vous référer à l'attribut `v.state`

```js
"state" : {
    "_t" : "stateinc",
        "t" : NumberLong(1524667892),
        "a" : "centreon.NOM_CONN",
        "m" : "ERROR: Interface Status Request : No response from remote host \"IP_MACHINE\"",
        "val" : NumberLong(1)
}
```

* `_t` : type d'évolution de la criticité (`stateinc`, `statedec`) : en croissance ou en décroissance
* `t` : timestamp associé
* `a` : auteur
* `m` : message de l'alarme
* `val` (0, 1, 2, 3) : criticité de l'alarme : Info, Mineure, Majeure, Critique.

Notez que l'attribut `step` conserve l'historique des changements sur l'alarme.

En parallèle, il existe également un attribut **status** sous forme de step qui permet de savoir si une alarme est en cours, en bagot, annulée, etc.  
Voici sa description :

```js
"status" : {
	"_t" : "statusinc",
	"t" : NumberLong(1529046262),
	"a" : "AUTEUR",
	"m" : "MESSAGE",
	"val" : NumberLong(1)
},
```

* `_t` : type de statut (`statusinc`, `statusdec`) : en croissance ou en décroissance
* `t` : timestamp associé
* `a` : auteur
* `m` : message de l'alarme
* `val` (0, 1, 2, 3, 4) : statut de l'alarme : Off, En cours, Furtif, Bagot, Annulé.

**Suppression**

Supprimer l'élément de la base MongoDB :
```js
db.periodical_alarm.remove({"v.component" : "NOM_MACHINE", "v.resource" : "Traffic - eth0"})
```

Supprimer une alarme qui a été annulée :
```js
db.periodical_alarm.remove({"v.component" : "XXXX", "v.resource" : "Ping", "v.resolved" : null, "v.status.val" : 4})
```

## État des composants de Canopsis

[Rendez-vous ici](etat-des-composants.md).
