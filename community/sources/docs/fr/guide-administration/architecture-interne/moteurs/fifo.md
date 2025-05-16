# Moteur FIFO

## Introduction

Le moteur FIFO a pour objectif de garantir la chronologie des événements et d'appliquer les règles de transformations d'entités. Ce moteur est disponible en édition Community.

Pour plus d'informations sur la fonctionnalité de transformation d'entités, consultez la [documentation sur les filtres d'événements](../../../../guide-utilisation/menu-exploitation/filtres-evenements/#type-change-entity).

## Options de démarrage

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

## Schéma d'interactions

Nous proposons des schémas d'interactions entre ce moteur et les autres composants de Canopsis.  
[EN - View Engine Interaction Schemas](../../../../guide-developpement/schemas/engine-fifo/)
