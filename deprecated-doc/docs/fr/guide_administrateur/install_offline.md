# Installation sans réseau

Il vous sera nécessaire de télécharger préalablement tous les dépôts Canopsis :

 * canopsis
 * canopsis-externals
 * canopsis-webcore
 * ... 

## Système

Selon la distribution, des dépendances sont nécessaires. Voir les fichiers `sources/extra/dependencies/<distro>_<version>`.

## Python

Pour installer les paquets Canopsis sans avoir besoin du réseau, il est nécessaire de passer des options à `pip` et de créer un fichier de configuration :

```
pip --no-index -f file:///chemin/vers/les/sources/des/paquets/python/ <cmd>
```

```
cat > ~/.pydistutils.cfg << EOF
[easy_install]
allow_hosts = ''
find_links = file:///chemin/vers/les/sources/des/paquets/python
EOF
```

Le chemin des paquets python est en réalité le chemin vers le dépôts `canopsis-externals`, dans le dossier `python-libs`.
