# Mail vers Canopsis

!!! attention
    Ce connecteur n'est disponible que dans l'édition CAT de Canopsis.

## Introduction

Le connecteur email2canopsis permet de lire des emails dans une boite aux lettres POP3 pour les convertir en événements Canopsis (grâce à un système de template).

## Fonctionnement

Les événements générés embarquent des attributs dont les valeurs sont contenues dans les emails et parsés par un mécanisme de template.

### Configuration

Le fichier de configuration du connecteur est `/opt/canopsis_connectors/email2canopsis/etc/email2canopsis.ini`.

```ini
# Chaine de connexion au bus AMQP
[amqp]
url=amqp://cpsrabbit:canopsis@vip-dvp-hypervision-rabbitmq.si3si.int/canopsis

# URL d'un Redis pour activer le renvoi automatique du dernier état connu des événements
[redis]
url=redis://redis:6379/0

[mail]
# Adresse du serveur POP3
popserver=webmail.test.net
# Port du serveur POP3 (110 ou 995 généralement)
port=995
# Doit-on négocier en SSL ou en PLAIN ?
overssl=True
# Identifiant de la BAL
user=email2canopsis
# Mot de passe associé
password=
# Dossier d'archivage des emails (active si présent)
archive_folder=

[event]
connector.constant=mail2canopsis
connector_name.constant=instance1
event_type.constant=check
source_type.constant=resource
component.value=param1
resource.value=param2
output.value=paramn
state.constant=0

[template]
# Expéditeur des emails à traiter
template1.sender=sender@mail.net
# Template à appliquer sur ces emails
template1.path=/opt/canopsis_connectors/email2canopsis/etc/template_1.conf
```

#### Configuration de la connexion RabbitMQ et Redis

Les blocs `[amqp]` et `[redis]` contiennent respectivement la configuration pour la connexion au bus RabbitMQ et au cache Redis.

Il faut donc vérifier que les URL qui y figurent sont les bonnes.

#### Configuration de la Boîte aux Lettres

Le bloc `[mail]` contient la configuration pour la connexion à la boîte aux lettres POP3 dont les emails seront convertis en événements Canopsis.

Il faut donc vérifier que les différents champs sont bien remplis.

#### Configuration des templates

Le bloc `template` contient la configuration des templates.

Pour la recette, il faut s'assurer que l'adresse depuis laquelle on va envoyer un email se trouve bien dans le bloc et que le contenu de l'email corresponde bien au template appliqué à cette même adresse email.

Ici, pour

```ini
[template]
# Expéditeur des emails à traiter
template1.sender=sender@mail.net
# Template à appliquer sur ces emails
template1.path=/opt/canopsis_connectors/email2canopsis/etc/template_1.conf
```

Il faut envoyer un email depuis l'adresse `sender@mail.net` et son contenu doit correspondre au template `/opt/canopsis_connectors/email2canopsis/etc/template_1.conf`.

### Configuration du template d'email

#### Principe du template

Les templates permettent d'analyser un email en suivant quelques régles de transformations simples.

Un fichier template devrait ressembler à quelque chose comme suit :

```ini
[tpl]
component=MAIL_SENDER
output=MAIL_BODY.line(6).word(3).untilword()
long_output=MAIL_BODY.line(7).after(le).untilword()
state=MAIL_SUBJECT.line(0).word(1)
state.converter=Mineur>1,Majeur>2,Critique>3
timestamp=MAIL_DATE
timestamp.output=timestamp
```

Les clefs sont divisées en deux parties, séparées par un point :

- à gauche, une *racine* : le nom de la variable à insérer dans l'événement généré
- à droite, une *méthode* : c'est-à-dire l'action de transformation à appliquer

Les actions peuvent être les suivantes :

* *selector* (utilisé par défaut ; implicite) : applique simplement le template à droite et copie la valeur traduite dans l'événement.
* *converter* : remplace une chaîne de caractères par une autre (insensiblement à la casse), les deux étant séparés par le symbole '>'. Plusieurs conversions sont applicables à la suite en les séparant par des virgules. Dans l'exemple ci-dessus, 'Mineur' sera remplacé par 1, 'Majeur' par 2…

La partie droite décrit les régles de transformations (où a, b et c sont des entiers, et d, e des chaînes de caractères) :

- `MAIL_BODY`, `MAIL_DATE`, `MAIL_ID`, `MAIL_SENDER`, et `MAIL_SUBJECT` sont les différentes parties de l'email
- `line(a)` sélectionne une ligne entière numéro a
- `line(a).word(b)` sélectionne le b-ième mot de la ligne a
- `line(a).word(b).untilword(c)` sélectionne tous les mots entre le b-ième et le c-ième sur la ligne a
- `line(a).word(b).untilword()` sélectionne tous les mots à partir du b-ième jusqu'à la fin de la ligne a
- `line(a).after(d)` sélectionne tous les mots après le mot d
- `line(a).after(d).untilword(c)` sélectionne les mots après le mot d et c mots ensuite
- `line(a).after(d).word(b).untilword(c)` sélectionne les mots après le mot d, à partir du b-ième jusqu'au c-ième
- `line(a).before(e)` sélectionne tous les mots avant e
- `line(a).before(e).word(c)` sélectionne tous les mots avant e, mais en commençant au c-ième

MAIL\_DATE est automatiquement converti en objet date, inutile d'appliquer une action 'dateformat' dessus.

!!! attention
    Les numéros de lignes et de mots commencent à partir de 0, non de 1.

Exemple : La séquence à la 1° ligne située entre les 5° et le 18° mots sont donc sélectionnables avec la ligne `line(0).word(4).untilword(17)`.

## Dépannage

#### Appliquer des changements de configuration

Pour appliquer un changement (modification de la configuration, ajout de templates, etc.), il faut redémarrer le connecteur.
