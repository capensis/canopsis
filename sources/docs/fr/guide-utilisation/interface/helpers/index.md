# Helpers Handlebars disponibles dans l'interface Canopsis

## Qu'est-ce qu'un « helper Handlebars » ?

À différents endroits dans l'interface de Canopsis, il est possible de personnaliser le texte affiché (par exemple, les tuiles de Météo de service).

La configuration de ces affichages se fait dans les options des widgets, grâce à des templates utilisant le langage [Handlebars](https://handlebarsjs.com/).

Le langage Handlebars propose, par défaut, un ensemble de *helpers* permettant de personnaliser les templates en fonction de diverses conditions. Par exemple, le helper `{{#if condition}}texte{{/if}}` permet d'afficher le texte contenu à l'intérieur de ce bloc seulement si la `condition` indiquée est vraie. Voyez [la documentation des helpers Handlebars officiels](https://handlebarsjs.com/guide/builtin-helpers.html) afin d'en savoir plus.

En plus de ces helpers officiels, Canopsis met à disposition certains helpers supplémentaires qui lui sont propres. La documentation suivante décrit leur fonctionnement.

## Liste des helpers Handlebars propres à Canopsis

L'ensemble des helpers suivants sont génériques, et peuvent être utilisés à chaque endroit où la syntaxe Handlebars est permise.

### Helper `compare`

!!! note
    Disponible depuis Canopsis 3.33.0.

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

### Helper `duration`

!!! note
    Disponible depuis Canopsis 3.38.0.

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

!!! note
    Disponible depuis Canopsis 3.33.0.

Le helper `internal-link` permet d'afficher un lien cliquable vers une autre page interne à Canopsis.

!!! attention
    Ce helper permet uniquement d'inclure un lien vers une autre page de l'interface de Canopsis actuellement utilisée. Il ne peut pas être utilisé pour construire des liens vers une page externe.

Utilisation générique :

```handlebars
{{internal-link href=lien-canopsis text=description-lien style=instructions-css}}
```

Ce helper accepte les attributs suivants :

*  `href` (obligatoire). La page interne à Canopsis vers laquelle on souhaite créer un lien.
    *  `lien-canopsis` peut être remplacé par une chaîne de caractères contenant une adresse complète, telle que `"http://10.127.186.12/en/static/canopsis-next/dist/index.html/views/123456"`.
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

!!! note
    Disponible depuis Canopsis 3.33.0.

Le helper `timestamp` transforme un *timestamp* Unix (UTC) en une date au format numérique `jour/mois/année heures:minutes:secondes` (adapté au fuseau horaire du navigateur de l'utilisateur). Si la date correspond au jour actuel, seule la partie `heures:minutes:secondes` sera affichée.

Utilisation générique :

```handlebars
{{timestamp nombre-de-secondes}}
```

Ce helper attend un unique paramètre :

*  `nombre-de-secondes` (obligatoire). Une variable, ou un entier positif correspondant à un [timestamp Unix](https://fr.wikipedia.org/wiki/Heure_Unix) réglé sur UTC. Si la variable est vide, le helper n'affichera rien.

#### Exemple d'utilisation du helper `timestamp`

Afficher le timestamp Unix de `544089600` secondes (correspondant au 30 mars 1987, à 8 heures précises du matin, UTC) sous la forme « 30/03/1987 10:00:00 » (avec un navigateur sur le fuseau `Europe/Paris`, en heure d'hiver).

```handlebars
{{timestamp 544089600}}
```

### Helper `request`

!!! note
    Disponible depuis Canopsis 3.38.0.

Le helper `request` permet d'effectuer des requêtes vers des API REST (internes ou externes) et de manipuler les résultats renvoyés.

Si la requête échoue, un message d'erreur sera affiché à la place du résultat attendu.

!!! attention
    L'API est appelée depuis le navigateur du client, et non pas depuis le serveur hébergeant Canopsis.

Ce helper accepte les attributs suivants :

*  `method` (optionnel, `GET` par défaut). La méthode HTTP à utiliser pour effectuer la requête.
*  `url` (obligatoire). URL de l'API JSON à interroger.
*  `headers` (optionnel). Entêtes HTTP, au format JSON (`{"Nom-Entete": "valeur"}`), à intégrer lors de l'envoi de la requête.
    *  **Note :** tout entête envoyé doit apparaître dans la directive `Access-Control-Allow-Headers` du [serveur Nginx intégré à Canopsis](../../../guide-administration/administration-avancee/reverse-proxy.md#configuration-de-nginx).
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

Afficher une liste des noms d'utilisateurs inscrits à GitHub :

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

## Helper `state`

!!! note
    Disponible depuis Canopsis 3.38.0.

Le helper `state` permet de transformer [l'état de criticité](../../vocabulaire/index.md#etat) d'une alarme (sous sa forme numérique 0, 1, 2 ou 3) en une pastille de couleur associée à cette criticité, telle qu'elle peut apparaître dans un Bac à alarmes. Le texte de ces pastilles est actuellement toujours affiché en anglais.

Utilisation générique :

```handlebars
{{state numero-etat-criticite}}
```

Ce helper attend un unique paramètre :

*  `numero-etat-criticite` (obligatoire). Une variable, ou un entier représentant l'état de criticité d'une alarme, tel que défini par Canopsis (0, 1, 2 ou 3). Un état de criticité différent de ces valeurs provoquera l'affichage d'une pastille d'erreur.

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

Afficher une pastille « Invalid val » (l'état de criticité étant invalide) :

```handlebars
{{state 9}}
```
