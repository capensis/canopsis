import {
  COLUMNS_WIDGET_TEMPLATES_TYPES,
  QUICK_ACTIONS_WIDGET_TEMPLATES_TYPES,
  WIDGET_TEMPLATES_TYPES,
} from '@/constants';

import { widgetColumnsToForm } from '../column/form';
import { widgetQuickActionsToForm } from '../quick-action/form';
import { widgetSortColumnsToForm } from '../sort-column/form';

/**
 * Prepare API widget templates for settings forms.
 *
 * @param {WidgetTemplate[]} [templates=[]]
 * @returns {WidgetTemplateForm[]}
 */
export const prepareWidgetTemplates = (templates = []) => templates.map((template) => {
    if (COLUMNS_WIDGET_TEMPLATES_TYPES.includes(template.type)) {
      return { ...template, columns: widgetColumnsToForm(template.columns) };
    }

    if (QUICK_ACTIONS_WIDGET_TEMPLATES_TYPES.includes(template.type)) {
      return { ...template, actions: widgetQuickActionsToForm(template.actions) };
    }

    if (WIDGET_TEMPLATES_TYPES.alarmSortColumns === template.type) {
      return { ...template, sort_columns: widgetSortColumnsToForm(template.sort_columns) };
    }

    return template;
  });
