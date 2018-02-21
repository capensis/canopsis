# Lifeline

Les baseline permettent de s'assurer qu'une entité recoit toujours des événements.
Si aucun évènement n'apparait après un certain délai, une alarme est levée.

Il y a besoin de faire évoluer *alerts* pour qu'il puisse prendre en compte un nouvel état: **Dead**.

## Entrée des données (évènement)

Un évènement arrive dans l'engine cleaner. Une fois validé, il est déposé sur une nouvelle queue rabbit (en plus du traitement habituel par les engines).

La queue *baseline* est consommée par un nouvel engine *Lifeline* (en Golang), qui fait prend l'évènement et le publie l'entity_id et le timestamp dans un redis.

## Traitement des données (beat process)

Dans le nouvel engine *Lifeline* (en Golang), au beat processing :
* On lit la configuration des lifelines
* Pour chaque entitée surveillée, on regarde la dernier timestamp présent dans redis. Si le différenciel entre timestamp et maintenant dépasse la limite, on lève une alarme avec notre nouveau state

Pour les alarmes en state *Dead*, il suffit qu'un évènement arrive pour que le state se mette à jour.

## Frontend

À terme, le client doit pouvoir modifier le paramétrage des lifelines.

### Créer/modifier/supprimer une lifeline

Exemple de design de création, avec l'explorateur d'entitées :
* Sélectionner les entitées ciblées
* Cliquer sur un nouveau bouton « Créer une ligne de vie »
* Dans une popup, rentrer la durée du déclenchement de la ligne de vie

## Format de configuration

Une lifeline est un document dans la collection *default_lifeline*.

Exemple :

```json
{
    'entity_filter': {'$in': ['entity_id']},
    'countdown': 60
}
```

### Description

* **entity_filter** (dict) : filtre d'entitées mongo
* **countdown** (int) : durée en secondes avant de paniquer

## Intégration

Il y a besoin d'un redis HA.
