# Connexion à la base de données

!!! attention
    Ce document est encore une ébauche. Il sera complété dans une future mise à jour.

## Connexion à MongoDB en ligne de commande

Depuis le nœud où est installé MongoDB, exécuter la commande shell suivante, pour se connecter avec les identifiants par défaut :

```sh
mongo -u cpsmongo -p canopsis canopsis
```

On arrive alors dans le prompt de MongoDB, permettant d'exécuter des requêtes ou des fonctions.

!!! note
    Dans le cadre d'une installation multi-nœuds avec un *Replica Set* MongoDB, veiller à se connecter au nœud primaire (`PRIMARY`).

## Connexion à MongoDB avec l'interface graphique Robo3T

Sur un poste client, installer [Robo3T](https://robomongo.org) (anciennement RoboMongo).

Vérifier que les flux sont bien ouverts entre le nœud MongoDB et le poste client.

Configurer l'interface Robo3T avec les informations données précédemment.

## Accès administrateur à MongoDB

Il est aussi possible de se connecter à la base d'administrateur de MongoDB avec les identifiants suivants :

```sh
mongo -u admin -p admin
```
