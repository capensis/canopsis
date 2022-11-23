# Méthodes d'authentification avancées (LDAP, CAS, SAML2)

!!! information "Information"

    Quelque soit le mécanisme d'authentification utilisé, vous pouvez configurer 2 paramètres concernant les expirations d'authentification : `inactivity_interval` et `expiration_interval`.
    
    * `inactivity_interval` : Délai de session utilisateur avant expiration sur inactivité
    * `expiration_interval` : Délai de session utilisateur avant expiration
    
    Exemple : si `expiration_interval` vaut **1 mois** et que `inactivity_interval` vaut **24 heures**
    
    Cas 1:
    
    * L'utilisateur se loggue
    * L'utilisateur ouvre son interface Canopsis chaque jour
    * L'utilisateur sera déloggué après 1 mois
    
    Cas 2:
    
    * L'utilisateur se loggue
    * L'utilisateur ouvre son interface Canopsis chaque jour ouvré
    * L'utilisateur sera déloggué dimanche
    * L'utilisateur se loggue à nouveau lundi

## Authentification LDAP

Les fonctionnalités actuellement implémentées permettent l'authentification des utilisateurs sur n'importe quel annuaire LDAP, tant que celui-ci respecte la [RFC 4510](https://tools.ietf.org/html/rfc4510) et ses déclinaisons.


### Configuration de LDAP

La configuration de l'authentification se fait au travers du fichier de configuration de l'API `/opt/canopsis/share/config/api/security/config.yml`.  

Tout d'abord, vous devez activer le mécanisme d'authentification LDAP en lui-même :

```yaml
security:
  auth_providers:
    - ldap
    ...
```

Puis vous devez renseigner les différents paramètres d'authentification LDAP.

```yaml
  ldap:
    inactivity_interval: 24h
    expiration_interval: 1M
    url: ldap://ldap.local
    admin_dn: uid=svccanopsis,ou=Special,dc=example,dc=com
    admin_passwd:
    user_dn: ou=People,dc=example,dc=com
    ufilter: uid=%s
    username_attr: uid
    attrs:
      mail: mail
      firstname: givenName
      lastname: sn
    default_role: Visualisation
    insecure_skip_verify: false
    max_tls_ver:
```

Définition des paramètres :

| Attribut        | Description                                        | Exemple                                                         |
|-----------------|----------------------------------------------------|-----------------------------------------------------------------|
| `url`        | Chaîne de connexion LDAP                                | `ldaps://ldap.example.com`                                      |
| `admin_dn`        | Bind DN : DN du compte utilisé pour lire l'annuaire | `uid=svccanopsis,ou=Special,dc=example,dc=com`                 |
| `admin_passwd`    | Bind password : mot de passe pour authentifier le Bind DN sur l'annuaire  |                                          |
| `user_dn`         | DN de base où rechercher les utilisateurs          | `ou=People,dc=example,dc=com`                                   |
| `ufilter`         | Filtre de recherche pour les utilisateurs <br> La valeur de l'utilisateur est présentée dans une variable notée `%s` | `uid=%s`    |
| `username_attr`   | Attribut portant l'identifiant utilisateur dans l'objet de l'annuaire  | `uid`                                       |
| `attrs`           | Association d'attributs pour les infos de l'utilisateur <br> Un utilisateur Canopsis dispose des attributs `firstname`, `lastname`, `mail` | `{"mail": "mail", "firstname": "givenName", "lastname": "sn"}` |
| `default_role`    | Rôle Canopsis par défaut au moment de la première connexion   | `Visualisation`                                      |
| `insecure_skip_verify` | Permet de ne pas vérifier la validité d'un certificat TLS fourni par le serveur (auto-signé, etc.)   | `true`   |
| `max_tls_ver` (optionnel) | La version maximale de TLS qui est acceptable      | `tls10` ou `tls11` ou `tls12` ou `tls13`                        |

Vous devez ensuite **obligatoirement** redémarrer le service API.

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

### Utilisation de LDAP

À ce stade, vous êtes en mesure de vous authentifier sur l'interface de Canopsis. Le profil d'affectation sera celui spécifié dans la configuration.

## Authentification CAS ( édition community )

Les fonctionnalités actuellement implémentées permettent l'authentification des utilisateurs via WebSSO.

### Configuration de CAS

La configuration de l'authentification se fait au travers du fichier de configuration de l'API `/opt/canopsis/share/config/api/security/config.yml`.  

Tout d'abord, vous devez activer le mécanisme d'authentification CAS en lui-même :

```yaml
security:
  auth_providers:
    - cas
    ...
```

Puis vous devez renseigner les différents paramètres d'authentification CAS.

```yaml
  cas:
    inactivity_interval: 24h
    expiration_interval: 1M
    # title defines label of UI login form.
    title: Connexion
    # login_url defines CAS login url to which UI is redirected to authenticate.
    login_url: http://cas.local/login
    # validate_url defines CAS validate url which is used to validate received ticket.
    validate_url: http://cas.local/serviceValidate
    # default_role defines role of new users which are created on successful CAS login.
    default_role: Visualisation
```

Définition des paramètres :

| Attribut       |                    Description                               |            Exemple             |
| -------------- | ------------------------------------------------------------ | ------------------------------ |
| `login_url`    | URL du serveur CAS sur laquelle le navigateur web va être redirigé pour s'authentifier        |   http://canopsis.info.local/  |
| `default_role` | Rôle par défaut au moment de la première connexion           | Visualisation                  |
| `title`        | Label sur le formulaire de connexion                         | Connexion                      |
| `validate_url`  | URL de validation du serveur CAS à laquelle l'API va accéder | https://cas.info.local/websso/ |

Vous devez ensuite **obligatoirement** redémarrer le service API.

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

### Utilisation de CAS

À ce stade, vous êtes en mesure de vous authentifier sur l'interface de Canopsis.
Le profil d'affectation sera celui spécifié dans la configuration.


## Authentification SAMLv2

Intégration de l’authentification avec le protocole SAMLV2

### Configuration et Paramétrage en lien avec l'Identity Provider (IDP )

La configuration de l'authentification se fait au travers du fichier de configuration de l'API `/opt/canopsis/share/config/api/security/config.yml`.  

Tout d'abord, vous devez activer le mécanisme d'authentification SAML en lui-même :

```yaml
security:
  auth_providers:
    - saml
    ...
```

Puis vous devez renseigner les différents paramètres d'authentification SAML.

```yaml
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

Définition des paramètres :

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
| `insecure_skip_verify`      | Permet de ne pas vérifier la validité d'un certificat TLS fourni par le serveur (auto-signé, etc.)   | `false`   |

Vous devez ensuite **obligatoirement** redémarrer le service API.

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

