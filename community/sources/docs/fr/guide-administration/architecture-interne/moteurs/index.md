# Fonctionnement des moteurs et services Canopsis

Canopsis repose sur une chaîne de moteurs interconnectés qui assurent le traitement des événements depuis les sources de supervision jusqu’à l’exécution d’actions ou d’appels externes.

Chaque moteur a un rôle bien défini dans ce pipeline. Certains sont disponibles en version Community, d'autres en Pro uniquement, et certains sont hybrides.

D'autres `services` sont mis à disposition par Canopsis : API, Connecteur Junit, Import Context Graph, Recorder.

## Enchaînement des moteurs

```mermaid
graph TD
linkStyle default interpolate basis

%% Supervision : tous les connecteurs, y compris SNMP
subgraph Supervision[Outils de supervision]
    sup[Connecteurs supervision<br/>Prometheus<br/>Zabbix<br/>Centreon...]
    snmp_src[Connecteur SNMP]
end

%% Exchanges RabbitMQ
subgraph Exchanges[Exchanges RabbitMQ]
    exch.events{canopsis.events}
    exch.snmp{canopsis.snmp}
end

%% Flux depuis supervision
sup --> exch.events
snmp_src --> exch.snmp
exch.snmp --> snmp["engine-snmp<br/>→ décodage traps SNMP"]
snmp --> exch.events

%% Chaîne de traitement
exch.events --> fifo["engine-fifo<br/>→ Chronologie des événements"]
fifo --> che["engine-che<br/>→ Filtrage/Enrichissement"]
che --> axe["engine-axe<br/>→ Gestion des alarmes"]

axe --> pbh["engine-pbehavior<br/>→ Comportements périodiques"]
axe --> correl["engine-correlation<br/>→ Corrélation"]
axe --> remediation["engine-remediation<br/>→ Remédiation"]

correl --> info["engine-dynamic-infos<br/>→ Enrichissement des Alarmes"]
info --> action["engine-action<br/>→ Scénarios"]
action --> webhook["engine-webhook<br/>→ Appels externes"]

%% Clicks
click snmp "./snmp/"
click fifo "./fifo/"
click che "./che/"
click pbh "./pbehavior/"
click axe "./axe/"
click correl "./correlation/"
click info "./dynamic-infos/"
click action "./action/"
click webhook "./webhook/"
click remediation "./remediation/"
click sup "../../../../../interconnexions/"

%% Styles
classDef grey font-weight:normal,font-size:12pt,color:#fff,fill:#878787,stroke:#222,stroke-width:3px;
classDef community-green font-weight:normal,font-size:12pt,color:#fff,fill:#2fab63,stroke:#222,stroke-width:3px;
classDef pro-blue font-weight:normal,font-size:12pt,color:#fff,fill:#2b3e4f,stroke:#222,stroke-width:3px;
classDef rabbit-orange font-weight:normal,font-size:12pt,color:#fff,fill:#ff6600,stroke:#222,stroke-width:3px;

class sup,snmp_src grey
class fifo,che,pbh,axe,action community-green
class snmp,correl,info,remediation,webhook pro-blue
class exch.events,exch.snmp rabbit-orange
```

!!! info "Informations techniques supplémentaires"

    Nous publions à titre d'informations des schémas d'interactions entre moteurs.  
    [EN - View Engine Interaction Schemas](../../../../guide-developpement/schemas/all-engines/)

{%
    include-markdown "./liste-moteurs-et-services.md"
%}
