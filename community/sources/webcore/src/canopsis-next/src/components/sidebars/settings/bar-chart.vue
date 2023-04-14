<template lang="pug">
  widget-settings(:submitting="submitting", divider, @submit="submit")
    bar-chart-widget-form(
      v-model="form",
      with-periodic-refresh,
      with-filters
    )
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';

import WidgetSettings from './partials/widget-settings.vue';
import BarChartWidgetForm from './forms/widgets/bar-chart-widget-form.vue';

export default {
  name: SIDE_BARS.barChartSettings,
  components: {
    WidgetSettings,
    BarChartWidgetForm,
  },
  mixins: [widgetSettingsMixin],
};
</script>
