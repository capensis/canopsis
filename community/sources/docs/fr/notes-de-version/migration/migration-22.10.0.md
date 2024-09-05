# Guide de migration vers Canopsis 22.10.0

Ce guide donne des instructions vous permettant de mettre à jour Canopsis 4.6 vers [la version 22.10.0](../22.10.0.md).


!!! warning "Avertissement"

    Seules les migrations en environnement Docker sont traitées dans ce guide.  
    Les migrations en environnement CentOS 7 seront traitées dans un second temps.

## Prérequis

L'ensemble de cette procédure doit être lu avant son exécution.

Ce document ne prend en compte que Canopsis Community et Canopsis Pro : tout développement personnalisé dont vous pourriez bénéficier ne fait pas partie du cadre de ce Guide de migration.


## Procédure de mise à jour

### Réalisation d'une sauvegarde

Des sauvegardes sont toujours recommandées, qu'elles soient régulières ou lors de modifications importantes.

La restructuration apportée dans les bases de données pour cette version de Canopsis nous amène à insister d'autant plus sur ce point. Il est donc fortement recommandé de réaliser une **sauvegarde complète** des VM hébergeant vos services Canopsis, avant cette mise à jour.


### Authentification

Les configurations des mécanismes d'authentification LDAP, CAS, SAML sont désormais définies dans le fichier de configuration du service d'API Canopsis.  
Vous devez donc reporter vos configurations dans `/opt/canopsis/share/config/api/security/config.yml`.  
Si vous êtes sur une installation Docker et que vous avez paramétré une authentification externe, alors vous avez déjà présenté le fichier `config.yml` dans un volume. Il vous reste alors à le configurer correctement.  

Pour connaitre l'ensemble des paramètres disponibles, reportez vous à [la documentation](https://doc.canopsis.net/guide-administration/administration-avancee/methodes-authentification-avancees/).


### Vérification MongoDB

!!! warning "Vérification"

    Avant de démarrer la procédure de mise à jour, vous devez vérifier que la valeur de `featureCompatibilityVersion` est bien positionnée à **4.2**  
    ```sh
    docker compose -f 00-data.docker-compose.yml exec mongodb bash
    mongo -u root -p root
    > db.adminCommand( { getParameter: 1, featureCompatibilityVersion: 1 } )
    > exit
    ```
    
    Le retour doit être de la forme `"featureCompatibilityVersion" : { "version" : "4.2" }`
    Si ce n'est pas le cas, vous ne pouvez pas continuer la mise à jour.

### Arrêt de l'environnement en cours de lancement

Vous devez prévoir une interruption du service afin de procéder à la mise à jour qui va suivre.

=== "Docker Compose"

    ```sh
    docker compose -f 00-data.docker-compose.yml -f 01-prov.docker-compose.yml -f 02-app.docker-compose.yml down
    ```

=== "Paquets CentOS 7"

    À venir


### Mise à jour Canopsis

=== "Docker Compose"

    !!! information "Information"

        Canopsis 22.10 est livré avec un nouveau jeu de configurations de référence.  
        Vous devez télécharger ces configurations et y reporter vos personnalisations.

    La configuration des versions 4.x était basée sur 3 fichiers docker-compose.  
    À présent, **un seul fichier `docker-compose.yml` est nécessaire**.

    Vous devez absolument reporter vos personnalisations vers ce nouveau fichier de configuration ou – préférablement – vers un fichier de surcharge (`docker-compose.override.yml`).

    ---

    Si vous êtes utilisateur de l'édition `community`, voici les étapes à suivre.

    ```sh
    export CPS_EDITION=community
    mkdir 22.10
    cd 22.10
    git clone https://git.canopsis.net/canopsis/canopsis-community.git -b release-22.10
    cd canopsis-community/community/deployment/canopsis/docker/
    ```

    Si vous êtes utilisateur de l'édition `pro`, voici les étapes à suivre.

    ```sh
    export CPS_EDITION=pro
    mkdir 22.10
    cd 22.10
    git clone https://git.canopsis.net/sources/canopsis-pro-sources.git -b release-22.10
    cd canopsis-pro-sources/pro/deployment/canopsis/docker/
    ```

    À ce stade, vous devez synchroniser les modifications réalisées sur vos anciens fichiers de configuration Docker Compose avec les fichiers `docker-compose.yml` et/ou `docker-compose.override.yml`.

=== "Paquets CentOS 7"

    À venir

### Mise à jour de MongoDB

Dans cette version de Canopsis, la base de données MongoDB passe de la version 4.2 à 4.4.

=== "Docker Compose"


    Démarrez le conteneur `mongodb` :

    ```sh
    docker compose up -d mongodb
    ```

    Entrez ensuite à l'intérieur de ce conteneur, afin de compléter la mise à jour vers MongoDB 4.4 :

    ```sh
    docker compose exec mongodb bash
    mongo -u root -p root
    > db.adminCommand( { setFeatureCompatibilityVersion: "4.4" } )
    exit
    ```

