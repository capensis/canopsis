import { ENTITY_TYPES } from '@/constants';

/**
 * Get icon data by geo map point
 *
 * @param {MapGeoPoint} point
 * @param {number} size
 * @returns {Object}
 */
export const getGeomapMarkerIcon = (point, size) => {
  const halfIconSize = size / 2;
  const pixelSize = `${size}px`;

  return {
    name: point.entity ? 'location_on' : 'link',
    style: {
      width: pixelSize,
      height: pixelSize,
      maxWidth: 'unset',
      maxHeight: 'unset',
    },
    size,
    tooltipAnchor: point.entity ? [0, 0] : [0, halfIconSize],
    anchor: point.entity
      ? [halfIconSize, size]
      : [halfIconSize, halfIconSize],
  };
};

/**
 * Get text for displaying in tree of dependencies components for entity
 *
 * @param {Entity} entity
 * @returns {string}
 */
export const getTreeOfDependenciesEntityText = entity => (
  entity.type === ENTITY_TYPES.resource
    ? entity._id
    : entity.name
);
