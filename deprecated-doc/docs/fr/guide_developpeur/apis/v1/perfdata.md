# Récupération des métriques du contexte

## URL
`POST` `/api/context/metric`

## Paramêtres
La route attend jusqu'à trois paramètres encapsulés dans un fichier JSON
dans le corps de la requête :

  * **limit** pour ne retourner que les *n* première métriques. Si le paramètre
  n'est pas transmis, la valeur 100 est prise par défaut.
  * **start** pour ne retourner que les métriques à partir de la *n*-ième. 
  Si le paramètre n'est pas transmis, la valeur 0 est prise par défaut.
  * **_filter** un dictionnaire avec un champ **name** qui contient lui même
  un champ **$regex** de type chaine de caractères. Ce champ va permettre
  de filtrer les métriques qui ont la valeur du champ *$regex* dans leur *eid*.
  Si aucun filtre n'est donné, toutes entités seront retournées.

Example :

```json
{
  "limit":5,
  "start":0,
  "_filter":
    {
    "name":
      {
        "$regex":"cps"
      }
    }
}
```

## Résultat
La route retourne un résultat similaire à l'exemple ci-dessous.

```json
{
  "total": 87,
  "data": [
    {
      "connector": "Engine",
      "connector_name": "engine",
      "resource": "task_linklist",
      "name": "cps_sec_per_evt",
      "component": "pcv-arthur",
      "_id": "/metric/Engine/engine/pcv-arthur/task_linklist/cps_sec_per_evt",
      "type": "metric",
      "internal": false
    }
  ],
  "success": true
}
```

  * Le champ **total** indique le nombre de métriques qui correspond au
  filtre appliqué lors de la recherche. Attention, il ne correspond
  pas au nombre de métriques contenues dans le champ *data* ;
  * Le champ **data** contient une liste de métriques ;
  * Le champ **success** contient un booléen. Si la requête s'est
  exécutée correctement, elle sera à *true*, sinon *false*.


En cas d'erreur lors du traitement de la requête :

```json
{
  "total": 0,
  "data": {
    "msg": "HTTPConnectionPool(host='localhost', port=8086): Max retries exceeded with url: /query?q=SHOW+SERIES%3B&db=canopsis (Caused by <class 'socket.error'>: [Errno 111] Connection refused)",
    "traceback": "TODO",
    "type": "<class 'requests.exceptions.ConnectionError'>"
  },
  "success": false
}
```
  * Le champ **success** vaut *false*
  * Le champ **data** contient 3 champs :
    * un champ **msg** qui contient une description de l'erreur
	* un champ **traceback** qui contient la traceback liée à l'erreur
	* un champ **type** qui contient le type d'exception survenue.
