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
    anchor: point.entity
      ? [halfIconSize, size]
      : [halfIconSize, halfIconSize],
  };
};
