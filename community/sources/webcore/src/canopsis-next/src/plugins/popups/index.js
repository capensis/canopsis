import popupsStoreModule from './store';

import ThePopups from './components/the-popups.vue';

export default {
  install(Vue, { store, moduleName = 'popups', componentName = 'the-popups' } = {}) {
    if (!store) {
      throw new Error('Missing store option');
    }

    Vue.component(componentName, ThePopups);

    store.registerModule(moduleName, popupsStoreModule);

    Object.defineProperty(Vue.prototype, '$popups', {
      get() {
        return {
          moduleName,

          add: popup => store.dispatch(`${moduleName}/add`, popup),
          remove: popup => store.dispatch(`${moduleName}/remove`, popup),
          success: popup => store.dispatch(`${moduleName}/success`, popup),
          info: popup => store.dispatch(`${moduleName}/info`, popup),
          warning: popup => store.dispatch(`${moduleName}/warning`, popup),
          error: popup => store.dispatch(`${moduleName}/error`, popup),
          setDefaultCloseTime: popup => store.dispatch(`${moduleName}/setDefaultCloseTime`, popup),
        };
      },
    });
  },
};
