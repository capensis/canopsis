import { omit } from 'lodash';

import { DEFAULT_PERIODIC_REFRESH } from '@/constants';

import { uuid } from '@/helpers/uuid';
import { durationWithEnabledToForm } from '@/helpers/date/duration';

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
 * @property {boolean} is_private
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
 * @typedef {Object} ExportedViewRequest
 * @property {string} _id
 */

/**
 * @typedef {ExportedViewRequest} ExportedViewGroupRequest
 * @property {ExportedViewRequest[]} views
 */

/**
 * @typedef {View} ImportedView
 * @property {boolean} imported
 */

/**
 * @typedef {ViewGroup} ImportedViewGroupWithViews
 * @property {boolean} imported
 * @property {ImportedView[] | View[]} views
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
  is_private: view.is_private ?? false,
  enabled: view.enabled ?? true,
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
export const viewToRequest = view => ({
  ...omit(view, ['_id', 'group', 'group_id', 'created', 'updated']),

  group: view.group._id,
});

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
 * Convert Group to export group entity.
 *
 * @param {ViewGroupWithViews} group
 * @param {string[]} exportedViewIds
 * @return {ExportedViewGroupRequest}
 */
export const groupToExportedGroup = (group, exportedViewIds = []) => ({
  _id: group._id,
  views: group.views.reduce((acc, view) => {
    if (exportedViewIds.includes(view._id)) {
      acc.push(view._id);
    }

    return acc;
  }, []),
});

/**
 * Get exported group and views
 *
 * @param {ViewGroup[]} groups - groups with selected views.
 * @param {View[]} views - views without group.
 * @return {{ groups: ExportedViewGroupRequest[], views: ExportedViewRequest[] }}
 */
export const exportedGroupsAndViewsToRequest = ({ groups, views }) => {
  const viewsIds = views.map(({ _id }) => _id);
  const groupsIds = groups.map(({ _id }) => _id);

  return {
    groups: groups.map(group => groupToExportedGroup(group, viewsIds)),
    views: views.reduce((acc, view) => {
      if (!groupsIds.includes(view.group._id)) {
        acc.push(view._id);
      }

      return acc;
    }, []),
  };
};

/**
 * Get all views in one collection from groups
 *
 * @param {ViewGroupWithViews[]} groups
 * @return {View[]}
 */
export const getAllViewsFromGroups = groups => groups.reduce((acc, { views }) => {
  acc.push(...views);

  return acc;
}, []);

/**
 * Prepare current groups for importing process
 *
 * @param {ViewGroupWithViews[]} groups
 * @return {ViewGroupWithViews[]}
 */
export const prepareCurrentGroupsForImporting = groups => (
  groups.map(group => ({ ...group, views: [...group.views] }))
);

/**
 * Prepare imported views
 *
 * @param {View[]} views
 * @return {ImportedView[]}
 */
export const prepareImportedViews = views => views.map(view => ({
  ...view,

  imported: true,
  _id: uuid(),
}));

/**
 * Prepare imported groups
 *
 * @param {ViewGroupWithViews[]} groups
 * @return {ImportedViewGroupWithViews[]}
 */
export const prepareImportedGroups = groups => groups.map(({ views, ...group }) => ({
  ...group,

  imported: true,
  _id: uuid(),
  views: prepareImportedViews(views),
}));

export const prepareViewsForImportRequest = views => views.map(view => (
  view.imported
    ? omit(view, ['_id', 'imported'])
    : view
));

/**
 * Prepare imported groups
 *
 * @param {ImportedViewGroupWithViews[] | ViewGroupWithViews[]} groups
 * @return {ViewGroupWithViews[]}
 */
export const prepareViewGroupsForImportRequest = groups => groups.map((group) => {
  let preparedGroup = group;

  if (group.imported) {
    preparedGroup = omit(group, ['_id', 'imported']);
  }

  return {
    ...preparedGroup,

    views: prepareViewsForImportRequest(group.views),
  };
});