=== "Paquets CentOS 7"

    À venir

### Mise à jour de TimescaleDB

Dans cette version de Canopsis, la base de données TimescaleDB passe de la version 2.5.1 à 2.7.2.

=== "Docker Compose"

    Relancez le conteneur `timescaledb` :

    ```sh
    docker compose up -d timescaledb
    ```

    Puis mettez à jour l'extension timescaledb (La chaîne de connexion doit être adaptée à votre environnement)

    ```sh
    docker compose exec timescaledb psql postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis
    canopsis=# ALTER EXTENSION timescaledb UPDATE;
    ```

    Ensuite, vérifiez que l'extension en elle-même est à présent bien à jour

    ```sh
    \dx
    ...
    timescaledb | 2.7.2   | public     | Enables scalable inserts and complex queries for time-series data
    ...
    ```
   

=== "Paquets CentOS 7"

    À venir


### Remise à 0 du cache Redis

Dans cette version de Canopsis, le cache de Canopsis doit repartir à 0.

=== "Docker Compose"

    ```sh
    docker compose up -d redis
    docker compose exec redis /usr/local/bin/redis-cli flushall
    ```

=== "Paquets CentOS 7"

    À venir

### Suppression d'options de lancement de certains moteurs

#### Suppression des options `-featureHideResources`, `-postProcessorsDirectory`, `-ignoreDefaultTomlConfig` dans **engine-axe**

=== "Docker Compose"

    Supprimez toute éventuelle utilisation des options citées dans vos fichiers de référence Docker Compose.

=== "Paquets CentOS 7"

    À venir

#### Suppression des options `-dataSourceDirectory`, `enrichContext`, `enrichExclude`, `enrichInclude` dans **engine-che** 

=== "Docker Compose"

    Supprimez toute éventuelle utilisation des options citées dans vos fichiers de référence Docker Compose.

=== "Paquets CentOS 7"

    À venir

#### Suppression des options `-dataSourceDirectory`, `enableMetaAlarmProcessing` dans **engine-fifo** 

=== "Docker Compose"

    Supprimez toute éventuelle utilisation des options citées dans vos fichiers de référence Docker Compose.

=== "Paquets CentOS 7"

    À venir

#### Suppression des options `-publishQueue`, `-consumeQueue`, `-fifoAckExchange` dans **engine-pbehavior**

=== "Docker Compose"

    Supprimez toute éventuelle utilisation des options citées dans vos fichiers de référence Docker Compose.

=== "Paquets CentOS 7"

    À venir

#### Suppression des options `-c` dans **engine-remediation**

=== "Docker Compose"

    Supprimez toute éventuelle utilisation des options citées dans vos fichiers de référence Docker Compose.

=== "Paquets CentOS 7"

    À venir

### Lancement du provisioning `canopsis-reconfigure`

!!! important "Important"

    Le système de provisioning (conteneur provisioning) de Canopsis 4.x disparait au profit du système `canopsis-reconfigure` qui existait déjà.  
    Exécuter les scripts de migration MongoDB est à présent une des missions de `canopsis-reconfigure`

#### Synchronisation du fichier de configuration `canopsis.toml` ou fichier de surcharge

Si vous avez modifié le fichier `canopsis.toml` (vous le voyez via une définition de volume dans votre fichier `docker-compose.yml`), vous devez vérifier qu'il soit bien à jour par rapport au fichier de référence.  

* [`canopsis.toml` pour Canopsis Community 22.10.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/22.10.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-community.toml)
* [`canopsis.toml` pour Canopsis Pro 22.10.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/22.10.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-pro.toml)

!!! information "Information"

    Pour éviter ce type de synchronisation fastidieuse, la bonne pratique est d'utiliser [un fichier de surcharge de cette configuration](../../guide-administration/administration-avancee/modification-canopsis-toml.md). 


Si vous avez utilisé un fichier de surcharge, alors vous n'avez rien à faire, uniquement continuer à le présenter dans un volume.

#### Reconfiguration de Canopsis

=== "Docker Compose"

    !!! Attention
    
    Si vous avez personnalisé la ligne de commande de l'outil `canopsis-reconfigure`, nous vous conseillons de supprimer cette personnalisation.
    L'outil est en effet pré paramétré pour fonctionner naturellement.

    ```sh
    docker compose up -d reconfigure
    ```

    !!! information "Information"

        Cette opération peut prendre plusieurs minutes pour s'exécuter.

    Vous pouvez ensuite vérifier que le mécanisme de provisioning/reconfigure s'est correctement déroulé. Le conteneur doit présenté un "exit 0"

    ```sh
    docker compose ps
    canopsis-pro-reconfigure-1            "/canopsis-reconfigu…"   reconfigure            exited (0)
    ```

=== "Paquets CentOS 7"

    À venir

