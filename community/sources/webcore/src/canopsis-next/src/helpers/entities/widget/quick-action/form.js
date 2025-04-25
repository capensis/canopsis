import { uid } from '@/helpers/uid';

/**
 * @typedef {string} WidgetQuickAction
 */

/**
 * @typedef {ObjectKey} WidgetQuickActionForm
 * @property {WidgetQuickAction} value
 */

/**
 * Converts an array of quick actions to form-compatible format by adding unique keys
 *
 * @param {WidgetQuickAction[]} [quickActions=[]] - Array of quick actions to convert
 * @returns {WidgetQuickActionForm[]} Array of quick actions with unique keys for form usage
 */
export const widgetQuickActionsToForm = (quickActions = []) => quickActions.map(quickAction => ({
  key: uid(),
  value: quickAction,
}));

/**
 * Converts form data back to quick actions format by extracting only the values
 *
 * @param {WidgetQuickActionForm[]} [form=[]] - Array of form items to convert
 * @returns {WidgetQuickAction[]} Array of quick action values
 */
export const formToWidgetQuickActions = (form = []) => form.map(({ value }) => value);
