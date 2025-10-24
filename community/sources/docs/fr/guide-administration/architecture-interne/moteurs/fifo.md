# Moteur FIFO

## Introduction

Le moteur FIFO a pour objectif de garantir la chronologie des événements et d'appliquer les règles de transformations d'entités. Ce moteur est disponible en édition Community.

Afin de garantir l’ordre strict de traitement des événements, **une seule instance du moteur FIFO** doit être active à un instant donné. 
Pour faire respecter cette contrainte, un **mécanisme de verrouillage distribué** a été mis en place pour empêcher tout lancement simultané de plusieurs instances.

Ce mécanisme, basé sur Redis, empêche le lancement simultané de plusieurs instances du moteur. 

Lors du démarrage, le moteur tente d'acquérir un verrou exclusif dans Redis. Si ce verrou n’est pas obtenu après plusieurs tentatives, le moteur s'arrête automatiquement. Le verrou est maintenu pendant l’exécution et libéré à l’arrêt propre du moteur. 

En cas d’arrêt forcé, le verrou expirera au bout de 2 minutes et 15 secondes, ou après la durée définie par `-periodicalWaitTime` + `15s` si cette dernière est plus longue. 

Ce mécanisme garantit qu’une seule instance du moteur FIFO peut s’exécuter à un instant donné.

Pour plus d'informations sur la fonctionnalité de transformation d'entités, consultez la [documentation sur les filtres d'événements](../../../../guide-utilisation/menu-exploitation/filtres-evenements/#type-change-entity).

## Options de démarrage

L'option `-h` permet d'afficher toutes les options disponibles au lancement du moteur.

| Option | Description |
|--------|------------|
| `-consumeQueue string` | *Obsolète* : File d'attente pour la consommation des événements |
| `-d` | Active le mode debug |
| `-eventsStatsFlushInterval duration` | *Obsolète* : Intervalle entre les sauvegardes des statistiques de Redis vers MongoDB (défaut : 1m0s) |
| `-externalDataApiTimeout duration` | Délai d'attente pour les requêtes HTTP vers l'API externe (défaut : 30s) |
| `-lockTtl int` | Temps de vie (TTL) du verrou Redis en secondes (défaut : 10) |
| `-periodicalWaitTime duration` | Durée d'attente entre deux exécutions du processus périodique (défaut : 1m0s) |
| `-printEventOnError` | Affiche l'événement en cas d'erreur de traitement |
| `-publishQueue string` | *Obsolète* : File d'attente pour la publication des événements |
| `-version` | Affiche les informations de version |
| `-workers int` | Nombre de workers pour traiter les événements fifo_ack (défaut : 10) |

## Exemple d'utilisation

```bash
/engine-fifo -d -periodicalWaitTime 2m0s -workers 15
```

Cette commande lance le moteur FIFO en mode debug, avec un intervalle de 2 minutes entre les exécutions périodiques et 15 workers pour traiter les événements.

**Acquisition du verrou**

Lorsque le verrou est acquis avec succès, le moteur démarre normalement et initialise ses workers :

```bash
DBG ../lib/canopsis/scheduler/scheduler.go > subscribed
INF ../lib/canopsis/engine/engine.go > engine started consumers=2 periodical_workers=5 routines=6
```

**Echec d'acquision du verrou**

Si le verrou est déjà détenu par une autre instance, le moteur échoue à démarrer et se termine immédiatement avec une erreur du type :

```bash
ERR ../cmd/engine-fifo/main.go > exit with error error="cannot init engine: redislock: not obtained"
```

## Schéma d'interactions

Nous proposons des schémas d'interactions entre ce moteur et les autres composants de Canopsis.  
[EN - View Engine Interaction Schemas](../../../../guide-developpement/schemas/engine-fifo/)
