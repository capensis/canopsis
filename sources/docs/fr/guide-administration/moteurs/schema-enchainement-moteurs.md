# Enchaînement des moteurs Canopsis

Canopsis est constitué d'un enchaînement de moteurs Go et Python. Vous trouverez sur cette page les détails de la configuration et une représentation visuelle de cet enchaînement.

Les informations sur le rôle des différents moteurs sont dans [la liste des moteurs](index.md#liste-des-moteurs).

## Enchaînement des moteurs Go

L'enchaînement des moteurs Go de Canopsis se configure à leur lancement via l'option `-publishQueue`.

## Enchaînement des moteurs Python

L'enchaînement des moteurs Python de Canopsis se configure dans le fichier `/opt/canopsis/etc/amqp2engines.conf`.

De façon générique sur un moteur Python, on aura :
```ini
[engine:nom_du_moteur]
event_processing = canopsis.[nom_du_moteur].process.event_processing
beat_processing = canopsis.[nom_du_moteur].process.beat_processing
next = [moteur_suivant],[moteur_suivant2]
```

Dans le fichier `amqp2engines.conf` il y a `event.processing` et `beat.processing` : le premier permet de lire les évènements, le second permet de configurer leur traitement périodique.

## Représentation de l'enchaînement des moteurs

*Vous pouvez cliquer sur le nom des moteurs pour être redirigé vers la page de documentation dédiée.*

```mermaid
graph TD
linkStyle default interpolate basis
sup[Supervision] -- connecteurs :<br />centreon2canopsis<br />zabbix2canopsis<br />shinken2canopsis<br />etc--> exch.events{canopsis.events}
exch.events --> heart(engine-heartbeat)
exch.events --> fifo(engine-fifo &#40failover&#41)
fifo --> che(engine-che &#40multi-instanciable&#41)
che --> filter(event_filter)
filter --> metric(metric)
filter --> pbh(pbehavior)
pbh --> axe(engine-axe &#40multi-instanciable&#41)
axe --> correl(engine-correlation &#40multi-instanciable&#41)
axe --> watcher(engine-watcher)
correl --> watcher
watcher --> info(engine-dynamic-infos)
watcher --> action
info --> webh(engine-webhook)
webh --> action(engine-action)

click snmp "http://doc.canopsis.net/guide-administration/moteurs/moteur-snmp/"
click heart "http://doc.canopsis.net/guide-administration/moteurs/moteur-heartbeat/"
click fifo "http://doc.canopsis.net/guide-administration/moteurs/moteur-fifo/"
click che "http://doc.canopsis.net/guide-administration/moteurs/moteur-che/"
click filter "http://doc.canopsis.net/guide-administration/moteurs/moteur-che-event_filter/"
click metric "http://doc.canopsis.net/guide-administration/moteurs/moteur-metric/"
click pbh "http://doc.canopsis.net/guide-administration/moteurs/moteur-pbehavior/"
click axe "http://doc.canopsis.net/guide-administration/moteurs/moteur-axe/"
click correl "http://doc.canopsis.net/guide-administration/moteurs/moteur-correlation/"
click watcher "http://doc.canopsis.net/guide-administration/moteurs/moteur-watcher/"
click info "http://doc.canopsis.net/guide-administration/moteurs/moteur-dynamic-infos/"
click webh "http://doc.canopsis.net/guide-administration/moteurs/moteur-webhook/"
click action "http://doc.canopsis.net/guide-administration/moteurs/moteur-action/"

classDef grey font-weight:normal,font-size:12pt,color:#fff,fill:#878787,stroke:#222,stroke-width:3px;
classDef core-green font-weight:normal,font-size:12pt,color:#fff,fill:#2fab63,color:#fff,stroke:#222,stroke-width:3px;
classDef cat-blue font-weight:normal,font-size:12pt,color:#fff,fill:#2b3e4f,color:#fff,stroke:#222,stroke-width:3px;
classDef rabbit-orange font-weight:normal,font-size:12pt,color:#fff,fill:#ff6600,color:#fff,stroke:#222,stroke-width:3px;
class sup grey
class heart,fifo,che,filter,metric,pbh,axe,watcher,action core-green
class snmp,kpi,correl,info,webh cat-blue
class exch.snmp,exch.events rabbit-orange
```

Légende :
```mermaid
graph TD
leg-rabbit{Exchange<br />RabbitMQ}
leg-core(Moteur Core)
leg-cat(Moteur CAT)

classDef grey font-weight:normal,font-size:12pt,color:#fff,fill:#878787,stroke:#222,stroke-width:3px;
classDef core-green font-weight:normal,font-size:12pt,color:#fff,fill:#2fab63,color:#fff,stroke:#222,stroke-width:3px;
classDef cat-blue font-weight:normal,font-size:12pt,color:#fff,fill:#2b3e4f,color:#fff,stroke:#222,stroke-width:3px;
classDef rabbit-orange font-weight:normal,font-size:12pt,color:#fff,fill:#ff6600,color:#fff,stroke:#222,stroke-width:3px;
class leg grey
class leg-core core-green
class leg-cat cat-blue
class leg-rabbit rabbit-orange
```

!!! Note
    Certains moteurs ne sont pas représentés sur ce diagramme car leur fonctionnement est indépendant de l'enchaînement des moteurs de base. Par exemple [`snmp`](moteur-snmp.md) ou `import_ctx`.
