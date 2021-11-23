import { omit } from 'lodash';

import { DEFAULT_PERIODIC_REFRESH } from '@/constants';

import uuid from '../uuid';

import { durationWithEnabledToForm } from '../date/duration';
import { generateCopyOfViewTab, generateViewTab } from '../entities';

import { enabledToForm } from './shared/common';

/**
 * @typedef {Object} ViewTab
 * @property {string} _id
 * @property {string} title
 * @property {Widget[]} widgets
 */

/**
 * @typedef {Object} ViewGroupRequest
 * @property {string} author
 * @property {string} title
 */

/**
 * @typedef {ViewGroupRequest} ViewGroup
 * @property {string} _id
 * @property {number} created
 * @property {number} updated
 */

/**
 * @typedef {Object} ViewForm
 * @property {string} title
 * @property {string} description
 * @property {boolean} enabled
 * @property {string[]} tags
 * @property {ViewGroup} [group]
 * @property {DurationWithEnabled} periodic_refresh
 */

/**
 * @typedef {ViewGroup & ViewForm} View
 * @property {ViewTab[]} tabs
 */

/**
 * @typedef {View} ViewRequest
 * @property {string} group
 */

/**
 * @typedef {ViewGroup} ViewGroupWithViews
 * @property {View[]} views
 */

/**
 * @typedef {Object} ViewGroupWithViewsPosition
 * @property {string} _id
 * @property {string[]} views
 */

/**
 * @typedef {View} ExportedView
 * @property {boolean} exported
 */

/**
 * @typedef {ViewGroupWithViews} ExportedViewGroup
 * @property {boolean} exported
 * @property {ExportedView[]} views
 */

/**
 * @typedef {Object} ExportedViewWrapper
 * @property {View} view
 * @property {number} position
 */

/**
 * @typedef {Object} ExportedViewGroupWrapper
 * @property {ViewGroup} group
 * @property {number} position
 */

/**
 * Convert view to form
 *
 * @param {View | {}} [view = {}]
 * @returns {ViewForm}
 */
export const viewToForm = (view = {}) => ({
  title: view.title ?? '',
  description: view.description ?? '',
  enabled: enabledToForm(view.enabled),
  tags: view.tags ? [...view.tags] : [],
  group: view.group ? { ...view.group } : null,
  periodic_refresh: durationWithEnabledToForm(view.periodic_refresh ?? DEFAULT_PERIODIC_REFRESH),
});

/**
 * Convert view to request
 *
 * @param {View} view
 * @returns {ViewRequest}
 */
export const viewToRequest = (view) => {
  const request = {
    ...omit(view, ['_id', 'group', 'group_id', 'created', 'updated']),

    group: view.group._id,
  };

  if (!view.tabs) {
    request.tabs = [generateViewTab('Default')];
  }

  return request;
};

/**
 * Convert view group to request
 *
 * @param {ViewGroup} group
 * @returns {ViewRequest}
 */
export const groupToRequest = group => omit(group, ['_id', 'views', 'created', 'deletable', 'updated']);

/**
 * Convert view groups with views to view group with view positions
 *
 * @param {ViewGroupWithViews | {}} [group = {}]
 * @return {ViewGroupWithViewsPosition}
 */
export const groupWithViewsToPositions = (group = {}) => ({ _id: group._id, views: group.views.map(view => view._id) });

/**
 * Convert view groups with views to view group with view positions
 *
 * @param {ViewGroupWithViews[]} [groups = []]
 * @return {ViewGroupWithViewsPosition[]}
 */
export const groupsWithViewsToPositions = (groups = []) => groups.map(groupWithViewsToPositions);

/**
 * Convert view to exported view.
 *
 * @param {View} view
 * @return {ExportedView}
 */
export const viewToExportedView = view => ({
  ...view,

  exported: true,
});

/**
 * Convert Group to export group entity.
 *
 * @param {ViewGroupWithViews} group
 * @param {string[]} exportedViewIds
 * @return {ExportedViewGroup}
 */
export const groupToExportedGroup = (group, exportedViewIds = []) => ({
  ...group,

  exported: true,
  views: group.views.reduce((acc, view) => {
    if (exportedViewIds.includes(view._id)) {
      acc.push(viewToExportedView(view));
    }

    return acc;
  }, []),
});

/**
 * Get exported group and views
 *
 * @param {ViewGroup[]} groups - groups with selected views.
 * @param {View[]} views - views without group.
 * @return {{ groups: ExportedViewGroup[], views: ExportedView[] }}
 */
export const getExportedGroupsAndViews = ({ groups, views }) => {
  const viewsIds = views.map(({ _id }) => _id);
  const groupsIds = groups.map(({ _id }) => _id);

  return {
    groups: groups.map(group => groupToExportedGroup(group, viewsIds)),
    views: views.reduce((acc, view) => {
      if (!groupsIds.includes(view.group._id)) {
        acc.push(viewToExportedView(view));
      }

      return acc;
    }, []),
    viewsIds,
  };
};

/**
 * Get exported views wrappers
 *
 * @param {ExportedView[]} [views = []]
 * @param {ViewGroup} group
 * @param {number} groupIndex
 * @return {ExportedViewWrapper[]}
 */
export const getExportedViewsWrappersFromGroup = ({ views = [], ...group }, groupIndex) => views
  .reduce((acc, { exported, ...view }, index) => {
    if (exported) {
      acc.push({ view: { ...view, group }, path: `${groupIndex}.views.${index}` });
    }

    return acc;
  }, []);

/**
 * Get exported views wrappers from groups
 *
 * @param {ExportedViewGroup[]} [groups = []]
 * @return {ExportedViewWrapper[]}
 */
export const getExportedViewsWrappersFromGroups = (groups = []) => groups.reduce((acc, group, index) => {
  acc.push(...getExportedViewsWrappersFromGroup(group, index));

  return acc;
}, []);

/**
 * Get exported views wrappers
 *
 * @param {ExportedViewGroup[]} [groups = []]
 * @return {ExportedViewGroupWrapper[]}
 */
export const getExportedGroupsWrappers = (groups = []) => groups.reduce((acc, { exported, ...group }, index) => {
  if (exported) {
    acc.push({ group, path: index });
  }

  return acc;
}, []);

/**
 * Prepare imported view tabs
 *
 * @param {ViewTab[]} tabs
 * @return {ViewTab[]}
 */
export const prepareImportedViewTabs = (tabs = []) => tabs.map(tab => generateCopyOfViewTab(tab));

/**
 * Prepare imported views
 *
 * @param {View[]} views
 * @param {ViewGroup} [group]
 * @return {View[]}
 */
export const prepareImportedViews = (views, group) => views.map(view => ({
  ...view,
  _id: uuid(),
  group: group || view.group,
  tabs: prepareImportedViewTabs(view.tabs),
}));

/**
 * Prepare imported groups
 *
 * @param {ViewGroupWithViews[]} groups
 * @return {ViewGroupWithViews[]}
 */
export const prepareImportedGroups = groups => groups.map(({ views, ...group }) => {
  const preparedGroup = { ...group, _id: uuid() };

  return { ...preparedGroup, views: prepareImportedViews(views, preparedGroup) };
});
