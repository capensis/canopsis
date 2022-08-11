# Guide de migration vers Canopsis 4.6.0

Ce guide donne des instructions vous permettant de mettre à jour Canopsis 4.5 vers [la version 4.6.0](../4.6.0.md).

## Prérequis

L'ensemble de cette procédure doit être lu avant son exécution.

Ce document ne prend en compte que Canopsis Community et Canopsis Pro : tout développement personnalisé dont vous pourriez bénéficier ne fait pas partie du cadre de ce Guide de migration.

## Procédure de mise à jour

### Réalisation d'une sauvegarde

Des sauvegardes sont toujours recommandées, qu'elles soient régulières ou lors de modifications importantes.

La restructuration apportée dans les bases de données pour cette version de Canopsis nous amène à insister d'autant plus sur ce point. Il est donc fortement recommandé de réaliser une **sauvegarde complète** des VM hébergeant vos services Canopsis, avant cette mise à jour.

### Arrêt de l'environnement en cours de lancement

Vous devez prévoir une interruption du service afin de procéder à la mise à jour qui va suivre.

=== "Paquets CentOS 7"

    ```sh
    canoctl stop
    ```

=== "Docker Compose"

    ```sh
    docker-compose -f 00-data.docker-compose.yml -f 01-prov.docker-compose.yml -f 02-app.docker-compose.yml down
    ```
    
    Ou bien, si vous utilisez encore l'ancien procédé :
    
    ```sh
    docker-compose down
    ```

### Changement des ports d'écoute du reverse proxy nginx

Pour des raisons de sécurité, le reverse proxy `nginx` fourni avec Canopsis dans les images Docker n'écoute plus les ports **80** et **443**.
A présent, le port 8080 est utilisé pour l'écoute http et le 8443 pour l'écoute https.

=== "Paquets CentOS 7"

    Cette méthode d'installation n'est pas affectée par ce changement.
    
=== "Docker Compose"

Synchronisez vos fichiers de configuration yaml depuis le dépôt ou modifiez le fichier `02-app.docker-compose.yml` comme suit :

```diff
  nginx:
    image: ${DOCKER_REPOSITORY}${COMMUNITY_BASE_PATH}nginx:${CANOPSIS_IMAGE_TAG}
    ports:
-      - "80:80"
-      - "443:443"
+      - "80:8080"
+      - "443:8443"
    env_file:
      - compose.env
    restart: unless-stopped
    volumes:
      - nginxcerts:/etc/nginx/ssl
```

### Suppression des options `-enrichContext` et `-enrichInclude`

