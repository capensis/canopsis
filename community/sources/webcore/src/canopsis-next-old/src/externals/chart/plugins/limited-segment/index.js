import { isObject } from 'lodash';

/**
 * @type {Object} ChartLimitOptions
 * @property {number} value
 * @property {string} backgroundColor
 * @property {string} borderColor
 * @property {number} borderWidth
 * @property {Array} borderDash
 */

export const ChartLimitedSegmentPlugin = {
  getPointValue(point) {
    return isObject(point) ? point.y : point;
  },

  clipArea(ctx, area) {
    ctx.save();
    ctx.beginPath();
    ctx.rect(area.left, area.top, area.right - area.left, area.bottom - area.top);
    ctx.clip();
  },

  unclipArea(ctx) {
    ctx.restore();
  },

  /**
   * Function for draw limited segments
   * @param chart
   */
  afterDatasetsDraw(chart) {
    const { ctx, chartArea } = chart;
    const { datasets } = chart.data;
    const { limit = {} } = chart.options.plugins || {};

    if (!limit.enabled) {
      return;
    }

    this.clipArea(ctx, chartArea);

    datasets.forEach((dataset, datasetIndex) => {
      const meta = chart.getDatasetMeta(datasetIndex);
      const xAxis = chart.scales[limit.scaleID];

      if (!meta.hidden && limit && dataset.data) {
        let openedPath = false;

        dataset.data.forEach((pointData, dataIndex, points) => {
          const point = meta.data[dataIndex];
          const nextPoint = meta.data[dataIndex + 1];

          if (!point) {
            return;
          }

          const pointValue = this.getPointValue(pointData);
          const nextPointValue = nextPoint && this.getPointValue(points[dataIndex + 1]);
          const isClosablePath = (pointValue <= limit.value && nextPointValue <= limit.value) || !nextPoint;

          if (openedPath && isClosablePath) {
            /* right bottom corner */
            ctx.lineTo(point.x, xAxis.top);

            ctx.closePath();
            ctx.fillStyle = limit.backgroundColor;
            ctx.fill();

            openedPath = false;
          }

          if (!nextPoint) {
            return;
          }

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
          }
        });

        this.unclipArea(ctx);
      }
    });
  },
};
