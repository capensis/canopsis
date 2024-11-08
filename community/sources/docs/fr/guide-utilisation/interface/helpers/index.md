# Helpers Handlebars disponibles dans l'interface Canopsis

## Qu'est-ce qu'un « helper Handlebars » ?

À différents endroits dans l'interface de Canopsis, il est possible de personnaliser le texte affiché (par exemple, les tuiles de Météo de service).

La configuration de ces affichages se fait dans les options des widgets, grâce à des templates utilisant le langage [Handlebars](https://handlebarsjs.com/).

Le langage Handlebars propose, par défaut, un ensemble de *helpers* permettant de personnaliser les templates en fonction de diverses conditions. Par exemple, le helper `{{#if condition}}texte{{/if}}` permet d'afficher le texte contenu à l'intérieur de ce bloc seulement si la `condition` indiquée est vraie. Voyez [la documentation des helpers Handlebars officiels](https://handlebarsjs.com/guide/builtin-helpers.html) afin d'en savoir plus.

En plus de ces helpers officiels, Canopsis met à disposition certains helpers supplémentaires qui lui sont propres. La documentation suivante décrit leur fonctionnement.

## Liste des helpers Handlebars propres à Canopsis

L'ensemble des helpers suivants sont génériques, et peuvent être utilisés à chaque endroit où la syntaxe Handlebars est permise.

### Helper `compare`

Le helper `compare` permet de n'afficher un contenu que si une comparaison entre deux objets est vraie.

Utilisation générique :

```handlebars
{{#compare objet1 operateur objet2 options}}Texte à afficher si la comparaison est vraie{{/compare}}
```

Ce helper accepte quatre paramètres, dans l'ordre suivant :

*  `objet1` (obligatoire). Une variable, un entier ou une chaîne de caractères à comparer.
*  `operateur` (obligatoire). Un opérateur de comparaison. Il s'agit de l'un des éléments suivants :
    *  `'<'` : si `objet1` est strictement inférieur à `objet2`.
    *  `'>'` : si `objet1` est strictement supérieur à `objet2`.
    *  `'<='` : si `objet1` est inférieur ou égal à `objet2`.
    *  `'>='` : si `objet1` est supérieur ou égal à `objet2`.
    *  `'=='` : si `objet1` est égal à `objet2`.
    *  `'!='` : si `objet1` est différent de `objet2`.
    *  `'==='` : si `objet1` est *strictement* égal à `objet2` (cf. [différence entre une égalité stricte et une égalité faible en JavaScript](https://developer.mozilla.org/fr/docs/Web/JavaScript/Les_différents_tests_d_égalité)).
    *  `'!=='` : si `objet1` est *strictement* différent de `objet2`.
    *  `'regex'` : si l'expression régulière contenue dans `objet2` s'applique au contenu de `objet1`.
*  `objet2` (obligatoire). Une autre variable, un autre entier ou une autre chaîne de caractères.
*  `options` (optionnel). Permet de définir des options propres à un opérateur de comparaison.
    *  Seul l'opérateur `regex` est concerné par ce paramètre optionnel, pour l'instant. Il s'agit d'une chaîne regroupant l'ensemble des [flags de regex](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Regular_Expressions#Advanced_searching_with_flags_2) à appliquer lors de son évaluation.

Les opérateurs de comparaison `==` et `===` supportent le flag `i` dans leur évaluation.  

#### Exemples d'utilisation du helper `compare`

Afficher `Test` uniquement si le nombre d'`essais` est supérieur à 10 :

```handlebars
{{#compare essais '>' 10}}Test{{/compare}}
```

Afficher `Test` uniquement si le contenu de `variable` est égal à « important » (exactement tel qu'il est écrit en lettres minuscules) :

```handlebars
{{#compare variable '==' 'important'}}Test{{/compare}}
```

Afficher `Test` uniquement si une chaîne se *termine* par « motif », en ignorant la casse minuscules/majuscules :

```handlebars
{{#compare 'Cette phrase se termine par MOTIF' 'regex' 'motif$' flags='i'}}Test{{/compare}}
```

Afficher `Test` uniquement si le contenu de `variable` est égal à `Cette phrase insensible à la CASSE` sans tenir compte de la casse :

```handlebars
{{#compare variable '==' 'Cette phrase insensible à la CASSE'  flags='i'}}Test{{/compare}}
```

### Helper `duration`

Le helper `duration` transforme un nombre de secondes en une durée facilement compréhensible par un humain. Son affichage sera adapté en fonction de la langue définie dans le profil de l'utilisateur.

Utilisation générique :

```handlebars
{{duration secondes}}
```

Ce helper attend un unique paramètre :

*  `secondes` (obligatoire). Une variable, ou un entier positif représentant un nombre de secondes. Si la variable est vide, le helper n'affichera rien.

#### Exemples d'utilisation du helper `duration`

Transformer 120 secondes en « 2 mins » :

```handlebars
{{duration 120}}
```

Transformer 133742 secondes en « 1 jour, 13 hrs, 9 mins, 2 secs » :

```handlebars
{{duration 133742}}
```

### Helper `internal-link`

Le helper `internal-link` permet d'afficher un lien cliquable vers une autre page interne à Canopsis.

!!! attention
    Ce helper permet uniquement d'inclure un lien vers une autre page de l'interface de Canopsis actuellement utilisée. Il ne peut pas être utilisé pour construire des liens vers une page externe.

Utilisation générique :

```handlebars
{{internal-link href=lien-canopsis text=description-lien style=instructions-css}}
```

Ce helper accepte les attributs suivants :

*  `href` (obligatoire). La page interne à Canopsis vers laquelle on souhaite créer un lien.
    *  `lien-canopsis` peut être remplacé par une chaîne de caractères contenant une adresse complète, telle que `"http://10.127.186.12/views/123456"`.
    *  `lien-canopsis` peut aussi être une variable correspondant à un champ de l'entité. Par exemple : `entity.infos.weather-link.value`. Cette valeur sera remplacée par la valeur trouvée dans l'entité.
    *  **Attention :** les guillemets sont obligatoires dans le cas d'une chaîne de caractères, et ne doivent pas apparaître dans le cas d'une variable.
*  `text` (obligatoire). Le texte qui sera affiché pour construire le lien cliquable.
    *  `description-lien` suit, ici aussi, la syntaxe de `variable` et de `"chaîne de caractères"` décrite au point précédent.
*  `style` (optionnel). Des instructions de style CSS.
    *  `instructions-css` peut être une chaîne de caractères telle que `"color: blue;"` qui stylise le lien cliquable en bleu.

#### Exemples d'utilisation du helper `compare`

Créer un lien (en gras) redirigeant vers la page de météo liée à une entité :

```handlebars
{{internal-link href=entity.infos.weather-link.value text="Aller sur la météo" style="font-weight: bold;"}}
```

### Helper `timestamp`

Le helper `timestamp` transforme un *timestamp* Unix (UTC) en une date au format numérique `jour/mois/année heures:minutes:secondes` (adapté au fuseau horaire du navigateur de l'utilisateur). Si la date correspond au jour actuel, seule la partie `heures:minutes:secondes` sera affichée.

Utilisation générique :

```handlebars
{{timestamp nombre-de-secondes}}
```

Ce helper attend un unique paramètre :

*  `nombre-de-secondes` (obligatoire). Une variable, ou un entier positif correspondant à un [timestamp Unix](https://fr.wikipedia.org/wiki/Heure_Unix) réglé sur UTC. Si la variable est vide, le helper n'affichera rien.
*  `format` (optionnel). Il s'agit du format [moment](https://momentjscom.readthedocs.io/en/latest/moment/04-displaying/01-format/). De plus, le format `long` est également pris en charge. Il s'agit d'afficher la date au format long `30/03/1987 10:00:00` même si le timestamp correspond à la date du jour.

#### Exemple d'utilisation du helper `timestamp`

Afficher le timestamp Unix de `544089600` secondes (correspondant au 30 mars 1987, à 8 heures précises du matin, UTC) sous la forme « 30/03/1987 10:00:00 » (avec un navigateur sur le fuseau `Europe/Paris`, en heure d'hiver).

```handlebars
{{timestamp 544089600}}
{{timestamp 544089600 format='long'}}
```

Afficher la date à partir du timestamp lorsque la date est aujourd'hui (07:07:17)

```handlebars
{{timestamp 1673932037}}
```

Afficher la date à partir du timestamp sur un format long (17/01/2023 07:07:17)

```handlebars
{{timestamp 1673932037 format='long'}}
```

Afficher la date à partir du timestamp sur un format personnalisé (January 17th 2023, 07:07:17 am)

```handlebars
{{timestamp 1673932037 format='MMMM Do YYYY, h:mm:ss a'}}
```

### Helper `request`

Le helper `request` permet d'effectuer des requêtes vers des API REST (internes ou externes) et de manipuler les résultats renvoyés.

Si la requête échoue, un message d'erreur sera affiché à la place du résultat attendu.

!!! attention
    L'API est appelée depuis le navigateur du client, et non pas depuis le serveur hébergeant Canopsis.

Ce helper accepte les attributs suivants :

*  `method` (optionnel, `GET` par défaut). La méthode HTTP à utiliser pour exécuter la requête.
*  `url` (obligatoire). URL de l'API JSON à interroger.
*  `headers` (optionnel). Entêtes HTTP, au format JSON (`{"Nom-Entete": "valeur"}`), à intégrer lors de l'envoi de la requête.
    *  **Note :** tout entête envoyé doit apparaître dans la directive `Access-Control-Allow-Headers` du [serveur Nginx intégré à Canopsis](../../../guide-administration/administration-avancee/configuration-composants/reverse-proxy-nginx.md#configuration-de-nginx).
*  `data` (optionnel). Payload d'une requête POST, au format JSON (`{"Nom": "valeur"}`).
*  `username` (optionnel). Utilisateur pour l'authentification basique.
*  `password` (optionnel). Mot de passe pour l'authentification basique.
    *  **Attention :** la requête étant exécutée par le navigateur client, ces identifiants peuvent être interceptés par un utilisateur.
*  `variable` (optionnel). Nom de la variable à utiliser en interne afin de stocker le résultat de la requête.
*  `path` (optionnel). Chemin JSON des champs à récupérer lors de la requête.

#### Exemples d'utilisation du helper `request`

Afficher une liste d'alarmes depuis une instance de Canopsis, en interrogeant la route `get-alarms`, et en itérant sur les champs `alarms.d` et `alarms.v.state.val` de cette API :

```handlebars
<ul>
{{#request
  url="http://localhost:8082/alerts/get-alarms?opened=true&resolved=false&skip=0&limit=10"
  variable="alarms"
  path="data[0].alarms"
  username="root"
  password="root"}}
    {{#each alarms}}
      <li>État de l'alarme {{d}} : {{state v.state.val}}</li>
    {{/each}}
{{/request}}
</ul>
```

Afficher une liste des identifiants d'utilisateurs inscrits à GitHub :

```handlebars
<ul>
{{#request
  url="https://api.github.com/users"
  variable="users"}}
     {{#each users}}
       <li>{{login}}</li>
     {{/each}}
{{/request}}
</ul>
```

Ajouter un `todo` sur le site `https://jsonplaceholder.typicode.com`

```handlebars
<ul>
{{#request 
  method="post" 
  url="https://jsonplaceholder.typicode.com/todos" 
  variable="post" 
  headers='{ "Content-Type": "application/json" }' 
  data='{ "userId": "1", "title": "TEST321", "completed": false }'}}
    {{#each post}}
        <li><strong>{{@key}}</strong>: {{this}}</li>
    {{/each}}
{{/request}}
</ul>
```

## Helper `state`

Le helper `state` permet de transformer [la criticité](../../vocabulaire/index.md#criticite) d'une alarme (sous sa forme numérique 0, 1, 2 ou 3) en une pastille de couleur associée à cette criticité, telle qu'elle peut apparaître dans un Bac à alarmes. Le texte de ces pastilles est actuellement toujours affiché en anglais.

Utilisation générique :

```handlebars
{{state numero-criticite}}
```

Ce helper attend un unique paramètre :

*  `numero-criticite` (obligatoire). Une variable, ou un entier représentant la criticité d'une alarme, telle que définie par Canopsis (0, 1, 2 ou 3). Une criticité différente de ces valeurs provoquera l'affichage d'une pastille d'erreur.

#### Exemples d'utilisation du helper `state`

Afficher une pastille « ok » verte :

```handlebars
{{state 0}}
```

Afficher une pastille « minor » jaune :

```handlebars
{{state 1}}
```

Afficher une pastille « major » orange :

```handlebars
{{state 2}}
```

Afficher une pastille « critical » rouge :

```handlebars
{{state 3}}
```

Afficher une pastille « Invalid val » (la criticité étant invalide) :

```handlebars
{{state 9}}
```


## Helpers mathématiques basiques

4 helpers sont disponibles pour les opérations mathématiques de base.  

1. sum : renvoie la somme des nombres passés en paramètre
2. minus : renvoie la différence entre 2 nombres
3. mul : renvoie le produit de 2 nombres
4. divide : renvoie le résultat de la division entre 2 nombres

### Helper `sum`

```handlebars
{{sum 1 2 3}}
```

Ce helper attend en paramètre un ensemble de nombres et renvoie leur somme

#### Exemple d'utilisation du helper `sum`

Afficher la somme de 1, 2, et 3 :

```handlebars
{{sum 1 2 3}}
```

### Helper `minus`

```handlebars
{{minus 10 1}}
```

Ce helper attend en paramètre 2 nombres et renvoie leur différence

#### Exemple d'utilisation du helper `minus`

Afficher la différence entre 10 et 1 :

```handlebars
{{minus 10 1}}
```

### Helper `mul`

```handlebars
{{mul 5 6}}
```

Ce helper attend en paramètre 2 nombres et renvoie leur produit

#### Exemple d'utilisation du helper `mul`

Afficher le produit entre 5 et 6 :

```handlebars
{{mul 5 6}}
```

### Helper `divide`

```handlebars
{{divde 10 2}}
```

Ce helper attend en paramètre 2 nombres et renvoie le résultat de leur division

#### Exemple d'utilisation du helper `divide`

Afficher la division 10 par 2 :

```handlebars
{{divide 10 2}}
```

!!! Note
    Une division par 0 affichera `infinity`

## Helpers chaînes de caractères

5 helpers sont disponibles pour réaliser des opérations sur des chaînes de caractères.

1. concat : concatène des chaînes de caractères
1. lowercase : convertit une chaîne de caractères en minuscule
1. uppercase : convertit une chaîne de caractères en majuscule
3. capitalize : ajoute une majuscule en début de chaîne de caractères
4. capitalize-all : ajoute une majuscule sur la première lettre de tous les mots d'une chaîne de caractères

### Helper `concat`

Le helper `concat` permet de concaténer des chaînes de caractères.  

```handlebars
{{concat "chaine1" "chaine2" "chaineN"}}
```

Ce helper attend en paramètre un ensemble de chaînes de caractères à concaténer.  
Ces chaînes peuvent être également des variables de l'alarme ou de l'entité.  

#### Exemples d'utilisation du helper `concat`

Afficher la concaténation de "chaine1" et "chaine2" :

```handlebars
{{concat "chaine1" "chaine2"}}
```

Afficher la concaténation de "http://wiki.local/?co=" et du composant d'une alarme :

```handlebars
{{concat "http://wiki.local/?co=" alarm.v.component}}
```

Utiliser la concaténation dans le helper [request](#helper-request) pour bâtir une URL dynamique :

```handlebars
{{#request
  url=(concat "https://wiki.local/?co=" alarm.v.component)
  variable="items"}}
     {{#each items}}
       <li>{{name}}</li>
     {{/each}}
{{/request}}
```

### Helper `lowercase`

```handlebars
{{lowercase "chaine1"}}
```

Ce helper attend en paramètre une chaîne de caractères et la renvoie en minuscule.

#### Exemple d'utilisation du helper `lowercase`

Afficher la chaîne "CHAINE" en minuscule :

```handlebars
{{lowercase "CHAINE"}}
```

### Helper `uppercase`

```handlebars
{{uppercase "chaine1"}}
```

Ce helper attend en paramètre une chaîne de caractères et la renvoie en majuscule.

#### Exemple d'utilisation du helper `uppercase`

Afficher la chaîne "chaine" en majuscule :

```handlebars
{{uppercase "chaine"}}
```

### Helper `capitalize`

```handlebars
{{capitalize "chaine1"}}
```

Ce helper attend en paramètre une chaîne de caractères et transforme la première lettre en majuscule.

#### Exemple d'utilisation du helper `capitalize`

Ajouter une majuscule sur la première lettre de la chaîne "chaine" :

```handlebars
{{capitalize "chaine"}}
```

### Helper `capitalize-all`

```handlebars
{{capitalize-all "mot1 mot2 mot3"}}
```

Ce helper attend en paramètre une chaîne de caractères et transforme la première lettre de chaque mot en majuscule.

#### Exemple d'utilisation du helper `capitalize-all`

Ajouter une majuscule sur la première lettre de chaque mot de la chaîne "mot1 mot2 mot3" :

```handlebars
{{capitalize-all "mot1 mot2 mot3"}}
```


### Helper `replace`

```handlebars
{{replace 'Mot1 Mot2 Mot3' '(Mot1) (Mot2) (Mot3)' '$3 $2 $1' flags='g'}}
```

Ce helper attend :

- en premier paramètre une chaîne de caractères
- en deuxième paramètre une chaîne de type `capture group` permettant de matcher et d'extraire des patterns de la chaîne initiale avec des patenthèses (`()` )
- en troisième paramètre la nouvelle chaîne constituée de variables récupérées dans les groupes provenant du `capture group` précédent
- en dernier paramètre les `flags` comme ceux que l'on peut par exemple trouver dans `sed` ( ici le `g`  permet d'effectuer le traitement plusieurs fois dans une même ligne : sans cela la substitution s'arrêtera à la première occurence trouvée )



#### Exemple d'utilisation du helper `replace`

Changer l'ordre des mots dans une chaine de caractères :

```handlebars
{{replace 'Ubuntu Debian Linux Fedora' '(Ubuntu) (Debian) (Linux)' '$3 $2 $1' flags='g'}}
```

Donnera la chaine finale : `'Linux Debian Ubuntu Fedora'`

### Helper `tags`

```handlebars
{{tags}}
```

Ce helper permet d'afficher les `tags` d'une alarme sous forme de badge. Il n'attend pas de paramètre.  

### Helper `links`

```handlebars
{{links}}
```

Ce helper permet d'afficher les `liens` d'une alarme ou d'une entité.

### Helper `copy`

```handlebars
{{#copy 'Valeur à copier'}}Label{{/copy}}
```

Ce helper permet de copier le contenu de sa valeur dans le clipboard pour utilsiation ultérieure.

#### Exemple d'utilisation du helper `copy`

Copier dans le presse papier la valeur de l'identifiant d'une alarme :

{{#copy alarm.v.display_name}}<button>Cliquer pour copier</button>{{/copy}}

Copier dans le presse papier la structure d'un événement à partir d'une alarme :

```handlebars
{{#copy (concat '{
  "connector" : "' alarm.v.connector '",
  "connector_name" : "' alarm.v.connector_name '",
  "component" : "' alarm.v.component '",
  "resource" : "' alarm.v.resource '",
  "source_type" : "resource",
  "event_type" : "check",
  "output" : "' alarm.v.output '",
  "state" : ' alarm.v.state.val '
}') }}Cliquez pour copier la base de l'événement d'origine{{/copy}}
```

### Helper `json`

```handlebars
{{json alarm.v 'display_name'}}
```

Ce helper permet de renvoyer des informations au format json

#### Exemple d'utilisation du helper `json`

Renvoyer au format json une alarme complète :

```handlebars
{{json alarm.v }}
```

