# Enchainement des moteurs

L'enchainement des moteurs Canopsis se configure dans le fichier `/opt/canopsis/etc/amqp2engines.conf`.

De façon générique, on aura :
```ini
[engine:nom_du_moteur]
event_processing = canopsis.[nom_du_moteur].process.event_processing
beat_processing = canopsis.[nom_du_moteur].process.beat_processing
next = [moteur_suivant],[moteur_suivant2]
```

Dans le fichier `amqp2engines.conf` il y a `event.processing` et `beat.processing` : le premier permet de lire les événements, le second permet de configurer leur traitement périodique.

Tous les moteurs peuvent communiquer avec MongoDB et InfluxDB.

## Représentation

Lorsqu'un évènement entre dans le processus de traitement, il passe par la première vague de moteurs qui vont traiter et renvoyer l'information vers une seconde série de moteurs et ainsi de suite.

Le schéma suivant représente un *exemple* de configuration d'enchaînement de moteurs dans Canopsis.

![schema_moteurs](img/schema_moteurs_V3.png)

Le détail du rôle des différents moteurs est dans [la liste des moteurs](index.md#liste-des-moteurs).
