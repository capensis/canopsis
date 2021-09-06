import { isObject } from 'lodash';

/**
 * @type {Object} ChartLimitOptions
 * @property {number} value
 * @property {string} backgroundColor
 * @property {string} borderColor
 * @property {number} borderWidth
 * @property {Array} borderDash
 */

export const limitedSegmentPlugin = {
  getPointValue(point) {
    return isObject(point) ? point.y : point;
  },

  /**
   * Function for draw limited segments
   * @param chart
   */
  afterDatasetsDraw(chart) {
    const { ctx } = chart;
    const { datasets } = chart.data;
    const { limit = {} } = chart.options.plugins || {};

    if (!limit.enabled) {
      return;
    }

    ctx.save();

    datasets.forEach((dataset, datasetIndex) => {
      const meta = chart.getDatasetMeta(datasetIndex);
      const xAxis = chart.scales[limit.scaleID];

      if (!meta.hidden && limit && dataset.data) {
        let openedPath = false;

        dataset.data.forEach((pointData, dataIndex, points) => {
          const point = meta.data[dataIndex];
          const nextPoint = meta.data[dataIndex + 1];

          if (!point || !nextPoint) {
            return;
          }

          const pointValue = this.getPointValue(pointData);
          const nextPointValue = this.getPointValue(points[dataIndex + 1]);

          if (pointValue > limit.value || nextPointValue > limit.value) {
            if (!openedPath) {
              ctx.beginPath();

              /* Left bottom corner */
              ctx.moveTo(point.x, xAxis.top);

              /* Left top corner */
              ctx.lineTo(point.x, point.y);

              openedPath = true;
            }

            /* Right top corner */
            if (point.steppedLine === true) {
              ctx.lineTo(point.x, point.y);
              ctx.lineTo(nextPoint.x, nextPoint.y);
            } else if (nextPoint.tension === 0) {
              ctx.lineTo(nextPoint.x, nextPoint.y);
            } else {
              ctx.bezierCurveTo(
                point.cp2x,
                point.cp2y,
                nextPoint.cp1x,
                nextPoint.cp1y,
                nextPoint.x,
                nextPoint.y,
              );
            }
          } else if (openedPath && pointValue < limit.value) {
            /* right bottom corner */
            ctx.lineTo(point.x, xAxis.top);

            ctx.closePath();
            ctx.fillStyle = limit.backgroundColor;
            ctx.fill();

            openedPath = false;
          }
        });

        ctx.restore();
      }
    });
  },
};
