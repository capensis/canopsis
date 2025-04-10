# Guide de migration vers Canopsis 25.04.0

Ce guide donne les instructions vous permettant de mettre à jour Canopsis 24.04 (dernière version disponible) vers [la version 25.04.0](../25.04.0.md).

## Prérequis

L'ensemble de cette procédure doit être lu avant son exécution.

Ce document ne prend en compte que Canopsis Community et Canopsis Pro : tout développement personnalisé dont vous pourriez bénéficier ne fait pas partie du cadre de ce Guide de migration.

Les fichiers de référence qui sont mentionnés dans ce guide sont disponibles à ces adresses

| Édition           | Sources                                                                                                                              |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| Édition Community | [https://git.canopsis.net/canopsis/canopsis-community/-/releases](https://git.canopsis.net/canopsis/canopsis-community/-/releases)   |
| Édition pro       | [https://git.canopsis.net/sources/canopsis-pro-sources/-/releases](https://git.canopsis.net/sources/canopsis-pro-sources/-/releases) |

[TOC]

## Procédure de mise à jour

### Réalisation d'une sauvegarde

Des sauvegardes sont toujours recommandées, qu'elles soient régulières ou lors de modifications importantes.

La restructuration apportée dans les bases de données pour cette version de Canopsis nous amène à insister d'autant plus sur ce point. Il est donc fortement recommandé de réaliser une **sauvegarde complète** des VM hébergeant vos services Canopsis, avant cette mise à jour.

### Vérification MongoDB

!!! warning "Vérification"

    Avant de démarrer la procédure de mise à jour, vous devez vérifier que la valeur de `featureCompatibilityVersion` est bien positionnée à **7.0**  

    === "Docker Compose"
        ```sh
        CPS_EDITION=pro docker compose exec mongodb bash
        mongosh -u root -p root
        > db.adminCommand( { getParameter: 1, featureCompatibilityVersion: 1 } )
        > exit
        ```

    === "Paquets RHEL 8"

        ```sh
        mongosh -u root -p root
        > db.adminCommand( { getParameter: 1, featureCompatibilityVersion: 1 } )
        > exit
        ```

    === "Helm"

        Les commandes ci-dessous s'appliquent uniquement si votre instance MongoDB est hébergée sur votre cluster Kubernetes.
        Dans le cas contraire, veuillez vous référer à l'onglet "Paquets RHEL 8".

        ```sh
        export MONGODB_ROOT_PASSWORD=$(kubectl get secret canopsis-mongodb -o jsonpath='{.data.mongodb-root-password}' | base64 --decode)
        
        kubectl exec canopsis-mongodb-0 -- mongosh -u root -p $MONGODB_ROOT_PASSWORD --eval 'db.adminCommand({ getParameter: 1, featureCompatibilityVersion: 1 })'
        ```

        Le retour doit être de la forme `{ "featureCompatibilityVersion" : { "version" : "7.0" }, "ok" : 1 }`
        Si ce n'est pas le cas, vous ne pouvez pas continuer la mise à jour.


### Arrêt de l'environnement en cours de lancement

Vous devez prévoir une interruption du service afin de procéder à la mise à jour qui va suivre.

=== "Docker Compose"

    ```sh
    CPS_EDITION=pro docker compose down
    ```

=== "Paquets RHEL 8"

    ```sh
    systemctl stop canopsis
    systemctl stop mongod
    systemctl stop postgresql-15
    systemctl stop rabbitmq-server
    systemctl stop redis
    ```

=== "Helm"

    ```sh
    kubectl get deployments | grep -v timescaledb | awk 'NR>1 {print $1}' | xargs kubectl delete deployment
    ```


## Mise à jour Canopsis

!!! information "Information"

    Canopsis 25.04 est livré avec un nouveau jeu de configurations de référence.
    Vous devez télécharger ces configurations et y reporter vos personnalisations.  

=== "Docker Compose"

    Si vous êtes utilisateur de l'édition `community`, voici les étapes à suivre.

    Télécharger le paquet de la version 24.10.0 (canopsis-community-docker-compose-25.04.0.tar.gz) disponible à cette adresse [https://git.canopsis.net/canopsis/canopsis-community/-/releases](https://git.canopsis.net/canopsis/canopsis-community/-/releases).

    ```sh
    export CPS_EDITION=community
    tar xvfz canopsis-community-docker-compose-25.04.0.tar.gz
    cd canopsis-community-docker-compose-25.04.0
    ```

    Si vous êtes utilisateur de l'édition `pro`, voici les étapes à suivre.

    Télécharger le paquet de la version 25.04.0 (canopsis-pro-docker-compose-25.04.0.tar.gz) disponible à cette adresse [https://git.canopsis.net/sources/canopsis-pro-sources/-/releases](https://git.canopsis.net/sources/canopsis-pro-sources/-/releases).

    ```sh
    export CPS_EDITION=pro
    tar xvfz canopsis-pro-docker-compose-25.04.0.tar.gz
    cd canopsis-pro-docker-compose-25.04.0
    ```

    À ce stade, vous devez synchroniser les modifications réalisées sur vos anciens fichiers de configuration `docker-compose` avec les fichiers `docker-compose.yml` et/ou `docker-compose.override.yml`.

=== "Paquets RHEL 8"

    Non concerné car ces configurations sont livrées directemement dans les paquets RPM.

=== "Helm"

    Non concerné car ces configurations sont livrées directement dans les charts Helm.

### Mise à jour de RabbitMQ

Dans cette version de Canopsis, le bus de données RabbitMQ passe de la version 3 à 4.

=== "Docker Compose"

=== "Paquets RHEL 8"

=== "Helm"


### Lancement du provisioning `canopsis-reconfigure`

#### Synchronisation du fichier de configuration `canopsis.toml` ou fichier de surcharge

Si vous avez modifié le fichier `canopsis.toml` (vous le voyez via une définition de volume dans votre fichier docker-compose.yml), vous devez vérifier qu'il soit bien à jour par rapport au fichier de référence.  

* [`canopsis.toml` pour Canopsis Community 25.04.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/25.04.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-community.toml)
* [`canopsis.toml` pour Canopsis Pro 25.04.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/25.04.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-pro.toml)

!!! information "Information"

    Pour éviter ce type de synchronisation fastidieuse, la bonne pratique est d'utiliser [un fichier de surcharge de cette configuration](../../../guide-administration/administration-avancee/modification-canopsis-toml/). 

Si vous avez utilisé un fichier de surcharge, alors vous n'avez rien à faire, uniquement continuer à le présenter dans un volume.

#### Séparation des flux d’événements par initiateur

À partir de la version 25.04, les moteurs Canopsis traitent les événements selon trois flux distincts, en fonction de la valeur du champ initiator :

* user : événements générés par les utilisateurs via l’interface (UI)
* system : événements internes du moteur (ex. : PBH, triggers, alarm updates…)
* external : événements émis par les connecteurs

Chaque moteur dispose désormais de processeurs dédiés à chacun de ces flux, permettant leur exécution en parallèle. De plus, chaque processeur embarque un pool de workers permettant un traitement concurrent efficace.

Objectif : éviter que les événements système ne bloquent les actions utilisateurs ou les événements de connecteurs. Par exemple, même en cas de nombreux pbh_enter, l’interface reste réactive et les nouvelles alarmes peuvent être créées sans latence.

Dépréciation de certains flags
Les flags suivants sont désormais obsolètes :

* -publishQueue
* -consumeQueue
* -workers (remplacé par des flags spécifiques à chaque type de flux)

| Moteur               | Flags (nouveaux)   | Valeur par défaut | Flags obsolètes                             |
|----------------------|--------------------|-------------------|---------------------------------------------|
| engine-fifo          | -workers           | 10                | -publishQueue, -consumeQueue                |
| engine-che           | -externalWorkers   | 4                 | -workers, -publishQueue, -consumeQueue      |
|                      | -systemWorkers     | 4                 |                                             |
|                      | -userWorkers       | 2                 |                                             |
| engine-axe           | -externalWorkers   | 4                 | -workers, -publishQueue                     |
|                      | -systemWorkers     | 4                 |                                             |
|                      | -userWorkers       | 2                 |                                             |
|                      | -rpcWorkers        | 4                 |                                             |
| engine-correlation   | -externalWorkers   | 4                 | -workers, -publishQueue, -consumeQueue      |
|                      | -systemWorkers     | 4                 |                                             |
|                      | -userWorkers       | 2                 |                                             |
|                      | -rpcWorkers        | 4                 |                                             |
| engine-dynamic-infos | -externalWorkers   | 4                 | -workers, -publishQueue                     |
|                      | -systemWorkers     | 4                 |                                             |
|                      | -userWorkers       | 2                 |                                             |
|                      | -rpcWorkers        | 4                 |                                             |
| engine-action        | -externalWorkers   | 4                 | -workers                                    |
|                      | -systemWorkers     | 4                 |                                             |
|                      | -userWorkers       | 2                 |                                             |
|                      | -rpcWorkers        | 4                 |                                             |


#### Reconfiguration de Canopsis

=== "Docker Compose"

    !!! Attention

        Si vous avez personnalisé la ligne de commande de l'outil `canopsis-reconfigure`, nous vous conseillons de supprimer cette personnalisation.
        L'outil est en effet pré paramétré pour fonctionner naturellement.

    ```sh
    CPS_EDITION=pro docker compose up -d reconfigure
    ```

    !!! information "Information"

        Cette opération peut prendre plusieurs minutes pour s'exécuter.

    Vous pouvez ensuite vérifier que le mécanisme de provisioning/reconfigure s'est correctement déroulé. Le conteneur doit présenté un "exit 0"

    ```sh
    CPS_EDITION=pro docker compose ps -a|grep reconfigure
    canopsis-pro-reconfigure-1            "/canopsis-reconfigu…"   reconfigure            exited (0)
    ```

=== "Paquets RHEL 8"

    La commande `canopsis-reconfigure` doit être exécutée après mise à jour de Canopsis dans le cadre d'installation par paquets RPM.

=== "Helm"

   Non concerné, `canopsis-reconfigure` est lancé automatiquement lors de l'upgrade.

### Mise à jour et démarrage final de Canopsis

Enfin, il vous reste à mettre à jour et à démarrer tous les composants applicatifs de Canopsis

=== "Docker Compose"

    ```sh
    CPS_EDITION=pro docker compose up -d
    ```

    Vous pouvez ensuite vérifier que l'ensemble des conteneurs soient correctement exécutés.

    ```sh
    CPS_EDITION=pro docker compose ps
    ```

=== "Paquets RHEL 8"

    Mise à jour de Canopsis

    ```sh
    dnf install canopsis-pro-25.04.0 canopsis-webui-25.04.0
    ```

    Reconfiguration de Canopsis

    !!! Attention

        Si vous avez personnalisé la ligne de commande de l'outil `canopsis-reconfigure`, nous vous conseillons de supprimer cette personnalisation.
        L'outil est en effet pré paramétré pour fonctionner naturellement.


    Si vous utilisez un fichier d'override du canopsis.toml, veuillez ajouter à la ligne de commande suivante l'option `-override` suivie du chemin du fichier en question.

    ```sh
    systemctl start postgresql-15
    set -o allexport ; source /opt/canopsis/etc/go-engines-vars.conf
    /opt/canopsis/bin/canopsis-reconfigure -migrate-postgres=true -migrate-mongo=true -edition pro
    ```

    !!! information "Information"

        Cette opération peut prendre plusieurs minutes pour s'exécuter.

    Vous pouvez ensuite vérifier que le mécanisme de reconfigure s'est correctement déroulé en lisant les logs sur la sortie standard de la commande.

    Redémarrage de Canopsis

    ```sh
    systemctl restart canopsis
    ```

    Vous pouvez ensuite vérifier que l'ensemble des services soient correctement exécutés.

    ```sh
    systemctl status canopsis
    ```

=== "Helm"

    Définir le nom de votre instance :
    
    ```sh
    export RELEASE_NAME="canopsis-prod"
    ```

    Mise à jour de Canopsis :

    ```sh
    helm upgrade ${RELEASE_NAME} canopsis/canopsis-pro -f customer-values.yaml
    ```

Par ailleurs, le mécanisme de bilan de santé intégré à Canopsis ne doit pas présenter d'erreur.  

![Healthcheck](./img/25.04.0-healthcheck.png)

