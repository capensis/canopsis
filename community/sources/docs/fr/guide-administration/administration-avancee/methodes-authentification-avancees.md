# Méthodes d'authentification avancées (LDAP, CAS, SAML2)

## Authentification LDAP ( édition community )

La configuration LDAP par l'interface web n'est pas prise en charge pour le moment.

Les fonctionnalités actuellement implémentées permettent l'authentification des utilisateurs sur n'importe quel annuaire LDAP, tant que celui-ci respecte la [RFC 4510](https://tools.ietf.org/html/rfc4510) et ses déclinaisons.

### Configuration de LDAP

La configuration de l'authentification se fait au moyen d'une requête sur l’API. Vous devez préparer un fichier de configuration et l'envoyer sur l'API.

Voici la liste de paramètres nécessaires à la configuration LDAP :

| Attribut        | Description                                                                                                                                | Exemple                                                        |
|-----------------|--------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------|
| ldap_uri        | Chaîne de connexion LDAP                                                                                                                   | `ldaps://ldap.example.com`                                     |
| max_tls_ver     | La version maximale de TLS qui est acceptable                                                                                              | `tls10` ou `tls11` ou `tls12` ou `tls13`                       |
| host            | Adresse du serveur LDAP <br> *Attribut obsolète conservé pour la rétrocompatibilité des configurations*                                    | `ldap.example.com`                                             |
| port            | Port d'écoute du serveur LDAP <br> *Attribut obsolète conservé pour la rétrocompatibilité des configurations*                              | `389`                                                          |
| admin_dn        | Bind DN : DN du compte utilisé pour lire l'annuaire                                                                                        | `uid=svccanopsis,ou=Special,dc=example,dc=com`                 |
| admin_passwd    | Bind password : mot de passe pour authentifier le Bind DN sur l'annuaire                                                                   |                                                                |
| user_dn         | DN de base où rechercher les utilisateurs                                                                                                  | `ou=People,dc=example,dc=com`                                  |
| ufilter         | Filtre de recherche pour les utilisateurs <br> La valeur de l'utilisateur est présentée dans une variable notée `%s`                       | `uid=%s`                                                       |
| username_attr   | Attribut portant l'identifiant utilisateur dans l'objet de l'annuaire                                                                      | `uid`                                                          |
| attrs           | Association d'attributs pour les infos de l'utilisateur <br> Un utilisateur Canopsis dispose des attributs `firstname`, `lastname`, `mail` | `{"mail": "mail", "firstname": "givenName", "lastname": "sn"}` |
| default_role    | Rôle Canopsis par défaut au moment de la première connexion                                                                                | `Visualisation`                                                |
| skip_verify     | Permet de ne pas vérifier la validité d'un certificat TLS fourni par le serveur (auto-signé, etc.)                                         | `true`                                                         |

#### Récupération de la configuration courante

Si une configuration est déjà présente en base, elle peut être récupérée à l'aide de la commande suivante :

```sh
curl -u root:root -X GET http://localhost:8082/rest/object/ldapconfig/cservice.ldapconfig
```

#### Envoi d'une configuration

La configuration se fait dans un fichier JSON : **ldapconfig.json**

```json
{
    "id": "cservice.ldapconfig",
    "crecord_type": "ldapconfig",
    "crecord_name": "ldapconfig",
    "ldap_uri": "ldap://ldap.example.com",
    "admin_dn": "uid=svccanopsis,ou=Special,dc=example,dc=com",
    "admin_passwd": "********",
    "user_dn": "ou=People,dc=example,dc=com",
    "ufilter": "uid=%s",
    "username_attr": "uid",
    "attrs": {
        "mail": "mail",
        "firstname": "givenName",
        "lastname": "sn"
    },
    "default_role": "Visualisation"
}
```

!!! note
    Vous pouvez remplacer les attributs `host` et `port` par `ldap_uri`.

La requête suivante permet d'envoyer cette configuration :

```sh
curl -u root:root -X POST \
  -H "Content-type: application/json" -d @ldapconfig.json \
  http://localhost:8082/rest/object/ldapconfig/cservice.ldapconfig
```

Le résultat renvoyé doit être de type :

```json
{"total": 1, "data": [{"..."}], "success": true}
```

Enfin, vous devez vous assurer que le fichier `/opt/canopsis/share/config/api/security/config.yml` contienne bien `ldap` dans la liste des mécanismes d'authentification :
```yaml
security:
  auth_providers:
    - ldap
    ...
```

Vous devez ensuite **obligatoirement** redémarrer les API (`api` et `oldapi` en Docker).

```sh
systemctl restart canopsis-service@canopsis-api canopsis-service@canopsis-oldapi
```

### Utilisation de LDAP

À ce stade, vous êtes en mesure de vous authentifier sur l'interface de Canopsis. Le profil d'affectation sera celui spécifié dans la configuration.

## Authentification CAS ( édition community )

La configuration de CAS par l'interface web n'est pas prise en charge pour le moment.

Les fonctionnalités actuellement implémentées permettent l'authentification des utilisateurs via WebSSO.

### Configuration de CAS

La configuration de l'authentification se fait au moyen d'un requête sur l’API. Vous devez préparer un fichier de configuration et l'envoyer sur l'API.

Voici la liste des paramètres nécessaires à la configuration CAS :

|   Attribut   |                    Description                     |            Exemple             |
| ------------ | -------------------------------------------------- | ------------------------------ |
|  `login_url`      |         URL du serveur CAS sur laquelle le navigateur web va être redirigé pour s'authentifier        |   http://canopsis.info.local/  |
| `default_role`    | Rôle par défaut au moment de la première connexion                                                    |         Visualisation          |
|   `title`         |        Label sur le formulaire de connexion                                                           |           Connexion            |
|   `validate_url`  |            URL de validation du serveur CAS à laquelle l'API va accéder                               | https://cas.info.local/websso/ |

La configuration se fait dans un fichier JSON : **casconfig.json**

```json
{
    "_id": "cservice.casconfig",
    "crecord_type": "cservice",
    "crecord_name": "casconfig",
    "enable": true,
    "default_role": "Visualisation",
    "title": "Connexion",
    "login_url": "http://localhost:8443/cas/login",
    "validate_url": "http://cas:8443/cas/serviceValidate",
    "server": "",
    "service": ""
}
```

La requête suivante permet d'envoyer cette configuration.

```sh
curl -u root:root -X POST -H "Content-type: application/json" -d @casconfig.json 'http://localhost:8082/rest/object/casconfig/cservice.casconfig'
```

Le résultat renvoyé doit être de type :

```json
{"total": 1, "data": [{"..."}], "success": true}
```

Enfin, vous devez vous assurer que le fichier `/opt/canopsis/share/config/api/security/config.yml` contienne bien `cas` dans la liste des mécanismes d'authentification :
```yaml
security:
  auth_providers:
    - cas
    ...
```

Vous devez ensuite **obligatoirement** redémarrer les API (`api` et `oldapi` en Docker).

```sh
systemctl restart canopsis-service@canopsis-api canopsis-service@canopsis-oldapi
```

### Utilisation de CAS

À ce stade, vous êtes en mesure de vous authentifier sur l'interface de Canopsis.
Le profil d'affectation sera celui spécifié dans la configuration.



## Authentification SAMLV2 ( édition community )

Intégration de l’authentification avec le protocole SAMLV2

### Configuration et Paramétrage en lien avec l'Identity Provider (IDP )

Le fichier de configuration à utiliser est le fichier `/opt/canopsis/share/config/api/security/config.yml`  utilisé par le service `api` de Canopsis. Dans le cadre d'une configuration SAMLV2, voici un exemple de fichier qui pourra être présenté en volume au conteneur Docker `api` ou directement sur le filesystem d'une installation via paquets.

```yaml
security:
  # auth_providers defines enabled authentication methods.
  # Possible values:
  # - basic Auth by username-password.
  # - apikey Auth by token.
  # - ldap Auth using LDAP service. Define LDAP config in object collection by cservice.ldapconfig id.
  # - cas Auth using CAS service. Define CAS config in object collection by cservice.casconfig id.
  # - saml Auth using SAML service. Define SAML config below(commented saml section)
  auth_providers:
    - basic
    - apikey
    - saml

  saml:
    x509_cert: /certs/saml.cert
    x509_key:  /certs/saml.key
    idp_metadata_url: <http(s)://IDP_METADATA_URL>
  # idp_metadata_xml: </path/to/xml>
    idp_attributes_map:
       email: email
       name: uid
       firstname: uid
       lastname: uid
    canopsis_saml_url: http(s)://<IP_MACHINE>/api/v4/saml
    default_role: "admin"
    insecure_skip_verify: false
    canopsis_sso_binding: redirect
    canopsis_acs_binding: redirect
    sign_auth_request: false
    name_id_format: urn:oasis:names:tc:SAML:2.0:nameid-format:persistent
    skip_signature_validation: true
    acs_index: 1
    auto_user_registration: true
```



La paire de certificats relatifs aux directives `x509_cert` et `x509_key` doit être générée en amont.

Exemple pour des certificats auto-signés :

```shell
$ openssl req -x509 -newkey rsa:2048 -keyout saml.key -out saml.cert -days 365 -nodes -subj "/CN=canopsis-saml.example.com"
```

Les options suivantes doivent ensuite être adaptées au contexte de l'IDP

| Directive                   | Définition                                                   |
| --------------------------- | ------------------------------------------------------------ |
| `idp_metadata_url`          | URL permettant de récupérer les Metadatas XML de l'IDP ( si les metadatas XML sont fournies via un service accessible ) |
| `idp_metadata_xml`          | Fichier XML contenant les Metadatas XML de l'IDP ( si les metadatas XML ne sont pas fournies via un service accessible ) |
| `idp_attributes_map`        | Tableau de correspondance entre les attributs utilisateurs de Canopsis ( colonne de gauche ) et les attributs fournis par l'IDP ( colonne de droite ) |
| `canopsis_saml_url`         | URL du service SAML fourni par Canopsis qui sera configuré côté IDP |
| `insecure_skip_verify`      | Permet de bypasser la vérification du certification de l'IDP si configuré à `true` |
| `canopsis_sso_binding`      | Type de binding HTTP pour le service SSO parmi `redirect` ou `post` |
| `canopsis_acs_binding`      | Type de binding HTTP pour le service ACS parmi `redirect` ou `post` |
| `sign_auth_request`         | Permet de signer les requêtes authentification si positionné à `true` |
| `name_id_format`            | Format du `NameIDPolicy`                                     |
| `skip_signature_validation` | Permet de bypasser la validation de la signature de l'idp lors du décodage des réponses envoyées par l'idp si positionné à `true` |
| `acs_index`                 | Valeur entière à utiliser lorsque l'on configure le service ACS Index dans les Metadata XML |
| `auto_user_registration`    | Permet de créer automatiquement les utilisateurs dans Canopsis ( s'ils n'existent pas déjà ) si cette valeur est mise à `true`|
| `default_role`              | Rôle Canopsis par défaut à attribuer pour l'utilisateur à sa création |

### Activation de l’authentification SAML2

Redémarrer le service `api` de Canopsis

* Installation via Docker Compose

=== "Canopsis Pro"
	```sh
	CPS_EDITION=pro docker-compose restart api
	```

=== "Canopsis Community"
	```sh
	CPS_EDITION=community docker-compose restart api
	```

* Installation Paquets

```sh
systemctl restart canopsis-service@canopsis-api.service
```

### Test de connexion

La mire de connexion de Canopsis doit maintenant proposer un nouveau menu de login SAML qui devra vous rediriger vers l'IDP configuré.

![saml2_login](img/saml2_login.png)

Pour tester l’authentification, il faudra vous authentifier avec un compte valide de votre IDP et ainsi arriver dans Canopsis sans erreur.

### Troubleshooting

Observer les logs du service `api` et vérifier la non présence de pattern de type `ERR`

Redémarrer le service `api` de Canopsis

* Installation via Docker Compose
=== "Canopsis Pro"
	```sh
	CPS_EDITION=pro docker-compose logs -f api
	```

=== "Canopsis Community"
	```sh
	CPS_EDITION=community docker-compose logs -f api
	```

* Installation Paquets

```sh
$ journactl -fu canopsis-service@canopsis-api.service
```

