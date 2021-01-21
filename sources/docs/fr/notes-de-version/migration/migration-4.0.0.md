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

### Note importante pour les utilisateurs de paquets Debian 9

Concernant les formats d'installation, les prochaines versions de Canopsis se recentreront sur CentOS 7 et Docker Compose : les paquets Debian 9 ne seront donc bientôt plus fournis ou pris en charge.

Des paquets Debian 9 sont encore disponibles pour Canopsis 4.0.0, mais ceux-ci sont **dépréciés** et seront totalement supprimés dans une future mise à jour de Canopsis 4.1 ou 4.2. Si vous utilisez les paquets Debian 9, vous devez préparer une migration vers une des [méthodes d'installation prises en charge](../../guide-administration/installation/index.md#methodes-dinstallation-de-canopsis), à savoir CentOS 7 ou Docker Compose.

Ce Guide de migration ne prend pas en charge la migration d'un environnement Debian 9 vers une autre méthode d'installation.

## Étape 1 : vérification de votre version actuelle de Canopsis

Sur votre installation actuelle de Canopsis, rendez-vous sur la [page de connnexion](../../guide-utilisation/interface/parametres-de-linterface/index.md#3-page-de-connexion-avance), et observez le numéro de version de Canopsis dans le coin inférieur droit de l'interface. Ce numéro de version est aussi affiché à droite du logo de l'application, une fois que vous êtes connecté.

Ce numéro de version doit **obligatoirement être 3.48.0**. Si vous disposez d'une version plus ancienne de Canopsis, vous devez obligatoirement avoir [réalisé toutes les mises à jour consécutives](../../guide-administration/mise-a-jour/index.md) jusqu'à [Canopsis 3.48.0](../3.48.0.md) au préalable.

Ce Guide de migration ne prend pas en charge les environnements n'ayant pas déjà été mis à jour vers Canopsis 3.48.0.

## Étape 2 : mise à jour des dépôts et registres d'installation

Choisissez un onglet ci-dessous, en fonction de votre environnement (paquets CentOS 7, Docker Compose ou Debian 9).

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

    Ce document sera mis à jour, et une communication sera effectuée auprès des utilisateurs connus de nos images DockerHub, lorsque 

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
