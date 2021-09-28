# Prérequis des versions

L'usage de versions différentes de celles indiquées ici n'est pas supporté.

## Prérequis systèmes

Solution       | Versions   |
---------------|------------|
Ansible        | = 2.8.7    |
Architecture   | x86-64     |
Docker CE      | ≥ 19.03.5  |
Docker-compose | ≥ 1.24     |
Kernel         | ≥ 4.4 (Docker uniquement) |
OS CentOS      | = 7        |
Python         | 2.7 et 3   |

Pour rappel, SELinux n'est pas supporté. 

Le support de l'IPv6 est possible à condition de le configurer.

## Prérequis composants Canopsis

Composant | Versions         |
----------|------------------|
InfluxDB  | 1.5              |
MongoDB   | 3.6              |
Nginx     | stable           |
RabbitMQ  | 3.7 (recommandé) |
Redis     | ≥ 5.0            |

## Prérequis composants externes

Pour le support des navigateurs, se référer à la page des [limitations](../../guide-utilisation/limitations/index.md#compatibilite-des-anciens-navigateurs).

Enfin, consulter la page des [interconnexions](../../interconnexions/index.md), pour le support des composants externes et leurs versions.
