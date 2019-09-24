import popupsStoreModule from './store';

import ThePopups from './components/the-popups.vue';

export default {
  install(Vue, {
    store, i18n, moduleName = 'popups', componentName = 'the-popups',
  } = {}) {
    if (!store || !i18n) {
      throw new Error('Missing required options.');
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

    /**
     * TODO: Update it to Vue.config.errorHandler after updating to 2.6.0+ Vue version
     */
    window.addEventListener('unhandledrejection', () => {
      store.dispatch(`${moduleName}/add`, { type: 'error', text: i18n.t('errors.default') });
    });
  },
};
