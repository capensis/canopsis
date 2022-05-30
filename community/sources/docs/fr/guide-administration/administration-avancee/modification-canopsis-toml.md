# Modification du fichier de configuration `canopsis.toml`

## Description

Le fichier `canopsis.toml`regroupe la plupart des réglages fondamentaux des différents moteurs et services de Canopsis.

!!! note
    Les réglages d'exploitation « du quotidien » se situent plutôt dans l'interface web de Canopsis.

    D'autres réglages propres à certains moteurs se font au travers de leurs options de lancement (voir la documentation de chaque moteur à ce sujet) et de [variables d'environnement](variables-environnement.md).

## Emplacement

L'emplacement du fichier de configuration diffère entre les différents types d'environnement d'installation proposés par Canopsis.

| Type d'environnement | Emplacement du fichier            |
|----------------------|-----------------------------------|
| Paquets RPM          | `/opt/canopsis/etc/canopsis.toml` |
| Docker Compose       | `/canopsis.toml`                  |

### Variables d'environnement associées

La [variable d'environnement `CPS_DEFAULT_CFG`](variables-environnement.md) permet de définir un autre emplacement à utiliser pour charger ce fichier de configuration.

Il est recommandé de ne pas modifier cette valeur.

## Liste des différentes options de configuration

Certaines des valeurs pouvant être modifiées dans le fichier `canopsis.toml` sont détaillées dans d'autres pages de cette plateforme de documentation. Lancez une recherche de `canopsis.toml` dans la barre de recherche de [doc.canopsis.net](../../index.md) afin d'identifier ces diverses références.

Il n'existe, à ce jour, pas de documentation répertoriant et décrivant la totalité de ces variables.
<!-- XXX: à faire -->

## Modification et maintenance du fichier

=== "En environnement paquets RPM"

    Éditez directement le fichier `/opt/canopsis/etc/canopsis.toml`, et suivez le reste de cette procédure.

    Lors de la mise à jour de Canopsis, vos modifications seront préservées par le gestionnaire de paquets `yum`. Vous devrez alors effectuer une synchronisation manuelle entre vos modifications passées et toute éventuelle nouvelle mise à jour du fichier.

=== "En environnement Docker Compose"

    Surchargez la totalité du fichier `/canopsis.toml` existant du conteneur `reconfigure`, à l'aide d'un volume.

    Lors des mises à jour de Canopsis, puisque vous effectuez une surcharge complète du fichier et que Docker Compose ne gère pas la mise à jour de fichiers de configuration, veillez tout particulièrement à comparer votre `canopsis.toml` surchargé localement avec la dernière version de `canopsis.toml` présente dans l'image de base.

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
    docker-compose restart reconfigure
    docker-compose down
    docker-compose up -d
    ```

## Description des options

[Canopsis.global]
PrefetchCount = 10000
PrefetchSize = 0
ReconnectTimeoutMilliseconds = 8
ReconnectRetries = 3

### [Canopsis.file]

| Attribut       | Exemple de valeur                        | Description                          |
| :------------- | :--------------------------------------- | :----------------------------------- |
| Upload         | "/opt/canopsis/var/lib/upload-files"     | Emplacement des fichiers uploadés. Utilisé pour le module de [remédiation](../../remediation) et des paramètres de l'interface graphique  |
| UploadMaxSize  | 314572800 # 300Mb                        | Taille maximale d'un fichier à uploader (en octet) |
| Junit          | "/opt/canopsis/var/lib/junit-files"      | Emplacement des fichiers traités par le module Junit |
| JunitApi       | "/tmp/canopsis/junit"                    | Emplacement des fichiers temporaires uploadés par le module Junit (via API) |



### [Canopsis.alarm]

| Attribut                        | Exemple de valeur     | Description                          |
| :------------------------------ | :---------------------| :----------------------------------- |
| StealthyInterval                |                       | Encore utilisé ?          |
| EnableLastEventDate             | true,false            | Active la mise à jour du champ `last_event_date` d'une alarme à chaque événement        | 
| CancelAutosolveDelay            | "1h"                  | Délai de résolution effective d'une alarme après annulation depuis l'interface graphiqe |
| DisplayNameScheme               | "{{ rand_string 3 }}-{{ rand_string 3 }}-{{ rand_string 3 }}" | Schéma utilisé pour générer le champ `display_name` d'une alarme |
| OutputLength                    | 255                   | Nombre maximum de caractères du champ `output` avant troncage | 
| LongOutputLength                | 1024                  | Nombre maximum de caractères du champ `long_output` avant troncage | 
| DisableActionSnoozeDelayOnPbh   | true,false            | Si `vrai` alors le délai du snooze n'est pas ajouté à un comportement périodique |
| TimeToKeepResolvedAlarms        | "720h"                | Délai durant lequel les alarmes résolues sont conservées dans la collection principale des alarmes |
| AllowDoubleAck                  | true,false            | Permet d'acquitter plusieurs fois une alarme |


### [Canopsis.timezone]

| Attribut | Exemple de valeur | Description                           |
| :------- | :-----------------| :------------------------------------ |
| Timezone | "Europe/Paris"    | Timezone générale du produit Canopsis |


### [Canopsis.data_storage]

| Attribut      | Exemple de valeur | Description                           |
| :------------ | :-----------------| :------------------------------------ |
| TimeToExecute | "Sunday,23"       | Jour et heure d'exécution de la politique de rotation des données définie dans le module `Data Storage` | 


### [Canopsis.import_ctx]

| Attribut            | Exemple de valeur     | Description                           |
| :------------------ | :---------------------| :------------------------------------ |
| ThdWarnMinPerImport | "30m"                 | A compléter |
| ThdCritMinPerImport | "60m"                 | A compléter |
| FilePattern         | "/tmp/import_s.json" | A compléter |

### [Canopsis.api]

| Attribut            | Exemple de valeur  | Description                           |
| :------------------ | :------------------| :------------------------------------ |
| TokenExpiration     | "24h"              | Durée de validité d'un token d'authentification |
| TokenSigningMethod  | "HS256"            | ?? |
| BulkMaxSize         | 1000               | ?? |

[Canopsis.metrics]

| Attribut            | Exemple de valeur  | Description                           |
| :------------------ | :------------------| :------------------------------------ |
| SliInterval         | "1h"               | Les longs intervalles de SLI sont découpés en plus petits intervalles définis par cet attribut. <br />Une valeur faible augmente la précision des métriques mais nécessite plus d'espace disque. <br />Une valeur élevée diminue la précision des métriques mais nécessaite moins d'espace disque. <br /> "1h" est la valeur recommandée dans la mesure où l'intervalle le plus petit gérée par l'interface graphique correspond à 1 heure |

