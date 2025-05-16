# Moteur WEBHOOK

## Introduction

Le moteur WEBHOOK a pour objectif de gérer les appels externes via webhook. Ce moteur permet d'envoyer des notifications ou des informations à des systèmes externes via des appels HTTP. Ce moteur est disponible uniquement en édition Pro.

Pour plus d'informations sur les fonctionnalités, consultez la [documentation sur les webhooks](../../../../guide-utilisation/menu-exploitation/scenarios/#webhook).

## Options de démarrage

| Option | Description |
|--------|------------|
| `-configPath string` | Chemin du fichier de configuration du moteur Webhook (défaut : "/opt/canopsis/etc/webhook.conf.toml") |
| `-d` | Active le mode debug |
| `-lastRetryInterval duration` | Intervalle de réessai pour l'exécution d'un webhook (défaut : 1m0s) |
| `-periodicalWaitTime duration` | Durée d'attente entre deux exécutions du processus périodique (défaut : 1m0s) |
| `-version` | Affiche les informations de version |
| `-workers int` | Nombre de workers pour traiter chaque flux d'événements (défaut : 10) |

## Exemple d'utilisation

```bash
/engine-webhook -d -workers 15 -lastRetryInterval 2m0s -configPath "/opt/canopsis/etc/custom-webhook.conf.toml"
```

Cette commande lance le moteur WEBHOOK en mode debug, avec 15 workers pour traiter les événements, un intervalle de réessai de 2 minutes pour l'exécution d'un webhook et un fichier de configuration personnalisé.

## Schéma d'interactions

Nous proposons des schémas d'interactions entre ce moteur et les autres composants de Canopsis.  
[EN - View Engine Interaction Schemas](../../../../guide-developpement/schemas/engine-webhook/)
