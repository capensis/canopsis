<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile
        div {{ title }}
    v-container
      v-range-slider(
        v-field="value",
        :min="min",
        :max="max",
        :step="step",
        class="mt-3",
        ticks="always",
        always-dirty,
        thumb-label,
        hide-details
      )
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
