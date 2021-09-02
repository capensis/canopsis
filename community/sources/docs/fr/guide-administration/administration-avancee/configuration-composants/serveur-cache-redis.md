# Configuration avancée du serveur de cache Redis intégré à Canopsis

Le serveur de cache [Redis](https://redis.io) permet un accès rapide aux données les plus utilisées dans la modélisation interne que Canopsis construit afin de représenter votre périmètre de surveillance. Il est essentiel pour la performance de Canopsis.

## Optimisations système

En dehors d'une utilisation pour un périmètre réduit, il est recommandé que Redis soit installé sur une instance dédiée.

Si tel est le cas, veuillez appliquer les [optimisations système officiellement recommandées par Redis](https://redis.io/topics/admin).

## Adaptation du paramètre `maxmemory` par rapport à votre besoin

En dehors d'une utilisation pour un périmètre réduit, il est crucial que vous adaptiez la valeur `maxmemory` de Redis en fonction de l'étendue et de la complexité de la supervision que vous souhaitez intégrer à Canopsis.

Assurez-vous tout d'abord que l'instance hébergeant Redis dispose bien des ressources nécessaires en RAM.

Vous devez ensuite ajuster la valeur suivante du fichier `/etc/redis.conf` :

```python
maxmemory "512mb" 
```

La valeur par défaut est de 512 Mio. Vous devez ajuster cette valeur de façon raisonnable, en fonction de votre environnement. Une augmentation de plusieurs Gio peut être nécessaire, sur les plus larges périmètres.

Vous devez ensuite [redémarrer Redis](../../gestion-composants/arret-relance-composants.md#redis).

!!! attention
    Afin de conserver la cohérence des données en cache, et afin de détecter immédiatement le problème, Redis est configuré avec une politique `noeviction` : les moteurs seront incapables d'ajouter de nouvelles données en cache lorsque la limite `maxmemory` est atteinte, ce qui bloquera les traitements.

    Vous devez impérativement laisser une marge suffisante à votre valeur `maxmemory` afin d'éviter cela.
