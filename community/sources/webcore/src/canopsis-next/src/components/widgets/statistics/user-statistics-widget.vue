<template>
  <statistics-widget
    :widget="widget"
    :tab-id="tabId"
    :type="$constants.KPI_RATING_SETTINGS_TYPES.user"
    :show-filter="hasAccessToUserFilter"
    :show-interval="hasAccessToInterval"
  />
</template>

<script>
import { permissionsWidgetsUserStatisticsFilters } from '@/mixins/permissions/widgets/statistics/user/filters';
import { permissionsWidgetsUserStatisticsInterval } from '@/mixins/permissions/widgets/statistics/user/interval';

import StatisticsWidget from './partials/statistics-widget.vue';

export default {
  components: { StatisticsWidget },

  mixins: [permissionsWidgetsUserStatisticsInterval, permissionsWidgetsUserStatisticsFilters],

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
