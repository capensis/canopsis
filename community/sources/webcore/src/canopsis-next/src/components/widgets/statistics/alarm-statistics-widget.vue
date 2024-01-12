<template>
  <statistics-widget
    :widget="widget"
    :tab-id="tabId"
    :show-filter="hasAccessToUserFilter"
    :show-interval="hasAccessToInterval"
    :filter-disabled="!hasAccessToListFilters"
    :filter-addable="hasAccessToAddFilter"
    :filter-editable="hasAccessToEditFilter"
  />
</template>

<script>
import { permissionsWidgetsAlarmStatisticsInterval } from '@/mixins/permissions/widgets/statistics/alarm/interval';
import { permissionsWidgetsAlarmStatisticsFilters } from '@/mixins/permissions/widgets/statistics/alarm/filters';

import StatisticsWidget from './partials/statistics-widget.vue';

export default {
  components: { StatisticsWidget },
  mixins: [
    permissionsWidgetsAlarmStatisticsInterval,
    permissionsWidgetsAlarmStatisticsFilters,
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
};
</script>
