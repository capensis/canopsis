/**
 * Function for help with drag/drop. Return a new array with changed item.
 * @param {Array} entities
 * @param {Object} moved
 * @param {Object} added
 * @param {Object} removed
 * @return {Array}
 */
export const dragDropChangePositionHandler = (entities, { moved, added, removed }) => {
  const copiedEntities = [...entities];

  if (moved) {
    const [item] = copiedEntities.splice(moved.oldIndex, 1);

    copiedEntities.splice(moved.newIndex, 0, item);
  } else if (added) {
    copiedEntities.splice(added.newIndex, 0, added.element);
  } else if (removed) {
    copiedEntities.splice(removed.oldIndex, 1);
  }

  return copiedEntities;
};

/**
 *
 * @param {Object[]} [entities = {}]
 * @param {Object[]} [anotherEntities = []]
 * @param {string | number} [key = '_id]
 * @param {Function} [callback = () => false]
 * @returns {boolean|boolean}
 */
export const isDeepOrderChanged = (
  entities = [],
  anotherEntities = [],
  key = '_id',
  callback = () => false,
) => entities.length !== anotherEntities.length || entities.some((entity, index) => {
  const anotherEntity = anotherEntities[index] || {};

  return entity[key] !== anotherEntity[key] || callback(entity, anotherEntity);
});
