# API healthcheck

L'API healthcheck permet d'obtenir l'état de fonctionnement de Canopsis ; par exemple elle permet de savoir si les services nécessaires sont disponibles.

## Récupérer l'état

```
GET /api/v2/healthcheck/?criticals=amqp,cache,database,engines,time_series
```

Renvoie un résumé de l'état.

```json
{
    "amqp": "",
    "cache": "",
    "database": "",
    "engines": "",
    "time_series": "",
    "timestamp": 1541774674,
    "overall": true
}
```

### Description du résumé

La réponse contient les champs suivants :
* amqp, cache, database, engines, timeseries : un message d'erreur associé à chacque service ; rien s'il n'y a pas d'erreur détectée ;
* timestamp : le moment de création du résultat ;
* overall : un booléen pour savoir si l'état global est bon ou mauvais.

Concernant `overall`, par défaut, tous les services sont pris en compte pour calculer l'état global. Il est toutefois possible de sélectionner les services à considérer comme indispensable en utilisant le paramètre `criticals` dans l'url GET.
`criticals` est une liste de services séparés par des virgules.

### Ce qui est vérifié

* amqp : la connection est ouverte, le channel aussi, et il est possible de publier un message. On vérifie aussi qu'une liste de queues existe bien, qu'il y a au moins un Consumer dessus, et que la queue est active et non saturée (> 100 000 messages en attente) ;
* cache: la connection est fonctionnelle et il est possible de faire un ECHO ;
* database: la connection fonctionne et il est possible de lire dans une liste de collections ;
* engines : hors docker, vérifie par systemctl que les engines de bases (python) sont "running" ;
* time series : vérifie que la database existe et que l'on peut lire des measurements.
