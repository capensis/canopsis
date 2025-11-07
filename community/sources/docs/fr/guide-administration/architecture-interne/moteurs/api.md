# Service API

## Introduction

Le service API est le point d'entrée principal pour interagir avec Canopsis. Il expose les interfaces nécessaires pour gérer les entités, les alarmes, les configurations et autres fonctionnalités de la plateforme.

## Options de démarrage

L'option `-h` permet d'afficher toutes les options disponibles au lancement de l'api.

| Option | Description |
|--------|------------|
| `-c string` | Répertoire des fichiers de configuration (défaut : "/opt/canopsis/share/config") |
| `-cps.logger string` | Destination de sortie du logger. Remplace le paramètre "Canopsis.logger.Writer" du fichier de configuration TOML |
| `-d` | Active le mode debug |
| `-docs` | Active la documentation Swagger. Disponible sur l'url `/swagger` |
| `-enableSameServiceNames` | Permet d'avoir des noms de service identiques, par défaut les services ont des noms uniques |
| `-entityCategoryMetaPeriodicalWaitTime duration` | Durée d'attente entre deux exécutions du processus périodique pour mettre à jour les métadonnées de catégorie d'entité (défaut : 1m0s) |
| `-externalDataAPITimeout duration` | Délai d'attente pour les requêtes HTTP vers l'API externe (défaut : 30s) |
| `-instructionRateNotificationPeriodicalWaitTime duration` | Durée pour vérifier les instructions et créer des notifications de taux (défaut : 1h0m0s) |
| `-integrationPeriodicalWaitTime duration` | Durée d'attente pour vérifier périodiquement les résultats des tâches des moteurs (défaut : 5s) |
| `-logBody` | Active la journalisation des corps de requêtes et réponses |
| `-logBodyOnError` | Active la journalisation des corps de requêtes et réponses en cas d'erreur |
| `-periodicalWaitTime duration` | Durée d'attente entre deux exécutions du processus périodique (défaut : 1m0s) |
| `-port int` | Port du serveur (défaut : 8082) |
| `-secure` | Sécurise les sessions |
| `-stateSettingRecomputeDelay duration` | Durée minimum d'attente avant d'envoyer un événement de recalcul pour les services et composants (défaut : 1s) |
| `-version` | Affiche les informations de version |

## Exemple d'utilisation

```bash
/canopsis-api -d -port 8083 -docs -externalDataAPITimeout 45s
```

Cette commande lance le service API en mode debug sur le port 8083, avec la documentation Swagger activée et un délai d'attente de 45 secondes pour les requêtes API externes.
