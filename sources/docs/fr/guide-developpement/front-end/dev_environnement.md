# Mise en place d'un environnement de développement UI

## Pré requis

* Disposer d'un webserver Canopsis fonctionnel (généralement sur le port 8082)
* Disposer d'un accès web pour récuperer les sources ainsi que les librairies complémentaires

## Récupération des sources

Les sources de l'interface Canopsis v3 se trouvent dans le sous dossier ```sources/webcore/src/canopsis-next/``` de Canopsis

Pour commencer, clonez le dépôt Canopsis :

````
$ git clone https://git.canopsis.net/canopsis/canopsis.git
````

Déplacez-vous dans le dossier contenant l'interface

````
$ cd sources/webcore/src/canopsis-next/
````

Vous pouvez choisir votre branche de développement, par exemple **develop**

````
$ git checkout develop
````

## Configuration et lancement

Avant de lancer le serveur de développement, celui-ci à besoin de connaître l'adresse du webserver Canopsis.

On suppose que votre webserver Canopsis écoute le port 8082 de votre machine locale.

Modifiez le fichier ```.env.local``` (le créer s'il n'existe pas)

````
$ echo "VUE_APP_API_HOST=http://localhost:8082" > .env.local
````

Il ne reste ensuite qu'à lancer le serveur de développement, en s'assurant que les dépendances nécessaires sont installées :

````
$ yarn && yarn serve
````

Une fois le serveur de développement démarré, votre navigateur par défaut devrait se lancer, sur l'adresse locale du serveur de développement.

## Utilisation

A ce stade, vous pouvez vous rendre avec un navigateur à l'URL indiquée par la dernière commande et utiliser l'UI Canopsis
