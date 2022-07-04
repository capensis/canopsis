import { POINT_TYPES } from '@/constants';

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
 * @param {string} type
 * @returns {Point}
 */
export const generatePoint = ({
  x,
  y,
  type,
}) => ({
  x,
  y,
  _id: uid(),
  type: type ?? '',
});

/**
 * Check is curves control point
 *
 * @param {string} type
 * @returns {boolean}
 */
export const isCurvesControl = type => type === POINT_TYPES.curvesControl;

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
    acc.push(generatePoint({
      ...calculateCenterBetweenPoint(point, nextPoint),
      type: isCurvesControl(point.type) || isCurvesControl(nextPoint.type)
        ? POINT_TYPES.curvesControl
        : POINT_TYPES.point,
    }));
  }

  return acc;
}, []);
