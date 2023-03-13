<template lang="pug">
  widget-settings-item(:title="$t('settings.chart.sharesType')")
    v-radio-group.mt-0(v-field="value", :name="name", :label="label")
      v-radio(
        v-for="showMode in showModes",
        :key="showMode.value",
        :label="showMode.label",
        :value="showMode.value",
        color="primary"
      )
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

import { KPI_PIE_CHART_SHOW_MODES } from '@/constants';

export default {
  components: { WidgetSettingsItem },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      required: true,
    },
    label: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      default: 'show_mode',
    },
  },
  computed: {
    showModes() {
      return [
        {
          value: KPI_PIE_CHART_SHOW_MODES.numbers,
          label: this.$tc('common.number', 2),
        },
        {
          value: KPI_PIE_CHART_SHOW_MODES.percent,
          label: this.$tc('common.percent', 2),
        },
      ];
    },
  },
};
</script>
