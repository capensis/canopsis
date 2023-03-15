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
      pre {{ alarmsMetrics }}
</template>

<script>
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { permissionsWidgetsBarChartInterval } from '@/mixins/permissions/widgets/chart/bar/interval';
import { permissionsWidgetsBarChartSampling } from '@/mixins/permissions/widgets/chart/bar/sampling';
import { permissionsWidgetsBarChartFilters } from '@/mixins/permissions/widgets/chart/bar/filters';
import { widgetIntervalFilterMixin } from '@/mixins/widget/chart/interval';
import { widgetSamplingFilterMixin } from '@/mixins/widget/chart/sampling';
import { widgetFetchMetricsMixin } from '@/mixins/widget/chart/fetch-metrics';

import ChartWidgetFilters from '@/components/widgets/chart/partials/chart-widget-filters.vue';

export default {
  components: {
    ChartWidgetFilters,
  },
  mixins: [
    widgetFilterSelectMixin,
    widgetIntervalFilterMixin,
    widgetSamplingFilterMixin,
    widgetFetchMetricsMixin,
    permissionsWidgetsBarChartInterval,
    permissionsWidgetsBarChartSampling,
    permissionsWidgetsBarChartFilters,
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
    fetchList() {
      this.fetchVectorMetrics();
    },
  },
};
</script>
