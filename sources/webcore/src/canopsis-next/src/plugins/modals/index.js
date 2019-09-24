import { get } from 'lodash';

import modalsStoreModule from './store';

import TheModals from './components/the-modals.vue';

export default {
  install(Vue, { store, moduleName = 'modals', componentName = 'the-modals' }) {
    if (!store) {
      throw new Error('Missing store option');
    }

    Vue.component(componentName, TheModals);

    store.registerModule(moduleName, modalsStoreModule);

    Object.defineProperty(Vue.prototype, '$modals', {
      get() {
        return {
          show: modal => store.dispatch(`${moduleName}/show`, modal),
          hide: ({ id } = {}) => store.dispatch(`${moduleName}/hide`, { id: id || get(this.modal, 'id') }),
        };
      },
    });
  },
};
