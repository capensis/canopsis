import Vuetify from 'vuetify/lib';
import * as VuetifyComponents from 'vuetify/lib/components';
import * as VuetifyDirectives from 'vuetify/lib/directives';

import icons from './components/icons';
import ClickOutside from './directives/click-outside';
// import VTabs from './components/v-tabs/v-tabs.vue';
// import VSelect from './components/v-select';
// import VTextarea from './components/v-textarea/v-textarea.vue';
// import VTooltip from './components/v-tooltip/v-tooltip.vue';
// import VListGroup from './components/v-list-group/v-list-group.vue';
// import VExpansionPanelContent from './components/v-expansion-panel-content/v-expansion-panel-content.vue';
// import VCombobox from './components/v-combobox/v-combobox.vue';
// import VChipGroup from './components/v-chip-group/v-chip-group.vue';
// import VChip from './components/v-chip/v-chip.vue';
import VMenu from './components/v-menu/v-menu.vue';
import VDialog from './components/v-dialog/v-dialog.vue';
import VNavigationDrawer from './components/v-navigation-drawer/v-navigation-drawer.vue';
import VDataTable from './components/v-data-table/v-data-table.vue';
import VSpeedDial from './components/v-speed-dial/v-speed-dial.vue';
import VCalendar from './components/v-calendar/v-calendar.vue';

import './styles/vuetify.scss';

export const createVuetify = (Vue, options) => {
  Vue.use(Vuetify, {
    components: {
      ...VuetifyComponents,
      VMenu,
      VNavigationDrawer,
      VDialog,
      VSpeedDial,
      VDataTable,
      VCalendar,
    },

    directives: {
      ...VuetifyDirectives,
      ClickOutside,
    },
  });

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
