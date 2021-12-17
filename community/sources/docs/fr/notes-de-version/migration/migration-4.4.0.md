# Guide de migration vers Canopsis 4.4.0

Ce guide donne des instructions vous permettant de mettre à jour Canopsis 4.3 vers [la version 4.4.0](../4.4.0.md).

## Prérequis

L'ensemble de cette procédure doit être lu avant son exécution.

Ce document ne prend en compte que Canopsis Community et Canopsis Pro : tout développement personnalisé dont vous pourriez bénéficier ne fait pas partie du cadre de ce Guide de migration.

## Docker : migration de Dockerhub vers `docker.canopsis.net`

Les conteneurs Docker de Canopsis étaient jusqu'à présents hébergés sur [Dockerhub](https://hub.docker.com/u/canopsis/).

À partir de Canopsis 4.4.0, les images doivent être récupérées depuis le registre `docker.canopsis.net`, réparti entre une section Community (publique) et Pro (privée).

Si vous bénéficiez d'une souscription à Canopsis Pro, un compte de service (lié au service Gitlab [git.canopsis.net](https://git.canopsis.net)) a normalement déjà été mis en place afin que vous puissiez récupérer ces images. Sinon, rapprochez-vous de votre contact habituel.

Les nouveaux fichiers de référence Docker Compose (décrits plus bas) comportent de nouvelles variables permettant de gérer cette nouvelle organisation :

```yaml
axe:
  image: ${DOCKER_REPOSITORY}${COMMUNITY_BASE_PATH}engine-axe:${CANOPSIS_IMAGE_TAG}
```

Pour récupérer ces images, assurez-vous que les flux réseau HTTPS (TCP/443) vers `docker.canopsis.net` soient bien autorisés, puis connectez-vous avec votre compte de service Gitlab avec la commande suivante :

```sh
docker login docker.canopsis.net
```

Vous devriez ensuite pouvoir récupérer les conteneurs Canopsis avec [les commandes habituelles](../../guide-administration/mise-a-jour/index.md#mise-a-jour-en-environnement-docker-compose).

!!! note "Notes"
    * Les conteneurs des services tiers (tels que MongoDB, RabbitMQ, Redis…) restent hébergés sur Dockerhub.
    * Les versions de Canopsis antérieures à Canopsis 4.4.0 ont aussi été importées dans le registre `docker.canopsis.net`.
    * À partir de Canopsis 4.5.0, plus aucune nouvelle version de Canopsis ne sera publiée sur Dockerhub. À terme, le projet Canopsis ne sera plus du tout disponible sur Dockerhub.
    * Les accès à `repositories.canopsis.net` (pour les paquets RPM) ne sont pas concernés pas ces modifications et leur utilisation reste donc inchangée.

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

### Suppression de composants obsolètes

Les moteurs `engine-heartbeat` ainsi que les composants `scheduler` et `importctx` ne sont plus pris en charge à partir de cette version. Ils doivent donc être supprimés.

=== "Paquets CentOS 7"

    Exécutez la commande suivante pour désactiver ces services :

    ```sh
    systemctl list-units -a --type=service --plain --no-legend "canopsis*" | awk '/heartbeat/ || /scheduler/ || /importctx/ { print $1 }' | xargs -r systemctl disable
    ```

=== "Docker Compose"

    Supprimez toute éventuelle section `heartbeat` `scheduler` ou `task_importctx` de votre fichier Docker Compose. Toute référence à un éventuel volume `importctxdata` peut aussi être supprimée.

### Passage de `external-job-executor` à `engine-remediation` (Pro)

Si vous utilisez Canopsis Pro, la fonctionnalité de remédiation n'est plus exécutée par le service `external-job-executor` mais par le nouveau moteur `engine-remediation`.

Le moteur `engine-axe` doit aussi être lancé avec l'option `-withRemediation=true` pour bénéficier de cette fonctionnalité.

=== "Paquets CentOS 7"

    Exécutez les commandes suivantes :

    ```sh
    systemctl disable canopsis-service@external-job-executor.service
    systemctl enable canopsis-engine-go@engine-remediation.service

    mkdir -p /etc/systemd/system/canopsis-engine-go@engine-axe.service.d
    cat > /etc/systemd/system/canopsis-engine-go@engine-axe.service.d/axe.conf << EOF
    [Service]
    ExecStart=
    ExecStart=/usr/bin/env /opt/canopsis/bin/%i -publishQueue Engine_correlation -withRemediation=true
    EOF
    systemctl daemon-reload
    ```

=== "Docker Compose"

    Dans votre Docker Compose, supprimez toute section `external-job-executor` similaire à la suivante :

    ```yaml
    external-job-executor:
      image: ${DOCKER_REPOSITORY}${PRO_BASE_PATH}external-job-executor:${CANOPSIS_IMAGE_TAG}
      env_file:
        - compose.env
      restart: unless-stopped
      command: /external-job-executor
    ```

    et remplacez-la par la nouvelle section `remediation` suivante :

    ```yaml
    remediation:
      image: ${DOCKER_REPOSITORY}${PRO_BASE_PATH}engine-remediation:${CANOPSIS_IMAGE_TAG}
      env_file:
        - compose.env
      restart: unless-stopped
      command: /engine-remediation
    ```

    Puis, assurez-vous que le moteur `engine-axe` soit bien lancé avec l'option `-withRemediation=true` :

    ```yaml hl_lines="3"
    axe:
      # ...
      command: /engine-axe -publishQueue Engine_correlation -withRemediation=true
    ```

### Docker Compose : réorganisation des environnements de référence

Les fichiers de référence Docker Compose ont été complétement revus dans cette version :

<https://git.canopsis.net/canopsis/canopsis-pro/-/tree/release-4.4/pro/deployment/canopsis/docker>

L'environnement a notamment été découpé en 3 parties à lancer successivement (`00-data` pour les données persistantes et briques externes ; `01-prov` pour le provisioning ; `02-app` pour l'application Canopsis en elle-même). De nouvelles variables, telles que `DOCKER_REPOSITORY`, ou `CPS_SERVER_NAME` ont aussi été introduites. Certains conteneurs ont aussi été renommés ou déplacés entre Community et Pro.

