# Moteur Corrélation

## Introduction

Le moteur Corrélation a pour objectif de gérer les méta-alarmes à partir de règles définies. Ce moteur permet de créer des corrélations entre différentes alarmes pour générer des méta-alarmes plus pertinentes pour les utilisateurs. Ce moteur est disponible uniquement en édition Pro.

Pour plus d'informations sur les fonctionnalités, consultez la [documentation sur les règles de méta-alarme](../../../../guide-utilisation/menu-exploitation/regles-metaalarme/).

## Options de démarrage

L'option `-h` permet d'afficher toutes les options disponibles au lancement du moteur.

| Option | Description |
|--------|------------|
| `-cps.logger string` | Destination de sortie du logger. Remplace le paramètre "Canopsis.logger.Writer" du fichier de configuration TOML |
| `-d` | Active le mode debug |
| `-externalWorkers int` | Nombre de workers pour traiter le flux d'événements externes (défaut : 4) |
| `-periodicalWaitTime duration` | Durée d'attente entre deux exécutions du processus périodique (défaut : 1m0s) |
| `-printEventOnError` | Affiche l'événement en cas d'erreur de traitement |
| `-rpcWorkers int` | Nombre de workers pour traiter le flux d'événements RPC (défaut : 4) |
| `-systemWorkers int` | Nombre de workers pour traiter le flux d'événements système (défaut : 4) |
| `-userWorkers int` | Nombre de workers pour traiter le flux d'événements utilisateur (défaut : 2) |
| `-version` | Affiche les informations de version |

## Exemple d'utilisation

```bash
/engine-correlation -d -externalWorkers 6 -periodicalWaitTime 30s
```

Cette commande lance le moteur Corrélation en mode debug, avec 6 workers pour les événements externes et une exécution du processus périodique toutes les 30 secondes.

## Schéma d'interactions

Nous proposons des schémas d'interactions entre ce moteur et les autres composants de Canopsis.  
[EN - View Engine Interaction Schemas](../../../../guide-developpement/schemas/engine-correlation/)