#### Démarrage final de Canopsis

Enfin, il vous reste à démarrer tous les composants applicatifs de Canopsis

=== "Docker Compose"

    ```sh
    docker compose up -d
    ```

    Vous pouvez ensuite vérifier que l'ensemble des conteneurs soient correctement exécutés.

    ```sh
    docker compose ps
    NAME                                  COMMAND                  SERVICE                STATUS              PORTS
    canopsis-pro-action-1                 "/engine-action -wit…"   action                 running             
    canopsis-pro-api-1                    "/bin/sh -c /${CMD}"     api                    running (healthy)   0.0.0.0:8082->8082/tcp, :::8082->8082/tcp
    canopsis-pro-axe-1                    "/engine-axe -publis…"   axe                    running             
    canopsis-pro-che-1                    "/engine-che"            che                    running             
    canopsis-pro-connector-junit-1        "/bin/sh -c /${CMD}"     connector-junit        running             
    canopsis-pro-correlation-1            "/bin/sh -c /${CMD}"     correlation            running             
    canopsis-pro-dynamic-infos-1          "/bin/sh -c /${CMD}"     dynamic-infos          running             
    canopsis-pro-fifo-1                   "/bin/sh -c /${CMD}"     fifo                   running             
    canopsis-pro-migrate-metrics-meta-1   "/bin/true /migrate-…"   migrate-metrics-meta   exited (0)          
    canopsis-pro-mongodb-1                "docker-entrypoint.s…"   mongodb                running (healthy)   0.0.0.0:27027->27017/tcp, :::27027->27017/tcp
    canopsis-pro-nginx-1                  "/bin/sh -c /entrypo…"   nginx                  running             80/tcp, 0.0.0.0:80->8080/tcp, :::80->8080/tcp, 0.0.0.0:443->8443/tcp, :::443->8443/tcp
    canopsis-pro-pbehavior-1              "/bin/sh -c /${CMD}"     pbehavior              running             
    canopsis-pro-rabbitmq-1               "docker-entrypoint.s…"   rabbitmq               running (healthy)   4369/tcp, 5671/tcp, 0.0.0.0:5672->5672/tcp, :::5672->5672/tcp, 15671/tcp, 25672/tcp, 0.0.0.0:15672->15672/tcp, :::15672->15672/tcp
    canopsis-pro-reconfigure-1            "/canopsis-reconfigu…"   reconfigure            exited (0)          
    canopsis-pro-redis-1                  "docker-entrypoint.s…"   redis                  running (healthy)   0.0.0.0:6379->6379/tcp, :::6379->6379/tcp
    canopsis-pro-remediation-1            "/bin/sh -c /${CMD}"     remediation            running             
    canopsis-pro-service-1                "/engine-service -pu…"   service                running             
    canopsis-pro-timescaledb-1            "docker-entrypoint.s…"   timescaledb            running (healthy)   0.0.0.0:5432->5432/tcp, :::5432->5432/tcp
    canopsis-pro-webhook-1                "/bin/sh -c /${CMD}"     webhook                running             
    ```

=== "Paquets CentOS 7"

    À venir

Par ailleurs, le mécanisme de bilan de santé intégré à Canopsis ne doit pas présenter d'erreur.  

![Healthcheck](./img/healthcheck.png)

### Règles métier

!!! warning "Avertissement"

    Dans cette version, l'ensemble des filtres existant ont été migrés vers le nouveau format.  
    Vous devez absolument vérifier le contenu de chacun des filtres pour valider que les cibles sont les bonnes.  
    Les anciens filtres ont été conservés au cas où une des migrations ne se serait pas correctement déroulée.


### Moteur SNMP et Linkbuilder

Toutes les parties Python de Canopsis ont été supprimées dans la version 22.10.  
Les fonctionnalités SNMP et Linkbuilder n'ont pas encore été migrées en GO et réintégrées dans Canopsis.  
Pour continuer à les utiliser, vous devez utiliser les images Docker des versions 4.x.  
Cela est temporaire mais nécessaire.

Que ce soit sur une installation par paquets RPM ou par image Docker, nous vous livrons des configurations à exécuter avec Docker.

Vous devez simplement définir la variable d'environnement `CPS_OLD_API` dans `/opt/canopsis/etc/go-engines-vars.conf` pour une installation paquets ou `compose.env` en installation docker.

!!! warning "Avertissement"

    Les configurations Docker Compose livrées ne doivent absolument pas être modifiées, notamment les versions d'images utilisées.  
    Par ailleurs, si vous utilisiez déjà la partie SNMP dans votre installation alors la collection `schema` existe très certainement dans mongoDB.  
    Si ce n'était pas le cas, vous allez devoir la créer avec ces instructions.  
    ```
    docker compose exec snmp /bin/bash
    schema2db
    ```

```sh
cd ../../../mock/external-services/snmp/docker
docker compose up -d
```
