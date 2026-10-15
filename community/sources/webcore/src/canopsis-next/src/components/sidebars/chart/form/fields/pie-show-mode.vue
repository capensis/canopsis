<template>
  <widget-settings-item :title="$t('settings.chart.sharesType')">
    <c-radio-group
      :model-value="modelValue"
      :name="name"
      :label="label"
      class="mt-0"
      @update:model-value="$emit('update:modelValue', $event)"
    >
      <v-radio
        v-for="showMode in showModes"
        :key="showMode.value"
        :label="showMode.label"
        :value="showMode.value"
        color="primary"
      />
    </c-radio-group>
  </widget-settings-item>
</template>

<script>
import { KPI_PIE_CHART_SHOW_MODS } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
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

  emits: ['update:modelValue'],

  computed: {
    showModes() {
      return [
        {
          value: KPI_PIE_CHART_SHOW_MODS.numbers,
          label: this.$t('common.number', 2),
        },
        {
          value: KPI_PIE_CHART_SHOW_MODS.percent,
          label: this.$t('common.percent', 2),
        },
      ];
    },
  },
};
</script>
