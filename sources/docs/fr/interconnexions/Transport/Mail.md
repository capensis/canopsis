# Mail vers Canopsis

!!! attention
    Ce connecteur n'est disponible que dans l'édition CAT de Canopsis.

## Introduction

Le connecteur email2canopsis permet de lire des emails dans une boîte aux lettres POP3 pour les convertir en événements Canopsis (grâce à un système de template).

## Fonctionnement

Les événements générés embarquent des attributs dont les valeurs sont contenues dans les emails et parsées par un mécanisme de template.

### Configuration

Le fichier de configuration du connecteur est `/opt/canopsis_connectors/email2canopsis/etc/email2canopsis.ini`.

```ini
# Adresse de connexion à RabbitMQ
[amqp]
url=amqp://cpsrabbit:canopsis@rabbitmq/canopsis

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
# Depuis la 3.21.0 : option leavemails pour conserver sur le serveur les emails lus
leavemails=False

[event]
connector.constant=mail2canopsis
connector_name.constant=instance1
event_type.constant=check
source_type.constant=resource
component.value=param1
resource.value=param2
output.value=paramn
state.constant=0

[event_error]
connector.constant=email2canopsis
event_type.constant=check
source_type.constant=resource
connector_name.constant=resource
component.constant=connection
resource.constant=resource
output.constant=output
state.constant=2

[template]
# Expéditeur des emails à traiter
template1.sender=sender@mail.net
# Vous pouvez définir une expression régulière pour lier des expéditeurs génériques à un template (à partir de la version 3.39.0)
template1.regex=sender\d*@mail.net
# Template à appliquer sur ces emails
template1.path=/opt/canopsis_connectors/email2canopsis/etc/template_1.conf


```

#### Configuration de la connexion RabbitMQ et Redis

Les blocs `[amqp]` et `[redis]` contiennent respectivement la configuration pour la connexion au bus RabbitMQ et au cache Redis.

Il faut donc vérifier que les URL qui y figurent sont les bonnes.

#### Configuration de la Boîte aux Lettres

Le bloc `[mail]` contient la configuration pour la connexion à la boîte aux lettres POP3 dont les emails seront convertis en événements Canopsis.

À partir de la `3.21.0`, l'option `leavemails` permet de conserver sur le serveur les emails lus. Sa valeur est booléenne (`False` ou `True`). Par défaut, l'option vaut `False` et les emails sont supprimés après leur lecture.

Il faut donc vérifier que les différents champs sont bien remplis.

#### Configuration de l'événement envoyé

Le bloc `[event]` contient la configuration de l'événement envoyé à Canopsis.

On peut y définir les différents champs d'un [événement de type check](../../guide-developpement/struct-event.md#event-check-structure).

On peut définir les champs `component`, `resource` et `output` de manière dynamique en faisant appel aux templates. Pour cela, on utilise `param1`, `param2` ou `paramn`.

