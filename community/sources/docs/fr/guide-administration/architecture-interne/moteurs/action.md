# Moteur ACTION

## Introduction

Le moteur ACTION a pour objectif d'élaborer des scénarios en fonction de déclencheurs. Ce moteur permet de mettre en place des réactions automatisées suite à des événements dans Canopsis. Ce moteur est disponible en édition Community.

Pour plus d'informations sur les fonctionnalités, consultez la [documentation sur les scénarios](../../../../guide-utilisation/menu-exploitation/scenarios/).

## Options de démarrage

| Option | Description |
|--------|------------|
| `-d` | Active le mode debug |
| `-externalWorkers int` | Nombre de workers pour traiter le flux d'événements "external" (défaut : 4) |
| `-fifoAckExchange string` | Exchange pour la publication des événements FIFO Ack |
| `-lastRetryInterval duration` | Intervalle de réessai de la dernière étape d'un scénario en cours d'exécution (défaut : 1m0s) |
| `-periodicalWaitTime duration` | Durée d'attente entre deux exécutions du processus périodique (défaut : 1m0s) |
| `-printEventOnError` | Affiche l'événement en cas d'erreur de traitement |
| `-rpcWorkers int` | Nombre de workers pour traiter le flux d'événements RPC (défaut : 4) |
| `-systemWorkers int` | Nombre de workers pour traiter le flux d'événements "system" (défaut : 4) |
| `-userWorkers int` | Nombre de workers pour traiter le flux d'événements "user" (défaut : 2) |
| `-version` | Affiche les informations de version |
| `-withWebhook` | *Obsolète* : Gère les actions webhook |
| `-workerPoolSize int` | Nombre de workers pour les exécutions de scénarios (défaut : 10) |
| `-workers int` | *Obsolète* : Nombre de workers pour traiter chaque flux d'événements |

## Exemple d'utilisation

```bash
/engine-action -d -workerPoolSize 15 -lastRetryInterval 2m0s
```

Cette commande lance le moteur ACTION en mode debug, avec 15 workers pour l'exécution des scénarios et un intervalle de réessai de 2 minutes pour la dernière étape d'un scénario.

## Schéma d'interactions

Nous proposons des schémas d'interactions entre ce moteur et les autres composants de Canopsis.  
[EN - View Engine Interaction Schemas](../../../../guide-developpement/schemas/engine-action/)
