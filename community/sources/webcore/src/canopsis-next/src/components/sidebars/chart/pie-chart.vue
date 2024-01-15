<template>
  <widget-settings
    :submitting="submitting"
    divider
    @submit="submit"
  >
    <pie-chart-widget-form
      v-model="form"
      :with-filters="hasAccessToListFilters"
      with-periodic-refresh
      with-preset
    />
  </widget-settings>
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { permissionsWidgetsLineChartFilters } from '@/mixins/permissions/widgets/chart/line/filters';

import WidgetSettings from '../partials/widget-settings.vue';

import PieChartWidgetForm from './form/pie-chart-widget-form.vue';

export default {
  name: SIDE_BARS.pieChartSettings,
  components: {
    WidgetSettings,
    PieChartWidgetForm,
  },
  mixins: [
    widgetSettingsMixin,
    permissionsWidgetsLineChartFilters,
  ],
};
</script>
