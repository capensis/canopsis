# Task ackcentreon

Le moteur `task_ackcentreon` permet de **descendre** les acks positionnés depuis Canopsis vers l'outil Centreon.  
Cela est valable aussi bien pour les poses d'ACK que pour les suppressions d'ACK.  

Ainsi, lorsqu'un ACK est posé sur Canopsis, l'information est **répliquée** sur le Poller Centreon qui avai généré l'alarme.  
En utilisation conjointe du [connecteur Centreon](../../guide-connecteurs/Supervision/Centreon.md), la communication est bi-directionnelle.  


## Workflow

Voici les différentes étapes permettant d'obtenir le résulat souhaité :

*  Un ACK est posé ou retiré depuis Canopsis
*  Le moteur [`event_filter`](moteur-event_filter.md) est configuré pour exécuter un job
*  Le job `ackcentreon` est exécuté et suit les étapes suivantes :
    *  Connexion SSH vers le serveur Centreon Central
    *  Demande d'informations via CLAPI (ID du poller qui a généré l'alarme à l'origine)
    *  Génération d'une commande externe conforme à la pose ou la suppression de l'ACK
    *  POST de la commande dans le fichier de commande externe


### Infrastructure

**Activation du moteur**

````
systemctl enable canopsis-engine-cat@task_ackcentreon-task_ackcentreon.service
systemctl start canopsis-engine-cat@task_ackcentreon-task_ackcentreon.service
````

**SSH**

Etant donné que les commandes vont circuler over SSH, un échange de clés est nécessaire.  
C'est l'utilisateur `canopsis` qui sera à l'origine des exécutions over ssh.  
Sa clé publique doit donc être placée dans .authorized_keys de l'utilisateur ̀`centreon` sur la machine centrale Centreon.  

**CLAPI**

`CLAPI` doit être disponible sur l'hôte Centreon.  
Vous devez donc vous en assurer. Une authentification sera demandée par le moteur `task_ackcentreon`.  

La commande finale qui sera utilisée est la suivante :

````
/path/to/centreon -u un_utilisateur -p un_mdp -a POLLERLIST"
````

### Job Ack Centreon

Vous devez créer un job **notification** de type **ack_centreon** dans Canopsis.  

* xtype : ackcentreon, type : notification

Renseignez les informations demandées :

* Hôte Centreon
* Utilisateur/port SSH
* Path Clapi
* Authentification Clapi

### Règle Event Filter

Le moteur `event_filter` doit exécuter le job ack_centreon lorsqu'il reçoit un événement de type **ack** ou **ackremove**.  
Il vous faut créer une règle de event_filter avec comme paramètres :

Filtre :

````json
{
  "$and": [
    {
      "connector": "centreon"
    },
    {
      "connector_name": "Central"
    },
    {
      "$or": [
        {
          "event_type": "ack"
        },
        {
          "event_type": "ackremove"
        }
      ]
    },
    {
      "extra.origin": "canopsis"
    }
  ]
}
````

Notez bien le `"extra.origin": "canopsis"` qui précise que seuls les ACK/ACKREMOVE en provenance de Canopsis doivent être transmis à Centreon.  
Cela permet d'éviter des boucles entre Canopsis et Centreon.

Action :

`execjob` qui pointe sur le job précédemment créé.

## Cas d'un environnement avec moteurs GO

Lorsque vous utilisez un environnement Canopsis avec moteurs en GO, vous devez spécifier en plus des éléments précédent une règle d'enrichissement pour le moteur **che**.  
Cette règle permet d'ajouter aux événements qui circulent les informations de l'entité correspondante disponible dans Canopsis.  
Le moteur **ack_centreon** en a besoin pour savoir s'il y a bien une alarme en cours et pour laquelle on souhaite descendre un ACK vers Centreon.  

````
$ cat ../enrich/enrichentity.json 
{
    "type": "enrichment",
    "pattern": {},
    "external_data": {
        "entity": {
            "type": "entity"
        }
    },
    "actions": [
        {
            "type": "copy",
            "from": "ExternalData.entity",
            "to": "Entity"
        }
    ],
    "on_success": "pass",
    "on_failure": "pass",
    "priority": 1
}
````

Cette règle est à poster sur l'API `eventfilter` de cette manière : 

````
curl -X POST -u root:root -H "Content-Type: application/json" -d @enrichentity.json 'http://localhost:28082/api/v2/eventfilter/rules'
````

## Tests de bout en bout

A ce stade, vous pouvez poser un ACK dans Canopsis et vérifier sur l'interface de Centreon qui l'a bien été transmis et même chose pour la suppression d'un ACK.
