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
