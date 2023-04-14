<template lang="pug">
  widget-settings(:submitting="submitting", divider, @submit="submit")
    line-chart-widget-form(
      v-model="form",
      with-periodic-refresh,
      with-filters
    )
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';

import WidgetSettings from './partials/widget-settings.vue';
import LineChartWidgetForm from './forms/widgets/line-chart-widget-form.vue';

export default {
  name: SIDE_BARS.lineChartSettings,
  components: {
    WidgetSettings,
    LineChartWidgetForm,
  },
  mixins: [widgetSettingsMixin],
};
</script>
