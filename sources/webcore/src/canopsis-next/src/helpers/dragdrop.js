/**
 * Function for help with drag/drop. Return a new array with changed item.
 * @param {Array} entities
 * @param {Object} moved
 * @param {Object} added
 * @param {Object} removed
 * @param {Function} prepareEntity
 * @return {Array}
 */
export const dragDropChangePositionHandler = (entities, { moved, added, removed }, prepareEntity = v => v) => {
  const copiedEntities = [...entities];

  if (moved) {
    const [item] = copiedEntities.splice(moved.oldIndex, 1);

    copiedEntities.splice(moved.newIndex, 0, item);
  } else if (added) {
    copiedEntities.splice(added.newIndex, 0, prepareEntity(added.element));
  } else if (removed) {
    copiedEntities.splice(removed.oldIndex, 1);
  }

  return copiedEntities;
};
