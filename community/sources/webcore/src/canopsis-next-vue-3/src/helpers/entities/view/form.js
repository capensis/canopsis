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
 * @typedef {Object} ExportedViewRequest
 * @property {string} _id
 */

/**
 * @typedef {View} ImportedView
 * @property {boolean} imported
 */

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
 * @typedef {Object} ViewGroupWithViewsPosition
 * @property {string} _id
 * @property {string[]} views
 */

/**
 * Convert view groups with views to view group with view positions
 *
 * @param {ViewGroupWithViews | {}} [group = {}]
 * @return {ViewGroupWithViewsPosition}
 */
export const groupWithViewsToPositions = (group = {}) => ({
  _id: group._id,
  views: group.views.map((view) => view._id),
});

/**
 * Convert view groups with views to view group with view positions
 *
 * @param {ViewGroupWithViews[]} [groups = []]
 * @return {ViewGroupWithViewsPosition[]}
 */
export const groupsWithViewsToPositions = (groups = []) => groups.map(groupWithViewsToPositions);
