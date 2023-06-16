import { kebabCase } from 'lodash';

import * as buttonsComponents from './buttons';
import * as chipsComponents from './chips';
import * as commonComponents from './common';
import * as iconsComponents from './icons';
import * as overlayComponents from './overlay';
import * as tableComponents from './table';
import * as fieldsComponents from './fields';

/**
 * @param {import('vue').VueConstructor | import('vue').Vue} Vue
 */
export const registerApplicationComponents = (Vue) => {
  Object.entries({
    ...buttonsComponents,
    ...chipsComponents,
    ...commonComponents,
    ...iconsComponents,
    ...overlayComponents,
    ...tableComponents,
    ...fieldsComponents,
  }).forEach(([name, component]) => {
    Vue.component(kebabCase(name), component);
  });
};
