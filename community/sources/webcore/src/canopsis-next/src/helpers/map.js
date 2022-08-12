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
    anchor: point.entity
      ? [halfIconSize, size]
      : [halfIconSize, halfIconSize],
  };
};
