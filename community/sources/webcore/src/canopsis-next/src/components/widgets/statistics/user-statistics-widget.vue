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
import { permissionsWidgetsUserStatisticsInterval } from '@/mixins/permissions/widgets/statistics/user/interval';
import { permissionsWidgetsUserStatisticsFilters } from '@/mixins/permissions/widgets/statistics/user/filters';

import StatisticsWidget from './partials/statistics-widget.vue';

export default {
  components: { StatisticsWidget },
  mixins: [
    permissionsWidgetsUserStatisticsInterval,
    permissionsWidgetsUserStatisticsFilters,
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
