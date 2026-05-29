import { keyBy } from 'lodash';

import {
  DEFAULT_BROADCAST_MESSAGE_COLOR,
  BROADCAST_MESSAGE_VIEWS,
  BROADCAST_MESSAGE_VIEWS_FORM_BLOCKS,
  BROADCAST_MESSAGE_PAGES_VIEW_VALUES,
  DEFAULT_BROADCAST_MESSAGE_VIEWS_FORM,
} from '@/constants';

import { convertDateToDateObject, convertDateToTimestamp } from '@/helpers/date/date';

/**
 * @typedef {Object} BroadcastMessageViewsForm
 * @property {string[]} pages
 * @property {string[]} views
 * @property {string[]} playlists
 */

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
 * @property {BroadcastMessageViewsForm} views
 */

/**
 * @param {Object[]} treeViews
 * @returns {string[]}
 */
const flattenTreeValues = (treeViews = []) => treeViews.flatMap(({ value, children }) => [
  value,
  ...(children ? flattenTreeValues(children) : []),
]);

/**
 * @param {string[]} views
 * @param {Object[]} treeViews
 * @returns {string[]}
 */
const filterViewsByTree = (views = [], treeViews = []) => {
  const treeValues = new Set(flattenTreeValues(treeViews));

  return views.filter(view => treeValues.has(view));
};

/**
 * Convert views array to views form object
 *
 * @param {string[]} views
 * @param {Object} [treeItems={}]
 * @returns {BroadcastMessageViewsForm}
 */
export const viewsArrayToViewsForm = (views = [], treeItems = {}) => {
  const result = {
    pages: filterViewsByTree(views, treeItems.pages),
    views: filterViewsByTree(views, treeItems.views),
    playlists: filterViewsByTree(views, treeItems.playlists),
  };

  const assigned = new Set([...result.pages, ...result.views, ...result.playlists]);
  const unassigned = views.filter(view => !assigned.has(view));

  unassigned.forEach((view) => {
    if (BROADCAST_MESSAGE_PAGES_VIEW_VALUES.includes(view)) {
      result.pages.push(view);
    } else if (view === BROADCAST_MESSAGE_VIEWS.allViews) {
      result.views.push(view);
    } else if (view === BROADCAST_MESSAGE_VIEWS.allPlaylists) {
      result.playlists.push(view);
    } else {
      result.views.push(view);
    }
  });

  return result;
};

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
  views: broadcastMessage?.views?.length
    ? viewsArrayToViewsForm(broadcastMessage.views)
    : {
      pages: [...DEFAULT_BROADCAST_MESSAGE_VIEWS_FORM.pages],
      views: [...DEFAULT_BROADCAST_MESSAGE_VIEWS_FORM.views],
      playlists: [...DEFAULT_BROADCAST_MESSAGE_VIEWS_FORM.playlists],
    },
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
 * Convert views form object to views array
 *
 * @param {BroadcastMessageViewsForm} viewsForm
 * @param {Object} treeItems
 * @returns {string[]}
 */
export const viewsFormToMessage = (viewsForm = {}, treeItems = {}) => (
  Object.values(BROADCAST_MESSAGE_VIEWS_FORM_BLOCKS).flatMap(block => (
    viewsFormToViews(keyBy(viewsForm[block] || []), treeItems[block] || [])
  ))
);

/**
 * Convert broadcast form to broadcast object
 *
 * @param {BroadcastForm} form
 * @param {Object} treeItems
 * @return {Broadcast}
 */
export const formToMessage = (form = {}, treeItems = {}) => ({
  ...form,

  views: viewsFormToMessage(form.views, treeItems),
  start: convertDateToTimestamp(form.start),
  end: convertDateToTimestamp(form.end),
});

/**
 * Get selected views with all their children from tree structure
 *
 * @param {Object} viewsMap
 * @param {Object[]} treeViews
 * @return {string[]}
 */
export const getViewsWithChildren = (viewsMap = {}, treeViews = []) => {
  const result = [];

  for (const { value, children } of treeViews) {
    if (viewsMap[value]) {
      result.push(value);
    }

    if (children) {
      result.push(...getViewsWithChildren(viewsMap, children));
    }
  }

  return result;
};

/**
 * Prepare views form by getting all children for each block
 *
 * @param {BroadcastMessageViewsForm} viewsForm
 * @param {Object} treeItems
 * @return {BroadcastMessageViewsForm}
 */
export const prepareMessageViews = (viewsForm = {}, treeItems = {}) => (
  Object.values(BROADCAST_MESSAGE_VIEWS_FORM_BLOCKS).reduce((acc, block) => {
    acc[block] = getViewsWithChildren(
      keyBy(viewsForm[block] || []),
      treeItems[block] || [],
    );

    return acc;
  }, {})
);
