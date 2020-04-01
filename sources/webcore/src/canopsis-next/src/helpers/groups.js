import uuid from '@/helpers/uuid';

export const mapViewToExport = view => ({
  ...view,
  _id: uuid(),
  exported: true,
});

export const mapGroupToExport = ({ views, name }, exportedViewIds = []) => {
  const groupId = uuid();

  return ({
    views: views.reduce((acc, { _id: viewId, ...view }) => {
      if (exportedViewIds.includes(viewId)) {
        acc.push(mapViewToExport({ ...view, group_id: groupId }));
      }

      return acc;
    }, []),
    _id: groupId,
    exported: true,
    name,
  });
};

/**
 * Prepare group and views to export data object.
 * @param {Array} groups - groups with selected views.
 * @param {Array} views - views without group.
 * @return {{ groups: Array, views: Array }}
 */
export const prepareGroupsAndViewsToImport = ({ groups, views }) => {
  const viewsIds = views.map(({ _id }) => _id);
  const groupsIds = groups.map(({ _id }) => _id);

  return {
    groups: groups.map(group => mapGroupToExport(group, viewsIds)),
    views: views.reduce((acc, { group_id: groupId, ...view }) => {
      if (!groupsIds.includes(groupId)) {
        acc.push(mapViewToExport(view));
      }

      return acc;
    }, []),
    viewsIds,
  };
};
