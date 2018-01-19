## SAML2

Intégration de l’authentification avec SAML2

Nécessite l’installation de la brique `CAT`.

### Paramétrage IdP

 * RelayState : `http://url-canopsis/`
 * Audience : `http://url-canopsis/auth/saml2/metadata/`
 * Recipient : `http://url-canopsis/auth/saml2/acs/`
 * ACS Consumer URL : `http://url-canopsis/auth/saml2/acs/`
 * Single Logout URL : `http://url-canopsis/auth/saml2/sls/`

L’IdP doit impérativement fournir dans les réponses d’authentification une valeur normalisée `NameID`. Il suffit de créer un *mapping* entre ce champs normalisé et une information unique dans le backend utilisé par l’IdP. Dans le cas contraire l’authentification côté Canopsis ne **pourra pas fonctionner**.

### Création des paramètres - Côté Canopsis

Vous pouvez suivre cette documentation : https://github.com/onelogin/python-saml#knowing-the-toolkit

En particulier la génération des clefs et des paramètres :

```
mkdir certs
openssl req -new -x509 -days 3652 -nodes -out certs/sp.crt -keyout certs/sp.key
```

`settings.json` :

```json
{
    "strict": true,
    "debug": false,
    "sp": {
        "entityId": "http://canopsis.fqdn:port/auth/saml2/metadata/",
        "assertionConsumerService": {
            "url": "http://canopsis.fqdn:port/auth/saml2/acs/",
            "binding": "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
        },
        "singleLogoutService": {
            "url": "http://canopsis.fqdn:port/auth/saml2/sls/",
            "binding": "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect"
        },
        "NameIDFormat": "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
        "x509cert": "",
        "privateKey": ""
    },
    "idp": {
        "entityId": "https://app.onelogin.com/saml/metadata/<appid>",
        "singleSignOnService": {
            "url": "https://app.onelogin.com/trust/saml2/http-post/sso/<appid>",
            "binding": "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect"
        },
        "singleLogoutService": {
            "url": "https://app.onelogin.com/trust/saml2/http-redirect/slo/<appid>",
            "binding": "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect"
        },
        "x509cert": "<certificat public de l'identity provider>"
    }
}
```

Remplacer les occurrences de :

 * Les url "app.onelogin", elle sont a remplacer intégralement par les données de configuration de L'IdP sur lequel vous allez vous brancher.
 * `canopsis.fqdn:port`
 * `<certificat public de l'identity provider>`

Pour le certificat public, vous devrez le télécharger au format PEM et exécuter la commande suivante afin de récupérer le contenu :

```
cat idp_cert.pem | grep -v "BEGIN CERTIFICATE" | grep -v "END CERTIFICATE" | tr '\n' ' ' | sed -e 's/ //g'
```

`advanced_settings.json` :

```json
{
    "security": {
        "nameIdEncrypted": false,
        "authnRequestsSigned": false,
        "logoutRequestSigned": false,
        "logoutResponseSigned": false,
        "signMetadata": false,
        "wantMessagesSigned": false,
        "wantAssertionsSigned": false,
        "wXoantNameId" : true,
        "wantNameIdEncrypted": false,
        "wantAssertionsEncrypted": false,
        "signatureAlgorithm": "http://www.w3.org/2000/09/xmldsig#rsa-sha1",
        "digestAlgorithm": "http://www.w3.org/2000/09/xmldsig#sha1"
    },
    "contactPerson": {
        "technical": {
            "givenName": "technical_name",
            "emailAddress": "technical@example.com"
        },
        "support": {
            "givenName": "support_name",
            "emailAddress": "support@example.com"
        }
    },
    "organization": {
        "en-US": {
            "name": "sp_test",
            "displayname": "SP test",
            "url": "http://sp.example.com"
        }
    }
}
```

Ces fichiers de configuration sont à adapter.

### Intégration des paramètres en base

Créer cette structure :

```
saml2/
    certs/
        sp.crt
        sp.key
    settings.json
    advanced_settings.json
    conf_path
    secret_key
```

Le fichier `conf_path` devra contenir le chemin de destination de la configuration SAML2 lorsqu’elle sera utilisée par le webserver canopsis.

Le fichier `secret_key` permettra de déchiffrer les données SAML2 en cas de chiffrement. Si vous n’activez pas le chiffrement, créez quand même ce fichier.

Ensuite, dans l’environnement Canopsis, exécutez ceci dans un shell Python :

```
python -c 'from canopsis_cat.saml2 import SAML2Conf; SAML2Conf.insert_conf("<path to saml2 source conf directory>", SAML2Conf.provide_default_collection())'
```

Vous pouvez relancer cette commande autant de fois que nécessaire : la configuration en place sera tout simplement écrasée intégralement.

Donc si vous voulez apporter une modification de la configuration, pas besoin de passer par la base de donnée : modifiez les fichiers "source" sur disque, exécutez la commande, c’est fini.

### Activation de l’authentification SAML2

Éditer le fichier de configuration canopsis `etc/webserver.conf` :

```ini
[auth]
; version pré-monopackage
providers = saml2
; version monopackage
providers = canopsis_cat.auth.saml2

[webservices]
; version pré-monopackage
saml2 = 1
; version monopackage
canopsis_cat.webcore.services.saml2 = 1
```

Puis exécutez :

```
su - canopsis -c "service webserver restart"
```

### Tests et log

Le fichier de log `var/log/saml2.log` contiendra les erreurs SAML2 s’il y en a.

Pour tester l’authentification :

 * Rendez-vous sur la page de login de canopsis
 * Entrez un utilisateur autre que ceux présents dans Canopsis, et n’importe quoi en mot de passe (changements à venir)
 * Vous devez être redirigé vers la page de login de l’IdP SAML2
 * Une fois authentifié via l’IdP, vous devez être redirigé vers Canopsis sans erreur.
