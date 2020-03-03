# Enchaînement des moteurs Canopsis

## Moteurs Python

L'enchaînement des moteurs Python de Canopsis se configure dans le fichier `/opt/canopsis/etc/amqp2engines.conf`.

De façon générique sur une stack Python, on aura :
```ini
[engine:nom_du_moteur]
event_processing = canopsis.[nom_du_moteur].process.event_processing
beat_processing = canopsis.[nom_du_moteur].process.beat_processing
next = [moteur_suivant],[moteur_suivant2]
```

Dans le fichier `amqp2engines.conf` il y a `event.processing` et `beat.processing` : le premier permet de lire les évènements, le second permet de configurer leur traitement périodique.

## Moteurs Go

L'enchaînement des moteurs Go de Canopsis se configure à leur lancement via l'option `-publishQueue`.

En environnement paquets CAT Go 3.34.0 et supérieur, [une procédure spéciale](../../notes-de-version/3.34.0.md#cat-activation-des-nouveaux-moteurs-engine-webhook-et-engine-dynamic-infos) doit être exécutée afin d'assurer le bon fonctionnement des moteurs `engine-webhook` et `engine-dynamic-infos`.

## Interactions avec les bases de données

Tous les moteurs peuvent communiquer avec MongoDB.

Seul `stat` communique avec InfluxDB.

Seuls `action`, `axe` et `heartbeat` communiquent avec Redis.

## Représentation

Lorsqu'un évènement entre dans le processus de traitement, il passe par la première vague de moteurs qui vont traiter et renvoyer l'information vers une seconde série de moteurs et ainsi de suite.

Le schéma suivant représente un *exemple* de configuration d'enchaînement de moteurs dans Canopsis.

![schema_moteurs](img/schema_moteurs_V3.png)

Le détail du rôle des différents moteurs est dans [la liste des moteurs](index.md#liste-des-moteurs).
