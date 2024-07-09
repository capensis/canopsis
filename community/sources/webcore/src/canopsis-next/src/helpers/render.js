import { chunk } from 'lodash';

/**
 * A recursive function that calls the provided callback function using requestAnimationFrame.
 *
 * @param {Function} [callback] - The callback function to be called recursively.
 * @param {number} [depth = 0] - The depth of recursion.
 */
export const recursiveRaf = (callback, depth = 0) => {
  if (depth <= 0) {
    callback?.();
    return;
  }

  window.requestAnimationFrame(() => recursiveRaf(callback, depth - 1));
};

/**
 * Calculate the nearest viewport indexes bound based on the provided parameters.
 *
 * @param {Object} options - The options object.
 * @param {HTMLElement} options.wrapperEl - The table body element.
 * @param {number} [options.length = 0] - The array of ids.
 * @param {number} [options.height = 24] - The height of each item.
 * @param {number} [options.threshold = 0.2] - The threshold value.
 * @returns {{ start: number, end: number }} The start and end indexes for the viewport.
 */
export const getNearestViewportIndexesBound = ({
  wrapperEl,
  length = 0,
  height = 24,
  threshold = 0.2,
} = {}) => {
  let start = 0;
  let end = length - 1;

  if (!wrapperEl || !length) {
    return { start, end };
  }

  const itemsPerViewport = Math.round(window.innerHeight / height);
  const batchSize = Math.round(itemsPerViewport + (itemsPerViewport * threshold * 2));
  const halfBatchSize = Math.round(batchSize / 2);

  end = Math.min(batchSize, end);

  const originalMiddleOfVisible = window.scrollY + (window.innerHeight / 2);
  let middleOfVisible = originalMiddleOfVisible;

  const { top: wrapperTop } = wrapperEl.getBoundingClientRect();

  wrapperEl.querySelectorAll('.v-data-table__expanded__content').forEach((element) => {
    const { height: expandPanelHeight, top: expandPanelTop } = element.getBoundingClientRect();

    if (expandPanelTop < originalMiddleOfVisible) {
      middleOfVisible -= expandPanelHeight;
    }
  });

  if (wrapperTop < window.scrollY) {
    const middleItemIndex = Math.abs(Math.round((middleOfVisible - (wrapperTop + window.scrollY)) / height));

    const startDiff = middleItemIndex - halfBatchSize;

    start = Math.max(startDiff, 0);
    end = Math.min(middleItemIndex + halfBatchSize + (startDiff < 0 ? -startDiff : 0), length - 1);
  }

  return { start, end };
};

/**
 * Get nearest and farthest indexes from a given range.
 *
 * @param {Object} options - The options object.
 * @param {number} [options.start = 0] - The starting index.
 * @param {number} [options.end = 0] - The ending index.
 * @param {Array<string | number>} [options.ids = []] - The array of ids.
 * @returns {{ nearest: Array<string | number>, farthest: Array<string | number> }} An object containing nearest and
 * farthest indexes.
 */
export const getNearestAndFarthestIndexes = ({
  start = 0,
  end = 0,
  ids = [],
} = {}) => {
  const { length } = ids;
  const nearest = [];
  const farthest = [];

  for (let i = start; i <= end; i += 1) {
    nearest.push(ids[i]);
  }

  for (let i = start - 1, j = end + 1; i >= 0 || j < length; i -= 1, j += 1) {
    if (i >= 0) {
      farthest.push(ids[i]);
    }

    if (j <= length) {
      farthest.push(ids[j]);
    }
  }

  return { nearest, farthest };
};

/**
 * Splits an array of ids into chunks of a specified length.
 *
 * @param {Array} [ids = []] - The array of ids to split into chunks.
 * @param {number} [length = 10] - The length of each chunk.
 * @returns {Array} An array of chunks containing the ids.
 */
export const splitIdsToChunk = (ids = [], length = 10) => chunk(ids, length);
