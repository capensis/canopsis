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
 * @returns {{ rect: Rect, direction: string }}
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

  directionArray.forEach((singleDirection, index) => {
    switch (singleDirection) {
      case DIRECTIONS.south: {
        const newHeight = cursorY - newRect.y;

        if (newHeight > 0) {
          newRect.height = newHeight;
        } else {
          newRect.height = Math.abs(newHeight);
          newRect.y -= newRect.height;

          directionArray[index] = DIRECTIONS.north;
        }

        break;
      }
      case DIRECTIONS.north: {
        const newHeight = newRect.height + newRect.y - cursorY;

        if (newHeight > 0) {
          newRect.height = newHeight;
          newRect.y = cursorY;
        } else {
          newRect.height = Math.abs(newHeight);
          newRect.y = cursorY - newRect.height;

          directionArray[index] = DIRECTIONS.south;
        }

        break;
      }
      case DIRECTIONS.east: {
        const newWidth = cursorX - newRect.x;

        if (newWidth > 0) {
          newRect.width = newWidth;
        } else {
          newRect.width = Math.abs(newWidth);
          newRect.x = cursorX;

          directionArray[index] = DIRECTIONS.west;
        }

        break;
      }
      case DIRECTIONS.west: {
        const newWidth = newRect.width + newRect.x - cursorX;

        if (newWidth > 0) {
          newRect.width = newWidth;
          newRect.x = cursorX;
        } else {
          newRect.width = Math.abs(newWidth);
          newRect.x = cursorX - newRect.width;

          directionArray[index] = DIRECTIONS.east;
        }

        break;
      }
    }
  });

  return {
    rect: newRect,
    direction: directionArray.join(''),
  };
};

/**
 * Resize rect by direction and cursor position
 *
 * @param {Rect} rect
 * @param {number} cursorX
 * @param {number} cursorY
 * @param {string} direction
 * @param {number} ratio
 * @returns {{ rect: Rect, direction: string }}
 */
export const resizeRectangleWithAspectRation = ({
  rect,
  cursorX,
  cursorY,
  direction,
  ratio,
}) => {
  const newRect = { ...rect };
  let newDirection = direction;

  switch (direction) {
    case DIRECTIONS.west:
    case DIRECTIONS.southWest: {
      const newWidth = rect.width + rect.x - cursorX;
      newRect.height = Math.abs(newWidth * ratio);

      if (newWidth >= 0) {
        newRect.width = newWidth;
        newRect.x = cursorX;
      } else {
        const absoluteNewWidth = Math.abs(newWidth);

        newRect.x += rect.width;
        newRect.y -= newRect.height;
        newRect.width = absoluteNewWidth;

        newDirection = DIRECTIONS.northEast;
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
      } else {
        const absoluteNewWidth = Math.abs(newWidth);

        newRect.y += rect.height;
        newRect.x -= absoluteNewWidth;
        newRect.width = absoluteNewWidth;

        newDirection = DIRECTIONS.southWest;
      }

      break;
    }
    case DIRECTIONS.south:
    case DIRECTIONS.southEast: {
      const newHeight = cursorY - rect.y;
      newRect.width = Math.abs(newHeight / ratio);

      if (newHeight >= 0) {
        newRect.height = newHeight;
      } else {
        const absoluteNewHeight = Math.abs(newHeight);

        newRect.height = absoluteNewHeight;
        newRect.x -= newRect.width;
        newRect.y -= absoluteNewHeight;

        newDirection = DIRECTIONS.northWest;
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
      } else {
        const absoluteNewHeight = Math.abs(newHeight);

        newRect.y += rect.height;
        newRect.x += rect.width;
        newRect.height = absoluteNewHeight;

        newDirection = DIRECTIONS.southEast;
      }

      break;
    }
  }

  return {
    rect: newRect,
    direction: newDirection,
  };
};

/**
 * Resize rect by direction and cursor position
 *
 * @param {Rect} rect
 * @param {number} cursorX
 * @param {number} cursorY
 * @param {string} direction
 * @param {number} ratio
 * @param {boolean} [aspectRatio]
 * @returns {{ rect: Rect, direction: string }}
 */
export const resizeRectangleShape = ({
  rect,
  cursorX,
  cursorY,
  direction,
  ratio,
  aspectRatio,
}) => {
  const func = aspectRatio ? resizeRectangleWithAspectRation : resizeRectangleShapeWithoutAspectRatio;

  return func({
    rect,
    cursorX,
    cursorY,
    direction,
    ratio,
  });
};
