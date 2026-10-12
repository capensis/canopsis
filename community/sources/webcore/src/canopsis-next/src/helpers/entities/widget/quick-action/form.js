import { uid } from '@/helpers/uid';

/**
 * @typedef {string} WidgetQuickAction
 */

/**
 * @typedef {ObjectKey} WidgetQuickActionForm
 * @property {WidgetQuickAction} value
 */

/**
 * Converts one quick action to its editable form representation.
 *
 * @param {WidgetQuickAction} [quickAction='']
 * @returns {WidgetQuickActionForm}
 */
export const widgetQuickAlarmActionToForm = (quickAction = '') => ({
  key: uid(),
  value: quickAction,
});

/**
 * Converts widget quick actions to editable form items.
 *
 * @param {WidgetQuickAction[]} [quickActions=[]]
 * @returns {WidgetQuickActionForm[]}
 */
export const widgetQuickActionsToForm = (quickActions = []) =>
  quickActions.map(widgetQuickAlarmActionToForm);

/**
 * Converts editable quick action items back to widget values.
 *
 * @param {WidgetQuickActionForm[]} [form=[]]
 * @returns {WidgetQuickAction[]}
 */
export const formToWidgetQuickActions = (form = []) => form.map(({ value }) => value);
