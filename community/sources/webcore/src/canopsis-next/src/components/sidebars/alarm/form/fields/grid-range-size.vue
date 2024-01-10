<template>
  <widget-settings-item :title="title">
    <v-range-slider
      v-field="value"
      :min="min"
      :max="max"
      :step="step"
      ticks="always"
      always-dirty
      thumb-label
    />
  </widget-settings-item>
</template>

<script>
import { GRID_SIZES } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: Array,
      default: () => [GRID_SIZES.min, GRID_SIZES.max],
    },
    min: {
      type: Number,
      default: GRID_SIZES.min,
    },
    max: {
      type: Number,
      default: GRID_SIZES.max,
    },
    step: {
      type: Number,
      default: GRID_SIZES.step,
    },
    title: {
      type: String,
      default: '',
    },
  },
};
</script>
