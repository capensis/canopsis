# Statsng

Le moteur statsng est un moteur qui reçoit des événements statistiques envoyés
par les autres moteurs et les utilise pour calculer des statistiques.

La configuration du moteur statsng est dans le fichier
`etc/statsng/engine.conf` et a la structure suivante :

```
[ENGINE]
send_events = True
entity_tags = parent_service
pbehavior_tags = Maintenance,Pause
```

## Envoi d'événements

Les événements statistiques ne sont envoyés que si la valeur de l'option
`send_events` est `True`. Le moteur statsng ne va pas calculer de statistiques
si `send_event` vaut `False`.

## Entity Tags

L'option `entity_tags` est une liste d'ids d'informations d'entités séparés par
des virgules. Ces ids vont être enregistrés avec les statistiques, permettant
de les utiliser dans les filtres de l'[API stats](../guide_developpeur/apis/v2/stats.md).
