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
exch.events --> fifo["engine-fifo (failover)"]
fifo --> che["engine-che (multi-instanciable)"]
che --> axe["engine-axe (multi-instanciable)"]
axe --> pbh["engine-pbehavior (multi-instanciable)"]
axe --> correl["engine-correlation (multi-instanciable)"]
axe --> remediation["engine-remediation (multi-instanciable)"]
correl --> service["engine-service (multi-instanciable)"]
service --> info["engine-dynamic-infos"]
info --> action["engine-action"]
action --> webhook["engine-webhook"]

click snmp "https://doc.canopsis.net/guide-administration/moteurs/moteur-snmp/"
click heart "https://doc.canopsis.net/guide-administration/moteurs/moteur-heartbeat/"
click fifo "https://doc.canopsis.net/guide-administration/moteurs/moteur-fifo/"
click che "https://doc.canopsis.net/guide-administration/moteurs/moteur-che/"
click pbh "https://doc.canopsis.net/guide-administration/moteurs/moteur-pbehavior/"
click axe "https://doc.canopsis.net/guide-administration/moteurs/moteur-axe/"
click correl "https://doc.canopsis.net/guide-administration/moteurs/moteur-correlation/"
click service "https://doc.canopsis.net/guide-administration/moteurs/moteur-service/"
click info "https://doc.canopsis.net/guide-administration/moteurs/moteur-dynamic-infos/"
click action "https://doc.canopsis.net/guide-administration/moteurs/moteur-action/"

classDef grey font-weight:normal,font-size:12pt,color:#fff,fill:#878787,stroke:#222,stroke-width:3px;
classDef community-green font-weight:normal,font-size:12pt,color:#fff,fill:#2fab63,color:#fff,stroke:#222,stroke-width:3px;
classDef pro-blue font-weight:normal,font-size:12pt,color:#fff,fill:#2b3e4f,color:#fff,stroke:#222,stroke-width:3px;
classDef rabbit-orange font-weight:normal,font-size:12pt,color:#fff,fill:#ff6600,color:#fff,stroke:#222,stroke-width:3px;
class sup grey
class fifo,che,filter,pbh,axe,service,action community-green
class snmp,kpi,correl,info,remediation,webhook pro-blue
class exch.snmp,exch.events rabbit-orange
```

Légende :
```mermaid
graph TD
leg-rabbit{Exchange<br />RabbitMQ}
leg-community(Moteur Community)
leg-pro(Moteur Pro)

classDef grey font-weight:normal,font-size:12pt,color:#fff,fill:#878787,stroke:#222,stroke-width:3px;
classDef community-green font-weight:normal,font-size:12pt,color:#fff,fill:#2fab63,color:#fff,stroke:#222,stroke-width:3px;
classDef pro-blue font-weight:normal,font-size:12pt,color:#fff,fill:#2b3e4f,color:#fff,stroke:#222,stroke-width:3px;
classDef rabbit-orange font-weight:normal,font-size:12pt,color:#fff,fill:#ff6600,color:#fff,stroke:#222,stroke-width:3px;
class leg grey
class leg-community community-green
class leg-pro pro-blue
class leg-rabbit rabbit-orange
```

!!! Note
    Certains moteurs ne sont pas représentés sur ce diagramme car leur fonctionnement est indépendant de l'enchaînement des moteurs de base. Par exemple [`snmp`](moteur-snmp.md) ou `import_ctx`.
