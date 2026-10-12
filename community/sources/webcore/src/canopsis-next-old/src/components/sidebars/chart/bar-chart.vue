<template>
  <widget-settings
    :submitting="submitting"
    :dirty="hasChanges"
    divider
    @submit="submit"
  >
    <bar-chart-widget-form
      v-model="form"
      :with-filters="hasAccessToFilter"
      with-periodic-refresh
      with-preset
    />
  </widget-settings>
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { permissionsWidgetsBarChartFilters } from '@/mixins/permissions/widgets/chart/bar/filters';

import WidgetSettings from '../partials/widget-settings.vue';

import BarChartWidgetForm from './form/bar-chart-widget-form.vue';

export default {
  name: SIDE_BARS.barChartSettings,
  components: {
    WidgetSettings,
    BarChartWidgetForm,
  },
  mixins: [
    widgetSettingsMixin,
    permissionsWidgetsBarChartFilters,
  ],
};
</script>
