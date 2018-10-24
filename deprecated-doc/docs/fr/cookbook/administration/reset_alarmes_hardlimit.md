# Procédure de remise à zéro des alarmes en Hard Limit

** Attention: cette procédure provoque une destruction de données, effectuer impérativement une sauvegarde de la BDD avant son exécution**


Cette procédure vise à réinitialiser les alarmes ayant atteint leur hard limit afin de les rendre modifiables. 

La procédure s'effectue en 2 etapes: 

- sauvegarde de la base de données
- réinitiaisation des alarmes

La réinitiaisation des alarmes provoque la suppression de l'ensemble de l'historique de toutes les alarmes en hard limit


## Sauvegarde de la BDD 


- Se connecter à un serveur Master MongoDB
- lancer la sauvegarde : 

```bash
$ mongodump -u cpsmongo -p canopsis -d canopsis --gzip --archive=dump.gz
```

stocker le fichier `dump.gz` qui sera utile en cas de restauration

## Réinitialiser les alarmes en hardlimit

- Se connecter à un serveur Master MongoDB
- ouvrir un shell mongo : 

```bash
$ mongo -u cpsmongo -p canopsis canopsis
```

Vérifier que certaines alarmes sont bien éligibles à la réinitialisation, en effectuant la requête suivante :

```javascript

db.getCollection('periodical_alarm').find({"v.hard_limit": {$ne: null}})

```
Si cette requête ne retourne aucun résultat, inutile de continuer: toutes les alarmes sont modifiables. 

Sinon,exécuter la requête suivante:

```javascript
db.getCollection('periodical_alarm').updateMany({"v.hard_limit": {$ne: null}}, {$set: {"v.hard_limit": null,"v.steps": []}})
```

Le résultat attendu indique le nombre d'alarmes mises à jour

```json
{
    "acknowledged" : true,
    "matchedCount" : 6.0,
    "modifiedCount" : 6.0
}
```

Les attributs `matchedCount` et `modifiedCount` indiquent respectivement le nombre d'alarmes correspondant au filtre et le nombre d'alarmes mises à jour


Exécuter à nouveau la requête de contrôle pour s'assurer que toutes les alarmes ont bien été réinitialisées:

```javascript
db.getCollection('periodical_alarm').find({"v.hard_limit": {$ne: null}})

```

La requête ne doit retourner aucun résultat. Si elle en retourne, recommencer l'étape puis contacter le support.


Les alarmes sont alors complètement vidé de leur historique, mais les actions d'acquittement et d'annulation pourront être appliquées. 


## En cas de problème: restauration de la BDD

Pour restaurer la base de données, effectuer les étapes suivantes: 

- se connecter au master mongoDB 

depuis le dossier où le fichier dump.gz a été conservé, effectuer la commande: 

```bash
$  mongorestore -u cpsmongo -p canopsis -d canopsis --drop --gzip --archive=dump.gz
```
