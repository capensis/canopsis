# Mise en place d'un environnement de développement UI

## Pré requis

* Disposer d'un webserver Canopsis fonctionnel (généralement sur le port 8082)
* Disposer d'un accès web pour récuperer les sources ainsi que les librairies complémentaires

## Récupération des sources

````
$ git clone https://git.canopsis.net/canopsis/canopsis-next.git
````

Vous pouvez choisir votre branche de développement, par exemple **develop**

````
$ cd canopsis-next
$ git checkout develop
````

## Configuration et lancement

On suppose que votre webserver Canopsis écoute le port 8082 de votre machine locale

````
$ echo "VUE_APP_API_HOST=http://localhost:8082" > .env.local
$ yarn install && yarn serve
````

## Utilisation

A ce stade, vous pouvez vous rendre avec un navigateur à l'URL indiquée par la dernière commande et utiliser l'UI Canopsis
