# Moteur AXE

## Introduction

Le moteur AXE a pour objectif de gérer les alarmes. Il est au cœur du système de gestion d'alarmes de Canopsis. Ce moteur est disponible en édition Community.

Pour plus d'informations sur les fonctionnalités, consultez :

- [Guide d'utilisation](../../../guide-utilisation/index.md)
- [Règles d'inactivité](../../../guide-utilisation/menu-exploitation/regles-inactivite.md)
- [Règles de bagot](../../../guide-utilisation/menu-exploitation/regles-bagot.md)
- [Règles de résolution](../../../guide-utilisation/menu-exploitation/regles-resolution.md)
- [Gestion des tags](../../../guide-utilisation/menu-administration/gestion-des-tags.md)

## Options de démarrage

L'option `-h` permet d'afficher toutes les options disponibles au lancement du moteur.

| Option | Description |
|--------|------------|
| `-cps.logger string` | Destination de sortie du logger. Remplace le paramètre "Canopsis.logger.Writer" du fichier de configuration TOML |
| `-d` | Active le mode debug |
| `-externalWorkers int` | Nombre de workers pour traiter le flux d'événements "external" (défaut : 4) |
| `-fifoAckExchange string` | Exchange pour la publication des événements FIFO Ack |
| `-periodicalWaitTime duration` | Durée d'attente entre deux exécutions du processus périodique (défaut : 1m0s) |
| `-printEventOnError` | Affiche l'événement en cas d'erreur de traitement |
| `-recomputeAllOnInit` | Recalcule les entités de type service à l'initialisation |
| `-rpcWorkers int` | Nombre de workers pour traiter le flux d'événements RPC (défaut : 4) |
| `-sliPeriodicalWaitTime duration` | Durée d'attente entre deux exécutions du processus périodique pour mettre à jour les métriques SLI (défaut : 5m0s) |
| `-softDeleteCorrPeriodicalWaitTime duration` | Durée d'attente entre deux exécutions du processus périodique pour supprimer les règles de méta-alarmes et les méta-alarmes correspondantes (défaut : 1m0s) |
| `-systemWorkers int` | Nombre de workers pour traiter le flux d'événements "system" (défaut : 4) |
| `-tagsPeriodicalWaitTime duration` | Durée d'attente entre deux exécutions du processus périodique pour mettre à jour les tags d'alarmes (défaut : 5s) |
| `-userWorkers int` | Nombre de workers pour traiter le flux d'événements "user" (défaut : 2) |
| `-version` | Affiche les informations de version |

## Exemple d'utilisation

```bash
/engine-axe -d -externalWorkers 8 -sliPeriodicalWaitTime 10m0s -tagsPeriodicalWaitTime 10s
```

Cette commande lance le moteur AXE en mode debug, avec 8 workers pour les événements externes, une mise à jour des métriques SLI toutes les 10 minutes et une mise à jour des tags d'alarmes toutes les 10 secondes.

## Schéma d'interactions

Nous proposons des schémas d'interactions entre ce moteur et les autres composants de Canopsis.  
[EN - View Engine Interaction Schemas](../../../guide-developpement/schemas/engine-axe.md)
