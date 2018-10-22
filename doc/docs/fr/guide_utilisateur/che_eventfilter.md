# Filtre à événement

Cette page décrit le filtre à événement implémenté dans le moteur go `che`. Le
moteur python `event-filter` est documenté
[ici](../../en/user_guide/event_filter.md).


## Règles

Une règle est un document JSON contenant les paramètres suivants :

 - `type` (requis) : le type de règle (`enrichment`, `drop`, `break`).
 - `pattern` (optionnel) : les événements correspondant à ce pattern seront
   traités par la règle. Si le pattern n'est pas précisé, la règle est
   appliquée à tous les événements.
 - `priority` (optionnel, 0 par défaut) : la priorité de la règle (les règles
   sont appliquées par ordre croissant de priorité).
 - `enabled` (optionnel, `true` par défaut) : `false` pour désactiver la règle.

Par exemple, la règle suivante supprime les événements dont la ressource vaut
`invalid_resource` :

```json
{
    "type": "drop",
    "pattern": {
        "resource": "invalid_resource"
    },
    "priority": 10
}
```


### Application des règles

Lors de la réception d'un événement, les règles sont parcourues par ordre de
priorité croissante. Si l'événement est reconnu par le pattern d'une règle
active, elle est appliquée.

 - Si le résultat de la règle est `pass`, l'événement passe à la prochaine
   règle.
 - Si le résultat de la règle est `break`, on n'applique plus de règles à
   l'événement.
 - Si le résultat de la règle est `drop`, l'événement est supprimé. Les
   événements supprimés sont loggés.

Le résultat d'une règle de type `break` est toujours `break`, et celui d'une
règle de type `drop` est toujours `drop`. Le résultat des règles de type
`enrichment` est spécifié dans leur documentation.

Si l'événement est invalide à la fin de l'exécution des règles, il est
supprimé. Ces événements sont loggés. Un événement est valide si :

 - son champ `source_type` vaut `component`, et son champ `component` est
   défini; *ou*
 - son champ `source_type` vaut `resource`, et ses champs `component` et
   `resource` sont définis.

Si le champ `debug` d'un événement vaut `true`, le passage de l'événement dans
les règles est loggé. Ce champ peut être défini avec une règle d'enrichissement.

### Patterns

Un pattern permet de sélectionner les événements auxquels une règle s'applique.
Un pattern est défini par un dictionnaire contenant les valeurs de certaines
clés d'un événement. Par exemple :

```json
"pattern": {
    "component": "component_name",
    "resource": "resource_name"
}
```

Une règle contenant ce pattern sera appliquée aux événements dont le composant
vaut `component_name` et la ressource vaut `resource_name`.

Pour plus d'expressivité, il est possible d'associer à une clé un dictionnaire
contenant des couples `operateur: valeur`. Les opérateurs disponibles sont :

 - `>=`, `>`, `<`, `<=` : compare une valeur numérique à une autre valeur.
 - `regex_match` : filtre la valeur d'une clé selon une expression régulière.

Par exemple, le pattern suivant sélectionne les événements dont l'état est
compris entre 1 et 3 et dont l'output vérifie une expression régulière :

```json
"pattern": {
    "state": {">=": 1, "<=": 3},
    "output": {"regex_match": "Warning: CPU Load is critical \(.*\)"}
}
```

Si la règle est appliquée après l'enrichissement avec l'entité, l'entité
correspondant à l'événement est disponible dans le champ `current_entity`. On
peut alors ajouter un dictionnaire dans `current_entity` pour filtrer sur les
champs de l'entité. Par exemple, le pattern suivant sélectionne les
événements dont l'entité est active et n'a pas d'information
`service_description` définie.

```json
"pattern": {
    "current_entity": {
        "enabled": true,
        "infos": {
            "service_description": null
        }
    }
}
```


## Exemples de règles

### Suppression d'événements

La règle suivante supprime les événements dont la ressource vaut
`invalid_resource` :

```json
{
    "type": "drop",
    "pattern": {
        "resource": "invalid_resource"
    },
    "priority": 10
}
```

La règle suivante supprime les événements majeurs et critiques sur les
resources dont le nom commence par "cpu-" :

```json
{
    "type": "drop",
    "pattern": {
        "state": {">=": 2},
        "resource": {"regex_match": "cpu-.*"}
    },
    "priority": 10
}
```

### Sortie de l'eventfilter

La règle suivante fait sortir tous les événements de type `pbehavior` de
l'eventfilter :

```json
{
    "type": "break",
    "pattern": {
        "event_type": "pbehavior"
    },
    "priority": 0
}
```
