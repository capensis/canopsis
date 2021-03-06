# Sommaire et présentation des moteurs Canopsis

Les [évènements](../../guide-utilisation/vocabulaire/index.md#evenement) envoyés par des [connecteurs](../../guide-utilisation/vocabulaire/index.md#connecteur) à Canopsis sont traités à l'aide de [moteurs](../../guide-utilisation/vocabulaire/index.md#moteur). En voici la liste.

Par défaut, ces moteurs sont open-source. Les moteurs marqués « CAT » ne sont en revanche disponibles qu'auprès d'une souscription commerciale à [Canopsis CAT](https://www.capensis.fr/canopsis/).

## Liste des moteurs Canopsis

### Enchaînement des moteurs

L'organisation de [l'enchaînement des moteurs Canopsis](schema-enchainement-moteurs.md) est décrite dans un document dédié.

### Liste des moteurs Go

La plupart des moteurs « nouvelle génération » de Canopsis sont écrits en Go.

| Moteur | Rôle | CAT ? |
|--------|------|:-----:|
| [`engine-action`](moteur-action.md) | Applique des actions définies par l'utilisateur | |
| [`engine-axe`](moteur-axe.md) | Gère le cycle de vie des alarmes | |
| [`engine-che`](moteur-che.md) | Supprime les évènements invalides, gère le contexte, et enrichit les évènements via sa fonctionnalité d'[event-filter](moteur-che-event_filter.md) | |
| [`engine-che-cat`](moteur-che.md#activation-des-plugins-denrichissement-externe-datasource) | Variante d'`engine-che`, ajoutant des plugins d'enrichissement externe | ✅ |
| [`engine-correlation`](moteur-correlation.md) | Applique et gère les règles de corrélation | ✅ |
| [`engine-dynamic-infos`](moteur-dynamic-infos.md)| Enrichit les alarmes | ✅ |
| [`engine-fifo`](moteur-fifo.md) | Garantit la cohérence et l'ordre des évènements entrant dans Canopsis | |
| [`engine-heartbeat`](moteur-heartbeat.md)  | Surveille des entités, et lève des alarmes en cas d'absence d'information | |
| [`engine-pbehavior`](moteur-pbehavior.md) | Gère les comportements périodiques | |
| [`engine-watcher`](moteur-watcher.md)| Calcule les états des observateurs | |
| [`engine-webhook`](moteur-webhook.md) | Gère le système de webhooks vers des services externes | ✅ |
<!-- Note : maintenir ce tableau dans l'ordre alphabétique -->

L'ensemble des moteurs Go se trouve dans `/opt/canopsis/bin/` lors d'une installation paquets, ou à la racine du conteneur de ce moteur, lors d'une installation Docker.

Les moteurs Go acceptent au minimum les options suivantes :

* `-d` : passage du moteur en mode *debug* ;
* `-help` : obtenir la liste complète des options acceptées par ce moteur (voyez aussi pour cela la documentation associée à chaque moteur) ;
* `-version` : obtenir les informations de version et de compilation du moteur.

### Liste des moteurs Python

Certains moteurs et composants historiques de Canopsis sont écrits en Python.

| Moteur | Rôle | CAT ? |
|--------|------|:-----:|
| [`kpi`](moteur-kpi.md) | Mise en place de statistiques sur les alarmes, entités et sessions | ✅ |
| `scheduler` | Coordonne le travail destiné aux différents moteurs `task_*` | |
| [`snmp`](moteur-snmp.md) | Gère les traps SNMP | ✅ |
| [`task_ackcentreon`](moteur-task_ackcentreon.md) | Envoi d'ACK de Canopsis vers Centreon | ✅ |
| `task_importctx` | Gestionnaire des imports de données en masse | |
<!-- Note : maintenir ce tableau dans l'ordre alphabétique -->

### Liste des anciens moteurs (non supportés)

Les moteurs suivants sont obsolètes et ne sont plus maintenus, documentés ou pris en charge.

| Moteur obsolète | Remplacé par |
|-----------------|--------------|
| `acknowledgement` (Python) | `engine-axe` (Go) |
| `alerts` (Python) | `engine-axe`(Go)  |
| `cancel` (Python) | `engine-axe` (Go) |
| `cleaner_alerts` (Python) | `engine-che` (Go) |
| `cleaner_events` (Python) | `engine-che` (Go) |
| `context` (Python) | `engine-che` (Go) |
| `context-graph` (Python) | `engine-che` (Go) |
| `engine-stat` (Go) | `statsng` (Python) ⇒ n/a |
| `eventstore` (Python) | n/a |
| `event_filter` (Python) | `engine-che` (Go) |
| `metric` (Python) | n/a |
| `pbehavior` (Python) | `engine-pbehavior` (Go) |
| `perfdata` (Python) | `metric` (Python) ⇒ n/a |
| `statsng` (Python) | n/a |
| `task_dataclean` (Python) | n/a |
| `task_linklist` (Python) | Utilisation du [linkbuilder](../linkbuilder/index.md) |
| `task_mail` (Python) | Utilisation d'un [Webhook](moteur-webhook.md) (CAT) vers un service d'envoi d'e-mails |
| `ticket` | `engine-axe`(Go) |
| `watcher` (Python) | `engine-watcher` (Go) |
<!-- Note : maintenir ce tableau dans l'ordre alphabétique -->

## Variables d'environnement liées aux moteurs

Certains comportements des moteurs de Canopsis peuvent être ajustés à l'aide de variables.

Consultez la [liste des variables d'environnement de Canopsis](../administration-avancee/variables-environnement.md) pour en savoir plus.
