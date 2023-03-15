import { DATETIME_FORMATS, SAMPLINGS } from '@/constants';

import {
  convertStartDateIntervalToTimestampByTimezone,
  convertStopDateIntervalToTimestampByTimezone,
} from '@/helpers/date/date-intervals';
import { convertDateToStartOfDayTimestampByTimezone } from '@/helpers/date/date';
import { convertMetricsToTimezone } from '@/helpers/metrics';

import { entitiesMetricsMixin } from '@/mixins/entities/metrics';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';

export const widgetFetchMetricsMixin = {
  inject: ['$system'],
  mixins: [widgetFetchQueryMixin, entitiesMetricsMixin],
  data() {
    return {
      alarmsMetrics: [],
      minAvailableDate: null,
    };
  },
  mounted() {
    this.setQuery();
  },
  methods: {
    getIntervalQuery() {
      if (!this.query.interval) {
        return {};
      }

      return {
        from: convertStartDateIntervalToTimestampByTimezone(
          this.query.interval.from,
          DATETIME_FORMATS.datePicker,
          SAMPLINGS.day,
          this.$system.timezone,
        ),
        to: convertStopDateIntervalToTimestampByTimezone(
          this.query.interval.to,
          DATETIME_FORMATS.datePicker,
          SAMPLINGS.day,
          this.$system.timezone,
        ),
      };
    },

    getVectorQuery() {
      return {
        ...this.getIntervalQuery(),

        parameters: this.widget.parameters.metrics.map(({ metric }) => metric),
        sampling: this.query.sampling,
        filter: this.query.filter,
      };
    },

    getAggregatedQuery() {
      return {
        ...this.getIntervalQuery(),

        parameters: this.widget.parameters.metrics.map(({ metric, aggregate_func: aggregateFunc }) => ({
          metric,
          aggregate_func: aggregateFunc ?? this.widget.parameters.aggregate_func,
        })),
        sampling: this.query.sampling,
        filter: this.query.filter,
      };
    },

    async fetchVectorMetrics() {
      try {
        this.pending = true;

        const {
          data: alarmsMetrics,
          meta: { min_date: minDate },
        } = await this.fetchAlarmsMetricsWithoutStore({ params: this.getVectorQuery() });

        this.alarmsMetrics = convertMetricsToTimezone(alarmsMetrics, this.$system.timezone);
        this.minAvailableDate = convertDateToStartOfDayTimestampByTimezone(minDate, this.$system.timezone);
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },

    async fetchAggregatedMetrics() {
      try {
        this.pending = true;

        const { data: alarmsMetrics } = await this.fetchAggregateMetricsWithoutStore({
          params: this.getAggregatedQuery(),
        });

        this.alarmsMetrics = alarmsMetrics;
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
