# Moteur `engine-webhook` (Pro)

!!! info
    Disponible uniquement en édition Pro.

Le moteur `engine-webhook` permet d'automatiser la gestion de la vie des tickets vers un service externe en fonction de l'état des évènements ou des alarmes.

Des exemples pratiques d'utilisation des webhooks sont disponibles dans la partie [Exemples](#exemples).

## Utilisation

### Options du moteur

La commande `engine-webhook -help` liste toutes les options acceptées par le moteur.

## Fonctionnement

À l'arrivée dans sa file, le moteur va vérifier si l'événement correspond à un ou plusieurs de ces Webhooks.

Si oui, il va alors immédiatement appliquer les Webhooks correspondant (il n'y a pas de *beat*).

En cas d'échec, il existe un mécanisme de réémission du webhook.

La documentation d'utilisation se trouve dans le [guide d'utilisation](../../../guide-utilisation/menu-exploitation/scenarios/)
