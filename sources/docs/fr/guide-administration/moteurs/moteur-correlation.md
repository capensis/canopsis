# Moteur `engine-correlation` (Go, Cat)

Le moteur `engine-correlation` permet de créer des méta alarmes à partir de [`règles de gestion`](../metaalarm.md).

!!! note
    Ce moteur est disponible à partir de Canopsis 3.40.0.


## Utilisation


### Options du moteur

La commande `engine-correlation -help` liste toutes les options acceptées par le moteur.

```
  -consumeQueue string
    	Consume events from this queue. (default "Engine_correlation")
  -d	debug
  -printEventOnError
    	Print event on processing error
  -publishQueue string
    	Publish event to this queue. (default "Engine_watcher")
  -version
    	version infos
```

### Multi-instanciation

!!! note
    Cette fonctionnalité sera disponible à partir de Canopsis 3.41.0. Elle ne doit pas être utilisée sur les versions antérieures.

Il est possible, à partir de **Canopsis 3.41.0**, de lancer plusieurs instances du moteur `engine-correlation`, afin d'améliorer sa performance de traitement et sa résilience.

En environnement Docker, il vous suffit par exemple de lancer Docker Compose avec `docker-compose up -d --scale correlation=2` pour que le moteur `engine-correlation` soit lancé avec 2 instances.

Cette fonctionnalité sera aussi disponible en installation par paquets lors d'une prochaine mise à jour.

## Fonctionnement

Tous les événements qui circulent dans Canopsis sont transmis à la file `Engine_correlation` et seront donc lus par le moteur `engine-correlation`.

Le moteur va alors vérifier si une règle de corrélation doit s'appliquer et si c'est le cas, il générera ou modifiera une meta alarme.

La documentation sur les règles de gestion est disponible [`ici`](../metaalarm.md)

## Collection MongoDB associée

Les entités sont stockées dans la collection MongoDB `meta_alarm_rules`.


```json
{
	"_id" : "73da1ad7-058e-46af-8442-7ea3f246eb68",
	"patterns" : null,
	"config" : null,
	"name" : "Relation-composant-ressource",
	"type" : "relation"
}
```