Il vous est donc recommandé de partir de ce nouveau référentiel, et d'y appliquer toute modification locale que vous pouviez faire jusqu'à présent. En cas de nécessité, rapprochez-vous de votre contact habituel pour un accompagnement.

### Synchronisation du fichier de configuration `canopsis.toml`

Vérifiez que votre fichier `canopsis.toml` soit bien à jour par rapport au fichier de référence, notamment dans le cas où vous auriez apporté des modifications locales à ce fichier :

* [`canopsis.toml` pour Canopsis Community 4.4.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/4.4.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-community.toml)
* [`canopsis.toml` pour Canopsis Pro 4.4.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/4.4.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-pro.toml)

=== "Paquets CentOS 7"

    Le fichier à synchroniser est `/opt/canopsis/etc/canopsis.toml`.

=== "Docker Compose"

    Si vous n'avez pas apporté de modification locale, ce fichier est directement intégré et mise à jour dans les conteneurs.

    Si vous le surchargez à l'aide d'un volume pour y apporter des modifications, c'est ce fichier local qui doit être synchronisé.

## Fin de la mise à jour

Une fois ces changements apportés, suivez la [procédure standard de mise à jour de Canopsis](../../guide-administration/mise-a-jour/index.md) et redémarrez l'environnement.

Vous devez ensuite contrôler la bonne mise à jour de la configuration Nginx

### Migrations

Sur une machine disposant d'un accès à `git.canopsis.net` ainsi que d'un client MongoDB, assurez-vous que le service MongoDB soit bien lancé et exécutez les commandes suivantes, en adaptant les identifiants MongoDB ci-dessous si nécessaire :

```sh
git clone --depth 1 --single-branch -b release-4.4 https://git.canopsis.net/canopsis/canopsis-community.git
cd canopsis-community/community/go-engines-community/database/migrations
for file in $(find release4.4 -type f -name "*.js" | sort -n); do
   mongo -u cpsmongo -p canopsis canopsis < "$file"
done
```

!!! attention
    Ces scripts essaient de gérer le plus de cas d'usage possible, mais la bonne exécution de ces scripts en toute condition ne peut être garantie.

    N'hésitez pas à nous signaler tout problème d'exécution que vous pourriez rencontrer lors de cette étape.

### Mise à jour de la configuration de Nginx

!!! information
    Canopsis 4.4.0 propose maintenant une configuration HTTPS, non activée par défaut. Consultez le [Guide d'activation de HTTPS](../../guide-administration/administration-avancee/configuration-composants/reverse-proxy-nginx-https.md) pour en savoir plus.

Plusieurs changements ont été apportés à la configuration de Nginx.

=== "Paquets CentOS 7"

    ```sh
    cp -f /opt/canopsis/deploy-ansible/playbook/roles/canopsis/templates/nginx/cors.j2 /etc/nginx/cors.inc
    cp -f /opt/canopsis/deploy-ansible/playbook/roles/canopsis/templates/nginx/https.j2 /etc/nginx/https.inc
    sed -e 's,{{ CPS_API_URL }},http://127.0.0.1:8082,g' -e 's,{{ CPS_SERVER_NAME }},"localhost",g' /opt/canopsis/deploy-ansible/playbook/roles/canopsis/templates/nginx/default.j2 > /etc/nginx/conf.d/default.conf

    systemctl restart nginx
    ```

    !!! attention
        Si vous accédez à l'interface web de Canopsis au travers d'un nom de domaine (par exemple `canopsis.mon-si.fr`), vous devrez **obligatoirement** configurer la ligne `set $canopsis_server_name` du fichier `/etc/nginx/conf.d/default.conf` avec cette valeur.

=== "Docker Compose"

    Si vous n'avez pas surchargé la configuration Nginx à l'aide d'un volume, vous n'avez rien à faire.

    En revanche, si vous mainteniez vos propres versions modifiées de ces fichiers de configuration, vous devez manuellement vous synchroniser avec la totalité des modifications ayant été apportées dans `/etc/nginx/`.

    !!! attention
        Si vous accédez à l'interface web de Canopsis au travers d'un nom de domaine (par exemple `canopsis.mon-si.fr`), vous devrez **obligatoirement** configurer la ligne `CPS_SERVER_NAME` du fichier `compose.env` associé à votre Docker Compose avec cette valeur :

        ```ini
        CPS_SERVER_NAME=canopsis.mon-si.fr
        ```
