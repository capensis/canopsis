# Connecteur Zabbix vers Canopsis (connector-zabbix2canopsis)

Notifie les évènements des triggers Zabbix sous forme d'évènements Canopsis.

## Prérequis

- Zabbix 5.0 ou plus récent

!!! note
    Une alternative existe pour les versions plus anciennes de Zabbix, qui ne
    disposent pas de la fonctionnalité *webhook*. Voir le dépôt
    [connector-zabbix2canopsis][conn-z2c] sur GitLab pour toutes les solutions.

[conn-z2c]: https://git.canopsis.net/canopsis-connectors/connector-zabbix2canopsis

## Introduction

Le connecteur Zabbix consiste en un *Media type* de type *webhook*, qui peut
être importé dans Zabbix.

L'envoi de l'évènement à Canopsis est fait directement par le serveur Zabbix
via l'API HTTP de Canopsis, sans intermédiaire, dans les conditions définies
par l'*action* configurée.

## Mise en place

### Importer le media type

Selon votre version de Zabbix, récupérer le fichier XML ou YAML à partir du
[dépôt connector-zabbix2canopsis][conn-z2c-webhook].

Importer le média dans « Administration » > « Media types », bouton « Import ».

[conn-z2c-webhook]: https://git.canopsis.net/canopsis-connectors/connector-zabbix2canopsis/-/tree/master/webhook

### Paramétrer le media type

Dans la liste des *Media types*, vous trouvez à présent « Canopsis ».

Quelques paramètres doivent être renseignés systématiquement pour correspondre
à l'instance Canopsis cible. Modifier le nouveau média pour définir au moins
ces quatre paramètres :

- `canopsis_url` : URL de l'API Canopsis (de la forme http://canopsis:8082/)
- `canopsis_user` : nom d'utilisateur Canopsis utilisé pour l'API
- `canopsis_password` : mot de passe associé à l'utilisateur Canopsis
- `connector_name` : nom pour ce connecteur Zabbix
  (sera utilisé dans les évènements)

### Compléter la configuration

Comme pour tous les webhooks et autres types de médias dans Zabbix,
l'utilisation d'un média au sein d'une action nécessite une liaison à un
utilisateur.

À ce sujet, la [documentation Zabbix][doc-zab-webhook] recommande de créer un
utilisateur Zabbix dédié, qui représente le webhook.

Dans Zabbix on créera alors (exemple) :

- Un *User group* « Canopsis »

    À créer dans « Administration » > « User groups ».

    - Frontend access : disabled
    - Permissions en lecture sur les host groups à propos desquels il faut
    envoyer des évènements

- Un *User* « canopsis »

    À créer dans « Administration » > « Users ».

    - Membre du groupe « Canopsis »
    - Media : ajouter le média « Canopsis »

        Dans le champ « Send to », renseigner une adresse email factice comme
        « canopsis@localhost.localdomain ».

        Il est possible de personnaliser les plages horaires ou les sévérités
        pour lesquelles le média est utilisé.

    - Permissions – User type : Zabbix User

- Une *Action* (trigger action)

    À créer dans « Configuration » > « Actions », « Trigger actions ».

    Dans cet exemple minimal, on propose une action qui envoie tous les
    problèmes à Canopsis, ainsi que les résolutions des problèmes.

    - Name : « Report to Canopsis »
    - Conditions : (aucune)
    - Operations :

        - Operations : step 1, « Send message to users »
        choisir l'utilisateur canopsis créé précédemment et le média Canopsis
        - Recovery operations : « Notify all involved »

## Configuration avancée

Les paramètres du *Media type* Canopsis permettent d'effectuer différentes
configurations plus avancées pour adapter les évènements produits par le
connecteur à vos cas d'usage.

#### Attributs standard

Vous pouvez modifier le contenu qui sera mis par Zabbix dans chacun des
attributs de l'évènement (`component`, `resource`, `output`, …) via les
paramètres du même nom.

#### Traduction des sévérités

Vous pouvez personnaliser la correspondance entre les sévérités Zabbix et les
sévérités Canopsis via les paramètres `state_for_*`.

#### Attributs extra

Vous pouvez aussi passer des attributs « extra » dans l'évènement Canopsis.

Pour ce faire, ajoutez vos propres paramètres avec un nom de la forme `extra_*`
et la valeur souhaitée (texte statique, [macros Zabbix][zbx-macros] supportées
pour « Trigger-based notifications and commands », combinaison de statique et
de macros…).

Exemples :

```
extra_customnumber: 42
extra_eventid: {EVENT.ID}
extra_opdata: {EVENT.OPDATA}
```

Ceci ajouterait les attributs suivants à vos évènements (exemples) :

```json
{
  "customnumber":"42",
  "opdata":"not available (0)",
  "eventid":"52",
  "...": "..."
}
```

Rappel : pour être utiles, ces attributs « extra » doivent être exploités lors
du traitement des évènements Canopsis via des
[règles d'enrichissement][enrich].

[doc-zab-webhook]: https://www.zabbix.com/documentation/6.0/en/manual/config/notifications/media/webhook#user-media
[zbx-macros]: https://www.zabbix.com/documentation/6.0/en/manual/appendix/macros/supported_by_location
[enrich]: ../../guide-utilisation/menu-exploitation/filtres-evenements.md
