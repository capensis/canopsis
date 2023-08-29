import { isNumber, uniq } from 'lodash';

/**
 * @typedef GridLayoutItem
 * @property {number | string} i
 * @property {number} x
 * @property {number} y
 * @property {number} w
 * @property {number} h
 * @property {boolean} [moved]
 * @property {boolean} [autoHeight]
 * @property {Widget} [widget]
 */

/**
 * @typedef {GridLayoutItem[]} GridLayout
 */

/**
 * Calculate the bottom available Y value
 *
 * @param {GridLayout} [layout = []]
 * @returns {number}
 */
export const calculateLayoutBottom = (layout = []) => layout.reduce((acc, item) => Math.max(acc, item.y + item.h), 0);

/**
 * Calculate count of the unique rows with content
 *
 * @param {GridLayout} [layout = []]
 * @returns {number}
 */
export const calculateLayoutRowsCount = (layout = []) => uniq(layout.map(({ y }) => y)).length;

/**
 * Given two layout items, check if they collide.
 *
 * @param {GridLayoutItem} first
 * @param {GridLayoutItem} second
 * @returns {boolean}
 */
const collides = (first, second) => !(
  first.i === second.i // same element
  || first.x + first.w <= second.x // first is left of second
  || first.x >= second.x + second.w // first is right of second
  || first.y + first.h <= second.y // first is above second
  || first.y >= second.y + second.h // first is below second
);

/**
 * Get all collisions for layout
 *
 * @param {GridLayout} layout
 * @param {GridLayoutItem} layoutItem
 * @returns {GridLayout}
 */
const getAllCollisions = (layout, layoutItem) => layout.filter(l => collides(layoutItem, l));

/**
 * Get first collision for the item in the layout
 *
 * @param {GridLayout} layout
 * @param {GridLayoutItem} layoutItem
 * @returns {GridLayoutItem || undefined}
 */
const getFirstCollision = (layout, layoutItem) => layout.find(item => collides(item, layoutItem));

/**
 * Get sorted layout by rows and columns
 *
 * @param {GridLayout} layout
 * @returns {GridLayout}
 */
const getSortedLayout = (layout = []) => (
  [...layout].sort((a, b) => Number(a.y > b.y || (a.y === b.y && a.x > b.x)) || -1)
);

/**
 * Replace the layout item in the layout
 *
 * @param {GridLayout} layout
 * @param {GridLayoutItem} layoutItem
 * @returns {GridLayout}
 */
export const replaceLayoutItemInLayout = (layout = [], layoutItem) => (
  layout.map(item => (item.i === layoutItem?.i ? layoutItem : item))
);

/**
 * Immutable move item in the layout
 *
 * @param {GridLayout} layout
 * @param {GridLayoutItem} layoutItem
 * @param {number} x
 * @param {number} y
 * @param {boolean} isUserAction
 * @returns {GridLayout}
 */
export const moveLayoutItem = (layout, layoutItem, x, y, isUserAction) => {
  const newLayoutItem = {
    ...layoutItem,
    x,
    y,
    moved: true,
  };
  const newLayout = replaceLayoutItemInLayout(layout, newLayoutItem);
  const movingUp = y && layoutItem.y > y;

  // If this collides with anything, move it.
  // When doing this comparison, we have to sort the items we compare with
  // to ensure, in the case of multiple collisions, that we're getting the
  // nearest collision.
  const sorted = getSortedLayout(newLayout);

  if (movingUp) {
    sorted.reverse();
  }

  const collisions = getAllCollisions(sorted, newLayoutItem);

  // Move each item that collides away from this element.
  return collisions.reduce((acc, collision) => {
    if (collision.moved) {
      return acc;
    }

    // This makes it feel a bit more precise by waiting to swap for just a bit when moving up.
    if (newLayoutItem.y > collision.y && newLayoutItem.y - collision.y > collision.h / 4) {
      return acc;
    }

    if (isUserAction) {
      const fakeItem = {
        x: collision.x,
        w: collision.w,
        h: collision.h,
        i: '-1',
        y: Math.max(newLayoutItem.y - collision.h, 0),
      };

      if (!getFirstCollision(acc, fakeItem)) {
        return moveLayoutItem(acc, collision, collision.x, fakeItem.y, false);
      }
    }

    return moveLayoutItem(acc, collision, collision.x, collision.y + 1, false);
  }, newLayout);
};

