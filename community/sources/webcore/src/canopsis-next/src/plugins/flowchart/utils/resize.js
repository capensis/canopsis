import { DIRECTIONS } from '@/plugins/flowchart/constants';

/**
 * @typedef {Object} Rect
 * @property {number} width
 * @property {number} height
 * @property {number} x
 * @property {number} y
 */

/**
 * @typedef {Object} Square
 * @property {number} size
 * @property {number} x
 * @property {number} y
 */

/**
 * @typedef {Object} Cursor
 * @property {number} x
 * @property {number} y
 */

/**
 * Resize rectangle by direction and cursor position
 *
 * @param {Rect} rect
 * @param {Cursor} cursor
 * @param {string} direction
 * @returns {{ rect: Rect, direction: string }}
 */
export const resizeRectangleShapeByDirection = (rect, cursor, direction) => {
  const { x: cursorX, y: cursorY } = cursor;
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
 * Resize square by direction and cursor position
 *
 * @param {Square} square
 * @param {Cursor} cursor
 * @param {string} direction
 * @returns {{ square: Square, direction: string }}
 */
export const resizeSquareShapeByDirection = (square, cursor, direction) => {
  const { x: cursorX, y: cursorY } = cursor;
  const newSquare = { ...square };
  let newDirection = direction;

  switch (direction) {
    case DIRECTIONS.west:
    case DIRECTIONS.southWest: {
      const newSize = square.size + square.x - cursorX;

      if (newSize >= 0) {
        newSquare.size = newSize;
        newSquare.x = cursorX;
      } else {
        const absoluteNewSize = Math.abs(newSize);

        newSquare.x += square.size;
        newSquare.y -= absoluteNewSize;
        newSquare.size = absoluteNewSize;

        newDirection = DIRECTIONS.northEast;
      }

      break;
    }
    case DIRECTIONS.east:
    case DIRECTIONS.northEast: {
      const newSize = cursorX - square.x;

      if (newSize >= 0) {
        newSquare.y += square.size - newSize;
        newSquare.size = newSize;
      } else {
        const absoluteNewSize = Math.abs(newSize);

        newSquare.y += square.size;
        newSquare.x -= absoluteNewSize;
        newSquare.size = absoluteNewSize;

        newDirection = DIRECTIONS.southWest;
      }

      break;
    }
    case DIRECTIONS.south:
    case DIRECTIONS.southEast: {
      const newSize = cursorY - square.y;

      if (newSize >= 0) {
        newSquare.size = newSize;
      } else {
        const absoluteNewSize = Math.abs(newSize);

        newSquare.size = absoluteNewSize;
        newSquare.x -= absoluteNewSize;
        newSquare.y -= absoluteNewSize;

        newDirection = DIRECTIONS.northWest;
      }

      break;
    }
    case DIRECTIONS.north:
    case DIRECTIONS.northWest: {
      const newSize = square.y + square.size - cursorY;

      if (newSize >= 0) {
        const diffBetweenSizes = square.size - newSize;

        newSquare.size = newSize;
        newSquare.y += diffBetweenSizes;
        newSquare.x += diffBetweenSizes;
      } else {
        const absoluteNewSize = Math.abs(newSize);

        newSquare.y += square.size;
        newSquare.x += square.size;
        newSquare.size = absoluteNewSize;

        newDirection = DIRECTIONS.southEast;
      }

      break;
    }
  }

  return {
    square: newSquare,
    direction: newDirection,
  };
};
