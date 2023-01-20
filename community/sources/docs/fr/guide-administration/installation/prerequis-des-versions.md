# Prérequis des versions

L'usage de versions différentes de celles indiquées ici n'est pas supporté.

Les versions indiquées ici se réfèrent à la dernière version publiée de Canopsis.

## Prérequis systèmes

Solution       | Version    |
---------------|------------|
Architecture   | x86-64     |
Docker CE      | ≥ 20.10.17 |
Docker Compose | ≥ 2.12     |
Noyau Linux    | ≥ 4.4 (uniquement pour l'installation via Docker Compose)             |
OS             | = CentOS 7 ou RHEL 8 (uniquement pour l'installation via paquets RPM) |
Python         | 3          |

Pour rappel, SELinux n'est pas supporté. 

## Prérequis composants Canopsis

Composant   | Version          |
------------|------------------|
MongoDB     | 4.4              |
Nginx       | 1.20 (uniquement pour l'installation via paquets RPM) |
PostgreSQL  | 13               |
TimescaleDB | 2.7.2            |
RabbitMQ    | 3.10             |
Redis       | ≥ 5.0, < 7.0     |

## Prérequis composants externes

Pour le support des navigateurs, se référer à la page des [limitations](../../guide-utilisation/limitations/index.md#compatibilite-des-anciens-navigateurs).

Enfin, consulter la page des [interconnexions](../../interconnexions/index.md), pour le support des composants externes et leurs versions.
