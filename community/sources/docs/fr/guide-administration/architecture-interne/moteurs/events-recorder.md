# Service Recorder

## Introduction

Le service Recorder est responsable de l'enregistrement des événements dans Canopsis. Il lit les événements à partir d'un exchange RabbitMQ et les enregistre pour une utilisation ultérieure et pour l'historisation.

Pour plus d'informations sur les fonctionnalités, consultez :  
- [Documentation les enregistrements d'événements](../../../../guide-utilisation/menu-administration/enregistrements-d-evenements/)


## Options de démarrage

| Option | Description |
|--------|------------|
| `-d` | Active le mode debug |
| `-exchange string` | Nom de l'exchange à partir duquel lire les événements (défaut : "canopsis.events") |
| `-lock duration` | Durée de verrouillage pour les instances multiples. Définissez 0s pour exécuter en tant qu'instance unique. Des valeurs courtes nécessitent des rafraîchissements fréquents. (défaut : 1m30s) |
| `-version` | Affiche la version et quitte |

## Exemple d'utilisation

```bash
/events-recorder -d -exchange "custom.events" -lock 2m0s
```

Cette commande lance le service Recorder en mode debug, en lisant les événements à partir de l'exchange "custom.events" et avec une durée de verrouillage de 2 minutes pour les instances multiples.
