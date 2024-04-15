# Modification du fichier de configuration toml `canopsis.toml`

## Description

Le fichier `canopsis.toml` regroupe la plupart des réglages fondamentaux des différents moteurs et services de Canopsis.

!!! note
    Les réglages d'exploitation « du quotidien » se situent plutôt dans l'interface web de Canopsis.

    D'autres réglages propres à certains moteurs se font au travers de leurs options de lancement (voir la documentation de chaque moteur à ce sujet) et de [variables d'environnement](variables-environnement.md).

## Emplacement

L'emplacement du fichier de configuration diffère entre les différents types d'environnement d'installation proposés par Canopsis.

| Type d'environnement | Emplacement du fichier            |
|----------------------|-----------------------------------|
| Paquets RPM                                         | `/opt/canopsis/etc/canopsis.toml` |
| Docker Compose ( Canopsis < 4.4.0 )                 | `/canopsis.toml` dans le service `reconfigure` |
| Docker Compose ( Canopsis Pro >= 4.4.0 )            | `/canopsis-pro.toml` dans le service `reconfigure` |
| Docker Compose ( Canopsis Community >= 4.4.0 )      | `/canopsis-community.toml` dans le service `reconfigure` |

!!! tip "Astuce"
    Le fichier de configuration `canopsis.toml` peut être surchargé par un autre fichier défini grâce à l'option `-override` de la commande `canopsis-reconfigure`.

### Variables d'environnement associées

