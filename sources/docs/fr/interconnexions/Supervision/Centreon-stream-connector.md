# Connecteur Centreon "Stream Connector"

## Description

Le connecteur convertit des évènements envoyés par le Broker Centreon en 
évènements Canopsis.
Ce connecteur utilise les fonctionnalités du Stream Connector de Centreon.

Liens connexes :

- [README des sources du connecteur][readme]
- [Documentation du stream connector de Centreon][centreon-stream-connector]

## Principe de fonctionnement

Le connecteur est développé en `lua`, le langage imposé par le mécanisme du Stream Connector.
Tous les évènements filtrés par le connecteur sont traduits au format JSON et 
envoyés sur l'**API** Canopsis via le protocole **HTTP**.

### Évènements filtrés

Les évènements de type "NEB" suivants sont actuellement gérés par le connecteur 
et correspondent à une catégorie et un élément du protocole BBDO de Centreon :

- Acquittement ou "Acknowledgment" (category 1, element 1)
- Plages de maintenance ou "Downtime" (category 1, element 5)
- Hôtes ou "Host status" (category 1, element 14)
- Services ou "Service status" (category 1, element 24)

Nous ajoutons des informations "extra" supplémentaires aux évènements hôtes et
services :

- action_url
- notes_url
- hostgroups
- servicegroups (pour les services uniquement)

#### Acquittement (ack)

Deux sortes d'actions sont envoyées à Canopsis :

- Création d'un ack
- Suppression d'un ack

L'ack est positionné sur le couple resource/component concerné.

#### Plages de maintenance (downtimes)

Deux sortes d'actions sont envoyées à Canopsis :

- Création d'un downtime
- Annulation d'un downtime

Pour chaque downtime, un identifiant unique est généré afin que l'action 
d'annulation puisse être fonctionnelle en retrouvant le downtime précèdemment
créé.

!!! warning
    Les downtimes récurrents ne sont actuellement pas gérés par le connecteur.

#### Hosts

Seuls les évènements de type HARD lors d'un changement d'état sont envoyés à Canopsis.
La traduction des états entre Centreon et Canopsis est la suivante :

| CENTREON        | CANOPSIS    |
|-----------------|-------------|
| UP (0)          | INFO (0)    |
| DOWN (1)        | CRITICAL (3)|
| UNREACHABLE (2) | MAJOR (2)   |

#### Services

Seuls les évènements de type HARD lors d'un changement d'état sont envoyés à Canopsis.
La traduction des états entre Centreon et Canopsis est la suivante :

| CENTREON        | CANOPSIS    |
|-----------------|-------------|
| OK (0)          | INFO (0)    |
| WARNING (1)     | MINOR (1)   |
| CRITICAL (2)    | CRITICAL (3)|
| UNKNOWN (3)     | MAJOR (2)   |

## Intégration du connecteur

### Prérequis

- lua version >= 5.1.4
- lua-socket library >= 3.0rc1-2
- centreon-broker version >= 20.04.12 ou >= 20.10.3

### Installation

#### Par les paquets

!!! warning
    Uniquement valable pour une version de centreon-broker version >= 20.04.12 ou >= 20.10.3

**Installation du dépôt Canopsis :**

```
echo "[canopsis]
name = canopsis
baseurl=https://repositories.canopsis.net/pulp/repos/centos7-canopsis/
gpgcheck=0
enabled=1" > /etc/yum.repos.d/canopsis.repo
```

**Installation du paquet :**

   * Pour Centreon 20.04 
   ```
   yum install canopsis-connector-centreon-stream-connector-2004
   ```
   * Pour Centreon 20.10
   ```
   yum install canopsis-connector-centreon-stream-connector-2010
   ```


!!! warning
    Si une précédente version du connecteur à été installé, il faudra la désinstaller ua préalable

   ```
   yum remove install canopsis-connector-centreon-stream-connector
   ```

#### Par les sources

!!! warning
    Compatible avec la version >= 20.04.12 ou >= 20.10.3 :

0. Récupérer les [sources du connecteur][sources]
1. Copier le script sur le serveur Centreon central dans `/usr/share/centreon-broker/lua/bbdo2canopsis.lua`.
2. Ajouter les permissions suivantes : `chown centreon-engine:centreon-engine /usr/share/centreon-broker/lua/bbdo4canopsis.lua`

#### Activation du connecteur

1. Ajout d'une nouvelle entrée ["Generic - Stream connector"][configure-centreon-broker]
   Voir les détails de la configuration dans la [section Configuration](#configuration)
2. Export de la [configuration du poller][configure-centreon-broker]
3. Redémarrage des services `systemctl restart cbd centengine gorgoned`

### Configuration

Toute la configuration du connecteur peut se faire au travers de l'interface
Centreon.

**Voici les principaux paramètres :**

| VARIABLE          | DESCRIPTION                   | VALEUR PAR DÉFAUT       |
|-------------------|-------------------------------|-------------------------|
| connector_name    | Nom du connecteur             | centreon-stream-central |
| canopsis_user     | Utilisateur de l'API          | root                    |
| canopsis_password | Mot de passe de l'utilisateur | root                    |
| canopsis_host     | Hôte Canopsis                 | localhost               |
| canopsis_port     | Port d'écoute de Canopsis     | 8082                    |

**Il est possible de modifier les paramètres de file d'attente :**

| VARIABLE          | DESCRIPTION                                   | VALEUR PAR DÉFAUT |
|-------------------|-----------------------------------------------|-------------------|
| max_buffer_age    | Durée (en secondes) de rétention des évènements avant envoi | 60  |
| max_buffer_size   | Nombre d'évènements en attente avant envoi    | 10                |

**Temps de propagation et convergence des évènements :**

| VARIABLE          | DESCRIPTION                                                    | VALEUR PAR DÉFAUT |
|-------------------|----------------------------------------------------------------|-------------------|
| init_spread_timer | Temps de propagation (en secondes) des évènements, quelque soit leur état, au démarrage du connecteur | 360 |

Étant donné que seuls les changements d'état sont transmis à Canopsis,
au moment du démarrage du connecteur, s'il existe déjà des alarmes Centreon,
alors elles ne seront pas transmises à Canopsis car aucun changement d'état ne
sera détecté.

Pour limiter ce phénomène, nous proposons une option qui permet d'envoyer à 
Canopsis tous les évènements qui circulent sans qu'il y ait forcément de changement 
d'état et ce pendant la durée du "init_spread_timer".

!!! warning
    Cela implique un pic de charge  lors de l'activation du connecteur pendant la durée du "init_spread_timer".

#### Exemple de configuration

Dans : Configuration > Pollers > Broker configuration > central-broker-master >
Output > Select "Generic - Stream connector" > Add

![centreon-configuration-screenshot](img/centreon-configuration-screenshot.png)

### Contrôle du bon fonctionnement

Connectez-vous l'interface de Canopsis et assurez-vous que les évènements et leurs
états affichés du côté de Centreon correspondent avec les évènements côté Canopsis.

Pour rappel, seules les alarmes sont envoyées (état différent de "ok").

[readme]: https://git.canopsis.net/canopsis-connectors/connector-centreon-stream-connector/-/blob/master/README.md
[sources]: https://git.canopsis.net/canopsis-connectors/connector-centreon-stream-connector
[configure-centreon-broker]: https://docs.centreon.com/current/en/developer/developer-stream-connector.html#configure-centreon-broker
[centreon-stream-connector]: https://docs.centreon.com/current/en/developer/developer-stream-connector.html#docsNav
