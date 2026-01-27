# Service Connector-JUnit

## Introduction

Le service Connector-JUnit est un connecteur spécialisé qui permet d'intégrer les résultats de tests JUnit dans Canopsis. Il permet de suivre et d'analyser les résultats de tests automatisés et de générer des événements à partir de ces résultats.

## Options de démarrage

L'option `-h` permet d'afficher toutes les options disponibles au lancement du connecteur.

| Option | Description |
|--------|------------|
| `-cps.logger string` | Destination de sortie du logger. Remplace le paramètre "Canopsis.logger.Writer" du fichier de configuration TOML |
| `-d` | Active le mode debug |
| `-periodicalWaitTime duration` | Durée d'attente entre deux exécutions du processus périodique (défaut : 1m0s) |
| `-processArtifacts` | Active le traitement des artefacts de cas de test : captures d'écran et vidéos (défaut : true) |
| `-version` | Affiche les informations de version |

## Exemple d'utilisation

```bash
/connector-junit -d -periodicalWaitTime 2m0s
```

Cette commande lance le service Connector-JUnit en mode debug, avec une exécution du processus périodique toutes les 2 minutes.