La [variable d'environnement `CPS_DEFAULT_CFG`](variables-environnement.md) permet de définir un autre emplacement à utiliser pour charger ce fichier de configuration.

Il est recommandé de ne pas modifier cette valeur.

## Modification et maintenance du fichier

### Canopsis < 4.6.0

=== "En environnement paquets RPM"

    Éditez directement le fichier `/opt/canopsis/etc/canopsis.toml` (ou le fichier de surcharge `/opt/canopsis/etc/conf.d/canopsis-override.toml`), et suivez le reste de cette procédure.

    Lors de la mise à jour de Canopsis, vos modifications seront préservées par le gestionnaire de paquets `yum`. Vous devrez alors effectuer une synchronisation manuelle entre vos modifications passées et toute éventuelle nouvelle mise à jour du fichier.

=== "En environnement Docker Compose"

    Surchargez la totalité du fichier `/canopsis.toml` existant du conteneur `reconfigure`, à l'aide d'un volume.

    Lors des mises à jour de Canopsis, puisque vous effectuez une surcharge complète du fichier et que Docker Compose ne gère pas la mise à jour de fichiers de configuration, veillez tout particulièrement à comparer votre `canopsis.toml` surchargé localement avec la dernière version de `canopsis.toml` présente dans l'image de base.

### Canopsis >= 4.6.0

Depuis Canopsis 4.6.0, le fichier `canopsis-override.toml` permet de surcharger la configuration par défaut.
Ce fichier ne contiens donc que les configuration qui diffèrent avec la configuration par défaut.

=== "En environnement paquets RPM"

    Le fichier est situé au chemin suivant : `/opt/canopsis/etc/conf.d/canopsis-override.toml`.

    Lors de la mise à jour de Canopsis, vos modifications seront préservées par le gestionnaire de paquets `yum`.

=== "En environnement Docker Compose"

    Le fichier est situé dans le conteneur `reconfigure` au chemin suivant : `/opt/canopsis/etc/conf.d/canopsis-override.toml`.
    Montez y votre fichier personnalisé a l'aide d'un volume.

    Lors de la mise à jour de Canopsis, vos modifications seront préservées.



## Étape obligatoire pour la prise en compte des modifications

Après toute modification d'une valeur présente dans `canopsis.toml`, `canopsis-reconfigure` doit être relancé et les services et moteurs de Canopsis doivent être redémarrés.

=== "En environnement paquets RPM"

    Exécuter les commandes suivantes :

    ```bash
    set -o allexport ; source /opt/canopsis/etc/go-engines-vars.conf
    /opt/canopsis/bin/canopsis-reconfigure
    canoctl restart
    ```

=== "En environnement Docker Compose"

    Exécuter les commandes suivantes :

    ```sh
    docker compose restart reconfigure
    docker compose restart
    ```

## Description des options

### Section [Canopsis.global]

| Attribut                             | Exemple de valeur          | Description                          |
| :----------------------------------- | :------------------------- | :----------------------------------- |
| PrefetchCount                        | 10000                      |
| PrefetchSize                         | 0                          |
| ReconnectTimeoutMilliseconds         | 8                          | Délai de reconnexion auprès des services tiers (redis, mongodb, rabbitmq, ...)  |
| ReconnectRetries                     | 3                          | Nombre de tentative de reconnexion aux services tiers |

### Section [Canopsis.file]

| Attribut       | Exemple de valeur                        | Description                          |
| :------------- | :--------------------------------------- | :----------------------------------- |
| Upload         | "/opt/canopsis/var/lib/upload-files"     | Emplacement des fichiers uploadés. Utilisé pour le module de [remédiation](../../guide-utilisation/remediation/index.md) et des paramètres de l'interface graphique  |
| UploadMaxSize  | 314572800 # 300Mb                        | Taille maximale d'un fichier à uploader (en octet) |
| Junit          | "/opt/canopsis/var/lib/junit-files"      | Emplacement des fichiers traités par le module Junit |
| JunitApi       | "/tmp/canopsis/junit"                    | Emplacement des fichiers temporaires uploadés par le module Junit (via API) |
| SnmpMib        | ["/usr/share/snmp/mibs"]                 | Emplacement des fichiers MIB qui seront utilisés par le module SNMP |
| Icon           | "/opt/canopsis/var/lib/icons"            | Emplacement des fichiers d'icônes |
| IconMaxSize    | 10240 # 10Kb                             | Taille max des fichiers d'icônes |


### Section [Canopsis.alarm]

| Attribut                          | Exemple de valeur     | Description                          |
| :-------------------------------- | :---------------------| :----------------------------------- |
| StealthyInterval                  |                       | Encore utilisé ?          |
| :warning: Obsolète :warning: EnableLastEventDate               | true,false            | Active la mise à jour du champ `last_event_date` d'une alarme à chaque événement :warning: Depuis Canopsis 23.10, la date de dernier changement est nécessairement calculée :warning:  | 
| CancelAutosolveDelay              | "1h"                  | Délai de résolution effective d'une alarme après annulation depuis l'interface graphiqe |
| DisplayNameScheme                 | "{{ rand_string 3 }}-{{ rand_string 3 }}-{{ rand_string 3 }}" | Schéma utilisé pour générer le champ `display_name` d'une alarme |
| OutputLength                      | 255                   | Nombre maximum de caractères du champ `output` avant troncage | 
| LongOutputLength                  | 1024                  | Nombre maximum de caractères du champ `long_output` avant troncage | 
| DisableActionSnoozeDelayOnPbh     | true,false            | Si `vrai` alors le délai du snooze n'est pas ajouté à un comportement périodique |
| TimeToKeepResolvedAlarms          | "720h"                | Délai durant lequel les alarmes résolues sont conservées dans la collection principale des alarmes |
| AllowDoubleAck                    | true,false            | Permet d'acquitter plusieurs fois une alarme |
| ActivateAlarmAfterAutoRemediation | true,false            | Permet de décaler l'activation d'une alarme après l'exécution de la remédiation automatique |
| EnableArraySortingInEntityInfos   | true,false            | Active ou désactive le tri dans les listes utilisées dans les attributs d'événements. Par exemple, si un événement contient `info1=["item2", "item1"]` et que l'option est activée alors info1 vaudra en sortie `info1=["item1", "item2"]` |

### Section [Canopsis.timezone]

| Attribut | Exemple de valeur | Description                           |
| :------- | :-----------------| :------------------------------------ |
| Timezone | "Europe/Paris"    | Timezone générale du produit Canopsis |


### Section [Canopsis.data_storage]

| Attribut      | Exemple de valeur | Description                           |
| :------------ | :-----------------| :------------------------------------ |
| TimeToExecute | "Sunday,23"       | Jour et heure d'exécution de la politique de rotation des données définie dans le module `Data Storage` | 


### Section [Canopsis.import_ctx]

| Attribut            | Exemple de valeur     | Description                           |
| :------------------ | :---------------------| :------------------------------------ |
| ThdWarnMinPerImport | "30m"                 | Durée d'import au délà de laquelle une alarme mineure sera générée |
| ThdCritMinPerImport | "60m"                 | Durée d'import au délà de laquelle une alarme critique sera générée |
| FilePattern         | "/tmp/import_s.json"  | Pattern de nommage des fichiers temporaires d'import  |

### Section [Canopsis.api]

| Attribut                 | Exemple de valeur  | Description                           |
| :----------------------- | :------------------| :------------------------------------ |
| TokenSigningMethod       | "HS256"            | Méthode de signature d'un token d'authentification |
| BulkMaxSize              | 1000               | Taille maximum d'un batch (api endpoint) de changement en données en base |
| ExportMongoClientTimeout | "1m"               | Durée maximum d'un export au format CSV |
| AuthorScheme             | ["$username"]      | Permet de définir la manière de représenter l'auteur d'une action dans Canopsis. Ex : `["$username", " ", "$firstname", " ", "$lastname", " ", "$email", " ", "$_id"] ` |
| MetricsCacheExpiration   | "24h"              | Durée de validité du cache des API liées aux métriques |

### Section [Canopsis.logger]

| Attribut            | Exemple de valeur  | Description                                             |
| :------------------ | :------------------| :------------------------------------------------------ |
| Writer              | "stdout"           | Canal de sortie du logger. **`stdout`** ou **`stderr`** |

### Sous-section [Canopsis.logger.console_writer]

| Attribut            | Exemple de valeur                           | Description                                             |
| :------------------ | :-------------------------------------------| :------------------------------------------------------ |
| Enabled             | true                                        | Active ou désactive le mode [ConsoleWriter](https://github.com/rs/zerolog#pretty-logging). Si désactivé alors les messages sont loggués en JSON. |
| NoColor             | true                                        | Active ou désactive les couleurs dans les logs |
| TimeFormat          | "2006-01-02T15:04:05Z07:00"                 | Format des dates des messages de logs au format [GO](../../guide-utilisation/templates-go/index.md) |
| PartsOrder          | ["time", "level", "caller", "message"]      | Ordre des parties des messages de logs parmi "time", "level", "message", "caller", "error" |

### Section [Canopsis.metrics]

| Attribut            | Exemple de valeur  | Description                           |
| :------------------ | :------------------| :------------------------------------ |
| SliInterval         | "1h"               | Les longs intervalles de SLI sont découpés en plus petits intervalles définis par cet attribut. <br />Une valeur faible augmente la précision des métriques mais nécessite plus d'espace disque. <br />Une valeur élevée diminue la précision des métriques mais nécessaite moins d'espace disque. <br /> "1h" est la valeur recommandée dans la mesure où l'intervalle le plus petit gérée par l'interface graphique correspond à 1 heure |


### Section [Canopsis.tech_metrics]

| Attribut            | Exemple de valeur  | Description                           |
| :------------------ | :------------------| :------------------------------------ |
| Enabled             | false|true         | Active ou non la collecte des [métriques techniques](../../guide-de-depannage/metriques-techniques/index.md) |
| DumpKeepInterval    | "1h"               | Détermine le temps durant lequel les dumps seront disponibles avant leur suppression                    |


### Section [Canopsis.template.vars]

| Attribut                | Exemple de valeur  | Description                           |
| :---------------------- | :------------------| :------------------------------------ |
| system_env_var_prefixes | ["ENV_"]           | Les variables d'environnement peuvent être utilisées dans des [templates Go](../../guide-utilisation/templates-go/index.md) sous la forme `{{ .Env.System.ENV_var }}` ou dans l'interface graphique en [Handlebars](../../guide-utilisation/cas-d-usage/template_handlebars.md) sous la forme `{{ env.System.ENV_var }}`.<br />Seules les variables dont le prefixe est mentionné dans ce paramètre seront lues. |
| var1                    | "valeur1"          | Ces variables peuvent être utilisées dans des [templates Go](../../guide-utilisation/templates-go/index.md) sous la forme `{{ .Env.var }}` ou dans l'interface graphique en [Handlebars](../../guide-utilisation/cas-d-usage/template_handlebars.md) sous la forme `{{ env.var1 }}` |
 

### Section [Remediation]

| Attribute | Example | Description
| ------ | ------ | ------ |
| http_timeout | "1m" | Timeout de connexion au serveur distant |
| launch_job_retries_amount | 3 | Nombre de tentatives d'exécution du job sur le serveur distant |
| launch_job_retries_interval | "5s" | Intervalle de temps entre 2 tentative d'exécution d'un job |
| wait_job_complete_retries_amount | 12 | Nombre par défaut de tentatives de récupération du statut d'un job |
| wait_job_complete_retries_interval | 5s | Intervalle par défaut entre 2 tentatives de récupération du statut d'un job |
| pause_manual_instruction_interval | 15s | Délai d'inactivité de l'utilisateur après lequel une consigne manuelle est mise en pause |

**Exemples**

1. Rundeck est défaillant. Le moteur `remediation` essaie de se connecter à Rundeck. Après le délai `http_timeout`, la requête est considérée en échec.
1. Le moteur `remédiation` émet une requête vers Rundeck pour déclencher un job. Rundeck renvoie une erreur 500. Le moteur tente de déclencher le job `launch_job_retries_amount` fois toutes les `launch_job_retries_interval`.  
1. Le moteur `remediation` récupère le statut d'un job Rundeck. Rundeck renvoie un statut **running**. Le moteur répète cette requête `wait_job_complete_retries_amount` fois toutes les `wait_job_complete_retries_interval`.
1. Un utilisateur exécute une consigne manuelle, il ferme son navigateur. Après `pause_manual_instruction_interval`, la consigne est mise en pause

### Section [HealthCheck]

| Attribut            | Exemple de valeur  | Description                           |
| :------------------ | :------------------| :------------------------------------ |
| update_interval     | "10s"              | Intervalle de mise à jour des informations de HealthCheck | 

