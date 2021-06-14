import uuid from '@/helpers/uuid';

/**
 * Convert view to export view entity.
 * @param {Object} view
 * @return {Object}
 */
export const prepareViewToExport = view => ({
  ...view,
  _id: uuid(),
  exported: true,
});

/**
 * Convert Group to export group entity.
 * @param {Array} views
 * @param {string} name
 * @param {Object} rest
 * @param {Array} exportedViewIds
 * @return {Object}
 */
export const prepareGroupToExport = ({ views, name, ...rest }, exportedViewIds = []) => {
  const groupId = uuid();

  return {
    ...rest,

    name,
    _id: groupId,
    exported: true,
    views: views.reduce((acc, { _id: viewId, ...view }) => {
      if (exportedViewIds.includes(viewId)) {
        acc.push(prepareViewToExport({ ...view, group_id: groupId }));
      }

      return acc;
    }, []),
  };
};

/**
 * Prepare group and views to export data object.
 * @param {Array} groups - groups with selected views.
 * @param {Array} views - views without group.
 * @return {{ groups: Array, views: Array }}
 */
export const prepareGroupsAndViewsToExport = ({ groups, views }) => {
  const viewsIds = views.map(({ _id }) => _id);
  const groupsIds = groups.map(({ _id }) => _id);

  return {
    groups: groups.map(group => prepareGroupToExport(group, viewsIds)),
    views: views.reduce((acc, { group_id: groupId, ...view }) => {
      if (!groupsIds.includes(groupId)) {
        acc.push(prepareViewToExport(view));
      }

      return acc;
    }, []),
    viewsIds,
  };
};
