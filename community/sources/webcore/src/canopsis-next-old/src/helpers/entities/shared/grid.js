import { GRID_SIZES } from '@/constants';

/**
 * Generates flex classes for a given grid range size.
 *
 * @param {Array<number>} [gridRangeSize] - An array of two numbers where the first number is the start of the range
 *                                          and the second number is the end of the range. Optional.
 * @returns {Array<string>} An array of two strings: the first is the offset class, and the second is the size class.
 *
 * @example
 * // Returns ['offset-xs0`, `xs12`]
 * getFlexClassesForGridRangeSize();
 *
 * @example
 * // Returns ['offset-xs2`, `xs4`]
 * getFlexClassesForGridRangeSize([2, 6]);
 */
export const getFlexClassesForGridRangeSize = (gridRangeSize) => {
  const [start, end] = gridRangeSize ?? [GRID_SIZES.min, GRID_SIZES.max];

  return [
    `offset-xs${start}`,
    `xs${end - start}`,
  ];
};
