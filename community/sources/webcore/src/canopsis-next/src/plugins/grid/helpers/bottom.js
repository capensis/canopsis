/**
 * Calculate rows count
 *
 * @param {Object} layout
 * @return {number}
 */
export const bottom = (layout) => {
  let max = 0;
  let bottomY;

  for (let i = 0, len = layout.length; i < len; i += 1) {
    bottomY = layout[i].y + layout[i].h;

    if (bottomY > max) {
      max = bottomY;
    }
  }
  return max;
};
