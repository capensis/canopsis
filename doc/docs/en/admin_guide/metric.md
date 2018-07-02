# Metric

Le moteur metric est un moteur qui enregistre les données de performances
envoyés avec les événements (dans performance data
base influxdb.

La configuration du moteur metric est dans le fichier
`etc/metric/engine.conf` et a la structure suivante :

```
[ENGINE]
tags = parent_service
```

## Tags

L'option `tags` est une liste d'ids d'informations d'entités séparés par des
virgules. Ces ids vont être enregistrés avec les données de performances.
