<template>
  <v-list-item v-if="inline" class="overflow-visible">
    <v-list-item-content class="overflow-visible">
      <v-list-item-title>{{ title }}</v-list-item-title>
      <v-list-item-action class="ml-0">
        <v-range-slider
          v-field="value"
          :min="min"
          :max="max"
          :step="step"
          class="mt-3"
          ticks="always"
          always-dirty
          thumb-label
          hide-details
        />
      </v-list-item-action>
    </v-list-item-content>
  </v-list-item>
  <widget-settings-item v-else :title="title">
    <v-range-slider
      v-field="value"
      :min="min"
      :max="max"
      :step="step"
      class="mt-3"
      ticks="always"
      always-dirty
      thumb-label
      hide-details
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
    inline: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss">

</style>
