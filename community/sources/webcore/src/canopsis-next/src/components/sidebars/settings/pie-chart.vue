<template lang="pug">
  widget-settings(:submitting="submitting", divider, @submit="submit")
    pie-chart-widget-form(
      v-model="form",
      with-periodic-refresh,
      with-filters
    )
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';

import WidgetSettings from './partials/widget-settings.vue';
import PieChartWidgetForm from './forms/widgets/pie-chart-widget-form.vue';

export default {
  name: SIDE_BARS.pieChartSettings,
  components: {
    WidgetSettings,
    PieChartWidgetForm,
  },
  mixins: [widgetSettingsMixin],
};
</script>
