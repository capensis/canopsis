import { cloneDeep } from 'lodash';

import { WIDGET_COLORS_MODES } from '@/constants';

import {
  formWidgetColorsToTheme,
  tableColorsToTheme,
  themeTableColorsToForm,
  themeWidgetColorsToForm,
} from '@/helpers/entities/theme/form';

/**
 * Convert widget color overrides to a form object.
 *
 * @param {Object} [colors={}] Widget color overrides.
 * @returns {Object}
 */
export const widgetColorsToForm = (colors = {}) => ({
  mode: colors.mode ?? WIDGET_COLORS_MODES.theme,
  table: themeTableColorsToForm(colors.table),
  widget: themeWidgetColorsToForm(colors.widget),
});

/**
 * Convert widget color overrides form to the API object.
 *
 * @param {Object} [form={}] Widget color overrides form.
 * @returns {Object}
 */
export const formToWidgetColors = (form = {}) => {
  const colors = {
    mode: form.mode ?? WIDGET_COLORS_MODES.theme,
  };

  if (colors.mode === WIDGET_COLORS_MODES.custom) {
    colors.table = tableColorsToTheme(form.table);
    colors.widget = formWidgetColorsToTheme(form.widget);
  }

  return cloneDeep(colors);
};
