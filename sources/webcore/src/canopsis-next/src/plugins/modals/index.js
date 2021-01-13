import { get, isFunction } from 'lodash';

import { setField } from '@/helpers/immutable';

import modalsStoreModule from './store';

import TheModals from './components/the-modals.vue';
import ModalBase from './components/modal-base.vue';

import innerMixin from './mixins/inner';

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

    const preparedComponents = Object.entries(components).reduce((acc, [key, value]) => {
      acc[key] = isFunction(value) ? value : setField(value, 'mixins', (mixins = []) => [innerMixin, ...mixins]);

      return acc;
    }, {});

    Vue.component(componentName, TheModals);
    Vue.component('modal-base', {
      components: preparedComponents,

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
          minimize: ({ id } = {}) => store.dispatch(`${moduleName}/minimize`, { id: id || get(this.modal, 'id') }),
          maximize: ({ id } = {}) => store.dispatch(`${moduleName}/maximize`, { id: id || get(this.modal, 'id') }),
        };
      },
    });
  },
};
