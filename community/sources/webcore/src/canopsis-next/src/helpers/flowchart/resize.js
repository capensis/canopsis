import { DIRECTIONS } from '@/constants';

/**
 * @typedef {Object} Rect
 * @property {number} width
 * @property {number} height
 * @property {number} x
 * @property {number} y
 */

/**
 * Resize rectangle by direction and cursor position
 *
 * @param {Rect} rect
 * @param {number} cursorX
 * @param {number} cursorY
 * @param {string} direction
 * @returns {Rect}
 */
export const resizeRectangleShapeWithoutAspectRatio = ({
  rect,
  cursorX,
  cursorY,
  direction,
}) => direction
  .split('')
  .reduce((acc, singleDirection) => {
    switch (singleDirection) {
      case DIRECTIONS.south:
        acc.height = Math.max(0, cursorY - acc.y);
        break;
      case DIRECTIONS.east:
        acc.width = Math.max(0, cursorX - acc.x);
        break;
      case DIRECTIONS.north:
        acc.height = Math.max(0, acc.height + acc.y - cursorY);
        acc.y += rect.height - acc.height;
        break;
      case DIRECTIONS.west:
        acc.width = Math.max(0, acc.width + acc.x - cursorX);
        acc.x += rect.width - acc.width;
        break;
    }

    return acc;
  }, { ...rect });

/**
 * Resize rect by direction and cursor position
 *
 * @param {Rect} rect
 * @param {number} cursorX
 * @param {number} cursorY
 * @param {string} direction
 * @param {number} ratio
 * @returns {Rect}
 */
export const resizeRectangleWithAspectRation = ({
  rect,
  cursorX,
  cursorY,
  direction,
  ratio,
}) => {
  const newRect = { ...rect };

  switch (direction) {
    case DIRECTIONS.west:
    case DIRECTIONS.southWest:
      newRect.width = Math.max(0, rect.width + rect.x - cursorX);
      newRect.height = newRect.width * ratio;
      newRect.x += rect.width - newRect.width;
      break;

    case DIRECTIONS.east:
    case DIRECTIONS.northEast:
      newRect.width = Math.max(0, cursorX - rect.x);
      newRect.height = newRect.width * ratio;
      newRect.y += rect.height - newRect.height;
      break;

    case DIRECTIONS.south:
    case DIRECTIONS.southEast:
      newRect.height = Math.max(0, cursorY - rect.y);
      newRect.width = newRect.height / ratio;
      break;

    case DIRECTIONS.north:
    case DIRECTIONS.northWest:
      newRect.height = Math.max(0, rect.y + rect.height - cursorY);
      newRect.width = newRect.height / ratio;
      newRect.y += rect.height - newRect.height;
      newRect.x += rect.width - newRect.width;
      break;
  }

  return newRect;
};

/**
 * Resize rect by direction and cursor position
 *
 * @param {Rect} rect
 * @param {number} options.cursorX
 * @param {number} options.cursorY
 * @param {string} options.direction
 * @param {number} options.ratio
 * @param {boolean} [aspectRatio]
 * @returns {Rect}
 */
export const resizeRectangleShape = ({ aspectRatio, ...options }) => {
  const func = aspectRatio
    ? resizeRectangleWithAspectRation
    : resizeRectangleShapeWithoutAspectRatio;

  return func(options);
};
