import { isObject } from 'lodash';

export const limitedSegmentPlugin = {
  getPointValue(point) {
    return isObject(point) ? point.y : point;
  },

  /**
   * Function for draw limited segments
   * @param chart
   */
  beforeDatasetsDraw(chart) {
    const { ctx } = chart.chart;
    const { datasets } = chart.data;

    ctx.save();

    datasets.forEach((dataset, datasetIndex) => {
      const meta = chart.getDatasetMeta(datasetIndex);
      const xAxis = chart.scales[meta.xAxisID];

      if (!meta.hidden && dataset.limit && dataset.data) {
        let openedPath = false;

        dataset.data.forEach((pointData, dataIndex, points) => {
          const point = meta.data[dataIndex];
          const nextPoint = meta.data[dataIndex + 1];

          if (!point || !nextPoint) {
            return;
          }

          const pointValue = this.getPointValue(pointData);
          const nextPointValue = this.getPointValue(points[dataIndex + 1]);
          /* eslint-disable no-underscore-dangle */
          const pointPosition = point._view;
          const nextPointPosition = nextPoint._view;
          /* eslint-enable no-underscore-dangle */

          if (pointValue > dataset.limit.value || nextPointValue > dataset.limit.value) {
            if (!openedPath) {
              ctx.beginPath();

              /* Left bottom corner */
              ctx.moveTo(pointPosition.x, xAxis.top);

              /* Left top corner */
              ctx.lineTo(pointPosition.x, pointPosition.y);

              openedPath = true;
            }

            /* Right top corner */
            if (pointPosition.steppedLine === true) {
              ctx.lineTo(pointPosition.x, pointPosition.y);
              ctx.lineTo(nextPointPosition.x, nextPointPosition.y);
            } else if (nextPointPosition.tension === 0) {
              ctx.lineTo(nextPointPosition.x, nextPointPosition.y);
            } else {
              ctx.bezierCurveTo(
                pointPosition.controlPointNextX,
                pointPosition.controlPointNextY,
                nextPointPosition.controlPointPreviousX,
                nextPointPosition.controlPointPreviousY,
                nextPointPosition.x,
                nextPointPosition.y,
              );
            }
          } else if (openedPath && pointValue < dataset.limit.value) {
            /* right bottom corner */
            ctx.lineTo(pointPosition.x, xAxis.top);

            ctx.closePath();
            ctx.fillStyle = dataset.limit.backgroundColor;
            ctx.fill();

            openedPath = false;
          }
        });

        ctx.restore();
      }
    });
  },
};
