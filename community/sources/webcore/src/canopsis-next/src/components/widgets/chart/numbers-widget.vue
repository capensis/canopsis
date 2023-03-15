<template lang="pug">
  v-layout.px-3.pb-3(column)
    chart-widget-filters(
      :widget-id="widget._id",
      :user-filters="userPreference.filters",
      :widget-filters="widget.filters",
      :locked-value="lockedFilter",
      :filters="mainFilter",
      :interval="query.interval",
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
      pre {{ aggregatedMetrics }}
</template>

<script>
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { permissionsWidgetsNumbersInterval } from '@/mixins/permissions/widgets/chart/numbers/interval';
import { permissionsWidgetsNumbersSampling } from '@/mixins/permissions/widgets/chart/numbers/sampling';
import { permissionsWidgetsNumbersFilters } from '@/mixins/permissions/widgets/chart/numbers/filters';
import { widgetIntervalFilterMixin } from '@/mixins/widget/chart/interval';
import { widgetSamplingFilterMixin } from '@/mixins/widget/chart/sampling';
import { entitiesAggregatedMetricsMixin } from '@/mixins/entities/aggregated-metrics';

import ChartWidgetFilters from '@/components/widgets/chart/partials/chart-widget-filters.vue';

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
    permissionsWidgetsNumbersInterval,
    permissionsWidgetsNumbersSampling,
    permissionsWidgetsNumbersFilters,
    entitiesAggregatedMetricsMixin,
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
  methods: {
    getQuery() {
      return {
        ...this.getIntervalQuery(),

        parameters: this.widget.parameters.metrics.map(({ metric, aggregate_func: aggregateFunc }) => ({
          metric,
          aggregate_func: aggregateFunc,
        })),
        sampling: this.query.sampling,
        filter: this.query.filter,
      };
    },

    fetchList() {
      this.fetchAggregatedMetricsList({ widgetId: this.widget._id, params: this.getQuery() });
    },
  },
};
</script>
