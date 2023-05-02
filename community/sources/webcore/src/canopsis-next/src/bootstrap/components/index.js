import { kebabCase } from 'lodash';

import CEnabled from '@/components/icons/c-enabled.vue';
import AlarmsListTable from '@/components/widgets/alarm/partials/alarms-list-table.vue';

import * as commonComponents from './common';
import * as fieldsComponents from './fields';

/**
 * @param {import('vue').VueConstructor | import('vue').Vue} Vue
 */
export const registerApplicationComponents = (Vue) => {
  Object.entries({
    CEnabled,
    AlarmsListTable,
    ...commonComponents,
    ...fieldsComponents,
  }).forEach(([name, component]) => {
    Vue.component(kebabCase(name), component);
  });
};
