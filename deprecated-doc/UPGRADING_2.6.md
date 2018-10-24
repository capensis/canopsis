# Mise à jour Canopsis 2.6

## Migration des vues

Canopsis 2.6 permet de renommer les colonnes du bac à alarmes. C'est complètement transparent, mais nécessite une migration en base de données
pour les vues existantes. Celle-ci doit être faite manuellement car la structure actuelle d'une vue ne permet pas de migrer automatiquement ces données sans risque.

**Avant de commencer** : faire un dump de la base mongo pour éviter toute perte de données en cas de problème de migration

La migration sera faite avant la mise à jour de Canopsis : cela permet de maintenir un canopsis opérationnel au long de la mise à jour:

- ouvrir canopsis de préprod.
- ouvrir chacun des liens vers les bacs à alarmes pour récupérer leurs IDs dans l'URL (dernière partie de l'URL)
- noter les IDs et les noms des vues
- ouvrir un shell MongoDB
- Exécuter la commande pour chaque vue :

```javascript
db.object.find({'_id': 'id trouvé dans lurl'})
```

- Modifier le document obtenu, en ajoutant la structure suivante au même niveau que le tableau `columns`, juste après ce dernier :

```javascript
"widget_columns": [
    {
        "value": "resource",
        "label": "resourcessss"
    },
    {
        "value": "component",
        "label": "componentssss"
    },
    {
        "value": "extra_details",
        "label": "suivi des actions"
    }
],
```

Le contenu de chaque objet est le suivant :

- `value` : le nom technique de la colonne
- `label` : le nom affiché dans le bac à alarmes

Une fois l'objet inséré, sauvegarder le document modifié.

*Répéter l'opération pour chaque vue*

recharger l'ui via F5 : le bac à alarme doit rester fonctionnel

## Mise à jour de Canopsis

installer le rpm `canopsis-cat-2.5.16-1.e17.centos.x86_64.rpm`


Mettre à jour les dépendances systemd et redémarrer Canopsis (à exécuter en root)

```shell
systemctl daemon-reload
/opt/canopsis/bin/canopsis-systemd restart
```

Vérifier que les bacs à alarmes sont toujours fonctionnels

## Mise à jour de la gestion des droits

### Pré-requis :  effectuer une sauvegarde de la base de données

La gestion des droits nécessite une mise à jour de la base. L'exécutable `canopsinit` a été créé pour gérer automatiquement les mises à jour, et intégré à Canopsis 2.6.0

Pour les installations existantes, il est nécessaire d'ajouter un flag d'initialisation en base avant de commencer à utiliser `canopsinit` :

- se connecter via un shell mongodb au Primary Mongodb
- Vérifier si la collection `initialized` existe :

```javascript
db.getCollectionNames()
```

- Si elle n'existe pas, la créer :

```javascript
db.initialized.insert({"_id": "initialized", "at": "WED, 18 APR 2018 11:50:00 +0000"})
// Vous pouvez changer la date à la date du jour.
```

### Mise à jour de la base de données des  droits

se connecter à un serveur canopsis en tant qu'utilisateur canopsis et exécuter `canopsinit`

```bash
su - canopsis
canopsinit
```

## Correction d'un problème de performances des moteurs

Se connecter au primary MongoDB via un shell
exécuter la commande :

```javascript
db.getCollection('default_entities').ensureIndex({"type":1});
```

## Création d'un filtre sur les pbehaviors

- pour n'obtenir que les alarmes ayant un pbehavior actif, ajouter un filtre dans le bac à alarmes avec le champ `has_active_pb`: `true`
- pour n'obtenir que les alarmes n'ayant pas de pbehavior actif, ajouter un filtre dans le bac à alarmes avec le champ `has_active_pb`: `false`
