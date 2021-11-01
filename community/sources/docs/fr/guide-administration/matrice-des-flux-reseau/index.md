# Matrice des flux réseau

## Liste des ports Canopsis

Composant     | Description                                 | Port                  |
--------------|---------------------------------------------|-----------------------|
MongoDB       | Base de données                             | TCP/27017             |
RabbitMQ      | Passage de messages                         | TCP/5672              |
RabbitMQ UI   | Interface web de RabbitMQ                   | TCP/15672             |
API Canopsis  | API REST de Canopsis                        | TCP/8082              |
Nginx         | Accès à l'interface web et API              | TCP/80,443            |
Redis         | Serveur de cache                            | TCP/6739              |
InfluxDB      | Métriques                                   | TCP/8086              |
SNMP          | Passage des traps SNMP                      | UDP/162               |

## Matrice des flux

Ci-dessous la matrice des flux réseaux des différents composants de Canopsis. Cette matrice ne comprend pas les différentes [interconnexions avec les autres applications](../../interconnexions/index.md) avec lesquelles Canopsis peut communiquer. Il faudra donc compléter cette liste avec les différents composants additionnels, par exemple l'accès aux outils de remédiations ou de ticketing.

Certains flux de cette liste sont nécessaires pour l'installation ou la mise à jour de Canopsis. D'autres concernent l'administration de Canopsis ainsi que les accès utilisateurs et sources d'évènements.

Source | Destination | Port | Description |
-------|-------------|------|-------------|
Canopsis | git.canopsis.net, hub.docker.com | TCP/443 | Récupération des paquets d'installation (Utilisation possible à travers un proxy) |
Utilisateurs | Canopsis | TCP/80,443 | Accès à l'interface web et API de Canopsis |
Administrateurs | Canopsis, MongoDB | TCP/22 | Accès aux systèmes via SSH |
Administrateurs | Canopsis | TCP/15672 | Accès à l'interface web du bus AMQP. Permet de suivre l'activité des files d'attente |
Sources d'événements AMQP | Canopsis | TCP/5672 | Permet la publication d'événements dans le bus de données |
Sources d'événements API | Canopsis | TCP/80,443 ou TCP/8082 | Permet la publication d'événements dans l'API |
Sources d'événements trap SNMP | Canopsis | UDP/162 | Permet la publication de trap SNMP vers Canopsis |
Canopsis | LDAP | TCP/389,636 | Permet l'authentification à Canopsis via un identifiant LDAP |
Canopsis | MongoDB | TCP/27017 | Permet l'accès à la base de données MongoDB depuis Canopsis |

Définition des objets:

 * Utilisateurs : Postes de travail des utilisateurs de la solution
 * Administrateurs : Postes de travail des administrateurs de la solution ou des bastions associés
 * Sources d'événements : Machines qui produisent des événements au format AMQP/JSON (supervision, scripts, curl, etc)
 * Sources d'événements trap SNMP : Machines qui produisent des événements au format Trap SNMP
 * Canopsis : Machine qui héberge Canopsis
 * MongoDB : Machine(s) qui héberge(nt) MongoDB

