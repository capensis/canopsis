<template>
  <statistics-widget
    :widget="widget"
    :tab-id="tabId"
    :type="$constants.KPI_RATING_SETTINGS_TYPES.entity"
    :show-filter="hasAccessToUserFilter"
    :show-interval="hasAccessToInterval"
  />
</template>

<script>
import { permissionsWidgetsAlarmStatisticsFilters } from '@/mixins/permissions/widgets/statistics/alarm/filters';
import { permissionsWidgetsAlarmStatisticsInterval } from '@/mixins/permissions/widgets/statistics/alarm/interval';

import StatisticsWidget from './partials/statistics-widget.vue';

export default {
  components: { StatisticsWidget },

  mixins: [permissionsWidgetsAlarmStatisticsInterval, permissionsWidgetsAlarmStatisticsFilters],

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
