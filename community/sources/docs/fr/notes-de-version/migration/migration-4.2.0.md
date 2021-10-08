# Guide de migration vers Canopsis 4.2.0

Ce guide donne des instructions vous permettant de mettre à jour un Canopsis 4.1.0 ou 4.1.1 vers la version 4.2.0.

!!! note
    **18/05/2021 :** Mise à jour des [scripts de migration](#migration-des-actions-webhooks-et-idlerules-existants).

## Prérequis

L'ensemble de cette procédure doit être lu avant son exécution.

Cette procédure ne doit être suivie que dans le cadre d'une mise à jour. Elle n'est pas utile dans le cadre d'une nouvelle installation.

Vous devez avoir suivi les notes de version et appliqué tous les guides de mise à jour de toutes les versions comprises entre votre version actuelle de Canopsis et Canopsis 4.2.0.

### Fin de la production de paquets pour Debian 9

Comme annoncé depuis plusieurs mois, Canopsis ne produit et n'héberge désormais plus aucun paquet pour Debian 9 (« *stretch* »).

Les deux seules méthodes d'installation officiellement prises en charge (sauf disposition particulière) sont donc :

*  les paquets RPM pour CentOS 7 ;
*  les conteneurs Docker au travers de Docker Compose.

### Note auprès des utilisateurs de Dockerhub

Canopsis 4.4.0 migrera ses images Docker de Dockerhub vers le nouveau registre `registry.canopsis.net`.

Si vous bénéficiez de Canopsis Pro, et que vous utilisez Canopsis au travers des images Docker, vous avez dû recevoir une communication vous invitant à nous transmettre une adresse e-mail à utiliser afin de mettre en place ce nouvel accès. Si ce n'est pas le cas, veuillez vous rapprocher de votre contact habituel chez Capensis.

Les détails de cette migration seront publiés après la sortie de Canopsis 4.3.0. Canopsis reste encore disponible sur Dockerhub sans changement particulier pour le moment.

### Suppression du moteur `engine-stat`

Le moteur `engine-stat`, désactivé et déprécié depuis [Canopsis 3.31.0](../3.31.0.md), a été complètement supprimé.

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

### Déclenchement des webhooks par le moteur `engine-action` (Pro)

À partir de Canopsis 4.2.0, les webhooks (disponibles uniquement dans l'édition Pro) sont déclenchés par le moteur `engine-action`.

=== "Paquets CentOS 7"

    Configurez l'unité systemd `canopsis-engine-go@engine-action.service` afin qu'elle comporte dorénavant l'option `-withWebhook` :

    ```sh
    mkdir -p /etc/systemd/system/canopsis-engine-go@engine-action.service.d
    cat > /etc/systemd/system/canopsis-engine-go@engine-action.service.d/action.conf << EOF
    [Service]
    ExecStart=
    ExecStart=/usr/bin/env /opt/canopsis/bin/%i -withWebhook
    EOF
    systemctl daemon-reload
    ```

=== "Docker Compose"

    Dans votre fichier `docker-compose.yml`, modifiez la ligne d'exécution du conteneur `engine-action` afin qu'elle contienne l'option `-withWebhook` :

    ```yaml
    command: /engine-action -withWebhook
    ```

Notez que le moteur `engine-webhook` est cependant toujours nécessaire et ne doit donc pas être désactivé.

### Migration des actions, webhooks et idlerules existants

Toutes vos configurations d'actions, webhooks ou *idle rules* doivent être migrées à l'aide de scripts, afin d'être fonctionnels avec les nouveaux scénarios de Canopsis 4.2.

Sur une machine disposant d'un accès à `git.canopsis.net` ainsi que d'un client MongoDB, assurez-vous que le service MongoDB soit bien lancé et exécutez les scripts suivants, en adaptant les identifiants MongoDB ci-dessous si nécessaire :

```sh
for file in 000_migrate_actions.js 001_migrate_webhook.js 002_migrate_idlerules.js ; do
    curl -O -L "https://git.canopsis.net/canopsis/canopsis-community/-/raw/release-4.3/community/go-engines-community/database/migrations/release4.2/$file"
    mongo -u cpsmongo -p canopsis canopsis < "$file"
done
```

Notez au passage que les nouveaux scénarios **ne disposent plus d'_event patterns_**. Si vous aviez écrit des règles d'actions ou de webhooks comportant des *event patterns*, vous devez les transposer vers des *alarm patterns* ou des *entity patterns*, selon votre cas d'utilisation.

### Ajout d'une nouvelle option de configuration pour la remédiation (Pro)

Si vous utilisez la fonctionnalité de remédiation, une nouvelle option a été ajoutée au fichier `canopsis.toml` :

```toml
[Canopsis.remediation]
JobExecutorFetchTimeoutSeconds = 30
```

Vérifiez que vous disposez bien de cette option, et ajustez-la si nécessaire.

### Synchronisation du fichier de configuration `canopsis.toml`

Vérifiez que votre fichier `canopsis.toml` soit bien à jour par rapport au fichier de référence, notamment dans le cas où vous auriez apporté des modifications locales à ce fichier :

* [`canopsis.toml` pour Canopsis Community 4.2.0](https://git.canopsis.net/canopsis/go-engines/-/blob/4.2.0/cmd/canopsis-reconfigure/canopsis-core.toml.example)
* [`canopsis.toml` pour Canopsis Pro 4.2.0](https://git.canopsis.net/canopsis/go-engines/-/blob/4.2.0/cmd/canopsis-reconfigure/canopsis-cat.toml.example)

### Fin de la mise à jour

Une fois ces changements apportés, suivez la [procédure standard de mise à jour de Canopsis](../../guide-administration/mise-a-jour/index.md) et redémarrez l'environnement.
