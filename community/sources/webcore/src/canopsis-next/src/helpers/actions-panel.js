import { ALARM_LIST_ACTIONS_TYPES, ALARM_LIST_TOGGLE_ACTIONS_TYPES_MAP } from '@/constants';

/**
 * Finds the index of an action in the quickActions array.
 * Special handling is included for link-type actions and toggle actions.
 *
 * @param {Array<string>} [quickActions=[]] - Array of quick action types
 * @param {Object} action - The action object to find
 * @param {string} [action.type] - The type of the action
 * @param {boolean} [action.link] - Whether the action is a link type
 * @returns {number} The index of the action in quickActions array, or -1 if not found
 */
export const findQuickActionIndex = (quickActions = [], action) => quickActions.findIndex((quickAction) => {
  if (quickAction === ALARM_LIST_ACTIONS_TYPES.links) {
    return action.link;
  }

  return quickAction === action.type || ALARM_LIST_TOGGLE_ACTIONS_TYPES_MAP[quickAction] === action.type;
});

/**
 * Sorts an array of actions based on their position in the quickActions array.
 * Actions that are not in quickActions are placed at the end.
 *
 * @param {Array<Object>} [actions=[]] - Array of action objects to sort
 * @param {Array<string>} [quickActions=[]] - Array of quick action types that defines the sort order
 * @returns {Array<Object>} A new sorted array of actions
 */
export const sortActionsByQuickActions = (actions = [], quickActions = []) => [...actions].sort((a, b) => {
  const aIndex = findQuickActionIndex(quickActions, a);
  const bIndex = findQuickActionIndex(quickActions, b);

  if (aIndex === -1) {
    if (bIndex === -1) {
      return 0;
    }

    return 1;
  }

  if (bIndex === -1) {
    return -1;
  }

  return aIndex - bIndex;
});

/**
 * Calculates the number of inline actions to display in the actions panel.
 *
 * @param {Array<Object>} actions - The list of all possible actions.
 * @param {Array<string>} quickActions - The list of quick actions to match against.
 * @returns {number} The number of inline actions to display (including menu button if needed).
 */
export const getActionsInlineCount = (actions = [], quickActions = []) => {
  const filteredActions = actions.filter(action => (
    findQuickActionIndex(quickActions, action) !== -1
  ));

  if (!filteredActions.length) {
    return 0;
  }

  if (filteredActions.length === actions.length) {
    return filteredActions.length;
  }

  /**
   * We need to add +1 because we have a menu button in the actions panel
   */
  return filteredActions.length + 1;
};
