import { ENTITIES_TYPES } from '@/constants';

import { widgetColumnsToForm, formToWidgetColumns } from './shared/widget-column';

/**
 * @typedef {'alarm' | 'entity'} WidgetTemplateType
 */

/**
 * @typedef {Object} WidgetTemplate
 * @property {string} title
 * @property {WidgetTemplateType} type
 * @property {WidgetColumn[]} columns
 */

/**
 * @typedef {WidgetTemplate} WidgetTemplateForm
 * @property {WidgetColumnForm[]} columns
 */

/**
 * Convert widget template to form
 *
 * @param {WidgetTemplate} widgetTemplate
 * @returns {WidgetTemplateForm}
 */
export const widgetTemplateToForm = (widgetTemplate = {}) => ({
  title: widgetTemplate.title ?? '',
  type: widgetTemplate.type ?? ENTITIES_TYPES.alarm,
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
