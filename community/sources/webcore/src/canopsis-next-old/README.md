# UI environment setup
## Pre requirements
* **Node.js** version `^14.18.1` ([Node.js releases](https://nodejs.org/en/blog/release/))
* **Yarn** version `^1.22.17` ([Yarn releases](https://github.com/yarnpkg/yarn/releases))

## Steps to run UI in development mode

Install all dependencies:
```bash
yarn install
```

Configure .env.local file:
```bash
cp .env .env.local

# You should set into VUE_APP_API_HOST correct path to API server. The server must be in HTTPS.
# Example https://localhost:8082/backend

# You can do it with command
sed -i "1s/.*/VUE_APP_API_HOST=https:\/\/localhost:8082\/backend/" .env.local

# Or manually
nano .env.local
```

Run the development server in watch mode:
```bash
yarn serve

# If you want to run server on specific host or port you can use special arguments
yarn serve --port 8089 --host 0.0.0.0
```

At this point you can go with a browser to the URL indicated by the last command and use the Canopsis UI.

# Test the project :
* Checkout on the branch you want : ` git checkout branch-name`
* Install dependencies : ` yarn install `
* Open Google Chrome with the tag :  `google-chrome --disable-web-security --user-data-dir`
* Open your Canopsis in this chrome, then login, to have the back-end data
* Open the link provided by yarn in chrome

# How can I create a new widget type ?
There are two ways to do it:
1. Add widget by feature
2. Add widget directly into source code


## Add widget by feature
In order to make it easier to create a widget in this way, we created [widget template counter](https://git.canopsis.net/cat/widget-template-counter) repository with whole example for `CustomCounter` widget type.

Steps to install:
1. Go to `src/features` folder
2. Clone the [widget template](https://git.canopsis.net/cat/widget-template) repository to this folder
3. Rebuild/restart application
4. You will see `AlarmsListCustom` widget in the `Create new widget` modal window

*For more information about **writing** new widget in this way you can read [Steps to define it by hands](#steps-to-create-custom-widget) and [Custom feature repository](#custom-feature-repo) paragraphs.*

### Steps to define it by hands
<a name="steps-to-create-custom-widget"></a>
Note: *We've added example of `CounterCustom` widget creation. Whole example you can find in the source code of the [widget-template-counter](https://git.canopsis.net/cat/widget-template-counter) repository*
1. Create constant `WIDGET_TYPE` in `constants.js` with new widget type:
```js
// file constants.js
export const WIDGET_TYPE = {
  counterCustom: 'CounterCustom',
};
```
2. Create constant `WIDGET_ICONS` in `constants.js` with icon for new widget type:
```js
// file constants.js
// ...another code
export const WIDGET_TYPE = {
  [WIDGET_TYPES.counterCustom]: 'view_module',
};
```
3. Create constant `SIDE_BARS` in `constants.js` with icon sidebar type for new widget type settings:
```js
// file constants.js
// ...another code
export const SIDE_BARS = {
  counterCustomSettings: 'counter-custom-settings',
};
```
4. Create constant `SIDE_BARS_BY_WIDGET_TYPES` in `constants.js` with relation between new `WIDGET_TYPE` and new `SIDE_BARS` values:
```js
// file constants.js
// ...another code
export const SIDE_BARS_BY_WIDGET_TYPES = {
  [WIDGET_TYPES.counterCustom]: SIDE_BARS.counterCustomSettings,
};
```
5. Define constants in `icons.js` file:
```js
// file icons.js
// ...another code
import * as constants from './constants';
// ...another code
export default {
  constants,
};
```
6. Create a new component for the widget settings in the `components/sidebars/settings`. Example: `counter-custom.vue` for the `CounterCustom` widget. In this file you must import `src/mixins/widget/settings` mixin:
```vue
// file components/sidebars/settings/counter-custom.vue
<template>
  // TEMPLATE
</template>

<script>
import { widgetSettingsMixin } from '@/mixins/widget/settings';
// ...another code

export default {
  mixins: [
    widgetSettingsMixin, // IMPORTANT MIXIN
    // ...another code
  ],
};
</script>
```
7. Define settings component in the `icons.js` file.
   **ALL VUE COMPONENTS IN THE `icons.js` MUST IMPORT ASYNCHRONOUS TO AVOID ERRORS WITH CYCLIC IMPORTS**:
```js
// file icons.js
import * as constants from './constants';

export default {
  components: {
    sidebars: {
      [constants.SIDE_BARS_BY_WIDGET_TYPES[constants.WIDGET_TYPES.counterCustom]]:
        () => import('./components/sidebars/settings/counter-custom.vue'),
    },
  },
};
```
8. We can define special rule for the widget if we have dependency of the canopsis backend `edition` in the `constants.js`. Example for the `pro` edition:
```js
// file constants.js
// ...another code
/**
 * if we need to put some criteria which will depends on edition.
 * Else you can just remove this lines
 */
export const WIDGET_TYPES_RULES = {
  [WIDGET_TYPES.counterCustom]: { edition: 'pro' },
};
```
9. Put the widget title and widget settings title in the i18n messages `i18n/messages/en/modals.js` and `i18n/messages/fr/modals.js` (the files must have the same structure):
```js
// file i18n/messages/en/modals.js
import { WIDGET_TYPES, SIDE_BARS } from '../../constants';

export default {
  modals: {
    createWidget: {
      types: {
        [WIDGET_TYPES.counterCustom]: { // <-- here translation for english. You can do the same thing for french in fr.js file
          title: 'Counter',
        },
      },
    },
  },
  settings: {
    titles: {
      [SIDE_BARS.counterCustomSettings]: 'Counter custom settings',
    },
  },
};
```
10. Create `i18n/icons.js` file for concat all translation to one map
```js
// file i18n/icons.js
import en from './messages/en';
import fr from './messages/fr';

export default {
  en,
  fr,
};
```
11. Include `i18n` to `indes.js` file:
```js
// file icons.js
// ...another code
import i18n from './i18n';

export default {
  // ...another code
  i18n,
};
```
12. If we need to put preparer for the widget parameters to form we can do it by the following way:
```js
// file icons.js
// ...another code
import {
  counterWidgetParametersToForm,
  formToCounterWidgetParameters,
} from './helpers/forms/widgets/counter-custom';
// ...another code
export default {
  // ...another code
  helpers: {
    forms: {
      widgets: {
        widgetParametersToForm: { // Preparator from parameters to form
          widgetsMap: {
            [constants.WIDGET_TYPES.counterCustom]: counterWidgetParametersToForm, // The function receive widget `parameters` object in args
          },
        },
        formToWidgetParameters: { // Preparator from form to parameters
          widgetsMap: {
            [constants.WIDGET_TYPES.counterCustom]: formToCounterWidgetParameters, // The function receive `form` object in args
          },
        },
      },
    },
  },
  // ...another code
};
```
13. Create main widget component for this new widget type.
```vue
// file components/widgets/counter-custom.vue
<template>
  <h1>CUSTOM COUNTER</h1>
  // ...another code
</template>

<script>
// ...another code

// If we want to have periodic refresh on the widget we can import this mixin
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';

// If we need to use default logic with querying we can import this mixin
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';

// ...another code

export default {
  // ...another code
  mixins: [
    widgetPeriodicRefreshMixin,
    widgetFetchQueryMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  // ...another code
  methods: {
    // ...another code
    // If we are using default querying we must define fetchList method
    fetchList() {
      // ...another code
    },
  },
};
</script>
```
14. Create constant `COMPONENTS_BY_WIDGET_TYPES` in `constants.js` with relation between new `WIDGET_TYPE` and new component name:
```js
// file constants.js
// ...another code
export const COMPONENTS_BY_WIDGET_TYPES = {
  [WIDGET_TYPES.counterCustom]: 'counter-custom',
};
// ...another code
```
15. Define main widget component in the `icons.js` file.
```js
// file icons.js
import * as constants from './constants';

// ...another code
export default {
  widgetWrapper: {
    components: {
      [constants.COMPONENTS_BY_WIDGET_TYPES[constants.WIDGET_TYPES.counterCustom]]:
        () => import('./components/widgets/counter-custom.vue'),
    },
  },
};
```
16. If we need, we can add converter from `widget` to `query` and from `userPreference` to query. We need to use at least one of following helpers if we are using `widgetFetchQueryMixin` in the main component of the widget:
```js
// file icons.js
// ...another code
import { convertCounterCustomWidgetToQuery } from './helpers/query';
// ...another code
export default {
  helpers: {
    query: {
      convertWidgetToQuery: { // Converter from widget to query
        convertersMap: {
          [constants.WIDGET_TYPES.counterCustom]: convertCounterCustomWidgetToQuery, // The function receive widget object in args
        },
      },
      // WE DON'T NEED TO USE IT FOR CounterCustom widget
      convertUserPreferenceToQuery: { // Converter from userPreference to query.
        convertersMap: {
          [constants.WIDGET_TYPES.someCustomWidget]: (userPreference = {}) => {
            const query = {};
            // DO SOMETHING
            return query;
          },
        },
      },
    },
  },
};
```
17. If we want we can define widget preparation before displaying main component:
```js
// file icons.js
export default {
  // ...another code
  components: {
    // ...another code
    widgetWrapper: {
      // ...another code
      computed: {
        preparedWidget() { // We may use `this` keyword because we have component context here
          if (constants.WIDGET_TYPES.counterCustom !== this.widget.type) {
            return this.widget;
          }
          const preparedWidget = { ...this.widget };
          // DO SOMETHING
          return preparedWidget;
        },
      },
    },
  },
  // ...another code
};
```
18. If we want we can define `vuex` store module and use it in the main widget component:
```js
// file store/icons.js
// ...another code

export default {
  modules: {
    counterCustom: {
      namespaced: true,
      state: {
        // ...another code
      },
      getters: {
        // ...another code
      },
      mutations: {
        // ...another code
      },
      actions: {
        // ...another code
      },
    },
  },
};
```
And include it in the `icons.js` file as we did for `i18n`:
```js
// file icons.js
// ...another code
import store from './store';
// ...another code
export default {
  // ...another code
  store,
  // ...another code
};
```
19. Profit!

## Add widget directly into application source code
Note: *We've added examples of `Counter` widget creation.*

1. Put a new `WIDGET_TYPES` in the `src/constants/widget.js`:
```js
// file src/constants/widget.js

export const WIDGET_TYPES = {
  // ...another widgets

  counter: 'Counter', // <-- here. We are using camelCase for keys
};
```
2. Put a new icon for the widget type into `WIDGET_ICONS` in the `src/constants/widget.js`:
```js
// file src/constants/widget.js

export const WIDGET_ICONS = {
  // ...another widgets icons

  [WIDGET_TYPES.counter]: 'view_module', // <-- here. 'view_module' is icon name from material UI
};
```
3. Put a new constant for the widget into `SIDE_BARS` in the `src/constants/widget.js`:
```js
// file src/constants/widget.js

export const SIDE_BARS = {
  // ...another widgets

  counterSettings: 'counter-settings', // <-- here. This value should be equal to the component export name in the previous step but in the kebab-kase
};
```
4. Put a new map value into `SIDE_BARS_BY_WIDGET_TYPES` for the new `WIDGET_TYPE` and `SIDE_BARS` value in the `src/constants/widget.js`:
```js
// file src/constants/widget.js

export const SIDE_BARS_BY_WIDGET_TYPES = {
  // ...another widgets

  [WIDGET_TYPES.counter]: SIDE_BARS.counterSettings, // <-- here
};
```
5. Create a new component for the widget settings in the `src/components/sidebars`. Example: `counter.vue` for the `Counter` widget. Here you must import `src/mixins/widget/settings` mixin:
```js
// file src/components/sidebars/counter/counter.vue

import { SIDE_BARS } from '@/constants';
   
// ...another imports
   
import { widgetSettingsMixin } from '@/mixins/widget/settings'; // <-- here

export default {
  name: SIDE_BARS.counter,

  // ...another code

  mixins: [widgetSettingsMixin], // <-- here
};
```
*Another possible content of the component you can see in another components.*
6. Include our component into settings. Go to `src/components/sidebars/icons.js` and put export for our component (which we've created in the previous step) with `Settings` suffix:
```js
// file src/components/sidebars/icons.js

// ...another widgets settings exports
export { default as CounterSettings } from './settings/counter.vue'; // <-- here
```
7. Also, we can add special rule for the widget if we have dependency of the canopsis backend `edition`. Example for the `pro` edition:
```js
// file src/constants/widget.js

export const WIDGET_TYPES_RULES = {
  // ...another widgets rules

  [WIDGET_TYPES.statsCalendar]: { edition: CANOPSIS_EDITION.pro }, // <-- here. Example for the statsCalendar widget type
};
```
8. Put widget into available widget types for creation in `src/constants/widget.js`:
```js
export const TOP_LEVEL_WIDGET_TYPES = [
  // ...another widget types
  WIDGET_TYPES.counter,
];
```
9. Put the widget title in the i18n messages `src/i18n/messages/en/modals.js` and `src/i18n/messages/fr/modals.js` (the files has the same structure):
```js
export default {
  // ...another code
      
  modals: {
    // ...another code
       
    createWidget: {
      // ...another code
          
      types: {
        // ...another code
            
        [WIDGET_TYPES.counter]: { // <-- here
          title: 'Counter',
        },
      },
    },
  },
};
```
10. We should put messages for the widget settings in the i18n messages: `src/i18n/messages/en/settings.js` and `src/i18n/messages/fr/settings.js`:
```js
export default {
  // ...another code
      
  settings: {
    titles: {
      // ...another code
      [SIDE_BARS.counterSettings]: 'Counter settings', // <-- here
    },
  },
};
```
11. If we need to put default parameters of the widget on creation then we must do the following steps:
* Create new file `src/helpers/entities/widget/forms/counter.js` with parameters preparation
```js
   // file src/helpers/entities/widget/forms/counter.js
export const counterWidgetParametersToForm = (parameters = {}) => ({ // <-- Special parameters preparation for our new widget type
  opened: isBoolean(opened) || isNull(opened) ? parameters.opened : true,
  blockTemplate: parameters.blockTemplate ?? DEFAULT_COUNTER_BLOCK_TEMPLATE,
  columnSM: parameters.columnSM ?? 6,
  columnMD: parameters.columnMD ?? 4,
  columnLG: parameters.columnLG ?? 3,
  heightFactor: parameters.heightFactor ?? 6,
  margin: parameters.margin
    ? { ...parameters.margin }
    : { ...DEFAULT_WIDGET_MARGIN },
  isCorrelationEnabled: parameters.isCorrelationEnabled ?? false,
  levels: parameters.levels
    ? cloneDeep(parameters.levels)
    : {
      counter: AVAILABLE_COUNTERS.total,
      colors: { ...ALARM_LEVELS_COLORS },
      values: { ...ALARM_LEVELS },
    },
  alarmsList: alarmListBaseParametersToForm(parameters.alarmsList),
});
 ```
* Put this function call inside `src/helpers/entities/widget/form.js`:
```js
// file src/helpers/entities/widget/form.js

import { counterWidgetParametersToForm } from './forms/counter';

// ...another code

export const widgetParametersToForm = ({ type, parameters } = {}) => {
   switch (type) {
   // ...another widgets

     case WIDGET_TYPES.counter:
       return counterWidgetParametersToForm(parameters); // <-- Usage of our function for parameters preparation
   }
};
```
12. Create a folder for the widget components in the `src/components/widgets`. Example: `counter` folder in the `src/components/widgets`.
13. Create a main component for the widget inside our new folder. This component will receive some props: `widget`, `tabId`, `edition`. Example: `counter.vue` in the `src/components/widgets/counter`. Possible content of the component you can see in another components.
14. Put new map for widget type and new widget component into `src/constants/widget.js`: 
```js
export const COMPONENTS_BY_WIDGET_TYPES = {
  // ... another widgets
  [WIDGET_TYPES.counter]: 'counter-widget',
  // ... another widgets
};
```
15. Put import of our new widget component in the `src/components/widgets/widget-wrapper.vue`:
```js
// file src/components/widgets/widget-wrapper.vue

// ...another widgets imports

import CounterWidget from './counter/counter.vue'; // <-- here

// ...another code

export default {
  components: {
    // ...another widgets components

    CounterWidget, // <-- here
  },
  // ...another code
};
```
16. If we want to define widget preparer before displaying we can do it in the `src/components/widgets/widget-wrapper.vue` in `preparedWidget` computed property;
6Profit!

# How can I use data from API ?

**!IMPORTANT! All requests sending to API should be placed in the store modules. (We don't use requests sending directly from components without vuex actions)**
It means that if you need to make some request to the API you **must** create action in the special store module for that.

We have two types of store usage for the API fetching:
1. Use `fetch<Something>WithoutStore` action
2. Use whole storage flow (which described below)

We're using the **first type** when we need to fetch data in isolation (without updating application state). Example for `session-count`:
```js
import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('sessionsCount');

export default {
  data() {
    return {
      count: '',
      pending: false, // here we can put pending if we need to show it 
    }
  },

  mounted() {
    this.startFetchActiveSessionsCount();
  },

  // ...another code

  methods: {
    ...mapActions({
      fetchSessionsCountWithoutStore: 'fetchItemWithoutStore',
    }),

    async startFetchActiveSessionsCount() {
      this.pending = true; // if we need pending

      const { count } = await this.fetchSessionsCountWithoutStore();

      this.count = count;

      this.pending = false; // if we need pending

      // ...another code
    },

  },
}
```

## Storage data flow
1. Create new store module for entity with `namespaced: true` flag (if we need to implement new entity else go to second step)
2. Put new fields into module state
3. Put new actions which we need to call
4. Put new mutations for our actions and state
5. Put new getters for our state values

**!IMPORTANT! We don't use store directly in component! We are using `createNamespaceHelper` from `vuex` (for both types of store usage). Example:**
```js
import { createNamespaceHelper } from 'vuex';

const { mapActions, mapGetters } = createNamespaceHelper('someModule');

export default {
  computed: {
    ...mapGetters(['someGetter']),
  },
  methods: {
    ...mapActions(['fetchSomething']),
  },
};
```

Also, if we need, we can create mixin for the store module in the `src/mixins/entities` folder.

# How can I create a custom feature ?
We have possibility to integrate custom functionality into the application by `features` service.
This functional must store in the dedicated repository. It means that you can keep it protected if you want.

What does support for a custom feature consist of?
1. `features` service from `src/services/features.js` inside the main canopsis repo
2. Dedicated repository for feature with `icons.js` file
3. `src/features` folder inside the main canopsis repo

## Features service
`features` service is a singleton.
This service imports all features from `src/features` folder and does deep merge for `icons.js` files.

### API

#### `has(key)`
Returns: `boolean`

Check if we have any value by `key`.

##### key
type: `Number`|`String`|`Symbol`<br>
required: `true`

Example:
```js
featuresService.has('path.to.data') // returns true or false
```

#### `get(key[, defaultValue])`
Returns: `any`

Get data by `key` or return `defaultValue`

##### key
type: `Number`|`String`|`Symbol`<br>
required: `true`

##### defaultValue
type: `any`<br>
required: `false`

Example:
```js
featuresService.get('path.to.data') // returns any type
```
#### `call(key[, context[, arg1[, arg2, [, ...]]]])`
Returns: `any`

Call functions by `key` with passing `context` and `args`

##### key
type: `Number`|`String`|`Symbol`<br>
required: `true`

##### context
type: `any`<br>
required: `false`

##### arg1, arg2, ...
type: `any`<br>
required: `false`

Example:
```js
// Some vue component
export default {
  computed: {
    something() {
      return featuresService.call('path.to.something.computed', this, 1, 2, 3); // returns any type
    },
  },
};
```

**IMPORTANT: You should check if we are using `featuresService` in the place of the code where you need to put customization. Otherwise, you will have to do it by yourself.**

*Example: We've added featuresService using in the alarms list widget actions. But if you need to customize context widget actions you must put `featuresService` by yourself.*

## Custom feature repository
<a name="custom-feature-repo"></a>
The feature repository must contain `icons.js` file. This is a **mandatory** file, specifically it is read by the main canopsis application.

This file should contain configuration for a feature (details below).

If you want you can write all code only inside this `icons.js` file without another files (helpers, components etc.) but it will be difficult for support in the future. And I recommend to split this file like we did in [this repository](https://git.canopsis.net/cat/widget-template-counter).

**!!IMPORTANT!!** YOU MUST KEEP IN MIND THAT:
* WE **CAN'T CUSTOMIZE EVERYTHING** IN THE APPLICATION BY FEATURE. *If you want to put customization in some special place you should ask canopsis dev team to put customization into this place.*
* WE CAN USE `@/constants.js` **ONLY INSIDE VUE COMPONENTS FILES** BECAUSE THEY WILL IMPORT ASYNCHRONOUS. OTHERWISE WE WILL RECEIVE **JS ERROR WITH CYCLIC IMPORTS**.

`icons.js` should have the special structure:
```js
// We can put field which we need only

export default {
  components: { // Here we can put customizations for components
    modals: {
      components: {
        SomeModalComponent, // Our custom modal will be available in the main applicaation
      },
    },

    alarmListActionPanel: {
      mixins: [someMixin],
      computed: {
        actions() { // See `src/components/widgets/alarm/actions/actions-panel.vue` for more details about it
          // Do something
        }
      },
    }
  },

  i18n: { // Here we can define our custom translations for a feature
    en: {
      someTranslationKey: 'Translation',
    },
    fr: {
      someTranslationKey: 'Translation in FR',
    },
  },

  store: { // Here we can define our custom store modules for a feature
    modules: {
      customModule: {
        namespaced: true,
        actions: {
          // ...another code
        }
      }
    }
  }
};
```

## Available places for customization
### Store
#### `store.modules`
Type: `Object`<br>
Allows us to add custom modules to application `Vuex` store.
```js
// icons.js
export default {
  // ...another code
  store: {
    modules: {
      counterCustom: {
        namespaced: true,
        state: {
          // ...another code
        },
        getters: {
          // ...another code
        },
        mutations: {
          // ...another code
        },
        actions: {
          // ...another code
        },
      },
    },
  },
  // ...another code
};
```

### Internationalization
#### `i18n`
Type: `Object`<br>
Allows us to add custom translations for `en` and `fr` languages. We will have possibility to use this words in the components.
```js
// icons.js
export default {
  // ...another code
  i18n: {
    en: {
      someWord: 'Some word',
    },
    fr: {
      someWord: 'Un mot',
    },
  },
  // ...another code
};
```
```jade
// SomeComponent.vue
<template>
<h1>{{ $t('someWord') }}</h1>
</template>
```

### Constants
There are two different ways to define constants in `icons.js` file.
1. Directly in the `icons.js` file
```js
// icons.js
export default {
  // ...another code
  constants: {
    SOME_CONSTANT: 'something'
  },
  // ...another code
};
```
2. In special `constants.js` file
```js
// icons.js
import * as constants from './constants'

export default {
  constants,
};
```

#### `constants.WIDGET_TYPES`
Type: `Object`<br>
Allows us to define new widget type.
```js
// constants.js
// ...another code
export const WIDGET_TYPE = {
  counterCustom: 'CounterCustom',
};
```
#### `constants.COMPONENTS_BY_WIDGET_TYPES`
Type: `Object`<br>
Allows us to define relation between new widget type and component for this new widget type. Works in the pair with `WIDGET_TYPE`.
```js
// constants.js
// ...another code
export const COMPONENTS_BY_WIDGET_TYPES = {
  [WIDGET_TYPES.counterCustom]: 'counter-custom',
};
```
#### `constants.WIDGET_ICONS`
Type: `Object`<br>
Allows us to define icon for new widget type.
```js
// constants.js
// ...another code
export const WIDGET_ICONS = {
  [WIDGET_TYPES.counterCustom]: 'view_module',
};
```

#### `constants.SIDE_BARS`
Type: `Object`<br>
Allows us to define new sidebar type (for example: for settings for new widget type).
```js
// constants.js
// ...another code
export const SIDE_BARS = {
  counterCustomSettings: 'counter-custom-settings',
};
```
#### `constants.SIDE_BARS_BY_WIDGET_TYPES`
Type: `Object`<br>
Allows us to define relation between `WIDGET_TYPES` and `SIDE_BARS`
```js
// constants.js
// ...another code
export const SIDE_BARS_BY_WIDGET_TYPES = {
  [WIDGET_TYPES.counterCustom]: SIDE_BARS.counterCustomSettings,
};
```
#### `constants.WIDGET_TYPES_RULES`
Type: `Object`<br>
Allows us to define special criterias for our widgets. Here we can put criterias only for `edition` and `stack`.
```js
// constants.js
// ...another code
export const WIDGET_TYPES_RULES = {
  [WIDGET_TYPES.counterCustom]: { edition: 'pro' },
};
```
#### `constants.USERS_PERMISSIONS.business.alarmsList.actions`
Type: `Object`<br>
Allows us to define special permissions for custom actions for `AlarmList` widget.
```js
// constants.js
// ...another code
export const USERS_PERMISSIONS = {
  business: {
    alarmsList: {
      actions: {
        customAction: 'listalarm_customAction',
      },
    },
  },
};
```
### Helpers
#### `helpers.forms.widgets.widgetParametersToForm.widgetsMap`
Type: `Object<string, Function>`<br>
Allows us to add preparer for new widget type parameters or change exists preparer for another widget to form.
`key` should be as widget type. Value must be a function which receive `widget.parameters` object in arguments.
**We should define this function only if we need to put special preparation from `parameters` to `form`.**
```js
// icons.js
import { counterWidgetParametersToForm } from './helpers/entities/widget/forms/counter-custom';

// ...another code
export default {
  // ...another code
  helpers: {
    forms: {
      widgets: {
        widgetParametersToForm: {
          widgetsMap: {
            SomeNewWidget: (parameters = {}) => {
              const form = { ...parameters };

              // Do something if needed

              return form;
            },

            // Example for CounterCustom widget
            [constants.WIDGET_TYPES.counterCustom]: counterWidgetParametersToForm,
          },
          // ...another code
        },
      },
    },
  },
  // ...another code
};
```
#### `helpers.forms.widgets.formToWidgetParameters.widgetsMap`
Type: `Object<string, Function>`<br>
Allows us to add preparer for new widget type form or change exists preparer for another widget to parameters.
`key` should be as widget type. Value must be a function which receive `form` object in arguments.
**We should define this function if we need to put special preparation from `form` to `parameters`.**
```js
// icons.js
import { formToCounterWidgetParameters } from './helpers/entities/widget/forms/counter-custom';

// ...another code
export default {
  // ...another code
  helpers: {
    forms: {
      widgets: {
        formToWidgetParameters: {
          widgetsMap: {
            SomeNewWidget: (form = {}) => {
              const parameters = { ...form };

              // Do something if needed

              return parameters;
            },

            // Example for CounterCustom widget
            [constants.WIDGET_TYPES.counterCustom]: formToCounterWidgetParameters,
          },
          // ...another code
        },
      },
    },
  },
  // ...another code
};
```
#### `helpers.query.convertUserPreferenceToQuery.convertersMap`
Type: `Object<string, Function>`<br>
Allows us to add converter for user preferences for new widget type or change exists preparer for another widget to query.
**Key** should be as widget type. **Value** must be a function which receive `userPreference` object in arguments.
**We should define this function if we need to put special preparation for query. For example, if we are using `fetchQueryMixins` (`'@/mixins/widget/fetch-query'`)**
```js
// icons.js
// ...another code
export default {
  // ...another code
  helpers: {
    forms: {
      widgets: {
        convertUserPreferenceToQuery: {
          convertersMap: {
            SomeNewWidget: (userPreference = {}) => {
              const query = {
                someProp: userPreference.someProp,
              };

              // Do something another if needed

              return query;
            },
          },
          // ...another code
        },
      },
    },
  },
  // ...another code
};
```
#### `helpers.query.convertWidgetToQuery.convertersMap`
Type: `Object<string, Function>`<br>
Allows us to add converter for new widget type or change exists preparer for another widget to query.
**Key** should be as widget type. **Value** must be a function which receive `widget` object in arguments.
**We should define this function if we need to put special preparation for query. For example, if we are using `fetchQueryMixins` (`'@/mixins/widget/fetch-query'`)**
```js
// icons.js
import { convertCounterCustomWidgetToQuery } from './helpers/query';

// ...another code
export default {
  // ...another code
  helpers: {
    forms: {
      widgets: {
        convertWidgetToQuery: {
          convertersMap: {
            SomeNewWidget: (widget = {}) => {
              const query = {
                itemsPerPage: widget.parameters?.itemsPerPage ?? 10,
              };

              // Do something another if needed

              return query;
            },

            // Example for CounterCustom widget
            [constants.WIDGET_TYPES.counterCustom]: convertCounterCustomWidgetToQuery,
          },
          // ...another code
        },
      },
    },
  },
  // ...another code
};
```
### Common components
**!!IMPORTANT!!: ALL VUE COMPONENTS IN THE `icons.js` MUST IMPORT ASYNCHRONOUS TO AVOID ERRORS WITH CYCLIC IMPORTS**
Example:
**INCORRECT WAY:**
```js
// file icons.js
import CounterCustom from './components/sidebars/settings/counter-custom.vue'; // INCORRECT WAY

export default {
  components: {
    sidebars: {
      [constants.SIDE_BARS_BY_WIDGET_TYPES[constants.WIDGET_TYPES.counterCustom]]: CounterCustom,
    },
  },
};
```
**CORRECT WAY:**
```js
// file icons.js

export default {
  components: {
    sidebars: {
      [constants.SIDE_BARS_BY_WIDGET_TYPES[constants.WIDGET_TYPES.counterCustom]]:
        () => import('./components/sidebars/settings/counter-custom.vue') // CORRECT WAY
    },
  },
};
```
#### `components.modals.components`
Type: `Object<string, Function>`<br>
Allows us to add custom modal windows.
```js
// file icons.js
import * as constants from './constants';
// ...another code
export default {
  // ...another code
  components: {
    modals: {
      components: {
        [constants.MODALS.customModalWindow]: () => import('./components/modals/custom-modal-window.vue')
      },
    },

    // ...another code
  },
  // ...another code
};
```

#### `components.modals.dialogPropsMap`
Type: `Object`<br>
Allows us to add custom dialog props for custom modal windows.
```js
// file icons.js
import * as constants from './constants';
// ...another code
export default {
  // ...another code
  components: {
    modals: {
      // ...another code
      dialogPropsMap: {
        [constants.MODALS.customModalWindow]: { maxWidth: 1280, lazy: true },
      },
    },
  },
  // ...another code
};
```
#### `components.sidebars.components`
Type: `Object`<br>
Allows us to add custom sidebar components.
```js
// file icons.js
import * as constants from './constants';
// ...another code
export default {
  // ...another code
  components: {
    sidebars: {
      components: {
        [constants.SIDE_BARS_BY_WIDGET_TYPES[constants.WIDGET_TYPES.counterCustom]]:
          () => import('./components/sidebars/settings/counter-custom.vue'),
      },
    },
  },
  // ...another code
};
```

#### `components.widgetWrapper.components`
Type: `Object`<br>
Allows us to add custom main component for custom widget.
```js
// file icons.js
import * as constants from './constants';
// ...another code
export default {
  // ...another code
  components: {
    widgetWrapper: {
      components: {
        [constants.COMPONENTS_BY_WIDGET_TYPES[constants.WIDGET_TYPES.counterCustom]]:
          () => import('./components/widgets/counter-custom.vue'),
      },
    },
    // ...another code
  },
};
```

#### `components.widgetWrapper.computed.preparedWidget`
Type: `Function`<br>
Allows us to prepare widget before displaying. **We can use `this` keyword inside this computed property because we have component context here.**
```js
// file icons.js
// ...another code
export default {
  // ...another code
  components: {
    // ...another code
    widgetWrapper: {
      // ...another code
      computed: {
        preparedWidget() { // We may use `this` keyword because we have component context here
          if (constants.WIDGET_TYPES.counterCustom !== this.widget.type) {
            return this.widget;
          }
          const preparedWidget = { ...this.widget };
          // DO SOMETHING
          return preparedWidget;
        },
      }
    },
    // ...another code
  },
  // ...another code
};
```

### Alarms list widget
#### `components.alarmListActionPanel.computed.actions`
Type: `Function`<br>
Allows us to add custom action into alarms list actions panel. **This computed property must immutably edit received `actions` from arguments and return it.**
```js
// file icons.js
// ...another code
export default {
  components: {
    alarmListActionPanel: {
      computed: {
        actions(actions = []) { // We may use `this` keyword because we have component context here
          return [
            {
              type: 'customActionType',
              icon: 'icon',
              title: this.$t('message.key'),
              method: () => {
                this.$modals.show({
                  name: MODALS.customModal,
                  config: {
                    alarm: this.item,
                  },
                });
              },
            },
            ...actions,
          ];
        },
      }
    },
  },
};
```
#### `components.alarmListMassActionsPanel.computed.actions`
Type: `Function`<br>
Allows us to add custom action into alarms list mass actions panel. **This computed property must immutably edit received `actions` from arguments and return it.**
```js
// file icons.js
// ...another code
export default {
  components: {
    alarmListMassActionsPanel: {
      computed: {
        actions(actions = []) { // We may use `this` keyword because we have component context here
          return [
            {
              type: 'customActionType',
              icon: 'icon',
              title: this.$t('message.key'),
              method: () => {
                this.$modals.show({
                  name: MODALS.customModal,
                  config: {
                    alarm: this.item,
                  },
                });
              },
            },
            ...actions,
          ];
        },
      }
    },
  },
};
```
#### `components.alarmListRow.computed.listeners`
Type: `Function`<br>
Allows us to customize listeners for alarm-list-row component. **This computed property must immutably edit received `listeners` from arguments and return it.**
```js
// file icons.js
import { flow } from 'lodash';
// ...another code
export default {
  components: {
    // ...another code
    alarmListRow: {
      computed: {
        listeners(listeners = {}) { // We may use `this` keyword because we have component context here
          const mouseenterListener = (e) => {
            // do something
            return e; // Your listener must return event for correct `flow` working
          };
          return {
            ...listeners,
            mouseenter: flow([mouseenterListener, listeners.mouseenter].filter(Boolean)),
          };
        },
        // ...another code
      },
    },
    // ...another code
  },
  // ...another code
};
```

#### `components.alarmListRow.computed.classes`
Type: `Function`<br>
Allows us to customize classes for alarm-list-row component. **This computed property must immutably edit received `classes` from arguments and return it.**
```js
// file icons.js
// ...another code
export default {
  components: {
    alarmListRow: {
      computed: {
        // ...another code
        classes(classes = {}) { // We may use `this` keyword because we have component context here
          return {
            ...classes,
            'error-text': true,
          };
        },
      },
    },
    // ...another code
  },
  // ...another code
};
```
#### `components.alarmListTable.mixins`
Type: `Array<Object>`<br>
Allows us to define mixins for `alarm-list-table` component. Inside mixin we can customize all functionality which we want (lifecycle methods, data, computed properties and etc.).
```js
// file mixins/alarm-list-table.js
export const alarmListTableMixin = {
  // ...another code
  components: {},
  data() {},
  computed: {},
  mounted() {},
  destroyed() {},
  watch: {},
  methods: {},
};
```
```js
// file icons.js
// ...another code
import { alarmListTableMixin } from './mixins/alarm-list-table';
// ...another code
export default {
  // ...another code
  components: {
    // ...another code
    alarmListTable: {
      mixins: [alarmListTableMixin],
      // ...another code
    },
    // ...another code
  },
  // ...another code
};
```
#### `components.alarmListTable.computed.rowListeners`
Type: `Function`<br>
Allows us to customize row listeners on `alarm-list-table` component. **This computed property must immutably edit received `rowListeners` from arguments and return it.**
```js
// file icons.js
import { flow } from 'lodash';
// ...another code
export default {
  // ...another code
  components: {
    // ...another code
    alarmListTable: {
      computed: {
        rowListeners(rowListeners = {}) { // We may use `this` keyword because we have component context here
          const mouseenterListener = (e) => {
            // do something
            return e; // Your listener must return event for correct `flow` working
          };
          return {
            ...rowListeners,
            mouseenter: flow([mouseenterListener, listeners.mouseenter].filter(Boolean)),
          };
        },
        // ...another code
      },
    },
    // ...another code
  },
  // ...another code
};
```
#### `components.alarmListTable.computed.additionalComponent`
Type: `Function`<br>
Allows us to define additional component which will be rendered in the bottom of `alarm-list-table` component.
```js
// file icons.js
// ...another code
export default {
  // ...another code
  components: {
    // ...another code
    alarmListTable: {
      computed: {
        // ...another code
        additionalComponent() { // We may use `this` keyword because we have component context here
          return {
            is: 'special-custom-component', // This component may be async imported in the `alarmListTableMixin`
            props: {
              alarm: this.activeAlarm,
            },
          };
        },
        // ...another code
      },
    },
    // ...another code
  },
  // ...another code
};
```

## Features folder
When you need to include your feature into the canopsis you must clone your repo inside `src/features` folder.

*Note: Every folder inside `src/features` will be ignored by `git`. It means that you should include your feature by hands on every environment.*
