<template>
  <widget-settings-item :title="title">
    <v-range-slider
      :model-value="modelValue"
      :min="min"
      :max="max"
      :step="step"
      color="primary"
      class="mt-3"
      show-ticks="always"
      thumb-label
      hide-details
      @update:model-value="$emit('update:modelValue', $event)"
    />
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';
import { GRID_SIZES } from '@/constants';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
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

  emits: ['update:modelValue'],
};
</script>
