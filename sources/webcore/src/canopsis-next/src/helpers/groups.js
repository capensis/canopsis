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
    groups: groups.map(({ views: groupViews, _id: groupId, ...group }) => ({
      views: groupViews.filter(({ _id: viewId }) => viewsIds.includes(viewId)),
      ...group,
    })),
    views: views.filter(({ group_id: groupId }) => !groupsIds.includes(groupId)),
  };
};
