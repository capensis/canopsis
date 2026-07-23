/**
 * @typedef {Object} TreeOfDependenciesMapEntity
 * @property {Entity} entity
 * @property {string[]} [dependencies]
 */

/**
 * Get text for displaying in map components for entity.
 *
 * @param {Object} entity
 * @returns {string}
 */
export const getMapEntityText = (entity) =>
  entity?.type === 'resource' ? entity._id : entity?.name;
