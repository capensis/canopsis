import Vuetify from 'vuetify/lib';
import * as VuetifyComponents from 'vuetify/lib/components';
import * as VuetifyDirectives from 'vuetify/lib/directives';

import icons from './components/icons';
// import ClickOutside from './directives/click-outside';
// import VTabs from './components/v-tabs/v-tabs.vue';
// import VSelect from './components/v-select';
// import VTextarea from './components/v-textarea/v-textarea.vue';
// import VTooltip from './components/v-tooltip/v-tooltip.vue';
// import VMenu from './components/v-menu/v-menu.vue';
// import VListGroup from './components/v-list-group/v-list-group.vue';
// import VNavigationDrawer from './components/v-navigation-drawer/v-navigation-drawer.vue';
// import VExpansionPanelContent from './components/v-expansion-panel-content/v-expansion-panel-content.vue';
// import VDialog from './components/v-dialog/v-dialog.vue';
// import VCombobox from './components/v-combobox/v-combobox.vue';
// import VChipGroup from './components/v-chip-group/v-chip-group.vue';
// import VChip from './components/v-chip/v-chip.vue';
import VDataTable from './components/v-data-table/v-data-table.vue';
import VSpeedDial from './components/v-speed-dial/v-speed-dial.vue';

import './styles/vuetify.scss';

export const createVuetify = (Vue, options) => {
  Vue.use(Vuetify, {
    components: {
      ...VuetifyComponents,

      VSpeedDial,
      VDataTable,
    },

    directives: {
      ...VuetifyDirectives,
    },
  });

  // Vue.component('v-chip', VChip);
  // Vue.component('v-chip-group', VChipGroup);
  // Vue.component('v-combobox', VCombobox);
  // Vue.component('v-dialog', VDialog);
  // Vue.component('v-expansion-panel-content', VExpansionPanelContent);
  // Vue.component('v-navigation-drawer', VNavigationDrawer);
  // Vue.component('v-list-group', VListGroup);
  // Vue.component('v-menu', VMenu);
  // Vue.component('v-tooltip', VTooltip);
  // Vue.component('v-data-table', VDataTable);
  // Vue.component('v-textarea', VTextarea);
  // Vue.component('v-select', VSelect);
  // Vue.component('v-tabs', VTabs);

  // Vue.directive('click-outside', ClickOutside);

  return new Vuetify({
    icons: {
      iconfont: 'md',
      values: icons,
    },
    theme: {
      options: {
        customProperties: true,
      },
      ...options.theme,
    },
  });
};
