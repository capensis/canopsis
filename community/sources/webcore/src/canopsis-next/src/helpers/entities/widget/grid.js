import { get, omit, pick } from 'lodash';

import { WIDGET_GRID_SIZES_KEYS } from '@/constants';

import { compactLayout } from '@/helpers/grid';

/**
 * @typedef {Object} GridLayoutWidgetGrid
 * @property {string} _id
 * @property {WidgetGridParameters} grid_parameters
 */

/**
 * Get layout for special size by widgets and oldLayouts
 *
 * @param {Widget[]} [widgets = []]
 * @param {GridLayout} [oldLayout = []]
 * @param {string} [size = WIDGET_GRID_SIZES_KEYS.desktop]
 * @returns {GridLayout}
 */
export const widgetsToLayoutBySize = (
  widgets = [],
  oldLayout = [],
  size = WIDGET_GRID_SIZES_KEYS.desktop,
) => widgets.map((widget) => {
  const oldLayoutItem = oldLayout.find(({ i }) => i === widget._id);
  const newLayoutItem = oldLayoutItem ? omit(oldLayoutItem, ['i', 'widget']) : { ...widget.grid_parameters[size] };

  newLayoutItem.i = widget._id;
  newLayoutItem.widget = widget;
  newLayoutItem.h = newLayoutItem.h || 1;

  return newLayoutItem;
});

/**
 * Get y positions for every sizes for new widget
 *
 * @param {Widget[]} [widgets = []]
 * @returns {{ mobile: number, tablet: number, desktop: number }}
 */
export const calculateNewWidgetGridParametersY = (widgets = []) => (
  widgets.reduce((acc, { grid_parameters: gridParameters }) => {
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
  }, { mobile: 0, tablet: 0, desktop: 0 })
);

/**
 * Convert widgets to layouts for all sizes with compact layout function
 *
 * @param {Widget[]} [widgets = []]
 * @param {GridLayout} [layouts = []]
 * @param {string[]} [sizes = Object.values(WIDGET_GRID_SIZES_KEYS)]
 * @returns {GridLayout}
 */
export const widgetsToLayoutsWithCompact = (
  widgets = [],
  layouts = [],
  sizes = Object.values(WIDGET_GRID_SIZES_KEYS),
) => (
  sizes.reduce((acc, size) => {
    const oldLayout = get(layouts, size, []);

    acc[size] = compactLayout(widgetsToLayoutBySize(widgets, oldLayout, size));

    return acc;
  }, {})
);

/**
 * Convert layouts to widgets grid
 *
 * @param {GridLayout} layouts
 * @returns {GridLayoutWidgetGrid[]}
 */
export const layoutsToWidgetsGrid = (layouts = []) => {
  const widgetsGridById = Object.entries(layouts).reduce((acc, [size, layout]) => {
    layout.forEach((layoutItem) => {
      if (!acc[layoutItem.i]) {
        acc[layoutItem.i] = {};
      }

      acc[layoutItem.i][size] = pick(layoutItem, ['x', 'y', 'w', 'h', 'autoHeight']);
    });

    return acc;
  }, {});

  return Object.entries(widgetsGridById).map(([id, gridParameters]) => ({
    _id: id,
    grid_parameters: gridParameters,
  }));
};
