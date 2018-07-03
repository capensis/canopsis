# Metric

Le moteur metric est un moteur qui enregistre les données de performances
envoyés avec les événements (dans `perf_data` et `perf_data_array`) dans la
base influxdb.

## Modèle de données

Les données de performances sont sauvegardées dans un measurement influxdb ayant
le même nom que la métrique. Ce measurement peut contenir trois champs :

 - `value` : la valeur de la métrique
 - `warn` : le seuil d'avertissement (peut valoir `null`)
 - `crit` : le seuil critique (peut valoir `null`)

Il contient également les tags `connector`, `connector_name`, `component` et
`resource`. Il est possible d'ajouter des informations sur l'entité avec
l'option `tags` de la configuration.

## Configuration

La configuration du moteur metric est dans le fichier
`etc/metric/engine.conf` et a la structure suivante :

```
[ENGINE]
tags = parent_service
```

### Tags

L'option `tags` est une liste d'ids d'informations d'entités séparés par des
virgules. Ces ids vont être enregistrés avec les données de performances.

Chaque `<information_id>` présent dans la liste est ajouté comme tag aux
measurements créés par le moteur `metric`. La valeur de ce tag est la valeur
stockée dans `<entity>.infos.<information_id>`.
