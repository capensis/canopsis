# Connexion à la base de données

## Connexion à MongoDB en ligne de commande

Depuis un nœud où est installé MongoDB, exécuter la commande suivante, pour se connecter avec les identifiants par défaut :
```sh
mongosh ${CPS_MONGO_URL}
```

Pour définir la valeur de `${CPS_MONGO_URL}`  :
```sh
CPS_MONGO_URL=URI_MONGODB
```
*(Se référer au contenu du fichier `go-engines-vars.conf` sur une installation RPM, `canopsis.env` sur une installation docker-compose ou le fichier `values.yml` pour une installation helm.)*

On arrive alors dans le prompt de MongoDB, permettant d'exécuter des requêtes ou des fonctions.

## Connexion à MongoDB avec l'interface graphique MongoDB Compass

Sur un poste client, installer [MongoDB Compass](https://www.mongodb.com/docs/compass/).

Vérifier que les flux sont bien ouverts entre le nœud MongoDB et le poste client.

Configurer l'interface MongoDB Compass avec les informations données précédemment.

## Accès administrateur à MongoDB

Il est aussi possible de se connecter à la base d'administrateur de MongoDB avec les identifiants suivants :

Par exemple, dans docker compose les variables sont définies ainsi: 
```sh
MONGO_INITDB_ROOT_USERNAME=root
MONGO_INITDB_ROOT_PASSWORD=root
```

Ce qui donnera pour la connexion :
```sh
mongosh -u root -p root admin
```

!!! note
    Dans le cadre d'une installation multi-nœuds avec un *Replica Set* MongoDB, veiller à bien être connecté au nœud primaire (`PRIMARY`).