<template>
  <widget-settings
    :submitting="submitting"
    divider
    @submit="submit"
  >
    <line-chart-widget-form
      v-model="form"
      :with-filters="hasAccessToListFilters"
      with-periodic-refresh
    />
  </widget-settings>
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { permissionsWidgetsLineChartFilters } from '@/mixins/permissions/widgets/chart/line/filters';

import WidgetSettings from '../partials/widget-settings.vue';

import LineChartWidgetForm from './form/line-chart-widget-form.vue';

export default {
  name: SIDE_BARS.lineChartSettings,
  components: {
    WidgetSettings,
    LineChartWidgetForm,
  },
  mixins: [
    widgetSettingsMixin,
    permissionsWidgetsLineChartFilters,
  ],
};
</script>