À partir de la `3.39.0`; le bloc `[event_error]` permet de définir l'événement envoyé en cas [d'erreur de connexion](#gestion-derreur-dans-la-connexion-pop3).

#### Configuration des templates

Le bloc `template` contient la configuration des templates.

Pour la recette, il faut s'assurer que l'adresse depuis laquelle on va envoyer un email se trouve bien dans le bloc (ou qu'elle respecte bien une expression régulière définie, à partir de Canopsis 3.39.0) et que le contenu de l'email correspond bien au template appliqué à cette même adresse email.

Ici, pour

```ini
[template]
# Expéditeur des emails à traiter
template1.sender=sender@mail.net
# Vous pouvez définir une expression régulière pour lier des expéditeurs génériques à un template (à partir de la version 3.39.0)
template1.regex=sender\d*@mail.net
# Template à appliquer sur ces emails
template1.path=/opt/canopsis_connectors/email2canopsis/etc/template_1.conf
```

Il faut envoyer un email depuis l'adresse `sender@mail.net` et son contenu doit correspondre au template `/opt/canopsis_connectors/email2canopsis/etc/template_1.conf`.

A partir de la `3.40.0` On peut assigner plusieurs templates à un expéditeur en fonction du sujet du mail. Pour cela il faut définir une expression régulière pour assigner un sujet a son template.

Exemple :

```ini
[template]
# Expéditeur des emails à traiter
template1.sender=sender@mail.net
# Template à appliquer sur ces emails
template1.path=/opt/canopsis_connectors/email2canopsis/etc/template_1.conf

# Expéditeur des emails à traiter
template2.sender=sender@mail.net
# Vous pouvez définir une  expression régulière pour lier un sujet de mail a un template (3.40.0+)
template2.subject=.*Datacenter.*
# Template à appliquer sur ces emails
template2.path=/opt/canopsis_connectors/email2canopsis/etc/template_2.conf
```

Dans cet exemple tous les mails de `sender@mail.net` qui ont dans le sujet le mot `Datacenter` seront liés au template `/opt/canopsis_connectors/email2canopsis/etc/template_2.conf`.

Le `subject` peut être utilisé avec les deux types de déclarations d'expéditeurs (`sender` ou `regex`).


### Configuration du template d'email

#### Principe du template

Les templates permettent d'analyser un email en suivant quelques règles de transformations simples.

Un fichier template devrait ressembler à quelque chose comme suit :

```ini
[tpl]
component=MAIL_SENDER
resource=MAIL_BODY.line(16).word(2).untilword()
# Depuis la 3.11.0 l'option trim permet de supprimer les espaces sur un champ, ici par exemple avec le champ ressource
resource.trim=left
output=MAIL_BODY.line(6).word(3).untilword()
long_output=MAIL_BODY.line(7).after(le).untilword()
state=MAIL_SUBJECT.line(0).word(1)
state.converter=Mineur>1,Majeur>2,Critique>3,Fin>0
timestamp=MAIL_DATE
timestamp.output=timestamp
```

Les clefs sont divisées en deux parties, séparées par un point :

- à gauche, une *racine* : le nom de la variable à insérer dans l'événement généré
- à droite, une *méthode* : c'est-à-dire l'action de transformation à appliquer

Les actions peuvent être les suivantes :

* *selector* (utilisé par défaut ; implicite) : applique simplement le template à droite et copie la valeur traduite dans l'événement.
* *converter* : remplace une chaîne de caractères par une autre (insensiblement à la casse), les deux étant séparés par le symbole '>'. Plusieurs conversions sont applicables à la suite en les séparant par des virgules. Dans l'exemple ci-dessus, 'Mineur' sera remplacé par 1, 'Majeur' par 2…

A partir de la `3.40.0` *converter* utilise des expressions régulières pour effectuer le remplacement.

Exemple :
```
    state=MAIL_SUBJECT
    state.converter=Mineur \?>1,^Majeur$>2,Critique>3,.*>0
```

On sélectionne dans cet exemple le sujet du mail pour définir la criticité de l’alarme.

- Les mails dont le sujet contient `Mineur ?` auront une criticité de 1. Le caractère `?` est un symbole utilisé dans l’écriture des expressions régulières, comme `*, {, }` etc. Il faut donc le protéger avec un `\`.
- Les mails dont le sujet est strictement `Majeur` auront une criticité de 2. Le caractère `^` défini le début de la chaîne de caractères et `$` la fin. On aurait donc pu définir comme  expression régulière `^Mineur` pour sélectionner les mails dont le sujet commence par `Mineur`. Inversement `Mineur$` pour la sélection des mails dont le sujet se termine par `Mineur`.
- Les mails dont le sujet contient `Critique` auront une criticité de 3.
- L'utilisation de l'expression régulière `.*` permet de définir un comportement par défaut. Les mails qui ne correspondent pas aux cas précédents auront donc une criticité par défaut de 0.

À partir de la `3.11.0`, l'option `trim` retire les espaces à gauche, à droite ou des 2 côtés du bloc de mots. Elle peut être appliquée à n'importe quelle *racine*. Par exemple, si la ressource dans le mail vaut "␣deux mots␣" avec un espace avant et après :  

- `resource.trim=left` donnera "deux mots␣" avec l'espace à gauche supprimé
- `resource.trim=right` donnera "␣deux mots" avec l'espace à droite supprimé
- `resource.trim=both` donnera "deux mots" avec les espaces à gauche et à droite supprimés

!!! info
    A partir de la `3.40.0` les options de `trim` deviennent des opérateurs :  
    
    - `trim_left`  
    - `trim_right`  
    - `trim_both`  

À partir de la `3.39.0`, l'option `print` assigne directement une valeur au champ à partir du template.

- `resource.print=Valeur`

La partie droite décrit les règles de transformations (où a, b et c sont des entiers, et d, e, f des chaînes de caractères) :

- `MAIL_BODY`, `MAIL_DATE`, `MAIL_ID`, `MAIL_SENDER`, et `MAIL_SUBJECT` sont les différentes parties de l'email
- `line(a)` sélectionne une ligne entière numéro a
- `line(a,b)` sélectionne les lignes entières entre a et b
- `line(*)` sélectionne l'ensemble des lignes
- `line(a).word(b)` sélectionne le b-ième mot de la ligne a
- `line(a).word(b).untilword(c)` sélectionne tous les mots entre le b-ième et le c-ième sur la ligne a
- `line(a).word(b).untilword()` sélectionne tous les mots à partir du b-ième jusqu'à la fin de la ligne a
- `line(a).after(d)` sélectionne tous les mots après le mot d
- `line(a).after_incl(d)` sélectionne tous les mots après le mot d, d inclus
- `line(a).after(d).untilword(c)` sélectionne les mots après le mot d et c mots ensuite
- `line(a).after(d).word(b).untilword(c)` sélectionne les mots après le mot d, à partir du b-ième jusqu'au c-ième
- `line(a).before(e)` sélectionne tous les mots avant e
- `line(a).before_incl(e)` sélectionne tous les mots avant e, e inclus
- `line(a).before(e).word(c)` sélectionne tous les mots avant e, mais en commençant au c-ième
- `line(a).replace(e,f)` remplace la chaîne de caractères e par la chaîne f dans la sélection (ici une ligne entière numéro a). Cette opération peut être répétée. (à partir de la `3.40.0`).
- `line(a).remove(e)` supprime la chaîne de caractères e dans la sélection (ici une ligne entière numéro a). Cette opération peut être répétée. (à partir de la `3.40.0`).
- `line(a).lowercase` passe la sélection (ici une ligne entière numéro a) en minuscules. (à partir de la `3.40.0`).
- `line(a).uppercase` passe la sélection (ici une ligne entière numéro a) en majuscules. (à partir de la `3.40.0`).
- `line(a).trim_left` supprime l'espace à gauche de la sélection (à partir de la `3.40.0`).
- `line(a).trim_right` supprime l'espace à droite de la sélection (à partir de la `3.40.0`).
- `line(a).trim_both` supprime les espaces à gauche et à droite de la sélection (à partir de la `3.40.0`).
- `and` permet d'effectuer une concaténation entre deux opérations (à partir de la `3.39.0`).
- `print(word)` permet d'assigner la valeur word dans le champ (à partir de la `3.39.0`).


MAIL\_DATE est automatiquement converti en objet date, inutile d'appliquer une action 'dateformat' dessus.

!!! info
    À partir de Canopsis 3.39.0, il existe la méthode `print` et la règle `print(word)`.  
    Exemple :  

        resource.print = valeur
        # et
        resource = print(valeur)  

    Il faut privilégier la méthode `print` pour l'insertion de valeur statique, et la règle `print(word)` lorsque l'on doit la combiner avec `and` et d'autres règles.

!!! attention
    Les numéros de lignes et de mots commencent à partir de 0, non de 1.

Exemple : La séquence à la 1° ligne située entre les 5° et le 18° mots sont donc sélectionnables avec la ligne `line(0).word(4).untilword(17)`.

Depuis la `3.44.0` il y a la possibilité de faire des templates pour des mails sous forme HTML sans prendre en compte le balisage. Il faut ajouter un `action_template=convert_html2text` au template.

## Dépannage

#### Appliquer des changements de configuration

Pour appliquer un changement (modification de la configuration, ajout de templates, etc.), il faut redémarrer le connecteur.

#### Gestion d'erreur dans la connexion POP3

En cas d'erreur de connexion au serveur mail, le connecteur envoie un événement à Canopsis. Vous pouvez paramétrer cette alerte avec la section `[event_error]` du fichier de configuration.

En cas de connexion normale du connecteur au serveur mail, le connecteur envoie l'événement avec une criticité de 0. Cela permet de fermer d'éventuelles alarmes.
