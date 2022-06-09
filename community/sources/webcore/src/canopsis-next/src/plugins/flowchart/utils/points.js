/**
 * Calculate points with ghosts points
 *
 * @param {Array} points
 * @returns {Array}
 */
export const getPointsWithGhosts = points => points.reduce((acc, point, index) => {
  const nextIndex = index + 1;
  const nextPoint = points[nextIndex];

  acc.push({ ...point, index });

  if (nextPoint) {
    acc.push({
      x: (point.x + nextPoint.x) / 2,
      y: (point.y + nextPoint.y) / 2,
      ghost: true,
      index: nextIndex,
    });
  }

  return acc;
}, []);
