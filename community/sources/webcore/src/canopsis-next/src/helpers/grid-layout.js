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
    const newLayoutItem = oldLayoutItem ? omit(oldLayoutItem, ['i', 'widget']) : { ...widget.grid_parameters[size] };

    newLayoutItem.i = widget._id;
    newLayoutItem.widget = widget;

    return newLayoutItem;
  });
}

/**
 * Get y positions for every sizes for new widget
 *
 * @param {Array} widgets
 * @returns {{ mobile: number, tablet: number, desktop: number }}
 */
export function getNewWidgetGridParametersY(widgets) {
  return widgets.reduce((acc, { grid_parameters: gridParameters }) => {
    const newMobileYPosition = gridParameters.mobile.y + gridParameters.mobile.h;
    const newTabletYPosition = gridParameters.tablet.y + gridParameters.tablet.h;
    const newDesktopYPosition = gridParameters.desktop.y + gridParameters.desktop.h;

    if (newMobileYPosition > acc.mobile) {
      acc.mobile = newMobileYPosition;
    }

    if (newTabletYPosition > acc.tablet) {
      acc.tablet = newTabletYPosition;
    }

    if (newDesktopYPosition > acc.desktop) {
      acc.desktop = newDesktopYPosition;
    }

    return acc;
  }, { mobile: 0, tablet: 0, desktop: 0 });
}

export default {
  getWidgetsLayoutBySize,
  getNewWidgetGridParametersY,
};
