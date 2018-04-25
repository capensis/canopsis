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

## Python Wheels

Le script de build docker que vous allez utiliser créer et démarre une image qui va créer des `wheels` python et les mettre en cache. Cela permet de considérablement accélérer la construction des images lorsqu’on fait une nouvelle release.

Les wheels sont reconstruites dans les cas suivants :

 * Dépôt inexistant pour la plateforme en cours de construction
 * Changement du contenu du fichier `requirements.txt` de canopsis

L’emplacement du cache des wheels peut être modifié :

```bash
export WHEEL_DIR="/tmp/canopsis-wheelrep"
```

Par défaut les wheels seront entreposée dans `docker/wheelbuild` lors de la construction, puis copiées dans `docker/wheels`.

## Docker

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

 * `<START_BRANCH>` doit correspondre à la branche du projet `canopsis/canopsis` que vous allez utiliser pour fabriquer les images. Cette valeur doit être égale à celle présente dans `catag.ini` pour ce projet (la branche doit être présente dans tous les dépots annexes, comme les différentes bricks).
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

(uniquement pour debian-9)

```bash
docker push canopsis/canopsis-core:${tag}
docker push canopsis/canopsis-cat:${tag}
docker push canopsislcanopsis-prov:{tag}
docker push canopsislcanopsis-cat-prov:{tag}
```

### Nettoyage

Ces actions ne sont pas nécessaires si vous allez relancer des *builds* :

```bash
rm -rf ${workdir}
ssh-agent -k
```

## Build des paquets

Se mettre à la racine du dépôt CAT puis :

```bash
./build-packages.sh ${tag}
```

Les paquets sont alors disponibles dans le dossier `packages`.

## Erreurs connues

### Failed to create image

```
failed to export image: failed to create image: failed to get layer sha256:51a946666f22f58babd6e3e642b9db0f262f761f7081997a1c2e71bcddcdf5d3: layer does not exist
```

Problème de synchronisation sur disque des données. Probablement BTRFS en système de fichier ?

Dans tous les cas :

```
sudo sync
```

Relancer le build, le cache docker sera disponible.
