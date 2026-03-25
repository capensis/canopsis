import { get } from 'lodash';

import TheSidebar from './components/the-sidebar.vue';
import SidebarBase from './components/sidebar-base.vue';
import sidebarStoreModule from './store';

export default {
  install(Vue, {
    store,
    components = {},
    moduleName = 'sidebar',
    componentName = 'the-sidebar',
  }) {
    if (!store) {
      throw new Error('Missing store option');
    }

    Vue.component(componentName, TheSidebar);
    Vue.component('sidebar-base', {
      components,

      extends: SidebarBase,
    });

    store.registerModule(moduleName, sidebarStoreModule);

    Object.defineProperty(Vue.prototype, '$sidebar', {
      get() {
        return {
          moduleName,

          show: payload => store.dispatch(`${moduleName}/show`, payload),
          hide: ({ id } = {}) => store.dispatch(`${moduleName}/hide`, { id: id || get(this.sidebar, 'id') }),
          minimize: ({ id } = {}) => store.dispatch(`${moduleName}/minimize`, { id: id || get(this.sidebar, 'id') }),
          maximize: ({ id } = {}) => store.dispatch(`${moduleName}/maximize`, { id: id || get(this.sidebar, 'id') }),
          updateConfig: ({ id, config } = {}) => store.dispatch(`${moduleName}/updateConfig`, { id: id || get(this.sidebar, 'id'), config }),
        };
      },
    });
  },
};
