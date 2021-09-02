# Sécurisation d'une installation de Canopsis et de ses composants

## SELinux

Ceci concerne majoritairement les environnements CentOS et assimilés.

SELinux n'est pas pris en charge par Canopsis.

Il est nécessaire de le mettre en mode permissif ou de le désactiver sur les nœuds où vous installez Canopsis :

```sh
setenforce 0
sed -i 's/^SELINUX=.*$/SELINUX=permissive/' /etc/selinux/config
```

Puis, redémarrer le système.

## Sécurisation des ports réseau

Cette section décrit les ports utilisés par Canopsis qui doivent être autorisés.

Voir aussi la [matrice des flux entre les composants](../installation/pre-requis-parefeu-et-selinux.md).

### Liste des ports

Composant     | Description                                 | Port                  |
--------------|---------------------------------------------|-----------------------|
MongoDB       | Base de données                             | 27017/TCP             |
RabbitMQ      | Passage de messages                         | 5672/TCP              |
RabbitMQ UI   | Interface web de RabbitMQ (recommandée)     | 15672/TCP             |
API Canopsis  | API REST de Canopsis                        | 8082/TCP              |
Nginx         | Accès à l'interface web                     | 80/TCP                |
Redis         | Serveur de cache                            | 6739/TCP              |
InfluxDB      | Métriques                                   | 8086/TCP              |

### Détails

#### MongoDB

Dans le cadre d'une installation mono-instance, fermez le port.

#### RabbitMQ

RabbitMQ doit être autorisé pour tous les connecteurs de Canopsis.

Il est utilisé pour transmettre des évènements dans Canopsis. Laisser le port ouvert dans votre pare-feu et restreindre les adresses IP autorisées.

#### Interface web RabbitMQ

Cette interface Web ne doit être utilisée que par les administrateurs et ne doit pas être ouverte aux utilisateurs.

Tout comme précédemment, laissez le port ouvert avec des restrictions sur les adresses IP.

#### Interface web Canopsis

Le serveur Web, fournissant l'interface web et l'API REST, doit être autorisé pour chaque utilisateur.

Laissez le port ouvert dans le pare-feu, avec les restrictions adéquates.
