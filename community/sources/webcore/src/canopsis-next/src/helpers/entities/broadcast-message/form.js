import { keyBy } from 'lodash';

import { DEFAULT_BROADCAST_MESSAGE_COLOR, BROADCAST_MESSAGE_VIEWS } from '@/constants';

import { convertDateToDateObject, convertDateToTimestamp } from '@/helpers/date/date';

/**
 * @typedef {Object} Broadcast
 * @property {string} message
 * @property {string} color
 * @property {number} start
 * @property {number} end
 * @property {string[]} views
 */

/**
 * @typedef {Broadcast} BroadcastForm
 * @property {Date} start
 * @property {Date} end
 */

/**
 * Convert broadcast object to broadcast form
 *
 * @param {Broadcast} broadcastMessage
 * @return {BroadcastForm}
 */
export const messageToForm = (broadcastMessage = {}) => ({
  message: broadcastMessage?.message || '',
  color: broadcastMessage?.color || DEFAULT_BROADCAST_MESSAGE_COLOR,
  start: convertDateToDateObject(broadcastMessage?.start),
  end: convertDateToDateObject(broadcastMessage?.end),
  views: [...(broadcastMessage?.views?.length ? broadcastMessage.views : Object.values(BROADCAST_MESSAGE_VIEWS))],
});

/**
 * Convert views form to optimized views array
 *
 * @param {Object} viewsMap
 * @param {Object[]} treeViews
 * @param {boolean} isChildren
 * @return {string[]|boolean}
 */
export const viewsFormToViews = (viewsMap = {}, treeViews = [], isChildren = false) => {
  const result = [];
  let hasAllItemsSelected = true;

  for (const { value, children } of treeViews) {
    if (children) {
      const childrenViews = viewsFormToViews(viewsMap, children, true);

      if (childrenViews === true) {
        result.push(value);
      } else {
        result.push(...childrenViews);
        hasAllItemsSelected = false;
      }
    } else if (viewsMap[value]) {
      result.push(value);
    } else {
      hasAllItemsSelected = false;
    }
  }

  return isChildren && hasAllItemsSelected
    ? true
    : result;
};

/**
 * Convert broadcast form to broadcast object
 *
 * @param {BroadcastForm} form
 * @param {Object[]} treeViews
 * @return {Broadcast}
 */
export const formToMessage = (form = {}, treeViews = []) => ({
  ...form,

  views: viewsFormToViews(keyBy(form.views), treeViews),
  start: convertDateToTimestamp(form.start),
  end: convertDateToTimestamp(form.end),
});

/**
 * Get selected views with all their children from tree structure
 *
 * @param {Object} viewsMap
 * @param {Object[]} views
 * @param {boolean} isChildren
 * @return {string[]}
 */
export const getViewsWithChildren = (viewsMap = {}, treeViews = [], isChildren = false) => {
  const result = [];

  for (const { value, children } of treeViews) {
    if (viewsMap[value] || isChildren) {
      if (children) {
        result.push(...getViewsWithChildren(viewsMap, children, true));
      } else {
        result.push(value);
      }
    }
  }

  return result;
};

/**
 * Prepare views by converting array to map and getting all children
 *
 * @param {string[]} views
 * @param {Object[]} treeViews
 * @return {string[]}
 */
export const prepareMessageViews = (views = [], treeViews = []) => getViewsWithChildren(keyBy(views), treeViews);
