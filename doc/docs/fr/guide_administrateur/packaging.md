# Packaging

La construction de paquets est supportée pour les distributions suivantes :

 * `centos-7` : CentOS 7 x86_64
 * `debian-8` : Debian 8 Jessie amd64
 * `debian-9` : Debian 9 Stretch amd64

La construction de ces paquets se repose sur Docker. Vous trouverez dans cette documentation le nécessaire pour effectuer un build Docker, nécessaire à la construction des paquets.

Pour le moment ces scripts ne sont disponibles que dans CAT et ne construisent que des paquets CAT.

Le travail pour avoir un paquet `core` et un paquet `cat` sera fait plus tard.

## Présentation globale

Les images Docker sont construites pour chaque plateforme supportée avec un tag `plateforme_id-id-tag`. Exemple : `centos-7-2.5.11`.

L’image Docker "officielle" se repose sur Debian 9 et sera taggée simplement avec le numéro de version.

Tous les scripts peuvent utiliser une variable d’environnement `SYSBASE` permettant de ne construire les images que sur un seul OS.

Pour activer cette fonctionnalité, dans votre shell exécutez :

```bash
# export SYSBASE="platform-version"
export SYSBASE="centos-7"
```

Si cette variable n’existe pas ou est vide, toutes les plateformes supportées seront alors construites.

## Build Docker

*Cette documentation regroupe la construction de `core` et `cat`.*

Les images Docker Canopsis se font en deux temps :

 * Une image *core* : `canopsis/canopsis-core:<tag>`
 * Une image *CAT* : `canopsis/canopsis-cat:<tag>`

Ces deux images embarquent à la fois les *engines* et le *webserver*.

Les variables `ENGINE_NAME` et `ENGINE_MODULE` doivent être positionnées pour chaque instance d’*engine* et la variable `CPS_WEBSERVER` doit valoir `0` pour démarrer un *engine*.

La variable d’environnement `CPS_WEBSERVER` doit être mise à `1` pour démarrer le *webserver*.

### Setup - Docker

 * Installer `docker-ce`
 * Ajouter votre utilisateur régulier au groupe `docker` si votre distribution supporte ce mode.
 * Installer `git`, les clients `openssh` et `ssh-agent`.

Versions de Docker testées : version Canopsis :

 * `17.xx.x-ce / canopsis <= 2.5.6`
 * `18.02.0-ce / canopsis >= 2.5.7`

### Setup - Environnement

Nous vous proposons ici un certain nombre de commandes destinées à vous faciliter la tâche. Conservez votre instance de shell tout au long du processus.

```bash
start_branch="<START_BRANCH>"
tag="<TAG_CANOPSIS>"
workdir="${HOME}/cps-docker"

mkdir -p ${workdir} && cd ${workdir}
```

 * `<START_BRANCH>` doit correspondre à la branche du projet `canopsis/canopsis` que vous allez utiliser pour fabriquer les images. Cette valeur doit être égale à celle présente dans `catag.ini` pour ce projet.
 * `<TAG_CANOPSIS>` la version à publier.
 * Vous pouvez changer le `workdir` à votre convenance. Ce doit être un dossier dans lequel vous pourrez écrire. Aucune permission particulière n’est nécessaire.

Afin d’éviter de saisir votre *passphrase* ssh :

```bash
eval $(ssh-agent -s)
ssh-add
```

### Setup - Tag sur les dépôts

Vous devez cloner le dépôt canopsis dans la branche que vous aller faire tagger par *catag*.

```bash
cd ${workdir}

git clone ssh://git@git.canopsis.net/canopsis/canopsis.git -b ${start_branch} canopsis
cd canopsis

git submodule update --init

cd tools/catag

go get github.com/vaughan0/go-ini
go build .

./catag -token "<token gitlab>" -tag ${tag}
```

### Build - Core

```bash
cd ${workdir}/canopsis
./build-docker.sh ${tag} ${start_branch}
```

### Build - CAT

Utiliser le même script pour cat :

```bash
cd ${workdir}

git clone ssh://git@git.canopsis.net/cat/canopsis-cat.git -b ${tag} canopsis-cat
cd canopsis-cat
./build-docker.sh ${tag} ${start_branch}
```

### Push

```bash
docker push canopsis/canopsis-core:${tag}
docker push canopsis/canopsis-cat:${tag}
```

### Nettoyage

Ces actions ne sont pas nécessaires si vous allez relancer des *builds* :

```bash
rm -rf ${workdir}
ssh-agent -k
```

## Build des paquets

Se mettre à la racine du dépôt CAT puis :

```
./build-packages.sh
```

Les paquets sont alors disponibles dans le dossier `packages`.

## Installation

### CentOS / RedHat 7

```
yum localinstall canopsis-cat-<version>-1.el7.centos.x86_64.rpm
```

### Debian 8 / 9

```
dpkg -i canopsis-cat-1-<version>.amd64.<platform>.deb
apt install -f
```

## Init

Des unités `systemd` sont disponibles :

 * `canopsis-engine@<module>-<name>.service`
 * `canopsis-webserver.service`

Voici tous les engines qui vous pouvez activer dans `core` et `cat` :

```bash
systemctl enable canopsis-engine@acknowledgement-acknowledgement.service
systemctl enable canopsis-engine@dynamic-alerts.service
systemctl enable canopsis-engine@cancel-cancel.service
systemctl enable canopsis-engine@cleaner-cleaner_alerts.service
systemctl enable canopsis-engine@cleaner-cleaner_events.service
systemctl enable canopsis-engine@dynamic-context-graph.service
systemctl enable canopsis-engine@eventduration-eventduration.service
systemctl enable canopsis-engine@event_filter-event_filter.service
systemctl enable canopsis-engine@eventstore-eventstore.service
systemctl enable canopsis-engine@linklist-linklist.service
systemctl enable canopsis-engine@dynamic-pbehavior.service
systemctl enable canopsis-engine@dynamic-perfdata.service
systemctl enable canopsis-engine@scheduler-scheduler.service
systemctl enable canopsis-engine@selector-selector.service
systemctl enable canopsis-engine@dynamic-serie.service
systemctl enable canopsis-engine@dynamic-stats.service
systemctl enable canopsis-engine@task_dataclean-task_dataclean.service
systemctl enable canopsis-engine@task_importctx-task_importctx.service
systemctl enable canopsis-engine@task_linklist-task_linklist.service
systemctl enable canopsis-engine@task_mail-task_mail.service
systemctl enable canopsis-engine@ticket-ticket.service
systemctl enable canopsis-engine@dynamic-watcher.service

systemctl enable canopsis-webserver.service
```

Le fichier `/opt/canopsis/etc/amqp2engines.conf` est toujours en vigeur.

### Nombre de process

Pour le moment le nombre de processus lancés via `engine-launcher` est fixé dans les unités.

Pour changer le nombre d’instances :

```bash
mkdir -p /etc/systemd/system/canopsis-engine@<module>-<name>
cat > /etc/systemd/system/canopsis-engine@<module>-<name>/workers.conf << EOF
[Service]
Environment=WORKERS=X
EOF
```

Remplacer `X` par le nombre de workers désiré. Par défaut `1`.

```bash
systemctl daemon-reload
systemctl restart canopsis-engine@<module>-<name>.service
```
