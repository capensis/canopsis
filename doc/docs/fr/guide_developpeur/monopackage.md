# Monopackage Python Canopsis

Depuis le commit canopsis/canopsis@fd4f67ae25a9e0f60073c81e95d15e101eb395a2 il est possible d’installer la partie Python de canopsis grâce à un paquet unique.

Cette étape singulière et absolument révolutionnaire dans l’approche de la distribution de Canopsis, en vue d’un monde meilleur et de réguler la souffrance sur Terre, nécessite quelques adaptations bégnines selon que vous migrez depuis un **package** éclaté, ou que vous vouliez installer ce monopackage dans un `virtualenv` python.

Veuillez trouver les incantations nécessaires à l’apparition de la paix sur Terre dans la suite de ce document.

## Virtualenv

Installez les paquets de dépendances que vous trouverez dans `sources/extra/dependencies/<distro>`. Étant donné que seules Debian et CentOS sont supportées pour le moment, il faudra adapter pour les autres distributions, mais le principe est le même : vous avez besoin de certains paquets, notemment de développement, pour pouvoir compiler les dépendances Python nécessitant certaines lib C/C++.

Une fois que c’est fait, exécutez ces commandes (en adaptant les chemins) :

```
virtualenv -p /usr/bin/python2 ~/venv-canopsis
source ~/venv-canopsis/bin/activate
pip install --no-index -f file:///${HOME}/canopsis/sources/externals/python-libs/ ~/canopsis/sources/canopsis/
```

Une fois fait, vous pouvez utiliser votre `venv` dans n’importe quel IDE/éditeur/OS-se-prenant-pour-un-éditeur afin de bénéficier (wait for it…) de la complétion !

## Upgrade vers le monopackage - build-install

```
# su - canopsis
$ hypcontrol stop
$ rm -rf ~/lib/python2.7/site-packages/canopsis.*
$ rm -rf ~/var/lib/pkgmgr/packages/python*
$ exit
# cd /path/to/canopsis/sources
# ./build-install.sh
```

## Upgrade vers le monopackage - Ansible

La partie Ansible n’est pas bien testée. La documentation suivra une fois ceci corrigé. Néanmoins une branche `develop-monopackage` est disponible pour le rôle `ansible-role-canopsis-backend`.

## Upgrade vers le monopackage - Docker

Les `Dockerfile` sont disponibles (privé).

## Installation de modifications

```
pip install --upgrade --no-index -f file:///path/to/canopsis/sources/externals/python-libs /path/to/canopsis/sources/canopsis/
```

Ou plus rapide, quand il n’y a pas besoin de MAJ les dépendances :

```
cd /path/to/canopsis/sources/canopsis
pip install --upgrade .
```

## CAT

Le même principe a été appliqué pour CAT. Le package est nommé `canopsis_cat` et toutes les références à une lib codée dans CAT doivent être mises à jour pour utiliser `canopsis_cat`. Exemple avec SNMP :

```python
# Ancien
from canopsis.snmp import mod

# Nouveau
from canopsis_cat.snmp import mod
```
