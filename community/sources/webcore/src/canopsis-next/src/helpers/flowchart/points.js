import { SHAPES } from '@/constants';

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
 * Check is area include point
 *
 * @param {Point} start
 * @param {Point} end
 * @param {Point} point
 * @returns {boolean}
 */
export const isAreaIncludePoint = (start, end, point) => point.x > start.x
    && point.x < end.x
    && point.y > start.y
    && point.y < end.y;

/**
 * Check is area include point
 *
 * @param {Point} start
 * @param {Point} end
 * @param {Shape} shape
 * @returns {boolean}
 */
export const isAreaIncludeShape = (start, end, shape) => {
  switch (shape.type) {
    case SHAPES.storage:
    case SHAPES.parallelogram:
    case SHAPES.image:
    case SHAPES.rhombus:
    case SHAPES.ellipse:
    case SHAPES.process:
    case SHAPES.document:
    case SHAPES.circle:
    case SHAPES.rect: {
      const leftTopCorner = { x: shape.x, y: shape.y };
      const bottomRightCorner = {
        x: shape.x + (shape.width ?? shape.diameter),
        y: shape.y + (shape.height ?? shape.diameter),
      };

      return isAreaIncludePoint(start, end, leftTopCorner)
        && isAreaIncludePoint(start, end, bottomRightCorner);
    }
    case SHAPES.arrowLine:
    case SHAPES.bidirectionalArrowLine:
    case SHAPES.line:
      return shape.points.every(
        point => isAreaIncludePoint(start, end, point),
      );
  }

  return false;
};
