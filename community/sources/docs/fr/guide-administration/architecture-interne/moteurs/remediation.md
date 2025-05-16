# Moteur REMEDIATION

## Introduction

Le moteur REMEDIATION a pour objectif d'appliquer les remédiations associées aux alarmes. Il permet d'exécuter des actions correctives automatiques pour résoudre des problèmes identifiés par des alarmes, avec des jobs ordonnancés si nécessaire. Ce moteur est disponible uniquement en édition Pro.

Pour plus d'informations sur les fonctionnalités, consultez la [documentation sur la remédiation](../../../../guide-utilisation/remediation/).

## Options de démarrage

| Option | Description |
|--------|------------|
| `-cleanUp` | Nettoie les données de remédiation |
| `-d` | Active le mode debug |
| `-lastRetryInterval duration` | Intervalle de réessai du dernier job d'une instruction de remédiation en cours d'exécution (défaut : 1m0s) |
| `-periodicalWaitTime duration` | Durée d'attente entre deux exécutions du processus périodique (défaut : 1m0s) |
| `-printEventOnError` | Affiche l'événement en cas d'erreur de traitement |
| `-version` | Affiche les informations de version |
| `-workers int` | Nombre de workers pour traiter chaque flux d'événements (défaut : 10) |

## Exemple d'utilisation

```bash
/engine-remediation -d -workers 15 -lastRetryInterval 2m0s
```

Cette commande lance le moteur REMEDIATION en mode debug, avec 15 workers pour traiter les événements et un intervalle de réessai de 2 minutes pour le dernier job d'une instruction de remédiation.

## Schéma d'interactions

Nous proposons des schémas d'interactions entre ce moteur et les autres composants de Canopsis.  
[EN - View Engine Interaction Schemas](../../../../guide-developpement/schemas/engine-remediation/)
