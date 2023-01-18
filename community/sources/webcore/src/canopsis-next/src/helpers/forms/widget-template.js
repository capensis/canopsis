import { WIDGET_TYPES } from '@/constants';

import { widgetColumnsToForm, formToWidgetColumns } from './shared/widget-column';

/**
 * @typedef {Object} WidgetTemplate
 * @property {string} name
 * @property {WidgetType} type
 * @property {WidgetColumn[]} columns
 */

/**
 * @typedef {Object} WidgetTemplateForm
 * @property {string} name
 * @property {WidgetType} type
 * @property {WidgetColumnForm[]} columns
 */

/**
 * Convert widget template to form
 *
 * @param {WidgetTemplate} widgetTemplate
 * @returns {WidgetTemplateForm}
 */
export const widgetTemplateToForm = (widgetTemplate = {}) => ({
  name: widgetTemplate.name ?? '',
  type: widgetTemplate.type ?? WIDGET_TYPES.alarmList,
  columns: widgetColumnsToForm(widgetTemplate.columns),
});

/**
 * Convert form to widget template
 *
 * @param {WidgetTemplateForm} form
 * @returns {WidgetTemplate}
 */
export const formToWidgetTemplate = form => ({
  ...form,

  columns: formToWidgetColumns(form.columns),
});
