# Moteur PBEHAVIOR

## Introduction

Le moteur PBEHAVIOR (Periodic Behavior) a pour objectif de gérer les changements de comportements en fonction de périodes temporelles. Il permet de gérer les maintenances, les mises en production, les périodes de service ou de non-service, etc. Ce moteur est disponible en édition Community.

Pour plus d'informations sur les fonctionnalités, consultez la [documentation sur les comportements périodiques](../../../../guide-utilisation/menu-exploitation/comportements-periodiques/).

## Options de démarrage

L'option `-h` permet d'afficher toutes les options disponibles au lancement du moteur.

| Option | Description |
|--------|------------|
| `-computeRruleEnd` | Calcule la fin des règles récurrentes (rrule) pour les comportements périodiques et quitte |
| `-cps.logger string` | Destination de sortie du logger. Remplace le paramètre "Canopsis.logger.Writer" du fichier de configuration TOML |
| `-d` | Active le mode debug |
| `-frameDuration int` | Le moteur calcule tous les comportements périodiques pour un intervalle futur dont la durée est contrôlée par ce paramètre. La valeur par défaut est de 120 minutes. Cette valeur peut être réduite si le pré-calcul utilise trop de ressources système (défaut : 120) |
| `-periodicalWaitTime duration` | Durée d'attente entre deux exécutions du processus périodique (défaut : 1m0s) |
| `-printEventOnError` | Affiche l'événement en cas d'erreur de traitement |
| `-version` | Affiche les informations de version |
| `-workers int` | Nombre de workers pour traiter chaque flux d'événements (défaut : 10) |

## Exemple d'utilisation

```bash
/engine-pbehavior -d -frameDuration 90 -periodicalWaitTime 2m0s
```

Cette commande lance le moteur PBEHAVIOR en mode debug, avec un intervalle de calcul de 90 minutes pour les comportements périodiques et une exécution du processus périodique toutes les 2 minutes.

## Schéma d'interactions

Nous proposons des schémas d'interactions entre ce moteur et les autres composants de Canopsis.  
[EN - View Engine Interaction Schemas](../../../../guide-developpement/schemas/engine-pbehavior/)
