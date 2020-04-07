# Configuration de pare-feu et de SELinux

Ce document décrit les ports utilisés par Canopsis qui doivent être autorisés.

Voir aussi la [matrice des flux entre les services](../installation/pre-requis-parefeu-et-selinux.md).

## Liste des ports

Service       | Description                                 | Port                  |
--------------|---------------------------------------------|-----------------------|
MongoDB       | Base de données                             | 27017/TCP             |
RabbitMQ      | Passage de messages                         | 5672/TCP              |
RabbitMQ UI   | Interface web de RabbitMQ (recommandée)     | 15672/TCP             |
Gunicorn      | Interface web et API REST de Canopsis       | 8082/TCP              |
Redis         | Serveur de cache                            | 6739/TCP              |
InfluxDB      | Métriques                                   | 8086/TCP              |

## Détails

### MongoDB

Dans le cadre d'une installation mono-instance, fermez le port.

### RabbitMQ

RabbitMQ doit être autorisé pour tous les connecteurs de Canopsis.

Il est utilisé pour transmettre des évènements dans Canopsis, laisser le port ouvert dans votre pare-feu et restreindre les adresses IP autorisées.

### Interface web RabbitMQ

Cette interface Web ne doit être utilisée que par les administrateurs et ne doit pas être ouverte aux utilisateurs.

Tout comme précédemment, laissez le port ouvert avec des restrictions sur les adresses IP.

### Interface web Canopsis

Le serveur Web, fournissant l'interface web et l'API REST, doit être autorisé pour chaque utilisateur.

Laissez le port ouvert dans le pare-feu, avec les restrictions adéquates.

## SELinux

SELinux n'est pas pris en charge par Canopsis.

Il est alors préférable de le désactiver :

```sh
setenforce 0
```

Ceci concerne majoritairement les environnements CentOS et assimilés.
