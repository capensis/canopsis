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

Tous les scripts peuvent utiliser une variable d’environnement `CANOPSIS_DISTRIBUTION` permettant de ne construire les images que sur un seul OS.

Pour activer cette fonctionnalité, dans votre shell exécutez :

```bash
# export CANOPSIS_DISTRIBUTION="platform-version"
export CANOPSIS_DISTRIBUTION="centos-7"
```

Si cette variable n’existe contient `all`, toutes les plateformes supportées seront alors construites.

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

## Setup - Environnement

Afin de correctement délivrer une *release*, certaines varibles d’environnement sont obligatoires :

```
# variables pour faciliter le boulot via cette doc
workdir="${HOME}/cps-release"
tag="version_a_produire"
```

**Prendre exemple sur le fichier `env-release.example.sh`** : en faire une copie que vous ne versionnerez **pas** dans canopsis contenant les paramètres désirés, puis :

```
source env-release.copie.sh
```

Il vous faudra aussi avoir un environnement Go utilisable et donc une variable `GOPATH` correctement positionnée.

Afin d’éviter de saisir votre *passphrase* ssh :

```bash
eval $(ssh-agent -s)
ssh-add
```

## Setup - Release

**Après avoir revérifié votre fichier d’environnement** et _sourcé_ ce fichier dans votre shell courant :

Lancer le script `build-release.sh` avec les variables d’environnement nécessaires.

Le script va faire le nécessaire pour ajouter les tags sur les projets définis dans `tools/catag/catag.ini`, cloner `canopsis-next` et faire tout autre travail nécessaire aux préparatifs d’une *release*.

Il vous faudra :

 * Avoir cloné le dépôt `canopsis/canopsis` dans la version qui va servir au build : ce dépôt n’est pas recloné par le script, tout comme CAT.
 * Un TOKEN d’accès à l’API GitLab de `git.canopsis.net`
 * Les droits de suppression et d’ajout de tag sur les dépôts configurés dans `catag.ini`
 * Le droit de *pull* et *push* sur les dépôts configurés dans `catag.ini`
 * Avoir installé `git`, `go` et `rsync` (présence vérifiée par le script)

```bash
cd ${workdir}

git clone ssh://git@git.canopsis.net/canopsis/canopsis -b <branch|commit>
./build-release.sh
```

## Build

Une fois les étapes `Setup` exécutées, les scripts `build-docker.sh` et `build-packages.sh` peuvent être lancés indépendemment.

Seul `build-packages.sh` requiert que `build-docker.sh` ai été exécuté préalablement.

En revanche nous allons utiliser `build-all.sh` qui va prendre cela en charge et assurer un lancement dans l’ordre.

À la moindre erreur, les scripts s’arrêteront.

### Build - Core

```bash
cd ${workdir}/canopsis
./build-all.sh
```

 * Images docker produites
 * Paquets disponibles dans le dossier `packages`

### Build - CAT

Utiliser le même script pour cat :

```bash
cd ${workdir}

git clone ssh://git@git.canopsis.net/cat/canopsis-cat.git -b ${tag} canopsis-cat
cd canopsis-cat
./build-all.sh
```

### Push

(uniquement pour debian-9)

```bash
docker push canopsis/canopsis-core:${tag}
docker push canopsis/canopsis-cat:${tag}
docker push canopsis/canopsis-prov:${tag}
docker push canopsis/canopsis-cat-prov:${tag}
```

### Nettoyage

Ces actions ne sont pas nécessaires si vous allez relancer des *builds* :

```bash
rm -rf ${workdir}
ssh-agent -k
```

## Build des paquets

Les scripts `build-all.sh` exécutent automatiquement la construction des paquets. Cependant s’il est nécessaire de le faire à la main :

Se mettre à la racine du dépôt CAT puis :

```bash
./build-packages.sh
```

Les paquets sont alors disponibles dans le dossier `packages`.

### RC release 

Dans le cas où la version utilisée en `CANOPSIS_TAG` n’est pas compatible avec le système de *build* des paquets, deux variables sont à disposition et utilisables pour toutes les plateformes :

```bash
CANOPSIS_PACKAGE_TAG=2.x.x CANOPSIS_PACKAGE_REL=2 ./build-packages.sh ${tag}

CANOPSIS_PACKAGE_TAG=2.x.x CANOPSIS_PACKAGE_REL=rc1 ./build-packages.sh ${tag}
```

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
