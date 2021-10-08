# Sommaire et présentation des moteurs Canopsis

Les [évènements](../../guide-utilisation/vocabulaire/index.md#evenement) envoyés par des [connecteurs](../../guide-utilisation/vocabulaire/index.md#connecteur) à Canopsis sont traités à l'aide de [moteurs](../../guide-utilisation/vocabulaire/index.md#moteur). En voici la liste.

Par défaut, ces moteurs sont open-source. Les moteurs marqués « Pro » ne sont en revanche disponibles qu'auprès d'une souscription commerciale à [Canopsis Pro](https://www.capensis.fr/canopsis/).

## Liste des moteurs Canopsis

### Enchaînement des moteurs

L'organisation de [l'enchaînement des moteurs Canopsis](schema-enchainement-moteurs.md) est décrite dans un document dédié.

### Liste des moteurs Go

La plupart des moteurs « nouvelle génération » de Canopsis sont écrits en Go.

| Moteur | Rôle | Exclusif à Canopsis Pro |
|--------|------|:-----:|
| [`engine-action`](moteur-action.md) | Applique des actions définies par l'utilisateur | |
| [`engine-axe`](moteur-axe.md) | Gère le cycle de vie des alarmes | |
| [`engine-che`](moteur-che.md) | Supprime les évènements invalides, gère le contexte, et enrichit les évènements via sa fonctionnalité d'[event-filter](moteur-che-event_filter.md) | |
| [`engine-che-cat`](moteur-che.md#activation-des-plugins-denrichissement-externe-datasource) | Variante d'`engine-che`, ajoutant des plugins d'enrichissement externe | ✅ |
| [`engine-correlation`](moteur-correlation.md) | Applique et gère les règles de corrélation | ✅ |
| [`engine-dynamic-infos`](moteur-dynamic-infos.md)| Enrichit les alarmes | ✅ |
| [`engine-fifo`](moteur-fifo.md) | Garantit la cohérence et l'ordre des évènements entrant dans Canopsis | |
| [`engine-pbehavior`](moteur-pbehavior.md) | Gère les comportements périodiques | |
| [`engine-service`](moteur-service.md)| Calcule les états des [services](../../guide-utilisation/vocabulaire/index.md#service) | |
| [`engine-webhook`](moteur-webhook.md) | Gère le système de webhooks vers des services externes | ✅ |
<!-- Note : maintenir ce tableau dans l'ordre alphabétique -->

L'ensemble des moteurs Go se trouve dans `/opt/canopsis/bin/` lors d'une installation paquets, ou à la racine du conteneur de ce moteur, lors d'une installation Docker.

Les moteurs Go acceptent au minimum les options suivantes :

* `-d` : passage du moteur en mode *debug* ;
* `-help` : obtenir la liste complète des options acceptées par ce moteur (voyez aussi pour cela la documentation associée à chaque moteur) ;
* `-version` : obtenir les informations de version et de compilation du moteur.

### Liste des moteurs Python

Certains moteurs et composants historiques de Canopsis sont écrits en Python.

| Moteur | Rôle | Exclusif à Canopsis Pro |
|--------|------|:-----:|
| [`kpi`](moteur-kpi.md) | Mise en place de statistiques sur les alarmes, entités et sessions | ✅ |
| [`snmp`](moteur-snmp.md) | Gère les traps SNMP | ✅ |
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
| `engine-heartbeat` (Go) | Utilisation des Idle rules |
| `engine-stat` (Go) | `statsng` (Python) ⇒ n/a |
| `engine-watcher` (Go) | `engine-service` (Go) |
| `metric` (Python) | n/a |
| `pbehavior` (Python) | `engine-pbehavior` (Go) |
| `perfdata` (Python) | `metric` (Python) ⇒ n/a |
| `scheduler` (Python) | n/a |
| `statsng` (Python) | n/a |
| `task_ackcentreon` (Python) | ? |
| `task_dataclean` (Python) | n/a |
| `task_importctx` (Python) | APIv4 d'import |
| `task_linklist` (Python) | Utilisation du [linkbuilder](../linkbuilder/index.md) |
| `task_mail` (Python) | Utilisation d'un [Webhook](moteur-webhook.md) (Pro) vers un service d'envoi d'e-mails |
<!-- Note : maintenir ce tableau dans l'ordre alphabétique -->

Tous les autres anciens moteurs Python sont aussi considérés obsolètes et non supportés.

## Gestion des erreurs

En cas d'erreur fatale, telle qu'une perte prolongée de la connexion à un service externe (tel que la base de données), les moteurs s'arrêtent avec un code d'erreur.

Les unités systemd (en installation paquets) et les lignes `restart: unless-stopped` (avec Docker Compose) sont configurées de manière à ce que les moteurs soient automatiquement relancés après une erreur fatale.

!!! attention
    Ceci est une partie essentielle de l'architecture de Canopsis et de la disponibilité du service. Il ne doit en aucun cas être modifié.

## Variables d'environnement liées aux moteurs

Certains comportements des moteurs de Canopsis peuvent être ajustés à l'aide de variables.

Consultez la [liste des variables d'environnement de Canopsis](../administration-avancee/variables-environnement.md) pour en savoir plus.
