# Compilation manuelle du module connector-centreon-engine

Cette procédure décrit les étapes nécessaires à la compilation de connector-centreon-engine sur CentOS 6 et CentOS 7.

Elle doit donc être réalisée sur un environnement de compilation identique au système cible, que ce soit pour le système (CentOS 6, CentOS 7…) ou la version de Centreon Broker.

> **Note :** il n'est pas recommandé d'appliquer cette procédure de compilation *directement* sur le système cible, car cela aura tendance à « polluer » le système hôte avec de nombreuses dépendances.

## Préparation de l'environnement de compilation

*  **Pour CentOS 6 :** une CES (depuis [download.centreon.com](https://download.centreon.com)) à base CentOS 6.9 est indispensable (une CentOS officielle ne suffira pas).
*  **Pour CentOS 7 :** une CentOS ≥ 7.4 officielle (depuis [www.centos.org](https://www.centos.org)) est suffisante.

Installer les dépendances suivantes, nécessaires pour les compilations qui vont suivre :
```shell
sudo yum install bzip2 autoconf automake binutils bison gcc gcc-c++ gettext libtool make patch cmake qt-devel gnutls gnutls-devel rrdtool rrdtool-devel git redhat-lsb-core
```

## Installation de Boost

> **Information :** il est généralement préférable d'embarquer une version statique de Boost, afin de faciliter la maintenance. La compatibilité Boost a en effet tendance à se casser régulièrement lors des mises à jour de paquets, si l'on prend la version du système.

On compile donc une version statique de Boost dans `$HOME/cbd/boost_static`. Compter quelques minutes seulement, car on ne compile qu'une petite partie de Boost (`system`, `chrono` et `date_time`).
```shell
mkdir -p $HOME/cbd
cd $HOME/cbd && wget https://dl.bintray.com/boostorg/release/1.65.1/source/boost_1_65_1.tar.bz2 && tar xjf boost_1_65_1.tar.bz2 && cd boost_1_65_1
./bootstrap.sh && ./b2 variant=release link=static cxxflags=-fPIC cflags=-fPIC linkflags=-fPIC --with-chrono --with-date_time --with-system --prefix=$HOME/cbd/boost_static -j4 install
```

## Récupération des sources de connector-centreon-engine

Récupérer les sources de connector-centreon-engine de Canopsis, et se positionner sur le branche `amqp-feature` :
```shell
# Note : les accès sont "External@2017"
cd $HOME/cbd && git clone https://external@git.canopsis.net/canopsis-connectors/connector-centreon-engine.git
cd connector-centreon-engine && git checkout amqp-feature
```

## Préparation des sources de Centreon Broker

On doit d'abord déterminer la version de Centreon Broker que l'on cible :
```shell
cbd -v
```

Ce qui donne, par exemple :
```
info:    Centreon Broker 3.0.14
```

Dans notre cas, il s'agit donc de Centreon Broker 3.0.14 :
```shell
export CBVER=3.0.14
```

Récupérer les sources de cette version :
```shell
cd $HOME/cbd && wget https://github.com/centreon/centreon-broker/archive/$CBVER.tar.gz && tar xzf $CBVER.tar.gz
```

Inclure les sources de notre connecteur dans celles de Centreon Broker :
```shell
cp -ru $HOME/cbd/connector-centreon-engine/build/ $HOME/cbd/connector-centreon-engine/src/amqp/ $HOME/cbd/centreon-broker-$CBVER
```

Éditer les fichiers de compilation de Centreon Broker :
```shell
vi centreon-broker-$CBVER/build/CMakeLists.txt
```

Y chercher le bloc suivant, déjà existant :
```
# TLS module.
option(WITH_MODULE_TLS "Build TLS module." ON)
if (WITH_MODULE_TLS)
  add_subdirectory("tls")
  list(APPEND MODULE_LIST "tls")
endif()
```

Et y ajouter le *nouveau* bloc suivant, à la suite :
```
# AMQP module.
option(WITH_MODULE_AMQP "Build AMQP module." ON)
if (WITH_MODULE_AMQP)
  add_subdirectory("amqp")
  list(APPEND MODULE_LIST "amqp")
endif()
```

## Compilation du module

Compiler notre module AMQP dans Centreon Broker, en forçant l'utilisation de la bibliothèque statique de Boost qui a été compilée plus tôt :
```shell
cd $HOME/cbd/centreon-broker-$CBVER/build && cmake -G "Unix Makefiles" -DCMAKE_BUILD_TYPE=Release -DBoost_INCLUDE_DIR=$HOME/cbd/boost_static/include -DWITH_MODULE_LUA=OFF
make -j3 85-amqp
```

(**Note :** tenter une compilation sans l'option `-DWITH_MODULE_LUA=OFF` si celle-ci est refusée)

Le fichier final à récupérer est `$HOME/cbd/centreon-broker-$CBVER/build/amqp/85-amqp.so`.

Utiliser le [guide d'installation du module connector-centreon-engine](doc-ce/Guide Administrateur/Quelques%20interconnexions/Centreon.md) afin de valider son bon fonctionnement.
