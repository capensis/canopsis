import { get } from 'lodash';

import modalsStoreModule from './store';

import TheModals from './components/the-modals.vue';
import ModalBase from './components/modal-base.vue';

export default {
  install(Vue, {
    store,
    components = {},
    dialogPropsMap = {},
    moduleName = 'modals',
    componentName = 'the-modals',
  }) {
    if (!store) {
      throw new Error('Missing store option');
    }

    Vue.component(componentName, TheModals);
    Vue.component('modal-base', {
      components,

      extends: ModalBase,
    });

    store.registerModule(moduleName, modalsStoreModule);

    Object.defineProperty(Vue.prototype, '$modals', {
      get() {
        return {
          moduleName,
          dialogPropsMap,

          show: modal => store.dispatch(`${moduleName}/show`, modal),
          hide: ({ id } = {}) => store.dispatch(`${moduleName}/hide`, { id: id || get(this.modal, 'id') }),
        };
      },
    });
  },
};
