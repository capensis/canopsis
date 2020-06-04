import { omit } from 'lodash';

import { WIDGET_GRID_SIZES_KEYS } from '@/constants';

/**
 * Get layout for special size by widgets and oldLayouts
 *
 * @param {Array} widgets
 * @param {Array} [oldLayout = []]
 * @param {string} [size = WIDGET_GRID_SIZES_KEYS.desktop]
 * @returns {Object}
 */
export function getWidgetsLayoutBySize(widgets, oldLayout = [], size = WIDGET_GRID_SIZES_KEYS.desktop) {
  return widgets.map((widget) => {
    const oldLayoutItem = oldLayout.find(({ i }) => i === widget._id);
    const newLayoutItem = oldLayoutItem ? omit(oldLayoutItem, ['i', 'widget']) : { ...widget.gridParameters[size] };

    newLayoutItem.i = widget._id;
    newLayoutItem.widget = widget;

    return newLayoutItem;
  });
}

export default {
  getWidgetsLayoutBySize,
};
