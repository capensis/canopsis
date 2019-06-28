# Liste des helpers handlebars disponibles dans l'interface Canopsis

## Qu'est-ce qu'un "helper" ?

A différents endroits, dans l'interface de Canopsis, il est possible de personnaliser le texte affiché.

Par example, il est possible de personnaliser le texte affiché pour les tuiles d'un widget de Météo de service.

La configuration de ces affichages se fait, dans les options des widgets, grâce à des templates utilisants le langage **Handlebars** ([documentation du langage Handlebars](https://handlebarsjs.com/)).

Les helpers permettent d'ajouter de la logique dans les templates Handlebars. Le langage en inclut certains par défaut. Par exemple, le helper ```{{#if condition}}{{/if}}``` permet de n'afficher le texte contenu à l'intérieur de ce bloc seulement si la condition est respectée. Pour plus de détails, ainsi que la liste des helpers inclus par défaut, cliquez [ici](https://handlebarsjs.com/builtin_helpers.html).

En plus des helpers inclus par défaut, Handlebars permet d'ajouter des helpers spécifiques. L'interface de Canopsis met à disposition certains helpers supplémentaires.

## Liste des helpers

### Global

#### internal-link

Le helper ```internal-link``` permet d'inclure, dans un template, un lien vers une autre page de l'interface Canopsis.

**Attention**: Ce helper permet spécifiquement d'inclure un lien vers une autre page de l'interface de Canopsis. Pour inclure un lien vers l'extérieur de canopsis, ce helper n'est pas adapté.

- Utilisation:

```
{{ internal-link href="adresse complète de la page de Canopsis souhaitée" text="texte affiché pour le lien" }}
```

L'option ```href``` peut contenir:

- Une addresse complète. Exemple: ```http://10.127.186.12/en/static/canopsis-next/dist/index.html/views/123456```
- Un chemin vers un champ de l'entité. Par exemple: ```entity.infos.weather-link.val```. Cette valeur sera remplacée par la valeur trouvée dans l'entité. Exemple: ```{{ internal-link href=entity.infos.weather-link text="Link" }}```. **Attention**: Pour insérer une variable, il ne faut pas entourer la valeur de guillemets 

En dehors des options ```href``` et ```text```, obligatoires pour un fonctionnement normal, il est possible d'ajouter tous les attributs souhaités.

Par exemple, un attribut style, pour personnaliser l'apparence du lien : ```{{ internal-link href=entity.infos.weather-link text="Link" style="color: 'blue';" }}```.

### Par widget

