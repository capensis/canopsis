import {
  COLUMNS_WIDGET_TEMPLATES_TYPES,
  CUSTOM_WIDGET_TEMPLATE,
  WIDGET_TEMPLATES_TYPES,
} from '@/constants';

import { widgetColumnsToForm, formToWidgetColumns } from '../column/form';

/**
 * @typedef {'alarm' | 'entity'} WidgetTemplateType
 */

/**
 * @typedef {Object} WidgetTemplate
 * @property {string} title
 * @property {WidgetTemplateType} type
 * @property {WidgetColumn[]} [columns]
 * @property {string} [content]
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
  type: widgetTemplate.type ?? WIDGET_TEMPLATES_TYPES.alarmMoreInfos,
  columns: widgetColumnsToForm(widgetTemplate.columns),
  content: widgetTemplate.content ?? '',
});

/**
 * Convert form to widget template
 *
 * @param {WidgetTemplateForm} form
 * @param {WidgetColumnForm[]} columns
 * @param {string} content
 * @returns {WidgetTemplate}
 */
export const formToWidgetTemplate = ({ columns, content, ...form }) => {
  const widgetTemplate = form;

  if (COLUMNS_WIDGET_TEMPLATES_TYPES.includes(form.type)) {
    widgetTemplate.columns = formToWidgetColumns(columns);
  } else {
    widgetTemplate.content = content;
  }

  return widgetTemplate;
};

/**
 * Convert widget column template to form
 *
 * @param {string} [template]
 * @returns {string}
 */
export const widgetTemplateValueToForm = template => template || CUSTOM_WIDGET_TEMPLATE;

/**
 * Convert form to widget column template
 *
 * @param {string} [template]
 * @returns {string}
 */
export const formToWidgetTemplateValue = template => (template === CUSTOM_WIDGET_TEMPLATE ? '' : template);
