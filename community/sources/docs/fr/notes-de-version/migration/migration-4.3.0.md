# Guide de migration vers Canopsis 4.3.0

Ce guide donne des instructions vous permettant de mettre à jour Canopsis 4.2 vers [la version 4.3.0](../4.3.0.md).

## Prérequis

L'ensemble de cette procédure doit être lu avant son exécution.

Cette procédure ne doit être suivie que dans le cadre d'une mise à jour. Elle n'est pas utile dans le cadre d'une nouvelle installation.

Vous devez avoir suivi les notes de version et appliqué tous les guides de mise à jour de toutes les versions comprises entre votre version actuelle de Canopsis et Canopsis 4.3.0.

### Plugin `datasource` indisponible dans cette version

Comme précisé dans les [notes de version de Canopsis 4.3.0](../4.3.0.md), le plugin `datasource` du moteur `engine-che` n'est pas fonctionnel avec Canopsis 4.3.0. Il est conseillé de ne pas effectuer de mise à jour pour le moment si vous dépendez de ce plugin.

Note : la majorité des utilisateurs n'utilisent pas ce plugin et ne sont donc pas concernés.

### Synchronisation obligatoire avec l'APIv4 pour tout script tiers

Si vous aviez développé ou fait développer des scripts tiers s'interfaçant avec les API de Canopsis, ceux-ci doivent obligatoirement être migrés vers l'APIv4 et ses nouveaux paradigmes, afin de rester fonctionnels.

Voyez les [détails apportés dans les notes de version](../4.3.0.md#migration-apiv2-vers-apiv4) pour en savoir plus à ce sujet.

## Procédure de mise à jour

### Arrêt de l'environnement en cours de lancement

Vous devez prévoir une interruption du service afin de procéder à la mise à jour qui va suivre.

=== "Paquets CentOS 7"

    ```sh
    canoctl stop
    ```

=== "Docker Compose"

    ```sh
    docker-compose down
    ```

### Ajustement possible de la configuration SELinux (paquets RPM)

=== "Paquets CentOS 7"

    Auparavant, dans le cadre d'une installation avec des paquets RPM, Canopsis nécessitait de complètement désactiver SELinux, SELinux n'étant pas pris en charge.

    À partir de Canopsis 4.3.0, ce prérequis est moins radical, en ce sens qu'une configuration de SELinux en mode *permissif* est aussi permise.

    Pour cela, si aucun autre applicatif ne vous l'empêche, vous pouvez modifier le fichier `/etc/selinux/config` de votre nœud Canopsis de la sorte :

    ```diff
    -SELINUX=disabled
    +SELINUX=permissive
    ```

    Puis, redémarrez votre système.

=== "Docker Compose"

    Cette partie ne s'applique pas aux environnements Docker Compose.

### Mise à jour de la configuration de Nginx

La configuration de Nginx a été revue afin de corriger des problèmes d'accès à certaines pages ou requêtes d'API pouvant contenir des caractères spéciaux.

Le fichier de configuration `/etc/nginx/conf.d/default.conf` doit être modifié de la sorte :

```diff
...

-location ~ ^/backend/(?<api_uri>(.*)) {
-        include /etc/nginx/cors.inc;
-        proxy_pass $canopsis_api_url/$api_uri$is_args$args;
+location ~ ^/backend(/.*) {
+        set $api_uri $1;
+        include /etc/nginx/cors.inc;
+        proxy_pass $canopsis_api_url$api_uri$is_args$args;

...
```

=== "Paquets CentOS 7"

    La mise à jour de ce fichier n'est pas encore automatisée lors d'une mise à jour de paquets. La modification doit obligatoirement être apportée à la main.

    Appliquez ces changements au fichier `/etc/nginx/conf.d/default.conf`, validez la modification avec `nginx -t`, puis rechargez la configuration (`systemctl reload nginx`).

=== "Docker Compose"

    Si vous n'avez pas surchargé la configuration de Nginx avec un volume de configuration personnalisé, vous n'avez rien à faire : la configuration sera automatiquement mise à jour avec Canopsis.

    Si vous avez surchargé ce fichier de configuration (à l'aide d'une ligne `volumes:` dans votre section `nginx:` de Docker Compose), vous devez ajuster votre fichier de surcharge en fonction des changements décrits ci-dessus.

### Passage d'`engine-watcher` à `engine-service`

Le moteur `engine-watcher` est remplacé par `engine-service`. Cette refonte lui ajoute notamment une capacité de multi-instanciation.

=== "Paquets CentOS 7"

    ```sh
    systemctl disable canopsis-engine-go@engine-watcher
    test -d /etc/systemd/system/canopsis-engine-go@engine-watcher.service.d && mv /etc/systemd/system/canopsis-engine-go@engine-{watcher,service}.service.d
    systemctl daemon-reload
    systemctl enable canopsis-engine-go@engine-service
    ```

=== "Docker Compose"

    Dans votre fichier `docker-compose.yml`, remplacez toutes les occurences de `watcher` par `service`.

    À titre d'exemple :

    ```diff
    -  watcher:
    -    image: engine-watcher:${CANOPSIS_IMAGE_TAG}
    +  service:
    +    image: engine-service:${CANOPSIS_IMAGE_TAG}
    ...
    ...
    -    command: /engine-watcher -publishQueue Engine_dynamic_infos
    +    command: /engine-service -publishQueue Engine_dynamic_infos
    ```

### Synchronisation du fichier de configuration `canopsis.toml`

Vérifiez que votre fichier `canopsis.toml` soit bien à jour par rapport au fichier de référence, notamment dans le cas où vous auriez apporté des modifications locales à ce fichier :

* [`canopsis.toml` pour Canopsis Community 4.3.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/4.3.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-core.toml.example)
* [`canopsis.toml` pour Canopsis Pro 4.3.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/4.3.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-cat.toml.example)

=== "Paquets CentOS 7"

    Le fichier à synchroniser est `/opt/canopsis/etc/canopsis.toml`.

=== "Docker Compose"

    Si vous n'avez pas apporté de modification locale, ce fichier est directement intégré et mise à jour dans les conteneurs.

    Si vous le surchargez à l'aide d'un volume pour y apporter des modifications, c'est ce fichier local qui doit être synchronisé.

### Migrations

Canopsis 4.3.0 apportant une importante refonte des API, plusieurs scripts de migration ont été écrits afin d'adapter vos données existantes en base de données vers les nouveaux formats attendus.

Sur une machine disposant d'un accès à `git.canopsis.net` ainsi que d'un client MongoDB, assurez-vous que le service MongoDB soit bien lancé et exécutez les commandes suivantes, en adaptant les identifiants MongoDB ci-dessous si nécessaire :

```sh
git clone --depth 1 --single-branch -b release-4.3 https://git.canopsis.net/canopsis/canopsis-community.git
cd canopsis-community/community/go-engines-community/database/migrations
for file in $(find release4.3 -type f \( -name "*.js" \) | sort); do     
   mongo -u cpsmongo -p canopsis canopsis < "$file"
done
```

### Fin de la mise à jour

Une fois ces changements apportés, suivez la [procédure standard de mise à jour de Canopsis](../../guide-administration/mise-a-jour/index.md) et redémarrez l'environnement.
