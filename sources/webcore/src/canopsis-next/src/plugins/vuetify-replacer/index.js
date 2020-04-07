import VDialog from './components/v-dialog/v-dialog.vue';
import VCheckboxFunctional from './components/v-checkbox-functional/v-checkbox-functional.vue';
import VExpansionPanelContent from './components/v-expansion-panel-content/v-expansion-panel-content.vue';
import VNavigationDrawer from './components/v-navigation-drawer/v-navigation-drawer.vue';

export default {
  install(Vue) {
    Vue.component('v-dialog', VDialog);
    Vue.component('v-checkbox-functional', VCheckboxFunctional);
    Vue.component('v-expansion-panel-content', VExpansionPanelContent);
    Vue.component('v-navigation-drawer', VNavigationDrawer);
  },
};
