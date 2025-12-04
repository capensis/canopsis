# Matrice des flux réseau

## Liste des ports Canopsis

Composant     | Description                                 | Port                  |
--------------|---------------------------------------------|-----------------------|
MongoDB       | Base de données                             | TCP/27017             |
Nginx         | Accès à l'interface web et API              | TCP/443 (installation via Docker) ou TCP/8443 (via paquets) |
RabbitMQ      | Passage de messages                         | TCP/5672              |
RabbitMQ UI   | Interface web de RabbitMQ                   | TCP/15672             |
API Canopsis  | API REST de Canopsis                        | TCP/8082              |
Redis         | Serveur de cache                            | TCP/6739              |
Redis Sentinel| Supervision et basculement de Redis (optionnel) | TCP/26739         |
SNMP          | Passage des traps SNMP                      | UDP/162               |
PostgreSQL    | Base de données, métriques (TimescaleDB)    | TCP/5432              |

## Matrice des flux

Ci-dessous la matrice des flux réseaux commune, pour les différents composants
d'une plate-forme Canopsis.
Cette matrice ne comprend pas les différentes [interconnexions avec les autres
applications](../../interconnexions/index.md) avec lesquelles Canopsis peut
communiquer. Il faudra donc compléter cette liste avec les différents
composants additionnels, par exemple l'accès aux outils de remédiation ou de
ticketing (pour plus de précisions sur cela, voir [Interactions avec des services
extérieurs](#interactions-avec-des-services-exterieurs)).

Certains flux de cette liste sont nécessaires pour l'installation ou la mise à
jour de Canopsis. D'autres concernent l'administration de Canopsis ainsi que
les accès utilisateurs et sources d'évènements.

Source | Destination | Port | Description |
-------|-------------|------|-------------|
Canopsis | `git.canopsis.net`, `nexus.canopsis.net`, `docker.canopsis.net` | TCP/443 | Récupération des paquets d'installation (Utilisation possible à travers un proxy) |
Utilisateurs | Canopsis | TCP/443 ou TCP/8443 | Accès à l'interface web et API de Canopsis |
Administrateurs | Canopsis, MongoDB, PostgreSQL, RabbitMQ, Redis | TCP/22 | Accès aux systèmes via SSH |
Administrateurs | Canopsis | TCP/15672 | Accès à l'interface web du bus AMQP. Permet de suivre l'activité des files d'attente |
Sources d'événements AMQP | Canopsis | TCP/5672 | Permet la publication d'événements dans le bus de données |
Sources d'événements API | Canopsis | TCP/443 ou TCP/8443 | Permet la publication d'événements dans l'API |
Sources d'événements trap SNMP | Canopsis | UDP/162 | Permet la publication de trap SNMP vers Canopsis |
Canopsis | LDAP | TCP/389,636 | Permet l'authentification à Canopsis via un identifiant LDAP |
Canopsis | MongoDB | TCP/27017 | Permet l'accès à la base de données MongoDB depuis Canopsis |
Canopsis | PostgreSQL | TCP/5432 | Permet l'accès à la base de données PostgreSQL depuis Canopsis |
Canopsis | RabbitMQ | TCP/5672 | Permet l'accès à l'agent de messages RabbitMQ depuis Canopsis |
Canopsis | Redis | TCP/6739 | Permet l'accès à la base de données Redis depuis Canopsis |
Canopsis | Redis Sentinel | TCP/26739 | Permet l'accès à la supervision Redis Sentinel depuis Canopsis (optionnel) |

Définition des objets:

 * Utilisateurs : Postes de travail des utilisateurs de la solution
 * Administrateurs : Postes de travail des administrateurs de la solution ou des bastions associés
 * Sources d'événements : Machines qui produisent des événements au format AMQP/JSON (supervision, scripts, curl, etc)
 * Sources d'événements trap SNMP : Machines qui produisent des événements au format Trap SNMP
 * Canopsis : Machine qui héberge Canopsis
 * MongoDB, PostgreSQL, RabbitMQ, Redis : Machine(s) qui héberge(nt) MongoDB, PostgreSQL, RabbitMQ et Redis

## Interactions avec des services extérieurs

Dans de nombreux cas d'installation il est utile de savoir précisément quels
services applicatifs de Canopsis pourront joindre des services extérieurs à
Canopsis, voire extérieurs à votre réseau interne.

La réalité dépendra évidemment des fonctionnalités utilisées et des règles ou
configurations que vous mettez en place. Nous pouvons toutefois lister de
manière générale qu'au plus, les appels vers des services extérieurs auront
comme origine les processus suivants :

- Le moteur `engine-che` :
  pour de la récupération de données externes via API depuis des [règles
  d'enrichissement][gu-eventfilter] ;

- Le moteur `engine-webhook` :
  pour tous les appels HTTP effectués du fait de vos [Scénarios][gu-scenario]
  et [Règles de déclaration de tickets][gu-ticketrules], ainsi que pour les
  appels d'obtention et rafraîchissement de [jetons d'authentification
  externes][gu-authtokens] ;

- Le moteur `remediation` :
  pour tout appel à une API d'un outil d'automatisation utilisé dans le cadre
  de la [Remédiation][ga-remediation] ;

- Le service `api` :
  il est susceptible de joindre un service extérieur si vous configurez une des
  méthodes d'[authentification externe][ga-externalauth] (LDAP, CAS, SAML,
  OAuth 2.0 ou OpenID).

Cette liste précise vous permet ainsi d'anticiper :

- L'ouverture des différents flux réseau ;

- L'utilisation d'un proxy HTTP(S) à renseigner dans les [variables
  d'environnement][ga-envvars-proxy] appropriées (cas d'un outil de ticketing
  en SaaS par exemple : nécessite Internet, souvent via un proxy) ;

- L'installation de vos autorités de certification racine internes,
  nécessaires au bon établissement des communications vers des
  protocoles sécurisés (HTTPS, LDAPS) :

    Selon le mode d'installation retenu, il va s'agir…

    - sur des systèmes RHEL, d'ajouter le certificat au magasin de confiance
      système : ref. [Adding new certificates to the system-wide
      truststore][rhel-ca-trust] (documentation Red Hat) ;
    - sur des déploiements conteneurisés, de préparer et présenter un fichier
      complet de toutes les CA approuvées (celles par défaut de l'image Docker
      Alpine Linux + les vôtres) dans les conteneurs concernés – le fichier qui
      en résulte devra être monté en tant que volume au chemin
      `/etc/ssl/certs/ca-certificates.crt`.

- La mise en place des bonnes règles de trafic sortant dans certains systèmes
  conteneurisés (Egress).

[gu-eventfilter]: ../../../guide-utilisation/menu-exploitation/filtres-evenements/
[gu-scenario]: ../../../guide-utilisation/menu-exploitation/scenarios/
[gu-ticketrules]: ../../../guide-utilisation/menu-exploitation/regles-declaration-tickets/
[gu-authtokens]: ../../../guide-utilisation/menu-administration/jetons-authentification-externe/
[ga-remediation]: ../remediation/
[ga-externalauth]: ../administration-avancee/methodes-authentification-avancees/
[ga-envvars-proxy]: ../administration-avancee/variables-environnement/#utilisation-dun-proxy-http-ou-https
[rhel-ca-trust]: https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/9/html/securing_networks/using-shared-system-certificates_securing-networks#adding-new-certificates_using-shared-system-certificates
