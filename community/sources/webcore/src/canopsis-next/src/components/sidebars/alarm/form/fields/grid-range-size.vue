<template>
  <v-list-group>
    <template #activator="">
      <v-list-item>
        <div>{{ title }}</div>
      </v-list-item>
    </template>
    <v-container>
      <v-range-slider
        v-field="value"
        :min="min"
        :max="max"
        :step="step"
        ticks="always"
        always-dirty
        thumb-label
      />
    </v-container>
  </v-list-group>
</template>

<script>
import { GRID_SIZES } from '@/constants';

export default {
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
