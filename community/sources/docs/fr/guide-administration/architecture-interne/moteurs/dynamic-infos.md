# Moteur DYNAMIC INFOS

## Introduction

Le moteur DYNAMIC INFOS a pour objectif d'enrichir les alarmes avec des informations dynamiques. Ce moteur permet d'ajouter des informations contextuelles et pertinentes aux alarmes en fonction de règles configurées. Ce moteur est disponible uniquement en édition Pro.

Pour plus d'informations sur les fonctionnalités, consultez la [documentation sur les informations dynamiques](../../../../guide-utilisation/menu-exploitation/informations-dynamiques/).

## Options de démarrage

L'option `-h` permet d'afficher toutes les options disponibles au lancement du moteur.

| Option | Description |
|--------|------------|
| `-cps.logger string` | Destination de sortie du logger. Remplace le paramètre "Canopsis.logger.Writer" du fichier de configuration TOML |
| `-d` | Active le mode debug |
| `-externalWorkers int` | Nombre de workers pour traiter le flux d'événements "external" (défaut : 4) |
| `-infosDictionaryWaitTime duration` | Durée d'attente entre deux exécutions du processus de mise à jour du dictionnaire d'informations dynamiques (défaut : 1h0m0s) |
| `-periodicalWaitTime duration` | Durée d'attente entre deux exécutions du processus périodique (défaut : 1m0s) |
| `-printEventOnError` | Affiche l'événement en cas d'erreur de traitement |
| `-rpcWorkers int` | Nombre de workers pour traiter le flux d'événements RPC (défaut : 4) |
| `-systemWorkers int` | Nombre de workers pour traiter le flux d'événements "system" (défaut : 4) |
| `-userWorkers int` | Nombre de workers pour traiter le flux d'événements "user" (défaut : 2) |
| `-version` | Affiche les informations de version |

## Exemple d'utilisation

```bash
/engine-dynamic-infos -d -infosDictionaryWaitTime 30m0s -userWorkers 4
```

Cette commande lance le moteur DYNAMIC INFOS en mode debug, avec une mise à jour du dictionnaire d'informations dynamiques toutes les 30 minutes et 4 workers pour traiter les événements utilisateur.

## Schéma d'interactions

Nous proposons des schémas d'interactions entre ce moteur et les autres composants de Canopsis.  
[EN - View Engine Interaction Schemas](../../../../guide-developpement/schemas/engine-dynamic-infos/)
