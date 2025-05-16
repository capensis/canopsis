# Moteur CHE

## Introduction

Le moteur CHE a pour objectif d'appliquer les règles d'enrichissement d'événements et d'entités. Ce moteur est disponible en édition Community et Pro, avec des fonctionnalités supplémentaires pour l'enrichissement externe en édition Pro.

Pour plus d'informations sur les fonctionnalités, consultez :  
- [Documentation sur les filtres d'événements](../../../../guide-utilisation/menu-exploitation/filtres-evenements/)  
- [Documentation sur les données externes](../../../../guide-utilisation/menu-exploitation/donnees-externes/) (édition Pro)

## Options de démarrage

| Option | Description |
|--------|------------|
| `-consumeQueue string` | *Obsolète* : File d'attente pour la consommation des événements |
| `-createContext` | Active la création du context graph (activé par défaut). ATTENTION : désactivez l'ancien moteur context-graph lorsque vous utilisez cette option (défaut : true) |
| `-d` | Active le mode debug |
| `-externalDataApiTimeout duration` | Délai d'attente pour les requêtes HTTP vers les API externes (défaut : 30s) |
| `-externalWorkers int` | Nombre de workers pour traiter le flux d'événements "external" (défaut : 4) |
| `-fifoAckExchange string` | Exchange pour la publication des événements FIFO Ack |
| `-infosDictionaryWaitTime duration` | Durée d'attente entre deux exécutions du processus de mise à jour des informations des entités (défaut : 1h0m0s) |
| `-periodicalWaitTime duration` | Durée d'attente entre deux exécutions du processus périodique (défaut : 1m0s) |
| `-printEventOnError` | Affiche l'événement en cas d'erreur de traitement |
| `-processEvent` | Active le traitement des événements (activé par défaut) (défaut : true) |
| `-publishQueue string` | *Obsolète* : File d'attente pour la publication des événements |
| `-purge` | Purge la/les file(s) d'attente de consommation avant de travailler |
| `-softDeleteWaitTime duration` | Durée pendant laquelle les entités supprimées logiquement sont conservées en base de données avant d'être supprimées définitivement (défaut : 1h0m0s) |
| `-systemWorkers int` | Nombre de workers pour traiter le flux d'événements "system" (défaut : 4) |
| `-userWorkers int` | Nombre de workers pour traiter le flux d'événements "user" (défaut : 2) |
| `-version` | Affiche les informations de version |
| `-workers int` | *Obsolète* : Nombre de workers pour traiter chaque flux d'événements |

## Exemple d'utilisation

```bash
/engine-che -d -externalWorkers 6 -externalDataApiTimeout 45s -infosDictionaryWaitTime 2h0m0s
```

Cette commande lance le moteur CHE en mode debug, avec 6 workers pour les événements externes, un délai d'attente de 45 secondes pour les requêtes API et une mise à jour du dictionnaire d'informations toutes les 2 heures.


## Schéma d'interactions

Nous proposons des schémas d'interactions entre ce moteur et les autres composants de Canopsis.  
[EN - View Engine Interaction Schemas](../../../../guide-developpement/schemas/engine-che/)