/**
 * Calculate new Y for the item depends on part of layout
 *
 * @param {GridLayout} compareWith
 * @param {GridLayoutItem} layoutItem
 * @returns {GridLayoutItem}
 */
const compactLayoutItem = (compareWith, layoutItem) => {
  const newLayoutItem = { ...layoutItem };

  while (newLayoutItem.y > 0 && !getFirstCollision(compareWith, newLayoutItem)) {
    newLayoutItem.y -= 1;
  }

  for (
    let collision = getFirstCollision(compareWith, newLayoutItem);
    collision;
    collision = getFirstCollision(compareWith, newLayoutItem)
  ) {
    newLayoutItem.y = collision.y + collision.h;
  }

  return newLayoutItem;
};

/**
 * Compact layout (reducing free spaces by Y)
 *
 * @param {GridLayout} layout
 * @returns {GridLayout}
 */
export const compactLayout = (layout = []) => {
  const indexesById = layout.reduce((acc, { i }, index) => {
    acc[i] = index;

    return acc;
  }, {});

  const newLayout = new Array(layout.length);
  const sortedLayout = getSortedLayout(layout);
  const compareWith = [];

  for (const layoutItem of sortedLayout) {
    const newLayoutItem = compactLayoutItem(compareWith, layoutItem);

    newLayoutItem.moved = false;

    compareWith.push(newLayoutItem);

    newLayout[indexesById[layoutItem.i]] = newLayoutItem;
  }

  return newLayout;
};

/**
 * Get count of items which above of the item
 *
 * @param {GridLayout} layout
 * @param {number} itemX
 * @param {number} itemY
 * @param {number} itemW
 * @returns {*}
 */
export const getCountAboveItems = (layout = [], itemX, itemY, itemW) => {
  const sortedByYLayout = layout
    .filter(({ y }) => y < itemY)
    .sort((a, b) => b.y - a.y);

  let count = 0;
  let x = itemX;
  let y = itemY;
  let w = itemW;

  for (const item of sortedByYLayout) {
    if (y !== item.y + item.h) {
      continue;
    }

    const diff = (x + w) - (item.x + item.w);
    const isInteraction = diff > 0
      ? diff < w
      : Math.abs(diff) < item.w;

    if (item.y < y && isInteraction) {
      count += 1;
      x = item.x;
      y = item.y;
      w = item.w;
    }
  }

  return count;
};

/**
 * Get delta for grid item
 *
 * @param {number} prevX
 * @param {number} prevY
 * @param {number} x
 * @param {number} y
 * @returns {{ deltaX: number, deltaY: number }}
 */
export const getItemDelta = (prevX, prevY, x, y) => (
  !isNumber(prevX)
    ? { deltaX: 0, deltaY: 0 }
    : { deltaX: x - prevX, deltaY: y - prevY }
);

/**
 * Get control position by event and layout HTML element
 *
 * @param event
 * @param layoutElement
 * @returns {{x: number, y: number}}
 */
export const getControlPosition = (event, layoutElement) => {
  if (!layoutElement || !event) {
    return { x: 0, y: 0 };
  }

  const layoutElementRect = layoutElement.getBoundingClientRect();

  const x = event.clientX + layoutElement.scrollLeft - layoutElementRect.left;
  const y = event.clientY + layoutElement.scrollTop - layoutElementRect.top;

  return { x, y };
};

/**
 * Find the layout item in the layout
 *
 * @param {GridLayout} layout
 * @param {string | number} id
 * @returns {GridLayoutItem | undefined}
 */
export const findLayoutItem = (layout = [], id) => layout.find(({ i }) => i === id);
