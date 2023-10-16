<template>
  <widget-settings
    :submitting="submitting"
    divider
    @submit="submit"
  >
    <statistics-form
      v-model="form"
      :widget="widget"
      :type="$constants.KPI_RATING_SETTINGS_TYPES.entity"
      :show-filter="hasAccessToListFilters"
      :show-interval="hasAccessToInterval"
      :filter-disabled="!hasAccessToListFilters"
      :filter-addable="hasAccessToAddFilter"
      :filter-editable="hasAccessToEditFilter"
    />
  </widget-settings>
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { permissionsWidgetsAlarmStatisticsInterval } from '@/mixins/permissions/widgets/statistics/alarm/interval';
import { permissionsWidgetsAlarmStatisticsFilters } from '@/mixins/permissions/widgets/statistics/alarm/filters';

import WidgetSettings from '../partials/widget-settings.vue';

import StatisticsForm from './form/statistics-form.vue';

export default {
  name: SIDE_BARS.alarmStatisticsSettings,
  components: {
    WidgetSettings,
    StatisticsForm,
  },
  mixins: [
    widgetSettingsMixin,
    permissionsWidgetsAlarmStatisticsInterval,
    permissionsWidgetsAlarmStatisticsFilters,
  ],
};
</script>
