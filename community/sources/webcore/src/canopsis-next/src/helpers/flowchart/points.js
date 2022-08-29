import uid from '@/helpers/uid';

/**
 * @typedef {Object} Point
 * @property {number} x
 * @property {number} y
 * @property {string} [_id]
 */

/**
 * Generate point
 *
 * @param {number} x
 * @param {number} y
 * @returns {Point}
 */
export const generatePoint = ({
  x,
  y,
}) => ({
  x,
  y,
  _id: uid(),
});

/**
 * Calculate center between points
 *
 * @param {Point} firstPoint
 * @param {Point} secondPoint
 * @returns {Point}
 */
export const calculateCenterBetweenPoint = (firstPoint, secondPoint) => ({
  x: (firstPoint.x + secondPoint.x) / 2,
  y: (firstPoint.y + secondPoint.y) / 2,
});
