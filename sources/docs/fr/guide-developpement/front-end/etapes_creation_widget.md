# Etapes création d'un widget (Cas Pratique: 'Calendrier watcher')

## Ajout du widget dans les constantes

Dans le fichier ```src/constants.js```.

Ajouter le widget à l'objet ```WIDGET_TYPES```
Ici, on ajoute à ```WIDGET_TYPES``` :

```JS 
weatherCalendar: 'Weather Calendar'
```

## Ajout du widget à la modale d'ajout de widget

Dans le fichier ```src/components/modals/view/create-widget.vue```.

La partie ```data``` contient un tableau ```widgetTypes```, contenant les différents types de widgets disponibles depuis cette modale d'ajout de widget.

Il faut ici ajouter notre widget à cette liste:

```JS
{
  title: this.$t('modals.widgetCreation.types.watcher.title'),
  value: WIDGET_TYPES.watcherCalendar,
  icon: 'calendar_today',
},
```

- ```title```: Nom du widget affiché dans la modale. L'appel à la fonction ```this.$t()``` permet la traduction de ce nom ([voir ci-après pour l'ajout des traductions](#traductions))
- ```value```: Nom interne du widget. Il faut ici indiquer le chemin du widget dans l'objet ```WIDGET_TYPES``` cité ci-dessus
- ```icon```: Icône présent dans la modale, à droite du nom du widget. Tout les icônes disponibles sont listés [ici](https://material.io/tools/icons/?style=baseline)


Le widget doit maintenant être listé dans a modale d'ajout de widgets. Au clic sur ce nouveau widget, l'application cherchera à ouvrir la barre latérale de paramètrages de ce widget. Il nous faut donc maintenant créer cette barre latérale

## Création de la barre de paramètres du widget

### Ajout dans les constantes

Dans le fichier ```src/constants.js```.

1 - Ajouter  à l'objet ```SIDE_BARS```:

```JS
watcherCalendarSettings: 'watcher-calendar-settings'
```

Afin de conserver une homogénéité, l'entrée ajoutée ici correspond au nom du widget, suivi du mot ```settings```

2 - Ajouter à l'objet ```SIDE_BARS_BY_WIDGET_TYPES```:

```JS
[WIDGET_TYPES.watcherCalendar]: SIDE_BARS.watcherCalendarSettings
```

Ce champ utilise les objets auxquels nous avons déjà ajouté notre widget. Il nous servira à faire le lien entre le type de widget, et le type de barre de paramètres à utiliser.

### Création du squelette

Tout d'abord, il faut créer le squelette du composant de notre barre de paramètres.

Dans le dossier ```src/components/side-bars/settings/widgets```, créez un nouveau fichier.

Pour notre exemple, nous appellerons ce fichier ```watcher-calendar.vue```.

Nous compléterons ce fichier par la suite, voici les éléments que nous pouvons cependant déjà commencer à compléter :

- Template: Afin de rester cohérent avec le style des autres barres de paramètres de widgets, nous pouvons déjà ajouter la base du template:

```html
<template lang="pug">
  div
    v-list.pt-0(expand)
    v-btn.primary(@click="submit") {{ $t('common.save') }}
</template>
```

- Script: Les composants de barres de paramètres de widgets sont tous nommés. Nous pouvons donc dès maintenant ajouter cet élément à notre script

```html
<script>
import { cloneDeep } from 'lodash';

import { SIDE_BARS } from '@/constants';

import widgetSettingsMixin from '@/mixins/widget/settings';

export default {
  name: SIDE_BARS.watcherCalendarSettings,
  mixins: [widgetSettingsMixin],
  data() {
    const { widget, rowId } = this.config;

    return {
      settings: {
        rowId,
        widget: cloneDeep(widget),
      },
    };
  },
};
</script>
```

Nous récupérons ici, depuis les constantes que nous avons ajoutés ci-dessus, le nom à utiliser pour ce composant.

### Ajout dans la liste des barres de paramètres

Dans le fichier ```src/components/side-bars/settings/index.vue```.

Il faut ici :

1. Importer le fichier du composant de la barre latérale créé ci-dessus

```JS
import WatcherCalendarSettings from './settings/widgets/watcher-calendar.vue';
```

2. Ajouter ce composant à la liste des composants

```JS
  components: {
    ...,
    WatcherCalendarSettings,
  },
```

Enfin, afin de pouvoir afficher le titre de la barre latérale, il faut l'ajouter dans le fichier ```src/components/side-bars/settings/side-bar-wrapper.vue```. Dans la propriété calculée ```title()``` :

```JS
const TITLES_MAP = {
  ...,
  [SIDE_BARS.watcherCalendarSettings]: this.$t('settings.titles.watcherCalendarSettings'),
};
```
Comme préciser ci-dessus, la fonction ```this.$t()``` permet la traduction ([voir ci-dessous pour l'ajout des traductions](#traductions))

## Génération du widget et ajout de paramètres

A cette étape, il nous est possible de sélectionner notre widget dans la modale d'ajout de widget. Au clic sur ce nouveau widget, une barre de paramètres s'ouvre. Celle-ci est néanmoins toujours vide.

Il nous faut ici ajouter les mécanismes de génération du widget. Cela nous fera passer par l'étape d'ajout des paramètres communs à tous les widgets, puis l'ajout d'éventuels paramètres propres à notre nouveau widget.

### Génération du widget

Cette étape se passe dans le fichier ```src/helpers/entities.js```. Ce fichier exporte une fonction ```generateWidgetByType```. Cette fonction à pour but de générer l'objet représentant le widget souhaité.

Cette fonction démarre par la déclaration de la constante :

```JS
const widget = {
  type,
  _id: uuid(`widget_${type}`),
  title: '',
  parameters: {},
  size: {
    sm: 3,
    md: 3,
    lg: 3,
  },
};
```

Cela représente en réalité le socle commun entre tous les types de widgets:

- ```type```: Le type du widget
- ```_id```: Identifiant unique du widget
- ```title```: Titre (optionnel) du widget. Il nous faudra ajouter ce paramètre à la barre latérale de paramètrage du widget
- ```parameters```: Objet qui contiendra les paramètres spécifiques du widget
- ```size```: Objet correspondant à la taille du widget dans la vue. Il nous faudra ajouter ce paramètre à la barre latérale de paramètrage du widget

Nous pouvons voir ici que 2 paramètres communs à tous les widgets seront déjà à ajouter à la barre de paramétrage créée ci -dessus. Nous y reviendrons plus tard dans cette partie.

Ce fichier contient ensuite les différents cas, pour chaque type de widget. Nous pouvons ajouter ici un nouveau cas, pour notre nouveau type de widget:

```JS
case WIDGET_TYPES.context:
  specialParameters = {
  };
  break;
```

Nous n'avons, pour le moment, pas encore identifié de paramètres spécifiques pour notre widget d'exemple. Ce qu'il faut retenir concernant ce fichier est qu'il nous faudra y ajouter tous les paramètres que l'on souhaite ajouter au widget par défaut, lors de sa création.

### Ajout des paramètres à la barre latérale

Nous avons vu ci-dessus que 2 paramètres sont à ajouter à la barre latérale :

- Taille du widget
- Titre

Les différents composants correspondants aux champs de paramétrages sont contenus dans le dossier ```src/components/side-bars/settings/widgets/fields```

Ajoutons les 2 champs souhaités. Dans le fichier de notre composant de barre latérale ```src/components/side-bars/settings/widgets/watcher-calendar.vue```:

1. Importation des composants des champs:

```JS
import FieldRowGridSize from './fields/common/row-grid-size.vue';
import FieldTitle from './fields/common/title.vue';
```

2. Ajout des composants importés à la liste des composants utilisés:

```JS
components: {
  FieldRowGridSize,
  FieldTitle,
},
```

3. Insertion des composants dans le template

```JS
<template lang="pug">
  div
    v-list.pt-0(expand)
      field-row-grid-size
      v-divider
      field-title
      v-divider
    v-btn.primary(@click="submit") {{ $t('common.save') }}
</template>
```

4. Ajout des propriétés nécessaires aux différents champs

Chaque composant de champ de paramétrage (ici taille du widget, et titre du widget) ajoutés ci-dessus à besoin d'unc ertain nombre de propriétés pour fonctionner correctement.

La liste de ces propriétés peut-être retrouvé dans le composant du champs lui-même.

Par exemple, pour le champs de titre, si nous allons voir le composant correspondant (fichier ```src/components/side-bars/settings/widgets/fields/common/title.vue```), nous pouvons constanter, dans le champ ```props```, que ce composant attends 2 propriétés: ```value``` (correspondant à la valeur du titre du widget, ```""``` par défaut, et ```title``` (correspondant au titre du champ lui-même, affiché dans la barre de paramétrage).

Ajoutons ces propriétés dans notre exemple:

```JS
<template lang="pug">
  div
    v-list.pt-0(expand)
      field-row-grid-size(
      :rowId.sync="settings.rowId",
      :size.sync="settings.widget.size",
      :availableRows="availableRows",
      @createRow="createRow"
      )
      v-divider
      field-title(v-model="settings.widget.title", :title="$t('common.title')")
      v-divider
    v-btn.primary(@click="submit") {{ $t('common.save') }}
</template>
```

## Le widget

Après avoir vu les étapes permettant de générer le widget, puis de lui ajouter les premiers paramètres communs à tous les widgets, il est temps de créer le composant du widget lui-même.

Selon le type le widget, l'emplacement du fichier peut varier. Les widgets sont contenus dans le dossier ```src/components/other```. 
Notre exemple porte sur un widget se rapprochant d'un widget de statistique, nous allons donc créer le fichier dans le répertoire ```src/components/other/stats```

Ce fichier contient le corps de notre widget. Rien n'est obligatoire quant au contenu de ce composant. Cependant, tous les composants de widgets recoivent comme propriété l'objet représentant le widget lui-même.
Nous pouvons donc ajouté cette props dans notre nouveau composant:

```JS
props: {
  widget: {
    type: Object,
    required: true,
  },
},
```

Cette propriété nous permettra, entre autre, d'accèder aux paramètres du widget (via ```widget.parameters```)

## Inscription du widget dans les vues

La dernière étape consiste à indiquer au composant gérant les vues quel composant il doit associer lorsque celui-ci doit afficher un widget du type que nous venons de créer (dans notre exemple, quel composant vue doit-il afficher lorsqu'un widget de type ```watcherCalendar``` lui ai demandé).

Pour cela, dans le fichier ```src/components/other/view/view-tab-rows.vue``` :

1. Importer le composant du widget

```JS
import WatcherCalendar from '@/components/other/stats/watcher-calendar.vue';
```

2. L'ajouter à la liste des composants utilisés

```JS
components: {
  ...,
  WatcherCalendar,
},
```

3. L'ajouter à la liste des composants ```widgetComponentsMap```, qui se trouve dans les ```data```

```JS
data() {
  return {
    widgetsComponentsMap: {
      ...,
      [WIDGET_TYPES.watcherCalendar]: 'watcher-calendar',
    },
  };
},
```

La base de notre nouveau widget est prête. Il ne reste plus qu'à y ajouter les fonctionnalités souhaitées, ainsi que les paramètres nécessaires.

## Traductions

L'interface de Canopsis est disponible en deux langues, français et anglais. Afin de permettre à l'utilisateur de passer d'une langue à l'autre, il nous faut ajouter les différentes traductions nécessaires à notre widget.

Ces traductions sont contenues dans les fichiers ```en.json``` et ```fr.json```, du dossier ```src/i18n/messages```.
Ces deux fichiers, correspondants respectivement aux traduction anglaise et française, sont des fichiers JSON. Chaque entrée de ces fichiers correspond à la traduction d'une chaîe de caractère utilisée dans l'interface. Une même entrée de traduction peut être utilisée à différent endroit de l'interface.

Comme nous l'avons vu dans les étapes précédentes, les traductions sont utilisées depuis nos composants grâce à la fonction ```$t()```. En réalité, cette fonction prend en argument le chemin vers la traduction souhaitée, dans le fichier JSON de la langue utilisée.

Par exemple, l'appel ```$t('common.submit')``` ira chercher dans le fichier JSON de la langue sélectionnée (les fichiers ```en.json``` et ```fr.json``` décrits ci-dessus), l'entrée ```submit``` du sous-objet ```common```. Puis remplacera la chaîne de caractère par la valeur de traduction trouvée.


```JSON
# fr.json
{
  "common": {
    "submit": "Soumettre" 
  }
}
```

```JSON
# en.json
{
  "common": {
    "submit": "Submit"
  }
}
```

Pour ajouter/utiliser une traduction dans notre composant, nous suivons les étapes suivantes:

1) Vérifier que la traduction souhaitée n'est pas déjà présente dans la liste de traductions de l'interface. Une recherche dans le fichier ```fr.json``` ou ```en.json``` nous permet de rapidement déterminer si la traduction existe ou non.
2) Si la traduction n'est pas présente, l'ajouter dans les fichiers JSON ```fr.json``` et ```en.json```. Il convient, pour chaque ajout, de juger de l'endroit le plus adéquat, à l'intérieur du fichier JSON.
3) Dans notre composant, appeler la fonction de traduction ```$t()```, en lui passant en argument le chemin, dans les fichiers JSON de traduction, vers la traduction souhaitée.