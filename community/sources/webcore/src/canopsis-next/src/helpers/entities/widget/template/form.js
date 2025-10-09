import {
  COLUMNS_WIDGET_TEMPLATES_TYPES,
  QUICK_ACTIONS_WIDGET_TEMPLATES_TYPES,
  CUSTOM_WIDGET_TEMPLATE,
  WIDGET_TEMPLATES_TYPES,
} from '@/constants';

import { widgetColumnsToForm, formToWidgetColumns } from '../column/form';
import { widgetQuickActionsToForm, formToWidgetQuickActions } from '../quick-action/form';
/**
 * @typedef {
 *  'alarm_columns'
 *  | 'entity_columns'
 *  | 'alarm_more_infos'
 *  | 'alarm_export_to_pdf'
 *  | 'alarm_quick_actions'
 *  | 'alarm_mass_quick_actions'
 *  | 'weather_item'
 *  | 'weather_modal'
 *  | 'weather_entity'
 * } WidgetTemplateType
 */

/**
 * @typedef {Object} WidgetTemplate
 * @property {string} title
 * @property {WidgetTemplateType} type
 * @property {WidgetColumn[]} [columns]
 * @property {WidgetQuickAction[]} [actions]
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
  actions: widgetQuickActionsToForm(widgetTemplate.actions),
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
export const formToWidgetTemplate = ({ columns, content, actions, ...form }) => {
  const widgetTemplate = form;

  if (COLUMNS_WIDGET_TEMPLATES_TYPES.includes(form.type)) {
    widgetTemplate.columns = formToWidgetColumns(columns);
  } else if (QUICK_ACTIONS_WIDGET_TEMPLATES_TYPES.includes(form.type)) {
    widgetTemplate.actions = formToWidgetQuickActions(actions);
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
