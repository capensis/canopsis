# Sommaire et présentation des moteurs Canopsis

Les [évènements](../../guide-utilisation/vocabulaire/index.md#evenement) envoyés par des [connecteurs](../../guide-utilisation/vocabulaire/index.md#connecteur) à Canopsis sont traités à l'aide de [moteurs](../../guide-utilisation/vocabulaire/index.md#moteur).

## Liste des moteurs

Les tableaux suivants décrivent l'ensemble des différents moteurs de Canopsis, dans la dernière version disponible.

Les moteurs historiques de Canopsis sont écrits en Python (il s'agit des « moteurs Python ancienne génération »). Depuis Canopsis 3, ces moteurs sont progressivement améliorés, réécrits et réarchitecturés en Go (il s'agit des « moteurs Go nouvelle génération »).

Par défaut, ces moteurs sont open-source. Les moteurs marqués « CAT » ne sont en revanche disponibles qu'auprès d'une souscription commerciale à Canopsis CAT.

### Liste des moteurs Go

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
| [`engine-watcher`](moteur-watcher.md)| Calcule les états des observateurs | |
| [`engine-webhook`](moteur-webhook.md) | Gère le système de webhooks vers des services externes | ✅ |
<!-- Note : maintenir ce tableau dans l'ordre alphabétique -->

En installation par paquets, l'ensemble des moteurs Go se trouvent dans le namespace systemd `canopsis-engine-go@`.

### Liste des moteurs Python

| Moteur | Rôle | CAT ? |
|--------|------|:-----:|
| `datametrie` | Gère le connecteur datametrie | ✅ |
| [`pbehavior`](moteur-pbehavior.md) | Gère les comportements périodiques | |
| [`event_filter` (Python)](moteur-event_filter.md) | Applique des règles de filtrage. Ne doit pas être confondu avec le nouvel `event-filter` Go, contenu dans `engine-che` | |
| `metric` | Stocke les données de métrologie des évènements | |
| `scheduler` | Coordonne le travail destiné aux différents moteurs `task_*` | |
| [`snmp`](moteur-snmp.md) | Gère les traps SNMP | ✅ |
| [`task_ackcentreon`](moteur-task_ackcentreon.md) | Envoi d'ACK de Canopsis vers Centreon | ✅ |
| `task_importctx` | Gestionnaire des imports de données en masse | |
<!-- Note : maintenir ce tableau dans l'ordre alphabétique -->

En installation par paquets, les moteurs Python se trouvent dans le namespace systemd `canopsis-engine@` ou `canopsis-engine-cat@`. Les moteurs `datametrie` et `snmp` n'ont pas de namespace.

### Liste des moteurs obsolètes

Les moteurs suivants sont obsolètes et ne sont plus maintenus, documentés ou pris en charge. Sauf indication contraire, ces anciens moteurs étaient en Python.

| Moteur obsolète | Remplacé par |
|-----------------|--------------|
| `acknowledgement` | `engine-axe` |
| `alerts` | `engine-action` |
| `cancel` | `engine-axe` |
| `cleaner_alerts` | `engine-che` |
| `cleaner_events` | `engine-che` |
| `context` | `engine-che` |
| `context-graph` | `engine-che` |
| `engine-stat` (Go) | `statsng` ⇒ n/a |
| `eventstore` | n/a |
| `perfdata` | `metric` |
| `statsng` | n/a |
| `task_dataclean` | n/a |
| `task_linklist` | Utilisation du [linkbuilder](../linkbuilder/index.md) |
| `task_mail` | Utilisation d'un [Webhook](moteur-webhook.md) (CAT) vers un service d'envoi d'e-mails |
| `ticket` | ? |
| `watcher` | `engine-watcher` |
<!-- Note : maintenir ce tableau dans l'ordre alphabétique -->

## Enchaînement des moteurs

L'organisation de [l'enchaînement des moteurs Canopsis](schema-enchainement-moteurs.md) est décrite dans un document dédié.

## Options génériques des moteurs

L'ensemble des moteurs Go se trouvent dans `/opt/canopsis/bin/` lors d'une installation paquets, ou à la racine du conteneur de ce moteur, lors d'une installation Docker.

Les moteurs Go acceptent au minimum les options suivantes :

* `-d` : passage du moteur en mode *debug* ;
* `-help` : obtenir la liste complète des options acceptées par ce moteur (voyez aussi pour cela la [documentation associée à chaque moteur](#liste-des-moteurs-go)) ;
* `-version` : obtenir les informations de version et de compilation du moteur.
