import { cloneDeep } from 'lodash';

import { ENTITY_TYPES } from '@/constants';

/**
 * Get icon data by geo map point
 *
 * @param {MapGeoPoint} point
 * @param {number} size
 * @returns {Object}
 */
export const getGeomapMarkerIconOptions = (point, size) => {
  const halfIconSize = size / 2;
  const pixelSize = `${size}px`;

  return {
    style: {
      width: pixelSize,
      height: pixelSize,
      maxWidth: 'unset',
      maxHeight: 'unset',
    },
    size,
    anchor: point.entity ? [halfIconSize, size] : [halfIconSize, halfIconSize],
  };
};

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
export const getMapEntityText = entity =>
  entity?.type === ENTITY_TYPES.resource ? entity._id : entity?.name;

/**
 *
 * @param {Entity[]} [entities = []]
 * @returns {Object<string, TreeOfDependenciesMapEntity>}
 */
export const normalizeTreeOfDependenciesMapEntities = (
  entities = [],
  childrenKey = 'dependencies',
) =>
  entities.reduce((acc, { entity, pinned_entities: pinnedEntities }) => {
    const newEntity = {
      entity: cloneDeep(entity),
      [childrenKey]: [],
    };

    pinnedEntities.forEach((pinnedEntity) => {
      const { _id: id } = pinnedEntity;

      newEntity[childrenKey].push(id);

      if (!acc[id]) {
        acc[id] = {
          entity: cloneDeep(pinnedEntity),
        };
      }
    });

    acc[entity._id] = newEntity;

    return acc;
  }, {});
