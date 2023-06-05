# Moteur `engine-axe` (Community)

Le moteur `engine-axe` permet de créer et d'enrichir les alarmes. Il permet également d'appliquer les actions entrées depuis le Bac à alarmes.

## Utilisation

### Options du moteur

La commande `engine-axe -help` liste toutes les options acceptées par le moteur.

### Multi-instanciation

Il est possible, à partir de Canopsis 3.40.0, de lancer plusieurs instances du moteur `engine-axe`, afin d'améliorer sa performance de traitement et sa résilience.

Modifier le ficher `docker-compose.yml` (ou en fonction de votre version`02-app.docker-compose.yml`) de la manière suivante :

```yml
  axe:
    ...
    deploy:
      replicas: 2
```

Exécutez ensuite `docker-compose up -d` (ou `docker-compose -f 02-app.docker-compose.yml`) pour appliquer cette modification.


Cette fonctionnalité sera aussi disponible en installation par paquets lors d'une prochaine mise à jour.

### Option `EnableLastEventDate`

!!! attention
    Activer cette option entraîne une action supplémentaire systématique dans le moteur qui a une incidence négative sur ses performances.

Les alarmes dans Canopsis incluent un champ `alarm.v.last_event_date`.

Cependant, la mise à jour de ce champ n'est pas activée par défaut. Sa valeur est celle de `alarm.v.creation_date`, soit la date de création de l'alarme par `engine-axe`.

Pour l'activer, passez le paramètre `EnableLastEventDate` du fichier de configuration à `true`.

### Option `StealthyInterval`

Si une alarme change de criticité (de stable vers alerte ou inversement), une ou plusieurs fois, durant ce délai, elle sera alors considérée comme `furtive` et se verra attribuer le [statut](../../guide-utilisation/vocabulaire/index.md#statut) correspondant.

### Option `CancelAutosolveDelay`

Lorsqu'une alarme est annulée manuellement, via l'interface web par exemple, elle prend le statut annulée et reste pendant 1h dans le bac des alarmes en cours. Passé le délai d'une heure, elle change de statut pour passer en résolue et bascule dans le bac des alarmes résolues tout en gardant le dernier niveau de criticité connu.

Vous pouvez agir sur ce délai en modifiant le paramètre `CancelAutosolveDelay`.

### Option `DisplayNameScheme`

Vous avez la possibilité de personnaliser le schéma de construction de l'attribut `display_name` d'une alarme par l'intermédiaire de l'option `DisplayNameScheme`.

L'attribut `display_name` d'une alarme permet d'identifier une alarme par une chaîne plus simple que son identifiant technique.

!!! attention
    Canopsis n'apporte pas la garantie que cet identifiant sera unique.
    Il vous appartient d'utiliser un schéma qui offre une probabilité suffisamment faible par rapport au nombre d'alarmes que vous allez traiter.

Par défaut, le schéma utilisé est le suivant : "{{ rand_string 2 }}-{{ rand_string 2 }}-{{ rand_string 2 }}"

Vous pouvez modifier cette valeur en utilisant une fonction du tableau ci-après (Une seule fonction à ce jour)

| Fonction | Description | Syntaxe
| ------ | ------ | ---- |
| `rand_string` | Lettre ou chiffre aléatoire | `rand_string ${longueur}`

Exemples :

```ini
[alarm]
...
DisplayNameScheme = "{{ rand_string 3 }}-{{ rand_string 3 }}-{{ rand_string 3 }}"
```

```ini
[alarm]
...
DisplayNameScheme = "{{ rand_string 4 }}_{{ rand_string 3 }}_{{ rand_string 2 }}"
```

### Option `DisableActionSnoozeDelayOnPbh`

Si cette option est mise à la valeur `true`, cela rend l'action `unsnooze indépendante d'un Pbehavior` et elle s’exécutera même pendant la période de maintenance si ce Pbehavior est toujours actif au moment où l'action doit normalement s’exécuter.

Si elle est mise à `false` (valeur par défaut), l'action `unsnooze sera dépendante` et attendra donc la fin du Pbehavior pour commencer à s'appliquer.

## Fonctionnement du moteur

La file du moteur est placée juste après le moteur [`engine-che`](moteur-che.md).

À l'arrivée dans sa file, le moteur `engine-axe` va transformer les événements en alarmes qu'il va créer et mettre à jour.

Lorsque la multi-instanciation est activée, une seule instance d'`engine-axe` s'occupe du *periodical process*. Ce mécanisme est automatique.

### Gestion des événements de type check

3 possibilités pour un événement de type `check` :

1. Il ne correspond à aucune alarme en cours : l'alarme va alors être créée
2. Il correspond à une alarme en cours et son champ `state` ne vaut pas `0` : l'alarme va alors être mise à jour
3. Il correspond à une alarme en cours et son champ `state` vaut `0` : l'alarme va alors passer en `OK`. Au [battement (beat)](../../guide-utilisation/vocabulaire/index.md#battement) suivant, si l'alarme n'a pas été rouverte par un nouvel événement de type `check`, elle est considérée comme résolue. Un champ `v.resolved` lui est alors ajouté avec le timestamp courant.

### Gestion des autres types d'événements

Si l'événement correspond à une action (comme la pose d'un acquittement), l'alarme va être mise à jour en appliquant l'action.

## Collection MongoDB associée

Les alarmes sont stockées dans la collection MongoDB `periodical_alarm`.

Le champ `_id` est généré automatiquement.

Le champ `d` correspond à l'`_id` de l'entité à laquelle l'alarme est rattachée.

```json
{
    "_id" : "aad73d0b-2e0e-453d-90c5-1c843cd196b2",
    "t" : 1567498879,
    "d" : "disk2/serveur_de_salle_machine_DHCP",
    "v" : {
        "state" : {
            "_t" : "stateinc",
            "t" : 1567498879,
            "a" : "superviseur1.superviseur1",
            "m" : "Disque plein a 98%, 50GO occupe",
            "val" : 2
        },
        "status" : {
            "_t" : "statusinc",
            "t" : 1567498879,
            "a" : "superviseur1.superviseur1",
            "m" : "Disque plein a 98%, 50GO occupe",
            "val" : 1
        },
        "steps" : [
            {
                "_t" : "stateinc",
                "t" : 1567498879,
                "a" : "superviseur1.superviseur1",
                "m" : "Disque plein a 98%, 50GO occupe",
                "val" : 2
            },
            {
                "_t" : "statusinc",
                "t" : 1567498879,
                "a" : "superviseur1.superviseur1",
                "m" : "Disque plein a 98%, 50GO occupe",
                "val" : 1
            }
        ],
        "component" : "serveur_de_salle_machine_DHCP",
        "connector" : "superviseur1",
        "connector_name" : "superviseur1",
        "creation_date" : 1567498879,
        "display_name" : "XA-KU-AQ",
        "extra" : {},
        "initial_output" : "Disque plein a 98%, 50GO occupe",
        "output" : "Disque plein a 98%, 50GO occupe",
        "initial_long_output" : "",
        "long_output" : "",
        "long_output_history" : [
            ""
        ],
        "last_update_date" : 1567498879,
        "last_event_date" : 1567498879,
        "resource" : "disk2",
        "state_changes_since_status_update" : 0,
        "tags" : [],
        "total_state_changes" : 1
    }
}
```