Les options `-enrichContext` et `-enrichIncludeè ont été retirées du moteur `engine-che`.

=== "Paquets CentOS 7"

    Lancez la commande suivante afin de savoir si cette option est utilisée :
    
    ```sh
    grep -lr "-enrich" /etc/systemd/system/canopsis-engine-go@engine-che.service.d/*
    ```
    
    Si cette commande affiche un résultat, éditez les fichiers qu'elle mentionne afin d'y retirer ces options.

=== "Docker Compose"

    Supprimez toute éventuelle utilisation des options `-enrichContext` et `-enrichInclude` dans vos fichiers de référence Docker Compose.

### Lancement des scripts de migration

Assurez-vous que le service MongoDB soit bien lancé et exécutez les commandes suivantes, en adaptant les identifiants MongoDB ci-dessous si nécessaire :

=== "Paquets CentOS 7"

    Sur la machine sur laquelle les paquets `canopsis*` sont installés :
    
    ```sh
    cd /opt/canopsis/share/migrations/mongodb/release4.6
    for file in $(find . -type f -name "*.js" | sort -n); do
       mongo -u cpsmongo -p canopsis canopsis < "$file"
    done
    ```

=== "Docker Compose"

    Depuis une machine qui a un client `mongo` installé et qui peut joindre le service `mongodb` d'un point de vue réseau :
    
    ```sh
    git clone --depth 1 --single-branch -b release-4.6 https://git.canopsis.net/canopsis/canopsis-community.git
    cd canopsis-community/community/go-engines-community/database/migrations
    for file in $(find release4.6 -type f -name "*.js" | sort -n); do
       mongo "mongodb://cpsmongo:canopsis@localhost:27017/canopsis" < "$file" # URI à adapter au besoin
    done
    ```
    
    Il est aussi possible de récupérer le répertoire `migrations` et de le présenter en volume dans le conteneur `mongodb` afin de réaliser le lancement du script depuis le conteneur `mongodb`.

!!! attention
    Ces scripts essaient de gérer le plus de cas d'usage possible, mais la bonne exécution de ces scripts en toute condition ne peut être garantie.

    Ils doivent obligatoirement être lancés **avant** le lancement des scripts de provisioning lors de l'étape suivante.
    
    N'hésitez pas à nous signaler tout problème d'exécution que vous pourriez rencontrer lors de cette étape.

### Synchronisation du fichier de configuration `canopsis.toml`

Vérifiez que votre fichier `canopsis.toml` soit bien à jour par rapport au fichier de référence, notamment dans le cas où vous auriez apporté des modifications locales à ce fichier :

* [`canopsis.toml` pour Canopsis Community 4.6.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/4.6.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-community.toml)
* [`canopsis.toml` pour Canopsis Pro 4.6.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/4.6.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-pro.toml)

=== "Paquets CentOS 7"

    Le fichier à synchroniser est `/opt/canopsis/etc/canopsis.toml`.

=== "Docker Compose"

    Si vous n'avez pas apporté de modification locale, ce fichier est directement intégré et mise à jour dans les conteneurs, et vous n'avez donc pas de modification à apporter.
    
    Si vous modifiez ce fichier à l'aide d'un volume surchargeant `canopsis.toml`, c'est ce fichier local qui doit être synchronisé.

### Lancement du provisioning et de `canopsis-reconfigure`

Le provisioning doit être lancé afin de mettre à jour certaines données en base, tandis que `canopsis-reconfigure` prend en compte les changements apportés au fichier `canopsis.toml`.

=== "Paquets CentOS 7"

    Lancez les scripts de provisioning :
    
    ```sh
    # si vous utilisez Canopsis Community
    su - canopsis -c "canopsinit --canopsis-edition core"
    # OU si vous utilisez Canopsis Pro
    su - canopsis -c "canopsinit --canopsis-edition cat"
    ```
    
    Puis, lancez `canopsis-reconfigure`. Attention, cette fois-ci de nouvelles options doivent lui être données :
    
    ```bash
    set -o allexport ; source /opt/canopsis/etc/go-engines-vars.conf
    /opt/canopsis/bin/canopsis-reconfigure -migrate-postgres=true -postgres-migration-mode=up -postgres-migration-directory=/opt/canopsis/share/migrations/postgres
    ```
    
=== "Docker Compose"

    Lancez à nouveau toute la partie `data` (MongoDB, RabbitMQ, Redis, PostgreSQL…) :

    ```sh
    docker-compose -f 00-data.docker-compose.yml up -d
    ```

    !!! Attention
        Si vous avez personnalisé la ligne de commande de l'outil `canopsis-reconfigure`, nous vous conseillons de supprimer cette persionnalisation.  
        L'outil est en effet pré paramétré pour fonctionner naturellement.

    Exécutez la commande suivante :
    
    ```sh
    docker-compose -f 01-prov.docker-compose.yml up -d
    ```
    
    Ou bien, si vous utilisez encore l'ancien procédé :
    
    ```sh
    docker-compose up -d provisioning reconfigure
    ```
    
### Remise en route des moteurs et des services de Canopsis

Si et seulement si les commandes précédentes n'ont pas renvoyé d'erreur, vous pouvez relancer l'intégralité des services.

=== "Paquets CentOS 7"

    Relancez la totalité de l'environnement :
    
    ```sh
    systemctl daemon-reload
    canoctl restart
    ```

=== "Docker Compose"

    Lancez maintenant la partie `02-app`, afin de bénéficier de l'application Canopsis en elle-même :
    
    ```sh
    docker-compose -f 02-app.docker-compose.yml up -d
    ```
    
    Ou bien, si vous utilisez encore l'ancien procédé :
    
    ```sh
    docker-compose up -d
    ```

### Fin de la mise à jour

Après quelques minutes, le service devrait être à nouveau accessible sur son interface web habituelle. En cas de problème, consultez l'ensemble des logs.

Suivez la [section « Après la mise à jour »](../../guide-administration/mise-a-jour/index.md#apres-la-mise-a-jour) du Guide d'administration afin d'en savoir plus.
