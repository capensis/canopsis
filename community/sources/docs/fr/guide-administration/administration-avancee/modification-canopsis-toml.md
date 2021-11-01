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
