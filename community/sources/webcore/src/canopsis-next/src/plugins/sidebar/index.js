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

          show: sidebar => store.dispatch(`${moduleName}/show`, sidebar),
          hide: () => store.dispatch(`${moduleName}/hide`),
        };
      },
    });
  },
};
