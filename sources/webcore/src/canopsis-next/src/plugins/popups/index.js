import popupsStoreModule from './store';

import ThePopups from './components/the-popups.vue';

export default {
  install(Vue, {
    store, moduleName = 'popups', componentName = 'the-popups',
  } = {}) {
    if (!store) {
      throw new Error('Missing store option.');
    }

    Vue.component(componentName, ThePopups);

    store.registerModule(moduleName, popupsStoreModule);

    Object.defineProperty(Vue.prototype, '$popups', {
      get() {
        return {
          moduleName,

          add(popup) {
            return store.dispatch(`${moduleName}/add`, popup);
          },
          remove({ id }) {
            return store.dispatch(`${moduleName}/remove`, { id });
          },

          success(popup) {
            return store.dispatch(`${moduleName}/success`, popup);
          },

          info(popup) {
            return store.dispatch(`${moduleName}/info`, popup);
          },

          warning(popup) {
            return store.dispatch(`${moduleName}/warning`, popup);
          },

          error(popup) {
            return store.dispatch(`${moduleName}/error`, popup);
          },
        };
      },
    });
  },
};
