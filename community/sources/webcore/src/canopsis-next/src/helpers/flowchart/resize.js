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
}) => {
  const newRect = {
    x: rect.x,
    y: rect.y,
    width: rect.width,
    height: rect.height,
  };
  const directionArray = direction.split('');

  directionArray.forEach((singleDirection) => {
    switch (singleDirection) {
      case DIRECTIONS.south: {
        const newHeight = cursorY - newRect.y;

        if (newHeight >= 0) {
          newRect.height = newHeight;
        }

        break;
      }
      case DIRECTIONS.north: {
        const newHeight = newRect.height + newRect.y - cursorY;

        if (newHeight >= 0) {
          newRect.height = newHeight;
          newRect.y = cursorY;
        }

        break;
      }
      case DIRECTIONS.east: {
        const newWidth = cursorX - newRect.x;

        if (newWidth >= 0) {
          newRect.width = newWidth;
        }

        break;
      }
      case DIRECTIONS.west: {
        const newWidth = newRect.width + newRect.x - cursorX;

        if (newWidth >= 0) {
          newRect.width = newWidth;
          newRect.x = cursorX;
        }

        break;
      }
    }
  });

  return newRect;
};

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
    case DIRECTIONS.southWest: {
      const newWidth = rect.width + rect.x - cursorX;
      newRect.height = Math.abs(newWidth * ratio);

      if (newWidth >= 0) {
        newRect.width = newWidth;
        newRect.x = cursorX;
      }

      break;
    }
    case DIRECTIONS.east:
    case DIRECTIONS.northEast: {
      const newWidth = cursorX - rect.x;

      newRect.height = Math.abs(newWidth * ratio);

      if (newWidth >= 0) {
        newRect.y += rect.height - newRect.height;
        newRect.width = newWidth;
      }

      break;
    }
    case DIRECTIONS.south:
    case DIRECTIONS.southEast: {
      const newHeight = cursorY - rect.y;
      newRect.width = Math.abs(newHeight / ratio);

      if (newHeight >= 0) {
        newRect.height = newHeight;
      }

      break;
    }
    case DIRECTIONS.north:
    case DIRECTIONS.northWest: {
      const newHeight = rect.y + rect.height - cursorY;
      newRect.width = Math.abs(newHeight / ratio);

      if (newHeight >= 0) {
        newRect.height = newHeight;
        newRect.y += rect.height - newHeight;
        newRect.x += rect.width - newRect.width;
      }

      break;
    }
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
