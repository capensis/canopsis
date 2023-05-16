<template lang="pug">
  widget-settings(:submitting="submitting", divider, @submit="submit")
    statistics-form(
      v-model="form",
      :widget="widget",
      :type="$constants.KPI_RATING_SETTINGS_TYPES.entity"
    )
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';

import WidgetSettings from '../partials/widget-settings.vue';

import StatisticsForm from './form/statistics-form.vue';

export default {
  name: SIDE_BARS.alarmStatisticsSettings,
  components: {
    WidgetSettings,
    StatisticsForm,
  },
  mixins: [widgetSettingsMixin],
};
</script>
