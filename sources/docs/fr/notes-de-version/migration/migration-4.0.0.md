# Guide de migration vers Canopsis 4.0.0

Ce guide donne, **à titre indicatif**, des instructions vous aidant à mettre un jour un environnement Canopsis 3.48.0 vers Canopsis 4.0.0. 

Canopsis 4.0.0 étant une nouvelle [version majeure](../../guide-administration/mise-a-jour/numeros-version-canopsis.md) de l'outil, et de profonds changements ayant eu lieu, ces notes ne sauraient être exhaustives ou garanties, comme cela peut être le cas lors des mises à jour standard de Canopsis.

!!! information
    Si vous n'effectuez pas une mise à jour, mais une installation de Canopsis v4, cette procédure ne s'applique pas, et seul le [Guide d'installation](../../guide-administration/installation/index.md) vous concerne.

## Prérequis

L'ensemble de cette procédure doit être lu avant son exécution.

Ce document ne prend en compte que Canopsis Core et Canopsis CAT : tout développement personnalisé dont vous pourriez bénéficier ne fait pas partie du cadre de ce Guide de migration.

Vous devez [réaliser une sauvegarde](../../guide-administration/administration-avancee/actions-base-donnees.md#sauvegarde) de votre base de données existante. Réalisez aussi une sauvegarde de tout fichier de configuration que vous auriez personnalisé (soit à l'aide d'un volume de configuration ajouté dans Docker Compose, soit en faisant une sauvegarde de `/opt/canopsis/etc` dans un environnement par paquets). Si vous utilisez des machines virtuelles, vous êtes fortement incités à y réaliser des *snapshots* de votre environnement v3 au complet.

Fonctionnellement, vous ne devez plus dépendre d'un [ancien moteur Canopsis](../../guide-administration/moteurs/index.md#liste-des-anciens-moteurs-non-supportes) : la procédure qui suit les désactive obligatoirement, et plus aucun support n'est assuré pour les environnements v4 où ces moteurs seraient encore activés.

!!! note
    Ainsi, à titre d'exemple, si vous utilisiez encore des règles d'event-filter Python, ces règles doivent au préalable avoir déjà toutes été migrées au format des event-filters Go, avant de migrer vers Canopsis v4.

D'autres prérequis ont aussi été mis à jour. Vérifiez que vous respectez toujours :

* les [prérequis d'utilisation de Docker Compose](../../guide-administration/installation/installation-conteneurs.md#prerequis) ;
* les [prérequis de version de navigateur pris en charge](../../guide-utilisation/limitations/index.md#compatibilite-des-anciens-navigateurs).

### Note importante pour les utilisateurs de paquets Debian 9

Concernant les formats d'installation, les prochaines versions de Canopsis se recentreront sur CentOS 7 et Docker Compose : les paquets Debian 9 ne seront donc bientôt plus fournis ou pris en charge.

Des paquets Debian 9 sont encore disponibles pour Canopsis 4.0.0, mais ceux-ci sont **dépréciés** et seront totalement supprimés dans une future mise à jour de Canopsis 4.1 ou 4.2. Si vous utilisez les paquets Debian 9, vous devez préparer une migration vers une des [méthodes d'installation prise en charge](../../guide-administration/installation/index.md#methodes-dinstallation-de-canopsis), à savoir CentOS 7 ou Docker Compose.

Ce Guide de migration ne prend pas en charge la migration d'un environnement Debian 9 vers une autre méthode d'installation.

## Étape 1 : vérification de votre version actuelle de Canopsis

Sur votre installation actuelle de Canopsis, rendez-vous sur la [page de connexion](../../guide-utilisation/interface/parametres-de-linterface/index.md#3-page-de-connexion-avance), et observez le numéro de version de Canopsis dans le coin inférieur droit de l'interface. Ce numéro de version est aussi affiché à droite du logo de l'application, une fois que vous êtes connecté.

Ce numéro de version doit **obligatoirement être 3.48.0 et déjà utiliser les moteurs Go**. Si vous disposez d'une version plus ancienne de Canopsis, vous devez obligatoirement avoir [réalisé toutes les mises à jour consécutives](../../guide-administration/mise-a-jour/index.md) jusqu'à [Canopsis 3.48.0](../3.48.0.md) au préalable.

Ce Guide de migration ne prend pas en charge les environnements n'ayant pas déjà été mis à jour vers Canopsis 3.48.0 avec des moteurs Go.

## Étape 2 : mise à jour des dépôts et registres d'installation

Choisissez un onglet ci-dessous, en fonction de votre environnement (paquets CentOS 7, Docker Compose ou Debian 9).

=== "CentOS 7"

    Les dépôts de paquets Canopsis v4 ont été déplacés dans une autre arborescence de `repositories.canopsis.net`. Exécutez les commandes suivantes pour appliquer cette mise à jour.

    ```sh
    rm -f /etc/yum.repos.d/canopsis*.repo

    echo "[canopsis]
    name = canopsis
    baseurl=https://repositories.canopsis.net/pulp/repos/centos7-canopsis4/
    gpgcheck=0
    enabled=1" > /etc/yum.repos.d/canopsis.repo
    ```

    Si vous bénéficiez d'une souscription à Canopsis CAT, vous devez aussi mettre à jour ses dépôts :

    ```sh
    echo "[canopsis-cat]
    name = canopsis-cat
    baseurl=https://repositories.canopsis.net/pulp/repos/centos7-canopsis4-cat/
    gpgcheck=0
    enabled=1" > /etc/yum.repos.d/canopsis-cat.repo
    ```

    Il est aussi recommandé de forcer une mise à jour de vos caches Yum, sauf si cela ne fait pas partie de vos pratiques de maintenance :

    ```sh
    yum clean all
    yum makecache
    ```
 
=== "Docker Compose"

    Le registre Docker de Canopsis sera migré de [DockerHub](https://hub.docker.com/u/canopsis/) vers un registre interne en février 2021.

    À la date de publication de ce Guide de migration, les URL et accès aux images Docker restent inchangés pour le moment.

    Ce document sera mis à jour, et une communication sera effectuée auprès des utilisateurs connus de nos images DockerHub, lorsque ce nouveau registre sera disponible.

=== "Debian 9"

    !!! attention
        Rappel important : [l'environnement Debian 9 est déprécié](#note-importante-pour-les-utilisateurs-de-paquets-debian-9).

    Les dépôts de paquets Canopsis v4 ont été déplacés dans une autre arborescence de `repositories.canopsis.net`. Exécutez les commandes suivantes pour appliquer cette mise à jour.

    ```sh
    rm -f /etc/apt/sources.list.d/canopsis*.list

    echo "deb [trusted=yes] https://repositories.canopsis.net/pulp/deb/debian9-canopsis4/ stable main" > /etc/apt/sources.list.d/canopsis.list
    ```

    Si vous bénéficiez d'une souscription à Canopsis CAT, vous devez aussi mettre à jour ses dépôts :

    ```sh
    echo "deb [trusted=yes] https://repositories.canopsis.net/pulp/deb/debian9-canopsis4-cat/ stable main" > /etc/apt/sources.list.d/canopsis-cat.list
    ```

## Étape 3 : coupure du service

Les changements architecturaux étant nombreux, une **coupure du service** doit être effectuée.

=== "Paquets"

    Arrêt des moteurs :

    ```sh
    canoctl stop
    ```

=== "Docker Compose"

    Arrêt de l'ensemble de l'environnement :

    ```sh
    docker-compose down
    ```

### Étape 4 : application de la procédure de mise à jour

=== "CentOS 7"
    Appliquez la mise à jour des paquets Canopsis :

    ```sh
    yum --disablerepo="*" --enablerepo="canopsis*" update
    ```

=== "Docker Compose"
    Passer directement à l'étape suivante.

=== "Debian 9 (déprécié)"
    Appliquez l'ensemble de vos mises à jour (ciblez uniquement les paquets `canopsis*` si nécessaire) :

    ```sh
    apt update
    apt upgrade
    ```

## Étape 5 : mise à jour de la liste des moteurs

=== "Paquets"

    Désactivation de tout ancien moteur maintenant obsolète :

    ```sh
    systemctl list-units -a --type=service --plain --no-legend "canopsis*" | awk '/ackcentreon/ || /dynamic-pbehavior/ || /event_filter/ || /metric/ || /webserver/ { print $1 }' | xargs systemctl disable
    ```

    Activation des nouveaux moteurs Core :

    ```sh
    systemctl enable canopsis-engine-go@engine-pbehavior
    systemctl enable canopsis-service@canopsis-api canopsis-service@canopsis-oldapi
    
    grep -q ^CPS_API_URL= /opt/canopsis/etc/go-engines-vars.conf || echo "CPS_API_URL=http://localhost:8082" >> /opt/canopsis/etc/go-engines-vars.conf
    grep -q ^CPS_OLD_API_URL= /opt/canopsis/etc/go-engines-vars.conf || echo "CPS_OLD_API_URL=http://localhost:8081" >> /opt/canopsis/etc/go-engines-vars.conf
    
    cp /opt/canopsis/etc/amqp2engines-core.conf.example /opt/canopsis/etc/amqp2engines.conf
    ```

    Puis, si et seulement si vous utilisez CAT :

    ```sh
    cp /opt/canopsis/etc/amqp2engines-cat.conf.example /opt/canopsis/etc/amqp2engines.conf
    
    systemctl enable canopsis-engine-go@engine-correlation
    systemctl enable canopsis-service@external-job-executor
    
    mkdir -p /etc/systemd/system/canopsis-engine-go@engine-axe.service.d
    cat > /etc/systemd/system/canopsis-engine-go@engine-axe.service.d/axe.conf << EOF
    [Service]
    ExecStart=
    ExecStart=/usr/bin/env /opt/canopsis/bin/%i -publishQueue Engine_correlation
    EOF
    
    mkdir -p /etc/systemd/system/canopsis-engine-go@engine-watcher.service.d
    cat > /etc/systemd/system/canopsis-engine-go@engine-axe.service.d/watcher.conf << EOF
    [Service]
    ExecStart=
    ExecStart=/usr/bin/env /opt/canopsis/bin/%i -publishQueue Engine_dynamic_infos
    EOF
    ```

=== "Docker Compose"

    Si vous utilisez Canopsis Core, vous devez manuellement synchroniser l'ensemble de vos fichiers Docker Compose avec les fichiers de référence disponibles ici : <https://git.canopsis.net/canopsis/canopsis/-/tree/4.0.0/docker-compose>. Prenez bien garde à synchroniser chacun des fichiers `env` et `yml`.
    
    Si vous bénéficiez d'une souscription Canopsis CAT, rapprochez-vous de votre contact habituel pour obtenir plus d'information sur la mise à jour de ces fichiers.
    
    Dans le fichier `.env`, assurez-vous de bien avoir `CANOPSIS_IMAGE_TAG=4.0.0`, ainsi que les nouvelles variables `CPS_API_URL` et `CPS_OLD_API_URL`. La variable `CPS_WEBSERVER=1` doit aussi être renommée en `CPS_OLD_API=1` là où elle était déjà utilisée.

    **Note :** si vous utilisiez le conteneur `canopsis/uiv3`, celui-ci n'est plus disponible et doit être remplacé par l'image `canopsis/nginx`. Faites aussi attention à la chaîne `provisionning` (deux *n*) qui a été corrigée en `provisioning` (un seul *n*) dans ce fichier.

## Étape 6 : mise à jour des fichiers de configuration principaux

### `webserver.conf` vers `oldapi.conf`

Dans Canopsis v3, la gestion des API historiques se faisait dans le fichier `webserver.conf`. En v4, ce fichier a été renommé en `oldapi.conf`.

Si vous aviez apporté des modifications locales à ce fichier, vous devez le renommer en `oldapi.conf`, après l'avoir resynchronisé avec le fichier de référence suivant : <https://git.canopsis.net/canopsis/canopsis/-/blob/4.0.0/sources/canopsis/etc/oldapi.conf>.

De la même façon, les fichiers de logs `webserver*.log` ont été renommés en `oldapi*.log`. Adaptez vos éventuels logrotates à cet effet.

### `default_configuration.toml` vers `canopsis.toml`

L'ancien fichier de configuration des moteurs de Canopsis, `default_configuration.toml` a été profondément revu et a été renommé en `canopsis.toml`. Ce nouveau fichier remplace aussi l'ancien fichier de configuration `initialisation.toml`, devenu obsolète.

Vous devez, ici aussi, vous baser sur le nouveau fichier `canopsis.toml` installé par défaut, si vous modifiez certains de ces réglages par défaut.

Les fichiers de référence (pour Core et CAT) sont aussi disponibles à cette adresse : <https://git.canopsis.net/canopsis/go-engines/-/tree/develop/cmd/canopsis-reconfigure>.

Après toute modification du fichier `canopsis.toml`, vous devez relancer l'outil `canopsis-reconfigure` afin que ces changements soient pris en compte.

### Configuration de Nginx

La configuration de Nginx a été revue en v4. La refonte qui a été opérée est indispensable au bon fonctionnement de Canopsis.

Vous devez utiliser cette nouvelle configuration, et n'y apporter des changements que s'ils sont strictement nécessaires.

=== "Paquets"

    Exécutez les commandes suivantes pour installer les nouveaux fichiers de référence :

    ```sh
    mv /etc/nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf.oldv3
    cp /opt/canopsis/deploy-ansible/playbook/roles/canopsis/templates/nginx/cors.j2 /etc/nginx/cors.inc
    cp /opt/canopsis/deploy-ansible/playbook/roles/canopsis/templates/nginx/resolvers.j2 /etc/nginx/resolvers.inc
    sed -e 's,{{ CPS_API_URL }},http://127.0.0.1:8082,g' /opt/canopsis/deploy-ansible/playbook/roles/canopsis/templates/nginx/ > /etc/nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf
    ```

    Puis, vérifiez si certaines de vos modifications locales de `/etc/nginx/conf.d/default.conf.oldv3` doivent être reprises dans le nouveau fichier `/etc/nginx/conf.d/default.conf`.

=== "Docker Compose"

    Si vous surchargiez la configuration de Nginx, veuillez repartir des fichiers de configuration par défaut de la v4, et appliquer toute modification qui serait encore nécessaire.

## Étape 7 : vérification de l'URL d'appel aux API Canopsis

Cette étape n'est à suivre que si vous utilisez des scripts tiers appelant les API Canopsis.

Si c'est le cas, vérifiez que ces scripts interrogent bien les API Canopsis au travers d'une URL de ce type :

```
http://localhost:8082/api/...
```

et non pas une URL de ce type :

```
http://localhost/api/
```

En effet, l'API Canopsis doit toujours être interrogée sur son port `8082`. Canopsis v3 tolérait les appels à l'API au travers de Nginx (port `80` par défaut), mais **cette utilisation n'est plus prise en charge avec Canopsis v4**. Vous pouvez aussi avoir besoin d'ajuster vos flux réseau en conséquence.

## Étape 8 : Fin de la mise à jour

=== "Paquets"

    Si vous utilisez Canopsis Core, exécutez la commande suivante :

    ```sh
    su - canopsis -c "canopsinit --canopsis-edition core"
    ```

    Si vous utilisez Canopsis CAT, exécutez :

    ```sh
    su - canopsis -c "canopsinit --canopsis-edition cat"
    ```

    Puis, dans tous les cas, exécutez :

    ```sh
    set -o allexport ; source /opt/canopsis/etc/go-engines-vars.conf
    /opt/canopsis/bin/canopsis-reconfigure
    
    systemctl daemon-reload
    canoctl restart
    systemctl restart nginx
    ```

=== "Docker Compose"

    Relancez l'ensemble de l'environnement Docker Compose :

    ```sh
    docker-compose up -d
    ```

## Connexion à l'interface web de Canopsis

Une fois votre environnement à jour, vous pouvez à nouveau vous connecter à l'interface web de Canopsis en vous rendant sur <http://localhost> (par défaut) avec l'utilisateur `root` de Canopsis. Ce nouvel accès simplifié remplace les anciennes adresses de type `http://localhost/en/static/canopsis-next/dist/index.html#`.

Il est aussi recommandé, en parallèle, de [vous rendre sur l'interface web RabbitMQ](../../guide-de-depannage/rabbitmq-webui.md) afin de vérifier que l'ensemble des moteurs dépilent bien l'ensemble de leurs évènements en attente.

## Mise à jour des configurations de type CAS, LDAP ou SAML2

Si votre installation utilise une connexion de type CAS, LDAP ou SAML2, vous devez consulter la [documentation des méthodes d'authentification avancées](../../guide-administration/administration-avancee/methodes-authentification-avancees.md) afin de vous assurer que cette configuration est bien à jour pour une utilisation avec Canopsis v4.
