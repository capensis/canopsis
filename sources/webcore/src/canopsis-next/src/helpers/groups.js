import uuid from '@/helpers/uuid';

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
      views: groupViews.reduce((acc, { _id: viewId, ...view }) => {
        if (viewsIds.includes(viewId)) {
          acc.push({
            ...view,
            _id: uuid(),
          });
        }

        return acc;
      }, []),
      _id: uuid(),
      ...group,
    })),
    views: views.reduce((acc, { group_id: groupId, ...view }) => {
      if (!groupsIds.includes(groupId)) {
        acc.push(view);
      }

      return acc;
    }, []),
  };
};
