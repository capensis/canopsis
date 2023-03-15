<template lang="pug">
  v-layout.px-3.pb-3(column)
    chart-widget-filters(
      :widget-id="widget._id",
      :user-filters="userPreference.filters",
      :widget-filters="widget.filters",
      :locked-value="lockedFilter",
      :filters="mainFilter",
      :interval="query.interval",
      :min-interval-date="minAvailableDate",
      :sampling="query.sampling",
      :show-filter="hasAccessToUserFilter",
      :show-interval="hasAccessToInterval",
      :show-sampling="hasAccessToSampling",
      :filter-disabled="!hasAccessToListFilters",
      :filter-addable="hasAccessToAddFilter",
      :filter-editable="hasAccessToEditFilter",
      @update:filters="updateSelectedFilter",
      @update:sampling="updateSampling",
      @update:interval="updateInterval"
    )
    v-layout
      pre {{ vectorMetrics }}
</template>

<script>
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { permissionsWidgetsLineChartInterval } from '@/mixins/permissions/widgets/chart/line/interval';
import { permissionsWidgetsLineChartSampling } from '@/mixins/permissions/widgets/chart/line/sampling';
import { permissionsWidgetsLineChartFilters } from '@/mixins/permissions/widgets/chart/line/filters';
import { widgetIntervalFilterMixin } from '@/mixins/widget/chart/interval';
import { widgetSamplingFilterMixin } from '@/mixins/widget/chart/sampling';
import { entitiesVectorMetricsMixin } from '@/mixins/entities/vector-metrics';

import ChartWidgetFilters from '@/components/widgets/chart/partials/chart-widget-filters.vue';
import { convertDateToStartOfDayTimestampByTimezone } from '@/helpers/date/date';

export default {
  inject: ['$system'],
  components: {
    ChartWidgetFilters,
  },
  mixins: [
    widgetFetchQueryMixin,
    widgetFilterSelectMixin,
    widgetIntervalFilterMixin,
    widgetSamplingFilterMixin,
    permissionsWidgetsLineChartInterval,
    permissionsWidgetsLineChartSampling,
    permissionsWidgetsLineChartFilters,
    entitiesVectorMetricsMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    tabId: {
      type: String,
      default: '',
    },
  },
  computed: {
    minAvailableDate() {
      const { min_date: minDate } = this.vectorMetricsMeta;

      return minDate
        ? convertDateToStartOfDayTimestampByTimezone(minDate, this.$system.timezone)
        : null;
    },
  },
  watch: {
    minAvailableDate() {
      const { from } = this.getIntervalQuery();

      if (from < this.minAvailableDate) {
        this.query = {
          ...this.query,
          interval: {
            from: this.minAvailableDate,
            to: this.query.interval.to,
          },
        };
      }
    },
  },
  methods: {
    getQuery() {
      return {
        ...this.getIntervalQuery(),

        parameters: this.widget.parameters.metrics.map(({ metric }) => metric),
        sampling: this.query.sampling,
        filter: this.query.filter,
      };
    },

    fetchList() {
      this.fetchVectorMetricsList({
        widgetId: this.widget._id,
        params: this.getQuery(),
      });
    },
  },
};
</script>
