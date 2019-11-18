# Guide de développement Canopsis

Vous trouverez ici toute la documentation nécessaire au développement sur Canopsis.

!!! tip "Note"
    Cette page contient le plan de la documentation de développement, qui est en cours d'écriture.

# Process de développement
## Organisation des dépôts
## Process de release
## Nomenclature des messages de commit
<!--  - specification des segments de canopsis (alerts, action, …) -->

# Installation d'un environnement de développement
## Python
### VM
### LXC
## Python et Go
### Docker
### VM
### LXC

# Backend
## Python
### Installation de nouvelle source python
### Structure du projet
<!--
  - organisation des packages
  - architecture à mettre en place : modele, adapter, api
-->
### Création d'engines
### Création d'API

# Golang
## Compilation
## Architecture du projet
## Création de moteurs
## Implémentation de source de données externes (pour l'event-filter)

# Base de données
<!--
## default_entities
### Présentation générale
### Présentation de la structure d'un document.
## periodical_alarms
### Présentation générale
### Présentation de la structure d'un document.
-->

# Front-end
## Mise en place de l'environnement de développement
## Technologies utilisées
## Structure du projet
## Règles de style
## Les mixins, helpers et filters
## Le store Vuex
## Guides de création nouvelle fonctionnalité
### Modal
### Vue
### Widget (+ Paramètres du widget)

# API

  * [Internal](api/api-internal.md)
  * [Action](api/api-v2-action.md)
  * [Alarm-filter](api/api-v2-alarm-filter.md)
  * [Event-filter](api/api-v2-event-filter.md)
  * [Healthcheck](api/api-v2-healthcheck.md)
  * [Heartbeat](api/api-v2-heartbeat.md)
  * [Import](api/api-v2-import.md)
  * [PBehavior](api/api-v2-pbehavior.md)
  * [Watchers](api/api-v2-watcherng.md)
  * [Webhooks](api/api-v2-webhooks.md)

# Gestion de la documentation
