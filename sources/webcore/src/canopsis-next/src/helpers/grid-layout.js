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

export function getNewWidgetGridParametersY(widgets) {
  return widgets.reduce((acc, { gridParameters }) => {
    if (gridParameters.mobile.y >= acc.mobile) {
      acc.mobile = gridParameters.mobile.y + gridParameters.mobile.h + 1;
    }

    if (gridParameters.tablet.y >= acc.tablet) {
      acc.tablet = gridParameters.tablet.y + gridParameters.mobile.h + 1;
    }

    if (gridParameters.desktop.y >= acc.desktop) {
      acc.desktop = gridParameters.desktop.y + gridParameters.mobile.h + 1;
    }

    return acc;
  }, { mobile: 0, tablet: 0, desktop: 0 });
}

export default {
  getWidgetsLayoutBySize,
};
