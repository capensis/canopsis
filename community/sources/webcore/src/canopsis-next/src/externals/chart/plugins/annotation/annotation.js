/* eslint-disable no-param-reassign, no-underscore-dangle, no-confusing-arrow */

import { get, isNaN, isFinite } from 'lodash';

import annotationGenerator from 'chartjs-plugin-annotation/src/annotation';
import annotationHelpersGenerator from 'chartjs-plugin-annotation/src/helpers';

function afterDataLimits(scale) {
  const annotations = get(scale, 'chart.annotation.elements', {});

  const ranges = Object.values(annotations).filter(annotation => !!annotation._model.ranges[scale.id])
    .map(annotation => annotation._model.ranges[scale.id]);

  const min = ranges.map(range => Number(range.min))
    .reduce((a, b) => isFinite(b) && !isNaN(b) && b < a ? b : a, Number.MAX_VALUE);

  const max = ranges.map(range => Number(range.max))
    .reduce((a, b) => isFinite(b) && !isNaN(b) && b > a ? b : a, Number.MIN_VALUE);

  const pixelsPerTick = scale.height / (scale.max - scale.min);
  const sizeOfLabel = 15;
  const calculatedMin = min - Math.floor((sizeOfLabel / pixelsPerTick));
  const calculatedMax = max + Math.floor((sizeOfLabel / pixelsPerTick));

  if (
    calculatedMin <= scale.min &&
    typeof scale.options.ticks.min === 'undefined' &&
    typeof scale.options.ticks.suggestedMin === 'undefined'
  ) {
    scale.min = calculatedMin;
  }

  if (
    calculatedMax >= scale.max &&
    typeof scale.options.ticks.max === 'undefined' &&
    typeof scale.options.ticks.suggestedMax === 'undefined'
  ) {
    scale.max = calculatedMax;
  }

  if (scale.handleTickRangeOptions) {
    scale.handleTickRangeOptions();
  }
}

export default function (Chart) {
  const helpers = annotationHelpersGenerator(Chart);
  const plugin = annotationGenerator(Chart);

  function setAfterDataLimitsHook(axisOptions) {
    helpers.decorate(axisOptions, 'afterDataLimits', (previous, scale) => {
      if (previous) {
        previous(scale);
      }

      afterDataLimits(scale);
    });
  }

  return {
    ...plugin,

    beforeInit(chartInstance) {
      const { options: chartOptions } = chartInstance;

      chartInstance.annotation = {
        elements: {},
        options: helpers.initConfig(chartOptions.annotation || {}),
        onDestroy: [],
        firstRun: true,
        supported: false,
      };

      // Add the annotation scale adjuster to each scale's afterDataLimits hook
      chartInstance.ensureScalesHaveIDs();

      if (chartOptions.scales) {
        chartInstance.annotation.supported = true;

        Chart.helpers.each(chartOptions.scales.xAxes, setAfterDataLimitsHook);
        Chart.helpers.each(chartOptions.scales.yAxes, setAfterDataLimitsHook);
      }
    },
  };
}
