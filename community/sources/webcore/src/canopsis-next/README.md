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
* Open google chrome with the tag :  `google-chrome --disable-web-security --user-data-dir`
* Open your Canopsis in this chrome, then login, to have the back-end data
* Open the link provided by yarn in chrome

# How can I create a new widget type ?
There are two ways to do it:
1. Add widget by feature
2. Add widget directly into source code


## Add widget by feature
In order to make it easier to create a widget in this way, we created [widget template](https://git.canopsis.net/cat/widget-template)

Steps to install:
1. Go to `src/features` folder
2. Clone the [widget template](https://git.canopsis.net/cat/widget-template) repository to this folder
3. Rebuild/restart application
4. You will see `AlarmsListCustom` widget in the `Create new widget` modal window

*For more information about **writing** new widget in this way you can read [Custom feature repository](#custom-feature-repo) paragraph and README.md in the [widget template](https://git.canopsis.net/cat/widget-template) repository.*

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

5. Create a new component for the widget settings in the `src/components/sidebars/settings`. Example: `counter.vue` for the `Counter` widget. Here you must import `src/mixins/widget/settings` mixin:
    ```js
    // file src/components/side-bars/settings/counter.vue

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

6. Include our component into settings. Go to `src/components/sidebars/index.js` and put export for our component (which we've created in the previous step) with `Settings` suffix:
    ```js
    // file src/components/sidebars/index.js

    // ...another widgets settings exports
    export { default as CounterSettings } from './settings/counter.vue'; // <-- here
    ```

7. Put new widget into `availableTypes` in the `src/components/modals/view/create-widget.vue`:
    ```js
    // file src/components/modals/view/create-widget.vue

    export default {
      // ...another code
   
      computed: {
        availableTypes() {
          return [
            // ...another widgets

            WIDGET_TYPES.counter, // <-- here
          ].filter((widgetType) => {
            // ...another code
          });
        },
      },

      // ...another code
    };
    ```

8. Also we can add special rule for the widget if we have dependency of the canopsis backend `edition`. Example for the `pro` edition:
    ```js
    // file src/constants/widget.js

    export const WIDGET_TYPES_RULES = {
      // ..another widgets rules

      [WIDGET_TYPES.statsCalendar]: { edition: CANOPSIS_EDITION.pro }, // <-- here. Example for the statsCalendar widget type
    };
    ```

9. Put the widget title in the i18n messages `src/i18n/messages/en.js` and `src/i18n/messages/fr.js` (the files has the same structure):
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

10. We should put messages for the widget settings in the i18n messages: `src/i18n/messages/en.js` and `src/i18n/messages/fr.js`:
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
    1. Create new file `src/helpers/forms/widgets/counter.js` with parameters preparation
    ```js
            // file src/helpers/forms/widgets/counter.js
    
    export const counterWidgetParametersToForm = (parameters = {}) => ({ // <-- Special parameters preparation for our new widget type
      opened: parameters.opened ?? true,
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
    2. Put this function call inside `src/helpers/forms/widgets/common.js`:
    ```js
        // file src/helpers/forms/widgets/common.js
    
        import { counterWidgetParametersToForm } from './counter';

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

14. Put import of our new widget component and put new map value into `widgetProps` in the `src/components/widgets/widget-wrapper.vue`:
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
      computed: {
        // ...another code
        widgetProps() {
          // ...another code
        
          const widgetComponentsMap = {
            // ...another widgets

            [WIDGET_TYPES.counter]: 'counter-widget', // <-- here
          };
        },
      },
    };
    ```

15. Profit!

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

Also, if we need we can create mixin for the store module in the `src/mixins/entities` folder.

# How can I create a custom feature ?

We have possibility to integrate custom functionality into the application by `features` service.
This functional must store in the dedicated repository. It means that you can keep it protected if you want.

What does support for a custom feature consist of?
1. `features` service from `src/services/features.js` inside the main canopsis repo
2. Dedicated repository for feature with `index.js` file
3. `src/features` folder inside the main canopsis repo

## Features service
`features` service is a singletone.
This service imports all features from `src/features` folder and does deep merge for `index.js` files.

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

*Example: We've added featuresService using in the alarms list widget actions. But if you need to customize context widget actions you must to put `featuresService` by yourself.*

## Custom feature repository
<a name="custom-feature-repo"></a>
Feature repository must contain `index.js` file with configurations. Here we should define the points which we want to customize.

`index.js` should have the special structure:
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
          // ...
        }
      }
    }
  }
};
```

*Note: You can split config to several files for better developer experience*

## Features folder
When you need to include your feature into the canopsis you must clone your repo inside `src/features` folder.

*Note: Every folder inside `src/features` will be ignored by `git`. It means that you should include your feature by hands on every environment.*
