# Moteur SNMP

## Introduction

Le moteur SNMP a pour objectif de convertir les trap SNMP reçus du connecteur SNMP en événements Canopsis. Ce moteur est disponible uniquement en édition Pro.

Pour plus d'informations sur les fonctionnalités, consultez la [documentation sur les règles SNMP](../../../guide-utilisation/menu-exploitation/regles-snmp/).

## Options de démarrage

Le moteur SNMP est lancé avec la commande `engine-launcher` et les paramètres suivants :

| Option | Description |
|--------|------------|
| `-h`, `--help` | Affiche l'aide et les options disponibles |
| `-m MODULE`, `--module MODULE` | Nom du module Python |
| `-n NAME`, `--name NAME` | Nom du moteur |
| `-l LOGLEVEL`, `--loglevel LOGLEVEL` | Niveau de log |
| `-q QUEUE`, `--queue QUEUE` | File d'attente pour la consommation des événements |
| `-p PUBQUEUE`, `--pubqueue PUBQUEUE` | File d'attente pour la publication des événements |

## Exemple d'utilisation

```bash
engine-launcher -m snmp -n snmp -l INFO -q snmp -p che
```

Ce lancement va démarrer le moteur SNMP qui va consommer les événements de la file d'attente "snmp" et publier les événements traités vers la file "che".