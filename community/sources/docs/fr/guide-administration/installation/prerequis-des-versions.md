# Prérequis des versions

L'usage de versions différentes de celles indiquées ici n'est pas supporté.

Les versions indiquées ici se réfèrent à la dernière version publiée de Canopsis.

## Prérequis systèmes

Solution       | Version    |
---------------|------------|
Ansible        | = 2.8.7    |
Architecture   | x86-64     |
Docker CE      | ≥ 20.10    |
Docker-compose | ≥ 1.24     |
Noyau Linux    | ≥ 4.4 (Docker uniquement) |
OS             | = CentOS 7 (paquets uniquement) |
Python         | 2.7 et 3   |

Pour rappel, SELinux n'est pas supporté. 

## Prérequis composants Canopsis

Composant   | Version          |
------------|------------------|
MongoDB     | 4.2              |
Nginx       | stable           |
PostgreSQL  | 13               |
RabbitMQ    | 3.7 (recommandé) |
Redis       | ≥ 5.0            |
TimescaleDB | 2.5              |

## Prérequis composants externes

Pour le support des navigateurs, se référer à la page des [limitations](../../guide-utilisation/limitations/index.md#compatibilite-des-anciens-navigateurs).

Enfin, consulter la page des [interconnexions](../../interconnexions/index.md), pour le support des composants externes et leurs versions.
