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
</template>

<script>
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { permissionsWidgetsBarChartInterval } from '@/mixins/permissions/widgets/chart/bar/interval';
import { permissionsWidgetsBarChartSampling } from '@/mixins/permissions/widgets/chart/bar/sampling';
import { permissionsWidgetsBarChartFilters } from '@/mixins/permissions/widgets/chart/bar/filters';
import { widgetIntervalFilterMixin } from '@/mixins/widget/interval';
import { widgetSamplingFilterMixin } from '@/mixins/widget/sampling';

import ChartWidgetFilters from '@/components/widgets/chart/partials/chart-widget-filters.vue';

export default {
  inject: ['$system'],
  components: {
    ChartWidgetFilters,
  },
  mixins: [
    widgetFilterSelectMixin,
    widgetFetchQueryMixin,
    widgetIntervalFilterMixin,
    widgetSamplingFilterMixin,
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
    getQuery() {
      return {
        ...this.getIntervalQuery(),

        sampling: this.query.sampling,
        filter: this.query.filter,
      };
    },

    fetchList() {},
  },
};
</script>
