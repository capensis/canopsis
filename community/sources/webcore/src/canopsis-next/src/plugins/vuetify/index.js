import Vuetify from 'vuetify/lib';
import * as VuetifyComponents from 'vuetify/lib/components';
import * as VuetifyDirectives from 'vuetify/lib/directives';

import icons from './components/icons';
import ClickOutside from './directives/click-outside';
import VCombobox from './components/v-combobox/v-combobox.vue';
import VMenu from './components/v-menu/v-menu.vue';
import VDialog from './components/v-dialog/v-dialog.vue';
import VNavigationDrawer from './components/v-navigation-drawer/v-navigation-drawer.vue';
import VDataTable from './components/v-data-table/v-data-table.vue';
import VSimpleTable from './components/v-data-table/v-simple-table.vue';
import VDataTableHeaderDesktop from './components/v-data-table/v-data-table-header-desktop.vue';
import VDataTableHeader from './components/v-data-table/v-data-table-header.vue';
import VSpeedDial from './components/v-speed-dial/v-speed-dial.vue';
import VCalendar from './components/v-calendar/v-calendar.vue';
import VTooltip from './components/v-tooltip/v-tooltip.vue';
import VSelect from './components/v-select/v-select.vue';
import VListGroup from './components/v-list/v-list-group.vue';
import VListItemMask from './components/v-list/v-list-item-mask.vue';
import VTabs from './components/v-tabs/v-tabs.vue';

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
      VSimpleTable,
      VDataTableHeaderDesktop,
      VDataTableHeader,
      VCalendar,
      VTooltip,
      VSelect,
      VListGroup,
      VListItemMask,
      VCombobox,
      VTabs,
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
