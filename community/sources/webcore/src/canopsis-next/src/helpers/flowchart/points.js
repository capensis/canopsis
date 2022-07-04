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

/**
 * Calculate points with ghosts points
 *
 * @param {Point[]} points
 * @returns {Point[]}
 */
export const getGhostPoints = points => points.reduce((acc, point, index) => {
  const nextIndex = index + 1;
  const nextPoint = points[nextIndex];

  if (nextPoint) {
    acc.push(generatePoint(calculateCenterBetweenPoint(point, nextPoint)));
  }

  return acc;
}, []);
