# Task ackcentreon

!!! attention
    Ce moteur n'est disponible que dans l'édition CAT de Canopsis.

Le moteur `task_ackcentreon` permet de *descendre* les ACK positionnés depuis Canopsis vers l'outil Centreon. Ceci est valable aussi bien pour les poses d'ACK que pour les suppressions d'ACK.

Ainsi, lorsqu'un ACK est posé sur Canopsis, l'information est *répliquée* sur le Poller Centreon qui avait généré l'alarme. En utilisation conjointe du [connecteur Centreon](../../interconnexions/Supervision/Centreon.md), la communication est bi-directionnelle.

## Fonctionnement

Voici les différentes étapes permettant d'obtenir le résultat souhaité :

*  Un ACK est posé ou retiré depuis Canopsis
*  Le moteur [`event_filter`] (**Python**) est configuré pour exécuter un job
*  Le job `ackcentreon` est exécuté et suit les étapes suivantes :
    *  Connexion SSH vers le serveur Centreon Central
    *  Demande d'informations via CLAPI (ID du poller qui a généré l'alarme à l'origine)
    *  Génération d'une commande externe conforme à la pose ou la suppression de l'ACK
    *  POST de la commande dans le fichier de commande externe

## Mise en place

### Activation du moteur

Sur le nœud des moteurs Canopsis :

```sh
systemctl enable canopsis-engine-cat@task_ackcentreon-task_ackcentreon.service
systemctl start canopsis-engine-cat@task_ackcentreon-task_ackcentreon.service
```

### Accès SSH entre Canopsis et Centreon

La remontée d'informations de Canopsis vers Centreon s'exécute via SSH. Il est donc nécessaire de transférer une clé publique SSH depuis le nœud Canopsis (avec l'utilisateur `canopsis`) vers le nœud central Centreon (avec l'utilisateur `centreon`).

Sur le nœud Canopsis, on crée une nouvelle clé RSA si nécessaire, et on récupère le contenu de la clé publique associée :

```sh
su - canopsis
[ ! -r ~/.ssh/id_rsa.pub ] && ssh-keygen -t rsa -N ''
cat ~/.ssh/id_rsa.pub
```

Le contenu du fichier `/opt/canopsis/.ssh/id_rsa.pub` du nœud Canopsis doit ensuite être ajouté au fichier `/var/spool/centreon/.ssh/authorized_keys` du nœud Centreon.

### Mise en place de CLAPI sur le nœud Centreon

CLAPI doit être disponible sur l'hôte Centreon.  Vous devez donc vous en assurer. Une authentification sera demandée par le moteur `task_ackcentreon`.

!!! note
    À titre d'information, la commande finale qui sera utilisée est la suivante :

    `/chemin/vers/centreon -u un_utilisateur -p un_mdp -a POLLERLIST`

### Ajout d'un job « ACK Centreon » dans Canopsis

Danos Canopsis, vous devez créer un job avec les paramètres `type: notification` et `xtype: ack_centreon` dans Canopsis.

Renseignez ensuite les informations demandées :

*  Hôte Centreon
*  Utilisateur/port SSH
*  Chemin CLAPI
*  Authentification CLAPI

### Ajout d'une règle dans l'event\_filter Python

Le moteur `event_filter` (Python) doit exécuter le job `ack_centreon` lorsqu'il reçoit un évènement de type `ack` ou `ackremove` (pause et suppression d'un ACK).

!!! note
    Le moteur `event_filter` Python est requis pour `task_ackcentreon`, même dans un environnement Go.

Il faut, pour cela, créer une règle d'`event_filter` Python, avec le filtre suivant :

```json
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
```

!!! note
    Notez bien le `"extra.origin": "canopsis"` qui précise que seuls les ACK/ACKREMOVE en provenance de Canopsis doivent être transmis à Centreon. Ceci est indispensable pour ne pas créer de boucles entre Canopsis et Centreon.

!!! attention
    Il est possible, dans un environnement pur Python, que le champ `"extra.origin": "canopsis"` doive être remplacé par `"origin": "canopsis"`. Tester l'une ou l'autre de ces valeurs, si la tâche semble ne pas se déclencher.

Après avoir configuré le filtre, créer l'action suivante :

*  `execjob`, qui pointe sur le job précédemment créé.

## Procédure supplémentaire en environnement Go

!!! attention
    Cette procédure est obligatoire dans le cas d'un environnement Go. Elle n'a pas besoin d'être exécutée dans un environnement pur Python.

Lorsque vous utilisez un environnement Canopsis avec moteurs en Go, vous devez spécifier une règle d'enrichissement supplémentaire, pour le moteur `che`.

Cette règle permet d'ajouter aux évènements qui circulent les informations de l'entité correspondante disponible dans Canopsis. Contrairement aux moteurs Python où cette information est ajoutée par défaut, les moteurs Go nécessitent une règle explicite. Sans cette règle, la tâche `ackcentreon` ne peut pas fonctionner en environnement Go.

Sur un nœud frontend Canopsis, créer un fichier `enrichentity.json` :

```json
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
```

Cette règle est à envoyer à l'API `eventfilter` **Go** de cette manière (adapter les identifiants Canopsis, l'URL et le nom du fichier JSON si nécessaire) :

```sh
curl -X POST -u root:root -H "Content-Type: application/json" -d @enrichentity.json 'http://localhost:8082/api/v2/eventfilter/rules'
```

## Procédure de test

À ce stade, vous pouvez poser un ACK dans Canopsis et vérifier sur l'interface de Centreon qu'il a bien été transmis, et même chose pour la suppression d'un ACK.
