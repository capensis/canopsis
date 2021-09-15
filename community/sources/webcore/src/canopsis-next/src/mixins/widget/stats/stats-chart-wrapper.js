import { STATS_TYPES } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';
import { convertNumberToRoundedPercentString } from '@/helpers/string';
import { convertDurationToString } from '@/helpers/date/duration';

import widgetStatsQueryMixin from './stats-query';

export default {
  mixins: [widgetStatsQueryMixin],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: true,
      serverErrorMessage: null,
      stats: null,
    };
  },
  computed: {
    labels() {
      if (this.stats) {
        const stats = Object.keys(this.stats);

        /**
         'start' correspond to the beginning timestamp.
         It's the same for all stats, that's why we can just take the first.
         We then give it to the date filter, to display it with a date format
         */
        return this.stats[stats[0]].sum.map((value) => {
          const start = convertDateToString(value.start, 'medium');

          return [start];
        });
      }

      return [];
    },

    datasets() {
      return [];
    },

    annotationLine() {
      const { annotationLine } = this.widget.parameters;

      return {
        annotations: {
          annotationLine: {
            drawTime: 'afterDatasetsDraw',
            display: annotationLine && annotationLine.enabled,
            type: 'line',
            mode: 'horizontal',
            scaleID: 'y',
            value: annotationLine.value,
            borderColor: annotationLine.lineColor,
            borderWidth: 2,
            label: {
              enabled: true,
              position: 'left',
              xPadding: 5,
              yPadding: 5,
              content: annotationLine.label,
              backgroundColor: annotationLine.labelColor,
              font: {
                size: 10,
              },
            },
          },
        },
      };
    },

    options() {
      return {
        annotation: this.annotationLine,
        tooltips: {
          callbacks: {
            label: this.tooltipLabel,
          },
        },
      };
    },
  },
  methods: {
    getQuery() {
      const {
        mfilter,
        tstart,
        tstop,
        periodUnit,
        periodValue,
        duration,
        stats = {},
      } = this.getStatsQuery();

      return {
        duration,
        mfilter,

        periods: Math.ceil((tstop.diff(tstart, periodUnit) + 1) / periodValue),
        tstop: tstop.startOf('h').unix(),
        stats: Object.entries(stats).reduce((acc, [key, value]) => {
          acc[key] = {
            ...value,

            aggregate: ['sum'],
          };

          return acc;
        }, {}),
      };
    },

    async fetchList() {
      try {
        this.pending = true;

        const { aggregations } = await this.fetchStatsEvolutionWithoutStore({
          params: this.getQuery(),
        });

        this.stats = aggregations;
        this.pending = false;
      } catch (err) {
        this.serverErrorMessage = err.description || this.$t('errors.statsRequestProblem');
      } finally {
        this.pending = false;
      }
    },

    tooltipLabel(tooltipItem, data) {
      const PROPERTIES_FILTERS_MAP = {
        [STATS_TYPES.stateRate.value]: convertNumberToRoundedPercentString,
        [STATS_TYPES.ackTimeSla.value]: convertNumberToRoundedPercentString,
        [STATS_TYPES.resolveTimeSla.value]: convertNumberToRoundedPercentString,
        [STATS_TYPES.timeInState.value]: convertDurationToString,
        [STATS_TYPES.mtbf.value]: convertDurationToString,
      };

      const { stats } = this.query;

      const statObject = stats ? stats[data.datasets[tooltipItem.datasetIndex].label] : null;
      let label = data.datasets[tooltipItem.datasetIndex].label || '';

      if (label) {
        label += ': ';
      }

      if (statObject && PROPERTIES_FILTERS_MAP[statObject.stat]) {
        label += PROPERTIES_FILTERS_MAP[statObject.stat](tooltipItem.yLabel);
      } else {
        label += tooltipItem.yLabel;
      }

      return label;
    },
  },
};
