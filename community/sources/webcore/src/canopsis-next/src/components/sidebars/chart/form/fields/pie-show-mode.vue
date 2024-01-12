<template>
  <widget-settings-item :title="$t('settings.chart.sharesType')">
    <v-radio-group
      class="mt-0"
      v-field="value"
      :name="name"
      :label="label"
    >
      <v-radio
        v-for="showMode in showModes"
        :key="showMode.value"
        :label="showMode.label"
        :value="showMode.value"
        color="primary"
      />
    </v-radio-group>
  </widget-settings-item>
</template>

<script>
import { KPI_PIE_CHART_SHOW_MODS } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

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
          value: KPI_PIE_CHART_SHOW_MODS.numbers,
          label: this.$tc('common.number', 2),
        },
        {
          value: KPI_PIE_CHART_SHOW_MODS.percent,
          label: this.$tc('common.percent', 2),
        },
      ];
    },
  },
};
</script>
